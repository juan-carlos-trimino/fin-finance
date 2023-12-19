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

type WfOaGaPages struct {
}

func (o WfOaGaPages) OaGaPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
    return
  }
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering OaGaPages/webfinances.",
  })
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    of := GetOaGaFields(sessions.GetUserName(sessionToken))
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
        of.fd1Interest = req.PostFormValue("fd1-interest")
        of.fd1Compound = req.PostFormValue("fd1-cp")
        of.fd1Grow = req.PostFormValue("fd1-grow")
        of.fd1Pmt = req.PostFormValue("fd1-pmt")
        var n float64
        var i float64
        var grow float64
        var pmt float64
        var err error
        if n, err = strconv.ParseFloat(of.fd1N, 64); err != nil {
          of.fd1Result = fmt.Sprintf("Error: %s -- %+v", of.fd1N, err)
        } else if i, err = strconv.ParseFloat(of.fd1Interest, 64); err != nil {
          of.fd1Result = fmt.Sprintf("Error: %s -- %+v", of.fd1Interest, err)
        } else if grow, err = strconv.ParseFloat(of.fd1Grow, 64); err != nil {
          of.fd1Result = fmt.Sprintf("Error: %s -- %+v", of.fd1Grow, err)
        } else if pmt, err = strconv.ParseFloat(of.fd1Pmt, 64); err != nil {
          of.fd1Result = fmt.Sprintf("Error: %s -- %+v", of.fd1Pmt, err)
        } else {
          var oa finances.Annuities
          of.fd1Result = fmt.Sprintf("Future Value: $%.2f", oa.O_GrowingAnnuityFutureValue(pmt, n,
            grow, i / 100.0, oa.GetCompoundingPeriod(of.fd1Compound[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, i = %s, cp = %s, grow = %s, pmt = %s, %s",
            of.fd1N, of.fd1Interest, of.fd1Compound, of.fd1Grow, of.fd1Pmt, of.fd1Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/ga/ga.html",
        "webfinances/templates/header.html",
        "webfinances/templates/ordinaryannuity/ga/FV.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "oagrowingannuity", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd1N string
        Fd1Interest string
        Fd1Compound string
        Fd1Grow string
        Fd1Pmt string
        Fd1Result string
      } { "Ordinary Annuity / Growing Annuity", m.DTF(), of.currentButton, newSession.CsrfToken,
          of.fd1N, of.fd1Interest, of.fd1Compound, of.fd1Grow, of.fd1Pmt, of.fd1Result,
        })
    } else if strings.EqualFold(of.currentPage, "rhs-ui2") {
      of.currentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        of.fd2N = req.PostFormValue("fd2-n")
        of.fd2Interest = req.PostFormValue("fd2-interest")
        of.fd2Compound = req.PostFormValue("fd2-cp")
        of.fd2Grow = req.PostFormValue("fd2-grow")
        of.fd2Pmt = req.PostFormValue("fd2-pmt")
        var n float64
        var i float64
        var grow float64
        var pmt float64
        var err error
        if n, err = strconv.ParseFloat(of.fd2N, 64); err != nil {
          of.fd2Result = fmt.Sprintf("Error: %s -- %+v", of.fd2N, err)
        } else if i, err = strconv.ParseFloat(of.fd2Interest, 64); err != nil {
          of.fd2Result = fmt.Sprintf("Error: %s -- %+v", of.fd2Interest, err)
        } else if grow, err = strconv.ParseFloat(of.fd2Grow, 64); err != nil {
          of.fd2Result = fmt.Sprintf("Error: %s -- %+v", of.fd2Grow, err)
        } else if pmt, err = strconv.ParseFloat(of.fd2Pmt, 64); err != nil {
          of.fd2Result = fmt.Sprintf("Error: %s -- %+v", of.fd2Pmt, err)
        } else {
          var oa finances.Annuities
          of.fd2Result = fmt.Sprintf("Present Value: $%.2f", oa.O_GrowingAnnuityPresentValue(pmt,
            n, grow, i / 100.0, oa.GetCompoundingPeriod(of.fd2Compound[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, i = %s, cp = %s, grow = %s, pmt = %s, %s",
            of.fd2N, of.fd2Interest, of.fd2Compound, of.fd2Grow, of.fd2Pmt, of.fd2Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/ga/ga.html",
        "webfinances/templates/header.html",
        "webfinances/templates/ordinaryannuity/ga/PV.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "oagrowingannuity", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd2N string
        Fd2Interest string
        Fd2Compound string
        Fd2Grow string
        Fd2Pmt string
        Fd2Result string
      } { "Ordinary Annuity / Growing Annuity", m.DTF(), of.currentButton, newSession.CsrfToken,
          of.fd2N, of.fd2Interest, of.fd2Compound, of.fd2Grow, of.fd2Pmt, of.fd2Result,
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
