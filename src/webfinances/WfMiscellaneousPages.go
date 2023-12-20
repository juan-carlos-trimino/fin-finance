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
    mf := getMiscellaneousFields(sessions.GetUserName(sessionToken))
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
        mf.fd1Nominal = req.PostFormValue("fd1-nominal")
        mf.fd1Compound = req.PostFormValue("fd1-compound")
        var nr float64
        var err error
        if nr, err = strconv.ParseFloat(mf.fd1Nominal, 64); err != nil {
          mf.fd1Result[1] = fmt.Sprintf("Error: %s -- %+v", mf.fd1Nominal, err)
        } else {
          var a finances.Annuities
          mf.fd1Result[1] = fmt.Sprintf("Effective Annual Rate: %.3f%%",
           a.NominalRateToEAR(nr / 100.0, a.GetCompoundingPeriod(mf.fd1Compound[0], false)) * 100.0)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("nominal rate = %s, cp = %s, %s", mf.fd1Nominal, mf.fd1Compound,
           mf.fd1Result[1]),
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
      } { "Miscellaneous", m.DTF(), mf.currentButton, newSession.CsrfToken,
          mf.fd1Nominal, mf.fd1Compound, mf.fd1Result,
        })
    } else if strings.EqualFold(mf.currentPage, "rhs-ui2") {
      mf.currentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        mf.fd2Effective = req.PostFormValue("fd2-effective")
        mf.fd2Compound = req.PostFormValue("fd2-compound")
        var ear float64
        var err error
        if ear, err = strconv.ParseFloat(mf.fd2Effective, 64); err != nil {
          mf.fd2Result[2] = fmt.Sprintf("Error: %s -- %+v", mf.fd2Effective, err)
        } else {
          var a finances.Annuities
          mf.fd2Result[2] = fmt.Sprintf("Nominal Rate: %.3f%% %s", a.EARToNominalRate(ear / 100.0,
            a.GetCompoundingPeriod(mf.fd2Compound[0], false)) * 100.0, mf.fd2Compound)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("effective rate = %s, cp = %s, %s", mf.fd2Effective, mf.fd2Compound,
           mf.fd2Result[2]),
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
      } { "Miscellaneous", m.DTF(), mf.currentButton, newSession.CsrfToken,
          mf.fd2Effective, mf.fd2Compound, mf.fd2Result,
        })
    } else if strings.EqualFold(mf.currentPage, "rhs-ui3") {
      mf.currentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        mf.fd3Nominal = req.PostFormValue("fd3-nominal")
        mf.fd3Inflation = req.PostFormValue("fd3-inflation")
        var nr float64
        var ir float64
        var err error
        if nr, err = strconv.ParseFloat(mf.fd3Nominal, 64); err != nil {
          mf.fd3Result[3] = fmt.Sprintf("Error: %s -- %+v", mf.fd3Nominal, err)
        } else if ir, err = strconv.ParseFloat(mf.fd3Inflation, 64); err != nil {
          mf.fd3Result[3] = fmt.Sprintf("Error: %s -- %+v", mf.fd3Inflation, err)
        } else {
          var a finances.Annuities
          mf.fd3Result[3] = fmt.Sprintf("Real Interest Rate: %.3f%%", a.RealInterestRate(
           nr / 100.0, ir / 100.0) * 100.0)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("nominal rate = %s, inflation rate = %s, %s", mf.fd3Nominal, mf.fd3Inflation,
           mf.fd3Result[3]),
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
      } { "Miscellaneous", m.DTF(), mf.currentButton, newSession.CsrfToken,
          mf.fd3Nominal, mf.fd3Inflation, mf.fd3Result,
        })
    } else if strings.EqualFold(mf.currentPage, "rhs-ui4") {
      mf.currentButton = "lhs-button4"
      if req.Method == http.MethodPost {
        mf.fd4Interest = req.PostFormValue("fd4-interest")
        mf.fd4Compound = req.PostFormValue("fd4-compound")
        mf.fd4Factor = req.PostFormValue("fd4-factor")
        var ir float64
        var factor float64
        var err error
        if ir, err = strconv.ParseFloat(mf.fd4Interest, 64); err != nil {
          mf.fd4Result = fmt.Sprintf("Error: %s -- %+v", mf.fd4Interest, err)
        } else if factor, err = strconv.ParseFloat(mf.fd4Factor, 64); err != nil {
          mf.fd4Result = fmt.Sprintf("Error: %s -- %+v", mf.fd4Factor, err)
        } else {
          var a finances.Annuities
          mf.fd4Result = fmt.Sprintf("Growth/Decay: %.3f %s", a.GrowthDecayOfFunds(factor,
           ir / 100.0, a.GetCompoundingPeriod(mf.fd4Compound[0], true)),
           a.TimePeriods(mf.fd4Compound))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("interest rate = %s, cp = %s, factor = %s, %s\n", mf.fd4Interest,
           mf.fd4Compound, mf.fd4Factor, mf.fd4Result),
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
      } { "Miscellaneous", m.DTF(), mf.currentButton, newSession.CsrfToken,
          mf.fd4Interest, mf.fd4Compound, mf.fd4Factor, mf.fd4Result,
        })
    } else if strings.EqualFold(mf.currentPage, "rhs-ui5") {
      mf.currentButton = "lhs-button5"
      if req.Method == http.MethodPost {
        mf.fd5Values = req.PostFormValue("fd5-values")
        split := strings.Split(mf.fd5Values, ";")
        values := make([]float64, len(split))
        var err error
        for i, s := range split {
          if values[i], err = strconv.ParseFloat(s, 64); err != nil {
            mf.fd5Result[1] = fmt.Sprintf("Error: %s -- %+v", s, err)
            break;
          }
        }
        //
        if err == nil {
          var a finances.Annuities
          mf.fd5Result[1] = fmt.Sprintf("Avg: %.3f%%", a.AverageRateOfReturn(values) * 100.0)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("values = [%s], %s\n", mf.fd5Values, mf.fd5Result[1]),
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
      } { "Miscellaneous", m.DTF(), mf.currentButton, newSession.CsrfToken,
          mf.fd5Values, mf.fd5Result,
        })
    } else if strings.EqualFold(mf.currentPage, "rhs-ui6") {
      mf.currentButton = "lhs-button6"
      if req.Method == http.MethodPost {
        mf.fd6Time = req.PostFormValue("fd6-time")
        mf.fd6TimePeriod = req.PostFormValue("fd6-tp")
        mf.fd6Rate = req.PostFormValue("fd6-rate")
        mf.fd6Compound = req.PostFormValue("fd6-compound")
        mf.fd6PV = req.PostFormValue("fd6-pv")
        var time float64
        var rate float64
        var pv float64
        var err error
        if time, err = strconv.ParseFloat(mf.fd6Time, 64); err != nil {
          mf.fd6Result = fmt.Sprintf("Error: %s -- %+v", mf.fd6Time, err)
        } else if rate, err = strconv.ParseFloat(mf.fd6Rate, 64); err != nil {
          mf.fd6Result = fmt.Sprintf("Error: %s -- %+v", mf.fd6Rate, err)
        } else if pv, err = strconv.ParseFloat(mf.fd6PV, 64); err != nil {
          mf.fd6Result = fmt.Sprintf("Error: %s -- %+v", mf.fd6PV, err)
        } else {
          var a finances.Annuities
          mf.fd6Result = fmt.Sprintf("Future Value: %.2f", a.Depreciation(pv, rate / 100.0,
           a.GetCompoundingPeriod(mf.fd6Compound[0], false), time,
           a.GetTimePeriod(mf.fd6TimePeriod[0], false)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("time = %s, tp = %s, rate = %s, cp = %s, pv = %s, %s\n", mf.fd6Time,
           mf.fd6TimePeriod, mf.fd6Rate, mf.fd6Compound, mf.fd6PV, mf.fd6Result),
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
      } { "Miscellaneous", m.DTF(), mf.currentButton, newSession.CsrfToken,
          mf.fd6Time, mf.fd6TimePeriod, mf.fd6Rate, mf.fd6Compound, mf.fd6PV, mf.fd6Result,
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
        mf.fd1Result[1] = ""
      } else if strings.EqualFold(mf.currentPage, "rhs-ui2") {
        mf.fd2Result[2] = ""
      } else if strings.EqualFold(mf.currentPage, "rhs-ui3") {
        mf.fd3Result[3] = ""
      } else if strings.EqualFold(mf.currentPage, "rhs-ui4") {
        mf.fd4Result = ""
      } else if strings.EqualFold(mf.currentPage, "rhs-ui5") {
        mf.fd5Result[1] = ""
      } else if strings.EqualFold(mf.currentPage, "rhs-ui6") {
        mf.fd6Result = ""
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
