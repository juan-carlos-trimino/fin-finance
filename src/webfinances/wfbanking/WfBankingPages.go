package wfbanking

import (
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/go-middlewares"
  // "github.com/juan-carlos-trimino/gpsessions"
  //Package template (html/template) implements data-driven templates for generating HTML output safe against code injection. It
  //provides the same interface as text/template and should be used instead of text/template whenever the output is HTML.
  "html/template"
  "net/http"
  "time"
)

/***
When handling authentication errors, the application should not disclose which part of the
authentication data was incorrect. Instead of "Invalid username" or "Invalid password", just use
"Invalid username and/or password" interchangeably.
***/
func invalidSession(res http.ResponseWriter) {
  err := tmpl.ExecuteTemplate(res, "login_page", struct {
    ErrMsg string
  } { "Invalid username and/or password" })
  //
  if err != nil {
    logger.LogError(fmt.Sprintf("%+v", err), "-1")
  }
}

var tmpl *template.Template

/***
In Go, the predefined init() function sets off a piece of code to run before any other part of the package; i.e., adding the
init() function tells the compiler that when the package is imported, it should run the init() function once. Unlike the main()
function that can only be declared once, the init() function can be declared multiple times throughout a package.
***/
func init() {
  logger.LogInfo("Entering wfbanking.init.", "-1")
  /***
  The Must function wraps around the ParseGlob function that returns a pointer to a template and an error, and it panics
  if the error is not nil.
  ***/
  tmpl = template.New("root")  //Initialize the root template.
  tmpl = template.Must(tmpl.ParseGlob("webfinances/templates/*.html"))
  tmpl = template.Must(tmpl.ParseGlob("webfinances/templates/banking/*.html"))
}

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
    err := tmpl.ExecuteTemplate(res, "banking_page", struct {
      Header string
      Datetime string
    } { "Banking", logger.DatetimeFormat() })
    //
    if err != nil {
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    }
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms\n", time.Since(startTime).Microseconds()), correlationId)
}














func (p WfBankingPages) PublicTableSelectRowFile(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfbanking.PublicTableSelectRowFile.", correlationId)
  http.ServeFile(res, req, "./webfinances/public/js/tableSelectRow.js")
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}
