package middlewares

import (
  "net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func ChainMiddlewares(originalHandler http.HandlerFunc, mw []Middleware) http.HandlerFunc {
  wrapHandler := originalHandler
  length := len(mw)
  for idx := length - 1; idx > -1; idx-- {
    wrapHandler = mw[idx](wrapHandler)
  }
  return wrapHandler
}
