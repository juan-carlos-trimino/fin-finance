// HTTP server.
package main

import (
  "context"
  "errors"
  "finance/webfinances"
  "fmt"
  "net/http"
  "finance/misc"
  // "mime"
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
var SERVER string = "localhost"

// func init() {
//   fmt.Printf("%s - Entering init/main.\n", time.Now().UTC().Format(time.RFC3339Nano))
//   // mime.AddExtensionType(".js", "application/javascript; charset=utf-8")
//   ct := mime.TypeByExtension(".css")
//   fmt.Printf("ct: %s\n", ct)
//   // mime.AddExtensionType(".css", "text/css; charset=utf-8")
// }

var m = misc.Misc{}

type handlers struct {
  mux map[string]func(http.ResponseWriter, *http.Request)
}

func (h *handlers) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("%s - Entering ServeHTTP/main.\n", m.DTF())
  fmt.Printf("%s - Method: %s, Request URI: %s\n", m.DTF(), req.Method, req.RequestURI)
  //Implement route forwarding.
  if handler, ok := h.mux[req.URL.Path]; ok {
    fmt.Printf("%s - URL Path: %s\n", m.DTF(), req.URL.Path)
    handler(res, req)
    return
  }
  http.NotFound(res, req)  //404 - page not found.
}

/***
How to kill a process using a port on localhost (Windows).
C:\> netstat -ano | findstr :<port>
C:\> taskkill /PID <PID> /F
or
C:\> npx kill-port <port>

To display the headers:
$ curl.exe -IL "http://localhost:8080"
***/
func main() {
  var exists bool = false
  /*** k8s
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
  k8s ***/
  _, exists = os.LookupEnv("PORT")
  if exists {
    PORT = os.Getenv("PORT")
  }
  fmt.Printf("%s - Using PORT: %s\n", m.DTF(), PORT)
  SERVER += ":" + PORT
  _, exists = os.LookupEnv("SHUTDOWN_TIMEOUT")
  if exists {
    sdto := os.Getenv("SHUTDOWN_TIMEOUT")
    tm, err := strconv.Atoi(sdto)
    if err == nil {
      SHUTDOWN_TIMEOUT = tm
    } else {
      fmt.Printf("%s - '%s' is not an int number.\n", m.DTF(), sdto)
    }
  }
  fmt.Printf("%s - Using SHUTDOWN_TIMEOUT: %d\n", m.DTF(), SHUTDOWN_TIMEOUT)
  var wfp webfinances.WebFinancesPages
  var wfa webfinances.Annuities
  var h handlers = handlers{}
  h.mux = make(map[string]func(http.ResponseWriter, *http.Request), 16)
  h.mux["/readiness"] =
  func (res http.ResponseWriter, req *http.Request) {
    fmt.Printf("\naaaaaaServer not ready. %s\n", SERVER)
    // req, err := http.NewRequest(http.MethodHead, SERVER, nil)
    // if err != nil {
    //   fmt.Println("Server not ready.")
    //   res.WriteHeader(http.StatusInternalServerError)
    //   return
    // }
    // resp, err := client.Do(req)
    // if err != nil {
    //   fmt.Printf("client: error making http request: %s\n", err)
    //   res.WriteHeader(http.StatusInternalServerError)
    //   return
    // }
    // resp.Body.Close()
    // fmt.Println("Server is ready.")
    // //https://go.dev/src/net/http/status.go
    res.WriteHeader(http.StatusOK)
  }
  h.mux["/"] = wfp.HomePage
  h.mux["/public/css/home.css"] = wfp.PublicHomeFile
  h.mux["/contact"] = wfp.ContactPage
  h.mux["/about"] = wfp.AboutPage
  h.mux["/finances"] = wfp.FinancesPage
  h.mux["/fin/simpleinterest"] = wfp.SimpleInterestPage
  h.mux["/fin/simpleinterest/ordinary"] = wfp.SimpleInterestOrdinaryPage
  h.mux["/fin/simpleinterest/ordinarycompute"] = wfp.SimpleInterestOrdinaryCompute
  h.mux["/fin/annuities/AverageRateOfReturn"] = wfa.AverageRateOfReturn
  h.mux["/fin/annuities/GrowthDecayOfFunds"] = wfa.GrowthDecayOfFunds
  server := &http.Server {  //https://pkg.go.dev/net/http#ServeMux
    /***
    By not specifying an IP address before the colon, the server will listen on every IP address
    associated with the computer, and it will listen on port PORT.
    ***/
    Addr: ":" + PORT,
    Handler: &h,
    MaxHeaderBytes: 1 << 20,  //1 MB.
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
    fmt.Printf("%s - Waiting for notification to shut down the server.\n", m.DTF())
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
      fmt.Printf("%s - Server shutdown failed: %+v\n", m.DTF(), err)  //https://pkg.go.dev/fmt
    }
  }()
  fmt.Printf("%s - Starting the server at port %s...\n", m.DTF(), PORT)
  /***
  ListenAndServe runs forever, or until the server fails (or fails to start) with an error,
  always non-nil, which it returns.

  The web server invokes each handler in a new goroutine, so handlers must take precautions such as
  locking when accessing variables that other goroutines, including other requests to the same
  handler, may be accessing.
  ***/
  err := server.ListenAndServe()
  if errors.Is(err, http.ErrServerClosed) {
    fmt.Printf("%s - Server has been closed.\n", m.DTF())
  } else if err != nil {
    fmt.Printf("%s - Server error: %+v\n", m.DTF(), err)
    signalChan <- syscall.SIGINT //Let the goroutine finish.
  }
  <- waitMainChan //Block until shutdown is done.
  return
}
