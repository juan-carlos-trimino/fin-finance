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

type WfOaInterestRatePages interface {
  OaInterestRatePages(http.ResponseWriter, *http.Request)
}

type wfOaInterestRatePages struct {
  currentPage string
  currentButton string
  //
  fd1N string
  fd1TimePeriod string
  fd1Compound string
  fd1PV string
  fd1FV string
  fd1Result string
}

func NewWfOaInterestRatePages() WfOaInterestRatePages {
  return &wfOaInterestRatePages {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1N: "1.0",
    fd1TimePeriod: "year",
    fd1Compound: "monthly",
    fd1PV: "1.00",
    fd1FV: "1.07",
    fd1Result: "",
  }
}

func (p *wfOaInterestRatePages) OaInterestRatePages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  sessionStatus, _ := ctxKey.GetSessionStatus(req.Context())
  if !sessionStatus {
    invalidSession(res)
    return
  }
  ctxKey = middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering OaInterestRatePages/webfinances.",
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
        p.fd1Compound = req.PostFormValue("fd1-cp")
        p.fd1PV = req.PostFormValue("fd1-pv")
        p.fd1FV = req.PostFormValue("fd1-fv")
        var n float64
        var pv float64
        var fv float64
        var err error
        if n, err = strconv.ParseFloat(p.fd1N, 64); err != nil {
          p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1N, err)
        } else if pv, err = strconv.ParseFloat(p.fd1PV, 64); err != nil {
          p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1PV, err)
        } else if fv, err = strconv.ParseFloat(p.fd1FV, 64); err != nil {
          p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1FV, err)
        } else {
          var oa finances.Annuities
          p.fd1Result = fmt.Sprintf("Interest: %.3f%% %s", oa.O_Interest_PV_FV(pv, fv, n,
                                    oa.GetTimePeriod(p.fd1TimePeriod[0], true),
                                    oa.GetCompoundingPeriod(p.fd1Compound[0], true)) * 100.0,
                                    p.fd1Compound)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, cp = %s, pv = %s, fv = %s, %s",
                      p.fd1N, p.fd1TimePeriod, p.fd1Compound, p.fd1PV, p.fd1FV, p.fd1Result),
        })
      }
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/interestrate/interestrate.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/ordinaryannuity/interestrate/n-PV-FV.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "oainterestrate", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd1N string
        Fd1TimePeriod string
        Fd1Compound string
        Fd1PV string
        Fd1FV string
        Fd1Result string
      } { "Ordinary Annuity / Interest Rate", m.DTF(), p.currentButton,
          p.fd1N, p.fd1TimePeriod, p.fd1Compound, p.fd1PV, p.fd1FV, p.fd1Result,
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
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
