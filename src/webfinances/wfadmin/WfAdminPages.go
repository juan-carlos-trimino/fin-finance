package wfadmin

import (
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/gpsessions"
  "net/http"
  "time"
)

/***
When handling authentication errors, the application should not disclose which part of the authentication data was incorrect.
Instead of "Invalid username" or "Invalid password", just use "Invalid username and/or password" interchangeably.
***/
func invalidSession(res http.ResponseWriter, correlationId string) {
  logger.LogInfo("Invalid session (wfadmin.invalidSession).", correlationId)
  templatesNeeded := []string{
    "webfinances/templates/layout.html",
    "webfinances/templates/login.html",
  }
  renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
    Data: struct{
      LayoutType string
      Header string
      ErrMsg string
    } { "std-wo-headers", "Login", "Invalid username and/or password" },
  })
}

type WfAdminPages struct{}

func (s WfAdminPages) WelcomePage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfadmin.WelcomePage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res, correlationId)
  } else {
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    cookie := sessions.CreateCookie(newSessionToken)
    http.SetCookie(res, cookie)
    templatesNeeded := []string{
      "webfinances/templates/layout.html",
      "webfinances/templates/admin/welcome.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/footer.html",
    }
    renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
      Data: struct{
        LayoutType string
        Header string
        Datetime string
        CsrfToken string
      } { "std-wo-nav-menu", "Investments - Admin", logger.DatetimeFormat(), newSession.CsrfToken },
    })
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (s WfAdminPages) AdminSettingsPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfadmin.AdminSettingsPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res, correlationId)
    return
  }
  newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
  cookie := sessions.CreateCookie(newSessionToken)
  http.SetCookie(res, cookie)
  templatesNeeded := []string{
    "webfinances/templates/layout.html",
    "webfinances/templates/admin/settings/settings.html",
    "webfinances/templates/title.html",
    "webfinances/templates/datetime.html",
    "webfinances/templates/footer.html",
  }
  renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
    Data: struct{
      LayoutType string
      Header string
      Datetime string
      CsrfToken string
    } { "std-wo-nav-menu", "Settings - Admin", logger.DatetimeFormat(), newSession.CsrfToken },
  })
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}
