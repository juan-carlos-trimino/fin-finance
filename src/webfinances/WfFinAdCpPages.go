package webfinances

import (
  "encoding/json"
  "finance/finances"
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/gposu"
  "github.com/juan-carlos-trimino/gpsessions"
  "net/http"
  "os"
  "strconv"
  "strings"
  "time"
)

type adCpFields struct {
  MenuPage string `json:"menuPage"`
  CurrentPage string `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
  // UI 1
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1PV string `json:"fd1PV"`
  Fd1FV string `json:"fd1FV"`
  Fd1Result string `json:"fd1Result"`
  // UI 2
  Fd2Interest string `json:"fd2Interest"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2Payment string `json:"fd2Payment"`
  Fd2PV string `json:"fd2PV"`
  Fd2Result string `json:"fd2Result"`
  // UI 3
  Fd3Interest string `json:"fd3Interest"`
  Fd3Compound string `json:"fd3Compound"`
  Fd3Payment string `json:"fd3Payment"`
  Fd3FV string `json:"fd3FV"`
  Fd3Result string `json:"fd3Result"`
}

func newAdCpFields(dir1, dir2, correlationId string) *adCpFields {
  //Default values returned if file is missing, empty, or JSON is corrupt.
  defaults := adCpFields{
    MenuPage: "",
    CurrentPage: "rhs-ui2",
    CurrentButton: "lhs-button2",
    //
    Fd1Interest: "1.00",
    Fd1Compound: "annually",
    Fd1PV: "1.00",
    Fd1FV: "1.00",
    Fd1Result: "",
    //
    Fd2Interest: "1.00",
    Fd2Compound: "annually",
    Fd2Payment: "1.00",
    Fd2PV: "1.00",
    Fd2Result: "",
    //
    Fd3Interest: "1.00",
    Fd3Compound: "annually",
    Fd3Payment: "1.00",
    Fd3FV: "1.00",
    Fd3Result: "",
  }
  return loadFieldsFromDisk(dir1, dir2, "adcp.txt", correlationId, &defaults)
}

func getAdCpFields(userName string) *adCpFields {
  return getFieldsGeneric(userName, func(s *UserSession) *adCpFields {
    return s.Fields.adCp
  })
}

type WfAdCpPages struct {}

func (a WfAdCpPages) AdCpPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.AdCpPages.", correlationId)
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
  fields := getAdCpFields(userName)
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
  if ui := req.FormValue("compute"); ui != "" {
    fields.CurrentPage = ui
  }
  //Dynamic variables determined by the routing route condition.
  var partialTemplate string
  var templateData interface{}
  //Core application state splitting via explicit page mapping.
  switch strings.ToLower(fields.CurrentPage) {
  case "rhs-ui2":
    fields.CurrentButton = "lhs-button2"
    if req.Method == http.MethodPost {
      a.processUi2Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "i-PMT-PV.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd2Interest string
      Fd2Compound string
      Fd2Payment string
      Fd2PV string
      Fd2Result string
    }{
      "standard",
      "Annuity Due / Compounding Periods",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd2Interest,
      fields.Fd2Compound,
      fields.Fd2Payment,
      fields.Fd2PV,
      fields.Fd2Result,
    }
  case "rhs-ui3":
    fields.CurrentButton = "lhs-button3"
    if req.Method == http.MethodPost {
      a.processUi3Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "i-PMT-FV.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd3Interest string
      Fd3Compound string
      Fd3Payment string
      Fd3FV string
      Fd3Result string
    }{
      "standard",
      "Annuity Due / Compounding Periods",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd3Interest,
      fields.Fd3Compound,
      fields.Fd3Payment,
      fields.Fd3FV,
      fields.Fd3Result,
    }
  default:
    errString := fmt.Sprintf("Unsupported page: %s", fields.CurrentPage)
    logger.LogError(errString, correlationId)
    panic(errString)
  }
  //Unified execution of templates.
  templatesNeeded := []string{
    "webfinances/templates/layout.html",
    "webfinances/templates/finances/annuitydue/cp/cp.html",
    "webfinances/templates/finances/annuitydue/cp/" + partialTemplate,
    "webfinances/templates/title.html",
    "webfinances/templates/datetime.html",
    "webfinances/templates/navbar.html",
    "webfinances/templates/footer.html",
  }
  renderer.Render(res, "layout", templatesNeeded, renderer.PageData{Data: templateData})
  data, err := json.Marshal(fields)  //Preserve current choices.
  if err != nil {
    //Don't crash the server (panic), but log it clearly so you can debug the serialization.
    logger.LogError(fmt.Sprintf("Failed to marshal fields to JSON for user %s: %+v", userName, err), correlationId)
  } else {
    /***
    Pass arguments explicitly into the goroutine closure: Notice the func(userData []byte, ...) signature followed by (data,
    userName, correlationId) at the end. In Go, passing variables into a goroutine as parameters ensures they are copied,
    preventing the background thread from accidentally reading shifting pointers from the main routine.

    Never use the fields pointer inside the goroutine: The fields object must never cross the boundary into the go func block.
    Only the finalized, immutable data byte slice should be touched by the background thread.

    Your OS lock is still required: Even in a background thread, osu.WriteAllExclusiveLock1 is vital because it stops multiple
    background goroutines from corrupting the same physical file on disk.
    ***/
    go func(userData []byte, uName, cId string) {
      filePath := fmt.Sprintf("%s/%s/adcp.txt", mainDir, uName)
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
func (a WfAdCpPages) processUi2Form(req *http.Request, fields *adCpFields, correlationId string) {
  fields.Fd2Interest = req.FormValue("fd2-interest")
  fields.Fd2Compound = req.PostFormValue("fd2-cp")
  fields.Fd2Payment = req.PostFormValue("fd2-payment")
  fields.Fd2PV = req.PostFormValue("fd2-pv")
  var i float64
  var pmt float64
  var pv float64
  var err error
  if i, err = strconv.ParseFloat(fields.Fd2Interest, 64); err != nil {
    fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Interest, err)
  } else if pmt, err = strconv.ParseFloat(fields.Fd2Payment, 64); err != nil {
    fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Payment, err)
  } else if pv, err = strconv.ParseFloat(fields.Fd2PV, 64); err != nil {
    fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2PV, err)
  } else {
    var oa finances.Annuities
    fields.Fd2Result = fmt.Sprintf("Compounding Period: %.5f %s", oa.D_Periods_PMT_PV(pmt, pv, i / 100.0,
    oa.GetCompoundingPeriod(fields.Fd2Compound[0], true)), oa.TimePeriods(fields.Fd2Compound))
  }
  logger.LogInfo(fmt.Sprintf("i = %s, cp = %s, pmt = %s, pv = %s, %s", fields.Fd2Interest, fields.Fd2Compound,
    fields.Fd2Payment, fields.Fd2PV, fields.Fd2Result), correlationId)
}

//Extraction helper for UI-3 calculations.
func (a WfAdCpPages) processUi3Form(req *http.Request, fields *adCpFields, correlationId string) {
  fields.Fd3Interest = req.FormValue("fd3-interest")
  fields.Fd3Compound = req.PostFormValue("fd3-cp")
  fields.Fd3Payment = req.PostFormValue("fd3-payment")
  fields.Fd3FV = req.PostFormValue("fd3-fv")
  var i float64
  var pmt float64
  var fv float64
  var err error
  if i, err = strconv.ParseFloat(fields.Fd3Interest, 64); err != nil {
    fields.Fd3Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd3Interest, err)
  } else if pmt, err = strconv.ParseFloat(fields.Fd3Payment, 64); err != nil {
    fields.Fd3Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd3Payment, err)
  } else if fv, err = strconv.ParseFloat(fields.Fd3FV, 64); err != nil {
    fields.Fd3Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd3FV, err)
  } else {
    var oa finances.Annuities
    fields.Fd3Result = fmt.Sprintf("Compounding Period: %.5f %s", oa.D_Periods_PMT_FV(pmt, fv, i / 100.0,
      oa.GetCompoundingPeriod(fields.Fd3Compound[0], true)), oa.TimePeriods(fields.Fd3Compound))
  }
  logger.LogInfo(fmt.Sprintf("i = %s, cp = %s, pmt = %s, fv = %s, %s", fields.Fd3Interest, fields.Fd3Compound,
    fields.Fd3Payment, fields.Fd3FV, fields.Fd3Result), correlationId)
}
