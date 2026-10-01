package wfadmin

import (
  "finance/renderer"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/go-logger"
  "net/http"
)

type WfAdminPages struct{}

func (s WfAdminPages) WelcomePage(res http.ResponseWriter, req *http.Request) {
  ck := middlewares.MwContextKey{}
  correlationId, _ := ck.GetCorrelationId(req.Context())
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
}

func (s WfAdminPages) AdminSettingsPage(res http.ResponseWriter, req *http.Request) {
  ck := middlewares.MwContextKey{}
  correlationId, _ := ck.GetCorrelationId(req.Context())
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
}
