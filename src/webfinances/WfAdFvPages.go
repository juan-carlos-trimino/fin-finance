package webfinances

import (
  "context"
  "encoding/json"
  "finance/finances"
  "finance/middlewares"
  "finance/misc"
  "finance/sessions"
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "html/template"
  "net/http"
  "os"
  "strconv"
  "strings"
)

type WfAdFvPages struct {
}

func (a WfAdFvPages) AdFvPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
    return
  }
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logger.LogInfo("Entering AdFvPages/webfinances.", correlationId)
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    af := getAdFvFields(userName)
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
    // if strings.EqualFold(p.CurrentPage, "rhs-ui1") {
      // p.CurrentButton = "lhs-button1"
      // if req.Method == http.MethodPost {
      //   p.fd1N = req.PostFormValue("fd1-n")
      //   p.fd1TimePeriod = req.PostFormValue("fd1-tp")
      //   p.fd1Interest = req.PostFormValue("fd1-interest")
      //   p.fd1Compound = req.PostFormValue("fd1-cp")
      //   p.fd1FV = req.PostFormValue("fd1-fv")
      //   var n float64
      //   var i float64
      //   var fv float64
      //   var err error
      //   if n, err = strconv.ParseFloat(p.fd1N, 64); err != nil {
      //     p.Fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1N, err)
      //   } else if i, err = strconv.ParseFloat(p.fd1Interest, 64); err != nil {
      //     p.Fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1Interest, err)
      //   } else if fv, err = strconv.ParseFloat(p.fd1FV, 64); err != nil {
      //     p.Fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1FV, err)
      //   } else {
      //     var oa finances.Annuities
      //     p.Fd1Result = fmt.Sprintf("Future Value: $%.2f", oa.O_FutureValue_PV(fv, i / 100.0,
      //                               oa.GetCompoundingPeriod(p.fd1Compound[0], true),
      //                               n, oa.GetTimePeriod(p.fd1TimePeriod[0], true)))
      //   }
      //   logEntry.Print(INFO, correlationId, []string {
      //     fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, fv = %s, %s",
      //                 p.fd1N, p.fd1TimePeriod, p.fd1Interest, p.fd1Compound, p.fd1FV, p.Fd1Result),
      //   })
      // }
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      // t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/fv/fv.html",
      //                                        "webfinances/templates/header.html",
      //                                        "webfinances/templates/ordinaryannuity/fv/n-i-PV.html",
      //                                        "webfinances/templates/footer.html"))
      // t.ExecuteTemplate(res, "oafuturevalue", struct {
      //   Header string
      //   Datetime string
      //   CurrentButton string
      //   Fd1N string
      //   Fd1TimePeriod string
      //   Fd1Interest string
      //   Fd1Compound string
      //   Fd1FV string
      //   Fd1Result string
      // } { "Ordinary Annuity / Future Value", m.DTF(), p.CurrentButton,
      //     p.fd1N, p.fd1TimePeriod, p.fd1Interest, p.fd1Compound, p.fd1FV, p.Fd1Result,
      //   })
    /*} else*/ if strings.EqualFold(af.CurrentPage, "rhs-ui2") {
      af.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        af.Fd2N = req.FormValue("fd2-n")
        af.Fd2TimePeriod = req.PostFormValue("fd2-tp")
        af.Fd2Interest = req.PostFormValue("fd2-interest")
        af.Fd2Compound = req.PostFormValue("fd2-cp")
        af.Fd2PMT = req.PostFormValue("fd2-pmt")
        var n float64
        var i float64
        var pmt float64
        var err error
        if n, err = strconv.ParseFloat(af.Fd2N, 64); err != nil {
          af.Fd2Result = fmt.Sprintf("Error: %s -- %+v", af.Fd2N, err)
        } else if i, err = strconv.ParseFloat(af.Fd2Interest, 64); err != nil {
          af.Fd2Result = fmt.Sprintf("Error: %s -- %+v", af.Fd2Interest, err)
        } else if pmt, err = strconv.ParseFloat(af.Fd2PMT, 64); err != nil {
          af.Fd2Result = fmt.Sprintf("Error: %s -- %+v", af.Fd2PMT, err)
        } else {
          var oa finances.Annuities
          af.Fd2Result = fmt.Sprintf("Future Value: $%.2f",
            oa.D_FutureValue_PMT(pmt, i / 100.0, oa.GetCompoundingPeriod(af.Fd2Compound[0], true),
            n, oa.GetTimePeriod(af.Fd2TimePeriod[0], true)))
        }
        logger.LogInfo(fmt.Sprintf("n = %s, tp = %s, interest = %s, cp = %s, pmt = %s, %s",
         af.Fd2N, af.Fd2TimePeriod, af.Fd2Interest, af.Fd2Compound, af.Fd2PMT, af.Fd2Result),
         correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/annuitydue/fv/fv.html",
        "webfinances/templates/header.html",
        "webfinances/templates/annuitydue/fv/n-i-PMT.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "adfuturevalue", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd2N string
        Fd2TimePeriod string
        Fd2Interest string
        Fd2Compound string
        Fd2PMT string
        Fd2Result string
      } { "Annuity Due / Future Value", m.DTF(), af.CurrentButton, newSession.CsrfToken,
          af.Fd2N, af.Fd2TimePeriod, af.Fd2Interest, af.Fd2Compound, af.Fd2PMT, af.Fd2Result,
        })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", af.CurrentPage)
      logger.LogError(errString, "-1")
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      logger.LogError("*** Request timeout ***", "-1")
      if strings.EqualFold(af.CurrentPage, "rhs-ui1") {
        af.Fd1Result = ""
      } else if strings.EqualFold(af.CurrentPage, "rhs-ui2") {
        af.Fd2Result = ""
      }
    }
    //
    if data, err := json.Marshal(af); err != nil {
      logger.LogError(fmt.Sprintf("%+v", err), "-1")
    } else {
      filePath := fmt.Sprintf("%s/%s/adfv.txt", mainDir, userName)
      if _, err := misc.WriteAllExclusiveLock1(filePath, data, os.O_CREATE | os.O_RDWR |
        os.O_TRUNC, 0o600); err != nil {
        logger.LogError(fmt.Sprintf("%+v", err), "-1")
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    logger.LogError(errString, "-1")
    panic(errString)
  }
}
