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

type WfOaPerpetuityPages interface {
  OaPerpetuityPages(http.ResponseWriter, *http.Request)
}

type wfOaPerpetuityPages struct {
  currentPage string
  currentButton string
  //
  fd1Interest string
  fd1Compound string
  fd1Pmt string
  fd1Result string
  //
  fd2Interest string
  fd2Compound string
  fd2Grow string
  fd2Pmt string
  fd2Result string
}

func NewWfOaPerpetuityPages() WfOaPerpetuityPages {
  return &wfOaPerpetuityPages {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1Interest: "1.00",
    fd1Compound: "annually",
    fd1Pmt: "1.00",
    fd1Result: "",
    //
    fd2Interest: "1.00",
    fd2Compound: "annually",
    fd2Grow: "1.00",
    fd2Pmt: "1.00",
    fd2Result: "",
  }
}

func (p *wfOaPerpetuityPages) OaPerpetuityPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
    return
  }
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering OaPerpetuityPages/webfinances.",
  })
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
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
      p.currentPage = ui
    }
    //
    if strings.EqualFold(p.currentPage, "rhs-ui1") {
      p.currentButton = "lhs-button1"
      if req.Method == http.MethodPost {
        p.fd1Interest = req.PostFormValue("fd1-interest")
        p.fd1Compound = req.PostFormValue("fd1-cp")
        p.fd1Pmt = req.PostFormValue("fd1-pmt")
        var i float64
        var pmt float64
        var err error
        if i, err = strconv.ParseFloat(p.fd1Interest, 64); err != nil {
          p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1Interest, err)
        } else if pmt, err = strconv.ParseFloat(p.fd1Pmt, 64); err != nil {
          p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1Pmt, err)
        } else {
          var oa finances.Annuities
          p.fd1Result = fmt.Sprintf("Present Value of Perpetuity: $%.2f", oa.O_Perpetuity(
                                  i / 100.0, pmt, oa.GetCompoundingPeriod(p.fd1Compound[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("i = %s, cp = %s, pmt = %s, %s",
                      p.fd1Interest, p.fd1Compound, p.fd1Pmt, p.fd1Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/perpetuity/perpetuity.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/ordinaryannuity/perpetuity/p.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "oaperpetuity", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd1Interest string
        Fd1Compound string
        Fd1Pmt string
        Fd1Result string
      } { "Ordinary Annuity / Perpetuities", m.DTF(), p.currentButton, newSession.CsrfToken,
          p.fd1Interest, p.fd1Compound, p.fd1Pmt, p.fd1Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui2") {
      p.currentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        p.fd2Interest = req.FormValue("fd2-interest")
        p.fd2Compound = req.PostFormValue("fd2-cp")
        p.fd2Grow = req.PostFormValue("fd2-grow")
        p.fd2Pmt = req.PostFormValue("fd2-pmt")
        var i float64
        var grow float64
        var pmt float64
        var err error
        if i, err = strconv.ParseFloat(p.fd2Interest, 64); err != nil {
          p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2Interest, err)
        } else if grow, err = strconv.ParseFloat(p.fd2Grow, 64); err != nil {
          p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2Grow, err)
        } else if pmt, err = strconv.ParseFloat(p.fd2Pmt, 64); err != nil {
          p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2Pmt, err)
        } else {
          var oa finances.Annuities
          p.fd2Result = fmt.Sprintf("Present Value of Perpetuity: $%.2f", oa.O_GrowingPerpetuity(
                            i / 100.0, grow, pmt, oa.GetCompoundingPeriod(p.fd2Compound[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("i = %s, cp = %s, grow = %s, pmt = %s, %s",
                      p.fd2Interest, p.fd2Compound, p.fd2Grow, p.fd2Pmt, p.fd2Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/perpetuity/perpetuity.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/ordinaryannuity/perpetuity/gp.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "oaperpetuity", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd2Interest string
        Fd2Compound string
        Fd2Grow string
        Fd2Pmt string
        Fd2Result string
      } { "Ordinary Annuity / Perpetuities", m.DTF(), p.currentButton, newSession.CsrfToken,
          p.fd2Interest, p.fd2Compound, p.fd2Grow, p.fd2Pmt, p.fd2Result,
        })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", p.currentPage)
      fmt.Printf("%s - %s\n", m.DTF(), errString)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      fmt.Println("*** Request timeout ***")
      if strings.EqualFold(p.currentPage, "rhs-ui1") {
        p.fd1Result = ""
      } else if strings.EqualFold(p.currentPage, "rhs-ui2") {
        p.fd2Result = ""
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
