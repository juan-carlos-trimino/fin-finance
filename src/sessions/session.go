package sessions

import (
  //The option -u instructs 'get' to update the module with dependencies.
  //go get -u github.com/google/uuid
  "github.com/google/uuid"
  //The option -u instructs 'get' to update the module with dependencies.
  //go get -u golang.org/x/crypto/bcrypt
  "golang.org/x/crypto/bcrypt"
  "net/http"
  "strings"
  "time"
)

//Store the session information for each user in memory.
var sessions = map[string]session{}  //key: sessionToken, value: session
//Store the username and password for each user.
var Users = map[string][]byte{}  //key: username, value: password

type session struct {
  Username string
  Expiry time.Time  //xxxxxxxxxxxxxxxxxxxx enforce periodic session termination as a way to prevent session hijacking.
//  lock sync.Mutex  //Protect session
  CsrfToken string
}

//Determine if a session has expired.
func IsSessionExpired(sessionToken string) bool {
  s, exists := sessions[sessionToken]
  if exists {
    var expired bool = s.Expiry.Before(time.Now())
    if expired {
      //Delete the session.
      delete(sessions, sessionToken)
    }
    return expired
  }
  return !exists
}

func SessionExists(sessionToken string) bool {
  _, exists := sessions[sessionToken]
  return exists
}

func HashSecret(secret string) ([]byte, error) {
  hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
  return hashedSecret, err
}

func CompareHashAndPassword(hashedPassword, password []byte) (bool, error) {
  err := bcrypt.CompareHashAndPassword(hashedPassword, password)
  return err == nil, err
}

func CompareUuids(csrf, sessionToken string) bool {
  s, exists := sessions[sessionToken]
  if exists {
    return strings.EqualFold(csrf, s.CsrfToken)
  }
  return exists
}

func AddEntryToSessions(userName string) (sessionToken string, s session) {
  /***
  Session based authentication keeps the users' sessions secure in a couple of ways:
  1. Since the session tokens are randomly generated, its near-impossible for a malicious user to
     brute-force his way into a user's session.
  2. If a user's session token is compromised somehow, it cannot be used after its expiry. This is
     why the expiry time is restricted to small intervals (a few seconds to a couple of minutes).
  ***/
  sessionToken = uuid.NewString()
  sessions[sessionToken] = session{
    Username: userName,
    Expiry: time.Now().Add(120 * time.Second),
    CsrfToken: uuid.NewString(),
  }
  s = sessions[sessionToken]
  return
}

func UpdateEntryInSessions(oldSessionToken string) (newSessionToken string, s session) {
  newSessionToken = uuid.NewString()
  sessions[newSessionToken] = session{
    Username: sessions[oldSessionToken].Username,
    Expiry: time.Now().Add(120 * time.Second),
    CsrfToken: uuid.NewString(),
  }
  delete(sessions, oldSessionToken)
  s = sessions[newSessionToken]
  return
}

func GetNewUuid() string {
  return uuid.NewString()
}

func CreateCookie(sessionToken string) (cookie *http.Cookie) {
  cookie = &http.Cookie{
    Name: "session_token",
    Value: sessionToken,
    /***
    It indicates the path that must exist in the requested URL for the browser to send the Cookie
    header. The forward slash (/) character is interpreted as a directory separator, and
    subdirectories are matched as well. For example, for Path=/docs:
    * The request paths /docs, /docs/, /docs/Web/, and /docs/Web/HTTP will all match.
    * The request paths /, /docsets, /fr/docs will not match.
    ***/
    Path: "/",
    //Expires: sessions[sessionToken].Expiry,
    /***
    An http-only cookie cannot be accessed by client-side APIs, such as JavaScript. This
    restriction eliminates the threat of cookie theft via Cross-Site Scripting (XSS). However, the
    cookie remains vulnerable to Cross-Site Tracing (XST) and Cross-Site Request Forgery (CSRF)
    attacks.
    ***/
    HttpOnly: true,
    /***
    The attribute SameSite can have a value of Strict, Lax or None. With attribute SameSite=Strict,
    the browsers would only send cookies to a target domain that is the same as the origin domain.
    This would effectively mitigate Cross-Site Request Forgery (CSRF) attacks. With SameSite=Lax,
    browsers would send cookies with requests to a target domain even it is different from the
    origin domain, but only for safe requests such as GET (POST is unsafe) and not third-party
    cookies (inside iframe). Attribute SameSite=None would allow third-party (cross-site) cookies,
    however, most browsers require secure attribute on SameSite=None cookies.
    ***/
    SameSite: http.SameSiteStrictMode,
  }
  return
}

func DeleteSession(sessionToken string) (cookie *http.Cookie) {
  delete(sessions, sessionToken)
  cookie = &http.Cookie{
    Name: "session_token",
    Value: "",
    Path: "/",
    // Expires: time.Now(),
    /***
    It indicates the number of seconds until the cookie expires. A zero or negative number will
    expire the cookie immediately. If both Expires and Max-Age are set, Max-Age has precedence.
    ***/
    MaxAge: -1,
    HttpOnly: true,
    SameSite: http.SameSiteStrictMode,
  }
  return
}
