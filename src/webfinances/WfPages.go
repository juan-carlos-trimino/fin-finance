package webfinances

import (
//   "finance/finances"
  "finance/misc"
  "fmt"
  "html/template"
  "net/http"
//   "strconv"
//   "strings"
)

var m = misc.Misc{}
var tmpl *template.Template

func init() {
  fmt.Printf("%s - Entering init/webfinances.\n", m.DTF())
  tmpl = template.Must(template.ParseGlob("webfinances/templates/*.html"))
}

type WfPages struct{}

func (p WfPages) HomePage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering HomePage/webfinances.\n", m.DTF())
  tmpl.ExecuteTemplate(res, "index.html", struct {
    Header string
    Datetime string
  } { "Investments", m.DTF() })
}

func (p WfPages) ContactPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering ContactPage/webfinances.\n", m.DTF())
  tmpl.ExecuteTemplate(res, "contact.html", struct {
    Header string
    Datetime string
  } { "Investments", m.DTF() })
}

func (p WfPages) AboutPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering AboutPage/webfinances.\n", m.DTF())
  tmpl.ExecuteTemplate(res, "about.html", struct {
    Header string
    Datetime string
  } { "Investments", m.DTF() })
}

func (p WfPages) PublicHomeFile(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering PublicHomeFile/webfinances.\n", m.DTF())
  http.ServeFile(res, req, "./webfinances/public/css/home.css")
}

func (p WfPages) FinancesPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering FinancesPage/webfinances.\n", m.DTF())
  tmpl.ExecuteTemplate(res, "finances.html", struct {
    Header string
    Datetime string
  } { "Finances", m.DTF() })
}

func (p WfPages) SimpleInterestPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering SimpleInterestPage/webfinances.\n", m.DTF())
  tmpl.ExecuteTemplate(res, "simpleinterest.html", struct {
    Header string
    Datetime string
  } { "Simple Interest", m.DTF() })
}
