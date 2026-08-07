package wfbanking

import (
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/go-middlewares"
  // "github.com/juan-carlos-trimino/gpsessions"
  //Package template (html/template) implements data-driven templates for generating HTML output safe against code injection. It
  //provides the same interface as text/template and should be used instead of text/template whenever the output is HTML.
  "net/http"
  "time"
)

var bankingMenuPage string = "home"
var contactMenuPage = "contact"
var aboutMenupage string = "about"

/***
When handling authentication errors, the application should not disclose which part of the
authentication data was incorrect. Instead of "Invalid username" or "Invalid password", just use
"Invalid username and/or password" interchangeably.
***/
func invalidSession(res http.ResponseWriter) {
  templatesNeeded := []string{
    "webfinances/templates/layout-simple.html",
    "webfinances/templates/login.html",
  }
  renderer.Render(res, "layout-simple", templatesNeeded, renderer.PageData{
    Data: struct{
      Header string
      ErrMsg string
    } { "Login", "Invalid username and/or password" },
  })
}

/***
In Go, the predefined init() function sets off a piece of code to run before any other part of the package; i.e., adding the
init() function tells the compiler that when the package is imported, it should run the init() function once. Unlike the main()
function that can only be declared once, the init() function can be declared multiple times throughout a package.
***/
// func init() {
// }

type WfBankingPages struct{}

func (p WfBankingPages) BankingPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfbanking.BankingPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    templatesNeeded := []string{
      "webfinances/templates/layout.html",
      "webfinances/templates/banking/banking.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/navbar.html",
      "webfinances/templates/footer.html",
    }
    renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
      Data: struct {
        Header string
        Datetime string
        MenuPage string
        } { "Bankig", logger.DatetimeFormat(), bankingMenuPage },
    })
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms\n", time.Since(startTime).Microseconds()), correlationId)
}
