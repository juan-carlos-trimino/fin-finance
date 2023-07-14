package webfinances

import (
  "finance/finances"
  "fmt"
  "net/http"
  "strconv"
  "strings"
)

type WfSiOrdinaryPages interface {
  SimpleInterestOrdinaryPages(http.ResponseWriter, *http.Request)
}

type wfSiOrdinaryPages struct {
  currentButton string
  fd1Time string
  fd1TimePeriod string
  fd1Interest string
  fd1Compound string
  fd1PV string
  fd1Result string
  fd2Time string
  fd2TimePeriod string
  fd2Amount string
  fd2Compound string
  fd2PV string
  fd2Result string
  fd3Time string
  fd3TimePeriod string
  fd3Interest string
  fd3Compound string
  fd3Amount string
  fd3Result string
  fd4TimePeriod string
  fd4Interest string
  fd4Compound string
  fd4Amount string
  fd4PV string
  fd4Result string
}

func NewWfSiOrdinaryPages() WfSiOrdinaryPages {
  return &wfSiOrdinaryPages {
    currentButton: "lhs-button1",
    fd1Time: "1",
    fd1TimePeriod: "year",
    fd1Interest: "1.00",
    fd1Compound: "annually",
    fd1PV: "1.00",
    fd1Result: "",
    fd2Time: "1",
    fd2TimePeriod: "year",
    fd2Amount: "1.00",
    fd2Compound: "annually",
    fd2PV: "1.00",
    fd2Result: "",
    fd3Time: "1",
    fd3TimePeriod: "year",
    fd3Interest: "1.0",
    fd3Compound: "annually",
    fd3Amount: "1.00",
    fd3Result: "",
    fd4TimePeriod: "year",
    fd4Interest: "1.00",
    fd4Compound: "annually",
    fd4Amount: "1.00",
    fd4PV: "1.00",
    fd4Result: "",
  }
}

func (p *wfSiOrdinaryPages) SimpleInterestOrdinaryPages(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering SimpleInterestOrdinaryPage/webfinances.\n", m.DTF())
  if req.Method == http.MethodPost {
    ui := req.FormValue("compute")
    if strings.EqualFold(ui, "rhs-ui1") {
      p.fd1Time = req.FormValue("fd1-time")
      p.fd1TimePeriod = req.FormValue("fd1-tp")
      p.fd1Interest = req.FormValue("fd1-interest")
      p.fd1Compound = req.FormValue("fd1-compound")
      p.fd1PV = req.FormValue("fd1-pv")
      p.currentButton = "lhs-button1"
      var n float64
      var i float64
      var pv float64
      var err error
      if n, err = strconv.ParseFloat(p.fd1Time, 64); err != nil {
        p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1Time, err)
      } else if i, err = strconv.ParseFloat(p.fd1Interest, 64); err != nil {
        p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1Interest, err)
      } else if pv, err = strconv.ParseFloat(p.fd1PV, 64); err != nil {
        p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1PV, err)
      } else {
        var si finances.SimpleInterest
        var periods finances.Periods
        p.fd1Result = fmt.Sprintf("Amount of Interest: $%.2f", si.OrdinaryInterest(pv, i / 100.0,
                                   periods.GetCompoundingPeriod(p.fd1Compound[0], false), n,
                                   periods.GetTimePeriod(p.fd1TimePeriod[0], false)))
      }
      fmt.Printf("%s - n = %s, tp = %s, i = %s, cp = %s, pv = %s, %s\n", m.DTF(), p.fd1Time,
                  p.fd1TimePeriod, p.fd1Interest, p.fd1Compound, p.fd1PV, p.fd1Result)
    } else if strings.EqualFold(ui, "rhs-ui2") {
      p.fd2Time = req.FormValue("fd2-time")
      p.fd2TimePeriod = req.FormValue("fd2-tp")
      p.fd2Amount = req.FormValue("fd2-amount")
      p.fd2Compound = req.FormValue("fd2-compound")
      p.fd2PV = req.FormValue("fd2-pv")
      p.currentButton = "lhs-button2"
      var n float64
      var a float64
      var pv float64
      var err error
      if n, err = strconv.ParseFloat(p.fd2Time, 64); err != nil {
        p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2Time, err)
      } else if a, err = strconv.ParseFloat(p.fd2Amount, 64); err != nil {
        p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2Amount, err)
      } else if pv, err = strconv.ParseFloat(p.fd2PV, 64); err != nil {
        p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2PV, err)
      } else {
        var si finances.SimpleInterest
        var periods finances.Periods
        p.fd2Result = fmt.Sprintf("Interest Rate: %.3f%%", si.OrdinaryRate(pv, a,
                                   periods.GetCompoundingPeriod(p.fd2Compound[0], false), n,
                                   periods.GetTimePeriod(p.fd2TimePeriod[0], false)))
      }
      fmt.Printf("%s - n = %s, tp = %s, a = %s, cp = %s, pv = %s, %s\n", m.DTF(), p.fd2Time,
                  p.fd2TimePeriod, p.fd2Amount, p.fd2Compound, p.fd2PV, p.fd2Result)
    } else if strings.EqualFold(ui, "rhs-ui3") {
      p.fd3Time = req.FormValue("fd3-time")
      p.fd3TimePeriod = req.FormValue("fd3-tp")
      p.fd3Interest = req.FormValue("fd3-interest")
      p.fd3Compound = req.FormValue("fd3-compound")
      p.fd3Amount = req.FormValue("fd3-amount")
      p.currentButton = "lhs-button3"
      var n float64
      var i float64
      var a float64
      var err error
      if n, err = strconv.ParseFloat(p.fd3Time, 64); err != nil {
        p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3Time, err)
      } else if i, err = strconv.ParseFloat(p.fd3Interest, 64); err != nil {
        p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3Interest, err)
      } else if a, err = strconv.ParseFloat(p.fd3Amount, 64); err != nil {
        p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3Amount, err)
      } else {
        var si finances.SimpleInterest
        var periods finances.Periods
        p.fd3Result = fmt.Sprintf("Principal: $%.2f", si.OrdinaryPrincipal(a, i / 100.0,
                                   periods.GetCompoundingPeriod(p.fd2Compound[0], false), n,
                                   periods.GetTimePeriod(p.fd2TimePeriod[0], false)))
      }
      fmt.Printf("%s - n = %s, tp = %s, i = %s, cp = %s, a = %s, %s\n", m.DTF(), p.fd3Time,
                  p.fd3TimePeriod, p.fd3Interest, p.fd3Compound, p.fd3Amount, p.fd3Result)
    } else if strings.EqualFold(ui, "rhs-ui4") {
      p.fd4TimePeriod = req.FormValue("fd4-tp")
      p.fd4Interest = req.FormValue("fd4-interest")
      p.fd4Compound = req.FormValue("fd4-compound")
      p.fd4Amount = req.FormValue("fd4-amount")
      p.fd4PV = req.FormValue("fd4-pv")
      p.currentButton = "lhs-button4"
      var i float64
      var a float64
      var pv float64
      var err error
      if i, err = strconv.ParseFloat(p.fd4Interest, 64); err != nil {
        p.fd4Result = fmt.Sprintf("Error: %s -- %+v", p.fd4Interest, err)
      } else if a, err = strconv.ParseFloat(p.fd4Amount, 64); err != nil {
        p.fd4Result = fmt.Sprintf("Error: %s -- %+v", p.fd4Amount, err)
      } else if pv, err = strconv.ParseFloat(p.fd4PV, 64); err != nil {
        p.fd4Result = fmt.Sprintf("Error: %s -- %+v", p.fd4PV, err)
      } else {
        var si finances.SimpleInterest
        var periods finances.Periods
        p.fd4Result = fmt.Sprintf("Time: %.2f %s(s)", si.OrdinaryTime(pv, a, i / 100.0,
                                   periods.GetCompoundingPeriod(p.fd4Compound[0], false),
                                   periods.GetTimePeriod(p.fd4TimePeriod[0], false)),
                                   p.fd4TimePeriod)
      }
      fmt.Printf("%s - tp = %s, i = %s, cp = %s, a = %s, i = %s, pv = %s, %s\n", m.DTF(), p.fd3Time,
                  p.fd4TimePeriod, p.fd4Interest, p.fd4Compound, p.fd4Amount, p.fd4PV, p.fd4Result)
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
  tmpl.ExecuteTemplate(res, "simpleinterestordinary.html", struct {
    Header string
    Datetime string
    CurrentButton string
    Fd1Time string
    Fd1TimePeriod string
    Fd1Interest string
    Fd1Compound string
    Fd1PV string
    Fd1Result string
    Fd2Time string
    Fd2TimePeriod string
    Fd2Amount string
    Fd2Compound string
    Fd2PV string
    Fd2Result string
    Fd3Time string
    Fd3TimePeriod string
    Fd3Interest string
    Fd3Compound string
    Fd3Amount string
    Fd3Result string
    Fd4TimePeriod string
    Fd4Interest string
    Fd4Compound string
    Fd4Amount string
    Fd4PV string
    Fd4Result string
  } { "Ordinary Interest", m.DTF(), p.currentButton,
      p.fd1Time, p.fd1TimePeriod, p.fd1Interest, p.fd1Compound, p.fd1PV, p.fd1Result,
      p.fd2Time, p.fd2TimePeriod, p.fd2Amount, p.fd2Compound, p.fd2PV, p.fd2Result,
      p.fd3Time, p.fd3TimePeriod, p.fd3Interest, p.fd3Compound, p.fd3Amount, p.fd3Result,
      p.fd4TimePeriod, p.fd4Interest, p.fd4Compound, p.fd4Amount, p.fd4PV, p.fd4Result })
}
