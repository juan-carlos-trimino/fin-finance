package webfinances

import (
  //"encoding/json"
  "finance/middlewares"
  "finance/misc"
  "finance/sessions"
  "fmt"
  //Package template (html/template) implements data-driven templates for generating HTML output
  //safe against code injection. It provides the same interface as text/template and should be used
  //instead of text/template whenever the output is HTML.
  "html/template"
  "net/http"
  "time"
)

var m = misc.Misc{}
var tmpl *template.Template

/***
In Go, the predefined init() function sets off a piece of code to run before any other part of the
package; i.e., adding the init() function tells the compiler that when the package is imported, it
should run the init() function once. Unlike the main() function that can only be declared once, the
init() function can be declared multiple times throughout a package.
***/
func init() {
  fmt.Printf("%s - Entering init/webfinances.\n", m.DTF())
  /***
  The Must function wraps around the ParseGlob function that returns a pointer to a template and an
  error, and it panics if the error is not nil.
  ***/
  tmpl = template.Must(template.ParseGlob("webfinances/templates/*.html"))
}

/***
When handling authentication errors, the application should not disclose which part of the
authentication data was incorrect. Instead of "Invalid username" or "Invalid password", just use
"Invalid username and/or password" interchangeably.
***/
func invalidSession(res http.ResponseWriter) {
  tmpl.ExecuteTemplate(res, "index_page", struct {
    Error string
  } { "Invalid username and/or password" })
}

type WfPages struct{}

func (p WfPages) IndexPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering IndexPage/webfinances.\n", m.DTF())
  tmpl.ExecuteTemplate(res, "index_page", nil)
}

func (p WfPages) LoginPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering LoginPage/webfinances.\n", m.DTF())
  un := req.PostFormValue("username")
  pw := req.PostFormValue("password")
  //Get the expected password from the in memory map.
  if hashedPassword, ok := sessions.Users[un]; !ok {
    invalidSession(res)
  } else if ok, _ := sessions.CompareHashAndPassword(hashedPassword, []byte(pw)); !ok {
    invalidSession(res)
  } else {
    sessionToken := sessions.AddEntryToSessions(un)
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
      Expires: sessions.Sessions[sessionToken].Expiry,
    })
    http.Redirect(res, req, "/welcome", http.StatusSeeOther)
  }
}

func (p WfPages) LogoutPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering LogoutPage/webfinances.\n", m.DTF())
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    delete(sessions.Sessions, sessionToken)
    http.SetCookie(res, &http.Cookie{
      Name: "session_token",
      Value: "",
      Path: "/",
      Expires: time.Now(),
    })
    http.Redirect(res, req, "/", http.StatusSeeOther)
  }
}

func (p WfPages) WelcomePage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering WelcomePage/webfinances.\n", m.DTF())
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    tmpl.ExecuteTemplate(res, "welcome_page", struct {
      Header string
      Datetime string
    } { "Investments", m.DTF() })
  }
}

func (p WfPages) ContactPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering ContactPage/webfinances.\n", m.DTF())
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    tmpl.ExecuteTemplate(res, "contact_page", struct {
      Header string
      Datetime string
    } { "Investments", m.DTF() })
  }
}

func (p WfPages) AboutPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering AboutPage/webfinances.\n", m.DTF())
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    /***
    Executing the template means that we take the content from the template files, combine it with
    data from another source, and generate the final HTML content.
    ***/
    tmpl.ExecuteTemplate(res, "about_page", struct {
      Header string
      Datetime string
    } { "Investments", m.DTF() })
  }
}

func (p WfPages) PublicHomeFile(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering PublicHomeFile/webfinances.\n", m.DTF())
  http.ServeFile(res, req, "./webfinances/public/css/home.css")
}

func (p WfPages) FinancesPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering FinancesPage/webfinances.\n", m.DTF())
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    tmpl.ExecuteTemplate(res, "finances_page", struct {
      Header string
      Datetime string
    } { "Finances", m.DTF() })
  }
}

func (p WfPages) SimpleInterestPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering SimpleInterestPage/webfinances.\n", m.DTF())
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    tmpl.ExecuteTemplate(res, "simple_interest_page", struct {
      Header string
      Datetime string
    } { "Simple Interest", m.DTF() })
  }
}

func (p WfPages) OrdinaryAnnuityPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering OrdinaryAnnuityPage/webfinances.\n", m.DTF())
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    tmpl.ExecuteTemplate(res, "ordinary_annuity_page", struct {
      Header string
      Datetime string
    } { "Ordinary Annuity", m.DTF() })
  }
}

func (p WfPages) AnnuityDuePage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering AnnuityDuePage/webfinances.\n", m.DTF())
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    tmpl.ExecuteTemplate(res, "annuity_due_page", struct {
      Header string
      Datetime string
    } { "Annuity Due", m.DTF() })
  }
}


