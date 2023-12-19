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

type WfOaEppPages struct {
}

func (o WfOaEppPages) OaEppPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
    return
  }
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering OaEppPages/webfinances.",
  })
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    of := GetOaEppFields(sessions.GetUserName(sessionToken))
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
        of.fd1N = req.PostFormValue("fd1-n")
        of.fd1TimePeriod = req.PostFormValue("fd1-tp")
        of.fd1Interest = req.PostFormValue("fd1-interest")
        of.fd1Compound = req.PostFormValue("fd1-cp")
        of.fd1FV = req.PostFormValue("fd1-fv")
        var n float64
        var i float64
        var fv float64
        var err error
        if n, err = strconv.ParseFloat(of.fd1N, 64); err != nil {
          of.fd1Result = fmt.Sprintf("Error: %s -- %+v", of.fd1N, err)
        } else if i, err = strconv.ParseFloat(of.fd1Interest, 64); err != nil {
          of.fd1Result = fmt.Sprintf("Error: %s -- %+v", of.fd1Interest, err)
        } else if fv, err = strconv.ParseFloat(of.fd1FV, 64); err != nil {
          of.fd1Result = fmt.Sprintf("Error: %s -- %+v", of.fd1FV, err)
        } else {
          var oa finances.Annuities
          of.fd1Result = fmt.Sprintf("Payment: $%.2f", oa.O_Payment_FV(fv, i / 100.0,
            oa.GetCompoundingPeriod(of.fd1Compound[0], true), n,
            oa.GetTimePeriod(of.fd1TimePeriod[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, fv = %s, %s", of.fd1N,
            of.fd1TimePeriod, of.fd1Interest, of.fd1Compound, of.fd1FV, of.fd1Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/epp/epp.html",
        "webfinances/templates/header.html",
        "webfinances/templates/ordinaryannuity/epp/n-i-FV.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "oaequalperiodicpayments", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd1N string
        Fd1TimePeriod string
        Fd1Interest string
        Fd1Compound string
        Fd1FV string
        Fd1Result string
      } { "Ordinary Annuity / Equal Periodic Payments", m.DTF(), of.currentButton,
          newSession.CsrfToken, of.fd1N, of.fd1TimePeriod, of.fd1Interest, of.fd1Compound,
          of.fd1FV, of.fd1Result,
        })
    } else if strings.EqualFold(of.currentPage, "rhs-ui2") {
      of.currentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        of.fd2N = req.PostFormValue("fd2-n")
        of.fd2TimePeriod = req.PostFormValue("fd2-tp")
        of.fd2Interest = req.FormValue("fd2-interest")
        of.fd2Compound = req.PostFormValue("fd2-cp")
        of.fd2PV = req.PostFormValue("fd2-pv")
        var n float64
        var i float64
        var pv float64
        var err error
        if n, err = strconv.ParseFloat(of.fd2N, 64); err != nil {
          of.fd2Result = fmt.Sprintf("Error: %s -- %+v", of.fd2N, err)
        } else if i, err = strconv.ParseFloat(of.fd2Interest, 64); err != nil {
          of.fd2Result = fmt.Sprintf("Error: %s -- %+v", of.fd2Interest, err)
        } else if pv, err = strconv.ParseFloat(of.fd2PV, 64); err != nil {
          of.fd2Result = fmt.Sprintf("Error: %s -- %+v", of.fd2PV, err)
        } else {
          var oa finances.Annuities
          of.fd2Result = fmt.Sprintf("Payment: $%.2f", oa.O_Payment_PV(pv, i / 100.0,
            oa.GetCompoundingPeriod(of.fd2Compound[0], true), n,
            oa.GetTimePeriod(of.fd2TimePeriod[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, pv = %s, %s", of.fd2N,
            of.fd2TimePeriod, of.fd2Interest, of.fd2Compound, of.fd2PV, of.fd2Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/epp/epp.html",
        "webfinances/templates/header.html",
        "webfinances/templates/ordinaryannuity/epp/n-i-PV.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "oaequalperiodicpayments", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd2N string
        Fd2TimePeriod string
        Fd2Interest string
        Fd2Compound string
        Fd2PV string
        Fd2Result string
      } { "Ordinary Annuity / Equal Periodic Payments", m.DTF(), of.currentButton,
          newSession.CsrfToken, of.fd2N, of.fd2TimePeriod, of.fd2Interest, of.fd2Compound,
          of.fd2PV, of.fd2Result,
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
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
