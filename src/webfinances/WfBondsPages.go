package webfinances

import (
  "finance/middlewares"
  "finance/finances"
  "fmt"
  "net/http"
  "strconv"
  "strings"
)

type WfBondsPages interface {
  BondsPages(http.ResponseWriter, *http.Request)
}

type wfBondsPages struct {
  currentButton string
  fd1TaxFree string
  fd1CityTax string
  fd1StateTax string
  fd1FederalTax string
  fd1Result string
  fd2FaceValue string
  fd2Time string
  fd2TimePeriod string
  fd2Coupon string
  fd2Current string
  fd2Compound string
  fd2Result string
  fd3FaceValue string
  fd3TimeCall string
  fd3TimePeriod string
  fd3Coupon string
  fd3Compound string
  fd3BondPrice string
  fd3CallPrice string
  fd3Result string
  fd4FaceValue string
  fd4Time string
  fd4TimePeriod string
  fd4Coupon string
  fd4Compound string
  fd4CurrentRadio string
  fd4CurInterest string
  fd4BondPrice string
  fd4Result string
}

func NewWfBondsPages() WfBondsPages {
  return &wfBondsPages {
    currentButton: "lhs-button1",
    fd1TaxFree: "3.5",
    fd1CityTax: "0.0",
    fd1StateTax: "1.0",
    fd1FederalTax: "23.0",
    fd1Result: "",
    fd2FaceValue: "1000.00",
    fd2Time: "5",
    fd2TimePeriod: "year",
    fd2Coupon: "3.00",
    fd2Current: "3.5",
    fd2Compound: "semiannually",
    fd2Result: "",
    fd3FaceValue: "1000.00",
    fd3TimeCall: "2",
    fd3TimePeriod: "year",
    fd3Coupon: "2.0",
    fd3Compound: "semiannually",
    fd3BondPrice: "990.00",
    fd3CallPrice: "1050.00",
    fd3Result: "",
    fd4FaceValue: "1000.00",
    fd4Time: "3",
    fd4TimePeriod: "year",
    fd4Coupon: "2.5",
    fd4Compound: "semiannually",
    fd4CurrentRadio: "fd4-curinterest",
    fd4CurInterest: "2.3",
    fd4BondPrice: "1000.00",
    fd4Result: "",
  }
}

func (p *wfBondsPages) BondsPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering BondsPages/webfinances.",
  })
  if req.Method == http.MethodPost {
    ui := req.FormValue("compute")
    if strings.EqualFold(ui, "rhs-ui1") {
      p.fd1TaxFree = req.FormValue("fd1-taxfree")
      p.fd1CityTax = req.FormValue("fd1-citytax")
      p.fd1StateTax = req.FormValue("fd1-statetax")
      p.fd1FederalTax = req.FormValue("fd1-federaltax")
      p.currentButton = "lhs-button1"
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
    } else if strings.EqualFold(ui, "rhs-ui2") {
      p.fd2FaceValue = req.FormValue("fd2-facevalue")
      p.fd2Time = req.FormValue("fd2-time")
      p.fd2TimePeriod = req.FormValue("fd2-tp")
      p.fd2Coupon = req.FormValue("fd2-coupon")
      p.fd2Current = req.FormValue("fd2-current")
      p.fd2Compound = req.FormValue("fd2-compound")
      p.currentButton = "lhs-button2"
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
    } else if strings.EqualFold(ui, "rhs-ui3") {
      p.fd3FaceValue = req.FormValue("fd3-facevalue")
      p.fd3TimeCall = req.FormValue("fd3-timecall")
      p.fd3TimePeriod = req.FormValue("fd3-tp")
      p.fd3Coupon = req.FormValue("fd3-coupon")
      p.fd3BondPrice = req.FormValue("fd3-bondprice")
      p.fd3CallPrice = req.FormValue("fd3-callprice")
      p.fd3Compound = req.FormValue("fd3-compound")
      p.currentButton = "lhs-button3"
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
    } else if strings.EqualFold(ui, "rhs-ui4") {
      p.fd4FaceValue = req.FormValue("fd4-facevalue")
      p.fd4Time = req.FormValue("fd4-time")
      p.fd4TimePeriod = req.FormValue("fd4-tp")
      p.fd4Coupon = req.FormValue("fd4-coupon")
      p.fd4Compound = req.FormValue("fd4-compound")
      p.fd4CurrentRadio = req.FormValue("fd4-choice")
      p.fd4CurInterest = req.FormValue("fd4-ci")
      p.fd4BondPrice = req.FormValue("fd4-bp")
      p.currentButton = "lhs-button4"
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
    } else if strings.EqualFold(ui, "rhs-ui5") {
      // p.fd5Values = req.FormValue("fd5-values")
      // p.currentButton = "lhs-button5"
      // split := strings.Split(p.fd5Values, ";")
      // values := make([]float64, len(split))
      // var err error
      // for i, s := range split {
      //   if values[i], err = strconv.ParseFloat(s, 64); err != nil {
      //     p.fd5Result[1] = fmt.Sprintf("Error: %s -- %+v", s, err)
      //     break;
      //   }
      // }
      // //
      // if err == nil {
      //   var a finances.Annuities
      //   p.fd5Result[1] = fmt.Sprintf("Avg: %.3f%%", a.AverageRateOfReturn(values) * 100.0)
      // }
      // fmt.Printf("%s - values = [%s], %s\n", m.DTF(), p.fd5Values, p.fd5Result[1])
    } else if strings.EqualFold(ui, "rhs-ui6") {
      // p.fd6Time = req.FormValue("fd6-time")
      // p.fd6TimePeriod = req.FormValue("fd6-tp")
      // p.fd6Rate = req.FormValue("fd6-rate")
      // p.fd6Compound = req.FormValue("fd6-compound")
      // p.fd6PV = req.FormValue("fd6-pv")
      // p.currentButton = "lhs-button6"
      // var time float64
      // var rate float64
      // var pv float64
      // var err error
      // if time, err = strconv.ParseFloat(p.fd6Time, 64); err != nil {
      //   p.fd6Result = fmt.Sprintf("Error: %s -- %+v", p.fd6Time, err)
      // } else if rate, err = strconv.ParseFloat(p.fd6Rate, 64); err != nil {
      //   p.fd6Result = fmt.Sprintf("Error: %s -- %+v", p.fd6Rate, err)
      // } else if pv, err = strconv.ParseFloat(p.fd6PV, 64); err != nil {
      //   p.fd6Result = fmt.Sprintf("Error: %s -- %+v", p.fd6PV, err)
      // } else {
      //   var a finances.Annuities
      //   p.fd6Result = fmt.Sprintf("Future Value: %.2f", a.Depreciation(pv, rate / 100.0,
      //                              a.GetCompoundingPeriod(p.fd6Compound[0], false),
      //                              time, a.GetTimePeriod(p.fd6TimePeriod[0], false)))
      // }
      // fmt.Printf("%s - time = %s, tp = %s, rate = %s, cp = %s, pv = %s, %s\n", m.DTF(), p.fd6Time,
      //             p.fd6TimePeriod, p.fd6Rate, p.fd6Compound, p.fd6PV, p.fd6Result)
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", ui)
      fmt.Printf("%s - %s\n", m.DTF(), errString)
      panic(errString)
    }
  } else if req.Method != http.MethodGet {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
  tmpl.ExecuteTemplate(res, "bonds.html", struct {
    Header string
    Datetime string
    CurrentButton string
    Fd1TaxFree string
    Fd1CityTax string
    Fd1StateTax string
    Fd1FederalTax string
    Fd1Result string
    Fd2FaceValue string
    Fd2Time string
    Fd2TimePeriod string
    Fd2Coupon string
    Fd2Current string
    Fd2Compound string
    Fd2Result string
    Fd3FaceValue string
    Fd3TimeCall string
    Fd3TimePeriod string
    Fd3Coupon string
    Fd3Compound string
    Fd3BondPrice string
    Fd3CallPrice string
    Fd3Result string
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
      p.fd1TaxFree, p.fd1CityTax, p.fd1StateTax, p.fd1FederalTax, p.fd1Result,
      p.fd2FaceValue, p.fd2Time, p.fd2TimePeriod, p.fd2Coupon, p.fd2Current, p.fd2Compound, p.fd2Result,
      p.fd3FaceValue, p.fd3TimeCall, p.fd3TimePeriod, p.fd3Coupon, p.fd3Compound, p.fd3BondPrice, p.fd3CallPrice, p.fd3Result,
      p.fd4FaceValue, p.fd4Time, p.fd4TimePeriod, p.fd4Coupon, p.fd4Compound, p.fd4CurrentRadio, p.fd4CurInterest, p.fd4BondPrice, p.fd4Result,
    })
}
