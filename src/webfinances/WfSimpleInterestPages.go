package webfinances

import (
  "finance/finances"
  "finance/misc"
  "fmt"
  "html/template"
  "net/http"
  "strconv"
  "strings"
)

var m = misc.Misc{}
var tmpl *template.Template

func init() {
  fmt.Printf("%s - Entering init/webfinances.\n", m.DTF())
  tmpl = template.Must(template.ParseGlob("webfinances/templates/*.html"))
}

type WfSimpleInterestPages interface {
  SimpleInterestPage(http.ResponseWriter, *http.Request)
  SimpleInterestOrdinaryPage(http.ResponseWriter, *http.Request)
  HomePage(http.ResponseWriter, *http.Request)
  PublicHomeFile(http.ResponseWriter, *http.Request)
  ContactPage(http.ResponseWriter, *http.Request)
  AboutPage(http.ResponseWriter, *http.Request)
  FinancesPage(http.ResponseWriter, *http.Request)
}

type wfSimpleInterestPages struct {
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
  fd2PV string
  fd2Result string
}

func NewWfSimpleInterestPages() WfSimpleInterestPages {
  return &wfSimpleInterestPages {
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
    fd2PV: "1.00",
    fd2Result: "",
  }
}

/***
PS> curl.exe "http://localhost:8080"
***/
func (p *wfSimpleInterestPages) HomePage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering HomePage/webfinances.\n", m.DTF())
  tmpl.ExecuteTemplate(res, "index.html", struct {
    Header string
    Datetime string
  } { "Investments", m.DTF() })
}

func (p *wfSimpleInterestPages) PublicHomeFile(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering PublicHomeFile/webfinances.\n", m.DTF())
  http.ServeFile(res, req, "./webfinances/public/css/home.css")
}

func (p *wfSimpleInterestPages) ContactPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering ContactPage/webfinances.\n", m.DTF())
  tmpl.ExecuteTemplate(res, "contact.html", struct {
    Header string
    Datetime string
  } { "Investments", m.DTF() })
}

func (p *wfSimpleInterestPages) AboutPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering AboutPage/webfinances.\n", m.DTF())
  tmpl.ExecuteTemplate(res, "about.html", struct {
    Header string
    Datetime string
  } { "Investments", m.DTF() })
}

func (p *wfSimpleInterestPages) FinancesPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering FinancesPage/webfinances.\n", m.DTF())
  tmpl.ExecuteTemplate(res, "finances.html", struct {
    Header string
    Datetime string
  } { "Finances", m.DTF() })
}

func (p *wfSimpleInterestPages) SimpleInterestPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering SimpleInterestPage/webfinances.\n", m.DTF())
  tmpl.ExecuteTemplate(res, "simpleinterest.html", struct {
    Header string
    Datetime string
  } { "Simple Interest", m.DTF() })
}

func (p *wfSimpleInterestPages) SimpleInterestOrdinaryPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering SimpleInterestOrdinaryPage/webfinances.\n", m.DTF())
  if req.Method == http.MethodGet {
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
      Fd2PV string
      Fd2Result string
    } { "Ordinary Interest", m.DTF(), p.currentButton,
        p.fd1Time, p.fd1TimePeriod, p.fd1Interest, p.fd1Compound, p.fd1PV, p.fd1Result,
        p.fd2Time, p.fd2TimePeriod, p.fd2Amount, p.fd2PV, p.fd2Result })
  } else if req.Method == http.MethodPost {
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
      p.currentButton = "lhs-button2"

    } else {
      errString := fmt.Sprintf("Unsupported page: %s", ui)
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
      Fd2PV string
      Fd2Result string
    } { "Ordinary Interest", m.DTF(), p.currentButton,
        p.fd1Time, p.fd1TimePeriod, p.fd1Interest, p.fd1Compound, p.fd1PV, p.fd1Result,
        p.fd2Time, p.fd2TimePeriod, p.fd2Amount, p.fd2PV, p.fd2Result })
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    fmt.Printf("%s - %s\n", m.DTF(), errString)
    panic(errString)
  }
}
