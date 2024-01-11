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

type WfAdEppPages struct {
}

func (a WfAdEppPages) AdEppPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
    return
  }
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering AdEppPages/webfinances.",
  })
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    af := getAdEppFields(userName)
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
      af.CurrentPage = ui
    }
    //
    if strings.EqualFold(af.CurrentPage, "rhs-ui1") {
      af.CurrentButton = "lhs-button1"
      if req.Method == http.MethodPost {
        af.Fd1N = req.PostFormValue("fd1-n")
        af.Fd1TimePeriod = req.PostFormValue("fd1-tp")
        af.Fd1Interest = req.PostFormValue("fd1-interest")
        af.Fd1Compound = req.PostFormValue("fd1-cp")
        af.Fd1FV = req.PostFormValue("fd1-fv")
        var n float64
        var i float64
        var fv float64
        var err error
        if n, err = strconv.ParseFloat(af.Fd1N, 64); err != nil {
          af.Fd1Result = fmt.Sprintf("Error: %s -- %+v", af.Fd1N, err)
        } else if i, err = strconv.ParseFloat(af.Fd1Interest, 64); err != nil {
          af.Fd1Result = fmt.Sprintf("Error: %s -- %+v", af.Fd1Interest, err)
        } else if fv, err = strconv.ParseFloat(af.Fd1FV, 64); err != nil {
          af.Fd1Result = fmt.Sprintf("Error: %s -- %+v", af.Fd1FV, err)
        } else {
          var oa finances.Annuities
          af.Fd1Result = fmt.Sprintf("Payment: $%.2f", oa.D_Payment_FV(fv, i / 100.0,
            oa.GetCompoundingPeriod(af.Fd1Compound[0], true), n,
            oa.GetTimePeriod(af.Fd1TimePeriod[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, fv = %s, %s", af.Fd1N,
            af.Fd1TimePeriod, af.Fd1Interest, af.Fd1Compound, af.Fd1FV, af.Fd1Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      t := template.Must(template.ParseFiles("webfinances/templates/annuitydue/epp/epp.html",
        "webfinances/templates/header.html",
        "webfinances/templates/annuitydue/epp/n-i-FV.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "adequalperiodicpayments", struct {
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
      } { "Annuity Due / Equal Periodic Payments", m.DTF(), af.CurrentButton, newSession.CsrfToken,
          af.Fd1N, af.Fd1TimePeriod, af.Fd1Interest, af.Fd1Compound, af.Fd1FV, af.Fd1Result,
        })
    } else if strings.EqualFold(af.CurrentPage, "rhs-ui2") {
      af.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        af.Fd2N = req.PostFormValue("fd2-n")
        af.Fd2TimePeriod = req.PostFormValue("fd2-tp")
        af.Fd2Interest = req.FormValue("fd2-interest")
        af.Fd2Compound = req.PostFormValue("fd2-cp")
        af.Fd2PV = req.PostFormValue("fd2-pv")
        var n float64
        var i float64
        var pv float64
        var err error
        if n, err = strconv.ParseFloat(af.Fd2N, 64); err != nil {
          af.Fd2Result = fmt.Sprintf("Error: %s -- %+v", af.Fd2N, err)
        } else if i, err = strconv.ParseFloat(af.Fd2Interest, 64); err != nil {
          af.Fd2Result = fmt.Sprintf("Error: %s -- %+v", af.Fd2Interest, err)
        } else if pv, err = strconv.ParseFloat(af.Fd2PV, 64); err != nil {
          af.Fd2Result = fmt.Sprintf("Error: %s -- %+v", af.Fd2PV, err)
        } else {
          var oa finances.Annuities
          af.Fd2Result = fmt.Sprintf("Payment: $%.2f", oa.D_Payment_PV(pv, i / 100.0,
            oa.GetCompoundingPeriod(af.Fd2Compound[0], true), n,
            oa.GetTimePeriod(af.Fd2TimePeriod[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, pv = %s, %s", af.Fd2N,
            af.Fd2TimePeriod, af.Fd2Interest, af.Fd2Compound, af.Fd2PV, af.Fd2Result),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/annuitydue/epp/epp.html",
        "webfinances/templates/header.html",
        "webfinances/templates/annuitydue/epp/n-i-PV.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "adequalperiodicpayments", struct {
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
      } { "Annuity Due / Equal Periodic Payments", m.DTF(), af.CurrentButton, newSession.CsrfToken,
          af.Fd2N, af.Fd2TimePeriod, af.Fd2Interest, af.Fd2Compound, af.Fd2PV, af.Fd2Result,
        })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", af.CurrentPage)
      fmt.Printf("%s - %s\n", m.DTF(), errString)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      fmt.Println("*** Request timeout ***")
      if strings.EqualFold(af.CurrentPage, "rhs-ui1") {
        af.Fd1Result = ""
      } else if strings.EqualFold(af.CurrentPage, "rhs-ui2") {
        af.Fd2Result = ""
      }
    }
    //
    if data, err := json.Marshal(af); err != nil {
      fmt.Printf("%s - %s\n", m.DTF(), err)
    } else {
      filePath := fmt.Sprintf("%s/%s/adepp.txt", mainDir, userName)
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
