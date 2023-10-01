package webfinances

import (
  "context"
  "finance/finances"
  "finance/middlewares"
  "fmt"
  "html/template"
  "net/http"
  "strconv"
  "strings"
)

var misc_notes = [...]string {
  "When comparing interest rates, use effective annual rates.",
  "Nominal returns are not adjusted for inflation.",
  "Real returns are useful while comparing returns over different time periods because of the differences in inflation rates.",
  "Real returns are adjusted for inflation.",
  "Values are semicolon (;) separated; e.g., 3;3.1;3.2;-1.01",
}

type WfMiscellaneousPages interface {
  MiscellaneousPages(http.ResponseWriter, *http.Request)
}

type wfMiscellaneousPages struct {
  currentPage string
  currentButton string
  //
  fd1Nominal string
  fd1Compound string
  fd1Result [2]string
  //
  fd2Effective string
  fd2Compound string
  fd2Result [3]string
  //
  fd3Nominal string
  fd3Inflation string
  fd3Result [4]string
  //
  fd4Interest string
  fd4Compound string
  fd4Factor string
  fd4Result string
  //
  fd5Values string
  fd5Result [2]string
  //
  fd6Time string
  fd6TimePeriod string
  fd6Rate string
  fd6Compound string
  fd6PV string
  fd6Result string
}

func NewWfMiscellaneousPages() WfMiscellaneousPages {
  return &wfMiscellaneousPages {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1Nominal: "3.5",
    fd1Compound: "monthly",
    fd1Result: [2]string { misc_notes[0], "" },
    //
    fd2Effective: "3.5",
    fd2Compound: "monthly",
    fd2Result: [3]string { misc_notes[0], misc_notes[1], "" },
    //
    fd3Nominal: "2.0",
    fd3Inflation: "2.0",
    fd3Result: [4]string { misc_notes[1], misc_notes[2], misc_notes[3], "" },
    //
    fd4Interest: "14.87",
    fd4Compound: "annually",
    fd4Factor: "2.0",
    fd4Result: "",
    //
    fd5Values: "2.0;1.5",
    fd5Result: [2]string { misc_notes[4], "" },
    //
    fd6Time: "1.0",
    fd6TimePeriod: "year",
    fd6Rate: "15.0",
    fd6Compound: "annually",
    fd6PV: "1.00",
    fd6Result: "",
  }
}

func (p *wfMiscellaneousPages) MiscellaneousPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering MiscellaneousPages/webfinances.",
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
        p.fd1Nominal = req.PostFormValue("fd1-nominal")
        p.fd1Compound = req.PostFormValue("fd1-compound")
        var nr float64
        var err error
        if nr, err = strconv.ParseFloat(p.fd1Nominal, 64); err != nil {
          p.fd1Result[1] = fmt.Sprintf("Error: %s -- %+v", p.fd1Nominal, err)
        } else {
          var a finances.Annuities
          p.fd1Result[1] = fmt.Sprintf("Effective Annual Rate: %.3f%%", a.NominalRateToEAR(nr / 100.0,
                                       a.GetCompoundingPeriod(p.fd1Compound[0], false)) * 100.0)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("nominal rate = %s, cp = %s, %s", p.fd1Nominal, p.fd1Compound, p.fd1Result[1]),
        })
      }
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      t := template.Must(template.ParseFiles("webfinances/templates/miscellaneous/miscellaneous.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/miscellaneous/nominalrate.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "miscellaneous", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd1Nominal string
        Fd1Compound string
        Fd1Result [2]string
      } { "Miscellaneous", m.DTF(), p.currentButton,
          p.fd1Nominal, p.fd1Compound, p.fd1Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui2") {
      p.currentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        p.fd2Effective = req.PostFormValue("fd2-effective")
        p.fd2Compound = req.PostFormValue("fd2-compound")
        var ear float64
        var err error
        if ear, err = strconv.ParseFloat(p.fd2Effective, 64); err != nil {
          p.fd2Result[2] = fmt.Sprintf("Error: %s -- %+v", p.fd2Effective, err)
        } else {
          var a finances.Annuities
          p.fd2Result[2] = fmt.Sprintf("Nominal Rate: %.3f%% %s", a.EARToNominalRate(ear / 100.0,
                                       a.GetCompoundingPeriod(p.fd2Compound[0], false)) * 100.0,
                                       p.fd2Compound)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("effective rate = %s, cp = %s, %s", p.fd2Effective, p.fd2Compound,
                      p.fd2Result[2]),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/miscellaneous/miscellaneous.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/miscellaneous/effectiveannualrate.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "miscellaneous", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd2Effective string
        Fd2Compound string
        Fd2Result [3]string
      } { "Miscellaneous", m.DTF(), p.currentButton,
          p.fd2Effective, p.fd2Compound, p.fd2Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui3") {
      p.currentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        p.fd3Nominal = req.PostFormValue("fd3-nominal")
        p.fd3Inflation = req.PostFormValue("fd3-inflation")
        var nr float64
        var ir float64
        var err error
        if nr, err = strconv.ParseFloat(p.fd3Nominal, 64); err != nil {
          p.fd3Result[3] = fmt.Sprintf("Error: %s -- %+v", p.fd3Nominal, err)
        } else if ir, err = strconv.ParseFloat(p.fd3Inflation, 64); err != nil {
          p.fd3Result[3] = fmt.Sprintf("Error: %s -- %+v", p.fd3Inflation, err)
        } else {
          var a finances.Annuities
          p.fd3Result[3] = fmt.Sprintf("Real Interest Rate: %.3f%%", a.RealInterestRate(nr / 100.0,
                                       ir / 100.0) * 100.0)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("nominal rate = %s, inflation rate = %s, %s", p.fd3Nominal, p.fd3Inflation,
                      p.fd3Result[3]),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/miscellaneous/miscellaneous.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/miscellaneous/nominalratevs.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "miscellaneous", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd3Nominal string
        Fd3Inflation string
        Fd3Result [4]string
      } { "Miscellaneous", m.DTF(), p.currentButton,
          p.fd3Nominal, p.fd3Inflation, p.fd3Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui4") {
      p.currentButton = "lhs-button4"
      if req.Method == http.MethodPost {
        p.fd4Interest = req.PostFormValue("fd4-interest")
        p.fd4Compound = req.PostFormValue("fd4-compound")
        p.fd4Factor = req.PostFormValue("fd4-factor")
        var ir float64
        var factor float64
        var err error
        if ir, err = strconv.ParseFloat(p.fd4Interest, 64); err != nil {
          p.fd4Result = fmt.Sprintf("Error: %s -- %+v", p.fd4Interest, err)
        } else if factor, err = strconv.ParseFloat(p.fd4Factor, 64); err != nil {
          p.fd4Result = fmt.Sprintf("Error: %s -- %+v", p.fd4Factor, err)
        } else {
          var a finances.Annuities
          p.fd4Result = fmt.Sprintf("Growth/Decay: %.3f %s", a.GrowthDecayOfFunds(factor,
                                    ir / 100.0, a.GetCompoundingPeriod(p.fd4Compound[0], true)),
                                    a.TimePeriods(p.fd4Compound))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("interest rate = %s, cp = %s, factor = %s, %s\n", p.fd4Interest,
                      p.fd4Compound, p.fd4Factor, p.fd4Result),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/miscellaneous/miscellaneous.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/miscellaneous/growthdecay.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "miscellaneous", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd4Interest string
        Fd4Compound string
        Fd4Factor string
        Fd4Result string
      } { "Miscellaneous", m.DTF(), p.currentButton,
          p.fd4Interest, p.fd4Compound, p.fd4Factor, p.fd4Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui5") {
      p.currentButton = "lhs-button5"
      if req.Method == http.MethodPost {
        p.fd5Values = req.PostFormValue("fd5-values")
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
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("values = [%s], %s\n", p.fd5Values, p.fd5Result[1]),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/miscellaneous/miscellaneous.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/miscellaneous/averagerate.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "miscellaneous", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd5Values string
        Fd5Result [2]string
      } { "Miscellaneous", m.DTF(), p.currentButton,
          p.fd5Values, p.fd5Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui6") {
      p.currentButton = "lhs-button6"
      if req.Method == http.MethodPost {
        p.fd6Time = req.PostFormValue("fd6-time")
        p.fd6TimePeriod = req.PostFormValue("fd6-tp")
        p.fd6Rate = req.PostFormValue("fd6-rate")
        p.fd6Compound = req.PostFormValue("fd6-compound")
        p.fd6PV = req.PostFormValue("fd6-pv")
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
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("time = %s, tp = %s, rate = %s, cp = %s, pv = %s, %s\n", p.fd6Time,
                      p.fd6TimePeriod, p.fd6Rate, p.fd6Compound, p.fd6PV, p.fd6Result),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/miscellaneous/miscellaneous.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/miscellaneous/depreciation.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "miscellaneous", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd6Time string
        Fd6TimePeriod string
        Fd6Rate string
        Fd6Compound string
        Fd6PV string
        Fd6Result string
      } { "Miscellaneous", m.DTF(), p.currentButton,
          p.fd6Time, p.fd6TimePeriod, p.fd6Rate, p.fd6Compound, p.fd6PV, p.fd6Result,
        })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", p.currentPage)
      fmt.Printf("%s - %s\n", m.DTF(), errString)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      fmt.Println("*** Request timeout ***")
      if strings.EqualFold(p.currentPage, "rhs-ui1") {
        p.fd1Result[1] = ""
      } else if strings.EqualFold(p.currentPage, "rhs-ui2") {
        p.fd2Result[2] = ""
      } else if strings.EqualFold(p.currentPage, "rhs-ui3") {
        p.fd3Result[3] = ""
      } else if strings.EqualFold(p.currentPage, "rhs-ui4") {
        p.fd4Result = ""
      } else if strings.EqualFold(p.currentPage, "rhs-ui5") {
        p.fd5Result[1] = ""
      } else if strings.EqualFold(p.currentPage, "rhs-ui6") {
        p.fd6Result = ""
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
