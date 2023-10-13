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
var sessions = map[string]session_token{}  //key: sessionToken, value: session

type session_token struct {
  Username string
  Expiry time.Time  //Enforce periodic session termination as a way to prevent session hijacking.
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
  session, exists := sessions[sessionToken]
  if exists {
    return strings.EqualFold(csrf, session.CsrfToken)
  }
  return exists
}

func AddEntryToSessions(userName string) (sessionToken string, session session_token) {
  /***
  Session based authentication keeps the users' sessions secure in a couple of ways:
  1. Since the session tokens are randomly generated, its near-impossible for a malicious user to
     brute-force his way into a user's session.
  2. If a user's session token is compromised somehow, it cannot be used after its expiry. This is
     why the expiry time is restricted to small intervals (a few seconds to a couple of minutes).
  ***/
  sessionToken = uuid.NewString()
  sessions[sessionToken] = session_token{
    Username: userName,
    Expiry: time.Now().Add(120 * time.Second),
    CsrfToken: uuid.NewString(),
  }
  session = sessions[sessionToken]
  return
}

func UpdateEntryInSessions(oldSessionToken string) (newSessionToken string, session session_token) {
  newSessionToken = uuid.NewString()
  sessions[newSessionToken] = session_token{
    Username: sessions[oldSessionToken].Username,
    Expiry: time.Now().Add(120 * time.Second),
    CsrfToken: uuid.NewString(),
  }
  delete(sessions, oldSessionToken)
  session = sessions[newSessionToken]
  return
}

func GetNewUuid() string {
  return uuid.NewString()
}

func CreateCookie(sessionToken string) (cookie *http.Cookie) {
  //https://en.wikipedia.org/wiki/HTTP_cookie
  //https://httpwg.org/specs/rfc6265.html
  cookie = &http.Cookie{
    Name: "session_token",
    Value: sessionToken,
    Path: "/",
    // Expires: sessions[sessionToken].Expiry,
    HttpOnly: true,
    SameSite: http.SameSiteStrictMode,
    Secure: false,
  }
  return
}

func DeleteSession(sessionToken string) (cookie *http.Cookie) {
  delete(sessions, sessionToken)
  cookie = &http.Cookie{
    Name: "session_token",
    Value: "",
    Path: "/",
    MaxAge: -1,
    HttpOnly: true,
    SameSite: http.SameSiteStrictMode,
    Secure: false,
  }
  return
}
