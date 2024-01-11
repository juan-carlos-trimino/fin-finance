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
    userName := sessions.GetUserName(sessionToken)
    of := getOaEppFields(userName)
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
        of.Fd1N = req.PostFormValue("fd1-n")
        of.Fd1TimePeriod = req.PostFormValue("fd1-tp")
        of.Fd1Interest = req.PostFormValue("fd1-interest")
        of.Fd1Compound = req.PostFormValue("fd1-cp")
        of.Fd1FV = req.PostFormValue("fd1-fv")
        var n float64
        var i float64
        var fv float64
        var err error
        if n, err = strconv.ParseFloat(of.Fd1N, 64); err != nil {
          of.Fd1Result = fmt.Sprintf("Error: %s -- %+v", of.Fd1N, err)
        } else if i, err = strconv.ParseFloat(of.Fd1Interest, 64); err != nil {
          of.Fd1Result = fmt.Sprintf("Error: %s -- %+v", of.Fd1Interest, err)
        } else if fv, err = strconv.ParseFloat(of.Fd1FV, 64); err != nil {
          of.Fd1Result = fmt.Sprintf("Error: %s -- %+v", of.Fd1FV, err)
        } else {
          var oa finances.Annuities
          of.Fd1Result = fmt.Sprintf("Payment: $%.2f", oa.O_Payment_FV(fv, i / 100.0,
            oa.GetCompoundingPeriod(of.Fd1Compound[0], true), n,
            oa.GetTimePeriod(of.Fd1TimePeriod[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, fv = %s, %s", of.Fd1N,
            of.Fd1TimePeriod, of.Fd1Interest, of.Fd1Compound, of.Fd1FV, of.Fd1Result),
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
      } { "Ordinary Annuity / Equal Periodic Payments", m.DTF(), of.CurrentButton,
          newSession.CsrfToken, of.Fd1N, of.Fd1TimePeriod, of.Fd1Interest, of.Fd1Compound,
          of.Fd1FV, of.Fd1Result,
        })
    } else if strings.EqualFold(of.CurrentPage, "rhs-ui2") {
      of.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        of.Fd2N = req.PostFormValue("fd2-n")
        of.Fd2TimePeriod = req.PostFormValue("fd2-tp")
        of.Fd2Interest = req.FormValue("fd2-interest")
        of.Fd2Compound = req.PostFormValue("fd2-cp")
        of.Fd2PV = req.PostFormValue("fd2-pv")
        var n float64
        var i float64
        var pv float64
        var err error
        if n, err = strconv.ParseFloat(of.Fd2N, 64); err != nil {
          of.Fd2Result = fmt.Sprintf("Error: %s -- %+v", of.Fd2N, err)
        } else if i, err = strconv.ParseFloat(of.Fd2Interest, 64); err != nil {
          of.Fd2Result = fmt.Sprintf("Error: %s -- %+v", of.Fd2Interest, err)
        } else if pv, err = strconv.ParseFloat(of.Fd2PV, 64); err != nil {
          of.Fd2Result = fmt.Sprintf("Error: %s -- %+v", of.Fd2PV, err)
        } else {
          var oa finances.Annuities
          of.Fd2Result = fmt.Sprintf("Payment: $%.2f", oa.O_Payment_PV(pv, i / 100.0,
            oa.GetCompoundingPeriod(of.Fd2Compound[0], true), n,
            oa.GetTimePeriod(of.Fd2TimePeriod[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, pv = %s, %s", of.Fd2N,
            of.Fd2TimePeriod, of.Fd2Interest, of.Fd2Compound, of.Fd2PV, of.Fd2Result),
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
      } { "Ordinary Annuity / Equal Periodic Payments", m.DTF(), of.CurrentButton,
          newSession.CsrfToken, of.Fd2N, of.Fd2TimePeriod, of.Fd2Interest, of.Fd2Compound,
          of.Fd2PV, of.Fd2Result,
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
      filePath := fmt.Sprintf("%s/%s/oaepp.txt", mainDir, userName)
      if _, err := misc.WriteAllExclusiveLock(filePath, data, os.O_CREATE | os.O_RDWR, 0o660); err != nil {
        fmt.Printf("%s - %s\n", m.DTF(), err)
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
