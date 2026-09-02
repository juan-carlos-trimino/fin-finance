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

type miscellaneousFields struct{
  MenuPage string `json:"menuPage"`
  CurrentPage string `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1Nominal string `json:"fd1Nominal"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1Result [2]string `json:"fd1Result"`
  //
  Fd2Effective string `json:"fd2Effective"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2Result [3]string `json:"fd2Result"`
  //
  Fd3Nominal string `json:"fd3Nominal"`
  Fd3Inflation string `json:"fd3Inflation"`
  Fd3Result [4]string `json:"fd3Result"`
  //
  Fd4CurrentRate string `json:"fd4CurrentRate"`
  Fd4CurrentCompound string `json:"fd4CurrentCompound"`
  Fd4NewCompound string `json:"fd4NewCompound"`
  Fd4Result string `json:"fd4Result"`
  //
  Fd5Interest string `json:"fd5Interest"`
  Fd5Compound string `json:"fd5Compound"`
  Fd5Factor string `json:"fd5Factor"`
  Fd5Result string `json:"fd5Result"`
  //
  Fd6Values string `json:"fd6Values"`
  Fd6Result [2]string `json:"fd6Result"`
  //
  Fd7Time string `json:"fd7Time"`
  Fd7TimePeriod string `json:"fd7TimePeriod"`
  Fd7Rate string `json:"fd7Rate"`
  Fd7Compound string `json:"fd7Compound"`
  Fd7PV string `json:"fd7PV"`
  Fd7Result string `json:"fd7Result"`
}

func newMiscellaneousFields(dir1, dir2, correlationId string) *miscellaneousFields {
  //Default values returned if file is missing, empty, or JSON is corrupt.
  defaults := miscellaneousFields{
    MenuPage: "",
    CurrentPage: "rhs-ui1",
    CurrentButton: "lhs-button1",
    //
    Fd1Nominal: "3.5",
    Fd1Compound: "monthly",
    Fd1Result: [2]string { misc_notes[0], "" },
    //
    Fd2Effective: "3.5",
    Fd2Compound: "monthly",
    Fd2Result: [3]string { misc_notes[0], misc_notes[1], "" },
    //
    Fd3Nominal: "2.0",
    Fd3Inflation: "2.0",
    Fd3Result: [4]string { misc_notes[1], misc_notes[2], misc_notes[3], "" },
    //
    Fd4CurrentRate: "9.00",
    Fd4CurrentCompound: "annually",
    Fd4NewCompound: "monthly",
    Fd4Result: "",
    //
    Fd5Interest: "14.87",
    Fd5Compound: "annually",
    Fd5Factor: "2.0",
    Fd5Result: "",
    //
    Fd6Values: "2.0;1.5",
    Fd6Result: [2]string { misc_notes[4], "" },
    //
    Fd7Time: "1.0",
    Fd7TimePeriod: "year",
    Fd7Rate: "15.0",
    Fd7Compound: "annually",
    Fd7PV: "1.00",
    Fd7Result: "",
  }
  return loadFieldsFromDisk(dir1, dir2, "miscellaneous.txt", correlationId, &defaults)
}

func getMiscellaneousFields(userName string) *miscellaneousFields {
  return getFieldsGeneric(userName, func(s *UserSession) *miscellaneousFields {
    return s.Fields.miscellaneous
  })
}

var misc_notes = [...]string {
  "When comparing interest rates, use effective annual rates.",
  "Nominal returns are not adjusted for inflation.",
  "Real returns are useful while comparing returns over different time periods because of the differences in inflation rates.",
  "Real returns are adjusted for inflation.",
  "Values are semicolon (;) separated; e.g., 3;3.1;3.2;-1.01",
}

type WfMiscellaneousPages struct {}

func (mp WfMiscellaneousPages) MiscellaneousPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.MiscellaneousPages.", correlationId)
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
  fields := getMiscellaneousFields(userName)
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
  case "rhs-ui1":
    fields.CurrentButton = "lhs-button1"
    if req.Method == http.MethodPost {
      mp.processUi1Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "nominalrate.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd1Nominal string
      Fd1Compound string
      Fd1Result [2]string
    }{
      "standard",
      "Miscellaneous",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd1Nominal,
      fields.Fd1Compound,
      fields.Fd1Result,
    }
  case "rhs-ui2":
    fields.CurrentButton = "lhs-button2"
    if req.Method == http.MethodPost {
      mp.processUi2Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "effectiveannualrate.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd2Effective string
      Fd2Compound string
      Fd2Result [3]string
    }{
      "standard",
      "Miscellaneous",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd2Effective,
      fields.Fd2Compound,
      fields.Fd2Result,
    }
  case "rhs-ui3":
    fields.CurrentButton = "lhs-button3"
    if req.Method == http.MethodPost {
      mp.processUi3Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "nominalratevs.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd3Nominal string
      Fd3Inflation string
      Fd3Result [4]string
    }{
      "standard",
      "Miscellaneous",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd3Nominal,
      fields.Fd3Inflation,
      fields.Fd3Result,
    }
  case "rhs-ui4":
    fields.CurrentButton = "lhs-button4"
    if req.Method == http.MethodPost {
      mp.processUi4Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "compfrequencyconv.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd4CurrentRate string
      Fd4CurrentCompound string
      Fd4NewCompound string
      Fd4Result string
    }{
      "standard",
      "Miscellaneous",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd4CurrentRate,
      fields.Fd4CurrentCompound,
      fields.Fd4NewCompound,
      fields.Fd4Result,
    }
  case "rhs-ui5":
    fields.CurrentButton = "lhs-button5"
    if req.Method == http.MethodPost {
      mp.processUi5Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "growthdecay.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd5Interest string
      Fd5Compound string
      Fd5Factor string
      Fd5Result string
    }{
      "standard",
      "Miscellaneous",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd5Interest,
      fields.Fd5Compound,
      fields.Fd5Factor,
      fields.Fd5Result,
    }
  case "rhs-ui6":
    fields.CurrentButton = "lhs-button6"
    if req.Method == http.MethodPost {
      mp.processUi6Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "averagerate.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd6Values string
      Fd6Result [2]string
    }{
      "standard",
      "Miscellaneous",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd6Values,
      fields.Fd6Result,
    }
  case "rhs-ui7":
    fields.CurrentButton = "lhs-button7"
    if req.Method == http.MethodPost {
      mp.processUi7Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "depreciation.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd7Time string
      Fd7TimePeriod string
      Fd7Rate string
      Fd7Compound string
      Fd7PV string
      Fd7Result string
    }{
      "standard",
      "Miscellaneous",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd7Time,
      fields.Fd7TimePeriod,
      fields.Fd7Rate,
      fields.Fd7Compound,
      fields.Fd7PV,
      fields.Fd7Result,
    }
  default:
    errString := fmt.Sprintf("Unsupported page: %s", fields.CurrentPage)
    logger.LogError(errString, correlationId)
    panic(errString)
  }
  //Unified execution of templates.
  templatesNeeded := []string{
    "webfinances/templates/layout.html",
    "webfinances/templates/finances/miscellaneous/miscellaneous.html",
    "webfinances/templates/finances/miscellaneous/" + partialTemplate,
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
      filePath := fmt.Sprintf("%s/%s/miscellaneous.txt", mainDir, uName)
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

//Extraction helper for UI-1 calculations.
func (mp WfMiscellaneousPages) processUi1Form(req *http.Request, fields *miscellaneousFields, correlationId string) {
  fields.Fd1Nominal = req.PostFormValue("fd1-nominal")
  fields.Fd1Compound = req.PostFormValue("fd1-compound")
  var nr float64
  var err error
  if nr, err = strconv.ParseFloat(fields.Fd1Nominal, 64); err != nil {
    fields.Fd1Result[1] = fmt.Sprintf("Error: %s -- %+v", fields.Fd1Nominal, err)
  } else {
    var a finances.Annuities
    fields.Fd1Result[1] = fmt.Sprintf("Effective Annual Rate: %.5f%%",
      a.NominalRateToEAR(nr / 100.0, a.GetCompoundingPeriod(fields.Fd1Compound[0], false)) * 100.0)
  }
  logger.LogInfo(fmt.Sprintf("nominal rate = %s, cp = %s, %s", fields.Fd1Nominal, fields.Fd1Compound,
    fields.Fd1Result[1]), correlationId)
}

//Extraction helper for UI-2 calculations.
func (mp WfMiscellaneousPages) processUi2Form(req *http.Request, fields *miscellaneousFields, correlationId string) {
  fields.Fd2Effective = req.PostFormValue("fd2-effective")
  fields.Fd2Compound = req.PostFormValue("fd2-compound")
  var ear float64
  var err error
  if ear, err = strconv.ParseFloat(fields.Fd2Effective, 64); err != nil {
    fields.Fd2Result[2] = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Effective, err)
  } else {
    var a finances.Annuities
    fields.Fd2Result[2] = fmt.Sprintf("Nominal Rate: %.5f%% %s",
      a.EARToNominalRate(ear / 100.0, a.GetCompoundingPeriod(fields.Fd2Compound[0], false)) * 100.0, fields.Fd2Compound)
  }
  logger.LogInfo(fmt.Sprintf("effective rate = %s, cp = %s, %s", fields.Fd2Effective, fields.Fd2Compound, fields.Fd2Result[2]),
    correlationId)
}

//Extraction helper for UI-3 calculations.
func (mp WfMiscellaneousPages) processUi3Form(req *http.Request, fields *miscellaneousFields, correlationId string) {
  fields.Fd3Nominal = req.PostFormValue("fd3-nominal")
  fields.Fd3Inflation = req.PostFormValue("fd3-inflation")
  var nr float64
  var ir float64
  var err error
  if nr, err = strconv.ParseFloat(fields.Fd3Nominal, 64); err != nil {
    fields.Fd3Result[3] = fmt.Sprintf("Error: %s -- %+v", fields.Fd3Nominal, err)
  } else if ir, err = strconv.ParseFloat(fields.Fd3Inflation, 64); err != nil {
    fields.Fd3Result[3] = fmt.Sprintf("Error: %s -- %+v", fields.Fd3Inflation, err)
  } else {
    var a finances.Annuities
    fields.Fd3Result[3] = fmt.Sprintf("Real Interest Rate: %.5f%%", a.RealInterestRate(nr / 100.0, ir / 100.0) * 100.0)
  }
  logger.LogInfo(fmt.Sprintf("nominal rate = %s, inflation rate = %s, %s", fields.Fd3Nominal, fields.Fd3Inflation,
    fields.Fd3Result[3]), correlationId)
}

//Extraction helper for UI-4 calculations.
func (mp WfMiscellaneousPages) processUi4Form(req *http.Request, fields *miscellaneousFields, correlationId string) {
  fields.Fd4CurrentRate = req.PostFormValue("fd4-currentrate")
  fields.Fd4CurrentCompound = req.PostFormValue("fd4-currentcompound")
  fields.Fd4NewCompound = req.PostFormValue("fd4-newcompound")
  var newDays, currentDays int
  var currentRate float64
  var err error
  if currentRate, err = strconv.ParseFloat(fields.Fd4CurrentRate, 64); err != nil {
    fields.Fd4Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd4CurrentRate, err)
  } else if strings.EqualFold(fields.Fd4CurrentCompound[0:1], "D") {
    if currentDays, err = strconv.Atoi(fields.Fd4CurrentCompound[5:len(fields.Fd4CurrentCompound)]); err != nil {
      fields.Fd4Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd4CurrentCompound, err)
    }
  }
  //
  if err == nil && strings.EqualFold(fields.Fd4NewCompound[0:1], "D") {
    if newDays, err = strconv.Atoi(fields.Fd4NewCompound[5:len(fields.Fd4NewCompound)]); err != nil {
      fields.Fd4Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd4NewCompound, err)
    }
  }
  //
  if err == nil {
    var isCurrentDaily365 bool = false
    if currentDays == finances.Daily365 {
      isCurrentDaily365 = true
    }
    var isNewDaily365 bool = false
    if newDays == finances.Daily365 {
      isNewDaily365 = true
    }
    var a finances.Annuities
    fields.Fd4Result = fmt.Sprintf("New Rate: %.5f%%",
      a.CompoundingFrequencyConversion(currentRate / 100.0,
      a.GetCompoundingPeriod(fields.Fd4CurrentCompound[0], isCurrentDaily365),
      a.GetCompoundingPeriod(fields.Fd4NewCompound[0], isNewDaily365)) * 100.0)
  }
  logger.LogInfo(fmt.Sprintf("current rate = %s, current compound = %s, new compound = %s, %s", fields.Fd4CurrentRate,
    fields.Fd4CurrentCompound, fields.Fd4NewCompound, fields.Fd4Result), correlationId)
}

//Extraction helper for UI-5 calculations.
func (mp WfMiscellaneousPages) processUi5Form(req *http.Request, fields *miscellaneousFields, correlationId string) {
  fields.Fd5Interest = req.PostFormValue("fd5-interest")
  fields.Fd5Compound = req.PostFormValue("fd5-compound")
  fields.Fd5Factor = req.PostFormValue("fd5-factor")
  var ir float64
  var factor float64
  var err error
  if ir, err = strconv.ParseFloat(fields.Fd5Interest, 64); err != nil {
    fields.Fd5Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd5Interest, err)
  } else if factor, err = strconv.ParseFloat(fields.Fd5Factor, 64); err != nil {
    fields.Fd5Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd5Factor, err)
  } else {
    var a finances.Annuities
    fields.Fd5Result = fmt.Sprintf("Growth/Decay: %.5f %s",
      a.GrowthDecayOfFunds(factor, ir / 100.0, a.GetCompoundingPeriod(fields.Fd5Compound[0], true)),
      a.TimePeriods(fields.Fd5Compound))
  }
  logger.LogInfo(fmt.Sprintf("interest rate = %s, cp = %s, factor = %s, %s\n",
    fields.Fd5Interest, fields.Fd5Compound, fields.Fd5Factor, fields.Fd5Result), correlationId)
}

//Extraction helper for UI-6 calculations.
func (mp WfMiscellaneousPages) processUi6Form(req *http.Request, fields *miscellaneousFields, correlationId string) {
  fields.Fd6Values = req.PostFormValue("fd6-values")
  split := strings.Split(fields.Fd6Values, ";")
  values := make([]float64, len(split))
  var err error
  for i, s := range split {
    if values[i], err = strconv.ParseFloat(s, 64); err != nil {
      fields.Fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", s, err)
      break;
    }
  }
  //
  if err == nil {
    var a finances.Annuities
    fields.Fd6Result[1] = fmt.Sprintf("Avg: %.5f%%", a.AverageRateOfReturn(values) * 100.0)
  }
  logger.LogInfo(fmt.Sprintf("values = [%s], %s\n", fields.Fd6Values, fields.Fd6Result[1]), correlationId)
}

//Extraction helper for UI-7 calculations.
func (mp WfMiscellaneousPages) processUi7Form(req *http.Request, fields *miscellaneousFields, correlationId string) {
  fields.Fd7Time = req.PostFormValue("fd7-time")
  fields.Fd7TimePeriod = req.PostFormValue("fd7-tp")
  fields.Fd7Rate = req.PostFormValue("fd7-rate")
  fields.Fd7Compound = req.PostFormValue("fd7-compound")
  fields.Fd7PV = req.PostFormValue("fd7-pv")
  var time float64
  var rate float64
  var pv float64
  var err error
  if time, err = strconv.ParseFloat(fields.Fd7Time, 64); err != nil {
    fields.Fd7Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd7Time, err)
  } else if rate, err = strconv.ParseFloat(fields.Fd7Rate, 64); err != nil {
    fields.Fd7Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd7Rate, err)
  } else if pv, err = strconv.ParseFloat(fields.Fd7PV, 64); err != nil {
    fields.Fd7Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd7PV, err)
  } else {
    var a finances.Annuities
    fields.Fd7Result = fmt.Sprintf("Future Value: %.5f",
      a.Depreciation(pv, rate / 100.0, a.GetCompoundingPeriod(fields.Fd7Compound[0], false), time,
      a.GetTimePeriod(fields.Fd7TimePeriod[0], false)))
  }
  logger.LogInfo(fmt.Sprintf("time = %s, tp = %s, rate = %s, cp = %s, pv = %s, %s\n", fields.Fd7Time,
    fields.Fd7TimePeriod, fields.Fd7Rate, fields.Fd7Compound, fields.Fd7PV, fields.Fd7Result), correlationId)
}
