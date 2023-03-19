// HTTP server.
package main

import (
	"context"
	"errors"
	"finance/webfinances"
	"fmt"
  "net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

//Environment variables.
var MAX_RETRIES int = 10
var SHUTDOWN_TIMEOUT int = 15
var PORT string = "8080"
var SVC_NAME string
var APP_NAME_VER string
var SERVER string

type handlers struct{
  mux map[string]func(http.ResponseWriter, *http.Request)
}

func (h *handlers) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  if handler, ok := h.mux[req.URL.Path]; ok {
    handler(res, req)
    return
  }
  http.NotFound(res, req) //404 - page not found
}

/***
How to kill a process using a port on localhost (Windows).
C:\> netstat -ano | findstr :<port>
C:\> taskkill /PID <PID> /F

or

C:\> npx kill-port <port>
***/
func main() {
  //http://localhost:8001/annuities/AverageRateOfReturn?ret=5.0&ret=-3.0&ret=12.0&ret=10

  var exists bool = false
  SVC_NAME, exists = os.LookupEnv("SVC_NAME")
  if !exists {
    fmt.Println("Missing environment parameter: SVC_NAME")
    return
  }
  APP_NAME_VER, exists = os.LookupEnv("APP_NAME_VER")
  if !exists {
    fmt.Println("Missing environment parameter: APP_NAME_VER")
    return
  }
  SERVER, exists = os.LookupEnv("SERVER")
  if !exists {
    fmt.Println("Missing environment parameter: SERVER")
    return
  }
  _, exists = os.LookupEnv("PORT")
  if exists {
    PORT = os.Getenv("PORT")
  }
  fmt.Printf("Using PORT: %s\n", PORT)
  SERVER += ":" + PORT
  //
  _, exists = os.LookupEnv("SHUTDOWN_TIMEOUT")
  if exists {
    sdto := os.Getenv("SHUTDOWN_TIMEOUT")
    tm, err := strconv.Atoi(sdto)
    if err == nil {
      SHUTDOWN_TIMEOUT = tm
    } else {
      fmt.Printf("'%s' is not an int number.\n", sdto)
    }
  }
  fmt.Printf("Using SHUTDOWN_TIMEOUT: %d\n", SHUTDOWN_TIMEOUT)
  //
  //Clients and Transports are safe for concurrent use by multiple goroutines and for efficiency should only be created once and re-used.
  transport := http.Transport {
    IdleConnTimeout: 1500 * time.Millisecond, //Close connection after 1500 milliseconds.
    MaxIdleConns: 2,
    MaxConnsPerHost: 2,
    MaxIdleConnsPerHost: 2,
    Dial: (&net.Dialer {
      Timeout: 1 * time.Second,
    }).Dial,
  }
  var client = &http.Client {
    Timeout: 500 * time.Millisecond, //Cancel request.
    Transport: &transport,
  }
  var a webfinances.Annuities
  var h handlers = handlers{}
  h.mux = make(map[string]func(http.ResponseWriter, *http.Request), 16)
  h.mux["/readiness"] =
  func (res http.ResponseWriter, req *http.Request) {
    fmt.Printf("\naaaaaaServer not ready. %s\n", "http://"+SERVER)
    req, err := http.NewRequest("HEAD", "http://"+SERVER, nil)
    if err != nil {
      fmt.Println("Server not ready.")
      res.WriteHeader(http.StatusInternalServerError)
      return
    }
    resp, err := client.Do(req)
    if err != nil {
      fmt.Printf("err: %v", err)
      res.WriteHeader(http.StatusInternalServerError)
      return
    }
    resp.Body.Close()
    fmt.Println("Server is ready.")
    //https://go.dev/src/net/http/status.go
    res.WriteHeader(http.StatusOK)
  }
  h.mux["/fin/annuities/AverageRateOfReturn"] = a.AverageRateOfReturn
  h.mux["/fin/annuities/GrowthDecayOfFunds"] = a.GrowthDecayOfFunds
  server := &http.Server {
    /***
    By not specifying an IP address before the colon, the server will listen on every IP address
    associated with the computer, and it will listen on port PORT.
    ***/
    Addr: ":" + PORT,
    Handler: &h,
  }
  /***
  A channel is a communication mechanism that lets one goroutine send values to another goroutine.
  Each channel is a conduit for values of a particular type, called the channel's element type.

  As with maps, a channel is a reference to the data structure created by make. When we copy a
  channel or pass one as an argument to a function, we are copying a reference, so caller and
  callee refer to the same data structure. As with other reference types, the zero value of a
  channel is nil.
  ***/
  signalChan := make(chan os.Signal, 1) //Buffered channel capacity 1; notifier will not block.
  /***
  When Shutdown is called, Serve, ListenAndServe, and ListenAndServeTLS immediately return
  ErrServerClosed. Make sure the program doesn't exit and waits instead for Shutdown to return.
  ***/
  waitMainChan := make(chan struct{})
  go func() {
    /***
    signal.Notify disables the default behavior for a given set of asynchronous signals and instead
    delivers them over one or more registered channels.
    https://pkg.go.dev/os/signal#hdr-Default_behavior_of_signals_in_Go_programs
    ***/
    signal.Notify(signalChan,
      syscall.SIGINT, //Ctrl-C
      syscall.SIGTERM, //Kubernetes sends a SIGTERM.
    )
    <- signalChan //Waiting for the signal; signal is discarded.
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(SHUTDOWN_TIMEOUT) * time.Second)
    defer func() {
      //Extra handling goes here...
      close(signalChan)
      //Calling cancel() releases the resources associated with the context.
      cancel()
      close(waitMainChan) //Shutdown is done; let the main goroutine terminate.
    }()
    //https://pkg.go.dev/net/http#Server.Shutdown
    if err := server.Shutdown(ctx); err != nil {
      fmt.Printf("Server shutdown failed: %+v", err) //https://pkg.go.dev/fmt
    }
  }()
  fmt.Printf("%s - Starting the server at port %s...\n", time.Now().UTC().Format(time.RFC3339Nano), PORT)
  /***
  ListenAndServe runs forever, or until the server fails (or fails to start) with an error,
  always non-nil, which it returns.

  The web server invokes each handler in a new goroutine, so handlers must take precautions such as
  locking when accessing variables that other goroutines, including other requests to the same
  handler, may be accessing.
  ***/
  err := server.ListenAndServe()
  if errors.Is(err, http.ErrServerClosed) {
    fmt.Println("Server has been closed.")
  } else if err != nil {
    fmt.Printf("Server error: %s\n", "err")
    signalChan <- syscall.SIGINT //Let the goroutine finish.
  }
  <- waitMainChan //Block until shutdown is done.
  return
}




/********
started := time.Now()
http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
    duration := time.Now().Sub(started)
    if duration.Seconds() > 10 {
        w.WriteHeader(500)
        w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
    } else {
        w.WriteHeader(200)
        w.Write([]byte("ok"))
    }
})


if err := redirectServer.Shutdown(ctx); err == context.DeadlineExceeded {
	return fmt.Errorf("%v timeout exceeded while waiting on HTTP shutdown", redirectTimeout)
}

************/