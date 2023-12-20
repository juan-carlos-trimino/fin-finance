package webfinances

import (
  "context"
  "finance/middlewares"
  "finance/finances"
	"finance/sessions"
  "fmt"
  "html/template"
  "net/http"
  "strconv"
  "strings"
)

var bond_notes = [...]string {
  "The modified duration of a bond is a measure of the sensitivity of the bond's price to " +
  "changes in interest rates. Since bond prices move in an inverse direction from interest " +
  "rates, for a one percent increase (decrease) in interest rates, the bond's price will " +
  "decrease (increase) by the percentage shown by the modified duration.",
  "The Macaulay duration is a measure of a bond's sensitivity to interest rate changes. The " +
  "duration is the weighed-average number of years the investor must hold a bond until the " +
  "present value of the bond's cash flows equals the amount paid for the bond.",
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
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering BondsPages/webfinances.",
  })
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    bf := getBondsFields(sessions.GetUserName(sessionToken))
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
      bf.currentPage = ui
    }
    //
    if strings.EqualFold(bf.currentPage, "rhs-ui1") {
      bf.currentButton = "lhs-button1"
      if req.Method == http.MethodPost {
        bf.fd1TaxFree = req.PostFormValue("fd1-taxfree")
        bf.fd1CityTax = req.PostFormValue("fd1-citytax")
        bf.fd1StateTax = req.PostFormValue("fd1-statetax")
        bf.fd1FederalTax = req.PostFormValue("fd1-federaltax")
        var taxFree float64
        var cityTax float64
        var stateTax float64
        var federalTax float64
        var err error
        if taxFree, err = strconv.ParseFloat(bf.fd1TaxFree, 64); err != nil {
          bf.fd1Result = fmt.Sprintf("Error: %s -- %+v", bf.fd1TaxFree, err)
        } else if cityTax, err = strconv.ParseFloat(bf.fd1CityTax, 64); err != nil {
          bf.fd1Result = fmt.Sprintf("Error: %s -- %+v", bf.fd1CityTax, err)
        } else if stateTax, err = strconv.ParseFloat(bf.fd1StateTax, 64); err != nil {
          bf.fd1Result = fmt.Sprintf("Error: %s -- %+v", bf.fd1StateTax, err)
        } else if federalTax, err = strconv.ParseFloat(bf.fd1FederalTax, 64); err != nil {
          bf.fd1Result = fmt.Sprintf("Error: %s -- %+v", bf.fd1FederalTax, err)
        } else {
          var b finances.Bonds
          bf.fd1Result = fmt.Sprintf("Taxable-Equivalent Yield: %.3f%%",
            b.TaxableVsTaxFreeYields(taxFree, cityTax, stateTax, federalTax) * 100.0)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("tax free = %s, city tax = %s, state tax = %s, federal tax = %s, %s",
            bf.fd1TaxFree, bf.fd1CityTax, bf.fd1StateTax, bf.fd1FederalTax, bf.fd1Result),
        })
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
      } { "Bonds", m.DTF(), bf.currentButton, newSession.CsrfToken, bf.fd1TaxFree, bf.fd1CityTax,
          bf.fd1StateTax, bf.fd1FederalTax, bf.fd1Result,
        })
    } else if strings.EqualFold(bf.currentPage, "rhs-ui2") {
      bf.currentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        bf.fd2FaceValue = req.FormValue("fd2-facevalue")
        bf.fd2Time = req.PostFormValue("fd2-time")
        bf.fd2TimePeriod = req.PostFormValue("fd2-tp")
        bf.fd2Coupon = req.PostFormValue("fd2-coupon")
        bf.fd2Current = req.PostFormValue("fd2-current")
        bf.fd2Compound = req.PostFormValue("fd2-compound")
        var fv float64
        var time float64
        var coupon float64
        var current float64
        var err error
        if fv, err = strconv.ParseFloat(bf.fd2FaceValue, 64); err != nil {
          bf.fd2Result = fmt.Sprintf("Error: %s -- %+v", bf.fd2FaceValue, err)
        } else if time, err = strconv.ParseFloat(bf.fd2Time, 64); err != nil {
          bf.fd2Result = fmt.Sprintf("Error: %s -- %+v", bf.fd2Time, err)
        } else if coupon, err = strconv.ParseFloat(bf.fd2Coupon, 64); err != nil {
          bf.fd2Result = fmt.Sprintf("Error: %s -- %+v", bf.fd2Coupon, err)
        } else if current, err = strconv.ParseFloat(bf.fd2Current, 64); err != nil {
          bf.fd2Result = fmt.Sprintf("Error: %s -- %+v", bf.fd2Current, err)
        } else {
          var b finances.Bonds
          cf := b.CashFlow(fv, coupon, b.GetCompoundingPeriod(bf.fd2Compound[0], true), time,
            b.GetTimePeriod(bf.fd2TimePeriod[0], true))
          currentPrice := b.CurrentPrice(cf, current, b.GetCompoundingPeriod(bf.fd2Compound[0],
            true))
          if fv > currentPrice {
            bf.fd2Result = fmt.Sprintf("Current Price: $%.2f (discount)", currentPrice)
          } else if fv < currentPrice {
            bf.fd2Result = fmt.Sprintf("Current Price: $%.2f (premium)", currentPrice)
          } else {
            bf.fd2Result = fmt.Sprintf("Current Price: $%.2f (par)", currentPrice)
          }
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon rate = %s, current interest = %s, cp = %s, %s",
            bf.fd2FaceValue, bf.fd2Time, bf.fd2TimePeriod, bf.fd2Coupon, bf.fd2Current,
            bf.fd2Compound, bf.fd2Result),
        })
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
        Fd2Current string
        Fd2Compound string
        Fd2Result string
      } { "Bonds", m.DTF(), bf.currentButton, newSession.CsrfToken, bf.fd2FaceValue, bf.fd2Time,
          bf.fd2TimePeriod, bf.fd2Coupon, bf.fd2Current, bf.fd2Compound, bf.fd2Result,
        })
    } else if strings.EqualFold(bf.currentPage, "rhs-ui3") {
      bf.currentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        bf.fd3FaceValue = req.PostFormValue("fd3-facevalue")
        bf.fd3TimeCall = req.PostFormValue("fd3-timecall")
        bf.fd3TimePeriod = req.PostFormValue("fd3-tp")
        bf.fd3Coupon = req.PostFormValue("fd3-coupon")
        bf.fd3BondPrice = req.PostFormValue("fd3-bondprice")
        bf.fd3CallPrice = req.PostFormValue("fd3-callprice")
        bf.fd3Compound = req.PostFormValue("fd3-compound")
        var fv float64
        var timeToCall float64
        var couponRate float64
        var bondPrice float64
        var callPrice float64
        var err error
        if fv, err = strconv.ParseFloat(bf.fd3FaceValue, 64); err != nil {
          bf.fd3Result = fmt.Sprintf("Error: %s -- %+v", bf.fd3FaceValue, err)
        } else if timeToCall, err = strconv.ParseFloat(bf.fd3TimeCall, 64); err != nil {
          bf.fd3Result = fmt.Sprintf("Error: %s -- %+v", bf.fd3TimeCall, err)
        } else if couponRate, err = strconv.ParseFloat(bf.fd3Coupon, 64); err != nil {
          bf.fd3Result = fmt.Sprintf("Error: %s -- %+v", bf.fd3Coupon, err)
        } else if bondPrice, err = strconv.ParseFloat(bf.fd3BondPrice, 64); err != nil {
          bf.fd3Result = fmt.Sprintf("Error: %s -- %+v", bf.fd3BondPrice, err)
        } else if callPrice, err = strconv.ParseFloat(bf.fd3CallPrice, 64); err != nil {
          bf.fd3Result = fmt.Sprintf("Error: %s -- %+v", bf.fd3CallPrice, err)
        } else {
          var b finances.Bonds
          bf.fd3Result = fmt.Sprintf("Yield to Call: %.3f%%", b.YieldToCall(fv, couponRate,
            b.GetCompoundingPeriod(bf.fd3Compound[0], true), timeToCall,
            b.GetTimePeriod(bf.fd3TimePeriod[0], true), bondPrice, callPrice))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("fv = %s, coupon rate = %s, cp = %s, time to call = %s, tp = %s, bond price = %s, call price = %s, %s\n",
            bf.fd3FaceValue, bf.fd3Coupon, bf.fd3Compound, bf.fd3TimeCall, bf.fd3TimePeriod,
            bf.fd3BondPrice, bf.fd3CallPrice, bf.fd3Result),
        })
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
      } { "Bonds", m.DTF(), bf.currentButton, newSession.CsrfToken, bf.fd3FaceValue,
          bf.fd3TimeCall, bf.fd3TimePeriod, bf.fd3Coupon, bf.fd3Compound, bf.fd3BondPrice,
          bf.fd3CallPrice, bf.fd3Result,
        })
    } else if strings.EqualFold(bf.currentPage, "rhs-ui4") {
      bf.currentButton = "lhs-button4"
      if req.Method == http.MethodPost {
        bf.fd4FaceValue = req.PostFormValue("fd4-facevalue")
        bf.fd4Time = req.PostFormValue("fd4-time")
        bf.fd4TimePeriod = req.PostFormValue("fd4-tp")
        bf.fd4Coupon = req.PostFormValue("fd4-coupon")
        bf.fd4Compound = req.PostFormValue("fd4-compound")
        bf.fd4CurrentRadio = req.PostFormValue("fd4-choice")
        bf.fd4CurInterest = req.PostFormValue("fd4-ci")
        bf.fd4BondPrice = req.PostFormValue("fd4-bp")
        var currentInterest bool = false
        if strings.EqualFold(bf.fd4CurrentRadio, "fd4-curinterest") {
          currentInterest = true
        }
        var fv float64
        var time float64
        var couponRate float64
        var curInterest float64
        var bondPrice float64
        var err error
        if fv, err = strconv.ParseFloat(bf.fd4FaceValue, 64); err != nil {
          bf.fd4Result = fmt.Sprintf("Error: %s -- %+v", bf.fd4FaceValue, err)
        } else if time, err = strconv.ParseFloat(bf.fd4Time, 64); err != nil {
          bf.fd4Result = fmt.Sprintf("Error: %s -- %+v", bf.fd4Time, err)
        } else if couponRate, err = strconv.ParseFloat(bf.fd4Coupon, 64); err != nil {
          bf.fd4Result = fmt.Sprintf("Error: %s -- %+v", bf.fd4Coupon, err)
        } else if curInterest, err = strconv.ParseFloat(bf.fd4CurInterest, 64); err != nil {
          bf.fd4Result = fmt.Sprintf("Error: %s -- %+v", bf.fd4CurInterest, err)
        } else if bondPrice, err = strconv.ParseFloat(bf.fd4BondPrice, 64); err != nil {
          bf.fd4Result = fmt.Sprintf("Error: %s -- %+v", bf.fd4BondPrice, err)
        } else {
          var b finances.Bonds
          var cp int = b.GetCompoundingPeriod(bf.fd4Compound[0], false)
          var tp = b.GetTimePeriod(bf.fd4TimePeriod[0], false)
          cf := b.CashFlow(fv, couponRate, cp, time, tp)
          if currentInterest {
            if cp != finances.Continuously {
              bf.fd4Result = fmt.Sprintf("Yield to Maturity: %.3f%%",
                b.YieldToMaturity(cf, b.CurrentPrice(cf, curInterest, cp), tp))
            } else {
              bf.fd4Result = fmt.Sprintf("Yield to Maturity: %.3f%%",
                b.YieldToMaturityContinuous(cf, b.CurrentPriceContinuous(cf, curInterest)))
            }
          } else {  //Bond price.
            if cp != finances.Continuously {
              bf.fd4Result = fmt.Sprintf("Yield to Maturity: %.3f%%",
                b.YieldToMaturity(cf, bondPrice, cp))
            } else {
              bf.fd4Result = fmt.Sprintf("Yield to Maturity: %.3f%%",
                b.YieldToMaturityContinuous(cf, bondPrice))
            }
          }
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur radio = %s, cur interest = %s, bond price = %s, %s",
            bf.fd4FaceValue, bf.fd4Time, bf.fd4TimePeriod, bf.fd4Coupon, bf.fd4Compound,
            bf.fd4CurrentRadio, bf.fd4CurInterest, bf.fd4BondPrice, bf.fd4Result),
        })
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
        Fd4Result string
      } { "Bonds", m.DTF(), bf.currentButton, newSession.CsrfToken, bf.fd4FaceValue, bf.fd4Time,
          bf.fd4TimePeriod, bf.fd4Coupon, bf.fd4Compound, bf.fd4CurrentRadio, bf.fd4CurInterest,
          bf.fd4BondPrice, bf.fd4Result,
        })
    } else if strings.EqualFold(bf.currentPage, "rhs-ui5") {
      bf.currentButton = "lhs-button5"
      if req.Method == http.MethodPost {
        bf.fd5FaceValue = req.PostFormValue("fd5-facevalue")
        bf.fd5Time = req.PostFormValue("fd5-time")
        bf.fd5TimePeriod = req.PostFormValue("fd5-tp")
        bf.fd5Coupon = req.PostFormValue("fd5-coupon")
        bf.fd5CurInterest = req.PostFormValue("fd5-current")
        bf.fd5Compound = req.PostFormValue("fd5-compound")
        var fv float64
        var time float64
        var couponRate float64
        var curInterest float64
        var err error
        if fv, err = strconv.ParseFloat(bf.fd5FaceValue, 64); err != nil {
          bf.fd5Result = fmt.Sprintf("Error: %s -- %+v", bf.fd5FaceValue, err)
        } else if time, err = strconv.ParseFloat(bf.fd5Time, 64); err != nil {
          bf.fd5Result = fmt.Sprintf("Error: %s -- %+v", bf.fd5Time, err)
        } else if couponRate, err = strconv.ParseFloat(bf.fd5Coupon, 64); err != nil {
          bf.fd5Result = fmt.Sprintf("Error: %s -- %+v", bf.fd5Coupon, err)
        } else if curInterest, err = strconv.ParseFloat(bf.fd5CurInterest, 64); err != nil {
          bf.fd5Result = fmt.Sprintf("Error: %s -- %+v", bf.fd5CurInterest, err)
        } else {
          var b finances.Bonds
          var cp int = b.GetCompoundingPeriod(bf.fd5Compound[0], false)
          var tp = b.GetTimePeriod(bf.fd5TimePeriod[0], false)
          cf := b.CashFlow(fv, couponRate, cp, time, tp)
          if cp != finances.Continuously {
            bf.fd5Result = fmt.Sprintf("Duration: %.3f%%",
              b.Duration(cf, cp, curInterest, b.CurrentPrice(cf, curInterest, cp)))
          } else {
            bf.fd5Result = "-1.00"
          }
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s",
            bf.fd5FaceValue, bf.fd5Time, bf.fd5TimePeriod, bf.fd5Coupon, bf.fd5Compound,
            bf.fd5CurInterest, bf.fd5Result),
        })
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
        Fd5CurInterest string
        Fd5Compound string
        Fd5Result string
      } { "Bonds", m.DTF(), bf.currentButton, newSession.CsrfToken, bf.fd5FaceValue, bf.fd5Time,
          bf.fd5TimePeriod, bf.fd5Coupon, bf.fd5CurInterest, bf.fd5Compound, bf.fd5Result,
        })
    } else if strings.EqualFold(bf.currentPage, "rhs-ui6") {
      bf.currentButton = "lhs-button6"
      if req.Method == http.MethodPost {
        bf.fd6FaceValue = req.PostFormValue("fd6-facevalue")
        bf.fd6Time = req.PostFormValue("fd6-time")
        bf.fd6TimePeriod = req.PostFormValue("fd6-tp")
        bf.fd6Coupon = req.PostFormValue("fd6-coupon")
        bf.fd6CurInterest = req.PostFormValue("fd6-current")
        bf.fd6Compound = req.PostFormValue("fd6-compound")
        var fv float64
        var time float64
        var couponRate float64
        var curInterest float64
        var err error
        if fv, err = strconv.ParseFloat(bf.fd6FaceValue, 64); err != nil {
          bf.fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.fd6FaceValue, err)
        } else if time, err = strconv.ParseFloat(bf.fd6Time, 64); err != nil {
          bf.fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.fd6Time, err)
        } else if couponRate, err = strconv.ParseFloat(bf.fd6Coupon, 64); err != nil {
          bf.fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.fd6Coupon, err)
        } else if curInterest, err = strconv.ParseFloat(bf.fd6CurInterest, 64); err != nil {
          bf.fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.fd6CurInterest, err)
        } else {
          var b finances.Bonds
          var cp int = b.GetCompoundingPeriod(bf.fd6Compound[0], false)
          var tp = b.GetTimePeriod(bf.fd6TimePeriod[0], false)
          cf := b.CashFlow(fv, couponRate, cp, time, tp)
          if cp != finances.Continuously {
            bf.fd6Result[1] = fmt.Sprintf("Macaulay Duration: %.3f year(s)",
              b.MacaulayDuration(cf, cp, b.CurrentPrice(cf, curInterest, cp)))
          } else {
            bf.fd6Result[1] = "-1.00"
          }
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s",
            bf.fd6FaceValue, bf.fd6Time, bf.fd6TimePeriod, bf.fd6Coupon, bf.fd6Compound,
            bf.fd6CurInterest, bf.fd6Result[1]),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
        "webfinances/templates/header.html",
        "webfinances/templates/bonds/macaulayduration.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd6FaceValue string
        Fd6Time string
        Fd6TimePeriod string
        Fd6Coupon string
        Fd6CurInterest string
        Fd6Compound string
        Fd6Result [2]string
      } { "Bonds", m.DTF(), bf.currentButton, newSession.CsrfToken, bf.fd6FaceValue, bf.fd6Time,
          bf.fd6TimePeriod, bf.fd6Coupon, bf.fd6CurInterest, bf.fd6Compound, bf.fd6Result,
        })
    } else if strings.EqualFold(bf.currentPage, "rhs-ui7") {
      bf.currentButton = "lhs-button7"
      if req.Method == http.MethodPost {
        bf.fd7FaceValue = req.PostFormValue("fd7-facevalue")
        bf.fd7Time = req.PostFormValue("fd7-time")
        bf.fd7TimePeriod = req.PostFormValue("fd7-tp")
        bf.fd7Coupon = req.PostFormValue("fd7-coupon")
        bf.fd7CurInterest = req.PostFormValue("fd7-current")
        bf.fd7Compound = req.PostFormValue("fd7-compound")
        var fv float64
        var time float64
        var couponRate float64
        var curInterest float64
        var err error
        if fv, err = strconv.ParseFloat(bf.fd7FaceValue, 64); err != nil {
          bf.fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.fd7FaceValue, err)
        } else if time, err = strconv.ParseFloat(bf.fd7Time, 64); err != nil {
          bf.fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.fd7Time, err)
        } else if couponRate, err = strconv.ParseFloat(bf.fd7Coupon, 64); err != nil {
          bf.fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.fd7Coupon, err)
        } else if curInterest, err = strconv.ParseFloat(bf.fd7CurInterest, 64); err != nil {
          bf.fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.fd7CurInterest, err)
        } else {
          var b finances.Bonds
          var cp int = b.GetCompoundingPeriod(bf.fd7Compound[0], false)
          var tp = b.GetTimePeriod(bf.fd7TimePeriod[0], false)
          cf := b.CashFlow(fv, couponRate, cp, time, tp)
          if cp != finances.Continuously {
            bf.fd7Result[1] = fmt.Sprintf("Modified Duration: %.3f%%",
              b.ModifiedDuration(cf, cp, b.CurrentPrice(cf, curInterest, cp)))
          } else {
            bf.fd7Result[1] = "-1.00"
          }
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s",
            bf.fd7FaceValue, bf.fd7Time, bf.fd7TimePeriod, bf.fd7Coupon, bf.fd7Compound,
            bf.fd7CurInterest, bf.fd7Result[1]),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
        "webfinances/templates/header.html",
        "webfinances/templates/bonds/modifiedduration.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd7FaceValue string
        Fd7Time string
        Fd7TimePeriod string
        Fd7Coupon string
        Fd7CurInterest string
        Fd7Compound string
        Fd7Result [2]string
      } { "Bonds", m.DTF(), bf.currentButton, newSession.CsrfToken, bf.fd7FaceValue, bf.fd7Time,
          bf.fd7TimePeriod, bf.fd7Coupon, bf.fd7CurInterest, bf.fd7Compound, bf.fd7Result,
        })
    } else if strings.EqualFold(bf.currentPage, "rhs-ui8") {
      bf.currentButton = "lhs-button8"
      if req.Method == http.MethodPost {
        bf.fd8FaceValue = req.PostFormValue("fd8-facevalue")
        bf.fd8Time = req.PostFormValue("fd8-time")
        bf.fd8TimePeriod = req.PostFormValue("fd8-tp")
        bf.fd8Coupon = req.PostFormValue("fd8-coupon")
        bf.fd8CurInterest = req.PostFormValue("fd8-current")
        bf.fd8Compound = req.PostFormValue("fd8-compound")
        var fv float64
        var time float64
        var couponRate float64
        var curInterest float64
        var err error
        if fv, err = strconv.ParseFloat(bf.fd8FaceValue, 64); err != nil {
          bf.fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.fd8FaceValue, err)
        } else if time, err = strconv.ParseFloat(bf.fd8Time, 64); err != nil {
          bf.fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.fd8Time, err)
        } else if couponRate, err = strconv.ParseFloat(bf.fd8Coupon, 64); err != nil {
          bf.fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.fd8Coupon, err)
        } else if curInterest, err = strconv.ParseFloat(bf.fd8CurInterest, 64); err != nil {
          bf.fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", bf.fd8CurInterest, err)
        } else {
          var b finances.Bonds
          var cp int = b.GetCompoundingPeriod(bf.fd8Compound[0], false)
          var tp = b.GetTimePeriod(bf.fd8TimePeriod[0], false)
          cf := b.CashFlow(fv, couponRate, cp, time, tp)
          if cp != finances.Continuously {
            bf.fd8Result[1] = fmt.Sprintf("Convexity: %.3f", b.Convexity(cf, curInterest, cp))
          } else {
            bf.fd8Result[1] = "-1.00"
          }
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s",
            bf.fd8FaceValue, bf.fd8Time, bf.fd8TimePeriod, bf.fd8Coupon, bf.fd8Compound,
            bf.fd8CurInterest, bf.fd8Result[1]),
        })
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
        "webfinances/templates/header.html",
        "webfinances/templates/bonds/convexity.html",
        "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Fd8FaceValue string
        Fd8Time string
        Fd8TimePeriod string
        Fd8Coupon string
        Fd8CurInterest string
        Fd8Compound string
        Fd8Result [2]string
      } { "Bonds", m.DTF(), bf.currentButton, newSession.CsrfToken, bf.fd8FaceValue, bf.fd8Time,
          bf.fd8TimePeriod, bf.fd8Coupon, bf.fd8CurInterest, bf.fd8Compound, bf.fd8Result,
        })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", bf.currentPage)
      fmt.Printf("%s - %s\n", m.DTF(), errString)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      fmt.Println("*** Request timeout ***")
      if strings.EqualFold(bf.currentPage, "rhs-ui1") {
        bf.fd1Result = ""
      } else if strings.EqualFold(bf.currentPage, "rhs-ui2") {
        bf.fd2Result = ""
      } else if strings.EqualFold(bf.currentPage, "rhs-ui3") {
        bf.fd3Result = ""
      } else if strings.EqualFold(bf.currentPage, "rhs-ui4") {
        bf.fd4Result = ""
      } else if strings.EqualFold(bf.currentPage, "rhs-ui5") {
        bf.fd5Result = ""
      } else if strings.EqualFold(bf.currentPage, "rhs-ui6") {
        bf.fd6Result[1] = ""
      } else if strings.EqualFold(bf.currentPage, "rhs-ui7") {
        bf.fd7Result[1] = ""
      } else if strings.EqualFold(bf.currentPage, "rhs-ui8") {
        bf.fd8Result[1] = ""
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
