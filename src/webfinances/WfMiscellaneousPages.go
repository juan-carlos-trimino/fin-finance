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
  result1 bool
  nominalRate1 string
  compoundPeriod1 string
  answer1 string
  result2 bool
  effectiveRate2 string
  compoundPeriod2 string
  answer2 string
}

func NewWfMiscellaneousPages() WfMiscellaneousPages {
  return &wfMiscellaneousPages {
    result1: false,
    nominalRate1: "3.5",
    compoundPeriod1: "monthly",
    answer1: "",
    result2: false,
    effectiveRate2: "3.5",
    compoundPeriod2: "monthly",
    answer2: "",
  }
}

func (p *wfMiscellaneousPages) MiscellaneousPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering MiscellaneousPage/webfinances.\n", m.DTF())
  if req.Method == http.MethodGet {

    fmt.Printf("%s\n", req.FormValue("buttons"))

    tmpl.ExecuteTemplate(res, "miscellaneous.html", struct {
      Header string
      Datetime string
      Result1 bool
      NominalRate1 string
      CompoundPeriod1 string
      Answer1 string
      Result2 bool
      EffectiveRate2 string
      CompoundPeriod2 string
      Answer2 string
    } { "Miscellaneous", m.DTF(), p.result1, p.nominalRate1, p.compoundPeriod1, p.answer1,
        p.result2, p.effectiveRate2, p.compoundPeriod2, p.answer2 })
  } else if req.Method == http.MethodPost {
    if strings.EqualFold(req.FormValue("compute"), "uinominal") {
      p.nominalRate1 = req.FormValue("nominal")
      p.compoundPeriod1 = req.FormValue("cp")
      var nr float64
      var err error
      p.result1 = true
      if nr, err = strconv.ParseFloat(p.nominalRate1, 64); err != nil {
        p.answer1 = fmt.Sprintf("Nominal Rate: %s -- %+v", p.nominalRate1, err)
      } else {
        var a finances.Annuities
        p.answer1 = fmt.Sprintf("Effective Annual Rate: %.3f%%", a.NominalRateToEAR(nr / 100.0,
                              a.GetCompoundingPeriod(p.compoundPeriod1[0], false)) * 100.0)
      }
      fmt.Printf("%s - nominal rate = %s, cp = %s, %s\n", m.DTF(), p.nominalRate1, p.compoundPeriod1, p.answer1)
      tmpl.ExecuteTemplate(res, "miscellaneous.html", struct {
        Header string
        Datetime string
        Result1 bool
        NominalRate1 string
        CompoundPeriod1 string
        Answer1 string
        Result2 bool
        EffectiveRate2 string
        CompoundPeriod2 string
        Answer2 string
      } { "Miscellaneous", m.DTF(), p.result1, p.nominalRate1, p.compoundPeriod1, p.answer1,
          p.result2, p.effectiveRate2, p.compoundPeriod2, p.answer2 })
    } else if strings.EqualFold(req.FormValue("compute"), "uieffective") {
      p.effectiveRate2 = req.FormValue("effective")
      p.compoundPeriod2 = req.FormValue("cp")
      var er float64
      var err error
      p.result2 = true
      if er, err = strconv.ParseFloat(p.effectiveRate2, 64); err != nil {
        p.answer2 = fmt.Sprintf("Effective Rate: %s -- %+v", p.effectiveRate2, err)
      } else {
        var a finances.Annuities
        p.answer2 = fmt.Sprintf("Nominal Rate: %.3f%%", a.EARToNominalRate(er / 100.0,
                             a.GetCompoundingPeriod(p.compoundPeriod2[0], false)) * 100.0)
      }
      fmt.Printf("%s - effective rate = %s, cp = %s, %s\n", m.DTF(), p.effectiveRate2, p.compoundPeriod2, p.answer2)
      tmpl.ExecuteTemplate(res, "miscellaneous.html", struct {
        Header string
        Datetime string
        Result1 bool
        NominalRate1 string
        CompoundPeriod1 string
        Answer1 string
        Result2 bool
        EffectiveRate2 string
        CompoundPeriod2 string
        Answer2 string
      } { "Miscellaneous", m.DTF(), p.result1, p.nominalRate1, p.compoundPeriod1, p.answer1,
          p.result2, p.effectiveRate2, p.compoundPeriod2, p.answer2 })
    }
  } else {
      fmt.Printf("zzzzzzzzzzzzxxxxxxxxxxxxxx\n")
  }



}
