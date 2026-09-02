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
  "math"
  "net/http"
  "os"
  "strconv"
  "strings"
  "time"
)

type bondsFields struct {
  MenuPage string `json:"menuPage"`
  CurrentPage string `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1TaxFree string `json:"fd1TaxFree"`
  Fd1CityTax string `json:"fd1CityTax"`
  Fd1StateTax string `json:"fd1StateTax"`
  Fd1FederalTax string `json:"fd1FederalTax"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2FaceValue string `json:"fd2FaceValue"`
  Fd2Time string `json:"fd2Time"`
  Fd2TimePeriod string `json:"Fd2TimePeriod"`
  Fd2Coupon string `json:"fd2Coupon"`
  Fd2CompoundCoupon string `json:"fd2CompoundCoupon"`
  Fd2Current string `json:"fd2Current"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2Result string `json:"fd2Result"`
  //
  Fd3FaceValue string `json:"fd3FaceValue"`
  Fd3TimeCall string `json:"fd3TimeCall"`
  Fd3TimePeriod string `json:"fd3TimePeriod"`
  Fd3Coupon string `json:"fd3Coupon"`
  Fd3Compound string `json:"fd3Compound"`
  Fd3BondPrice string `json:"fd3BondPrice"`
  Fd3CallPrice string `json:"fd3CallPrice"`
  Fd3Result string `json:"fd3Result"`
  //
  Fd4FaceValue string `json:"fd4FaceValue"`
  Fd4Time string `json:"fd4Time"`
  Fd4TimePeriod string `json:"fd4TimePeriod"`
  Fd4Coupon string `json:"fd4Coupon"`
  Fd4Compound string `json:"fd4Compound"`
  Fd4CurrentRadio string `json:"fd4CurrentRadio"`
  Fd4CurInterest string `json:"fd4CurInterest"`
  Fd4BondPrice string `json:"fd4BondPrice"`
  Fd4Result [2]string `json:"fd4Result"`
  //
  Fd5FaceValue string `json:"fd5FaceValue"`
  Fd5Time string `json:"fd5Time"`
  Fd5TimePeriod string `json:"fd5TimePeriod"`
  Fd5Coupon string `json:"fd5Coupon"`
  Fd5CompoundCoupon string `json:"fd5CompoundCoupon"`
  Fd5CurInterest string `json:"fd5CurInterest"`
  Fd5Compound string `json:"fd5Compound"`
  Fd5Result [7]string `json:"fd5Result"`
}

func newBondsFields(dir1, dir2, correlationId string) *bondsFields {
  //Default values returned if file is missing, empty, or JSON is corrupt.
  defaults := bondsFields{
    MenuPage: "",
    CurrentPage: "rhs-ui1",
    CurrentButton: "lhs-button1",
    //
    Fd1TaxFree: "3.5",
    Fd1CityTax: "0.0",
    Fd1StateTax: "1.0",
    Fd1FederalTax: "23.0",
    Fd1Result: "",
    //
    Fd2FaceValue: "1000.00",
    Fd2Time: "5",
    Fd2TimePeriod: "year",
    Fd2Coupon: "3.00",
    Fd2CompoundCoupon: "annually",
    Fd2Current: "3.5",
    Fd2Compound: "annually",
    Fd2Result: "",
    //
    Fd3FaceValue: "1000.00",
    Fd3TimeCall: "2",
    Fd3TimePeriod: "year",
    Fd3Coupon: "2.0",
    Fd3Compound: "semiannually",
    Fd3BondPrice: "990.00",
    Fd3CallPrice: "1050.00",
    Fd3Result: "",
    //
    Fd4FaceValue: "1000.00",
    Fd4Time: "3",
    Fd4TimePeriod: "year",
    Fd4Coupon: "2.5",
    Fd4Compound: "semiannually",
    Fd4CurrentRadio: "fd4-curinterest",
    Fd4CurInterest: "2.3",
    Fd4BondPrice: "1000.00",
    Fd4Result: [2]string { "", "" },
    //
    Fd5FaceValue: "1000.00",
    Fd5Time: "5",
    Fd5TimePeriod: "year",
    Fd5Coupon: "5.4",
    Fd5CompoundCoupon: "annually",
    Fd5CurInterest: "7.5",
    Fd5Compound: "annually",
    Fd5Result: [7]string { "", "", "", "", "", "", "" },
  }
  return loadFieldsFromDisk(dir1, dir2, "bonds.txt", correlationId, &defaults)
}

func getBondsFields(userName string) *bondsFields {
  return getFieldsGeneric(userName, func(s *UserSession) *bondsFields {
    return s.Fields.bonds
  })
}

var bond_notes = [...]string {
  "The Macaulay duration is a measure of a bond's sensitivity to interest rate changes. The duration is the weighed-average number of years the " +
  "investor must hold a bond until the present value of the bond's cash flows equals the amount paid for the bond.",
  "The modified duration of a bond is a measure of the sensitivity of the bond's price to changes in interest rates. Since bond prices move in " +
  "an inverse direction from interest rates, for a one percent increase (decrease) in interest rates, the bond's price will decrease (increase) " +
  "by the percentage shown by the modified duration.",
  "Convexity in bonds measures how sensitive the bond's duration is to changes in interest rates. The higher the convexity, the less the bond " +
  "price will increase when rates fall -- and the less the bond price will drop when rates rise.",
}

type WfBondsPages struct {}

func (b WfBondsPages) BondsPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.BondsPages.", correlationId)
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
  fields := getBondsFields(userName)
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
      b.processUi1Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "taxfree.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd1TaxFree string
      Fd1CityTax string
      Fd1StateTax string
      Fd1FederalTax string
      Fd1Result string
    }{
      "standard",
      "Bonds",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd1TaxFree,
      fields.Fd1CityTax,
      fields.Fd1StateTax,
      fields.Fd1FederalTax,
      fields.Fd1Result,
    }
  case "rhs-ui2":
    fields.CurrentButton = "lhs-button2"
    if req.Method == http.MethodPost {
      b.processUi2Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "currentprice.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd2FaceValue string
      Fd2Time string
      Fd2TimePeriod string
      Fd2Coupon string
      Fd2CompoundCoupon string
      Fd2Current string
      Fd2Compound string
      Fd2Result string
    }{
      "standard",
      "Bonds",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd2FaceValue,
      fields.Fd2Time,
      fields.Fd2TimePeriod,
      fields.Fd2Coupon,
      fields.Fd2CompoundCoupon,
      fields.Fd2Current,
      fields.Fd2Compound,
      fields.Fd2Result,
    }
  case "rhs-ui3":
    fields.CurrentButton = "lhs-button3"
    if req.Method == http.MethodPost {
      b.processUi3Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "yieldtocall.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd3FaceValue string
      Fd3TimeCall string
      Fd3TimePeriod string
      Fd3Coupon string
      Fd3Compound string
      Fd3BondPrice string
      Fd3CallPrice string
      Fd3Result string
    }{
      "standard",
      "Bonds",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd3FaceValue,
      fields.Fd3TimeCall,
      fields.Fd3TimePeriod,
      fields.Fd3Coupon,
      fields.Fd3Compound,
      fields.Fd3BondPrice,
      fields.Fd3CallPrice,
      fields.Fd3Result,
    }
  case "rhs-ui4":
    fields.CurrentButton = "lhs-button4"
    if req.Method == http.MethodPost {
      b.processUi4Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "yieldtomaturity.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd4FaceValue string
      Fd4Time string
      Fd4TimePeriod string
      Fd4Coupon string
      Fd4Compound string
      Fd4CurrentRadio string
      Fd4CurInterest string
      Fd4BondPrice string
      Fd4Result [2]string
    }{
      "standard",
      "Bonds",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd4FaceValue,
      fields.Fd4Time,
      fields.Fd4TimePeriod,
      fields.Fd4Coupon,
      fields.Fd4Compound,
      fields.Fd4CurrentRadio,
      fields.Fd4CurInterest,
      fields.Fd4BondPrice,
      fields.Fd4Result,
    }
  case "rhs-ui5":
    fields.CurrentButton = "lhs-button5"
    if req.Method == http.MethodPost {
      b.processUi5Form(req, fields, correlationId)
    }
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    http.SetCookie(res, sessions.CreateCookie(newSessionToken))
    partialTemplate = "duration.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd5FaceValue string
      Fd5Time string
      Fd5TimePeriod string
      Fd5Coupon string
      Fd5CompoundCoupon string
      Fd5CurInterest string
      Fd5Compound string
      Fd5Result [7]string
    }{
      "standard",
      "Bonds",
      logger.DatetimeFormat(),
      financesMenuPage,
      fields.CurrentButton,
      newSession.CsrfToken,
      fields.Fd5FaceValue,
      fields.Fd5Time,
      fields.Fd5TimePeriod,
      fields.Fd5Coupon,
      fields.Fd5CompoundCoupon,
      fields.Fd5CurInterest,
      fields.Fd5Compound,
      fields.Fd5Result,
    }
  default:
    errString := fmt.Sprintf("Unsupported page: %s", fields.CurrentPage)
    logger.LogError(errString, correlationId)
    panic(errString)
  }
  //Unified execution of templates.
  templatesNeeded := []string{
    "webfinances/templates/layout.html",
    "webfinances/templates/finances/bonds/bonds.html",
    "webfinances/templates/finances/bonds/" + partialTemplate,
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
      filePath := fmt.Sprintf("%s/%s/bonds.txt", mainDir, uName)
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
func (b WfBondsPages) processUi1Form(req *http.Request, fields *bondsFields, correlationId string) {
  fields.Fd1TaxFree = req.PostFormValue("fd1-taxfree")
  fields.Fd1CityTax = req.PostFormValue("fd1-citytax")
  fields.Fd1StateTax = req.PostFormValue("fd1-statetax")
  fields.Fd1FederalTax = req.PostFormValue("fd1-federaltax")
  var taxFree float64
  var cityTax float64
  var stateTax float64
  var federalTax float64
  var err error
  if taxFree, err = strconv.ParseFloat(fields.Fd1TaxFree, 64); err != nil {
    fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1TaxFree, err)
  } else if cityTax, err = strconv.ParseFloat(fields.Fd1CityTax, 64); err != nil {
    fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1CityTax, err)
  } else if stateTax, err = strconv.ParseFloat(fields.Fd1StateTax, 64); err != nil {
    fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1StateTax, err)
  } else if federalTax, err = strconv.ParseFloat(fields.Fd1FederalTax, 64); err != nil {
    fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1FederalTax, err)
  } else {
    var b finances.Bonds
    fields.Fd1Result = fmt.Sprintf("Taxable-Equivalent Yield: %.5f%%",
      b.TaxableVsTaxFreeYields(taxFree, cityTax, stateTax, federalTax) * 100.0)
  }
  logger.LogInfo(fmt.Sprintf("tax free = %s, city tax = %s, state tax = %s, federal tax = %s, %s", fields.Fd1TaxFree,
    fields.Fd1CityTax, fields.Fd1StateTax, fields.Fd1FederalTax, fields.Fd1Result), correlationId)
}

//Extraction helper for UI-2 calculations.
func (b WfBondsPages) processUi2Form(req *http.Request, fields *bondsFields, correlationId string) {
  fields.Fd2FaceValue = req.FormValue("fd2-facevalue")
  fields.Fd2Time = req.PostFormValue("fd2-time")
  fields.Fd2TimePeriod = req.PostFormValue("fd2-tp")
  fields.Fd2Coupon = req.PostFormValue("fd2-coupon")
  fields.Fd2CompoundCoupon = req.PostFormValue("fd2-compound-coupon")
  fields.Fd2Current = req.PostFormValue("fd2-current")
  fields.Fd2Compound = req.PostFormValue("fd2-compound")
  var fv float64
  var time float64
  var coupon float64
  var current float64
  var err error
  if fv, err = strconv.ParseFloat(fields.Fd2FaceValue, 64); err != nil {
    fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2FaceValue, err)
  } else if time, err = strconv.ParseFloat(fields.Fd2Time, 64); err != nil {
    fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Time, err)
  } else if coupon, err = strconv.ParseFloat(fields.Fd2Coupon, 64); err != nil {
    fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Coupon, err)
  } else if current, err = strconv.ParseFloat(fields.Fd2Current, 64); err != nil {
    fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Current, err)
  } else {
    var b finances.Bonds
    cf := b.CashFlow(fv, coupon, b.GetCompoundingPeriod(fields.Fd2CompoundCoupon[0], true), time,
      b.GetTimePeriod(fields.Fd2TimePeriod[0], true))
    var currentPrice float64
    switch fields.Fd2Compound[0] {
    case 'c', 'C':
      currentPrice = b.CurrentPriceContinuous(cf, current)
    default:
      currentPrice = b.CurrentPrice(cf, current, b.GetCompoundingPeriod(fields.Fd2Compound[0], true))
    }
    //
    if math.Abs(fv - currentPrice) < finances.Accuracy {
      fields.Fd2Result = fmt.Sprintf("Current Price: $%.5f (par)", currentPrice)
    } else if fv < currentPrice {
      fields.Fd2Result = fmt.Sprintf("Current Price: $%.5f (premium)", currentPrice)
    } else {
      fields.Fd2Result = fmt.Sprintf("Current Price: $%.5f (discount)", currentPrice)
    }
  }
  logger.LogInfo(fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon rate = %s, current interest = %s, cp = %s, %s", fields.Fd2FaceValue,
    fields.Fd2Time, fields.Fd2TimePeriod, fields.Fd2Coupon, fields.Fd2Current, fields.Fd2Compound, fields.Fd2Result), correlationId)
}

//Extraction helper for UI-3 calculations.
func (b WfBondsPages) processUi3Form(req *http.Request, fields *bondsFields, correlationId string) {
  fields.Fd3FaceValue = req.PostFormValue("fd3-facevalue")
  fields.Fd3TimeCall = req.PostFormValue("fd3-timecall")
  fields.Fd3TimePeriod = req.PostFormValue("fd3-tp")
  fields.Fd3Coupon = req.PostFormValue("fd3-coupon")
  fields.Fd3BondPrice = req.PostFormValue("fd3-bondprice")
  fields.Fd3CallPrice = req.PostFormValue("fd3-callprice")
  fields.Fd3Compound = req.PostFormValue("fd3-compound")
  var fv float64
  var timeToCall float64
  var couponRate float64
  var bondPrice float64
  var callPrice float64
  var err error
  if fv, err = strconv.ParseFloat(fields.Fd3FaceValue, 64); err != nil {
    fields.Fd3Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd3FaceValue, err)
  } else if timeToCall, err = strconv.ParseFloat(fields.Fd3TimeCall, 64); err != nil {
    fields.Fd3Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd3TimeCall, err)
  } else if couponRate, err = strconv.ParseFloat(fields.Fd3Coupon, 64); err != nil {
    fields.Fd3Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd3Coupon, err)
  } else if bondPrice, err = strconv.ParseFloat(fields.Fd3BondPrice, 64); err != nil {
    fields.Fd3Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd3BondPrice, err)
  } else if callPrice, err = strconv.ParseFloat(fields.Fd3CallPrice, 64); err != nil {
    fields.Fd3Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd3CallPrice, err)
  } else {
    var b finances.Bonds
    fields.Fd3Result = fmt.Sprintf("Yield to Call: %.5f%%", b.YieldToCall(fv, couponRate,
      b.GetCompoundingPeriod(fields.Fd3Compound[0], true), timeToCall,
      b.GetTimePeriod(fields.Fd3TimePeriod[0], true), bondPrice, callPrice))
  }
  logger.LogInfo(fmt.Sprintf("fv = %s, coupon rate = %s, cp = %s, time to call = %s, tp = %s, bond price = %s, call price = %s, %s",
    fields.Fd3FaceValue, fields.Fd3Coupon, fields.Fd3Compound, fields.Fd3TimeCall, fields.Fd3TimePeriod, fields.Fd3BondPrice,
    fields.Fd3CallPrice, fields.Fd3Result), correlationId)
}

//Extraction helper for UI-4 calculations.
func (b WfBondsPages) processUi4Form(req *http.Request, fields *bondsFields, correlationId string) {
  fields.Fd4FaceValue = req.PostFormValue("fd4-facevalue")
  fields.Fd4Time = req.PostFormValue("fd4-time")
  fields.Fd4TimePeriod = req.PostFormValue("fd4-tp")
  fields.Fd4Coupon = req.PostFormValue("fd4-coupon")
  fields.Fd4Compound = req.PostFormValue("fd4-compound")
  fields.Fd4CurrentRadio = req.PostFormValue("fd4-choice")
  fields.Fd4CurInterest = req.PostFormValue("fd4-ci")
  fields.Fd4BondPrice = req.PostFormValue("fd4-bp")
  var currentInterest bool = false
  if strings.EqualFold(fields.Fd4CurrentRadio, "fd4-curinterest") {
    currentInterest = true
  }
  var fv float64
  var time float64
  var couponRate float64
  var curInterest float64
  var bondPrice float64
  var err error
  if fv, err = strconv.ParseFloat(fields.Fd4FaceValue, 64); err != nil {
    fields.Fd4Result[0] = fmt.Sprintf("Error: %s -- %+v", fields.Fd4FaceValue, err)
  } else if time, err = strconv.ParseFloat(fields.Fd4Time, 64); err != nil {
    fields.Fd4Result[0] = fmt.Sprintf("Error: %s -- %+v", fields.Fd4Time, err)
  } else if couponRate, err = strconv.ParseFloat(fields.Fd4Coupon, 64); err != nil {
    fields.Fd4Result[0] = fmt.Sprintf("Error: %s -- %+v", fields.Fd4Coupon, err)
  } else if curInterest, err = strconv.ParseFloat(fields.Fd4CurInterest, 64); err != nil {
    fields.Fd4Result[0] = fmt.Sprintf("Error: %s -- %+v", fields.Fd4CurInterest, err)
  } else if bondPrice, err = strconv.ParseFloat(fields.Fd4BondPrice, 64); err != nil {
    fields.Fd4Result[0] = fmt.Sprintf("Error: %s -- %+v", fields.Fd4BondPrice, err)
  } else {
    var b finances.Bonds
    var cp int = b.GetCompoundingPeriod(fields.Fd4Compound[0], true)
    var tp = b.GetTimePeriod(fields.Fd4TimePeriod[0], true)
    cf := b.CashFlow(fv, couponRate, cp, time, tp)
    //Yield to Maturity.
    if currentInterest {
      if cp != finances.Continuously {
        bondPrice = b.CurrentPrice(cf, curInterest, cp)
        fields.Fd4Result[0] = fmt.Sprintf("Yield to Maturity: %.5f%%", b.YieldToMaturity(cf, bondPrice, tp))
      } else {
        bondPrice = b.CurrentPriceContinuous(cf, curInterest)
        fields.Fd4Result[0] = fmt.Sprintf("Yield to Maturity: %.5f%%", b.YieldToMaturityContinuous(cf, bondPrice))
      }
    } else {  //Bond price.
      if cp != finances.Continuously {
        fields.Fd4Result[0] = fmt.Sprintf("Yield to Maturity: %.5f%%", b.YieldToMaturity(cf, bondPrice, cp))
      } else {
        fields.Fd4Result[0] = fmt.Sprintf("Yield to Maturity: %.5f%%", b.YieldToMaturityContinuous(cf, bondPrice))
      }
    }
    //Current Yield.
    var a finances.Annuities
    var annualRate float64
    switch fields.Fd4Compound[0] {
    case 'a', 'A':
      annualRate = couponRate
    default:
      annualRate = a.CompoundingFrequencyConversion(couponRate / 100.0,
        a.GetCompoundingPeriod(fields.Fd4Compound[0], true), a.GetCompoundingPeriod('a', true)) * 100.0
    }
    fields.Fd4Result[1] = fmt.Sprintf("Current Yield: %.5f%%", b.CurrentYield(annualRate, fv, bondPrice) * 100.0)
  }
  logger.LogInfo(fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur radio = %s, cur interest = %s, bond price = %s, %s",
    fields.Fd4FaceValue, fields.Fd4Time, fields.Fd4TimePeriod, fields.Fd4Coupon, fields.Fd4Compound, fields.Fd4CurrentRadio,
    fields.Fd4CurInterest, fields.Fd4BondPrice, fields.Fd4Result), correlationId)
}

//Extraction helper for UI-5 calculations.
func (b WfBondsPages) processUi5Form(req *http.Request, fields *bondsFields, correlationId string) {
  fields.Fd5FaceValue = req.PostFormValue("fd5-facevalue")
  fields.Fd5Time = req.PostFormValue("fd5-time")
  fields.Fd5TimePeriod = req.PostFormValue("fd5-tp")
  fields.Fd5Coupon = req.PostFormValue("fd5-coupon")
  fields.Fd5CompoundCoupon = req.PostFormValue("fd5-compound-coupon")
  fields.Fd5CurInterest = req.PostFormValue("fd5-current")
  fields.Fd5Compound = req.PostFormValue("fd5-compound")
  var fv float64
  var time float64
  var coupon float64
  var current float64
  var err error
  if fv, err = strconv.ParseFloat(fields.Fd5FaceValue, 64); err != nil {
    fields.Fd5Result[0] = fmt.Sprintf("Error: %s -- %+v", fields.Fd5FaceValue, err)
  } else if time, err = strconv.ParseFloat(fields.Fd5Time, 64); err != nil {
    fields.Fd5Result[0] = fmt.Sprintf("Error: %s -- %+v", fields.Fd5Time, err)
  } else if coupon, err = strconv.ParseFloat(fields.Fd5Coupon, 64); err != nil {
    fields.Fd5Result[0] = fmt.Sprintf("Error: %s -- %+v", fields.Fd5Coupon, err)
  } else if current, err = strconv.ParseFloat(fields.Fd5CurInterest, 64); err != nil {
    fields.Fd5Result[0] = fmt.Sprintf("Error: %s -- %+v", fields.Fd5CurInterest, err)
  } else {
    var b finances.Bonds
    var cp int = b.GetCompoundingPeriod(fields.Fd5Compound[0], true)
    cf := b.CashFlow(fv, coupon, b.GetCompoundingPeriod(fields.Fd5CompoundCoupon[0], true), time,
      b.GetTimePeriod(fields.Fd5TimePeriod[0], true))
    //Duration.
    switch fields.Fd5Compound[0] {
    case 'c', 'C':
      fields.Fd5Result[0] = fmt.Sprintf("Duration: %.5f", b.DurationContinuous(cf, current, b.CurrentPriceContinuous(cf, current)))
    default:
      fields.Fd5Result[0] = fmt.Sprintf("Duration: %.5f", b.Duration(cf, cp, current, b.CurrentPrice(cf, current, cp)))
    }
    //Macaulay Duration.
    if len(fields.Fd5Result[2]) == 0 {
      fields.Fd5Result[1] = bond_notes[0]
    }
    switch fields.Fd5Compound[0] {
    case 'c', 'C':
      fields.Fd5Result[2] = fmt.Sprintf("Macaulay Duration: %.5f year(s)",
        b.MacaulayDurationContinuous(cf, b.CurrentPriceContinuous(cf, current)))
    default:
      fields.Fd5Result[2] = fmt.Sprintf("Macaulay Duration: %.5f year(s)",
        b.MacaulayDuration(cf, b.GetCompoundingPeriod(fields.Fd5CompoundCoupon[0], true),
        b.CurrentPrice(cf, current, b.GetCompoundingPeriod(fields.Fd5Compound[0], true))))
    }
    //Modified Duration.
    if len(fields.Fd5Result[4]) == 0 {
      fields.Fd5Result[3] = bond_notes[1]
    }
    fields.Fd5Result[4] = fmt.Sprintf("Modified Duration: %.5f%%",
      b.ModifiedDuration(cf, b.GetCompoundingPeriod(fields.Fd5CompoundCoupon[0], true),
      b.CurrentPrice(cf, current, b.GetCompoundingPeriod(fields.Fd5Compound[0], true))))
    //Convexity.
    if len(fields.Fd5Result[6]) == 0 {
      fields.Fd5Result[5] = bond_notes[2]
    }
    switch fields.Fd5Compound[0] {
    case 'c', 'C':
      fields.Fd5Result[6] = fmt.Sprintf("Convexity: %.5f", b.ConvexityContinuous(cf, current, b.CurrentPriceContinuous(cf, current)))
    default:
      fields.Fd5Result[6] = fmt.Sprintf("Convexity: %.5f", b.Convexity(cf, current,
        b.GetCompoundingPeriod(fields.Fd5Compound[0], true)))
    }
  }
  logger.LogInfo(fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s", fields.Fd5FaceValue,
    fields.Fd5Time, fields.Fd5TimePeriod, fields.Fd5Coupon, fields.Fd5Compound, fields.Fd5CurInterest, fields.Fd5Result), correlationId)
}
