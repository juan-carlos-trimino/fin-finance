package middlewares

import (
  "context"
  "finance/sessions"
  "net/http"
)

//Protect private pages.
func ValidateSessions(handler http.HandlerFunc) http.HandlerFunc {
  return func(res http.ResponseWriter, req *http.Request) {
    var ctx context.Context
    cookie, err := req.Cookie("session_token")
    if err != nil {
      ctx = context.WithValue(req.Context(), sessionStatusKey, false)
    } else if session, exists := sessions.Sessions[cookie.Value]; !exists {
      ctx = context.WithValue(req.Context(), sessionStatusKey, false)
    //If the session token is present, but has expired, delete the session and return
    //an unauthorized status.
    } else if session.IsExpired() {
      delete(sessions.Sessions, cookie.Value)
      ctx = context.WithValue(req.Context(), sessionStatusKey, false)
    } else {
      csrf := req.PostFormValue("csrf_token")
      if err := sessions.CompareHashAndPassword([]byte(csrf), sessions.Sessions[cookie.Value].CsrfToken); err == nil {
        ctx = context.WithValue(req.Context(), sessionStatusKey, true)
        ctx = context.WithValue(req.Context(), sessionTokenKey, cookie.Value)
      } else {
        ctx = context.WithValue(req.Context(), sessionStatusKey, false)
      }
    }
    handler.ServeHTTP(res, req.WithContext(ctx))
  }
}
