package webfinances

import (
  "finance/finances"
  "finance/misc"
  "fmt"
  "html/template"
  "net/http"
  "strconv"
)

var m = misc.Misc{}
var tmpl *template.Template

func init() {
  fmt.Printf("%s - Entering init/webfinances.\n", m.DTF())
  tmpl = template.Must(template.ParseGlob("webfinances/templates/*.html"))
}

type WebFinancesPages struct {
  Header string
  Datetime string
}

/***
PS> curl.exe "http://localhost:8080"
***/
func (p *WebFinancesPages) HomePage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering HomePage/webfinances.\n", m.DTF())
  params := WebFinancesPages {
    Header: "Investments",
    Datetime: m.DTF(),
  }
  tmpl.ExecuteTemplate(res, "index.html", params)
}

func (p *WebFinancesPages) PublicHomeFile(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering PublicHomeFile/webfinances.\n", m.DTF())
  http.ServeFile(res, req, "./webfinances/public/css/home.css")
}

func (p *WebFinancesPages) ContactPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Calling ContactPage/webfinances.\n", m.DTF())
  params := WebFinancesPages {
    Header: "Investments",
    Datetime: m.DTF(),
  }
  tmpl.ExecuteTemplate(res, "contact.html", params)
}

func (p *WebFinancesPages) AboutPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Calling AboutPage/webfinances.\n", m.DTF())
  params := WebFinancesPages {
    Header: "Investments",
    Datetime: m.DTF(),
  }
  tmpl.ExecuteTemplate(res, "about.html", params)
}

func (p *WebFinancesPages) FinancesPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Calling FinancesPage/webfinances.\n", m.DTF())
  params := WebFinancesPages {
    Header: "Finances",
    Datetime: m.DTF(),
  }
  tmpl.ExecuteTemplate(res, "finances.html", params)
}

func (p *WebFinancesPages) SimpleInterestPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Calling SimpleInterestPage/webfinances.\n", m.DTF())
  params := WebFinancesPages {
    Header: "Simple Interest",
    Datetime: m.DTF(),
  }
  tmpl.ExecuteTemplate(res, "simpleinterest.html", params)
}

func (p *WebFinancesPages) SimpleInterestOrdinaryPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Calling SimpleInterestOrdinaryPage/webfinances.\n", m.DTF())
  if req.Method == http.MethodGet {
    tmpl.ExecuteTemplate(res, "simpleinterestordinary.html", struct {
      Result bool
      Header string
      Datetime string
      Time string
      TimePeriod string
      Interest string
      CompoundingPeriod string
      PresentValue string
      Answer string
    } { false, "Ordinary Interest", m.DTF(), "30", "day", "2.5", "monthly", "1000.00", "" })
  } else if req.Method == http.MethodPost {
    sn := req.FormValue("n")
    stp := req.FormValue("tp")
    si := req.FormValue("i")
    scp := req.FormValue("cp")
    spv := req.FormValue("pv")
    var n float64
    var i float64
    var pv float64
    var answer string
    var err error
    if n, err = strconv.ParseFloat(sn, 64); err != nil {
      answer = fmt.Sprintf("Time: %s -- %+v", sn, err)
    } else if i, err = strconv.ParseFloat(si, 64); err != nil {
      answer = fmt.Sprintf("Interest: %s -- %+v", sn, err)
    } else if pv, err = strconv.ParseFloat(spv, 64); err != nil {
      answer = fmt.Sprintf("Present Value: %s -- %+v", sn, err)
    } else {
      var si finances.SimpleInterest
      var periods finances.Periods
      answer = fmt.Sprintf("Amount of Interest: $%.5f", si.OrdinaryInterest(pv, i / 100.0,
                           periods.GetCompoundingPeriod(scp[0], false), n,
                           periods.GetTimePeriod(stp[0], false)))
    }
    fmt.Printf("%s - n = %s, tp = %s, i = %s, cp = %s, pv = %s\n", m.DTF(), sn, stp, si, scp, spv)
    tmpl.ExecuteTemplate(res, "simpleinterestordinary.html", struct {
      Result bool
      Header string
      Datetime string
      Time string
      TimePeriod string
      Interest string
      CompoundingPeriod string
      PresentValue string
      Answer string
    } { true, "Ordinary Interest", m.DTF(), sn, stp, si, scp, spv, answer })
  }
}
/*
func (p *WebFinancesPages) SimpleInterestOrdinaryCompute(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Calling SimpleInterestOrdinaryCompute/webfinances.\n", m.DTF())
  if req.Method != http.MethodPost {
    //tmpl.Execute(res, nil)
    http.Redirect(res, req, "/fin/simpleinterest/ordinary", http.StatusSeeOther)
    return
  }
  sn := req.FormValue("n")
  stp := req.FormValue("tp")
  // tp := req.FormValue("tp")[0]
  si := req.FormValue("i")
  scp := req.FormValue("cp")
  spv := req.FormValue("pv")
  var n float64
  var i float64
  var pv float64
  var answer string
  var err error
  if n, err = strconv.ParseFloat(sn, 64); err != nil {
    answer = fmt.Sprintf("Time: %s -- %+v", sn, err)
  } else if i, err = strconv.ParseFloat(si, 64); err != nil {
    answer = fmt.Sprintf("Interest: %s -- %+v", sn, err)
  } else if pv, err = strconv.ParseFloat(spv, 64); err != nil {
    answer = fmt.Sprintf("Present Value: %s -- %+v", sn, err)
  } else {
    var si finances.SimpleInterest
    var periods finances.Periods

    cp1 := periods.GetCompoundingPeriod(scp[0], false)
    tp1 := periods.GetTimePeriod(stp[0], false)
    fmt.Printf("%s - n = %.5f, tp = %d, i = %.5f, cp = %d, pv = %.5f\n", m.DTF(), n, tp1, i, cp1, pv)

    c := si.OrdinaryInterest(pv, i / 100.0, periods.GetCompoundingPeriod(scp[0], false), n, periods.GetTimePeriod(stp[0], false))
    answer = fmt.Sprintf("Amount of Interest: $%.5f", c)
  }
  fmt.Printf("%s - n = %s, tp = %s, i = %s, cp = %s, pv = %s\n", m.DTF(), sn, stp, si, scp, spv)
  tmpl.ExecuteTemplate(res, "simpleinterestordinary.html", struct {
    Result bool
    Header string
    Datetime string
    Time string
    TimePeriod string
    Interest string
    CompoundingPeriod string
    PresentValue string
    Answer string
  } { true, "Ordinary Interest", m.DTF(), sn, stp, si, scp, spv, answer })
}
*/