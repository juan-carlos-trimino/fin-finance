package wfadmin

import (
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/go-logger"
  "net/http"
  "time"
)

type WfAdminPages struct{}

func (s WfAdminPages) WelcomePage(res http.ResponseWriter, req *http.Request) {
  ck := middlewares.MwContextKey{}
  correlationId, _ := ck.GetCorrelationId(req.Context())
  startTime, _ := ck.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfadmin.WelcomePage.", correlationId)
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
    }{
      "std-wo-nav-menu",
      "Investments - Admin",
      logger.DatetimeFormat(),
    },
  })
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (s WfAdminPages) AdminSettingsPage(res http.ResponseWriter, req *http.Request) {
  ck := middlewares.MwContextKey{}
  correlationId, _ := ck.GetCorrelationId(req.Context())
  startTime, _ := ck.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfadmin.AdminSettingsPage.", correlationId)
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
    }{
      "std-wo-nav-menu",
      "Settings - Admin",
      logger.DatetimeFormat(),
    },
  })
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}
