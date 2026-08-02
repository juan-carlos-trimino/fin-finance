package webfinances

import (
  bank "finance/databases/banking"  //Importing a package and assigning it a local alias.
  banking "finance/webfinances/wfbanking"  //Importing a package and assigning it a local alias.
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/gpsessions"
  //Package template (html/template) implements data-driven templates for generating HTML output
  //safe against code injection. It provides the same interface as text/template and should be used
  //instead of text/template whenever the output is HTML.
  "html/template"
  "net/http"
  "time"
)

var tmpl *template.Template
var tsia1 *template.Template
var tsia2 *template.Template
var tsia3 *template.Template
var tsia4 *template.Template

/***
In Go, the predefined init() function sets off a piece of code to run before any other part of the
package; i.e., adding the init() function tells the compiler that when the package is imported, it
should run the init() function once. Unlike the main() function that can only be declared once, the
init() function can be declared multiple times throughout a package.
***/
func init() {
  logger.LogInfo("Entering init/webfinances.", "-1")
  /***
  The Must function wraps around the ParseGlob function that returns a pointer to a template and an
  error, and it panics if the error is not nil.
  ***/
  tmpl = template.New("root")  //Initialize the root template.
  tmpl = template.Must(tmpl.ParseGlob("webfinances/templates/*.html"))
  tmpl = template.Must(tmpl.ParseGlob("webfinances/templates/finances/*.html"))
  tsia1 = template.Must(template.New("sia1").ParseFiles(
    "webfinances/templates/finances/simpleinterestaccurate/accurate.html",
    "webfinances/templates/title.html",
    "webfinances/templates/datetime.html",
    "webfinances/templates/navbar.html",
    "webfinances/templates/finances/simpleinterestaccurate/amountofinterest.html",
    "webfinances/templates/footer.html"))
  tsia2 = template.Must(template.New("sia2").ParseFiles(
    "webfinances/templates/finances/simpleinterestaccurate/accurate.html",
    "webfinances/templates/title.html",
    "webfinances/templates/datetime.html",
    "webfinances/templates/navbar.html",
    "webfinances/templates/finances/simpleinterestaccurate/interestrate.html",
    "webfinances/templates/footer.html"))
  tsia3 = template.Must(template.New("sia3").ParseFiles(
    "webfinances/templates/finances/simpleinterestaccurate/accurate.html",
    "webfinances/templates/title.html",
    "webfinances/templates/datetime.html",
    "webfinances/templates/navbar.html",
    "webfinances/templates/finances/simpleinterestaccurate/principal.html",
    "webfinances/templates/footer.html"))
  tsia4 = template.Must(template.New("sia4").ParseFiles(
    "webfinances/templates/finances/simpleinterestaccurate/accurate.html",
    "webfinances/templates/title.html",
    "webfinances/templates/datetime.html",
    "webfinances/templates/navbar.html",
    "webfinances/templates/finances/simpleinterestaccurate/time.html",
    "webfinances/templates/footer.html"))
}

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

type WfPages struct{}

func (p WfPages) IndexPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.IndexPage.", correlationId)
  tmpl.ExecuteTemplate(res, "index_page", nil)
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (p WfPages) LoginPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering LoginPage.", correlationId)
  tmpl.ExecuteTemplate(res, "login_page", nil)
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (p WfPages) VerifyLogin(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering VerifyLogin.", correlationId)
  //Only allow POST requests.
  if req.Method != http.MethodPost {
    logger.LogInfo("Method not allowed.", correlationId)
    http.Error(res, "Method not allowed.", http.StatusMethodNotAllowed)
    return
  }
  un := req.PostFormValue("uname")
  pw := req.PostFormValue("pwd")
  ok, isAdmin := bank.DbAuthenticateUser(req.Context(), un, pw, correlationId)
  if !ok {
    invalidSession(res)
  } else {
    sessionToken, session := sessions.AddEntryToSessions(un)
    AddSessionDataPerUser(un, correlationId)
    banking.AddSessionDataPerUser(un, correlationId)
    /***
    Once a cookie is set on a client, it is sent along with every subsequent request. Cookies store
    historical information (including user login information) on the client's computer. The
    client's browser sends these cookies everytime the user visits the same website, automatically
    completing the login step for the user.

    Sessions, on the other hand, store historical information on the server side. The server uses a
    session id to identify different sessions, and the session id that is generated by the server
    should always be random and unique. You can use cookies or URL arguments to get the client's
    identity.
    ***/
    http.SetCookie(res, &http.Cookie{
      Name: "session_token",
      Value: sessionToken,
      Expires: session.Expiry,
    })
    tokenString, err := middlewares.GenerateJwtToken(isAdmin)
    if err != nil {
      logger.LogError(fmt.Sprintf("Error signing token: %v+", err), correlationId)
      http.Error(res, "Error signing token: ", http.StatusInternalServerError)
      return
    }
    http.SetCookie(res, &http.Cookie{
      Name: "admin_token",
      Value: tokenString,
    })
    if isAdmin {
      http.Redirect(res, req, "/admin/welcome", http.StatusSeeOther)
    } else {
      http.Redirect(res, req, "/welcome", http.StatusSeeOther)
    }
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (p WfPages) LogoutPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.LogoutPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    DeleteSessionDataPerUser(sessions.GetUserName(sessionToken))
    cookie := sessions.DeleteSession(sessionToken)
    http.SetCookie(res, cookie)
    http.Redirect(res, req, "/", http.StatusSeeOther)
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (p WfPages) WelcomePage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.WelcomePage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    templatesNeeded := []string{
      "webfinances/templates/layout.html",
      "webfinances/templates/welcome.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/navbar.html",
      "webfinances/templates/footer.html",
    }
    renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
      Data: struct{
        Header string
        Datetime string
        CurrentPage string
      } { "Investments", logger.DatetimeFormat(), "welcome" },
    })
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (p WfPages) ContactPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.ContactPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    templatesNeeded := []string{
      "webfinances/templates/layout.html",
      "webfinances/templates/contact.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/navbar.html",
      "webfinances/templates/footer.html",
    }
    renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
      Data: struct{
        Header string
        Datetime string
        CurrentPage string
      } { "Contact Us", logger.DatetimeFormat(), "contact" },
    })
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (p WfPages) AboutPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.AboutPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    //Explicitly map out the entire structural block stack for this page.
    templatesNeeded := []string{
      "webfinances/templates/layout.html",
      "webfinances/templates/about.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/navbar.html",
      "webfinances/templates/footer.html",
    }
    renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
      Data: struct{
        Header string
        Datetime string
        CurrentPage string
      } { "About Us", logger.DatetimeFormat(), "about" },
    })
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (p WfPages) FinancesPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.FinancesPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    templatesNeeded := []string{
      "webfinances/templates/layout.html",
      "webfinances/templates/finances/finances.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/navbar.html",
      "webfinances/templates/footer.html",
    }
    renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
      Data: struct{
        Header string
        Datetime string
        CurrentPage string
      } { "Finances", logger.DatetimeFormat(), "welcome" },
    })
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (p WfPages) SimpleInterestPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.SimpleInterestPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    templatesNeeded := []string{
      "webfinances/templates/layout.html",
      "webfinances/templates/finances/simpleinterest.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/navbar.html",
      "webfinances/templates/footer.html",
    }
    renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
      Data: struct{
        Header string
        Datetime string
        CurrentPage string
      } { "Simple Interest", logger.DatetimeFormat(), "welcome" },
    })
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (p WfPages) OrdinaryAnnuityPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.OrdinaryAnnuityPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    templatesNeeded := []string{
      "webfinances/templates/layout.html",
      "webfinances/templates/finances/ordinaryannuity.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/navbar.html",
      "webfinances/templates/footer.html",
    }
    renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
      Data: struct{
        Header string
        Datetime string
        CurrentPage string
      } { "Ordinary Annuity", logger.DatetimeFormat(), "welcome" },
    })
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (p WfPages) AnnuityDuePage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.AnnuityDuePage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    templatesNeeded := []string{
      "webfinances/templates/layout.html",
      "webfinances/templates/finances/annuitydue.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/navbar.html",
      "webfinances/templates/footer.html",
    }
    renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
      Data: struct{
        Header string
        Datetime string
        CurrentPage string
      } { "Annuity Due", logger.DatetimeFormat(), "welcome" },
    })
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (p WfPages) PublicHomeFile(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.PublicHomeFile.", correlationId)
  http.ServeFile(res, req, "./webfinances/public/css/home.css")
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (p WfPages) PublicSetPageUIFile(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.PublicSetPageUIFile.", correlationId)
  http.ServeFile(res, req, "./webfinances/public/js/setPageUI.js")
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}
