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

type oaInterestRateFields struct {
  MenuPage string `json:"menuPage"`
  CurrentPage string `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1N string `json:"fd1N"`
  Fd1TimePeriod string `json:"fd1TimePeriod"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1PV string `json:"fd1PV"`
  Fd1FV string `json:"fd1FV"`
  Fd1Result string `json:"fd1Result"`
}

func newOaInterestRateFields(dir1, dir2, correlationId string) *oaInterestRateFields {
  //Default values returned if file is missing, empty, or JSON is corrupt.
  defaults := oaInterestRateFields{
    MenuPage: "",
    CurrentPage: "rhs-ui1",
    CurrentButton: "lhs-button1",
    //
    Fd1N: "1.0",
    Fd1TimePeriod: "year",
    Fd1Compound: "monthly",
    Fd1PV: "1.00",
    Fd1FV: "1.07",
    Fd1Result: "",
  }
  return loadFieldsFromDisk(dir1, dir2, "oainterestrate.txt", correlationId, &defaults)
}

func getOaInterestRateFields(userName string) *oaInterestRateFields {
  return getFieldsGeneric(userName, func(s *UserSession) *oaInterestRateFields {
    return s.Fields.oaInterestRate
  })
}

type WfOaInterestRatePages struct{}

func (o WfOaInterestRatePages) OaInterestRatePages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.OaInterestRatePages.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res, correlationId)
    return
  }
  //
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    fields := getOaInterestRateFields(userName)
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
        fields.Fd1TimePeriod = req.PostFormValue("fd1-tp")
        fields.Fd1Compound = req.PostFormValue("fd1-cp")
        fields.Fd1PV = req.PostFormValue("fd1-pv")
        fields.Fd1FV = req.PostFormValue("fd1-fv")
        var n float64
        var pv float64
        var fv float64
        var err error
        if n, err = strconv.ParseFloat(fields.Fd1N, 64); err != nil {
          fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1N, err)
        } else if pv, err = strconv.ParseFloat(fields.Fd1PV, 64); err != nil {
          fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1PV, err)
        } else if fv, err = strconv.ParseFloat(fields.Fd1FV, 64); err != nil {
          fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1FV, err)
        } else {
          var oa finances.Annuities
          fields.Fd1Result = fmt.Sprintf("Interest: %.5f%% %s", oa.O_Interest_PV_FV(pv, fv, n,
            oa.GetTimePeriod(fields.Fd1TimePeriod[0], true),
            oa.GetCompoundingPeriod(fields.Fd1Compound[0], true)) * 100.0, fields.Fd1Compound)
        }
        logger.LogInfo(fmt.Sprintf("n = %s, tp = %s, cp = %s, pv = %s, fv = %s, %s", fields.Fd1N, fields.Fd1TimePeriod,
          fields.Fd1Compound, fields.Fd1PV, fields.Fd1FV, fields.Fd1Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/ordinaryannuity/interestrate/interestrate.html",
        "webfinances/templates/finances/ordinaryannuity/interestrate/n-PV-FV.html",
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
          Fd1TimePeriod string
          Fd1Compound string
          Fd1PV string
          Fd1FV string
          Fd1Result string
        } { "standard", "Ordinary Annuity / Interest Rate", logger.DatetimeFormat(), financesMenuPage, fields.CurrentButton,
            newSession.CsrfToken, fields.Fd1N, fields.Fd1TimePeriod, fields.Fd1Compound, fields.Fd1PV, fields.Fd1FV, fields.Fd1Result },
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
      }
    }
    //
    if data, err := json.Marshal(fields); err != nil {
      logger.LogError(fmt.Sprintf("%+v", err), correlationId)
    } else {
      filePath := fmt.Sprintf("%s/%s/oainterestrate.txt", mainDir, userName)
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
