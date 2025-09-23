package webfinances

import (
  "context"
  "encoding/json"
  "finance/finances"
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/gposu"
  "github.com/juan-carlos-trimino/gpsessions"
  "html/template"
  "math"
  "net/http"
  "os"
  "strconv"
  "strings"
  "time"
)

var bond_notes = [...]string {
  "The Macaulay duration is a measure of a bond's sensitivity to interest rate changes. The " +
  "duration is the weighed-average number of years the investor must hold a bond until the " +
  "present value of the bond's cash flows equals the amount paid for the bond.",
  "The modified duration of a bond is a measure of the sensitivity of the bond's price to " +
  "changes in interest rates. Since bond prices move in an inverse direction from interest " +
  "rates, for a one percent increase (decrease) in interest rates, the bond's price will " +
  "decrease (increase) by the percentage shown by the modified duration.",
  "Convexity in bonds measures how sensitive the bond's duration is to changes in interest " +
  "rates. The higher the convexity, the less the bond price will increase when rates fall -- " +
  "and the less the bond price will drop when rates rise.",
}

type WfBondsPages struct {
}

func (b WfBondsPages) BondsPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
    return
  }
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logger.LogInfo("Entering BondsPages/webfinances.", correlationId)
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    bf := getBondsFields(userName)
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
      bf.CurrentPage = ui
    }
    //
    if strings.EqualFold(bf.CurrentPage, "rhs-ui1") {
      bf.CurrentButton = "lhs-button1"
      if req.Method == http.MethodPost {
        bf.Fd1TaxFree = req.PostFormValue("fd1-taxfree")
        bf.Fd1CityTax = req.PostFormValue("fd1-citytax")
        bf.Fd1StateTax = req.PostFormValue("fd1-statetax")
        bf.Fd1FederalTax = req.PostFormValue("fd1-federaltax")
        var taxFree float64
        var cityTax float64
        var stateTax float64
        var federalTax float64
        var err error
        if taxFree, err = strconv.ParseFloat(bf.Fd1TaxFree, 64); err != nil {
          bf.Fd1Result = fmt.Sprintf("Error: %s -- %+v", bf.Fd1TaxFree, err)
        } else if cityTax, err = strconv.ParseFloat(bf.Fd1CityTax, 64); err != nil {
          bf.Fd1Result = fmt.Sprintf("Error: %s -- %+v", bf.Fd1CityTax, err)
        } else if stateTax, err = strconv.ParseFloat(bf.Fd1StateTax, 64); err != nil {
          bf.Fd1Result = fmt.Sprintf("Error: %s -- %+v", bf.Fd1StateTax, err)
        } else if federalTax, err = strconv.ParseFloat(bf.Fd1FederalTax, 64); err != nil {
          bf.Fd1Result = fmt.Sprintf("Error: %s -- %+v", bf.Fd1FederalTax, err)
        } else {
          var b finances.Bonds
          bf.Fd1Result = fmt.Sprintf("Taxable-Equivalent Yield: %.3f%%",
            b.TaxableVsTaxFreeYields(taxFree, cityTax, stateTax, federalTax) * 100.0)
        }
        logger.LogInfo(fmt.Sprintf(
         "tax free = %s, city tax = %s, state tax = %s, federal tax = %s, %s", bf.Fd1TaxFree,
         bf.Fd1CityTax, bf.Fd1StateTax, bf.Fd1FederalTax, bf.Fd1Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
        "webfinances/templates/header.html",
        "webfinances/templates/bonds/taxfree.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd1TaxFree string
        Fd1CityTax string
        Fd1StateTax string
        Fd1FederalTax string
        Fd1Result string
      } { "Bonds", logger.DatetimeFormat(), bf.CurrentButton, newSession.CsrfToken, bf.Fd1TaxFree,
          bf.Fd1CityTax, bf.Fd1StateTax, bf.Fd1FederalTax, bf.Fd1Result,
        })
    } else if strings.EqualFold(bf.CurrentPage, "rhs-ui2") {
      bf.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        bf.Fd2FaceValue = req.FormValue("fd2-facevalue")
        bf.Fd2Time = req.PostFormValue("fd2-time")
        bf.Fd2TimePeriod = req.PostFormValue("fd2-tp")
        bf.Fd2Coupon = req.PostFormValue("fd2-coupon")
        bf.Fd2CompoundCoupon = req.PostFormValue("fd2-compound-coupon")
        bf.Fd2Current = req.PostFormValue("fd2-current")
        bf.Fd2Compound = req.PostFormValue("fd2-compound")
        var fv float64
        var time float64
        var coupon float64
        var current float64
        var err error
        if fv, err = strconv.ParseFloat(bf.Fd2FaceValue, 64); err != nil {
          bf.Fd2Result = fmt.Sprintf("Error: %s -- %+v", bf.Fd2FaceValue, err)
        } else if time, err = strconv.ParseFloat(bf.Fd2Time, 64); err != nil {
          bf.Fd2Result = fmt.Sprintf("Error: %s -- %+v", bf.Fd2Time, err)
        } else if coupon, err = strconv.ParseFloat(bf.Fd2Coupon, 64); err != nil {
          bf.Fd2Result = fmt.Sprintf("Error: %s -- %+v", bf.Fd2Coupon, err)
        } else if current, err = strconv.ParseFloat(bf.Fd2Current, 64); err != nil {
          bf.Fd2Result = fmt.Sprintf("Error: %s -- %+v", bf.Fd2Current, err)
        } else {
          var b finances.Bonds
          cf := b.CashFlow(fv, coupon, b.GetCompoundingPeriod(bf.Fd2CompoundCoupon[0], true), time,
            b.GetTimePeriod(bf.Fd2TimePeriod[0], true))
          var currentPrice float64
          switch bf.Fd2Compound[0] {
          case 'c', 'C':
            currentPrice = b.CurrentPriceContinuous(cf, current)
          default:
            currentPrice = b.CurrentPrice(cf, current, b.GetCompoundingPeriod(bf.Fd2Compound[0],
              true))
          }
          //
          if math.Abs(fv - currentPrice) < finances.Accuracy {
            bf.Fd2Result = fmt.Sprintf("Current Price: $%.2f (par)", currentPrice)
          } else if fv < currentPrice {
            bf.Fd2Result = fmt.Sprintf("Current Price: $%.2f (premium)", currentPrice)
          } else {
            bf.Fd2Result = fmt.Sprintf("Current Price: $%.2f (discount)", currentPrice)
          }
        }
        logger.LogInfo(fmt.Sprintf(
         "fv = %s, time = %s, tp = %s, coupon rate = %s, current interest = %s, cp = %s, %s",
         bf.Fd2FaceValue, bf.Fd2Time, bf.Fd2TimePeriod, bf.Fd2Coupon, bf.Fd2Current,
         bf.Fd2Compound, bf.Fd2Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
        "webfinances/templates/header.html",
        "webfinances/templates/bonds/currentprice.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
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
      } { "Bonds", logger.DatetimeFormat(), bf.CurrentButton, newSession.CsrfToken,
          bf.Fd2FaceValue, bf.Fd2Time, bf.Fd2TimePeriod, bf.Fd2Coupon, bf.Fd2CompoundCoupon,
          bf.Fd2Current, bf.Fd2Compound, bf.Fd2Result,
        })
    } else if strings.EqualFold(bf.CurrentPage, "rhs-ui3") {
      bf.CurrentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        bf.Fd3FaceValue = req.PostFormValue("fd3-facevalue")
        bf.Fd3TimeCall = req.PostFormValue("fd3-timecall")
        bf.Fd3TimePeriod = req.PostFormValue("fd3-tp")
        bf.Fd3Coupon = req.PostFormValue("fd3-coupon")
        bf.Fd3BondPrice = req.PostFormValue("fd3-bondprice")
        bf.Fd3CallPrice = req.PostFormValue("fd3-callprice")
        bf.Fd3Compound = req.PostFormValue("fd3-compound")
        var fv float64
        var timeToCall float64
        var couponRate float64
        var bondPrice float64
        var callPrice float64
        var err error
        if fv, err = strconv.ParseFloat(bf.Fd3FaceValue, 64); err != nil {
          bf.Fd3Result = fmt.Sprintf("Error: %s -- %+v", bf.Fd3FaceValue, err)
        } else if timeToCall, err = strconv.ParseFloat(bf.Fd3TimeCall, 64); err != nil {
          bf.Fd3Result = fmt.Sprintf("Error: %s -- %+v", bf.Fd3TimeCall, err)
        } else if couponRate, err = strconv.ParseFloat(bf.Fd3Coupon, 64); err != nil {
          bf.Fd3Result = fmt.Sprintf("Error: %s -- %+v", bf.Fd3Coupon, err)
        } else if bondPrice, err = strconv.ParseFloat(bf.Fd3BondPrice, 64); err != nil {
          bf.Fd3Result = fmt.Sprintf("Error: %s -- %+v", bf.Fd3BondPrice, err)
        } else if callPrice, err = strconv.ParseFloat(bf.Fd3CallPrice, 64); err != nil {
          bf.Fd3Result = fmt.Sprintf("Error: %s -- %+v", bf.Fd3CallPrice, err)
        } else {
          var b finances.Bonds
          bf.Fd3Result = fmt.Sprintf("Yield to Call: %.3f%%", b.YieldToCall(fv, couponRate,
            b.GetCompoundingPeriod(bf.Fd3Compound[0], true), timeToCall,
            b.GetTimePeriod(bf.Fd3TimePeriod[0], true), bondPrice, callPrice))
        }
        logger.LogInfo(fmt.Sprintf(
         "fv = %s, coupon rate = %s, cp = %s, time to call = %s, tp = %s, bond price = %s, call price = %s, %s",
         bf.Fd3FaceValue, bf.Fd3Coupon, bf.Fd3Compound, bf.Fd3TimeCall, bf.Fd3TimePeriod,
         bf.Fd3BondPrice, bf.Fd3CallPrice, bf.Fd3Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
        "webfinances/templates/header.html",
        "webfinances/templates/bonds/yieldtocall.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
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
      } { "Bonds", logger.DatetimeFormat(), bf.CurrentButton, newSession.CsrfToken, bf.Fd3FaceValue,
          bf.Fd3TimeCall, bf.Fd3TimePeriod, bf.Fd3Coupon, bf.Fd3Compound, bf.Fd3BondPrice,
          bf.Fd3CallPrice, bf.Fd3Result,
        })
    } else if strings.EqualFold(bf.CurrentPage, "rhs-ui4") {
      bf.CurrentButton = "lhs-button4"
      if req.Method == http.MethodPost {
        bf.Fd4FaceValue = req.PostFormValue("fd4-facevalue")
        bf.Fd4Time = req.PostFormValue("fd4-time")
        bf.Fd4TimePeriod = req.PostFormValue("fd4-tp")
        bf.Fd4Coupon = req.PostFormValue("fd4-coupon")
        bf.Fd4Compound = req.PostFormValue("fd4-compound")
        bf.Fd4CurrentRadio = req.PostFormValue("fd4-choice")
        bf.Fd4CurInterest = req.PostFormValue("fd4-ci")
        bf.Fd4BondPrice = req.PostFormValue("fd4-bp")
        var currentInterest bool = false
        if strings.EqualFold(bf.Fd4CurrentRadio, "fd4-curinterest") {
          currentInterest = true
        }
        var fv float64
        var time float64
        var couponRate float64
        var curInterest float64
        var bondPrice float64
        var err error
        if fv, err = strconv.ParseFloat(bf.Fd4FaceValue, 64); err != nil {
          bf.Fd4Result[0] = fmt.Sprintf("Error: %s -- %+v", bf.Fd4FaceValue, err)
        } else if time, err = strconv.ParseFloat(bf.Fd4Time, 64); err != nil {
          bf.Fd4Result[0] = fmt.Sprintf("Error: %s -- %+v", bf.Fd4Time, err)
        } else if couponRate, err = strconv.ParseFloat(bf.Fd4Coupon, 64); err != nil {
          bf.Fd4Result[0] = fmt.Sprintf("Error: %s -- %+v", bf.Fd4Coupon, err)
        } else if curInterest, err = strconv.ParseFloat(bf.Fd4CurInterest, 64); err != nil {
          bf.Fd4Result[0] = fmt.Sprintf("Error: %s -- %+v", bf.Fd4CurInterest, err)
        } else if bondPrice, err = strconv.ParseFloat(bf.Fd4BondPrice, 64); err != nil {
          bf.Fd4Result[0] = fmt.Sprintf("Error: %s -- %+v", bf.Fd4BondPrice, err)
        } else {
          var b finances.Bonds
          var cp int = b.GetCompoundingPeriod(bf.Fd4Compound[0], true)
          var tp = b.GetTimePeriod(bf.Fd4TimePeriod[0], true)
          cf := b.CashFlow(fv, couponRate, cp, time, tp)
          //Yield to Maturity.
          if currentInterest {
            if cp != finances.Continuously {
              bondPrice = b.CurrentPrice(cf, curInterest, cp)
              bf.Fd4Result[0] = fmt.Sprintf("Yield to Maturity: %.3f%%",
                b.YieldToMaturity(cf, bondPrice, tp))
            } else {
              bondPrice = b.CurrentPriceContinuous(cf, curInterest)
              bf.Fd4Result[0] = fmt.Sprintf("Yield to Maturity: %.3f%%",
                b.YieldToMaturityContinuous(cf, bondPrice))
            }
          } else {  //Bond price.
            if cp != finances.Continuously {
              bf.Fd4Result[0] = fmt.Sprintf("Yield to Maturity: %.3f%%",
                b.YieldToMaturity(cf, bondPrice, cp))
            } else {
              bf.Fd4Result[0] = fmt.Sprintf("Yield to Maturity: %.3f%%",
                b.YieldToMaturityContinuous(cf, bondPrice))
            }
          }
          //Current Yield.
          var a finances.Annuities
          var annualRate float64
          switch bf.Fd4Compound[0] {
          case 'a', 'A':
            annualRate = couponRate
          default:
            annualRate = a.CompoundingFrequencyConversion(couponRate / 100.0,
              a.GetCompoundingPeriod(bf.Fd4Compound[0], true), a.GetCompoundingPeriod('a', true)) * 100.0
          }
          bf.Fd4Result[1] = fmt.Sprintf("Current Yield = %.3f%%", b.CurrentYield(annualRate, fv, bondPrice) * 100.0)
        }
        logger.LogInfo(fmt.Sprintf(
         "fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur radio = %s, cur interest = %s, bond price = %s, %s",
         bf.Fd4FaceValue, bf.Fd4Time, bf.Fd4TimePeriod, bf.Fd4Coupon, bf.Fd4Compound,
         bf.Fd4CurrentRadio, bf.Fd4CurInterest, bf.Fd4BondPrice, bf.Fd4Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
        "webfinances/templates/header.html",
        "webfinances/templates/bonds/yieldtomaturity.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
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
      } { "Bonds", logger.DatetimeFormat(), bf.CurrentButton, newSession.CsrfToken, bf.Fd4FaceValue,
          bf.Fd4Time, bf.Fd4TimePeriod, bf.Fd4Coupon, bf.Fd4Compound, bf.Fd4CurrentRadio,
          bf.Fd4CurInterest, bf.Fd4BondPrice, bf.Fd4Result,
        })
    } else if strings.EqualFold(bf.CurrentPage, "rhs-ui5") {
      bf.CurrentButton = "lhs-button5"
      if req.Method == http.MethodPost {
        bf.Fd5FaceValue = req.PostFormValue("fd5-facevalue")
        bf.Fd5Time = req.PostFormValue("fd5-time")
        bf.Fd5TimePeriod = req.PostFormValue("fd5-tp")
        bf.Fd5Coupon = req.PostFormValue("fd5-coupon")
        bf.Fd5CompoundCoupon = req.PostFormValue("fd5-compound-coupon")
        bf.Fd5CurInterest = req.PostFormValue("fd5-current")
        bf.Fd5Compound = req.PostFormValue("fd5-compound")
        var fv float64
        var time float64
        var coupon float64
        var current float64
        var err error
        if fv, err = strconv.ParseFloat(bf.Fd5FaceValue, 64); err != nil {
          bf.Fd5Result[0] = fmt.Sprintf("Error: %s -- %+v", bf.Fd5FaceValue, err)
        } else if time, err = strconv.ParseFloat(bf.Fd5Time, 64); err != nil {
          bf.Fd5Result[0] = fmt.Sprintf("Error: %s -- %+v", bf.Fd5Time, err)
        } else if coupon, err = strconv.ParseFloat(bf.Fd5Coupon, 64); err != nil {
          bf.Fd5Result[0] = fmt.Sprintf("Error: %s -- %+v", bf.Fd5Coupon, err)
        } else if current, err = strconv.ParseFloat(bf.Fd5CurInterest, 64); err != nil {
          bf.Fd5Result[0] = fmt.Sprintf("Error: %s -- %+v", bf.Fd5CurInterest, err)
        } else {
          var b finances.Bonds
          var cp int = b.GetCompoundingPeriod(bf.Fd5Compound[0], true)
          cf := b.CashFlow(fv, coupon, b.GetCompoundingPeriod(bf.Fd5CompoundCoupon[0], true), time,
            b.GetTimePeriod(bf.Fd5TimePeriod[0], true))
          //Duration.
          switch bf.Fd5Compound[0] {
          case 'c', 'C':
            bf.Fd5Result[0] = fmt.Sprintf("Duration: %.3f",
              b.DurationContinuous(cf, current, b.CurrentPriceContinuous(cf, current)))
          default:
            bf.Fd5Result[0] = fmt.Sprintf("Duration: %.3f",
              b.Duration(cf, cp, current, b.CurrentPrice(cf, current, cp)))
          }
          //Macaulay Duration.
          if len(bf.Fd5Result[2]) == 0 {
            bf.Fd5Result[1] = bond_notes[0]
          }
          switch bf.Fd5Compound[0] {
          case 'c', 'C':
            bf.Fd5Result[2] = fmt.Sprintf("Macaulay Duration: %.3f year(s)",
              b.MacaulayDurationContinuous(cf, b.CurrentPriceContinuous(cf, current)))
          default:
            bf.Fd5Result[2] = fmt.Sprintf("Macaulay Duration: %.3f year(s)",
              b.MacaulayDuration(cf, b.GetCompoundingPeriod(bf.Fd5CompoundCoupon[0], true),
                b.CurrentPrice(cf, current, b.GetCompoundingPeriod(bf.Fd5Compound[0], true))))
          }
          //Modified Duration.
          if len(bf.Fd5Result[4]) == 0 {
            bf.Fd5Result[3] = bond_notes[1]
          }
          bf.Fd5Result[4] = fmt.Sprintf("Modified Duration: %.3f%%",
            b.ModifiedDuration(cf, b.GetCompoundingPeriod(bf.Fd5CompoundCoupon[0], true),
            b.CurrentPrice(cf, current, b.GetCompoundingPeriod(bf.Fd5Compound[0], true))))
          //Convexity.
          if len(bf.Fd5Result[6]) == 0 {
            bf.Fd5Result[5] = bond_notes[2]
          }
          switch bf.Fd5Compound[0] {
          case 'c', 'C':
            bf.Fd5Result[6] = fmt.Sprintf("Convexity: %.3f", b.ConvexityContinuous(cf, current,
              b.CurrentPriceContinuous(cf, current)))
          default:
            bf.Fd5Result[6] = fmt.Sprintf("Convexity: %.3f", b.Convexity(cf, current,
              b.GetCompoundingPeriod(bf.Fd5Compound[0], true)))
          }
        }
        logger.LogInfo(fmt.Sprintf(
         "fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s",
         bf.Fd5FaceValue, bf.Fd5Time, bf.Fd5TimePeriod, bf.Fd5Coupon, bf.Fd5Compound,
         bf.Fd5CurInterest, bf.Fd5Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
        "webfinances/templates/header.html",
        "webfinances/templates/bonds/duration.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
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
      } { "Bonds", logger.DatetimeFormat(), bf.CurrentButton, newSession.CsrfToken, bf.Fd5FaceValue,
          bf.Fd5Time, bf.Fd5TimePeriod, bf.Fd5Coupon, bf.Fd5CompoundCoupon, bf.Fd5CurInterest,
          bf.Fd5Compound, bf.Fd5Result,
        })
    // } else if strings.EqualFold(bf.CurrentPage, "rhs-ui6") {
    //   bf.CurrentButton = "lhs-button6"
    //   if req.Method == http.MethodPost {
    //     bf.Fd6FaceValue = req.PostFormValue("fd6-facevalue")
    //     bf.Fd6Time = req.PostFormValue("fd6-time")
    //     bf.Fd6TimePeriod = req.PostFormValue("fd6-tp")
    //     bf.Fd6Coupon = req.PostFormValue("fd6-coupon")
    //     bf.Fd6CompoundCoupon = req.PostFormValue("fd6-compound-coupon")
    //     bf.Fd6CurInterest = req.PostFormValue("fd6-current")
    //     bf.Fd6Compound = req.PostFormValue("fd6-compound")
    //     var fv float64
    //     var time float64
    //     var couponRate float64
    //     var current float64
    //     var err error
    //     if fv, err = strconv.ParseFloat(bf.Fd6FaceValue, 64); err != nil {
    //       bf.Fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.Fd6FaceValue, err)
    //     } else if time, err = strconv.ParseFloat(bf.Fd6Time, 64); err != nil {
    //       bf.Fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.Fd6Time, err)
    //     } else if couponRate, err = strconv.ParseFloat(bf.Fd6Coupon, 64); err != nil {
    //       bf.Fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.Fd6Coupon, err)
    //     } else if current, err = strconv.ParseFloat(bf.Fd6CurInterest, 64); err != nil {
    //       bf.Fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.Fd6CurInterest, err)
    //     } else {
    //       var b finances.Bonds
    //       cf := b.CashFlow(fv, couponRate, b.GetCompoundingPeriod(bf.Fd6CompoundCoupon[0], true),
    //         time, b.GetTimePeriod(bf.Fd6TimePeriod[0], true))
    //       switch bf.Fd6Compound[0] {
    //       case 'c', 'C':
    //         bf.Fd6Result[1] = fmt.Sprintf("Macaulay Duration: %.3f year(s)",
    //           b.MacaulayDurationContinuous(cf, b.CurrentPriceContinuous(cf, current)))
    //       default:
    //         bf.Fd6Result[1] = fmt.Sprintf("Macaulay Duration: %.3f year(s)",
    //           b.MacaulayDuration(cf, b.GetCompoundingPeriod(bf.Fd6CompoundCoupon[0], true),
    //             b.CurrentPrice(cf, current, b.GetCompoundingPeriod(bf.Fd6Compound[0], true))))
    //       }
    //     }
    //     logEntry.Print(INFO, correlationId, []string {
    //       fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s",
    //         bf.Fd6FaceValue, bf.Fd6Time, bf.Fd6TimePeriod, bf.Fd6Coupon, bf.Fd6Compound,
    //         bf.Fd6CurInterest, bf.Fd6Result[1]),
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
    //   } { "Bonds", m.DTF(), bf.CurrentButton, newSession.CsrfToken, bf.Fd6FaceValue, bf.Fd6Time,
    //       bf.Fd6TimePeriod, bf.Fd6Coupon, bf.Fd6CompoundCoupon, bf.Fd6CurInterest, bf.Fd6Compound,
    //       bf.Fd6Result,
    //     })
    // } else if strings.EqualFold(bf.CurrentPage, "rhs-ui7") {
    //   bf.CurrentButton = "lhs-button7"
    //   if req.Method == http.MethodPost {
    //     bf.Fd7FaceValue = req.PostFormValue("fd7-facevalue")
    //     bf.Fd7Time = req.PostFormValue("fd7-time")
    //     bf.Fd7TimePeriod = req.PostFormValue("fd7-tp")
    //     bf.Fd7Coupon = req.PostFormValue("fd7-coupon")
    //     bf.Fd7CompoundCoupon = req.PostFormValue("fd7-compound-coupon")
    //     bf.Fd7CurInterest = req.PostFormValue("fd7-current")
    //     bf.Fd7Compound = req.PostFormValue("fd7-compound")
    //     var fv float64
    //     var time float64
    //     var couponRate float64
    //     var current float64
    //     var err error
    //     if fv, err = strconv.ParseFloat(bf.Fd7FaceValue, 64); err != nil {
    //       bf.Fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.Fd7FaceValue, err)
    //     } else if time, err = strconv.ParseFloat(bf.Fd7Time, 64); err != nil {
    //       bf.Fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.Fd7Time, err)
    //     } else if couponRate, err = strconv.ParseFloat(bf.Fd7Coupon, 64); err != nil {
    //       bf.Fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.Fd7Coupon, err)
    //     } else if current, err = strconv.ParseFloat(bf.Fd7CurInterest, 64); err != nil {
    //       bf.Fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.Fd7CurInterest, err)
    //     } else {
    //       var b finances.Bonds
    //       cf := b.CashFlow(fv, couponRate, b.GetCompoundingPeriod(bf.Fd7CompoundCoupon[0], true),
    //         time, b.GetTimePeriod(bf.Fd7TimePeriod[0], true))
    //       bf.Fd7Result[1] = fmt.Sprintf("Modified Duration: %.3f%%",
    //         b.ModifiedDuration(cf, b.GetCompoundingPeriod(bf.Fd7CompoundCoupon[0], true),
    //         b.CurrentPrice(cf, current, b.GetCompoundingPeriod(bf.Fd7Compound[0], true))))
    //     }
    //     logEntry.Print(INFO, correlationId, []string {
    //       fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s",
    //         bf.Fd7FaceValue, bf.Fd7Time, bf.Fd7TimePeriod, bf.Fd7Coupon, bf.Fd7Compound,
    //         bf.Fd7CurInterest, bf.Fd7Result[1]),
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
    //   } { "Bonds", m.DTF(), bf.CurrentButton, newSession.CsrfToken, bf.Fd7FaceValue, bf.Fd7Time,
    //       bf.Fd7TimePeriod, bf.Fd7Coupon, bf.Fd7CompoundCoupon, bf.Fd7CurInterest, bf.Fd7Compound,
    //       bf.Fd7Result,
    //     })
    // } else if strings.EqualFold(bf.CurrentPage, "rhs-ui8") {
    //   bf.CurrentButton = "lhs-button8"
    //   if req.Method == http.MethodPost {
    //     bf.Fd8FaceValue = req.PostFormValue("fd8-facevalue")
    //     bf.Fd8Time = req.PostFormValue("fd8-time")
    //     bf.Fd8TimePeriod = req.PostFormValue("fd8-tp")
    //     bf.Fd8Coupon = req.PostFormValue("fd8-coupon")
    //     bf.Fd8CompoundCoupon = req.PostFormValue("fd8-compound-coupon")
    //     bf.Fd8CurInterest = req.PostFormValue("fd8-current")
    //     bf.Fd8Compound = req.PostFormValue("fd8-compound")
    //     var fv float64
    //     var time float64
    //     var couponRate float64
    //     var current float64
    //     var err error
    //     if fv, err = strconv.ParseFloat(bf.Fd8FaceValue, 64); err != nil {
    //       bf.Fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.Fd8FaceValue, err)
    //     } else if time, err = strconv.ParseFloat(bf.Fd8Time, 64); err != nil {
    //       bf.Fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.Fd8Time, err)
    //     } else if couponRate, err = strconv.ParseFloat(bf.Fd8Coupon, 64); err != nil {
    //       bf.Fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.Fd8Coupon, err)
    //     } else if current, err = strconv.ParseFloat(bf.Fd8CurInterest, 64); err != nil {
    //       bf.Fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.Fd8CurInterest, err)
    //     } else {
    //       var b finances.Bonds
    //       cf := b.CashFlow(fv, couponRate, b.GetCompoundingPeriod(bf.Fd8CompoundCoupon[0], true),
    //         time, b.GetTimePeriod(bf.Fd8TimePeriod[0], true))
    //       switch bf.Fd8Compound[0] {
    //       case 'c', 'C':
    //         bf.Fd8Result[1] = fmt.Sprintf("Convexity: %.3f", b.ConvexityContinuous(cf, current,
    //           b.CurrentPriceContinuous(cf, current)))
    //       default:
    //         bf.Fd8Result[1] = fmt.Sprintf("Convexity: %.3f", b.Convexity(cf, current,
    //           b.GetCompoundingPeriod(bf.Fd8Compound[0], true)))
    //       }
    //     }
    //     logEntry.Print(INFO, correlationId, []string {
    //       fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s",
    //         bf.Fd8FaceValue, bf.Fd8Time, bf.Fd8TimePeriod, bf.Fd8Coupon, bf.Fd8Compound,
    //         bf.Fd8CurInterest, bf.Fd8Result[1]),
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
    //   } { "Bonds", m.DTF(), bf.CurrentButton, newSession.CsrfToken, bf.Fd8FaceValue, bf.Fd8Time,
    //       bf.Fd8TimePeriod, bf.Fd8Coupon, bf.Fd8CompoundCoupon, bf.Fd8CurInterest, bf.Fd8Compound,
    //       bf.Fd8Result,
    //     })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", bf.CurrentPage)
      logger.LogError(errString, "-1")
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      logger.LogWarning("*** Request timeout ***", "-1")
      if strings.EqualFold(bf.CurrentPage, "rhs-ui1") {
        bf.Fd1Result = ""
      } else if strings.EqualFold(bf.CurrentPage, "rhs-ui2") {
        bf.Fd2Result = ""
      } else if strings.EqualFold(bf.CurrentPage, "rhs-ui3") {
        bf.Fd3Result = ""
      } else if strings.EqualFold(bf.CurrentPage, "rhs-ui4") {
        bf.Fd4Result[0] = ""
        bf.Fd4Result[1] = ""
      } else if strings.EqualFold(bf.CurrentPage, "rhs-ui5") {
        bf.Fd5Result[0] = ""
        bf.Fd5Result[1] = ""
        bf.Fd5Result[2] = ""
        bf.Fd5Result[3] = ""
        bf.Fd5Result[4] = ""
        bf.Fd5Result[5] = ""
        bf.Fd5Result[6] = ""
      // } else if strings.EqualFold(bf.CurrentPage, "rhs-ui6") {
      //   bf.Fd6Result[1] = ""
      // } else if strings.EqualFold(bf.CurrentPage, "rhs-ui7") {
      //   bf.Fd7Result[1] = ""
      // } else if strings.EqualFold(bf.CurrentPage, "rhs-ui8") {
      //   bf.Fd8Result[1] = ""
      }
    }
    //
    if data, err := json.Marshal(bf); err != nil {
      logger.LogError(fmt.Sprintf("%+v", err), "-1")
    } else {
      filePath := fmt.Sprintf("%s/%s/bonds.txt", mainDir, userName)
      if _, err := osu.WriteAllExclusiveLock1(filePath, data, os.O_CREATE | os.O_RDWR |
        os.O_TRUNC, 0o600); err != nil {
        logger.LogError(fmt.Sprintf("%+v", err), "-1")
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    logger.LogError(errString, "-1")
    panic(errString)
  }
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Request took %vms\n", time.Since(startTime).Microseconds()),
    correlationId)
}
