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

type WfOaPvPages interface {
  OaPvPages(http.ResponseWriter, *http.Request)
}

type wfOaPvPages struct {
  currentPage string
  currentButton string
  //
  fd1N string
  fd1TimePeriod string
  fd1Interest string
  fd1Compound string
  fd1FV string
  fd1Result string
  //
  fd2N string
  fd2TimePeriod string
  fd2Interest string
  fd2Compound string
  fd2PMT string
  fd2Result string
}

func NewWfOaPvPages() WfOaPvPages {
  return &wfOaPvPages {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1N: "1.0",
    fd1TimePeriod: "year",
    fd1Interest: "1.00",
    fd1Compound: "monthly",
    fd1FV: "1.00",
    fd1Result: "",
    //
    fd2N: "1.0",
    fd2TimePeriod: "year",
    fd2Interest: "1.00",
    fd2Compound: "monthly",
    fd2PMT: "1.00",
    fd2Result: "",
  }
}

func (p *wfOaPvPages) OaPvPages(res http.ResponseWriter, req *http.Request) {
  if !checkSession(res, req) {
    return
  }
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering OaPvPages/webfinances.",
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
        p.fd1N = req.PostFormValue("fd1-n")
        p.fd1TimePeriod = req.PostFormValue("fd1-tp")
        p.fd1Interest = req.PostFormValue("fd1-interest")
        p.fd1Compound = req.PostFormValue("fd1-cp")
        p.fd1FV = req.PostFormValue("fd1-fv")
        var n float64
        var i float64
        var fv float64
        var err error
        if n, err = strconv.ParseFloat(p.fd1N, 64); err != nil {
          p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1N, err)
        } else if i, err = strconv.ParseFloat(p.fd1Interest, 64); err != nil {
          p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1Interest, err)
        } else if fv, err = strconv.ParseFloat(p.fd1FV, 64); err != nil {
          p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1FV, err)
        } else {
          var oa finances.Annuities
          p.fd1Result = fmt.Sprintf("Present Value: $%.2f", oa.O_PresentValue_FV(fv, i / 100.0,
                                    oa.GetCompoundingPeriod(p.fd1Compound[0], true),
                                    n, oa.GetTimePeriod(p.fd1TimePeriod[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, fv = %s, %s",
                      p.fd1N, p.fd1TimePeriod, p.fd1Interest, p.fd1Compound, p.fd1FV, p.fd1Result),
        })
      }
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/pv/pv.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/ordinaryannuity/pv/n-i-FV.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "oapresentvalue", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd1N string
        Fd1TimePeriod string
        Fd1Interest string
        Fd1Compound string
        Fd1FV string
        Fd1Result string
      } { "Ordinary Annuity / Present Value", m.DTF(), p.currentButton,
          p.fd1N, p.fd1TimePeriod, p.fd1Interest, p.fd1Compound, p.fd1FV, p.fd1Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui2") {
      p.currentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        p.fd2N = req.FormValue("fd2-n")
        p.fd2TimePeriod = req.PostFormValue("fd2-tp")
        p.fd2Interest = req.PostFormValue("fd2-interest")
        p.fd2Compound = req.PostFormValue("fd2-cp")
        p.fd2PMT = req.PostFormValue("fd2-pmt")
        var n float64
        var i float64
        var pmt float64
        var err error
        if n, err = strconv.ParseFloat(p.fd2N, 64); err != nil {
          p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2N, err)
        } else if i, err = strconv.ParseFloat(p.fd2Interest, 64); err != nil {
          p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2Interest, err)
        } else if pmt, err = strconv.ParseFloat(p.fd2PMT, 64); err != nil {
          p.fd2Result = fmt.Sprintf("Error: %s -- %+v", p.fd2PMT, err)
        } else {
          var oa finances.Annuities
          p.fd2Result = fmt.Sprintf("Present Value: $%.2f", oa.O_PresentValue_PMT(pmt, i / 100.0,
                                    oa.GetCompoundingPeriod(p.fd2Compound[0], true),
                                    n, oa.GetTimePeriod(p.fd2TimePeriod[0], true)))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, interest = %s, cp = %s, pmt = %s, %s",
                      p.fd2N, p.fd2TimePeriod, p.fd2Interest, p.fd2Compound, p.fd2PMT,
                      p.fd2Result),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/pv/pv.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/ordinaryannuity/pv/n-i-PMT.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "oapresentvalue", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd2N string
        Fd2TimePeriod string
        Fd2Interest string
        Fd2Compound string
        Fd2PMT string
        Fd2Result string
      } { "Ordinary Annuity / Present Value", m.DTF(), p.currentButton,
          p.fd2N, p.fd2TimePeriod, p.fd2Interest, p.fd2Compound, p.fd2PMT, p.fd2Result,
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
        p.fd1Result = ""
      } else if strings.EqualFold(p.currentPage, "rhs-ui2") {
        p.fd2Result = ""
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
