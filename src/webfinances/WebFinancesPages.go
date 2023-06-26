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
    Header: "Finances",
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
    Header: "Finances",
    Datetime: m.DTF(),
  }
  tmpl.ExecuteTemplate(res, "contact.html", params)
}

func (p *WebFinancesPages) AboutPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Calling AboutPage/webfinances.\n", m.DTF())
  params := WebFinancesPages {
    Header: "Finances",
    Datetime: m.DTF(),
  }
  tmpl.ExecuteTemplate(res, "about.html", params)
}
