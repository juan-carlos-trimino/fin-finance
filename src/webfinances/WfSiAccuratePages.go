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

type WfSiAccuratePages struct {
}

func (s WfSiAccuratePages) SimpleInterestAccuratePages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
    return
  }
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering SimpleInterestAccuratePages/webfinances.",
  })
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    sif := GetSiAccurateFields(sessions.GetUserName(sessionToken))
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
      sif.currentPage = ui
    }
    //
    if strings.EqualFold(sif.currentPage, "rhs-ui1") {
      sif.currentButton = "lhs-button1"
      if req.Method == http.MethodPost {
        sif.fd1Time = req.PostFormValue("fd1-time")
        sif.fd1TimePeriod = req.PostFormValue("fd1-tp")
        sif.fd1Interest = req.PostFormValue("fd1-interest")
        sif.fd1Compound = req.PostFormValue("fd1-compound")
        sif.fd1PV = req.PostFormValue("fd1-pv")
        var n float64
        var i float64
        var pv float64
        var err error
        if n, err = strconv.ParseFloat(sif.fd1Time, 64); err != nil {
          sif.fd1Result = fmt.Sprintf("Error: %s -- %+v", sif.fd1Time, err)
        } else if i, err = strconv.ParseFloat(sif.fd1Interest, 64); err != nil {
          sif.fd1Result = fmt.Sprintf("Error: %s -- %+v", sif.fd1Interest, err)
        } else if pv, err = strconv.ParseFloat(sif.fd1PV, 64); err != nil {
          sif.fd1Result = fmt.Sprintf("Error: %s -- %+v", sif.fd1PV, err)
        } else {
          var si finances.SimpleInterest
          var periods finances.Periods
          sif.fd1Result = fmt.Sprintf("Amount of Interest: $%.2f",
            si.AccurateInterest(pv, i / 100.0,
            periods.GetCompoundingPeriod(sif.fd1Compound[0], true), n,
            periods.GetTimePeriod(sif.fd1TimePeriod[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, pv = %s, %s", sif.fd1Time,
            sif.fd1TimePeriod, sif.fd1Interest, sif.fd1Compound, sif.fd1PV, sif.fd1Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      t := template.Must(template.ParseFiles("webfinances/templates/simpleinterestaccurate/accurate.html",
        "webfinances/templates/header.html",
        "webfinances/templates/simpleinterestaccurate/amountofinterest.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "simpleinterestaccurate", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd1Time string
        Fd1TimePeriod string
        Fd1Interest string
        Fd1Compound string
        Fd1PV string
        Fd1Result string
      } { "Simple Interest / Accurate (Exact) Interest", m.DTF(), sif.currentButton,
          newSession.CsrfToken, sif.fd1Time, sif.fd1TimePeriod, sif.fd1Interest, sif.fd1Compound,
          sif.fd1PV, sif.fd1Result,
        })
    } else if strings.EqualFold(sif.currentPage, "rhs-ui2") {
      sif.currentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        sif.fd2Time = req.PostFormValue("fd2-time")
        sif.fd2TimePeriod = req.PostFormValue("fd2-tp")
        sif.fd2Amount = req.PostFormValue("fd2-amount")
        sif.fd2PV = req.PostFormValue("fd2-pv")
        var n float64
        var a float64
        var pv float64
        var err error
        if n, err = strconv.ParseFloat(sif.fd2Time, 64); err != nil {
          sif.fd2Result = fmt.Sprintf("Error: %s -- %+v", sif.fd2Time, err)
        } else if a, err = strconv.ParseFloat(sif.fd2Amount, 64); err != nil {
          sif.fd2Result = fmt.Sprintf("Error: %s -- %+v", sif.fd2Amount, err)
        } else if pv, err = strconv.ParseFloat(sif.fd2PV, 64); err != nil {
          sif.fd2Result = fmt.Sprintf("Error: %s -- %+v", sif.fd2PV, err)
        } else {
          var si finances.SimpleInterest
          var periods finances.Periods
          sif.fd2Result = fmt.Sprintf("Interest Rate: %.3f%%",
            si.AccurateRate(pv, a, n, periods.GetTimePeriod(sif.fd2TimePeriod[0], true)) * 100.0)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, a = %s, pv = %s, %s", sif.fd2Time, sif.fd2TimePeriod,
            sif.fd2Amount, sif.fd2PV, sif.fd2Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/simpleinterestaccurate/accurate.html",
        "webfinances/templates/header.html",
        "webfinances/templates/simpleinterestaccurate/interestrate.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "simpleinterestaccurate", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd2Time string
        Fd2TimePeriod string
        Fd2Amount string
        Fd2PV string
        Fd2Result string
      } { "Simple Interest / Accurate (Exact) Interest", m.DTF(), sif.currentButton,
          newSession.CsrfToken, sif.fd2Time, sif.fd2TimePeriod, sif.fd2Amount, sif.fd2PV,
          sif.fd2Result,
        })
    } else if strings.EqualFold(sif.currentPage, "rhs-ui3") {
      sif.currentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        sif.fd3Time = req.PostFormValue("fd3-time")
        sif.fd3TimePeriod = req.PostFormValue("fd3-tp")
        sif.fd3Interest = req.PostFormValue("fd3-interest")
        sif.fd3Compound = req.PostFormValue("fd3-compound")
        sif.fd3Amount = req.PostFormValue("fd3-amount")
        var n float64
        var i float64
        var a float64
        var err error
        if n, err = strconv.ParseFloat(sif.fd3Time, 64); err != nil {
          sif.fd3Result = fmt.Sprintf("Error: %s -- %+v", sif.fd3Time, err)
        } else if i, err = strconv.ParseFloat(sif.fd3Interest, 64); err != nil {
          sif.fd3Result = fmt.Sprintf("Error: %s -- %+v", sif.fd3Interest, err)
        } else if a, err = strconv.ParseFloat(sif.fd3Amount, 64); err != nil {
          sif.fd3Result = fmt.Sprintf("Error: %s -- %+v", sif.fd3Amount, err)
        } else {
          var si finances.SimpleInterest
          var periods finances.Periods
          sif.fd3Result = fmt.Sprintf("Principal: $%.2f", si.AccuratePrincipal(a, i / 100.0,
            periods.GetCompoundingPeriod(sif.fd3Compound[0], true), n,
            periods.GetTimePeriod(sif.fd3TimePeriod[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, a = %s, %s\n", sif.fd3Time,
            sif.fd3TimePeriod, sif.fd3Interest, sif.fd3Compound, sif.fd3Amount, sif.fd3Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/simpleinterestaccurate/accurate.html",
        "webfinances/templates/header.html",
        "webfinances/templates/simpleinterestaccurate/principal.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "simpleinterestaccurate", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd3Time string
        Fd3TimePeriod string
        Fd3Interest string
        Fd3Compound string
        Fd3Amount string
        Fd3Result string
      } { "Simple Interest / Accurate (Exact) Interest", m.DTF(), sif.currentButton,
          newSession.CsrfToken, sif.fd3Time, sif.fd3TimePeriod, sif.fd3Interest, sif.fd3Compound,
          sif.fd3Amount, sif.fd3Result,
        })
    } else if strings.EqualFold(sif.currentPage, "rhs-ui4") {
      sif.currentButton = "lhs-button4"
      if req.Method == http.MethodPost {
        sif.fd4Interest = req.FormValue("fd4-interest")
        sif.fd4Compound = req.FormValue("fd4-compound")
        sif.fd4Amount = req.FormValue("fd4-amount")
        sif.fd4PV = req.FormValue("fd4-pv")
        var i float64
        var a float64
        var pv float64
        var err error
        if i, err = strconv.ParseFloat(sif.fd4Interest, 64); err != nil {
          sif.fd4Result = fmt.Sprintf("Error: %s -- %+v", sif.fd4Interest, err)
        } else if a, err = strconv.ParseFloat(sif.fd4Amount, 64); err != nil {
          sif.fd4Result = fmt.Sprintf("Error: %s -- %+v", sif.fd4Amount, err)
        } else if pv, err = strconv.ParseFloat(sif.fd4PV, 64); err != nil {
          sif.fd4Result = fmt.Sprintf("Error: %s -- %+v", sif.fd4PV, err)
        } else {
          var si finances.SimpleInterest
          var periods finances.Periods
          sif.fd4Result = fmt.Sprintf("Time: %.3f %s", si.AccurateTime(pv, a, i / 100.0,
            periods.GetCompoundingPeriod(sif.fd4Compound[0], true)),
            periods.TimePeriods(sif.fd4Compound))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("i = %s, cp = %s, a = %s, pv = %s, %s\n", sif.fd4Interest, sif.fd4Compound,
            sif.fd4Amount, sif.fd4PV, sif.fd4Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/simpleinterestaccurate/accurate.html",
        "webfinances/templates/header.html",
        "webfinances/templates/simpleinterestaccurate/time.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "simpleinterestaccurate", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd4Interest string
        Fd4Compound string
        Fd4Amount string
        Fd4PV string
        Fd4Result string
      } { "Simple Interest / Accurate (Exact) Interest", m.DTF(), sif.currentButton,
          newSession.CsrfToken, sif.fd4Interest, sif.fd4Compound, sif.fd4Amount, sif.fd4PV,
          sif.fd4Result,
        })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", sif.currentPage)
      fmt.Printf("%s - %s\n", m.DTF(), errString)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      fmt.Println("*** Request timeout ***")
      if strings.EqualFold(sif.currentPage, "rhs-ui1") {
        sif.fd1Result = ""
      } else if strings.EqualFold(sif.currentPage, "rhs-ui2") {
        sif.fd2Result = ""
      } else if strings.EqualFold(sif.currentPage, "rhs-ui3") {
        sif.fd3Result = ""
      } else if strings.EqualFold(sif.currentPage, "rhs-ui4") {
        sif.fd4Result = ""
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
