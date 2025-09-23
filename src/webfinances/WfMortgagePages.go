package webfinances

import (
  "context"
  "encoding/json"
  "finance/finances"
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/gposu"
  "github.com/juan-carlos-trimino/gpsessions"
  "html/template"
  "net/http"
  "os"
  "strconv"
  "strings"
)

var mortgage_notes = [...]string {
  "Refinance mortgage and HELOC with one load.",
  "If the blended interest rate is higher than what you could get on a new fixed-rate mortgage, consider it.",
}

type Row struct { //Rows for the amortization table.
  PaymentNo string
  Payment, PmtPrincipal, PmtInterest, Balance string
}

type WfMortgagePages struct {
}

func (mp WfMortgagePages) MortgagePages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
    return
  }
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logger.LogInfo("Entering MortgagePages/webfinances.", correlationId)
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    mf := getMortgageFields(userName)
    /***
    The functions in Request that allow to extract data from the URL and/or the body revolve around
    the Form, PostForm, and MultipartForm fields; the data are in the form of key-value pairs.

    If the form and the URL have the same key name, both of them will be placed in a slice, with
    the form value always prioritized before the URL value.

    Since we want the form key-value pairs, we can ignore the URL key-value pairs. The PostForm
    field provides key-value pairs only for the form and not the URL. The PostForm field supports
    only application/x-www-form-urlencoded.

    The FormValue method lets you access the key-value pairs just like the Form field, except that
    it's for a specific key and there is no need to call the ParseForm method beforehand -- the
    FormValue method does it. The PostFormValue method does the same thing, except that it's for
    the PostForm field instead of the Form field.
    ***/
    if ui := req.FormValue("compute"); ui != "" {  //Values from form and URL.
      mf.CurrentPage = ui
    }
    //
    if strings.EqualFold(mf.CurrentPage, "rhs-ui1") {
      mf.CurrentButton = "lhs-button1"
      if req.Method == http.MethodPost {
        mf.Fd1N = req.PostFormValue("fd1-n")
        mf.Fd1TimePeriod = req.PostFormValue("fd1-tp")
        mf.Fd1Interest = req.PostFormValue("fd1-interest")
        mf.Fd1Compound = req.PostFormValue("fd1-compound")
        mf.Fd1Amount = req.PostFormValue("fd1-amount")
        var n float64
        var i float64
        var amount float64
        var err error
        mf.Fd1Result[1] = ""
        mf.Fd1Result[2] = ""
        if n, err = strconv.ParseFloat(mf.Fd1N, 64); err != nil {
          mf.Fd1Result[0] = fmt.Sprintf("Error: %s -- %+v", mf.Fd1N, err)
        } else if i, err = strconv.ParseFloat(mf.Fd1Interest, 64); err != nil {
          mf.Fd1Result[0] = fmt.Sprintf("Error: %s -- %+v", mf.Fd1Interest, err)
        } else if amount, err = strconv.ParseFloat(mf.Fd1Amount, 64); err != nil {
          mf.Fd1Result[0] = fmt.Sprintf("Error: %s -- %+v", mf.Fd1Amount, err)
        } else {
          var m finances.Mortgage
          payment, totalCost, totalInterest := m.CostOfMortgage(amount, i / 100.0,
            mf.Fd1Compound[0], n, mf.Fd1TimePeriod[0])
          mf.Fd1Result[0] = fmt.Sprintf("Payment: $%.2f", payment)
          mf.Fd1Result[1] = fmt.Sprintf("Total Interest: $%.2f", totalInterest)
          mf.Fd1Result[2] = fmt.Sprintf("Total Cost: $%.2f", totalCost)
        }
        logger.LogInfo(fmt.Sprintf("n = %s, tp = %s, interest = %s, cp = %s, amount = %s, %s",
         mf.Fd1N, mf.Fd1TimePeriod, mf.Fd1Interest, mf.Fd1Compound, mf.Fd1Amount, mf.Fd1Result[0]),
         correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      t := template.Must(template.ParseFiles("webfinances/templates/mortgage/mortgage.html",
        "webfinances/templates/header.html",
        "webfinances/templates/mortgage/costofmortgage.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "mortgage", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd1N string
        Fd1TimePeriod string
        Fd1Interest string
        Fd1Compound string
        Fd1Amount string
        Fd1Result [3]string
      } { "Mortgage", logger.DatetimeFormat(), mf.CurrentButton, newSession.CsrfToken,
          mf.Fd1N, mf.Fd1TimePeriod, mf.Fd1Interest, mf.Fd1Compound, mf.Fd1Amount, mf.Fd1Result,
        })
    } else if strings.EqualFold(mf.CurrentPage, "rhs-ui2") {
      mf.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        mf.Fd2N = req.FormValue("fd2-n")
        mf.Fd2TimePeriod = req.PostFormValue("fd2-tp")
        mf.Fd2Interest = req.PostFormValue("fd2-i")
        mf.Fd2Compound = req.PostFormValue("fd2-compound")
        mf.Fd2Amount = req.PostFormValue("fd2-amount")
        var n float64
        var i float64
        var amount float64
        var err error
        if n, err = strconv.ParseFloat(mf.Fd2N, 64); err != nil {
          mf.Fd2Result = append(mf.Fd2Result,
            Row {
              PaymentNo: fmt.Sprintf("Error: %s -- %+v", mf.Fd2N, err),
            })
        } else if i, err = strconv.ParseFloat(mf.Fd2Interest, 64); err != nil {
          mf.Fd2Result = append(mf.Fd2Result,
            Row {
              PaymentNo: fmt.Sprintf("Error: %s -- %+v", mf.Fd2Interest, err),
            })
        } else if amount, err = strconv.ParseFloat(mf.Fd2Amount, 64); err != nil {
          mf.Fd2Result = append(mf.Fd2Result,
            Row {
              PaymentNo: fmt.Sprintf("Error: %s -- %+v", mf.Fd2Amount, err),
            })
        } else {
          var m finances.Mortgage
          var at = m.AmortizationTable(amount, i / 100.0, mf.Fd1Compound[0], n,
            mf.Fd1TimePeriod[0])
          var numberOfRows = len(at.Rows)
          mf.Fd2Result = make([]Row, 0, numberOfRows + 1)
          mf.Fd2Result = append(mf.Fd2Result,
            Row {
              PaymentNo: "--",
              Payment: "--",
              PmtPrincipal: "--",
              PmtInterest: "--",
              Balance: fmt.Sprintf("%.2f", amount),
            })
          for idx := 0; idx < numberOfRows; idx++ {
            mf.Fd2Result = append(mf.Fd2Result,
              Row {
                PaymentNo: fmt.Sprintf("%v", idx + 1),
                Payment: fmt.Sprintf("%.2f", at.Rows[idx].Payment),
                PmtPrincipal: fmt.Sprintf("%.2f", at.Rows[idx].PmtPrincipal),
                PmtInterest: fmt.Sprintf("%.2f", at.Rows[idx].PmtInterest),
                Balance: fmt.Sprintf("%.2f", at.Rows[idx].Balance),
              })
          }
          mf.Fd2TotalCost = fmt.Sprintf("Total Cost: $%.2f", at.TotalCost)
          mf.Fd2TotalInterest = fmt.Sprintf("Total Interest: $%.2f", at.TotalInterest)
        }
        logger.LogInfo(fmt.Sprintf(
         "n = %s, tp = %s, interest = %s, cp = %s, amount = %s, total cost = %s, total interest = %s",
         mf.Fd2N, mf.Fd2TimePeriod, mf.Fd2Interest, mf.Fd2Compound, mf.Fd2Amount, mf.Fd2TotalCost,
         mf.Fd2TotalInterest), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/mortgage/mortgage.html",
        "webfinances/templates/header.html",
        "webfinances/templates/mortgage/amortizationtable.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "mortgage", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd2N string
        Fd2TimePeriod string
        Fd2Interest string
        Fd2Compound string
        Fd2Amount string
        Fd2TotalCost string
        Fd2TotalInterest string
        Fd2Result []Row
      } { "Mortgage", logger.DatetimeFormat(), mf.CurrentButton, newSession.CsrfToken, mf.Fd2N,
          mf.Fd2TimePeriod, mf.Fd2Interest, mf.Fd2Compound, mf.Fd2Amount, mf.Fd2TotalCost,
          mf.Fd2TotalInterest, mf.Fd2Result,
        })
    } else if strings.EqualFold(mf.CurrentPage, "rhs-ui3") {
      mf.CurrentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        mf.Fd3Mrate = req.PostFormValue("fd3-mrate")
        mf.Fd3Mbalance = req.PostFormValue("fd3-mbalance")
        mf.Fd3Hrate = req.PostFormValue("fd3-hrate")
        mf.Fd3Hbalance = req.PostFormValue("fd3-hbalance")
        var mRate float64
        var mBalance float64
        var hRate float64
        var hBalance float64
        var err error
        if mRate, err = strconv.ParseFloat(mf.Fd3Mrate, 64); err != nil {
          mf.Fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", mf.Fd3Mrate, err)
        } else if mBalance, err = strconv.ParseFloat(mf.Fd3Mbalance, 64); err != nil {
          mf.Fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", mf.Fd3Mbalance, err)
        } else if hRate, err = strconv.ParseFloat(mf.Fd3Hrate, 64); err != nil {
          mf.Fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", mf.Fd3Hrate, err)
        } else if hBalance, err = strconv.ParseFloat(mf.Fd3Hbalance, 64); err != nil {
          mf.Fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", mf.Fd3Hbalance, err)
        } else {
          var m finances.Mortgage
          mf.Fd3Result[2] = fmt.Sprintf("Blended Interest Rate: %.3f%%",
            m.BlendedInterestRate(mBalance, mRate, hBalance, hRate))
        }
        logger.LogInfo(fmt.Sprintf(
         "mortgage balance = %s, mortgage rate = %s, HELOC balance = %s, HELOC rate = %s, %s",
         mf.Fd3Mbalance, mf.Fd3Mrate, mf.Fd3Hbalance, mf.Fd3Hrate, mf.Fd3Result[2]), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/mortgage/mortgage.html",
        "webfinances/templates/header.html",
        "webfinances/templates/mortgage/heloc.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "mortgage", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd3Mrate string
        Fd3Mbalance string
        Fd3Hrate string
        Fd3Hbalance string
        Fd3Result [3]string
      } { "Mortgage", logger.DatetimeFormat(), mf.CurrentButton, newSession.CsrfToken, mf.Fd3Mrate,
          mf.Fd3Mbalance, mf.Fd3Hrate, mf.Fd3Hbalance, mf.Fd3Result,
        })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", mf.CurrentPage)
      logger.LogError(errString, "-1")
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      logger.LogWarning("*** Request timeout ***", "-1")
      if strings.EqualFold(mf.CurrentPage, "rhs-ui1") {
        mf.Fd1Result[0] = ""
        mf.Fd1Result[1] = ""
        mf.Fd1Result[2] = ""
      } else if strings.EqualFold(mf.CurrentPage, "rhs-ui2") {
        mf.Fd2Result = nil
        mf.Fd2TotalCost = ""
        mf.Fd2TotalInterest = ""
      } else if strings.EqualFold(mf.CurrentPage, "rhs-ui3") {
        mf.Fd3Result[2] = ""
      }
    }
    //
    if data, err := json.Marshal(mf); err != nil {
      logger.LogError(fmt.Sprintf("%+v", err), "-1")
    } else {
      filePath := fmt.Sprintf("%s/%s/mortgage.txt", mainDir, userName)
      if _, err := osu.WriteAllExclusiveLock1(filePath, data, os.O_CREATE | os.O_RDWR |
        os.O_TRUNC, 0o600); err != nil {
        logger.LogError(fmt.Sprintf("%+v", err), "-1")
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    logger.LogError(errString, "-1")
    panic(errString)
  }
}
