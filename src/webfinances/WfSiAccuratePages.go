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
    userName := sessions.GetUserName(sessionToken)
    sif := getSiAccurateFields(userName)
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
      sif.CurrentPage = ui
    }
    //
    if strings.EqualFold(sif.CurrentPage, "rhs-ui1") {
      sif.CurrentButton = "lhs-button1"
      if req.Method == http.MethodPost {
        sif.Fd1Time = req.PostFormValue("fd1-time")
        sif.Fd1TimePeriod = req.PostFormValue("fd1-tp")
        sif.Fd1Interest = req.PostFormValue("fd1-interest")
        sif.Fd1Compound = req.PostFormValue("fd1-compound")
        sif.Fd1PV = req.PostFormValue("fd1-pv")
        var n float64
        var i float64
        var pv float64
        var err error
        if n, err = strconv.ParseFloat(sif.Fd1Time, 64); err != nil {
          sif.Fd1Result = fmt.Sprintf("Error: %s -- %+v", sif.Fd1Time, err)
        } else if i, err = strconv.ParseFloat(sif.Fd1Interest, 64); err != nil {
          sif.Fd1Result = fmt.Sprintf("Error: %s -- %+v", sif.Fd1Interest, err)
        } else if pv, err = strconv.ParseFloat(sif.Fd1PV, 64); err != nil {
          sif.Fd1Result = fmt.Sprintf("Error: %s -- %+v", sif.Fd1PV, err)
        } else {
          var si finances.SimpleInterest
          var periods finances.Periods
          sif.Fd1Result = fmt.Sprintf("Amount of Interest: $%.2f",
            si.AccurateInterest(pv, i / 100.0,
            periods.GetCompoundingPeriod(sif.Fd1Compound[0], true), n,
            periods.GetTimePeriod(sif.Fd1TimePeriod[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, pv = %s, %s", sif.Fd1Time,
            sif.Fd1TimePeriod, sif.Fd1Interest, sif.Fd1Compound, sif.Fd1PV, sif.Fd1Result),
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
      } { "Simple Interest / Accurate (Exact) Interest", m.DTF(), sif.CurrentButton,
          newSession.CsrfToken, sif.Fd1Time, sif.Fd1TimePeriod, sif.Fd1Interest, sif.Fd1Compound,
          sif.Fd1PV, sif.Fd1Result,
        })
    } else if strings.EqualFold(sif.CurrentPage, "rhs-ui2") {
      sif.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        sif.Fd2Time = req.PostFormValue("fd2-time")
        sif.Fd2TimePeriod = req.PostFormValue("fd2-tp")
        sif.Fd2Amount = req.PostFormValue("fd2-amount")
        sif.Fd2PV = req.PostFormValue("fd2-pv")
        var n float64
        var a float64
        var pv float64
        var err error
        if n, err = strconv.ParseFloat(sif.Fd2Time, 64); err != nil {
          sif.Fd2Result = fmt.Sprintf("Error: %s -- %+v", sif.Fd2Time, err)
        } else if a, err = strconv.ParseFloat(sif.Fd2Amount, 64); err != nil {
          sif.Fd2Result = fmt.Sprintf("Error: %s -- %+v", sif.Fd2Amount, err)
        } else if pv, err = strconv.ParseFloat(sif.Fd2PV, 64); err != nil {
          sif.Fd2Result = fmt.Sprintf("Error: %s -- %+v", sif.Fd2PV, err)
        } else {
          var si finances.SimpleInterest
          var periods finances.Periods
          sif.Fd2Result = fmt.Sprintf("Interest Rate: %.3f%%",
            si.AccurateRate(pv, a, n, periods.GetTimePeriod(sif.Fd2TimePeriod[0], true)) * 100.0)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, a = %s, pv = %s, %s", sif.Fd2Time, sif.Fd2TimePeriod,
            sif.Fd2Amount, sif.Fd2PV, sif.Fd2Result),
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
      } { "Simple Interest / Accurate (Exact) Interest", m.DTF(), sif.CurrentButton,
          newSession.CsrfToken, sif.Fd2Time, sif.Fd2TimePeriod, sif.Fd2Amount, sif.Fd2PV,
          sif.Fd2Result,
        })
    } else if strings.EqualFold(sif.CurrentPage, "rhs-ui3") {
      sif.CurrentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        sif.Fd3Time = req.PostFormValue("fd3-time")
        sif.Fd3TimePeriod = req.PostFormValue("fd3-tp")
        sif.Fd3Interest = req.PostFormValue("fd3-interest")
        sif.Fd3Compound = req.PostFormValue("fd3-compound")
        sif.Fd3Amount = req.PostFormValue("fd3-amount")
        var n float64
        var i float64
        var a float64
        var err error
        if n, err = strconv.ParseFloat(sif.Fd3Time, 64); err != nil {
          sif.Fd3Result = fmt.Sprintf("Error: %s -- %+v", sif.Fd3Time, err)
        } else if i, err = strconv.ParseFloat(sif.Fd3Interest, 64); err != nil {
          sif.Fd3Result = fmt.Sprintf("Error: %s -- %+v", sif.Fd3Interest, err)
        } else if a, err = strconv.ParseFloat(sif.Fd3Amount, 64); err != nil {
          sif.Fd3Result = fmt.Sprintf("Error: %s -- %+v", sif.Fd3Amount, err)
        } else {
          var si finances.SimpleInterest
          var periods finances.Periods
          sif.Fd3Result = fmt.Sprintf("Principal: $%.2f", si.AccuratePrincipal(a, i / 100.0,
            periods.GetCompoundingPeriod(sif.Fd3Compound[0], true), n,
            periods.GetTimePeriod(sif.Fd3TimePeriod[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, a = %s, %s\n", sif.Fd3Time,
            sif.Fd3TimePeriod, sif.Fd3Interest, sif.Fd3Compound, sif.Fd3Amount, sif.Fd3Result),
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
      } { "Simple Interest / Accurate (Exact) Interest", m.DTF(), sif.CurrentButton,
          newSession.CsrfToken, sif.Fd3Time, sif.Fd3TimePeriod, sif.Fd3Interest, sif.Fd3Compound,
          sif.Fd3Amount, sif.Fd3Result,
        })
    } else if strings.EqualFold(sif.CurrentPage, "rhs-ui4") {
      sif.CurrentButton = "lhs-button4"
      if req.Method == http.MethodPost {
        sif.Fd4Interest = req.FormValue("fd4-interest")
        sif.Fd4Compound = req.FormValue("fd4-compound")
        sif.Fd4Amount = req.FormValue("fd4-amount")
        sif.Fd4PV = req.FormValue("fd4-pv")
        var i float64
        var a float64
        var pv float64
        var err error
        if i, err = strconv.ParseFloat(sif.Fd4Interest, 64); err != nil {
          sif.Fd4Result = fmt.Sprintf("Error: %s -- %+v", sif.Fd4Interest, err)
        } else if a, err = strconv.ParseFloat(sif.Fd4Amount, 64); err != nil {
          sif.Fd4Result = fmt.Sprintf("Error: %s -- %+v", sif.Fd4Amount, err)
        } else if pv, err = strconv.ParseFloat(sif.Fd4PV, 64); err != nil {
          sif.Fd4Result = fmt.Sprintf("Error: %s -- %+v", sif.Fd4PV, err)
        } else {
          var si finances.SimpleInterest
          var periods finances.Periods
          sif.Fd4Result = fmt.Sprintf("Time: %.3f %s", si.AccurateTime(pv, a, i / 100.0,
            periods.GetCompoundingPeriod(sif.Fd4Compound[0], true)),
            periods.TimePeriods(sif.Fd4Compound))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("i = %s, cp = %s, a = %s, pv = %s, %s\n", sif.Fd4Interest, sif.Fd4Compound,
            sif.Fd4Amount, sif.Fd4PV, sif.Fd4Result),
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
      } { "Simple Interest / Accurate (Exact) Interest", m.DTF(), sif.CurrentButton,
          newSession.CsrfToken, sif.Fd4Interest, sif.Fd4Compound, sif.Fd4Amount, sif.Fd4PV,
          sif.Fd4Result,
        })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", sif.CurrentPage)
      fmt.Printf("%s - %s\n", m.DTF(), errString)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      fmt.Println("*** Request timeout ***")
      if strings.EqualFold(sif.CurrentPage, "rhs-ui1") {
        sif.Fd1Result = ""
      } else if strings.EqualFold(sif.CurrentPage, "rhs-ui2") {
        sif.Fd2Result = ""
      } else if strings.EqualFold(sif.CurrentPage, "rhs-ui3") {
        sif.Fd3Result = ""
      } else if strings.EqualFold(sif.CurrentPage, "rhs-ui4") {
        sif.Fd4Result = ""
      }
    }
    //
    if data, err := json.Marshal(sif); err != nil {
      fmt.Printf("%s - %s\n", m.DTF(), err)
    } else {
      filePath := fmt.Sprintf("%s/%s/siaccurate.txt", mainDir, userName)
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
