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

type oaGaFields struct {
  MenuPage string `json:"menuPage"`
  CurrentPage string `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1N string `json:"fd1N"`
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1Grow string `json:"fd1Grow"`
  Fd1Pmt string `json:"fd1Pmt"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2N string `json:"fd2N"`
  Fd2Interest string `json:"fd2Interest"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2Grow string `json:"fd2Grow"`
  Fd2Pmt string `json:"fd2Pmt"`
  Fd2Result string `json:"fd2Result"`
}

func newOaGaFields(dir1, dir2, correlationId string) *oaGaFields {
  //Default values returned if file is missing, empty, or JSON is corrupt.
  defaults := oaGaFields{
    MenuPage: "",
    CurrentPage: "rhs-ui1",
    CurrentButton: "lhs-button1",
    //
    Fd1N: "1.00",
    Fd1Interest: "1.00",
    Fd1Compound: "annually",
    Fd1Grow: "1.00",
    Fd1Pmt: "1.00",
    Fd1Result: "",
    //
    Fd2N: "1.00",
    Fd2Interest: "1.00",
    Fd2Compound: "annually",
    Fd2Grow: "1.00",
    Fd2Pmt: "1.00",
    Fd2Result: "",
  }
  return loadFieldsFromDisk(dir1, dir2, "oaga.txt", correlationId, &defaults)
}

func getOaGaFields(userName string) *oaGaFields {
  return getFieldsGeneric(userName, func(s *UserSession) *oaGaFields {
    return s.Fields.oaGa
  })
}

type WfOaGaPages struct{}

func (o WfOaGaPages) OaGaPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.OaGaPages.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res, correlationId)
    return
  }
  //
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    fields := getOaGaFields(userName)
    /***
    The functions in Request that allow to extract data from the URL and/or the body revolve around the Form, PostForm, and
    MultipartForm fields; the data are in the form of key-value pairs.

    If the form and the URL have the same key name, both of them will be placed in a slice, with the form value always prioritized
    before the URL value.

    Since we want the form key-value pairs, we can ignore the URL key-value pairs. The PostForm field provides key-value pairs only
    for the form and not the URL. The PostForm field supports only application/x-www-form-urlencoded.

    The FormValue method lets you access the key-value pairs just like the Form field, except that it's for a specific key and there
    is no need to call the ParseForm method beforehand -- the FormValue method does it. The PostFormValue method does the same thing,
    except that it's for the PostForm field instead of the Form field.
    ***/
    if ui := req.FormValue("compute"); ui != "" {  //Values from form and URL.
      fields.CurrentPage = ui
    }
    //
    if strings.EqualFold(fields.CurrentPage, "rhs-ui1") {
      fields.CurrentButton = "lhs-button1"
      if req.Method == http.MethodPost {
        fields.Fd1N = req.PostFormValue("fd1-n")
        fields.Fd1Interest = req.PostFormValue("fd1-interest")
        fields.Fd1Compound = req.PostFormValue("fd1-cp")
        fields.Fd1Grow = req.PostFormValue("fd1-grow")
        fields.Fd1Pmt = req.PostFormValue("fd1-pmt")
        var n float64
        var i float64
        var grow float64
        var pmt float64
        var err error
        if n, err = strconv.ParseFloat(fields.Fd1N, 64); err != nil {
          fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1N, err)
        } else if i, err = strconv.ParseFloat(fields.Fd1Interest, 64); err != nil {
          fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1Interest, err)
        } else if grow, err = strconv.ParseFloat(fields.Fd1Grow, 64); err != nil {
          fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1Grow, err)
        } else if pmt, err = strconv.ParseFloat(fields.Fd1Pmt, 64); err != nil {
          fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1Pmt, err)
        } else {
          var oa finances.Annuities
          fields.Fd1Result = fmt.Sprintf("Future Value: $%.5f",
            oa.O_GrowingAnnuityFutureValue(pmt, n, grow, i / 100.0, oa.GetCompoundingPeriod(fields.Fd1Compound[0], true)))
        }
        logger.LogInfo(fmt.Sprintf("n = %s, i = %s, cp = %s, grow = %s, pmt = %s, %s", fields.Fd1N, fields.Fd1Interest, fields.Fd1Compound,
          fields.Fd1Grow, fields.Fd1Pmt, fields.Fd1Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/ordinaryannuity/ga/ga.html",
        "webfinances/templates/finances/ordinaryannuity/ga/FV.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/navbar.html",
        "webfinances/templates/footer.html",
      }
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
        Data: struct{
          LayoutType string
          Header string
          Datetime string
          MenuPage string
          CurrentButton string
          CsrfToken string
          Fd1N string
          Fd1Interest string
          Fd1Compound string
          Fd1Grow string
          Fd1Pmt string
          Fd1Result string
        } { "standard", "Ordinary Annuity / Growing Annuity", logger.DatetimeFormat(), financesMenuPage, fields.CurrentButton,
            newSession.CsrfToken, fields.Fd1N, fields.Fd1Interest, fields.Fd1Compound, fields.Fd1Grow, fields.Fd1Pmt, fields.Fd1Result },
      })
    } else if strings.EqualFold(fields.CurrentPage, "rhs-ui2") {
      fields.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        fields.Fd2N = req.PostFormValue("fd2-n")
        fields.Fd2Interest = req.PostFormValue("fd2-interest")
        fields.Fd2Compound = req.PostFormValue("fd2-cp")
        fields.Fd2Grow = req.PostFormValue("fd2-grow")
        fields.Fd2Pmt = req.PostFormValue("fd2-pmt")
        var n float64
        var i float64
        var grow float64
        var pmt float64
        var err error
        if n, err = strconv.ParseFloat(fields.Fd2N, 64); err != nil {
          fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2N, err)
        } else if i, err = strconv.ParseFloat(fields.Fd2Interest, 64); err != nil {
          fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Interest, err)
        } else if grow, err = strconv.ParseFloat(fields.Fd2Grow, 64); err != nil {
          fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Grow, err)
        } else if pmt, err = strconv.ParseFloat(fields.Fd2Pmt, 64); err != nil {
          fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Pmt, err)
        } else {
          var oa finances.Annuities
          fields.Fd2Result = fmt.Sprintf("Present Value: $%.5f",
            oa.O_GrowingAnnuityPresentValue(pmt, n, grow, i / 100.0, oa.GetCompoundingPeriod(fields.Fd2Compound[0], true)))
        }
        logger.LogInfo(fmt.Sprintf("n = %s, i = %s, cp = %s, grow = %s, pmt = %s, %s", fields.Fd2N, fields.Fd2Interest, fields.Fd2Compound,
          fields.Fd2Grow, fields.Fd2Pmt, fields.Fd2Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/ordinaryannuity/ga/ga.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/navbar.html",
        "webfinances/templates/finances/ordinaryannuity/ga/PV.html",
        "webfinances/templates/footer.html",
      }
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
        Data: struct{
          LayoutType string
          Header string
          Datetime string
          MenuPage string
          CurrentButton string
          CsrfToken string
          Fd2N string
          Fd2Interest string
          Fd2Compound string
          Fd2Grow string
          Fd2Pmt string
          Fd2Result string
        } { "standard", "Ordinary Annuity / Growing Annuity", logger.DatetimeFormat(), financesMenuPage, fields.CurrentButton,
            newSession.CsrfToken, fields.Fd2N, fields.Fd2Interest, fields.Fd2Compound, fields.Fd2Grow, fields.Fd2Pmt, fields.Fd2Result },
      })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", fields.CurrentPage)
      logger.LogError(errString, correlationId)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      logger.LogWarning("*** Request timeout ***", correlationId)
      if strings.EqualFold(fields.CurrentPage, "rhs-ui1") {
        fields.Fd1Result = ""
      } else if strings.EqualFold(fields.CurrentPage, "rhs-ui2") {
        fields.Fd2Result = ""
      }
    }
    //
    if data, err := json.Marshal(fields); err != nil {
      logger.LogError(fmt.Sprintf("%+v", err), correlationId)
    } else {
      filePath := fmt.Sprintf("%s/%s/oaga.txt", mainDir, userName)
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
