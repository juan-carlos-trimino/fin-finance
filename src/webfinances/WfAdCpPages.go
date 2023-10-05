package webfinances

import (
  "context"
  "finance/middlewares"
  "finance/finances"
  "fmt"
  "html/template"
  "net/http"
  "strconv"
  "strings"
)

type WfAdCpPages interface {
  AdCpPages(http.ResponseWriter, *http.Request)
}

type wfAdCpPages struct {
  currentPage string
  currentButton string
  //
  fd1Interest string
  fd1Compound string
  fd1PV string
  fd1FV string
  fd1Result string
  //
  fd2Interest string
  fd2Compound string
  fd2Payment string
  fd2PV string
  fd2Result string
  //
  fd3Interest string
  fd3Compound string
  fd3Payment string
  fd3FV string
  fd3Result string
}

func NewWfAdCpPages() WfAdCpPages {
  return &wfAdCpPages {
    currentPage: "rhs-ui2",
    currentButton: "lhs-button2",
    //
    fd1Interest: "1.00",
    fd1Compound: "annually",
    fd1PV: "1.00",
    fd1FV: "1.00",
    fd1Result: "",
    //
    fd2Interest: "1.00",
    fd2Compound: "annually",
    fd2Payment: "1.00",
    fd2PV: "1.00",
    fd2Result: "",
    //
    fd3Interest: "1.00",
    fd3Compound: "annually",
    fd3Payment: "1.00",
    fd3FV: "1.00",
    fd3Result: "",
  }
}

func (p *wfAdCpPages) AdCpPages(res http.ResponseWriter, req *http.Request) {
  if !checkSession(res, req) {
    return
  }
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering AdCpPages/webfinances.",
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
    // if strings.EqualFold(p.currentPage, "rhs-ui1") {
    //   p.currentButton = "lhs-button1"
    //   if req.Method == http.MethodPost {
    //     p.fd1Interest = req.PostFormValue("fd1-interest")
    //     p.fd1Compound = req.PostFormValue("fd1-cp")
    //     p.fd1PV = req.PostFormValue("fd1-pv")
    //     p.fd1FV = req.PostFormValue("fd1-fv")
    //     var i float64
    //     var pv float64
    //     var fv float64
    //     var err error
    //     if i, err = strconv.ParseFloat(p.fd1Interest, 64); err != nil {
    //       p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1Interest, err)
    //     } else if pv, err = strconv.ParseFloat(p.fd1PV, 64); err != nil {
    //       p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1PV, err)
    //     } else if fv, err = strconv.ParseFloat(p.fd1FV, 64); err != nil {
    //       p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1FV, err)
    //     } else {
    //       var oa finances.Annuities
    //       p.fd1Result = fmt.Sprintf("Compounding Period: %.5f %s", oa.O_Periods_PV_FV(pv, fv,
    //                                 i / 100.0, oa.GetCompoundingPeriod(p.fd1Compound[0], true)),
    //                                 oa.TimePeriods(p.fd1Compound))
    //     }
    //     logEntry.Print(INFO, correlationId, []string {
    //       fmt.Sprintf("i = %s, cp = %s, pv = %s, fv = %s, %s",
    //                   p.fd1Interest, p.fd1Compound, p.fd1PV, p.fd1FV, p.fd1Result),
    //     })
    //   }
    //   /***
    //   The Must function wraps around the ParseGlob function that returns a pointer to a template
    //   and an error, and it panics if the error is not nil.
    //   ***/
    //   t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/cp/cp.html",
    //                                          "webfinances/templates/header.html",
    //                                          "webfinances/templates/ordinaryannuity/cp/i-PV-FV.html",
    //                                          "webfinances/templates/footer.html"))
    //   t.ExecuteTemplate(res, "oacompoundingperiods", struct {
    //     Header string
    //     Datetime string
    //     CurrentButton string
    //     Fd1Interest string
    //     Fd1Compound string
    //     Fd1PV string
    //     Fd1FV string
    //     Fd1Result string
    //   } { "Ordinary Annuity / Compounding Periods", m.DTF(), p.currentButton,
    //       p.fd1Interest, p.fd1Compound, p.fd1PV, p.fd1FV, p.fd1Result,
    //     })
    /*} else*/ if strings.EqualFold(p.currentPage, "rhs-ui2") {
      p.currentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        p.fd2Interest = req.FormValue("fd2-interest")
        p.fd2Compound = req.PostFormValue("fd2-cp")
        p.fd2Payment = req.PostFormValue("fd2-payment")
        p.fd2PV = req.PostFormValue("fd2-pv")
        var i float64
        var pmt float64
        var pv float64
        var err error
        if i, err = strconv.ParseFloat(p.fd2Interest, 64); err != nil {
          p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2Interest, err)
        } else if pmt, err = strconv.ParseFloat(p.fd2Payment, 64); err != nil {
          p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2Payment, err)
        } else if pv, err = strconv.ParseFloat(p.fd2PV, 64); err != nil {
          p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2PV, err)
        } else {
          var oa finances.Annuities
          p.fd2Result = fmt.Sprintf("Compounding Period: %.5f %s", oa.D_Periods_PMT_PV(pmt, pv,
                                    i / 100.0, oa.GetCompoundingPeriod(p.fd2Compound[0], true)),
                                    oa.TimePeriods(p.fd2Compound))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("i = %s, cp = %s, pmt = %s, pv = %s, %s",
                      p.fd2Interest, p.fd2Compound, p.fd2Payment, p.fd2PV, p.fd2Result),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/annuitydue/cp/cp.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/annuitydue/cp/i-PMT-PV.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "adcompoundingperiods", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd2Interest string
        Fd2Compound string
        Fd2Payment string
        Fd2PV string
        Fd2Result string
      } { "Annuity Due / Compounding Periods", m.DTF(), p.currentButton,
          p.fd2Interest, p.fd2Compound, p.fd2Payment, p.fd2PV, p.fd2Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui3") {
      p.currentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        p.fd3Interest = req.FormValue("fd3-interest")
        p.fd3Compound = req.PostFormValue("fd3-cp")
        p.fd3Payment = req.PostFormValue("fd3-payment")
        p.fd3FV = req.PostFormValue("fd3-fv")
        var i float64
        var pmt float64
        var fv float64
        var err error
        if i, err = strconv.ParseFloat(p.fd3Interest, 64); err != nil {
          p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3Interest, err)
        } else if pmt, err = strconv.ParseFloat(p.fd3Payment, 64); err != nil {
          p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3Payment, err)
        } else if fv, err = strconv.ParseFloat(p.fd3FV, 64); err != nil {
          p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3FV, err)
        } else {
          var oa finances.Annuities
          p.fd3Result = fmt.Sprintf("Compounding Period: %.5f %s", oa.D_Periods_PMT_FV(pmt, fv,
                                    i / 100.0, oa.GetCompoundingPeriod(p.fd3Compound[0], true)),
                                    oa.TimePeriods(p.fd3Compound))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("i = %s, cp = %s, pmt = %s, fv = %s, %s", p.fd3Interest,
                      p.fd3Compound, p.fd3Payment, p.fd3FV, p.fd3Result),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/annuitydue/cp/cp.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/annuitydue/cp/i-PMT-FV.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "adcompoundingperiods", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd3Interest string
        Fd3Compound string
        Fd3Payment string
        Fd3FV string
        Fd3Result string
      } { "Annuity Due / Compounding Periods", m.DTF(), p.currentButton,
          p.fd3Interest, p.fd3Compound, p.fd3Payment, p.fd3FV, p.fd3Result,
        })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", p.currentPage)
      fmt.Printf("%s - %s\n", m.DTF(), errString)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      fmt.Println("*** Request timeout ***")
      if strings.EqualFold(p.currentPage, "rhs-ui2") {
        p.fd2Result = ""
      } else if strings.EqualFold(p.currentPage, "rhs-ui3") {
        p.fd3Result = ""
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
