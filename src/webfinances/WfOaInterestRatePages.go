package webfinances

import (
  "context"
  "encoding/json"
  "finance/finances"
  "finance/sessions"
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/gpmiddlewares"
  "github.com/juan-carlos-trimino/gposu"
  "html/template"
  "net/http"
  "os"
  "strconv"
  "strings"
)

type WfOaInterestRatePages struct {
}

func (o WfOaInterestRatePages) OaInterestRatePages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
    return
  }
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logger.LogInfo("Entering OaInterestRatePages/webfinances.", correlationId)
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    of := getOaInterestRateFields(userName)
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
        of.Fd1Compound = req.PostFormValue("fd1-cp")
        of.Fd1PV = req.PostFormValue("fd1-pv")
        of.Fd1FV = req.PostFormValue("fd1-fv")
        var n float64
        var pv float64
        var fv float64
        var err error
        if n, err = strconv.ParseFloat(of.Fd1N, 64); err != nil {
          of.Fd1Result = fmt.Sprintf("Error: %s -- %+v", of.Fd1N, err)
        } else if pv, err = strconv.ParseFloat(of.Fd1PV, 64); err != nil {
          of.Fd1Result = fmt.Sprintf("Error: %s -- %+v", of.Fd1PV, err)
        } else if fv, err = strconv.ParseFloat(of.Fd1FV, 64); err != nil {
          of.Fd1Result = fmt.Sprintf("Error: %s -- %+v", of.Fd1FV, err)
        } else {
          var oa finances.Annuities
          of.Fd1Result = fmt.Sprintf("Interest: %.3f%% %s", oa.O_Interest_PV_FV(pv, fv, n,
            oa.GetTimePeriod(of.Fd1TimePeriod[0], true),
            oa.GetCompoundingPeriod(of.Fd1Compound[0], true)) * 100.0, of.Fd1Compound)
        }
        logger.LogInfo(fmt.Sprintf("n = %s, tp = %s, cp = %s, pv = %s, fv = %s, %s", of.Fd1N,
         of.Fd1TimePeriod, of.Fd1Compound, of.Fd1PV, of.Fd1FV, of.Fd1Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/interestrate/interestrate.html",
        "webfinances/templates/header.html",
        "webfinances/templates/ordinaryannuity/interestrate/n-PV-FV.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "oainterestrate", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd1N string
        Fd1TimePeriod string
        Fd1Compound string
        Fd1PV string
        Fd1FV string
        Fd1Result string
      } { "Ordinary Annuity / Interest Rate", logger.DatetimeFormat(), of.CurrentButton,
          newSession.CsrfToken, of.Fd1N, of.Fd1TimePeriod, of.Fd1Compound, of.Fd1PV, of.Fd1FV,
          of.Fd1Result,
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
      }
    }
    //
    if data, err := json.Marshal(of); err != nil {
      logger.LogError(fmt.Sprintf("%+v", err), "-1")
    } else {
      filePath := fmt.Sprintf("%s/%s/oainterestrate.txt", mainDir, userName)
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
