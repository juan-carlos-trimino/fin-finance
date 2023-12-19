package webfinances

import (
  "context"
  "finance/middlewares"
  "finance/finances"
	"finance/sessions"
  "fmt"
  "html/template"
  "net/http"
  "strconv"
  "strings"
)

type WfOaCpPages struct {
}

func (o WfOaCpPages) OaCpPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
    return
  }
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering OaCpPages/webfinances.",
  })
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    of := GetOaCpFields(sessions.GetUserName(sessionToken))
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
      of.currentPage = ui
    }
    //
    if strings.EqualFold(of.currentPage, "rhs-ui1") {
      of.currentButton = "lhs-button1"
      if req.Method == http.MethodPost {
        of.fd1Interest = req.PostFormValue("fd1-interest")
        of.fd1Compound = req.PostFormValue("fd1-cp")
        of.fd1PV = req.PostFormValue("fd1-pv")
        of.fd1FV = req.PostFormValue("fd1-fv")
        var i float64
        var pv float64
        var fv float64
        var err error
        if i, err = strconv.ParseFloat(of.fd1Interest, 64); err != nil {
          of.fd1Result = fmt.Sprintf("Error: %s -- %+v", of.fd1Interest, err)
        } else if pv, err = strconv.ParseFloat(of.fd1PV, 64); err != nil {
          of.fd1Result = fmt.Sprintf("Error: %s -- %+v", of.fd1PV, err)
        } else if fv, err = strconv.ParseFloat(of.fd1FV, 64); err != nil {
          of.fd1Result = fmt.Sprintf("Error: %s -- %+v", of.fd1FV, err)
        } else {
          var oa finances.Annuities
          of.fd1Result = fmt.Sprintf("Compounding Period: %.5f %s", oa.O_Periods_PV_FV(pv, fv,
            i / 100.0, oa.GetCompoundingPeriod(of.fd1Compound[0], true)),
            oa.TimePeriods(of.fd1Compound))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("i = %s, cp = %s, pv = %s, fv = %s, %s",
            of.fd1Interest, of.fd1Compound, of.fd1PV, of.fd1FV, of.fd1Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/cp/cp.html",
        "webfinances/templates/header.html",
        "webfinances/templates/ordinaryannuity/cp/i-PV-FV.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "oacompoundingperiods", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd1Interest string
        Fd1Compound string
        Fd1PV string
        Fd1FV string
        Fd1Result string
      } { "Ordinary Annuity / Compounding Periods", m.DTF(), of.currentButton,
          newSession.CsrfToken, of.fd1Interest, of.fd1Compound, of.fd1PV, of.fd1FV, of.fd1Result,
        })
    } else if strings.EqualFold(of.currentPage, "rhs-ui2") {
      of.currentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        of.fd2Interest = req.FormValue("fd2-interest")
        of.fd2Compound = req.PostFormValue("fd2-cp")
        of.fd2Payment = req.PostFormValue("fd2-payment")
        of.fd2PV = req.PostFormValue("fd2-pv")
        var i float64
        var pmt float64
        var pv float64
        var err error
        if i, err = strconv.ParseFloat(of.fd2Interest, 64); err != nil {
          of.fd2Result = fmt.Sprintf("Error: %s -- %+v", of.fd2Interest, err)
        } else if pmt, err = strconv.ParseFloat(of.fd2Payment, 64); err != nil {
          of.fd2Result = fmt.Sprintf("Error: %s -- %+v", of.fd2Payment, err)
        } else if pv, err = strconv.ParseFloat(of.fd2PV, 64); err != nil {
          of.fd2Result = fmt.Sprintf("Error: %s -- %+v", of.fd2PV, err)
        } else {
          var oa finances.Annuities
          of.fd2Result = fmt.Sprintf("Compounding Period: %.5f %s", oa.O_Periods_PMT_PV(pmt, pv,
            i / 100.0, oa.GetCompoundingPeriod(of.fd2Compound[0], true)),
            oa.TimePeriods(of.fd2Compound))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("i = %s, cp = %s, pmt = %s, pv = %s, %s",
            of.fd2Interest, of.fd2Compound, of.fd2Payment, of.fd2PV, of.fd2Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/cp/cp.html",
        "webfinances/templates/header.html",
        "webfinances/templates/ordinaryannuity/cp/i-PMT-PV.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "oacompoundingperiods", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd2Interest string
        Fd2Compound string
        Fd2Payment string
        Fd2PV string
        Fd2Result string
      } { "Ordinary Annuity / Compounding Periods", m.DTF(), of.currentButton,
          newSession.CsrfToken, of.fd2Interest, of.fd2Compound, of.fd2Payment, of.fd2PV,
          of.fd2Result,
        })
    } else if strings.EqualFold(of.currentPage, "rhs-ui3") {
      of.currentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        of.fd3Interest = req.FormValue("fd3-interest")
        of.fd3Compound = req.PostFormValue("fd3-cp")
        of.fd3Payment = req.PostFormValue("fd3-payment")
        of.fd3FV = req.PostFormValue("fd3-fv")
        var i float64
        var pmt float64
        var fv float64
        var err error
        if i, err = strconv.ParseFloat(of.fd3Interest, 64); err != nil {
          of.fd3Result = fmt.Sprintf("Error: %s -- %+v", of.fd3Interest, err)
        } else if pmt, err = strconv.ParseFloat(of.fd3Payment, 64); err != nil {
          of.fd3Result = fmt.Sprintf("Error: %s -- %+v", of.fd3Payment, err)
        } else if fv, err = strconv.ParseFloat(of.fd3FV, 64); err != nil {
          of.fd3Result = fmt.Sprintf("Error: %s -- %+v", of.fd3FV, err)
        } else {
          var oa finances.Annuities
          of.fd3Result = fmt.Sprintf("Compounding Period: %.5f %s", oa.O_Periods_PMT_FV(pmt, fv,
            i / 100.0, oa.GetCompoundingPeriod(of.fd3Compound[0], true)),
            oa.TimePeriods(of.fd3Compound))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("i = %s, cp = %s, pmt = %s, fv = %s, %s", of.fd3Interest,
            of.fd3Compound, of.fd3Payment, of.fd3FV, of.fd3Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/cp/cp.html",
        "webfinances/templates/header.html",
        "webfinances/templates/ordinaryannuity/cp/i-PMT-FV.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "oacompoundingperiods", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd3Interest string
        Fd3Compound string
        Fd3Payment string
        Fd3FV string
        Fd3Result string
      } { "Ordinary Annuity / Compounding Periods", m.DTF(), of.currentButton,
          newSession.CsrfToken, of.fd3Interest, of.fd3Compound, of.fd3Payment, of.fd3FV,
          of.fd3Result,
        })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", of.currentPage)
      fmt.Printf("%s - %s\n", m.DTF(), errString)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      fmt.Println("*** Request timeout ***")
      if strings.EqualFold(of.currentPage, "rhs-ui1") {
        of.fd1Result = ""
      } else if strings.EqualFold(of.currentPage, "rhs-ui2") {
        of.fd2Result = ""
      } else if strings.EqualFold(of.currentPage, "rhs-ui3") {
        of.fd3Result = ""
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
