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
var Sessions = map[string]Session{}  //key: sessionToken, value: session
//Store the username and password for each user.
var Users = map[string][]byte{}  //key: username, value: password

type Session struct {
  Username string
  Expiry time.Time  //xxxxxxxxxxxxxxxxxxxx enforce periodic session termination as a way to prevent session hijacking.
//  lock sync.Mutex  //Protect session
  CsrfToken string
}

//Determine if a session has expired.
func (s *Session) IsExpired(sessionToken string) bool {
  var expired bool = s.Expiry.Before(time.Now())
  if expired {
    //Delete the session.
    delete(Sessions, sessionToken)
  }
  return expired
}

func HashSecret(secret string) ([]byte, error) {
  hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
  return hashedSecret, err
}

func CompareHashAndPassword(hashedPassword, password []byte) (bool, error) {
  err := bcrypt.CompareHashAndPassword(hashedPassword, password)
  return err == nil, err
}

func CompareUuids(u1, u2 string) bool {
  return strings.EqualFold(u1, u2)
}

func AddEntryToSessions(userName string) (sessionToken string) {
  /***
  Session based authentication keeps the users' sessions secure in a couple of ways:
  1. Since the session tokens are randomly generated, its near-impossible for a malicious user to
     brute-force his way into a user's session.
  2. If a user's session token is compromised somehow, it cannot be used after its expiry. This is
     why the expiry time is restricted to small intervals (a few seconds to a couple of minutes).
  ***/
  sessionToken = uuid.NewString()
  Sessions[sessionToken] = Session{
    Username: userName,
    Expiry: time.Now().Add(120 * time.Second),
    CsrfToken: uuid.NewString(),
  }
  return
}

func UpdateEntryInSessions(oldSessionToken string) (newSessionToken string) {
  newSessionToken = uuid.NewString()
  Sessions[newSessionToken] = Session{
    Username: Sessions[oldSessionToken].Username,
    Expiry: time.Now().Add(120 * time.Second),
    CsrfToken: uuid.NewString(),
  }
  delete(Sessions, oldSessionToken)
  return
}

func GetNewUuid() string {
  return uuid.NewString()
}

func CreateCookie(sessionToken string) (cookie *http.Cookie) {
  cookie = &http.Cookie{
    Name: "session_token",
    Value: sessionToken,
    Path: "/",
    Expires: Sessions[sessionToken].Expiry,
  }
  return
}





/***

type Provider interface {
  SessionInit(sessionToken string) (SessionManager, error)
  SessionRead(sessionToken string) (SessionManager, error)
  SessionDestroy(sessionToken string) error
  SessionGC(maxLifetime int64)
}

var providers = make(map[string]Provider, 8)

func Register(name string, provider Provider) {
  if provider == nil {
    panic("session: Register provider is nil")
  } else if _, dup := providers[name]; dup {
    panic("session: Register call twice for provider " + name)
  }
  providers[name] = provider
}


type Session1 interface {
  Set(key, value interface{}) error
  Get(key interface{}) interface{}
  Delete(key interface{}) error
  SessionToken() string
}

type SessionManager struct {
  sessionName string
  lock sync.Mutex
  provider Provider
  maxLifetime int64
}

func NewSessionManager(providerName, sessionName string, maxLifetime int64) (*SessionManager, error) {
  provider, ok := providers[providerName]
  if !ok {
    return nil, fmt.Errorf("Session: Unknown provider %q (forgotten import?)", providerName)
  }
  return &SessionManager{
    provider: provider,
    sessionName: sessionName,
    maxLifetime: maxLifetime,
  }, nil
}


func (sm *SessionManager) sessionId() string {
  b := make([]byte, 32)
  if _, err := io.ReadFull(rand.Reader, b); err != nil {
    return ""
  }
  return base64.URLEncoding.EncodeToString(b)
}


func (sm *SessionManager) SessionDestroy(res http.ResponseWriter, req *http.Request) {
  cookie, err := req.Cookie(sm.sessionName)
  if err != nil || cookie.Value == "" {
    return
  } else {
    sm.lock.Lock()
    defer sm.lock.Unlock()
    sm.provider.SessionDestroy(cookie.Value)
    expiration := time.Now()
    http.SetCookie(res, &http.Cookie{
      Name: sm.sessionName,
      Expires: expiration,
      // Path: "/",
      // HttpOnly: true,
      // MaxAge: -1,
    })
  }
}

**/