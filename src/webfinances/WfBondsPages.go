package webfinances

import (
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

  fd4Factor string
  fd4Result string
  fd5Values string
  fd5Result [2]string
  fd6Time string
  fd6TimePeriod string
  fd6Rate string
  fd6Compound string
  fd6PV string
  fd6Result string
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

    fd4Factor: "2.0",
    fd4Result: "",
    fd5Values: "2.0;1.5",
    fd5Result: [2]string { notes5[0], "" },
    fd6Time: "1.0",
    fd6TimePeriod: "year",
    fd6Rate: "15.0",
    fd6Compound: "annually",
    fd6PV: "1.00",
    fd6Result: "",
  }
}

func (p *wfBondsPages) BondsPages(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering BondsPages/webfinances.\n", m.DTF())
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
      fmt.Printf("%s - tax free = %s, city tax = %s, state tax = %s, federal tax = %s, %s\n",
       m.DTF(), p.fd1TaxFree, p.fd1CityTax, p.fd1StateTax, p.fd1FederalTax, p.fd1Result)
    } else if strings.EqualFold(ui, "rhs-ui2") {
      // p.fd2FaceValue = req.FormValue("fd2-facevalue")
      // p.fd2Time = req.FormValue("fd2-time")
      // p.fd2TimePeriod = req.FormValue("fd2-tp")
      // p.fd2Coupon = req.FormValue("fd2-coupon")
      // p.fd2Current = req.FormValue("fd2-current")
      // p.fd2Compound = req.FormValue("fd2-compound")
      // p.currentButton = "lhs-button2"
      // var fv float64
      // var time float64
      // var coupon float64
      // var current float64
      // var err error
      // if fv, err = strconv.ParseFloat(p.fd2FaceValue, 64); err != nil {
      //   p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2FaceValue, err)
      // } else if time, err = strconv.ParseFloat(p.fd2Time, 64); err != nil {
      //   p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2Time, err)
      // } else if coupon, err = strconv.ParseFloat(p.fd2Coupon, 64); err != nil {
      //   p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2Coupon, err)
      // } else if current, err = strconv.ParseFloat(p.fd2Current, 64); err != nil {
      //   p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2Current, err)
      // } else {
      //   var b finances.Bonds
      //   cf := b.CashFlow(fv, coupon, b.GetCompoundingPeriod(p.fd2Compound[0], false), time,
      //                    b.TimePeriods(p.fd4Compound))
        // p.fd2Result = fmt.Sprintf("Current Price: $%.2f %s", b.CurrentPrice(ear / 100.0)
      //                                 a.GetCompoundingPeriod(p.fd2Compound[0], false)) * 100.0,
      //                                 p.fd2Compound)
      //}
      // fmt.Printf("%s - effective rate = %s, cp = %s, %s\n", m.DTF(), p.fd2Effective, p.fd2Compound, p.fd2Result[1])
    } else if strings.EqualFold(ui, "rhs-ui3") {
      // p.fd3Nominal = req.FormValue("fd3-nominal")
      // p.fd3Inflation = req.FormValue("fd3-inflation")
      p.currentButton = "lhs-button3"
      // var nr float64
      // var ir float64
      // var err error
      // if nr, err = strconv.ParseFloat(p.fd3Nominal, 64); err != nil {
      //   p.fd3Result[1] = ""
      //   p.fd3Result[2] = ""
      //   p.fd3Result[3] = fmt.Sprintf("Error: %s -- %+v", p.fd3Nominal, err)
      // } else if ir, err = strconv.ParseFloat(p.fd3Inflation, 64); err != nil {
      //   p.fd3Result[1] = ""
      //   p.fd3Result[2] = ""
      //   p.fd3Result[3] = fmt.Sprintf("Error: %s -- %+v", p.fd3Inflation, err)
      // } else {
      //   var a finances.Annuities
      //   p.fd3Result[1] = notes3[1]
      //   p.fd3Result[2] = notes3[2]
      //   p.fd3Result[3] = fmt.Sprintf("Real Interest Rate: %.3f%%", a.RealInterestRate(nr / 100.0,
      //                                 ir / 100.0) * 100.0)
      // }
      // fmt.Printf("%s - nominal rate = %s, inflation rate = %s, %s\n", m.DTF(), p.fd3Nominal, p.fd3Inflation, p.fd3Result[3])
    } else if strings.EqualFold(ui, "rhs-ui4") {
      // p.fd4Interest = req.FormValue("fd4-interest")
      // p.fd4Compound = req.FormValue("fd4-compound")
      // p.fd4Factor = req.FormValue("fd4-factor")
      // p.currentButton = "lhs-button4"
      // var ir float64
      // var factor float64
      // var err error
      // if ir, err = strconv.ParseFloat(p.fd4Interest, 64); err != nil {
      //   p.fd4Result = fmt.Sprintf("Error: %s -- %+v", p.fd4Interest, err)
      // } else if factor, err = strconv.ParseFloat(p.fd4Factor, 64); err != nil {
      //   p.fd4Result = fmt.Sprintf("Error: %s -- %+v", p.fd4Factor, err)
      // } else {
      //   var a finances.Annuities
      //   p.fd4Result = fmt.Sprintf("Growth/Decay: %.3f %s", a.GrowthDecayOfFunds(factor, ir / 100.0,
      //     a.GetCompoundingPeriod(p.fd4Compound[0], false)), a.TimePeriods(p.fd4Compound))
      // }
      // fmt.Printf("%s - interest rate = %s, cp = %s, factor = %s, %s\n", m.DTF(), p.fd4Interest,
      //             p.fd4Compound, p.fd4Factor, p.fd4Result)
    } else if strings.EqualFold(ui, "rhs-ui5") {
      p.fd5Values = req.FormValue("fd5-values")
      p.currentButton = "lhs-button5"
      split := strings.Split(p.fd5Values, ";")
      values := make([]float64, len(split))
      var err error
      for i, s := range split {
        if values[i], err = strconv.ParseFloat(s, 64); err != nil {
          p.fd5Result[1] = fmt.Sprintf("Error: %s -- %+v", s, err)
          break;
        }
      }
      //
      if err == nil {
        var a finances.Annuities
        p.fd5Result[1] = fmt.Sprintf("Avg: %.3f%%", a.AverageRateOfReturn(values) * 100.0)
      }
      fmt.Printf("%s - values = [%s], %s\n", m.DTF(), p.fd5Values, p.fd5Result[1])
    } else if strings.EqualFold(ui, "rhs-ui6") {
      p.fd6Time = req.FormValue("fd6-time")
      p.fd6TimePeriod = req.FormValue("fd6-tp")
      p.fd6Rate = req.FormValue("fd6-rate")
      p.fd6Compound = req.FormValue("fd6-compound")
      p.fd6PV = req.FormValue("fd6-pv")
      p.currentButton = "lhs-button6"
      var time float64
      var rate float64
      var pv float64
      var err error
      if time, err = strconv.ParseFloat(p.fd6Time, 64); err != nil {
        p.fd6Result = fmt.Sprintf("Error: %s -- %+v", p.fd6Time, err)
      } else if rate, err = strconv.ParseFloat(p.fd6Rate, 64); err != nil {
        p.fd6Result = fmt.Sprintf("Error: %s -- %+v", p.fd6Rate, err)
      } else if pv, err = strconv.ParseFloat(p.fd6PV, 64); err != nil {
        p.fd6Result = fmt.Sprintf("Error: %s -- %+v", p.fd6PV, err)
      } else {
        var a finances.Annuities
        p.fd6Result = fmt.Sprintf("Future Value: %.2f", a.Depreciation(pv, rate / 100.0,
                                   a.GetCompoundingPeriod(p.fd6Compound[0], false),
                                   time, a.GetTimePeriod(p.fd6TimePeriod[0], false)))
      }
      fmt.Printf("%s - time = %s, tp = %s, rate = %s, cp = %s, pv = %s, %s\n", m.DTF(), p.fd6Time,
                  p.fd6TimePeriod, p.fd6Rate, p.fd6Compound, p.fd6PV, p.fd6Result)
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

    // Fd4Result string
    // Fd5Values string
    // Fd5Result [2]string
    // Fd6Time string
    // Fd6TimePeriod string
    // Fd6Rate string
    // Fd6Compound string
    // Fd6PV string
    // Fd6Result string
  } { "Bonds", m.DTF(), p.currentButton,
      p.fd1TaxFree, p.fd1CityTax, p.fd1StateTax, p.fd1FederalTax, p.fd1Result,
      p.fd2FaceValue, p.fd2Time, p.fd2TimePeriod, p.fd2Coupon, p.fd2Current, p.fd2Compound, p.fd2Result,
      // p.fd3Nominal, p.fd3Inflation, p.fd3Result,
      // p.fd4Interest, p.fd4Compound, p.fd4Factor, p.fd4Result,
      // p.fd5Values, p.fd5Result,
      /* p.fd6Time, p.fd6TimePeriod, p.fd6Rate, p.fd6Compound, p.fd6PV, p.fd6Result*/ })
}
