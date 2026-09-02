package webfinances

import (
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

type adFvFields struct {
  MenuPage string `json:"menuPage"`
  CurrentPage string `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1N string `json:"fd1N"`
  Fd1TimePeriod string `json:"fd1TimePeriod"`
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1FV string `json:"fd1FV"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2N string `json:"fd2N"`
  Fd2TimePeriod string `json:"fd2TimePeriod"`
  Fd2Interest string `json:"fd2Interest"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2PMT string `json:"fd2PMT"`
  Fd2Result string `json:"fd2Result"`
}

func newAdFvFields(dir1, dir2, correlationId string) *adFvFields {
  //Default values returned if file is missing, empty, or JSON is corrupt.
  defaults := adFvFields{
    MenuPage: "",
    CurrentButton: "lhs-button2",
    CurrentPage: "rhs-ui2",
    //
    Fd1N: "1.0",
    Fd1TimePeriod: "year",
    Fd1Interest: "1.00",
    Fd1Compound: "monthly",
    Fd1FV: "1.00",
    Fd1Result: "",
    //
    Fd2N: "1.0",
    Fd2TimePeriod: "year",
    Fd2Interest: "1.00",
    Fd2Compound: "monthly",
    Fd2PMT: "1.00",
    Fd2Result: "",
  }
  return loadFieldsFromDisk(dir1, dir2, "adfv.txt", correlationId, &defaults)
}

func getAdFvFields(userName string) *adFvFields {
  return getFieldsGeneric(userName, func(s *UserSession) *adFvFields {
    return s.Fields.adFv
  })
}

type WfAdFvPages struct {}

func (a WfAdFvPages) AdFvPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.AdFvPages.", correlationId)
  //Guard Clause 1: Validate HTTP Method.
  if req.Method != http.MethodPost && req.Method != http.MethodGet {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    logger.LogError(errString, correlationId)
    panic(errString)
  }
  //Guard Clause 2: Validate Session Token.
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res, correlationId)
    return
  }
  userName := sessions.GetUserName(sessionToken)
  fields := getAdFvFields(userName)
  //Every time a web request processes data for a user, update the timestamp under a lock.
  currentFieldsLock.Lock()
  if session, exists := currentFields[userName]; exists {
    session.LastAccessed = time.Now()
  }
  currentFieldsLock.Unlock()
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
  //Dynamic variables determined by the routing route condition.
  var partialTemplate string
  var templateData interface{}
  switch strings.ToLower(fields.CurrentPage) {
  case "rhs-ui2":
    fields.CurrentButton = "lhs-button2"
    if req.Method == http.MethodPost {
      a.processUi2Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "n-i-PMT.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd2N string
      Fd2TimePeriod string
      Fd2Interest string
      Fd2Compound string
      Fd2PMT string
      Fd2Result string
    }{
      "standard",
      "Annuity Due / Future Value",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd2N,
      fields.Fd2TimePeriod,
      fields.Fd2Interest,
      fields.Fd2Compound,
      fields.Fd2PMT,
      fields.Fd2Result,
    }
  default:
    errString := fmt.Sprintf("Unsupported page: %s", fields.CurrentPage)
    logger.LogError(errString, correlationId)
    panic(errString)
  }
  //Unified execution of templates.
  templatesNeeded := []string{
    "webfinances/templates/layout.html",
    "webfinances/templates/finances/annuitydue/fv/fv.html",
    "webfinances/templates/finances/annuitydue/fv/" + partialTemplate,
    "webfinances/templates/title.html",
    "webfinances/templates/datetime.html",
    "webfinances/templates/navbar.html",
    "webfinances/templates/footer.html",
  }
  renderer.Render(res, "layout", templatesNeeded, renderer.PageData{Data: templateData})
  data, err := json.Marshal(fields) //Preserve current choices.
  if err != nil {
    //Don't crash the server (panic), but log it clearly so you can debug the serialization.
    logger.LogError(fmt.Sprintf("Failed to marshal fields to JSON for user %s: %+v", userName, err), correlationId)
  } else {
    go func(userData []byte, uName, cId string) {
      filePath := fmt.Sprintf("%s/%s/adfv.txt", mainDir, uName)
      //The exclusive OS file lock handles goroutine collisions.
      _, err := osu.WriteAllExclusiveLock1(filePath, userData, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o600)
      //Check the error returned from the lock-writing function.
      if err != nil {
        logger.LogError(fmt.Sprintf("Goroutine file system error writing state to %s: %+v", filePath, err), cId)
        return
      }
      logger.LogInfo(fmt.Sprintf("Goroutine successfully persisted state to %s.", filePath), cId)
    }(data, userName, correlationId) // Pass variables into the closure to prevent scope races
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms\n", time.Since(startTime).Microseconds()), correlationId)
}

//Extraction helper for UI-2 calculations.
func (a WfAdFvPages) processUi2Form(req *http.Request, fields *adFvFields, correlationId string) {
  fields.Fd2N = req.FormValue("fd2-n")
  fields.Fd2TimePeriod = req.PostFormValue("fd2-tp")
  fields.Fd2Interest = req.PostFormValue("fd2-interest")
  fields.Fd2Compound = req.PostFormValue("fd2-cp")
  fields.Fd2PMT = req.PostFormValue("fd2-pmt")
  var n float64
  var i float64
  var pmt float64
  var err error
  if n, err = strconv.ParseFloat(fields.Fd2N, 64); err != nil {
    fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2N, err)
  } else if i, err = strconv.ParseFloat(fields.Fd2Interest, 64); err != nil {
    fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Interest, err)
  } else if pmt, err = strconv.ParseFloat(fields.Fd2PMT, 64); err != nil {
    fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2PMT, err)
  } else {
    var oa finances.Annuities
    fields.Fd2Result = fmt.Sprintf("Future Value: $%.5f",
      oa.D_FutureValue_PMT(pmt, i / 100.0, oa.GetCompoundingPeriod(fields.Fd2Compound[0], true), n,
      oa.GetTimePeriod(fields.Fd2TimePeriod[0], true)))
  }
  logger.LogInfo(fmt.Sprintf("n = %s, tp = %s, interest = %s, cp = %s, pmt = %s, %s", fields.Fd2N, fields.Fd2TimePeriod,
    fields.Fd2Interest, fields.Fd2Compound, fields.Fd2PMT, fields.Fd2Result), correlationId)
}
