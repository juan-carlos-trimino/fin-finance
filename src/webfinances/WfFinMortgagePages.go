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

type mortgageFields struct {
  MenuPage string `json:"menuPage"`
  CurrentPage string `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1N string `json:"fd1N"`
  Fd1TimePeriod string `json:"fd1TimePeriod"`
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1Amount string `json:"fd1Amount"`
  Fd1Result [3]string `json:"fd1Result"`
  //
  Fd2N string `json:"Fd2N"`
  Fd2TimePeriod string `json:"Fd2TimePeriod"`
  Fd2Interest string `json:"fd2Interest"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2Amount string `json:"fd2Amount"`
  Fd2TotalCost string `json:"fd2TotalCost"`
  Fd2TotalInterest string `json:"fd2TotalInterest"`
  Fd2Result []Row `json:"fd2Result"`
  //
  Fd3Mrate string `json:"fd3Mrate"`
  Fd3Mbalance string `json:"fd3Mbalance"`
  Fd3Hrate string `json:"fd3Hrate"`
  Fd3Hbalance string `json:"fd3Hbalance"`
  Fd3Result [3]string `json:"fd3Result"`
}

func newMortgageFields(dir1, dir2, correlationId string) *mortgageFields {
  //Default values returned if file is missing, empty, or JSON is corrupt.
  defaults := mortgageFields{
    MenuPage: "",
    CurrentButton: "lhs-button1",
    CurrentPage: "rhs-ui1",
    //
    Fd1N: "30.0",
    Fd1TimePeriod: "year",
    Fd1Interest: "7.50",
    Fd1Compound: "monthly",
    Fd1Amount: "100000.00",
    Fd1Result: [3]string { "", "", "" },
    //
    Fd2N: "30.0",
    Fd2TimePeriod: "year",
    Fd2Interest: "3.00",
    Fd2Compound: "monthly",
    Fd2Amount: "100000.00",
    Fd2TotalCost: "",
    Fd2TotalInterest: "",
    Fd2Result: []Row{},
    //
    Fd3Mrate: "3.375",
    Fd3Mbalance: "300000.00",
    Fd3Hrate: "2.875",
    Fd3Hbalance: "100000.00",
    Fd3Result: [3]string { mortgage_notes[0], mortgage_notes[1], "" },
  }
  return loadFieldsFromDisk(dir1, dir2, "mortgage.txt", correlationId, &defaults)
}

func getMortgageFields(userName string) *mortgageFields {
  return getFieldsGeneric(userName, func(s *UserSession) *mortgageFields {
    return s.Fields.mortgage
  })
}

var mortgage_notes = [...]string {
  "Refinance mortgage and HELOC with one load.",
  "If the blended interest rate is higher than what you could get on a new fixed-rate mortgage, consider it.",
}

type Row struct { //Rows for the amortization table.
  PaymentNo string
  Payment, PmtPrincipal, PmtInterest, Balance string
}

type WfMortgagePages struct {}

func (mp WfMortgagePages) MortgagePages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.MortgagePages.", correlationId)
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
  fields := getMortgageFields(userName)
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
    partialTemplate = "costofmortgage.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd1N string
      Fd1TimePeriod string
      Fd1Interest string
      Fd1Compound string
      Fd1Amount string
      Fd1Result [3]string
    }{
      "standard",
      "Mortgage",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd1N,
      fields.Fd1TimePeriod,
      fields.Fd1Interest,
      fields.Fd1Compound,
      fields.Fd1Amount,
      fields.Fd1Result,
    }
  case "rhs-ui2":
    fields.CurrentButton = "lhs-button2"
    if req.Method == http.MethodPost {
      mp.processUi2Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "amortizationtable.html"
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
      Fd2Amount string
      Fd2TotalCost string
      Fd2TotalInterest string
      Fd2Result []Row
    }{
      "standard",
      "Mortgage",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd2N,
      fields.Fd2TimePeriod,
      fields.Fd2Interest,
      fields.Fd2Compound,
      fields.Fd2Amount,
      fields.Fd2TotalCost,
      fields.Fd2TotalInterest,
      fields.Fd2Result,
    }
  case "rhs-ui3":
    fields.CurrentButton = "lhs-button3"
    if req.Method == http.MethodPost {
      mp.processUi3Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "heloc.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd3Mrate string
      Fd3Mbalance string
      Fd3Hrate string
      Fd3Hbalance string
      Fd3Result [3]string
    }{
      "standard",
      "Mortgage",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd3Mrate,
      fields.Fd3Mbalance,
      fields.Fd3Hrate,
      fields.Fd3Hbalance,
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
    "webfinances/templates/finances/mortgage/mortgage.html",
    "webfinances/templates/finances/mortgage/" + partialTemplate,
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
    go func(userData []byte, uName, cId string) {
      filePath := fmt.Sprintf("%s/%s/mortgage.txt", mainDir, uName)
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
func (mp WfMortgagePages) processUi1Form(req *http.Request, fields *mortgageFields, correlationId string) {
  fields.Fd1N = req.PostFormValue("fd1-n")
  fields.Fd1TimePeriod = req.PostFormValue("fd1-tp")
  fields.Fd1Interest = req.PostFormValue("fd1-interest")
  fields.Fd1Compound = req.PostFormValue("fd1-compound")
  fields.Fd1Amount = req.PostFormValue("fd1-amount")
  var n float64
  var i float64
  var amount float64
  var err error
  fields.Fd1Result[1] = ""
  fields.Fd1Result[2] = ""
  if n, err = strconv.ParseFloat(fields.Fd1N, 64); err != nil {
    fields.Fd1Result[0] = fmt.Sprintf("Error: %s -- %+v", fields.Fd1N, err)
  } else if i, err = strconv.ParseFloat(fields.Fd1Interest, 64); err != nil {
    fields.Fd1Result[0] = fmt.Sprintf("Error: %s -- %+v", fields.Fd1Interest, err)
  } else if amount, err = strconv.ParseFloat(fields.Fd1Amount, 64); err != nil {
    fields.Fd1Result[0] = fmt.Sprintf("Error: %s -- %+v", fields.Fd1Amount, err)
  } else {
    var m finances.Mortgage
    payment, totalCost, totalInterest := m.CostOfMortgage(amount, i / 100.0, fields.Fd1Compound[0], n, fields.Fd1TimePeriod[0])
    fields.Fd1Result[0] = fmt.Sprintf("Payment: $%.5f", payment)
    fields.Fd1Result[1] = fmt.Sprintf("Total Interest: $%.5f", totalInterest)
    fields.Fd1Result[2] = fmt.Sprintf("Total Cost: $%.5f", totalCost)
  }
  logger.LogInfo(fmt.Sprintf("n = %s, tp = %s, interest = %s, cp = %s, amount = %s, %s", fields.Fd1N, fields.Fd1TimePeriod,
    fields.Fd1Interest, fields.Fd1Compound, fields.Fd1Amount, fields.Fd1Result[0]), correlationId)
}

//Extraction helper for UI-2 calculations.
func (mp WfMortgagePages) processUi2Form(req *http.Request, fields *mortgageFields, correlationId string) {
  fields.Fd2N = req.FormValue("fd2-n")
  fields.Fd2TimePeriod = req.PostFormValue("fd2-tp")
  fields.Fd2Interest = req.PostFormValue("fd2-i")
  fields.Fd2Compound = req.PostFormValue("fd2-compound")
  fields.Fd2Amount = req.PostFormValue("fd2-amount")
  var n float64
  var i float64
  var amount float64
  var err error
  if n, err = strconv.ParseFloat(fields.Fd2N, 64); err != nil {
    fields.Fd2Result = append(fields.Fd2Result,
      Row {
        PaymentNo: fmt.Sprintf("Error: %s -- %+v", fields.Fd2N, err),
      })
  } else if i, err = strconv.ParseFloat(fields.Fd2Interest, 64); err != nil {
    fields.Fd2Result = append(fields.Fd2Result,
      Row {
        PaymentNo: fmt.Sprintf("Error: %s -- %+v", fields.Fd2Interest, err),
      })
  } else if amount, err = strconv.ParseFloat(fields.Fd2Amount, 64); err != nil {
    fields.Fd2Result = append(fields.Fd2Result,
      Row {
        PaymentNo: fmt.Sprintf("Error: %s -- %+v", fields.Fd2Amount, err),
      })
  } else {
    var m finances.Mortgage
    var at = m.AmortizationTable(amount, i / 100.0, fields.Fd1Compound[0], n, fields.Fd1TimePeriod[0])
    var numberOfRows = len(at.Rows)
    fields.Fd2Result = make([]Row, 0, numberOfRows + 1)
    fields.Fd2Result = append(fields.Fd2Result,
      Row {
        PaymentNo: "--",
        Payment: "--",
        PmtPrincipal: "--",
        PmtInterest: "--",
        Balance: fmt.Sprintf("%.5f", amount),
      })
    for idx := 0; idx < numberOfRows; idx++ {
      fields.Fd2Result = append(fields.Fd2Result,
        Row {
          PaymentNo: fmt.Sprintf("%d", idx + 1),
          Payment: fmt.Sprintf("%.5f", at.Rows[idx].Payment),
          PmtPrincipal: fmt.Sprintf("%.5f", at.Rows[idx].PmtPrincipal),
          PmtInterest: fmt.Sprintf("%.5f", at.Rows[idx].PmtInterest),
          Balance: fmt.Sprintf("%.5f", at.Rows[idx].Balance),
        })
    }
    fields.Fd2TotalCost = fmt.Sprintf("Total Cost: $%.5f", at.TotalCost)
    fields.Fd2TotalInterest = fmt.Sprintf("Total Interest: $%.5f", at.TotalInterest)
  }
  logger.LogInfo(fmt.Sprintf("n = %s, tp = %s, interest = %s, cp = %s, amount = %s, total cost = %s, total interest = %s",
    fields.Fd2N, fields.Fd2TimePeriod, fields.Fd2Interest, fields.Fd2Compound, fields.Fd2Amount, fields.Fd2TotalCost,
    fields.Fd2TotalInterest), correlationId)
}

//Extraction helper for UI-3 calculations.
func (mp WfMortgagePages) processUi3Form(req *http.Request, fields *mortgageFields, correlationId string) {
  fields.Fd3Mrate = req.PostFormValue("fd3-mrate")
  fields.Fd3Mbalance = req.PostFormValue("fd3-mbalance")
  fields.Fd3Hrate = req.PostFormValue("fd3-hrate")
  fields.Fd3Hbalance = req.PostFormValue("fd3-hbalance")
  var mRate float64
  var mBalance float64
  var hRate float64
  var hBalance float64
  var err error
  if mRate, err = strconv.ParseFloat(fields.Fd3Mrate, 64); err != nil {
    fields.Fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", fields.Fd3Mrate, err)
  } else if mBalance, err = strconv.ParseFloat(fields.Fd3Mbalance, 64); err != nil {
    fields.Fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", fields.Fd3Mbalance, err)
  } else if hRate, err = strconv.ParseFloat(fields.Fd3Hrate, 64); err != nil {
    fields.Fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", fields.Fd3Hrate, err)
  } else if hBalance, err = strconv.ParseFloat(fields.Fd3Hbalance, 64); err != nil {
    fields.Fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", fields.Fd3Hbalance, err)
  } else {
    var m finances.Mortgage
    fields.Fd3Result[2] = fmt.Sprintf("Blended Interest Rate: %.5f%%", m.BlendedInterestRate(mBalance, mRate, hBalance, hRate))
  }
  logger.LogInfo(fmt.Sprintf("mortgage balance = %s, mortgage rate = %s, HELOC balance = %s, HELOC rate = %s, %s",
    fields.Fd3Mbalance, fields.Fd3Mrate, fields.Fd3Hbalance, fields.Fd3Hrate, fields.Fd3Result[2]), correlationId)
}
