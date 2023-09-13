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
  fd3FaceValue string
  fd3TimeCall string
  fd3TimePeriod string
  fd3Coupon string
  fd3Compound string
  fd3BondPrice string
  fd3CallPrice string
  fd3Result string
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
    fd3FaceValue: "1000.00",
    fd3TimeCall: "2",
    fd3TimePeriod: "year",
    fd3Coupon: "2.0",
    fd3Compound: "semiannually",
    fd3BondPrice: "990.00",
    fd3CallPrice: "1050.00",
    fd3Result: "",
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
        p.fd3FaceValue = req.PostFormValue("fd3-facevalue")
        p.fd3TimeCall = req.PostFormValue("fd3-timecall")
        p.fd3TimePeriod = req.PostFormValue("fd3-tp")
        p.fd3Coupon = req.PostFormValue("fd3-coupon")
        p.fd3BondPrice = req.PostFormValue("fd3-bondprice")
        p.fd3CallPrice = req.PostFormValue("fd3-callprice")
        p.fd3Compound = req.PostFormValue("fd3-compound")
        var fv float64
        var timeToCall float64
        var couponRate float64
        var bondPrice float64
        var callPrice float64
        var err error
        if fv, err = strconv.ParseFloat(p.fd3FaceValue, 64); err != nil {
          p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3FaceValue, err)
        } else if timeToCall, err = strconv.ParseFloat(p.fd3TimeCall, 64); err != nil {
          p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3TimeCall, err)
        } else if couponRate, err = strconv.ParseFloat(p.fd3Coupon, 64); err != nil {
          p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3Coupon, err)
        } else if bondPrice, err = strconv.ParseFloat(p.fd3BondPrice, 64); err != nil {
          p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3BondPrice, err)
        } else if callPrice, err = strconv.ParseFloat(p.fd3CallPrice, 64); err != nil {
          p.fd3Result = fmt.Sprintf("Error: %s -- %+v", p.fd3CallPrice, err)
        } else {
          var b finances.Bonds
          p.fd3Result = fmt.Sprintf("Yield to Call: %.3f%%", b.YieldToCall(fv, couponRate,
                                     b.GetCompoundingPeriod(p.fd3Compound[0], true), timeToCall,
                                     b.GetTimePeriod(p.fd3TimePeriod[0], true), bondPrice,
                                     callPrice))
        }
        logEntry.Print(INFO, correlationId, []string {
          fmt.Sprintf("fv = %s, coupon rate = %s, cp = %s, time to call = %s, tp = %s, bond price = %s, call price = %s, %s\n",
                       p.fd3FaceValue, p.fd3Coupon, p.fd3Compound, p.fd3TimeCall, p.fd3TimePeriod,
                       p.fd3BondPrice, p.fd3CallPrice, p.fd3Result),
        })
      }
      t := template.Must(template.ParseFiles("webfinances/templates/bonds/bonds.html",
                                             "webfinances/templates/header.html",
                                             "webfinances/templates/bonds/yieldtocall.html",
                                             "webfinances/templates/footer.html"))
      t.ExecuteTemplate(res, "bonds", struct {
        Header string
        Datetime string
        CurrentButton string
        Fd3FaceValue string
        Fd3TimeCall string
        Fd3TimePeriod string
        Fd3Coupon string
        Fd3Compound string
        Fd3BondPrice string
        Fd3CallPrice string
        Fd3Result string
      } { "Bonds", m.DTF(), p.currentButton,
          p.fd3FaceValue, p.fd3TimeCall, p.fd3TimePeriod, p.fd3Coupon, p.fd3Compound, p.fd3BondPrice, p.fd3CallPrice, p.fd3Result,
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
