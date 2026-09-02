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
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res, correlationId)
    return
  }
  //
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    fields := getBondsFields(userName)
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
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/bonds/bonds.html",
        "webfinances/templates/finances/bonds/taxfree.html",
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
          Fd1TaxFree string
          Fd1CityTax string
          Fd1StateTax string
          Fd1FederalTax string
          Fd1Result string
        } { "standard", "Bonds", logger.DatetimeFormat(), financesMenuPage, fields.CurrentButton, newSession.CsrfToken, fields.Fd1TaxFree,
            fields.Fd1CityTax, fields.Fd1StateTax, fields.Fd1FederalTax, fields.Fd1Result },
      })
    } else if strings.EqualFold(fields.CurrentPage, "rhs-ui2") {
      fields.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {
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
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/bonds/bonds.html",
        "webfinances/templates/finances/bonds/currentprice.html",
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
          Fd2FaceValue string
          Fd2Time string
          Fd2TimePeriod string
          Fd2Coupon string
          Fd2CompoundCoupon string
          Fd2Current string
          Fd2Compound string
          Fd2Result string
        } { "standard", "Bonds", logger.DatetimeFormat(), financesMenuPage, fields.CurrentButton, newSession.CsrfToken,
            fields.Fd2FaceValue, fields.Fd2Time, fields.Fd2TimePeriod, fields.Fd2Coupon, fields.Fd2CompoundCoupon,
            fields.Fd2Current, fields.Fd2Compound, fields.Fd2Result },
      })
    } else if strings.EqualFold(fields.CurrentPage, "rhs-ui3") {
      fields.CurrentButton = "lhs-button3"
      if req.Method == http.MethodPost {
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
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/bonds/bonds.html",
        "webfinances/templates/finances/bonds/yieldtocall.html",
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
          Fd3FaceValue string
          Fd3TimeCall string
          Fd3TimePeriod string
          Fd3Coupon string
          Fd3Compound string
          Fd3BondPrice string
          Fd3CallPrice string
          Fd3Result string
        } { "standard", "Bonds", logger.DatetimeFormat(), financesMenuPage, fields.CurrentButton, newSession.CsrfToken,
            fields.Fd3FaceValue, fields.Fd3TimeCall, fields.Fd3TimePeriod, fields.Fd3Coupon, fields.Fd3Compound,
            fields.Fd3BondPrice, fields.Fd3CallPrice, fields.Fd3Result },
      })
    } else if strings.EqualFold(fields.CurrentPage, "rhs-ui4") {
      fields.CurrentButton = "lhs-button4"
      if req.Method == http.MethodPost {
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
          fields.Fd4FaceValue, fields.Fd4Time, fields.Fd4TimePeriod, fields.Fd4Coupon, fields.Fd4Compound, fields.Fd4CurrentRadio, fields.Fd4CurInterest, fields.Fd4BondPrice,
          fields.Fd4Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/bonds/bonds.html",
        "webfinances/templates/finances/bonds/yieldtomaturity.html",
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
          Fd4FaceValue string
          Fd4Time string
          Fd4TimePeriod string
          Fd4Coupon string
          Fd4Compound string
          Fd4CurrentRadio string
          Fd4CurInterest string
          Fd4BondPrice string
          Fd4Result [2]string
        } { "standard", "Bonds", logger.DatetimeFormat(), financesMenuPage, fields.CurrentButton, newSession.CsrfToken,
            fields.Fd4FaceValue, fields.Fd4Time, fields.Fd4TimePeriod, fields.Fd4Coupon, fields.Fd4Compound, fields.Fd4CurrentRadio,
            fields.Fd4CurInterest, fields.Fd4BondPrice, fields.Fd4Result },
      })
    } else if strings.EqualFold(fields.CurrentPage, "rhs-ui5") {
      fields.CurrentButton = "lhs-button5"
      if req.Method == http.MethodPost {
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
          fields.Fd5Time, fields.Fd5TimePeriod, fields.Fd5Coupon, fields.Fd5Compound, fields.Fd5CurInterest, fields.Fd5Result),
          correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/bonds/bonds.html",
        "webfinances/templates/finances/bonds/duration.html",
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
          Fd5FaceValue string
          Fd5Time string
          Fd5TimePeriod string
          Fd5Coupon string
          Fd5CompoundCoupon string
          Fd5CurInterest string
          Fd5Compound string
          Fd5Result [7]string
        } { "standard", "Bonds", logger.DatetimeFormat(), financesMenuPage, fields.CurrentButton, newSession.CsrfToken,
            fields.Fd5FaceValue, fields.Fd5Time, fields.Fd5TimePeriod, fields.Fd5Coupon, fields.Fd5CompoundCoupon,
            fields.Fd5CurInterest, fields.Fd5Compound, fields.Fd5Result },
      })
    /***
    // } else if strings.EqualFold(currentRHS, "rhs-ui6") {
    //   fields.CurrentButton = "lhs-button6"
    //   if req.Method == http.MethodPost {
    //     fields.Fd6FaceValue = req.PostFormValue("fd6-facevalue")
    //     fields.Fd6Time = req.PostFormValue("fd6-time")
    //     fields.Fd6TimePeriod = req.PostFormValue("fd6-tp")
    //     fields.Fd6Coupon = req.PostFormValue("fd6-coupon")
    //     fields.Fd6CompoundCoupon = req.PostFormValue("fd6-compound-coupon")
    //     fields.Fd6CurInterest = req.PostFormValue("fd6-current")
    //     fields.Fd6Compound = req.PostFormValue("fd6-compound")
    //     var fv float64
    //     var time float64
    //     var couponRate float64
    //     var current float64
    //     var err error
    //     if fv, err = strconv.ParseFloat(fields.Fd6FaceValue, 64); err != nil {
    //       fields.Fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", fields.Fd6FaceValue, err)
    //     } else if time, err = strconv.ParseFloat(fields.Fd6Time, 64); err != nil {
    //       fields.Fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", fields.Fd6Time, err)
    //     } else if couponRate, err = strconv.ParseFloat(fields.Fd6Coupon, 64); err != nil {
    //       fields.Fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", fields.Fd6Coupon, err)
    //     } else if current, err = strconv.ParseFloat(fields.Fd6CurInterest, 64); err != nil {
    //       fields.Fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", fields.Fd6CurInterest, err)
    //     } else {
    //       var b finances.Bonds
    //       cf := b.CashFlow(fv, couponRate, b.GetCompoundingPeriod(fields.Fd6CompoundCoupon[0], true),
    //         time, b.GetTimePeriod(fields.Fd6TimePeriod[0], true))
    //       switch fields.Fd6Compound[0] {
    //       case 'c', 'C':
    //         fields.Fd6Result[1] = fmt.Sprintf("Macaulay Duration: %.3f year(s)",
    //           b.MacaulayDurationContinuous(cf, b.CurrentPriceContinuous(cf, current)))
    //       default:
    //         fields.Fd6Result[1] = fmt.Sprintf("Macaulay Duration: %.3f year(s)",
    //           b.MacaulayDuration(cf, b.GetCompoundingPeriod(fields.Fd6CompoundCoupon[0], true),
    //             b.CurrentPrice(cf, current, b.GetCompoundingPeriod(fields.Fd6Compound[0], true))))
    //       }
    //     }
    //     logEntry.Print(INFO, correlationId, []string {
    //       fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s",
    //         fields.Fd6FaceValue, fields.Fd6Time, fields.Fd6TimePeriod, fields.Fd6Coupon, fields.Fd6Compound,
    //         fields.Fd6CurInterest, fields.Fd6Result[1]),
    //     })
    //   }
    //   newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    //   cookie := sessions.CreateCookie(newSessionToken)
    //   http.SetCookie(res, cookie)
    //   t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
    //     "webfinances/templates/header.html",
    //     "webfinances/templates/bonds/macaulayduration.html",
    //     "webfinances/templates/footer.html"))
    //   t.ExecuteTemplate(res, "bonds", struct {
    //     Header string
    //     Datetime string
    //     CurrentButton string
    //     CsrfToken string
    //     Fd6FaceValue string
    //     Fd6Time string
    //     Fd6TimePeriod string
    //     Fd6Coupon string
    //     Fd6CompoundCoupon string
    //     Fd6CurInterest string
    //     Fd6Compound string
    //     Fd6Result [2]string
    //   } { "Bonds", m.DTF(), fields.CurrentButton, newSession.CsrfToken, fields.Fd6FaceValue, fields.Fd6Time,
    //       fields.Fd6TimePeriod, fields.Fd6Coupon, fields.Fd6CompoundCoupon, fields.Fd6CurInterest, fields.Fd6Compound,
    //       fields.Fd6Result,
    //     })
    // } else if strings.EqualFold(currentRHS, "rhs-ui7") {
    //   fields.CurrentButton = "lhs-button7"
    //   if req.Method == http.MethodPost {
    //     fields.Fd7FaceValue = req.PostFormValue("fd7-facevalue")
    //     fields.Fd7Time = req.PostFormValue("fd7-time")
    //     fields.Fd7TimePeriod = req.PostFormValue("fd7-tp")
    //     fields.Fd7Coupon = req.PostFormValue("fd7-coupon")
    //     fields.Fd7CompoundCoupon = req.PostFormValue("fd7-compound-coupon")
    //     fields.Fd7CurInterest = req.PostFormValue("fd7-current")
    //     fields.Fd7Compound = req.PostFormValue("fd7-compound")
    //     var fv float64
    //     var time float64
    //     var couponRate float64
    //     var current float64
    //     var err error
    //     if fv, err = strconv.ParseFloat(fields.Fd7FaceValue, 64); err != nil {
    //       fields.Fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", fields.Fd7FaceValue, err)
    //     } else if time, err = strconv.ParseFloat(fields.Fd7Time, 64); err != nil {
    //       fields.Fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", fields.Fd7Time, err)
    //     } else if couponRate, err = strconv.ParseFloat(fields.Fd7Coupon, 64); err != nil {
    //       fields.Fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", fields.Fd7Coupon, err)
    //     } else if current, err = strconv.ParseFloat(fields.Fd7CurInterest, 64); err != nil {
    //       fields.Fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", fields.Fd7CurInterest, err)
    //     } else {
    //       var b finances.Bonds
    //       cf := b.CashFlow(fv, couponRate, b.GetCompoundingPeriod(fields.Fd7CompoundCoupon[0], true),
    //         time, b.GetTimePeriod(fields.Fd7TimePeriod[0], true))
    //       fields.Fd7Result[1] = fmt.Sprintf("Modified Duration: %.3f%%",
    //         b.ModifiedDuration(cf, b.GetCompoundingPeriod(fields.Fd7CompoundCoupon[0], true),
    //         b.CurrentPrice(cf, current, b.GetCompoundingPeriod(fields.Fd7Compound[0], true))))
    //     }
    //     logEntry.Print(INFO, correlationId, []string {
    //       fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s",
    //         fields.Fd7FaceValue, fields.Fd7Time, fields.Fd7TimePeriod, fields.Fd7Coupon, fields.Fd7Compound,
    //         fields.Fd7CurInterest, fields.Fd7Result[1]),
    //     })
    //   }
    //   newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    //   cookie := sessions.CreateCookie(newSessionToken)
    //   http.SetCookie(res, cookie)
    //   t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
    //     "webfinances/templates/header.html",
    //     "webfinances/templates/bonds/modifiedduration.html",
    //     "webfinances/templates/footer.html"))
    //   t.ExecuteTemplate(res, "bonds", struct {
    //     Header string
    //     Datetime string
    //     CurrentButton string
    //     CsrfToken string
    //     Fd7FaceValue string
    //     Fd7Time string
    //     Fd7TimePeriod string
    //     Fd7Coupon string
    //     Fd7CompoundCoupon string
    //     Fd7CurInterest string
    //     Fd7Compound string
    //     Fd7Result [2]string
    //   } { "Bonds", m.DTF(), fields.CurrentButton, newSession.CsrfToken, fields.Fd7FaceValue, fields.Fd7Time,
    //       fields.Fd7TimePeriod, fields.Fd7Coupon, fields.Fd7CompoundCoupon, fields.Fd7CurInterest, fields.Fd7Compound,
    //       fields.Fd7Result,
    //     })
    // } else if strings.EqualFold(currentRHS, "rhs-ui8") {
    //   fields.CurrentButton = "lhs-button8"
    //   if req.Method == http.MethodPost {
    //     fields.Fd8FaceValue = req.PostFormValue("fd8-facevalue")
    //     fields.Fd8Time = req.PostFormValue("fd8-time")
    //     fields.Fd8TimePeriod = req.PostFormValue("fd8-tp")
    //     fields.Fd8Coupon = req.PostFormValue("fd8-coupon")
    //     fields.Fd8CompoundCoupon = req.PostFormValue("fd8-compound-coupon")
    //     fields.Fd8CurInterest = req.PostFormValue("fd8-current")
    //     fields.Fd8Compound = req.PostFormValue("fd8-compound")
    //     var fv float64
    //     var time float64
    //     var couponRate float64
    //     var current float64
    //     var err error
    //     if fv, err = strconv.ParseFloat(fields.Fd8FaceValue, 64); err != nil {
    //       fields.Fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", fields.Fd8FaceValue, err)
    //     } else if time, err = strconv.ParseFloat(fields.Fd8Time, 64); err != nil {
    //       fields.Fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", fields.Fd8Time, err)
    //     } else if couponRate, err = strconv.ParseFloat(fields.Fd8Coupon, 64); err != nil {
    //       fields.Fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", fields.Fd8Coupon, err)
    //     } else if current, err = strconv.ParseFloat(fields.Fd8CurInterest, 64); err != nil {
    //       fields.Fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", fields.Fd8CurInterest, err)
    //     } else {
    //       var b finances.Bonds
    //       cf := b.CashFlow(fv, couponRate, b.GetCompoundingPeriod(fields.Fd8CompoundCoupon[0], true),
    //         time, b.GetTimePeriod(fields.Fd8TimePeriod[0], true))
    //       switch fields.Fd8Compound[0] {
    //       case 'c', 'C':
    //         fields.Fd8Result[1] = fmt.Sprintf("Convexity: %.3f", b.ConvexityContinuous(cf, current,
    //           b.CurrentPriceContinuous(cf, current)))
    //       default:
    //         fields.Fd8Result[1] = fmt.Sprintf("Convexity: %.3f", b.Convexity(cf, current,
    //           b.GetCompoundingPeriod(fields.Fd8Compound[0], true)))
    //       }
    //     }
    //     logEntry.Print(INFO, correlationId, []string {
    //       fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s",
    //         fields.Fd8FaceValue, fields.Fd8Time, fields.Fd8TimePeriod, fields.Fd8Coupon, fields.Fd8Compound,
    //         fields.Fd8CurInterest, fields.Fd8Result[1]),
    //     })
    //   }
    //   newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    //   cookie := sessions.CreateCookie(newSessionToken)
    //   http.SetCookie(res, cookie)
    //   t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
    //     "webfinances/templates/header.html",
    //     "webfinances/templates/bonds/convexity.html",
    //     "webfinances/templates/footer.html"))
    //   t.ExecuteTemplate(res, "bonds", struct {
    //     Header string
    //     Datetime string
    //     CurrentButton string
    //     CsrfToken string
    //     Fd8FaceValue string
    //     Fd8Time string
    //     Fd8TimePeriod string
    //     Fd8Coupon string
    //     Fd8CompoundCoupon string
    //     Fd8CurInterest string
    //     Fd8Compound string
    //     Fd8Result [2]string
    //   } { "Bonds", m.DTF(), fields.CurrentButton, newSession.CsrfToken, fields.Fd8FaceValue, fields.Fd8Time,
    //       fields.Fd8TimePeriod, fields.Fd8Coupon, fields.Fd8CompoundCoupon, fields.Fd8CurInterest, fields.Fd8Compound,
    //       fields.Fd8Result,
    //     })
    ***/
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
      } else if strings.EqualFold(fields.CurrentPage, "rhs-ui3") {
        fields.Fd3Result = ""
      } else if strings.EqualFold(fields.CurrentPage, "rhs-ui4") {
        fields.Fd4Result[0] = ""
        fields.Fd4Result[1] = ""
      } else if strings.EqualFold(fields.CurrentPage, "rhs-ui5") {
        fields.Fd5Result[0] = ""
        fields.Fd5Result[1] = ""
        fields.Fd5Result[2] = ""
        fields.Fd5Result[3] = ""
        fields.Fd5Result[4] = ""
        fields.Fd5Result[5] = ""
        fields.Fd5Result[6] = ""
        /***
      // } else if strings.EqualFold(currentRHS, "rhs-ui6") {
      //   fields.Fd6Result[1] = ""
      // } else if strings.EqualFold(currentRHS, "rhs-ui7") {
      //   fields.Fd7Result[1] = ""
      // } else if strings.EqualFold(currentRHS, "rhs-ui8") {
      //   fields.Fd8Result[1] = ""
        ***/
      }
    }
    //
    if data, err := json.Marshal(fields); err != nil {
      logger.LogError(fmt.Sprintf("%+v", err), correlationId)
    } else {
      filePath := fmt.Sprintf("%s/%s/bonds.txt", mainDir, userName)
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
