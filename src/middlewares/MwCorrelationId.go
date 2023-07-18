package middlewares

import (
	"context"
  "fmt"
  "net/http"
  "github.com/google/uuid"
  "time"
)

/***
To avoid context keys collisions, a best practice is to create an unexported custom type.
***/
//type ctxKey string
/***
The correlationIdKey constant is unexported. Hence, there's no risk that another package using the
same context could override the value that is already set. Even if another package creates the same
correlationIdKey based on a ctxKey type as well, it will be a different key.
***/
const correlationIdKey /*ctxKey*/ string = "correlationIdKey"

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
