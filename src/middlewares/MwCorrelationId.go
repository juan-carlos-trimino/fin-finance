package middlewares

import (
  "fmt"
  "net/http"
  "github.com/google/uuid"
  "time"
)

func CorrelationId(handler http.HandlerFunc) http.HandlerFunc {
  return func(res http.ResponseWriter, req *http.Request) {
    start := time.Now()
    
    uuid := uuid.New()
    

//    res.Header().Set("X-Correlation-Id", uuid)
    fmt.Printf("New request with correlation id: %s\n", uuid)
    handler.ServeHTTP(res, req)
    fmt.Printf("Time to complete request (ns): %v\nReturning response with correlation id: %s\n",
                time.Since(start).Nanoseconds(), uuid)
  }
}


