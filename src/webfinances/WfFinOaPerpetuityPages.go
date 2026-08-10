package webfinances

import (
  "context"
  "encoding/json"
  "finance/finances"
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/gposu"
  "github.com/juan-carlos-trimino/gpsessions"
  "net/http"
  "os"
  "strconv"
  "strings"
  "time"
)

type WfOaPerpetuityPages struct {}

func (o WfOaPerpetuityPages) OaPerpetuityPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.OaPerpetuityPages.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res, correlationId)
    return
  }
  //
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    of := getOaPerpetuityFields(userName)
    var currentRHS string = "rhs-ui1"  //Default.
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
      currentRHS = ui
    }
    //
    if strings.EqualFold(currentRHS, "rhs-ui1") {
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
          of.Fd1Result = fmt.Sprintf("Present Value of Perpetuity: $%.5f",
            oa.O_Perpetuity(i / 100.0, pmt, oa.GetCompoundingPeriod(of.Fd1Compound[0], true)))
        }
        logger.LogInfo(fmt.Sprintf("i = %s, cp = %s, pmt = %s, %s", of.Fd1Interest, of.Fd1Compound,
          of.Fd1Pmt, of.Fd1Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/ordinaryannuity/perpetuity/perpetuity.html",
        "webfinances/templates/finances/ordinaryannuity/perpetuity/p.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/navbar.html",
        "webfinances/templates/footer.html",
      }
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
        Data: struct{
          Header string
          Datetime string
          MenuPage string
          CurrentButton string
          CsrfToken string
          Fd1Interest string
          Fd1Compound string
          Fd1Pmt string
          Fd1Result string
        } { "Ordinary Annuity / Perpetuities", logger.DatetimeFormat(), financesMenuPage, of.CurrentButton,
            newSession.CsrfToken, of.Fd1Interest, of.Fd1Compound, of.Fd1Pmt, of.Fd1Result },
      })
    } else if strings.EqualFold(currentRHS, "rhs-ui2") {
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
          of.Fd2Result = fmt.Sprintf("Present Value of Perpetuity: $%.5f",
            oa.O_GrowingPerpetuity(i / 100.0, grow, pmt, oa.GetCompoundingPeriod(of.Fd2Compound[0], true)))
        }
        logger.LogInfo(fmt.Sprintf("i = %s, cp = %s, grow = %s, pmt = %s, %s", of.Fd2Interest,
          of.Fd2Compound, of.Fd2Grow, of.Fd2Pmt, of.Fd2Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/ordinaryannuity/perpetuity/perpetuity.html",
        "webfinances/templates/finances/ordinaryannuity/perpetuity/gp.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/navbar.html",
        "webfinances/templates/footer.html",
      }
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
        Data: struct{
          Header string
          Datetime string
          MenuPage string
          CurrentButton string
          CsrfToken string
          Fd2Interest string
          Fd2Compound string
          Fd2Grow string
          Fd2Pmt string
          Fd2Result string
        } { "Ordinary Annuity / Perpetuities", logger.DatetimeFormat(), financesMenuPage, of.CurrentButton,
            newSession.CsrfToken, of.Fd2Interest, of.Fd2Compound, of.Fd2Grow, of.Fd2Pmt, of.Fd2Result },
      })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", currentRHS)
      logger.LogError(errString, correlationId)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      logger.LogWarning("*** Request timeout ***", correlationId)
      if strings.EqualFold(currentRHS, "rhs-ui1") {
        of.Fd1Result = ""
      } else if strings.EqualFold(currentRHS, "rhs-ui2") {
        of.Fd2Result = ""
      }
    }
    //
    if data, err := json.Marshal(of); err != nil {
      logger.LogError(fmt.Sprintf("%+v", err), correlationId)
    } else {
      filePath := fmt.Sprintf("%s/%s/oaperpetuity.txt", mainDir, userName)
      if _, err := osu.WriteAllExclusiveLock1(filePath, data, os.O_CREATE | os.O_RDWR | os.O_TRUNC, 0o600); err != nil {
        logger.LogError(fmt.Sprintf("%+v", err), correlationId)
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    logger.LogError(errString, correlationId)
    panic(errString)
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms\n", time.Since(startTime).Microseconds()), correlationId)
}
