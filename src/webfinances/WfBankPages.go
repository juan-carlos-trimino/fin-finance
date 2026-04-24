package webfinances

import (
  //bank "finance/databases/banking" //Importing a package and assigning it a local alias.
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/go-middlewares"
  // "github.com/juan-carlos-trimino/gpsessions"
  //Package template (html/template) implements data-driven templates for generating HTML output
  //safe against code injection. It provides the same interface as text/template and should be used
  //instead of text/template whenever the output is HTML.
  // "html/template"
  "net/http"
  "time"
)

// var bankTmpl *template.Template

/***
In Go, the predefined init() function sets off a piece of code to run before any other part of the
package; i.e., adding the init() function tells the compiler that when the package is imported, it
should run the init() function once. Unlike the main() function that can only be declared once, the
init() function can be declared multiple times throughout a package.
***/
// func init() {
//   logger.LogInfo("Entering init/webfinances.", "-1")
//   /***
//   The Must function wraps around the ParseGlob function that returns a pointer to a template and an
//   error, and it panics if the error is not nil.
//   ***/
//   bankTmpl = template.Must(template.ParseGlob("webfinances/templates/banking/*.html"))
// }

type WfBankingPages struct{}

func (p WfBankingPages) BankingPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.",
    startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering BankingPage/webfinances.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    tmpl.ExecuteTemplate(res, "banking_page", struct {
      Header string
      Datetime string
    } { "Banking", logger.DatetimeFormat() })
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms\n", time.Since(startTime).Microseconds()),
    correlationId)
}

// func (p WfPages) SimpleInterestPage(res http.ResponseWriter, req *http.Request) {
//   ctxKey := middlewares.MwContextKey{}
//   correlationId, _ := ctxKey.GetCorrelationId(req.Context())
//   startTime, _ := ctxKey.GetStartTime(req.Context())
//   logger.LogInfo(fmt.Sprintf("Created correlationId at %s.",
//     startTime.UTC().Format(time.RFC3339Nano)), correlationId)
//   logger.LogInfo("Entering SimpleInterestPage/webfinances.", correlationId)
//   sessionToken, _ := ctxKey.GetSessionToken(req.Context())
//   if sessionToken == "" {
//     invalidSession(res)
//   } else {
//     tmpl.ExecuteTemplate(res, "simple_interest_page", struct {
//       Header string
//       Datetime string
//     } { "Simple Interest", logger.DatetimeFormat() })
//   }
//   logger.LogInfo(fmt.Sprintf("Request took %vms\n", time.Since(startTime).Microseconds()),
//     correlationId)
// }

// func (p WfPages) OrdinaryAnnuityPage(res http.ResponseWriter, req *http.Request) {
//   ctxKey := middlewares.MwContextKey{}
//   correlationId, _ := ctxKey.GetCorrelationId(req.Context())
//   startTime, _ := ctxKey.GetStartTime(req.Context())
//   logger.LogInfo(fmt.Sprintf("Created correlationId at %s.",
//     startTime.UTC().Format(time.RFC3339Nano)), correlationId)
//   logger.LogInfo("Entering OrdinaryAnnuityPage/webfinances.", correlationId)
//   sessionToken, _ := ctxKey.GetSessionToken(req.Context())
//   if sessionToken == "" {
//     invalidSession(res)
//   } else {
//     tmpl.ExecuteTemplate(res, "ordinary_annuity_page", struct {
//       Header string
//       Datetime string
//     } { "Ordinary Annuity", logger.DatetimeFormat() })
//   }
//   logger.LogInfo(fmt.Sprintf("Request took %vms\n", time.Since(startTime).Microseconds()),
//     correlationId)
// }

// func (p WfPages) AnnuityDuePage(res http.ResponseWriter, req *http.Request) {
//   ctxKey := middlewares.MwContextKey{}
//   correlationId, _ := ctxKey.GetCorrelationId(req.Context())
//   startTime, _ := ctxKey.GetStartTime(req.Context())
//   logger.LogInfo(fmt.Sprintf("Created correlationId at %s.",
//     startTime.UTC().Format(time.RFC3339Nano)), correlationId)
//   logger.LogInfo("Entering AnnuityDuePage/webfinances.", correlationId)
//   sessionToken, _ := ctxKey.GetSessionToken(req.Context())
//   if sessionToken == "" {
//     invalidSession(res)
//   } else {
//     tmpl.ExecuteTemplate(res, "annuity_due_page", struct {
//       Header string
//       Datetime string
//     } { "Annuity Due", logger.DatetimeFormat() })
//   }
//   logger.LogInfo(fmt.Sprintf("Request took %vms\n", time.Since(startTime).Microseconds()),
//     correlationId)
// }
