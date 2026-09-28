package wfbanking

import (
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/go-logger"
  "github.com/juan-carlos-trimino/go-middlewares"
  //Package template (html/template) implements data-driven templates for generating HTML output safe against code injection. It
  //provides the same interface as text/template and should be used instead of text/template whenever the output is HTML.
  "net/http"
  "time"
)

var bankingMenuPage string = "home"
var contactMenuPage = "contact"
var aboutMenupage string = "about"

/***
In Go, the predefined init() function sets off a piece of code to run before any other part of the package; i.e., adding the
init() function tells the compiler that when the package is imported, it should run the init() function once. Unlike the main()
function that can only be declared once, the init() function can be declared multiple times throughout a package.
***/
// func init() {
// }

type WfBankingPages struct{}

func (p WfBankingPages) BankingPage(res http.ResponseWriter, req *http.Request) {
  ck := middlewares.MwContextKey{}
  correlationId, _ := ck.GetCorrelationId(req.Context())
  startTime, _ := ck.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfbanking.BankingPage.", correlationId)
  //Declare dynamic variables based on the URL path.
  var bodyTemplate string
  var pageHeader string
  //Map the incoming route path to its respective template and title.
  switch req.URL.Path {
  case "/banking":
    bodyTemplate = "banking.html"
    pageHeader = "Banking"
  default:
    http.NotFound(res, req)
    return
  }
  //Construct the templates slice dynamically using the variables.
  templatesNeeded := []string{
    "webfinances/templates/layout.html",
    "webfinances/templates/banking/" + bodyTemplate,  //Dynamically loaded.,
    "webfinances/templates/title.html",
    "webfinances/templates/datetime.html",
    "webfinances/templates/navbar.html",
    "webfinances/templates/footer.html",
  }
  renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
    Data: struct {
      LayoutType string
      Header string
      Datetime string
      MenuPage string
    }{
      "standard",
      pageHeader,
      logger.DatetimeFormat(),
      bankingMenuPage,
    },
  })
  logger.LogInfo(fmt.Sprintf("Request took %vms\n", time.Since(startTime).Microseconds()), correlationId)
}
