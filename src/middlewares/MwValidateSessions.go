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
      ctx = context.WithValue(req.Context(), sessionTokenKey, "")
    } else if exists := sessions.SessionExists(cookie.Value); !exists {
      ctx = context.WithValue(req.Context(), sessionTokenKey, "")
    //If the session token is present, but has expired, delete the session and return
    //an unauthorized status.
    } else if sessions.IsSessionExpired(cookie.Value) {
      ctx = context.WithValue(req.Context(), sessionTokenKey, "")
    } else if req.Method == http.MethodPost {
      csrf := req.PostFormValue("csrf_token")
      if !sessions.CompareUuids(csrf, cookie.Value) {
        ctx = context.WithValue(req.Context(), sessionTokenKey, "")
      } else {
        //ctx = context.WithValue(context.Background(), sessionStatusKey, true)
        //ctx = context.WithValue(ctx, sessionTokenKey, cookie.Value)
        ctx = context.WithValue(req.Context(), sessionTokenKey, cookie.Value)
      }
    } else {
      ctx = context.WithValue(req.Context(), sessionTokenKey, cookie.Value)
    }
    handler.ServeHTTP(res, req.WithContext(ctx))
  }
}
