package webfinances

import (
  "context"
  "encoding/json"
  bank "finance/databases/banking" //Importing a package and assigning it a local alias.
  "finance/finances"
  "finance/renderer"
  "fmt"
  "net/http"
  "strconv"
  "strings"
  "github.com/juan-carlos-trimino/go-middlewares"
  logger "github.com/juan-carlos-trimino/go-logger"
  sess "github.com/juan-carlos-trimino/go-sessions"
)

type adEppFields struct {
  MenuPage      string `json:"menuPage"`
  CurrentPage   string `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1N          string `json:"fd1N"`
  Fd1TimePeriod string `json:"fd1TimePeriod"`
  Fd1Interest   string `json:"fd1Interest"`
  Fd1Compound   string `json:"fd1Compound"`
  Fd1FV         string `json:"fd1FV"`
  Fd1Result     string `json:"fd1Result"`
  //
  Fd2N          string `json:"fd2N"`
  Fd2TimePeriod string `json:"fd2TimePeriod"`
  Fd2Interest   string `json:"fd2Interest"`
  Fd2Compound   string `json:"fd2Compound"`
  Fd2PV         string `json:"fd2PV"`
  Fd2Result     string `json:"fd2Result"`
}

func defaultAdEppFields() *adEppFields {
  return &adEppFields{
    MenuPage:      "",
    CurrentPage:   "rhs-ui1",
    CurrentButton: "lhs-button1",
    //
    Fd1N:          "1.00",
    Fd1TimePeriod: "year",
    Fd1Interest:   "1.00",
    Fd1Compound:   "annually",
    Fd1FV:         "1.00",
    Fd1Result:     "",
    //
    Fd2N:          "1.00",
    Fd2TimePeriod: "year",
    Fd2Interest:   "1.00",
    Fd2Compound:   "annually",
    Fd2PV:         "1.00",
    Fd2Result:     "",
  }
}

type WfAdEppPages struct{}

func (a WfAdEppPages) AdEppPages(res http.ResponseWriter, req *http.Request) {
  ck := middlewares.MwContextKey{}
  correlationId, _ := ck.GetCorrelationId(req.Context())
  logger.LogInfo("Entering webfinances.AdEppPages.", correlationId)
  sessInfo, _ := ck.GetSessionInfo(req.Context())
  redisKey := "user:data:" + sessInfo.UserName
  partitionName := "ad_equal_prdc_pmnts"
  var fields *adEppFields = nil
  var jsonBytes []byte
  //Check Redis Cache First.
  jsonBytes, err := sess.HGetRedis(req.Context(), redisKey, partitionName)
  if err == nil {
    logger.LogInfo(fmt.Sprintf("Redis found key(%s)/partition(%s)...", redisKey, partitionName), correlationId)
    //Unmarshal unmarshals the JSON into the value pointed at by v. If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
    err = json.Unmarshal(jsonBytes, &fields)
    if err != nil {
      logger.LogInfo(fmt.Sprintf("Error (Redis) unmarshalling key(%s)/partition(%s): %v. Using default values.", redisKey, partitionName, err), correlationId)
      fields = defaultAdEppFields()
    }
  } else {
    //Cache Miss: Fall back to PostgreSQL database.
    logger.LogInfo(fmt.Sprintf("Redis does not have key(%s)/partition(%s): %v. Fetching from PostgreSQL...", redisKey, partitionName, err), correlationId)
    jsonBytes = bank.DbFetchUserData(req.Context(), sessInfo.UserName, partitionName, correlationId)
    if len(jsonBytes) == 0 {  //On error (nil) or not found (len() == 0), use default values.
      logger.LogInfo(fmt.Sprintf("Database did not find user(%s)/partition(%s). Using default values.", sessInfo.UserName, partitionName), correlationId)
      fields = defaultAdEppFields()
    } else {  //Use the values from the db.
      logger.LogInfo(fmt.Sprintf("Database found user(%s)/partition(%s)...", sessInfo.UserName, partitionName), correlationId)
      //Unmarshal unmarshals the JSON into the value pointed at by v. If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
      err = json.Unmarshal(jsonBytes, &fields)
      if err != nil {
        logger.LogInfo(fmt.Sprintf("Error (db) unmarshalling user(%s)/partition(%s): %v. Using default values.", sessInfo.UserName, partitionName, err), correlationId)
        fields = defaultAdEppFields()
      }
    }
  }
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
  if ui := req.FormValue("compute"); ui != "" { //Values from form and URL.
    fields.CurrentPage = ui
  }
  //Dynamic variables determined by the routing route condition.
  var partialTemplate string
  var templateData interface{}
  switch strings.ToLower(fields.CurrentPage) {
  case "rhs-ui1":
    fields.CurrentButton = "lhs-button1"
    if req.Method == http.MethodPost {
      a.processUi1Form(req, fields, correlationId)
    }
    partialTemplate = "n-i-FV.html"
    templateData = struct {
      LayoutType    string
      Header        string
      Datetime      string
      MenuPage      string
      CurrentButton string
      CsrfToken     string
      Fd1N          string
      Fd1TimePeriod string
      Fd1Interest   string
      Fd1Compound   string
      Fd1FV         string
      Fd1Result     string
    }{
      "standard",
      "Annuity Due / Equal Periodic Payments",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      sessInfo.CSRFToken,
      fields.Fd1N,
      fields.Fd1TimePeriod,
      fields.Fd1Interest,
      fields.Fd1Compound,
      fields.Fd1FV, fields.Fd1Result,
    }
  case "rhs-ui2":
    fields.CurrentButton = "lhs-button2"
    if req.Method == http.MethodPost {
      a.processUi2Form(req, fields, correlationId)
    }
    partialTemplate = "n-i-PV.html"
    templateData = struct {
      LayoutType    string
      Header        string
      Datetime      string
      MenuPage      string
      CurrentButton string
      CsrfToken     string
      Fd1N          string
      Fd1TimePeriod string
      Fd1Interest   string
      Fd1Compound   string
      Fd1FV         string
      Fd1Result     string
    }{
      "standard",
      "Annuity Due / Equal Periodic Payments",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      sessInfo.CSRFToken,
      fields.Fd1N,
      fields.Fd1TimePeriod,
      fields.Fd1Interest,
      fields.Fd1Compound,
      fields.Fd1FV,
      fields.Fd1Result,
    }
  default:
    errString := fmt.Sprintf("Unsupported page: %s", fields.CurrentPage)
    logger.LogError(errString, correlationId)
    panic(errString)
  }
  //Unified execution of templates.
  templatesNeeded := []string{
    "webfinances/templates/layout.html",
    "webfinances/templates/finances/annuitydue/epp/epp.html",
    "webfinances/templates/finances/annuitydue/epp/" + partialTemplate,
    "webfinances/templates/title.html",
    "webfinances/templates/datetime.html",
    "webfinances/templates/navbar.html",
    "webfinances/templates/footer.html",
  }
  renderer.Render(res, "layout", templatesNeeded, renderer.PageData{Data: templateData})
  data, err := json.Marshal(fields) //Preserve current choices.
  if err != nil {
    //Don't crash the server (panic), but log it clearly so you can debug the serialization.
    logger.LogError(fmt.Sprintf("Error marshalling user(%s)/partition(%s): %v", sessInfo.UserName, partitionName, err), correlationId)
  } else {
    go func(userData []byte, uName, redisKey, pName, cId string) {
      ctx := context.Background()
      err := sess.HSetRedis(ctx, uName, pName, userData)
      if err != nil {
        logger.LogInfo(fmt.Sprintf("Error (Redis) writing state for key(%s)/partition(%s).", redisKey, pName), cId)
      }
      ok := bank.DbSaveUserData(ctx, uName, pName, cId, userData)
      if ok {
        logger.LogInfo(fmt.Sprintf("Db successfully persisted state for user(%s)/partition(%s).", uName, pName), cId)
      }
    }(data, sessInfo.UserName, redisKey, partitionName, correlationId)
  }
}

//Extraction helper for UI-1 calculations.
func (a WfAdEppPages) processUi1Form(req *http.Request, fields *adEppFields, correlationId string) {
  fields.Fd1N = req.PostFormValue("fd1-n")
  fields.Fd1TimePeriod = req.PostFormValue("fd1-tp")
  fields.Fd1Interest = req.PostFormValue("fd1-interest")
  fields.Fd1Compound = req.PostFormValue("fd1-cp")
  fields.Fd1FV = req.PostFormValue("fd1-fv")
  var n float64
  var i float64
  var fv float64
  var err error
  if n, err = strconv.ParseFloat(fields.Fd1N, 64); err != nil {
    fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1N, err)
  } else if i, err = strconv.ParseFloat(fields.Fd1Interest, 64); err != nil {
    fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1Interest, err)
  } else if fv, err = strconv.ParseFloat(fields.Fd1FV, 64); err != nil {
    fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1FV, err)
  } else {
    var oa finances.Annuities
    fields.Fd1Result = fmt.Sprintf("Payment: $%.5f", oa.D_Payment_FV(fv, i/100.0,
      oa.GetCompoundingPeriod(fields.Fd1Compound[0], true), n, oa.GetTimePeriod(fields.Fd1TimePeriod[0], true)))
  }
  logger.LogInfo(fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, fv = %s, %s", fields.Fd1N,
    fields.Fd1TimePeriod, fields.Fd1Interest, fields.Fd1Compound, fields.Fd1FV, fields.Fd1Result), correlationId)
}

//Extraction helper for UI-2 calculations.
func (a WfAdEppPages) processUi2Form(req *http.Request, fields *adEppFields, correlationId string) {
  fields.Fd2N = req.PostFormValue("fd2-n")
  fields.Fd2TimePeriod = req.PostFormValue("fd2-tp")
  fields.Fd2Interest = req.FormValue("fd2-interest")
  fields.Fd2Compound = req.PostFormValue("fd2-cp")
  fields.Fd2PV = req.PostFormValue("fd2-pv")
  var n float64
  var i float64
  var pv float64
  var err error
  if n, err = strconv.ParseFloat(fields.Fd2N, 64); err != nil {
    fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2N, err)
  } else if i, err = strconv.ParseFloat(fields.Fd2Interest, 64); err != nil {
    fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Interest, err)
  } else if pv, err = strconv.ParseFloat(fields.Fd2PV, 64); err != nil {
    fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2PV, err)
  } else {
    var oa finances.Annuities
    fields.Fd2Result = fmt.Sprintf("Payment: $%.5f", oa.D_Payment_PV(pv, i/100.0,
      oa.GetCompoundingPeriod(fields.Fd2Compound[0], true), n, oa.GetTimePeriod(fields.Fd2TimePeriod[0], true)))
  }
  logger.LogInfo(fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, pv = %s, %s", fields.Fd2N,
    fields.Fd2TimePeriod, fields.Fd2Interest, fields.Fd2Compound, fields.Fd2PV, fields.Fd2Result), correlationId)
}
