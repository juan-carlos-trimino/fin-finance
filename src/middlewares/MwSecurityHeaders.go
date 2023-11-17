package middlewares

import (
  "net/http"
)

/***
Send HTTP HEAD request with curl:
$ curl -I http://localhost:8080
$ curl --head http://localhost:8080
***/
func SecurityHeaders(handler http.HandlerFunc) http.HandlerFunc {
  return func(res http.ResponseWriter, req *http.Request) {
    handler.ServeHTTP(res, req)
    res.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubdomains")
  }
}
