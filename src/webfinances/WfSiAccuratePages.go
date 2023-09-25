package webfinances

import (
  "finance/finances"
  "finance/middlewares"
  "fmt"
  "html/template"
  "net/http"
  "strconv"
  "strings"
)

type WfSiAccuratePages interface {
  SimpleInterestAccuratePages(http.ResponseWriter, *http.Request)
}

type wfSiAccuratePages struct {
  currentPage string
  currentButton string
  //
  fd1Time string
  fd1TimePeriod string
  fd1Interest string
  fd1Compound string
  fd1PV string
  fd1Result string
  //
  fd2Time string
  fd2TimePeriod string
  fd2Amount string
  fd2PV string
  fd2Result string
  //
  fd3Time string
  fd3TimePeriod string
  fd3Interest string
  fd3Compound string
  fd3Amount string
  fd3Result string
  //
  fd4Interest string
  fd4Compound string
  fd4Amount string
  fd4PV string
  fd4Result string
}

func NewWfSiAccuratePages() WfSiAccuratePages {
  return &wfSiAccuratePages {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1Time: "1",
    fd1TimePeriod: "year",
    fd1Interest: "1.00",
    fd1Compound: "annually",
    fd1PV: "1.00",
    fd1Result: "",
    //
    fd2Time: "1",
    fd2TimePeriod: "year",
    fd2Amount: "1.00",
    fd2PV: "1.00",
    fd2Result: "",
    //
    fd3Time: "1",
    fd3TimePeriod: "year",
    fd3Interest: "1.0",
    fd3Compound: "annually",
    fd3Amount: "1.00",
    fd3Result: "",
    //
    fd4Interest: "1.00",
    fd4Compound: "annually",
    fd4Amount: "1.00",
    fd4PV: "1.00",
    fd4Result: "",
  }
}

func (p *wfSiAccuratePages) SimpleInterestAccuratePages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering SimpleInterestAccuratePages/webfinances.",
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
        p.fd1Time = req.PostFormValue("fd1-time")
        p.fd1TimePeriod = req.PostFormValue("fd1-tp")
        p.fd1Interest = req.PostFormValue("fd1-interest")
        p.fd1Compound = req.PostFormValue("fd1-compound")
        p.fd1PV = req.PostFormValue("fd1-pv")
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
          p.fd1Result = fmt.Sprintf("Amount of Interest: $%.2f", si.AccurateInterest(pv, i / 100.0,
                                     periods.GetCompoundingPeriod(p.fd1Compound[0], true), n,
                                     periods.GetTimePeriod(p.fd1TimePeriod[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, pv = %s, %s", p.fd1Time, p.fd1TimePeriod,
                       p.fd1Interest, p.fd1Compound, p.fd1PV, p.fd1Result),
        })
      }
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      t := template.Must(template.ParseFiles("webfinances/templates/simpleinterestaccurate/accurate.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/simpleinterestaccurate/amountofinterest.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "simpleinterestaccurate", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd1Time string
        Fd1TimePeriod string
        Fd1Interest string
        Fd1Compound string
        Fd1PV string
        Fd1Result string
      } { "Simple Interest / Accurate (Exact) Interest", m.DTF(), p.currentButton,
          p.fd1Time, p.fd1TimePeriod, p.fd1Interest, p.fd1Compound, p.fd1PV, p.fd1Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui2") {
      p.currentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        p.fd2Time = req.PostFormValue("fd2-time")
        p.fd2TimePeriod = req.PostFormValue("fd2-tp")
        p.fd2Amount = req.PostFormValue("fd2-amount")
        p.fd2PV = req.PostFormValue("fd2-pv")
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
          p.fd2Result = fmt.Sprintf("Interest Rate: %.3f%%", si.AccurateRate(pv, a,
                                    n, periods.GetTimePeriod(p.fd2TimePeriod[0], true)) * 100.0)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, a = %s, pv = %s, %s", p.fd2Time, p.fd2TimePeriod,
                       p.fd2Amount, p.fd2PV, p.fd2Result),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/simpleinterestaccurate/accurate.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/simpleinterestaccurate/interestrate.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "simpleinterestaccurate", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd2Time string
        Fd2TimePeriod string
        Fd2Amount string
        Fd2PV string
        Fd2Result string
      } { "Simple Interest / Accurate (Exact) Interest", m.DTF(), p.currentButton,
          p.fd2Time, p.fd2TimePeriod, p.fd2Amount, p.fd2PV, p.fd2Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui3") {
      p.currentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        p.fd3Time = req.PostFormValue("fd3-time")
        p.fd3TimePeriod = req.PostFormValue("fd3-tp")
        p.fd3Interest = req.PostFormValue("fd3-interest")
        p.fd3Compound = req.PostFormValue("fd3-compound")
        p.fd3Amount = req.PostFormValue("fd3-amount")
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
          p.fd3Result = fmt.Sprintf("Principal: $%.2f", si.AccuratePrincipal(a, i / 100.0,
                                     periods.GetCompoundingPeriod(p.fd3Compound[0], true), n,
                                     periods.GetTimePeriod(p.fd3TimePeriod[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, a = %s, %s\n", p.fd3Time, p.fd3TimePeriod,
                       p.fd3Interest, p.fd3Compound, p.fd3Amount, p.fd3Result),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/simpleinterestaccurate/accurate.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/simpleinterestaccurate/principal.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "simpleinterestaccurate", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd3Time string
        Fd3TimePeriod string
        Fd3Interest string
        Fd3Compound string
        Fd3Amount string
        Fd3Result string
      } { "Simple Interest / Accurate (Exact) Interest", m.DTF(), p.currentButton,
          p.fd3Time, p.fd3TimePeriod, p.fd3Interest, p.fd3Compound, p.fd3Amount, p.fd3Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui4") {
      p.currentButton = "lhs-button4"
      if req.Method == http.MethodPost {
        p.fd4Interest = req.FormValue("fd4-interest")
        p.fd4Compound = req.FormValue("fd4-compound")
        p.fd4Amount = req.FormValue("fd4-amount")
        p.fd4PV = req.FormValue("fd4-pv")
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
          p.fd4Result = fmt.Sprintf("Time: %.3f %s", si.AccurateTime(pv, a, i / 100.0,
                                     periods.GetCompoundingPeriod(p.fd4Compound[0], true)),
                                     periods.TimePeriods(p.fd4Compound))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("i = %s, cp = %s, a = %s, pv = %s, %s\n", p.fd4Interest, p.fd4Compound,
                      p.fd4Amount, p.fd4PV, p.fd4Result),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/simpleinterestaccurate/accurate.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/simpleinterestaccurate/time.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "simpleinterestaccurate", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd4Interest string
        Fd4Compound string
        Fd4Amount string
        Fd4PV string
        Fd4Result string
      } { "Simple Interest / Accurate (Exact) Interest", m.DTF(), p.currentButton,
          p.fd4Interest, p.fd4Compound, p.fd4Amount, p.fd4PV, p.fd4Result,
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
