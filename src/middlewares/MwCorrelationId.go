package middlewares

import (
	"context"
  "fmt"
  "net/http"
  "github.com/google/uuid"
  "time"
)

func CorrelationId(handler http.HandlerFunc) http.HandlerFunc {
  return func(res http.ResponseWriter, req *http.Request) {
    start := time.Now()
    uuid := uuid.New()
    //Creating a new context from a parent context.
    ctx := context.WithValue(req.Context(), correlationIdKey, uuid.String())
    fmt.Printf("New request with correlation id: %s\n", uuid)
    //Calling the handler with the new context.
    handler.ServeHTTP(res, req.WithContext(ctx))
    fmt.Printf("Request took %vms\nRequest correlation id: %s\n",
      time.Since(start).Microseconds(), uuid)
  }
}
