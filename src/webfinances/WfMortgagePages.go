package webfinances

import (
  "finance/middlewares"
  "finance/finances"
  "fmt"
  "html/template"
  "net/http"
  "strconv"
  "strings"
)

var mortgage_notes = [...]string {
  "Refinance mortgage and HELOC with one load.",
  "If the blended interest rate is higher than what you could get on a new fixed-rate mortgage, consider it.",
}

type WfMortgagePages interface {
  MortgagePages(http.ResponseWriter, *http.Request)
}

type wfMortgagePages struct {
  currentPage string
  currentButton string
  //
  fd1N string
  fd1TimePeriod string
  fd1Interest string
  fd1Compound string
  fd1Amount string
  fd1Result [3]string
  //
  fd2FaceValue string
  fd2Time string
  fd2TimePeriod string
  fd2Coupon string
  fd2Current string
  fd2Compound string
  fd2Result string
  //
  fd3Mrate string
  fd3Mbalance string
  fd3Hrate string
  fd3Hbalance string
  fd3Result [3]string
}

func NewWfMortgagePages() WfMortgagePages {
  return &wfMortgagePages {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1N: "30.0",
    fd1TimePeriod: "year",
    fd1Interest: "7.50",
    fd1Compound: "monthly",
    fd1Amount: "100000.00",
    fd1Result: [3]string { "", "", "" },
    //
    fd2FaceValue: "1000.00",
    fd2Time: "5",
    fd2TimePeriod: "year",
    fd2Coupon: "3.00",
    fd2Current: "3.5",
    fd2Compound: "semiannually",
    fd2Result: "",
    //
    fd3Mrate: "3.375",
    fd3Mbalance: "300000.00",
    fd3Hrate: "2.875",
    fd3Hbalance: "100000.00",
    fd3Result: [3]string { mortgage_notes[0], mortgage_notes[1], "" },
  }
}

func (p *wfMortgagePages) MortgagePages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  logEntry := LogEntry{}
  logEntry.Print(INFO, correlationId, []string {
    "Entering MortgagePages/webfinances.",
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
        p.fd1Compound = req.PostFormValue("fd1-compound")
        p.fd1Amount = req.PostFormValue("fd1-amount")
        var n float64
        var i float64
        var amount float64
        var err error
        p.fd1Result[1] = ""
        p.fd1Result[2] = ""
        if n, err = strconv.ParseFloat(p.fd1N, 64); err != nil {
          p.fd1Result[0] = fmt.Sprintf("Error: %s -- %+v", p.fd1N, err)
        } else if i, err = strconv.ParseFloat(p.fd1Interest, 64); err != nil {
          p.fd1Result[0] = fmt.Sprintf("Error: %s -- %+v", p.fd1Interest, err)
        } else if amount, err = strconv.ParseFloat(p.fd1Amount, 64); err != nil {
          p.fd1Result[0] = fmt.Sprintf("Error: %s -- %+v", p.fd1Amount, err)
        } else {
          var m finances.Mortgage
          payment, totalCost, totalInterest := m.CostOfMortgage(amount, i / 100.0,
                                               p.fd1Compound[0], n, p.fd1TimePeriod[0])
          p.fd1Result[0] = fmt.Sprintf("Payment: $%.2f", payment)
          p.fd1Result[1] = fmt.Sprintf("Total Interest: $%.2f", totalInterest)
          p.fd1Result[2] = fmt.Sprintf("Total Cost: $%.2f", totalCost)
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("n = %s, tp = %s, interest = %s, cp = %s, amount = %s, %s", p.fd1N,
                      p.fd1TimePeriod, p.fd1Interest, p.fd1Compound, p.fd1Amount, p.fd1Result[0]),
        })
      }
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template
      and an error, and it panics if the error is not nil.
      ***/
      t := template.Must(template.ParseFiles("webfinances/templates/mortgage/mortgage.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/mortgage/costofmortgage.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "mortgage", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd1N string
        Fd1TimePeriod string
        Fd1Interest string
        Fd1Compound string
        Fd1Amount string
        Fd1Result [3]string
      } { "Mortgage", m.DTF(), p.currentButton,
          p.fd1N, p.fd1TimePeriod, p.fd1Interest, p.fd1Compound, p.fd1Amount, p.fd1Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui2") {
      p.currentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        p.fd2FaceValue = req.FormValue("fd2-facevalue")
        p.fd2Time = req.PostFormValue("fd2-time")
        p.fd2TimePeriod = req.PostFormValue("fd2-tp")
        p.fd2Coupon = req.PostFormValue("fd2-coupon")
        p.fd2Current = req.PostFormValue("fd2-current")
        p.fd2Compound = req.PostFormValue("fd2-compound")
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
      }
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/bonds/currentprice.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd2FaceValue string
        Fd2Time string
        Fd2TimePeriod string
        Fd2Coupon string
        Fd2Current string
        Fd2Compound string
        Fd2Result string
      } { "Bonds", m.DTF(), p.currentButton,
          p.fd2FaceValue, p.fd2Time, p.fd2TimePeriod, p.fd2Coupon, p.fd2Current, p.fd2Compound, p.fd2Result,
        })
    } else if strings.EqualFold(p.currentPage, "rhs-ui3") {
      p.currentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        p.fd3Mrate = req.PostFormValue("fd3-mrate")
        p.fd3Mbalance = req.PostFormValue("fd3-mbalance")
        p.fd3Hrate = req.PostFormValue("fd3-hrate")
        p.fd3Hbalance = req.PostFormValue("fd3-hbalance")
        var mRate float64
        var mBalance float64
        var hRate float64
        var hBalance float64
        var err error
        if mRate, err = strconv.ParseFloat(p.fd3Mrate, 64); err != nil {
          p.fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", p.fd3Mrate, err)
        } else if mBalance, err = strconv.ParseFloat(p.fd3Mbalance, 64); err != nil {
          p.fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", p.fd3Mbalance, err)
        } else if hRate, err = strconv.ParseFloat(p.fd3Hrate, 64); err != nil {
          p.fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", p.fd3Hrate, err)
        } else if hBalance, err = strconv.ParseFloat(p.fd3Hbalance, 64); err != nil {
          p.fd3Result[2] = fmt.Sprintf("Error: %s -- %+v", p.fd3Hbalance, err)
        } else {
          var m finances.Mortgage
          p.fd3Result[2] = fmt.Sprintf("Blended Interest Rate: %.3f%%",
                                       m.BlendedInterestRate(mBalance, mRate, hBalance, hRate))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("mortgage balance = %s, mortgage rate = %s, HELOC balance = %s, HELOC rate = %s, %s\n",
                      p.fd3Mbalance, p.fd3Mrate, p.fd3Hbalance, p.fd3Hrate, p.fd3Result[2]),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/mortgage/mortgage.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/mortgage/heloc.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "mortgage", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd3Mrate string
        Fd3Mbalance string
        Fd3Hrate string
        Fd3Hbalance string
        Fd3Result [3]string
      } { "Bonds", m.DTF(), p.currentButton,
          p.fd3Mrate, p.fd3Mbalance, p.fd3Hrate, p.fd3Hbalance, p.fd3Result,
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
