package webfinances

import (
  "finance/misc"
  "fmt"
  "html/template"
  "net/http"
)

var m = misc.Misc{}
var tmpl *template.Template

func init() {
  fmt.Printf("%s - Entering init/webfinances.\n", m.DTF())
  /***
  The Must function wraps around the ParseGlob function that returns a pointer to a template and an
  error, and it panics if the error is not nil.
  ***/
  tmpl = template.Must(template.ParseGlob("webfinances/templates/*.html"))
}

type WfPages struct{}

func (p WfPages) HomePage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering HomePage/webfinances.\n", m.DTF())
  tmpl.ExecuteTemplate(res, "home_page", struct {
    Header string
    Datetime string
  } { "Investments", m.DTF() })
}

func (p WfPages) ContactPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering ContactPage/webfinances.\n", m.DTF())
  tmpl.ExecuteTemplate(res, "contact_page", struct {
    Header string
    Datetime string
  } { "Investments", m.DTF() })
}

func (p WfPages) AboutPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering AboutPage/webfinances.\n", m.DTF())
  /***
  Executing the template means that we take the content from the template files, combine it with
  data from another source, and generate the final HTML content.
  ***/
  tmpl.ExecuteTemplate(res, "about_page", struct {
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
  tmpl.ExecuteTemplate(res, "finances_page", struct {
    Header string
    Datetime string
  } { "Finances", m.DTF() })
}

func (p WfPages) SimpleInterestPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering SimpleInterestPage/webfinances.\n", m.DTF())
  tmpl.ExecuteTemplate(res, "simple_interest_page", struct {
    Header string
    Datetime string
  } { "Simple Interest", m.DTF() })
}
