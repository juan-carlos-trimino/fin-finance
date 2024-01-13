package webfinances

import (
  "context"
	"encoding/json"
  "finance/finances"
  "finance/middlewares"
	"finance/misc"
	"finance/sessions"
  "fmt"
  "html/template"
  "net/http"
	"os"
  "strconv"
  "strings"
)

var misc_notes = [...]string {
  "When comparing interest rates, use effective annual rates.",
  "Nominal returns are not adjusted for inflation.",
  "Real returns are useful while comparing returns over different time periods because of the differences in inflation rates.",
  "Real returns are adjusted for inflation.",
  "Values are semicolon (;) separated; e.g., 3;3.1;3.2;-1.01",
}

type WfMiscellaneousPages struct {
}

func (mp WfMiscellaneousPages) MiscellaneousPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
    return
  }
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering MiscellaneousPages/webfinances.",
  })
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    mf := getMiscellaneousFields(userName)
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
        mf.Fd1Nominal = req.PostFormValue("fd1-nominal")
        mf.Fd1Compound = req.PostFormValue("fd1-compound")
        var nr float64
        var err error
        if nr, err = strconv.ParseFloat(mf.Fd1Nominal, 64); err != nil {
          mf.Fd1Result[1] = fmt.Sprintf("Error: %s -- %+v", mf.Fd1Nominal, err)
        } else {
          var a finances.Annuities
          mf.Fd1Result[1] = fmt.Sprintf("Effective Annual Rate: %.3f%%",
           a.NominalRateToEAR(nr / 100.0, a.GetCompoundingPeriod(mf.Fd1Compound[0], false)) * 100.0)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("nominal rate = %s, cp = %s, %s", mf.Fd1Nominal, mf.Fd1Compound,
           mf.Fd1Result[1]),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      t := template.Must(template.ParseFiles("webfinances/templates/miscellaneous/miscellaneous.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/miscellaneous/nominalrate.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "miscellaneous", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd1Nominal string
        Fd1Compound string
        Fd1Result [2]string
      } { "Miscellaneous", m.DTF(), mf.CurrentButton, newSession.CsrfToken,
          mf.Fd1Nominal, mf.Fd1Compound, mf.Fd1Result,
        })
    } else if strings.EqualFold(mf.CurrentPage, "rhs-ui2") {
      mf.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        mf.Fd2Effective = req.PostFormValue("fd2-effective")
        mf.Fd2Compound = req.PostFormValue("fd2-compound")
        var ear float64
        var err error
        if ear, err = strconv.ParseFloat(mf.Fd2Effective, 64); err != nil {
          mf.Fd2Result[2] = fmt.Sprintf("Error: %s -- %+v", mf.Fd2Effective, err)
        } else {
          var a finances.Annuities
          mf.Fd2Result[2] = fmt.Sprintf("Nominal Rate: %.3f%% %s", a.EARToNominalRate(ear / 100.0,
            a.GetCompoundingPeriod(mf.Fd2Compound[0], false)) * 100.0, mf.Fd2Compound)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("effective rate = %s, cp = %s, %s", mf.Fd2Effective, mf.Fd2Compound,
           mf.Fd2Result[2]),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/miscellaneous/miscellaneous.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/miscellaneous/effectiveannualrate.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "miscellaneous", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd2Effective string
        Fd2Compound string
        Fd2Result [3]string
      } { "Miscellaneous", m.DTF(), mf.CurrentButton, newSession.CsrfToken,
          mf.Fd2Effective, mf.Fd2Compound, mf.Fd2Result,
        })
    } else if strings.EqualFold(mf.CurrentPage, "rhs-ui3") {
      mf.CurrentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        mf.Fd3Nominal = req.PostFormValue("fd3-nominal")
        mf.Fd3Inflation = req.PostFormValue("fd3-inflation")
        var nr float64
        var ir float64
        var err error
        if nr, err = strconv.ParseFloat(mf.Fd3Nominal, 64); err != nil {
          mf.Fd3Result[3] = fmt.Sprintf("Error: %s -- %+v", mf.Fd3Nominal, err)
        } else if ir, err = strconv.ParseFloat(mf.Fd3Inflation, 64); err != nil {
          mf.Fd3Result[3] = fmt.Sprintf("Error: %s -- %+v", mf.Fd3Inflation, err)
        } else {
          var a finances.Annuities
          mf.Fd3Result[3] = fmt.Sprintf("Real Interest Rate: %.3f%%", a.RealInterestRate(
           nr / 100.0, ir / 100.0) * 100.0)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("nominal rate = %s, inflation rate = %s, %s", mf.Fd3Nominal, mf.Fd3Inflation,
           mf.Fd3Result[3]),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/miscellaneous/miscellaneous.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/miscellaneous/nominalratevs.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "miscellaneous", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd3Nominal string
        Fd3Inflation string
        Fd3Result [4]string
      } { "Miscellaneous", m.DTF(), mf.CurrentButton, newSession.CsrfToken,
          mf.Fd3Nominal, mf.Fd3Inflation, mf.Fd3Result,
        })
    } else if strings.EqualFold(mf.CurrentPage, "rhs-ui4") {
      mf.CurrentButton = "lhs-button4"
      if req.Method == http.MethodPost {
        mf.Fd4Interest = req.PostFormValue("fd4-interest")
        mf.Fd4Compound = req.PostFormValue("fd4-compound")
        mf.Fd4Factor = req.PostFormValue("fd4-factor")
        var ir float64
        var factor float64
        var err error
        if ir, err = strconv.ParseFloat(mf.Fd4Interest, 64); err != nil {
          mf.Fd4Result = fmt.Sprintf("Error: %s -- %+v", mf.Fd4Interest, err)
        } else if factor, err = strconv.ParseFloat(mf.Fd4Factor, 64); err != nil {
          mf.Fd4Result = fmt.Sprintf("Error: %s -- %+v", mf.Fd4Factor, err)
        } else {
          var a finances.Annuities
          mf.Fd4Result = fmt.Sprintf("Growth/Decay: %.3f %s", a.GrowthDecayOfFunds(factor,
           ir / 100.0, a.GetCompoundingPeriod(mf.Fd4Compound[0], true)),
           a.TimePeriods(mf.Fd4Compound))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("interest rate = %s, cp = %s, factor = %s, %s\n", mf.Fd4Interest,
           mf.Fd4Compound, mf.Fd4Factor, mf.Fd4Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/miscellaneous/miscellaneous.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/miscellaneous/growthdecay.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "miscellaneous", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd4Interest string
        Fd4Compound string
        Fd4Factor string
        Fd4Result string
      } { "Miscellaneous", m.DTF(), mf.CurrentButton, newSession.CsrfToken,
          mf.Fd4Interest, mf.Fd4Compound, mf.Fd4Factor, mf.Fd4Result,
        })
    } else if strings.EqualFold(mf.CurrentPage, "rhs-ui5") {
      mf.CurrentButton = "lhs-button5"
      if req.Method == http.MethodPost {
        mf.Fd5Values = req.PostFormValue("fd5-values")
        split := strings.Split(mf.Fd5Values, ";")
        values := make([]float64, len(split))
        var err error
        for i, s := range split {
          if values[i], err = strconv.ParseFloat(s, 64); err != nil {
            mf.Fd5Result[1] = fmt.Sprintf("Error: %s -- %+v", s, err)
            break;
          }
        }
        //
        if err == nil {
          var a finances.Annuities
          mf.Fd5Result[1] = fmt.Sprintf("Avg: %.3f%%", a.AverageRateOfReturn(values) * 100.0)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("values = [%s], %s\n", mf.Fd5Values, mf.Fd5Result[1]),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/miscellaneous/miscellaneous.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/miscellaneous/averagerate.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "miscellaneous", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd5Values string
        Fd5Result [2]string
      } { "Miscellaneous", m.DTF(), mf.CurrentButton, newSession.CsrfToken,
          mf.Fd5Values, mf.Fd5Result,
        })
    } else if strings.EqualFold(mf.CurrentPage, "rhs-ui6") {
      mf.CurrentButton = "lhs-button6"
      if req.Method == http.MethodPost {
        mf.Fd6Time = req.PostFormValue("fd6-time")
        mf.Fd6TimePeriod = req.PostFormValue("fd6-tp")
        mf.Fd6Rate = req.PostFormValue("fd6-rate")
        mf.Fd6Compound = req.PostFormValue("fd6-compound")
        mf.Fd6PV = req.PostFormValue("fd6-pv")
        var time float64
        var rate float64
        var pv float64
        var err error
        if time, err = strconv.ParseFloat(mf.Fd6Time, 64); err != nil {
          mf.Fd6Result = fmt.Sprintf("Error: %s -- %+v", mf.Fd6Time, err)
        } else if rate, err = strconv.ParseFloat(mf.Fd6Rate, 64); err != nil {
          mf.Fd6Result = fmt.Sprintf("Error: %s -- %+v", mf.Fd6Rate, err)
        } else if pv, err = strconv.ParseFloat(mf.Fd6PV, 64); err != nil {
          mf.Fd6Result = fmt.Sprintf("Error: %s -- %+v", mf.Fd6PV, err)
        } else {
          var a finances.Annuities
          mf.Fd6Result = fmt.Sprintf("Future Value: %.2f", a.Depreciation(pv, rate / 100.0,
           a.GetCompoundingPeriod(mf.Fd6Compound[0], false), time,
           a.GetTimePeriod(mf.Fd6TimePeriod[0], false)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("time = %s, tp = %s, rate = %s, cp = %s, pv = %s, %s\n", mf.Fd6Time,
           mf.Fd6TimePeriod, mf.Fd6Rate, mf.Fd6Compound, mf.Fd6PV, mf.Fd6Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/miscellaneous/miscellaneous.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/miscellaneous/depreciation.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "miscellaneous", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd6Time string
        Fd6TimePeriod string
        Fd6Rate string
        Fd6Compound string
        Fd6PV string
        Fd6Result string
      } { "Miscellaneous", m.DTF(), mf.CurrentButton, newSession.CsrfToken,
          mf.Fd6Time, mf.Fd6TimePeriod, mf.Fd6Rate, mf.Fd6Compound, mf.Fd6PV, mf.Fd6Result,
        })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", mf.CurrentPage)
      fmt.Printf("%s - %s\n", m.DTF(), errString)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      fmt.Println("*** Request timeout ***")
      if strings.EqualFold(mf.CurrentPage, "rhs-ui1") {
        mf.Fd1Result[1] = ""
      } else if strings.EqualFold(mf.CurrentPage, "rhs-ui2") {
        mf.Fd2Result[2] = ""
      } else if strings.EqualFold(mf.CurrentPage, "rhs-ui3") {
        mf.Fd3Result[3] = ""
      } else if strings.EqualFold(mf.CurrentPage, "rhs-ui4") {
        mf.Fd4Result = ""
      } else if strings.EqualFold(mf.CurrentPage, "rhs-ui5") {
        mf.Fd5Result[1] = ""
      } else if strings.EqualFold(mf.CurrentPage, "rhs-ui6") {
        mf.Fd6Result = ""
      }
    }
    //
    if data, err := json.Marshal(mf); err != nil {
      fmt.Printf("%s - %s\n", m.DTF(), err)
    } else {
      filePath := fmt.Sprintf("%s/%s/miscellaneous.txt", mainDir, userName)
      if _, err := misc.WriteAllExclusiveLock1(filePath, data, os.O_CREATE | os.O_RDWR, 0o600);
        err != nil {
        fmt.Printf("%s - %s\n", m.DTF(), err)
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
