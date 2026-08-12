package wfadmin

import (
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/gplogger"
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

func displayWelcomePage(res http.ResponseWriter) {
  templatesNeeded := []string{
    "webfinances/templates/layout-no-navbar.html",
    "webfinances/templates/admin/welcome.html",
    "webfinances/templates/title.html",
    "webfinances/templates/datetime.html",
    "webfinances/templates/footer.html",
  }
  renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
    Data: struct{
      Header string
      Datetime string
    } { "Investments - Admin", logger.DatetimeFormat() },
  })
}

type WfAdminPages struct {}

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
    displayWelcomePage(res)
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

/*
func (s WfAdminPages) UsersPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfadmin.UsersPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    cookie := sessions.CreateCookie(newSessionToken)
    http.SetCookie(res, cookie)
    templatesNeeded := []string{
      "webfinances/templates/layout-no-navbar.html",
      "webfinances/templates/admin/users/users.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/footer.html",
    }
    renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
      Data: struct{
        Header string
        Datetime string
        CsrfToken string
      } { "Users - Admin", logger.DatetimeFormat(), newSession.CsrfToken, },
    })
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}
*/





















func (s WfAdminPages) AdminSettingsPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering AdminSettingsPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res, correlationId)
  } else {
    templatesNeeded := []string{
      "webfinances/templates/layout-no-navbar.html",
      "webfinances/templates/admin/settings/settings.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/footer.html",
    }
    renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
      Data: struct{
        Header string
        Datetime string
      } { "Investments - Admin", logger.DatetimeFormat() },
    })
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}
