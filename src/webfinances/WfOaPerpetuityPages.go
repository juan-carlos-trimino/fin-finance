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

type WfOaPerpetuityPages struct {
}

func (o WfOaPerpetuityPages) OaPerpetuityPages(res http.ResponseWriter, req *http.Request) {
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
    userName := sessions.GetUserName(sessionToken)
    of := getOaPerpetuityFields(userName)
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
      of.CurrentPage = ui
    }
    //
    if strings.EqualFold(of.CurrentPage, "rhs-ui1") {
      of.CurrentButton = "lhs-button1"
      if req.Method == http.MethodPost {
        of.Fd1Interest = req.PostFormValue("fd1-interest")
        of.Fd1Compound = req.PostFormValue("fd1-cp")
        of.Fd1Pmt = req.PostFormValue("fd1-pmt")
        var i float64
        var pmt float64
        var err error
        if i, err = strconv.ParseFloat(of.Fd1Interest, 64); err != nil {
          of.Fd1Result = fmt.Sprintf("Error: %s -- %+v", of.Fd1Interest, err)
        } else if pmt, err = strconv.ParseFloat(of.Fd1Pmt, 64); err != nil {
          of.Fd1Result = fmt.Sprintf("Error: %s -- %+v", of.Fd1Pmt, err)
        } else {
          var oa finances.Annuities
          of.Fd1Result = fmt.Sprintf("Present Value of Perpetuity: $%.2f",
            oa.O_Perpetuity(i / 100.0, pmt, oa.GetCompoundingPeriod(of.Fd1Compound[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("i = %s, cp = %s, pmt = %s, %s",
            of.Fd1Interest, of.Fd1Compound, of.Fd1Pmt, of.Fd1Result),
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
      } { "Ordinary Annuity / Perpetuities", m.DTF(), of.CurrentButton, newSession.CsrfToken,
          of.Fd1Interest, of.Fd1Compound, of.Fd1Pmt, of.Fd1Result,
        })
    } else if strings.EqualFold(of.CurrentPage, "rhs-ui2") {
      of.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        of.Fd2Interest = req.FormValue("fd2-interest")
        of.Fd2Compound = req.PostFormValue("fd2-cp")
        of.Fd2Grow = req.PostFormValue("fd2-grow")
        of.Fd2Pmt = req.PostFormValue("fd2-pmt")
        var i float64
        var grow float64
        var pmt float64
        var err error
        if i, err = strconv.ParseFloat(of.Fd2Interest, 64); err != nil {
          of.Fd2Result = fmt.Sprintf("Error: %s -- %+v", of.Fd2Interest, err)
        } else if grow, err = strconv.ParseFloat(of.Fd2Grow, 64); err != nil {
          of.Fd2Result = fmt.Sprintf("Error: %s -- %+v", of.Fd2Grow, err)
        } else if pmt, err = strconv.ParseFloat(of.Fd2Pmt, 64); err != nil {
          of.Fd2Result = fmt.Sprintf("Error: %s -- %+v", of.Fd2Pmt, err)
        } else {
          var oa finances.Annuities
          of.Fd2Result = fmt.Sprintf("Present Value of Perpetuity: $%.2f",
            oa.O_GrowingPerpetuity(i / 100.0, grow, pmt,
            oa.GetCompoundingPeriod(of.Fd2Compound[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("i = %s, cp = %s, grow = %s, pmt = %s, %s",
            of.Fd2Interest, of.Fd2Compound, of.Fd2Grow, of.Fd2Pmt, of.Fd2Result),
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
      } { "Ordinary Annuity / Perpetuities", m.DTF(), of.CurrentButton, newSession.CsrfToken,
          of.Fd2Interest, of.Fd2Compound, of.Fd2Grow, of.Fd2Pmt, of.Fd2Result,
        })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", of.CurrentPage)
      fmt.Printf("%s - %s\n", m.DTF(), errString)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      fmt.Println("*** Request timeout ***")
      if strings.EqualFold(of.CurrentPage, "rhs-ui1") {
        of.Fd1Result = ""
      } else if strings.EqualFold(of.CurrentPage, "rhs-ui2") {
        of.Fd2Result = ""
      }
    }
    //
    if data, err := json.Marshal(of); err != nil {
      fmt.Printf("%s - %s\n", m.DTF(), err)
    } else {
      filePath := fmt.Sprintf("%s/%s/oaperpetuity.txt", mainDir, userName)
      if _, err := misc.WriteAllExclusiveLock(filePath, data, os.O_WRONLY, 0o220); err != nil {
        fmt.Printf("%s - %s\n", m.DTF(), err)
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
