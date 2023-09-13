package webfinances

import (
  "finance/middlewares"
  "finance/finances"
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

type WfBondsPages interface {
  BondsPages(http.ResponseWriter, *http.Request)
}

type wfBondsPages struct {
  currentPage string
  currentButton string
  //
  fd1TaxFree string
  fd1CityTax string
  fd1StateTax string
  fd1FederalTax string
  fd1Result string
  //
  fd2FaceValue string
  fd2Time string
  fd2TimePeriod string
  fd2Coupon string
  fd2Current string
  fd2Compound string
  fd2Result string
  //
  fd3FaceValue string
  fd3TimeCall string
  fd3TimePeriod string
  fd3Coupon string
  fd3Compound string
  fd3BondPrice string
  fd3CallPrice string
  fd3Result string
  //
  fd4FaceValue string
  fd4Time string
  fd4TimePeriod string
  fd4Coupon string
  fd4Compound string
  fd4CurrentRadio string
  fd4CurInterest string
  fd4BondPrice string
  fd4Result string
  //
  fd5FaceValue string
  fd5Time string
  fd5TimePeriod string
  fd5Coupon string
  fd5CurInterest string
  fd5Compound string
  fd5Result string
  //
  fd6FaceValue string
  fd6Time string
  fd6TimePeriod string
  fd6Coupon string
  fd6CurInterest string
  fd6Compound string
  fd6Result [2]string
  //
  fd7FaceValue string
  fd7Time string
  fd7TimePeriod string
  fd7Coupon string
  fd7CurInterest string
  fd7Compound string
  fd7Result [2]string
  //
  fd8FaceValue string
  fd8Time string
  fd8TimePeriod string
  fd8Coupon string
  fd8CurInterest string
  fd8Compound string
  fd8Result [2]string
}

func NewWfBondsPages() WfBondsPages {
  return &wfBondsPages {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1TaxFree: "3.5",
    fd1CityTax: "0.0",
    fd1StateTax: "1.0",
    fd1FederalTax: "23.0",
    fd1Result: "",
    //
    fd2FaceValue: "1000.00",
    fd2Time: "5",
    fd2TimePeriod: "year",
    fd2Coupon: "3.00",
    fd2Current: "3.5",
    fd2Compound: "semiannually",
    fd2Result: "",
    //
    fd3FaceValue: "1000.00",
    fd3TimeCall: "2",
    fd3TimePeriod: "year",
    fd3Coupon: "2.0",
    fd3Compound: "semiannually",
    fd3BondPrice: "990.00",
    fd3CallPrice: "1050.00",
    fd3Result: "",
    //
    fd4FaceValue: "1000.00",
    fd4Time: "3",
    fd4TimePeriod: "year",
    fd4Coupon: "2.5",
    fd4Compound: "semiannually",
    fd4CurrentRadio: "fd4-curinterest",
    fd4CurInterest: "2.3",
    fd4BondPrice: "1000.00",
    fd4Result: "",
    //
    fd5FaceValue: "1000.00",
    fd5Time: "5",
    fd5TimePeriod: "year",
    fd5Coupon: "5.4",
    fd5CurInterest: "7.5",
    fd5Compound: "semiannually",
    fd5Result: "",
    //
    fd6FaceValue: "1000.00",
    fd6Time: "5",
    fd6TimePeriod: "year",
    fd6Coupon: "5.4",
    fd6CurInterest: "7.5",
    fd6Compound: "semiannually",
    fd6Result: [2]string { bond_notes[1], "" },
    //
    fd7FaceValue: "1000.00",
    fd7Time: "5",
    fd7TimePeriod: "year",
    fd7Coupon: "5.4",
    fd7CurInterest: "7.5",
    fd7Compound: "semiannually",
    fd7Result: [2]string { bond_notes[0], "" },
    //
    fd8FaceValue: "1000.00",
    fd8Time: "5",
    fd8TimePeriod: "year",
    fd8Coupon: "5.4",
    fd8CurInterest: "7.5",
    fd8Compound: "semiannually",
    fd8Result: [2]string { bond_notes[2], "" },
  }
}

func (p *wfBondsPages) BondsPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering BondsPages/webfinances.",
  })
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
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
      p.currentPage = ui
    }
    //
    if strings.EqualFold(p.currentPage, "rhs-ui1") {
      p.currentButton = "lhs-button1"
      if req.Method == http.MethodPost {
        p.fd1TaxFree = req.PostFormValue("fd1-taxfree")
        p.fd1CityTax = req.PostFormValue("fd1-citytax")
        p.fd1StateTax = req.PostFormValue("fd1-statetax")
        p.fd1FederalTax = req.PostFormValue("fd1-federaltax")
        var taxFree float64
        var cityTax float64
        var stateTax float64
        var federalTax float64
        var err error
        if taxFree, err = strconv.ParseFloat(p.fd1TaxFree, 64); err != nil {
          p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1TaxFree, err)
        } else if cityTax, err = strconv.ParseFloat(p.fd1CityTax, 64); err != nil {
          p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1CityTax, err)
        } else if stateTax, err = strconv.ParseFloat(p.fd1StateTax, 64); err != nil {
          p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1StateTax, err)
        } else if federalTax, err = strconv.ParseFloat(p.fd1FederalTax, 64); err != nil {
          p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1FederalTax, err)
        } else {
          var b finances.Bonds
          p.fd1Result = fmt.Sprintf("Taxable-Equivalent Yield: %.3f%%", b.TaxableVsTaxFreeYields(
                                     taxFree, cityTax, stateTax, federalTax) * 100.0)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("tax free = %s, city tax = %s, state tax = %s, federal tax = %s, %s",
                      p.fd1TaxFree, p.fd1CityTax, p.fd1StateTax, p.fd1FederalTax, p.fd1Result),
        })
      }
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
        Fd1TaxFree string
        Fd1CityTax string
        Fd1StateTax string
        Fd1FederalTax string
        Fd1Result string
      } { "Bonds", m.DTF(), p.currentButton,
          p.fd1TaxFree, p.fd1CityTax, p.fd1StateTax, p.fd1FederalTax, p.fd1Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui2") {
      p.currentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        p.fd2FaceValue = req.FormValue("fd2-facevalue")
        p.fd2Time = req.PostFormValue("fd2-time")
        p.fd2TimePeriod = req.PostFormValue("fd2-tp")
        p.fd2Coupon = req.PostFormValue("fd2-coupon")
        p.fd2Current = req.PostFormValue("fd2-current")
        p.fd2Compound = req.PostFormValue("fd2-compound")
        var fv float64
        var time float64
        var coupon float64
        var current float64
        var err error
        if fv, err = strconv.ParseFloat(p.fd2FaceValue, 64); err != nil {
          p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2FaceValue, err)
        } else if time, err = strconv.ParseFloat(p.fd2Time, 64); err != nil {
          p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2Time, err)
        } else if coupon, err = strconv.ParseFloat(p.fd2Coupon, 64); err != nil {
          p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2Coupon, err)
        } else if current, err = strconv.ParseFloat(p.fd2Current, 64); err != nil {
          p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2Current, err)
        } else {
          var b finances.Bonds
          cf := b.CashFlow(fv, coupon, b.GetCompoundingPeriod(p.fd2Compound[0], true), time,
                           b.GetTimePeriod(p.fd2TimePeriod[0], true))
          currentPrice := b.CurrentPrice(cf, current, b.GetCompoundingPeriod(p.fd2Compound[0], true))
          if fv > currentPrice {
            p.fd2Result = fmt.Sprintf("Current Price: $%.2f (discount)", currentPrice)
          } else if fv < currentPrice {
            p.fd2Result = fmt.Sprintf("Current Price: $%.2f (premium)", currentPrice)
          } else {
            p.fd2Result = fmt.Sprintf("Current Price: $%.2f (par)", currentPrice)
          }
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon rate = %s, current interest = %s, cp = %s, %s",
                      p.fd2FaceValue, p.fd2Time, p.fd2TimePeriod, p.fd2Coupon, p.fd2Current,
                      p.fd2Compound, p.fd2Result),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/bonds/currentprice.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd2FaceValue string
        Fd2Time string
        Fd2TimePeriod string
        Fd2Coupon string
        Fd2Current string
        Fd2Compound string
        Fd2Result string
      } { "Bonds", m.DTF(), p.currentButton,
          p.fd2FaceValue, p.fd2Time, p.fd2TimePeriod, p.fd2Coupon, p.fd2Current, p.fd2Compound, p.fd2Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui3") {
      p.currentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        p.fd3FaceValue = req.PostFormValue("fd3-facevalue")
        p.fd3TimeCall = req.PostFormValue("fd3-timecall")
        p.fd3TimePeriod = req.PostFormValue("fd3-tp")
        p.fd3Coupon = req.PostFormValue("fd3-coupon")
        p.fd3BondPrice = req.PostFormValue("fd3-bondprice")
        p.fd3CallPrice = req.PostFormValue("fd3-callprice")
        p.fd3Compound = req.PostFormValue("fd3-compound")
        var fv float64
        var timeToCall float64
        var couponRate float64
        var bondPrice float64
        var callPrice float64
        var err error
        if fv, err = strconv.ParseFloat(p.fd3FaceValue, 64); err != nil {
          p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3FaceValue, err)
        } else if timeToCall, err = strconv.ParseFloat(p.fd3TimeCall, 64); err != nil {
          p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3TimeCall, err)
        } else if couponRate, err = strconv.ParseFloat(p.fd3Coupon, 64); err != nil {
          p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3Coupon, err)
        } else if bondPrice, err = strconv.ParseFloat(p.fd3BondPrice, 64); err != nil {
          p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3BondPrice, err)
        } else if callPrice, err = strconv.ParseFloat(p.fd3CallPrice, 64); err != nil {
          p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3CallPrice, err)
        } else {
          var b finances.Bonds
          p.fd3Result = fmt.Sprintf("Yield to Call: %.3f%%", b.YieldToCall(fv, couponRate,
                                     b.GetCompoundingPeriod(p.fd3Compound[0], true), timeToCall,
                                     b.GetTimePeriod(p.fd3TimePeriod[0], true), bondPrice,
                                     callPrice))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("fv = %s, coupon rate = %s, cp = %s, time to call = %s, tp = %s, bond price = %s, call price = %s, %s\n",
                       p.fd3FaceValue, p.fd3Coupon, p.fd3Compound, p.fd3TimeCall, p.fd3TimePeriod,
                       p.fd3BondPrice, p.fd3CallPrice, p.fd3Result),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/bonds/yieldtocall.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd3FaceValue string
        Fd3TimeCall string
        Fd3TimePeriod string
        Fd3Coupon string
        Fd3Compound string
        Fd3BondPrice string
        Fd3CallPrice string
        Fd3Result string
      } { "Bonds", m.DTF(), p.currentButton,
          p.fd3FaceValue, p.fd3TimeCall, p.fd3TimePeriod, p.fd3Coupon, p.fd3Compound, p.fd3BondPrice, p.fd3CallPrice, p.fd3Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui4") {
      p.currentButton = "lhs-button4"
      if req.Method == http.MethodPost {
        p.fd4FaceValue = req.PostFormValue("fd4-facevalue")
        p.fd4Time = req.PostFormValue("fd4-time")
        p.fd4TimePeriod = req.PostFormValue("fd4-tp")
        p.fd4Coupon = req.PostFormValue("fd4-coupon")
        p.fd4Compound = req.PostFormValue("fd4-compound")
        p.fd4CurrentRadio = req.PostFormValue("fd4-choice")
        p.fd4CurInterest = req.PostFormValue("fd4-ci")
        p.fd4BondPrice = req.PostFormValue("fd4-bp")
        var currentInterest bool = false
        if strings.EqualFold(p.fd4CurrentRadio, "fd4-curinterest") {
          currentInterest = true
        }
        var fv float64
        var time float64
        var couponRate float64
        var curInterest float64
        var bondPrice float64
        var err error
        if fv, err = strconv.ParseFloat(p.fd4FaceValue, 64); err != nil {
          p.fd4Result = fmt.Sprintf("Error: %s -- %+v", p.fd4FaceValue, err)
        } else if time, err = strconv.ParseFloat(p.fd4Time, 64); err != nil {
          p.fd4Result = fmt.Sprintf("Error: %s -- %+v", p.fd4Time, err)
        } else if couponRate, err = strconv.ParseFloat(p.fd4Coupon, 64); err != nil {
          p.fd4Result = fmt.Sprintf("Error: %s -- %+v", p.fd4Coupon, err)
        } else if curInterest, err = strconv.ParseFloat(p.fd4CurInterest, 64); err != nil {
          p.fd4Result = fmt.Sprintf("Error: %s -- %+v", p.fd4CurInterest, err)
        } else if bondPrice, err = strconv.ParseFloat(p.fd4BondPrice, 64); err != nil {
          p.fd4Result = fmt.Sprintf("Error: %s -- %+v", p.fd4BondPrice, err)
        } else {
          var b finances.Bonds
          var cp int = b.GetCompoundingPeriod(p.fd4Compound[0], false)
          var tp = b.GetTimePeriod(p.fd4TimePeriod[0], false)
          cf := b.CashFlow(fv, couponRate, cp, time, tp)
          if currentInterest {
            if cp != finances.Continuously {
              p.fd4Result = fmt.Sprintf("Yield to Maturity: %.3f%%", b.YieldToMaturity(cf,
                                        b.CurrentPrice(cf, curInterest, cp), tp))
            } else {
              p.fd4Result = fmt.Sprintf("Yield to Maturity: %.3f%%", b.YieldToMaturityContinuous(cf,
                                        b.CurrentPriceContinuous(cf, curInterest)))
            }
          } else {  //Bond price.
            if cp != finances.Continuously {
              p.fd4Result = fmt.Sprintf("Yield to Maturity: %.3f%%", b.YieldToMaturity(cf, bondPrice,
                                        cp))
            } else {
              p.fd4Result = fmt.Sprintf("Yield to Maturity: %.3f%%", b.YieldToMaturityContinuous(cf,
                                        bondPrice))
            }
          }
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur radio = %s, cur interest = %s, bond price = %s, %s",
                      p.fd4FaceValue, p.fd4Time, p.fd4TimePeriod, p.fd4Coupon, p.fd4Compound,
                      p.fd4CurrentRadio, p.fd4CurInterest, p.fd4BondPrice, p.fd4Result),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/bonds/yieldtomaturity.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd4FaceValue string
        Fd4Time string
        Fd4TimePeriod string
        Fd4Coupon string
        Fd4Compound string
        Fd4CurrentRadio string
        Fd4CurInterest string
        Fd4BondPrice string
        Fd4Result string
      } { "Bonds", m.DTF(), p.currentButton,
          p.fd4FaceValue, p.fd4Time, p.fd4TimePeriod, p.fd4Coupon, p.fd4Compound, p.fd4CurrentRadio, p.fd4CurInterest, p.fd4BondPrice, p.fd4Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui5") {
      p.currentButton = "lhs-button5"
      if req.Method == http.MethodPost {
        p.fd5FaceValue = req.PostFormValue("fd5-facevalue")
        p.fd5Time = req.PostFormValue("fd5-time")
        p.fd5TimePeriod = req.PostFormValue("fd5-tp")
        p.fd5Coupon = req.PostFormValue("fd5-coupon")
        p.fd5CurInterest = req.PostFormValue("fd5-current")
        p.fd5Compound = req.PostFormValue("fd5-compound")
        var fv float64
        var time float64
        var couponRate float64
        var curInterest float64
        var err error
        if fv, err = strconv.ParseFloat(p.fd5FaceValue, 64); err != nil {
          p.fd5Result = fmt.Sprintf("Error: %s -- %+v", p.fd5FaceValue, err)
        } else if time, err = strconv.ParseFloat(p.fd5Time, 64); err != nil {
          p.fd5Result = fmt.Sprintf("Error: %s -- %+v", p.fd5Time, err)
        } else if couponRate, err = strconv.ParseFloat(p.fd5Coupon, 64); err != nil {
          p.fd5Result = fmt.Sprintf("Error: %s -- %+v", p.fd5Coupon, err)
        } else if curInterest, err = strconv.ParseFloat(p.fd5CurInterest, 64); err != nil {
          p.fd5Result = fmt.Sprintf("Error: %s -- %+v", p.fd5CurInterest, err)
        } else {
          var b finances.Bonds
          var cp int = b.GetCompoundingPeriod(p.fd5Compound[0], false)
          var tp = b.GetTimePeriod(p.fd5TimePeriod[0], false)
          cf := b.CashFlow(fv, couponRate, cp, time, tp)
          if cp != finances.Continuously {
            p.fd5Result = fmt.Sprintf("Duration: %.3f%%", b.Duration(cf, cp, curInterest,
                                                              b.CurrentPrice(cf, curInterest, cp)))
          } else {
            p.fd5Result = "-1.00"
          }
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s",
                      p.fd5FaceValue, p.fd5Time, p.fd5TimePeriod, p.fd5Coupon, p.fd5Compound,
                      p.fd5CurInterest, p.fd5Result),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/bonds/duration.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd5FaceValue string
        Fd5Time string
        Fd5TimePeriod string
        Fd5Coupon string
        Fd5CurInterest string
        Fd5Compound string
        Fd5Result string
      } { "Bonds", m.DTF(), p.currentButton,
          p.fd5FaceValue, p.fd5Time, p.fd5TimePeriod, p.fd5Coupon, p.fd5CurInterest, p.fd5Compound, p.fd5Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui6") {
      p.currentButton = "lhs-button6"
      if req.Method == http.MethodPost {
        p.fd6FaceValue = req.PostFormValue("fd6-facevalue")
        p.fd6Time = req.PostFormValue("fd6-time")
        p.fd6TimePeriod = req.PostFormValue("fd6-tp")
        p.fd6Coupon = req.PostFormValue("fd6-coupon")
        p.fd6CurInterest = req.PostFormValue("fd6-current")
        p.fd6Compound = req.PostFormValue("fd6-compound")
        var fv float64
        var time float64
        var couponRate float64
        var curInterest float64
        var err error
        if fv, err = strconv.ParseFloat(p.fd6FaceValue, 64); err != nil {
          p.fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", p.fd6FaceValue, err)
        } else if time, err = strconv.ParseFloat(p.fd6Time, 64); err != nil {
          p.fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", p.fd6Time, err)
        } else if couponRate, err = strconv.ParseFloat(p.fd6Coupon, 64); err != nil {
          p.fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", p.fd6Coupon, err)
        } else if curInterest, err = strconv.ParseFloat(p.fd6CurInterest, 64); err != nil {
          p.fd6Result[1] = fmt.Sprintf("Error: %s -- %+v", p.fd6CurInterest, err)
        } else {
          var b finances.Bonds
          var cp int = b.GetCompoundingPeriod(p.fd6Compound[0], false)
          var tp = b.GetTimePeriod(p.fd6TimePeriod[0], false)
          cf := b.CashFlow(fv, couponRate, cp, time, tp)
          if cp != finances.Continuously {
            p.fd6Result[1] = fmt.Sprintf("Macaulay Duration: %.3f year(s)", b.MacaulayDuration(cf, cp,
                                                              b.CurrentPrice(cf, curInterest, cp)))
          } else {
            p.fd6Result[1] = "-1.00"
          }
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s",
                      p.fd6FaceValue, p.fd6Time, p.fd6TimePeriod, p.fd6Coupon, p.fd6Compound,
                      p.fd6CurInterest, p.fd6Result[1]),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/bonds/macaulayduration.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd6FaceValue string
        Fd6Time string
        Fd6TimePeriod string
        Fd6Coupon string
        Fd6CurInterest string
        Fd6Compound string
        Fd6Result [2]string
      } { "Bonds", m.DTF(), p.currentButton,
          p.fd6FaceValue, p.fd6Time, p.fd6TimePeriod, p.fd6Coupon, p.fd6CurInterest, p.fd6Compound, p.fd6Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui7") {
      p.currentButton = "lhs-button7"
      if req.Method == http.MethodPost {
        p.fd7FaceValue = req.PostFormValue("fd7-facevalue")
        p.fd7Time = req.PostFormValue("fd7-time")
        p.fd7TimePeriod = req.PostFormValue("fd7-tp")
        p.fd7Coupon = req.PostFormValue("fd7-coupon")
        p.fd7CurInterest = req.PostFormValue("fd7-current")
        p.fd7Compound = req.PostFormValue("fd7-compound")
        var fv float64
        var time float64
        var couponRate float64
        var curInterest float64
        var err error
        if fv, err = strconv.ParseFloat(p.fd7FaceValue, 64); err != nil {
          p.fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", p.fd7FaceValue, err)
        } else if time, err = strconv.ParseFloat(p.fd7Time, 64); err != nil {
          p.fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", p.fd7Time, err)
        } else if couponRate, err = strconv.ParseFloat(p.fd7Coupon, 64); err != nil {
          p.fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", p.fd7Coupon, err)
        } else if curInterest, err = strconv.ParseFloat(p.fd7CurInterest, 64); err != nil {
          p.fd7Result[1] = fmt.Sprintf("Error: %s -- %+v", p.fd7CurInterest, err)
        } else {
          var b finances.Bonds
          var cp int = b.GetCompoundingPeriod(p.fd7Compound[0], false)
          var tp = b.GetTimePeriod(p.fd7TimePeriod[0], false)
          cf := b.CashFlow(fv, couponRate, cp, time, tp)
          if cp != finances.Continuously {
            p.fd7Result[1] = fmt.Sprintf("Modified Duration: %.3f%%", b.ModifiedDuration(cf, cp,
                                                              b.CurrentPrice(cf, curInterest, cp)))
          } else {
            p.fd7Result[1] = "-1.00"
          }
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s",
                      p.fd7FaceValue, p.fd7Time, p.fd7TimePeriod, p.fd7Coupon, p.fd7Compound,
                      p.fd7CurInterest, p.fd7Result[1]),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
                                                 "webfinances/templates/header.html",
                                                 "webfinances/templates/bonds/modifiedduration.html",
                                                 "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd7FaceValue string
        Fd7Time string
        Fd7TimePeriod string
        Fd7Coupon string
        Fd7CurInterest string
        Fd7Compound string
        Fd7Result [2]string
      } { "Bonds", m.DTF(), p.currentButton,
          p.fd7FaceValue, p.fd7Time, p.fd7TimePeriod, p.fd7Coupon, p.fd7CurInterest, p.fd7Compound, p.fd7Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui8") {
      p.currentButton = "lhs-button8"
      if req.Method == http.MethodPost {
        p.fd8FaceValue = req.PostFormValue("fd8-facevalue")
        p.fd8Time = req.PostFormValue("fd8-time")
        p.fd8TimePeriod = req.PostFormValue("fd8-tp")
        p.fd8Coupon = req.PostFormValue("fd8-coupon")
        p.fd8CurInterest = req.PostFormValue("fd8-current")
        p.fd8Compound = req.PostFormValue("fd8-compound")
        var fv float64
        var time float64
        var couponRate float64
        var curInterest float64
        var err error
        if fv, err = strconv.ParseFloat(p.fd8FaceValue, 64); err != nil {
          p.fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", p.fd8FaceValue, err)
        } else if time, err = strconv.ParseFloat(p.fd8Time, 64); err != nil {
          p.fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", p.fd8Time, err)
        } else if couponRate, err = strconv.ParseFloat(p.fd8Coupon, 64); err != nil {
          p.fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", p.fd8Coupon, err)
        } else if curInterest, err = strconv.ParseFloat(p.fd8CurInterest, 64); err != nil {
          p.fd8Result[1] = fmt.Sprintf("Error: %s -- %+v", p.fd8CurInterest, err)
        } else {
          var b finances.Bonds
          var cp int = b.GetCompoundingPeriod(p.fd8Compound[0], false)
          var tp = b.GetTimePeriod(p.fd8TimePeriod[0], false)
          cf := b.CashFlow(fv, couponRate, cp, time, tp)
          if cp != finances.Continuously {
            p.fd8Result[1] = fmt.Sprintf("Convexity: %.3f", b.Convexity(cf, curInterest, cp))
          } else {
            p.fd8Result[1] = "-1.00"
          }
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("fv = %s, time = %s, tp = %s, coupon = %s, cp = %s, cur interest = %s, %s",
                          p.fd8FaceValue, p.fd8Time, p.fd8TimePeriod, p.fd8Coupon, p.fd8Compound,
                          p.fd8CurInterest, p.fd8Result[1]),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/bonds/convexity.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd8FaceValue string
        Fd8Time string
        Fd8TimePeriod string
        Fd8Coupon string
        Fd8CurInterest string
        Fd8Compound string
        Fd8Result [2]string
      } { "Bonds", m.DTF(), p.currentButton,
          p.fd8FaceValue, p.fd8Time, p.fd8TimePeriod, p.fd8Coupon, p.fd8CurInterest, p.fd8Compound, p.fd8Result,
        })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", p.currentPage)
      fmt.Printf("%s - %s\n", m.DTF(), errString)
      panic(errString)
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
