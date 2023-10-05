package webfinances

import (
  //"encoding/json"
  "finance/misc"
  "finance/sessions"
  "fmt"
  //The option -u instructs 'get' to update the module with dependencies.
  //go get -u github.com/google/uuid
  "github.com/google/uuid"
//  "golang.org/x/crypto/bcrypt"
  "html/template"
  "net/http"
  "time"
)

var m = misc.Misc{}
var tmpl *template.Template

func init() {
  fmt.Printf("%s - Entering init/webfinances.\n", m.DTF())
  /***
  The Must function wraps around the ParseGlob function that returns a pointer to a template and an
  error, and it panics if the error is not nil.
  ***/
  tmpl = template.Must(template.ParseGlob("webfinances/templates/*.html"))
}

func checkSession(res http.ResponseWriter, req *http.Request) bool {
  cookie, err := req.Cookie("session_token")
  if err != nil {
    if err == http.ErrNoCookie {
      tmpl.ExecuteTemplate(res, "index_page", struct {
        Error string
      } {
          "Please loggin",
        })
    } else {
      tmpl.ExecuteTemplate(res, "index_page", struct {
        Error string
      } {
          "Bad request",
        })
    }
    return false
  }
  session, exists := sessions.Sessions[cookie.Value]
  if !exists {
    tmpl.ExecuteTemplate(res, "index_page", struct {
      Error string
    } {
        "Invalid session token",
      })
    return false
  }
  //If the session token is present, but has expired, delete the session and return
  //an unauthorized status.
  if session.IsExpired() {
    delete(sessions.Sessions, cookie.Value)
    tmpl.ExecuteTemplate(res, "index_page", struct {
      Error string
    } {
        "Session has expired",
      })
    return false
  }
  return true
}

type WfPages struct{}

func (p WfPages) IndexPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering IndexPage/webfinances.\n", m.DTF())
  tmpl.ExecuteTemplate(res, "index_page", nil)
}


/*** signup user
  //Parse and decode the request body into a new `Credentials` instance.
  creds := &sessions.Credentials{}
  if err := json.NewDecoder(req.Body).Decode(creds); err != nil {
    tmpl.ExecuteTemplate(res, "index_page", struct {
      Error string
    } {
        "Invalid username and/or password",
      })
  }
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)

***/


func (p WfPages) LoginPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering LoginPage/webfinances.\n", m.DTF())
  un := req.PostFormValue("username")
  pw := req.PostFormValue("password")
  //Get the expected password from the in memory map.
  password, ok := sessions.Users[un]
  if !ok || password != pw {
    tmpl.ExecuteTemplate(res, "index_page", struct {
      Error string
    } {
        "Invalid username and/or password",
      })
    return
  }
  /***
  Session based authentication keeps the users' sessions secure in a couple of ways:
  1. Since the session tokens are randomly generated, its near-impossible for a malicious user to
     brute-force his way into a user's session.
  2. If a user's session token is compromised somehow, it cannot be used after its expiry. This is
     why the expiry time is restricted to small intervals (a few seconds to a couple of minutes).
  ***/
  sessionToken := uuid.NewString()
  expiresAt := time.Now().Add(120 * time.Second)
  sessions.Sessions[sessionToken] = sessions.Session{
    Username: un,
    Expiry: expiresAt,
  }
  //Once a cookie is set on a client, it is sent along with every subsequent request.
  http.SetCookie(res, &http.Cookie{
    Name: "session_token",
    Value: sessionToken,
    Expires: expiresAt,
  })
  http.Redirect(res, req, "/welcome", http.StatusSeeOther)
}

func (p WfPages) LogoutPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering LogoutPage/webfinances.\n", m.DTF())
  cookie, err := req.Cookie("session_token")
  if err != nil {
    if err == http.ErrNoCookie {
      tmpl.ExecuteTemplate(res, "index_page", struct {
        Error string
      } {
          "Please loggin",
        })
    } else {
      tmpl.ExecuteTemplate(res, "index_page", struct {
        Error string
      } {
          "Bad request",
        })
    }
    return
  }
  delete(sessions.Sessions, cookie.Value)
  http.SetCookie(res, &http.Cookie{
    Name: "session_token",
    Value: "",
    Expires: time.Now(),
  })
  http.Redirect(res, req, "/", http.StatusSeeOther)
}

func (p WfPages) WelcomePage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering WelcomePage/webfinances.\n", m.DTF())
  if checkSession(res, req) {
    tmpl.ExecuteTemplate(res, "welcome_page", struct {
      Header string
      Datetime string
    } { "Investments", m.DTF() })
  }
}

func (p WfPages) ContactPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering ContactPage/webfinances.\n", m.DTF())
  if checkSession(res, req) {
    tmpl.ExecuteTemplate(res, "contact_page", struct {
      Header string
      Datetime string
    } { "Investments", m.DTF() })
  }
}

func (p WfPages) AboutPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering AboutPage/webfinances.\n", m.DTF())
  if checkSession(res, req) {
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
  //if checkSession(res, req) {
  http.ServeFile(res, req, "./webfinances/public/css/home.css")
  //}
}

func (p WfPages) FinancesPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering FinancesPage/webfinances.\n", m.DTF())
  if checkSession(res, req) {
    tmpl.ExecuteTemplate(res, "finances_page", struct {
      Header string
      Datetime string
    } { "Finances", m.DTF() })
  }
}

func (p WfPages) SimpleInterestPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering SimpleInterestPage/webfinances.\n", m.DTF())
  if checkSession(res, req) {
    tmpl.ExecuteTemplate(res, "simple_interest_page", struct {
      Header string
      Datetime string
    } { "Simple Interest", m.DTF() })
  }
}

func (p WfPages) OrdinaryAnnuityPage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering OrdinaryAnnuityPage/webfinances.\n", m.DTF())
  if checkSession(res, req) {
    tmpl.ExecuteTemplate(res, "ordinary_annuity_page", struct {
      Header string
      Datetime string
    } { "Ordinary Annuity", m.DTF() })
  }
}

func (p WfPages) AnnuityDuePage(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering AnnuityDuePage/webfinances.\n", m.DTF())
  if checkSession(res, req) {
    tmpl.ExecuteTemplate(res, "annuity_due_page", struct {
      Header string
      Datetime string
    } { "Annuity Due", m.DTF() })
  }
}


/*
https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/06.2.html

//create a global session manager in the main() function
var globalSessions *session.Manager
//Initialize the session manager.
func init() {
  globalSessions = NewSessionManager("memory", "gosessionid", 3600)
}

-----------------
package sessionmanager

type SessionManager struct {
  cookieName string  //Private cookiename
  lock sync.Mutex  //Protect session
  provider Provider
  maxlifetime int64
}

func NewSessionManager(providerName, cookieName string, maxlifetime int64) *SessionManager, error) {
  provider, ok := provides[providerName]
  if !ok {
    return nil, fmt.Errorf(Session: unknown provider %q (forgotten import?), providerName)
  }
  return &SessionManager{ provider: provider, cookieName: cookieName, maxlifetime: maxlifetime }, nil
}
*/


