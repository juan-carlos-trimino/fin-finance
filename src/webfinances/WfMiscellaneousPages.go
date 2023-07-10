package webfinances

import (
  "finance/finances"
  // "finance/misc"
  "fmt"
  // "html/template"
  "net/http"
  "strconv"
  "strings"
)

type WfMiscellaneousPages interface {
  MiscellaneousPage(http.ResponseWriter, *http.Request)
}

type wfMiscellaneousPages struct {
  currentButton string
  fd1Nominal string
  fd1Compound string
  fd1Result [2]string
  fd2Effective string
  fd2Compound string
  fd2Result [2]string
  fd3Nominal string
  fd3Inflation string
  fd3Result [4]string
}

var notes1 = [...]string {
  "When comparing interest rates, use effective annual rates.",
}
var notes3 = [...]string {
  "Nominal returns are not adjusted for inflation.",
  "Real returns are useful while comparing returns over different time periods because of the differences in inflation rates.",
  "Real returns are adjusted for inflation.",
}

func NewWfMiscellaneousPages() WfMiscellaneousPages {
  return &wfMiscellaneousPages {
    currentButton: "lhs-button1",
    fd1Nominal: "3.5",
    fd1Compound: "monthly",
    fd1Result: [2]string { "", "" },
    fd2Effective: "3.5",
    fd2Compound: "monthly",
    fd2Result: [2]string { "", "" },
    fd3Nominal: "2.0",
    fd3Inflation: "2.0",
    fd3Result: [4]string { notes3[0], "", "", "" },
  }
}

func (p *wfMiscellaneousPages) MiscellaneousPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering MiscellaneousPage/webfinances.\n", m.DTF())
  if req.Method == http.MethodGet {
    tmpl.ExecuteTemplate(res, "miscellaneous.html", struct {
      Header string
      Datetime string
      CurrentButton string
      Fd1Nominal string
      Fd1Compound string
      Fd1Result [2]string
      Fd2Effective string
      Fd2Compound string
      Fd2Result [2]string
      Fd3Nominal string
      Fd3Inflation string
      Fd3Result [4]string
    } { "Miscellaneous", m.DTF(), p.currentButton,
        p.fd1Nominal, p.fd1Compound, p.fd1Result,
        p.fd2Effective, p.fd2Compound, p.fd2Result,
        p.fd3Nominal, p.fd3Inflation, p.fd3Result })
  } else if req.Method == http.MethodPost {
    ui := req.FormValue("compute")
    if strings.EqualFold(ui, "rhs-ui1") {
      p.fd1Nominal = req.FormValue("fd1-nominal")
      p.fd1Compound = req.FormValue("fd1-compound")
      p.currentButton = "lhs-button1"
      var nr float64
      var err error
      if nr, err = strconv.ParseFloat(p.fd1Nominal, 64); err != nil {
        p.fd1Result[0] = ""
        p.fd1Result[1] = fmt.Sprintf("Nominal Rate: %s -- %+v", p.fd1Nominal, err)
      } else {
        var a finances.Annuities
        p.fd1Result[0] = notes1[0]
        p.fd1Result[1] = fmt.Sprintf("Effective Annual Rate: %.3f%%", a.NominalRateToEAR(nr / 100.0,
                                      a.GetCompoundingPeriod(p.fd1Compound[0], false)) * 100.0)
      }
      fmt.Printf("%s - nominal rate = %s, cp = %s, %s\n", m.DTF(), p.fd1Nominal, p.fd1Compound,
                 p.fd1Result[1])
    } else if strings.EqualFold(ui, "rhs-ui2") {
      p.fd2Effective = req.FormValue("fd2-effective")
      p.fd2Compound = req.FormValue("fd2-compound")
      p.currentButton = "lhs-button2"
      var ear float64
      var err error
      if ear, err = strconv.ParseFloat(p.fd2Effective, 64); err != nil {
        p.fd2Result[0] = ""
        p.fd2Result[1] = fmt.Sprintf("Effective Rate: %s -- %+v", p.fd2Effective, err)
      } else {
        var a finances.Annuities
        p.fd2Result[0] = notes1[0]
        p.fd2Result[1] = fmt.Sprintf("Nominal Rate: %.3f%% %s", a.EARToNominalRate(ear / 100.0,
                                      a.GetCompoundingPeriod(p.fd2Compound[0], false)) * 100.0,
                                      p.fd2Compound)
      }
      fmt.Printf("%s - effective rate = %s, cp = %s, %s\n", m.DTF(), p.fd2Effective, p.fd2Compound, p.fd2Result[1])
    } else if strings.EqualFold(ui, "rhs-ui3") {
      p.fd3Nominal = req.FormValue("fd3-nominal")
      p.fd3Inflation = req.FormValue("fd3-inflation")
      p.currentButton = "lhs-button3"
      var nr float64
      var ir float64
      var err error
      if nr, err = strconv.ParseFloat(p.fd3Nominal, 64); err != nil {
        p.fd3Result[1] = ""
        p.fd3Result[2] = ""
        p.fd3Result[3] = fmt.Sprintf("Nominal Rate: %s -- %+v", p.fd3Nominal, err)
      } else if ir, err = strconv.ParseFloat(p.fd3Inflation, 64); err != nil {
        p.fd3Result[1] = ""
        p.fd3Result[2] = ""
        p.fd3Result[3] = fmt.Sprintf("Inflation Rate: %s -- %+v", p.fd3Inflation, err)
      } else {
        var a finances.Annuities
        p.fd3Result[1] = notes3[1]
        p.fd3Result[2] = notes3[2]
        p.fd3Result[3] = fmt.Sprintf("Real Interest Rate: %.3f%%", a.RealInterestRate(nr / 100.0,
                                      ir / 100.0) * 100.0)
      }
      fmt.Printf("%s - nominal rate = %s, inflation rate = %s, %s\n", m.DTF(), p.fd3Nominal, p.fd3Inflation, p.fd3Result[3])
    }
    tmpl.ExecuteTemplate(res, "miscellaneous.html", struct {
      Header string
      Datetime string
      CurrentButton string
      Fd1Nominal string
      Fd1Compound string
      Fd1Result [2]string
      Fd2Effective string
      Fd2Compound string
      Fd2Result [2]string
      Fd3Nominal string
      Fd3Inflation string
      Fd3Result [4]string
    } { "Miscellaneous", m.DTF(), p.currentButton,
        p.fd1Nominal, p.fd1Compound, p.fd1Result,
        p.fd2Effective, p.fd2Compound, p.fd2Result,
        p.fd3Nominal, p.fd3Inflation, p.fd3Result })
  } else {
    fmt.Printf("zzzzzzzzzzzzxxxxxxxxxxxxxx\n")
  }
}
