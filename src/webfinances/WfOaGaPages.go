package webfinances

import (
  "context"
  "encoding/json"
  "finance/finances"
  "finance/middlewares"
  "finance/sessions"
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/gposu"
  "html/template"
  "net/http"
  "os"
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
  logger.LogInfo("Entering OaGaPages/webfinances.", correlationId)
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    of := getOaGaFields(userName)
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
        of.Fd1Interest = req.PostFormValue("fd1-interest")
        of.Fd1Compound = req.PostFormValue("fd1-cp")
        of.Fd1Grow = req.PostFormValue("fd1-grow")
        of.Fd1Pmt = req.PostFormValue("fd1-pmt")
        var n float64
        var i float64
        var grow float64
        var pmt float64
        var err error
        if n, err = strconv.ParseFloat(of.Fd1N, 64); err != nil {
          of.Fd1Result = fmt.Sprintf("Error: %s -- %+v", of.Fd1N, err)
        } else if i, err = strconv.ParseFloat(of.Fd1Interest, 64); err != nil {
          of.Fd1Result = fmt.Sprintf("Error: %s -- %+v", of.Fd1Interest, err)
        } else if grow, err = strconv.ParseFloat(of.Fd1Grow, 64); err != nil {
          of.Fd1Result = fmt.Sprintf("Error: %s -- %+v", of.Fd1Grow, err)
        } else if pmt, err = strconv.ParseFloat(of.Fd1Pmt, 64); err != nil {
          of.Fd1Result = fmt.Sprintf("Error: %s -- %+v", of.Fd1Pmt, err)
        } else {
          var oa finances.Annuities
          of.Fd1Result = fmt.Sprintf("Future Value: $%.2f", oa.O_GrowingAnnuityFutureValue(pmt, n,
            grow, i / 100.0, oa.GetCompoundingPeriod(of.Fd1Compound[0], true)))
        }
        logger.LogInfo(fmt.Sprintf("n = %s, i = %s, cp = %s, grow = %s, pmt = %s, %s", of.Fd1N,
         of.Fd1Interest, of.Fd1Compound, of.Fd1Grow, of.Fd1Pmt, of.Fd1Result), correlationId)
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
      } { "Ordinary Annuity / Growing Annuity", logger.DatetimeFormat(), of.CurrentButton,
          newSession.CsrfToken, of.Fd1N, of.Fd1Interest, of.Fd1Compound, of.Fd1Grow,
          of.Fd1Pmt, of.Fd1Result,
        })
    } else if strings.EqualFold(of.CurrentPage, "rhs-ui2") {
      of.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        of.Fd2N = req.PostFormValue("fd2-n")
        of.Fd2Interest = req.PostFormValue("fd2-interest")
        of.Fd2Compound = req.PostFormValue("fd2-cp")
        of.Fd2Grow = req.PostFormValue("fd2-grow")
        of.Fd2Pmt = req.PostFormValue("fd2-pmt")
        var n float64
        var i float64
        var grow float64
        var pmt float64
        var err error
        if n, err = strconv.ParseFloat(of.Fd2N, 64); err != nil {
          of.Fd2Result = fmt.Sprintf("Error: %s -- %+v", of.Fd2N, err)
        } else if i, err = strconv.ParseFloat(of.Fd2Interest, 64); err != nil {
          of.Fd2Result = fmt.Sprintf("Error: %s -- %+v", of.Fd2Interest, err)
        } else if grow, err = strconv.ParseFloat(of.Fd2Grow, 64); err != nil {
          of.Fd2Result = fmt.Sprintf("Error: %s -- %+v", of.Fd2Grow, err)
        } else if pmt, err = strconv.ParseFloat(of.Fd2Pmt, 64); err != nil {
          of.Fd2Result = fmt.Sprintf("Error: %s -- %+v", of.Fd2Pmt, err)
        } else {
          var oa finances.Annuities
          of.Fd2Result = fmt.Sprintf("Present Value: $%.2f", oa.O_GrowingAnnuityPresentValue(pmt,
            n, grow, i / 100.0, oa.GetCompoundingPeriod(of.Fd2Compound[0], true)))
        }
        logger.LogInfo(fmt.Sprintf("n = %s, i = %s, cp = %s, grow = %s, pmt = %s, %s", of.Fd2N,
         of.Fd2Interest, of.Fd2Compound, of.Fd2Grow, of.Fd2Pmt, of.Fd2Result), correlationId)
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
      } { "Ordinary Annuity / Growing Annuity", logger.DatetimeFormat(), of.CurrentButton,
          newSession.CsrfToken, of.Fd2N, of.Fd2Interest, of.Fd2Compound, of.Fd2Grow, of.Fd2Pmt,
          of.Fd2Result,
        })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", of.CurrentPage)
      logger.LogError(errString, "-1")
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      logger.LogWarning("*** Request timeout ***", "-1")
      if strings.EqualFold(of.CurrentPage, "rhs-ui1") {
        of.Fd1Result = ""
      } else if strings.EqualFold(of.CurrentPage, "rhs-ui2") {
        of.Fd2Result = ""
      }
    }
    //
    if data, err := json.Marshal(of); err != nil {
      logger.LogError(fmt.Sprintf("%+v", err), "-1")
    } else {
      filePath := fmt.Sprintf("%s/%s/oaga.txt", mainDir, userName)
      if _, err := osu.WriteAllExclusiveLock1(filePath, data, os.O_CREATE | os.O_RDWR |
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
