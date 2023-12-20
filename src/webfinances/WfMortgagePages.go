package webfinances

import (
	"context"
	"finance/finances"
	"finance/middlewares"
	"finance/sessions"
	"fmt"
	"html/template"
	"net/http"
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
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering MortgagePages/webfinances.",
  })
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    mf := getMortgageFields(sessions.GetUserName(sessionToken))
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
      mf.currentPage = ui
    }
    //
    if strings.EqualFold(mf.currentPage, "rhs-ui1") {
      mf.currentButton = "lhs-button1"
      if req.Method == http.MethodPost {
        mf.fd1N = req.PostFormValue("fd1-n")
        mf.fd1TimePeriod = req.PostFormValue("fd1-tp")
        mf.fd1Interest = req.PostFormValue("fd1-interest")
        mf.fd1Compound = req.PostFormValue("fd1-compound")
        mf.fd1Amount = req.PostFormValue("fd1-amount")
        var n float64
        var i float64
        var amount float64
        var err error
        mf.fd1Result[1] = ""
        mf.fd1Result[2] = ""
        if n, err = strconv.ParseFloat(mf.fd1N, 64); err != nil {
          mf.fd1Result[0] = fmt.Sprintf("Error: %s -- %+v", mf.fd1N, err)
        } else if i, err = strconv.ParseFloat(mf.fd1Interest, 64); err != nil {
          mf.fd1Result[0] = fmt.Sprintf("Error: %s -- %+v", mf.fd1Interest, err)
        } else if amount, err = strconv.ParseFloat(mf.fd1Amount, 64); err != nil {
          mf.fd1Result[0] = fmt.Sprintf("Error: %s -- %+v", mf.fd1Amount, err)
        } else {
          var m finances.Mortgage
          payment, totalCost, totalInterest := m.CostOfMortgage(amount, i / 100.0,
            mf.fd1Compound[0], n, mf.fd1TimePeriod[0])
          mf.fd1Result[0] = fmt.Sprintf("Payment: $%.2f", payment)
          mf.fd1Result[1] = fmt.Sprintf("Total Interest: $%.2f", totalInterest)
          mf.fd1Result[2] = fmt.Sprintf("Total Cost: $%.2f", totalCost)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, interest = %s, cp = %s, amount = %s, %s", mf.fd1N,
            mf.fd1TimePeriod, mf.fd1Interest, mf.fd1Compound, mf.fd1Amount, mf.fd1Result[0]),
        })
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
      } { "Mortgage", m.DTF(), mf.currentButton, newSession.CsrfToken,
          mf.fd1N, mf.fd1TimePeriod, mf.fd1Interest, mf.fd1Compound, mf.fd1Amount, mf.fd1Result,
        })
    } else if strings.EqualFold(mf.currentPage, "rhs-ui2") {
      mf.currentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        mf.fd2N = req.FormValue("fd2-n")
        mf.fd2TimePeriod = req.PostFormValue("fd2-tp")
        mf.fd2Interest = req.PostFormValue("fd2-i")
        mf.fd2Compound = req.PostFormValue("fd2-compound")
        mf.fd2Amount = req.PostFormValue("fd2-amount")
        var n float64
        var i float64
        var amount float64
        var err error
        if n, err = strconv.ParseFloat(mf.fd2N, 64); err != nil {
          mf.fd2Result = append(mf.fd2Result,
            Row {
              PaymentNo: fmt.Sprintf("Error: %s -- %+v", mf.fd2N, err),
            })
        } else if i, err = strconv.ParseFloat(mf.fd2Interest, 64); err != nil {
          mf.fd2Result = append(mf.fd2Result,
            Row {
              PaymentNo: fmt.Sprintf("Error: %s -- %+v", mf.fd2Interest, err),
            })
        } else if amount, err = strconv.ParseFloat(mf.fd2Amount, 64); err != nil {
          mf.fd2Result = append(mf.fd2Result,
            Row {
              PaymentNo: fmt.Sprintf("Error: %s -- %+v", mf.fd2Amount, err),
            })
        } else {
          var m finances.Mortgage
          var at = m.AmortizationTable(amount, i / 100.0, mf.fd1Compound[0], n,
            mf.fd1TimePeriod[0])
          var numberOfRows = len(at.Rows)
          mf.fd2Result = make([]Row, 0, numberOfRows + 1)
          mf.fd2Result = append(mf.fd2Result,
            Row {
              PaymentNo: "--",
              Payment: "--",
              PmtPrincipal: "--",
              PmtInterest: "--",
              Balance: fmt.Sprintf("%.2f", amount),
            })
          for idx := 0; idx < numberOfRows; idx++ {
            mf.fd2Result = append(mf.fd2Result,
              Row {
                PaymentNo: fmt.Sprintf("%v", idx + 1),
                Payment: fmt.Sprintf("%.2f", at.Rows[idx].Payment),
                PmtPrincipal: fmt.Sprintf("%.2f", at.Rows[idx].PmtPrincipal),
                PmtInterest: fmt.Sprintf("%.2f", at.Rows[idx].PmtInterest),
                Balance: fmt.Sprintf("%.2f", at.Rows[idx].Balance),
              })
          }
          mf.fd2TotalCost = fmt.Sprintf("Total Cost: $%.2f", at.TotalCost)
          mf.fd2TotalInterest = fmt.Sprintf("Total Interest: $%.2f", at.TotalInterest)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, interest = %s, cp = %s, amount = %s, total cost = %s, total interest = %s",
            mf.fd2N, mf.fd2TimePeriod, mf.fd2Interest, mf.fd2Compound, mf.fd2Amount,
            mf.fd2TotalCost, mf.fd2TotalInterest),
        })
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
      } { "Mortgage", m.DTF(), mf.currentButton, newSession.CsrfToken, mf.fd2N, mf.fd2TimePeriod,
          mf.fd2Interest, mf.fd2Compound, mf.fd2Amount, mf.fd2TotalCost, mf.fd2TotalInterest,
          mf.fd2Result,
        })
    } else if strings.EqualFold(mf.currentPage, "rhs-ui3") {
      mf.currentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        mf.fd3Mrate = req.PostFormValue("fd3-mrate")
        mf.fd3Mbalance = req.PostFormValue("fd3-mbalance")
        mf.fd3Hrate = req.PostFormValue("fd3-hrate")
        mf.fd3Hbalance = req.PostFormValue("fd3-hbalance")
        var mRate float64
        var mBalance float64
        var hRate float64
        var hBalance float64
        var err error
        if mRate, err = strconv.ParseFloat(mf.fd3Mrate, 64); err != nil {
          mf.fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", mf.fd3Mrate, err)
        } else if mBalance, err = strconv.ParseFloat(mf.fd3Mbalance, 64); err != nil {
          mf.fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", mf.fd3Mbalance, err)
        } else if hRate, err = strconv.ParseFloat(mf.fd3Hrate, 64); err != nil {
          mf.fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", mf.fd3Hrate, err)
        } else if hBalance, err = strconv.ParseFloat(mf.fd3Hbalance, 64); err != nil {
          mf.fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", mf.fd3Hbalance, err)
        } else {
          var m finances.Mortgage
          mf.fd3Result[2] = fmt.Sprintf("Blended Interest Rate: %.3f%%",
            m.BlendedInterestRate(mBalance, mRate, hBalance, hRate))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("mortgage balance = %s, mortgage rate = %s, HELOC balance = %s, HELOC rate = %s, %s\n",
            mf.fd3Mbalance, mf.fd3Mrate, mf.fd3Hbalance, mf.fd3Hrate, mf.fd3Result[2]),
        })
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
      } { "Mortgage", m.DTF(), mf.currentButton, newSession.CsrfToken, mf.fd3Mrate, mf.fd3Mbalance,
          mf.fd3Hrate, mf.fd3Hbalance, mf.fd3Result,
        })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", mf.currentPage)
      fmt.Printf("%s - %s\n", m.DTF(), errString)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      fmt.Println("*** Request timeout ***")
      if strings.EqualFold(mf.currentPage, "rhs-ui1") {
        mf.fd1Result[0] = ""
        mf.fd1Result[1] = ""
        mf.fd1Result[2] = ""
      } else if strings.EqualFold(mf.currentPage, "rhs-ui2") {
        mf.fd2Result = nil
        mf.fd2TotalCost = ""
        mf.fd2TotalInterest = ""
      } else if strings.EqualFold(mf.currentPage, "rhs-ui3") {
        mf.fd3Result[2] = ""
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
