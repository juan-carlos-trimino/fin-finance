// HTTP server.
package main

/***
*****************************************************************************************
*** To run the app in a K8s environment, do NOT set the environment variable NOT_K8S. ***
*****************************************************************************************

To run the app as a standalone HTTP server:
Compile and run the app.
$ go build -o finance && NOT_K8S= ./finance

To change the PORT.
$ go build -o finance && NOT_K8S= PORT=8181 ./finance

Compile and run the app in the background.
$ go build -o finance && NOT_K8S= ./finance &

Force rebuilding of packages.
$ go build -o finance -a && NOT_K8S= ./finance

Compile and run (in the background) at the same time.
$ NOT_K8S= go run main.go &

***************************************************************************************************

How to kill a process using a port on localhost (Windows).
C:\> netstat -ano | findstr :<port>
C:\> taskkill /PID <PID> /F
or
C:\> npx kill-port <port>

Linux
$ ps -a
$ kill <PID>

To display the headers:
$ curl.exe -IL "http://localhost:8080"

PS> curl.exe "http://localhost:8080"
***/

import (
	"context"
	"errors"
	"finance/middlewares"
	"finance/misc"
  "finance/sessions"
	"finance/webfinances"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (  //Environment variables.
  MAX_RETRIES int = 10
  SHUTDOWN_TIMEOUT int = 15
  PORT string = "8080"
  SVC_NAME string
  APP_NAME_VER string
  SERVER string = "localhost"
)

var m = misc.Misc{}

/***
In Go, a handler is an interface (type Handler interface) that has a method named ServeHTTP with
two parameters: an http.ResponseWriter interface and a pointer to an http.Request struct. Hence,
any type that has a method called ServeHTTP with this method signature is a handler:
  ServeHTTP(http.ResponseWriter, *http.Request)

ServeMux is an HTTP request multiplexer; it accepts an HTTP request and redirects it to the
correct handler according to the URL in the request. ServeMux is a struct with a map of entries
that maps a URL to a handler, and it is also a handler because it implements the ServeHTTP method.
The ServeHTTP method finds the URL most closely mathing the requested one and calls the
corresponding handler.

Since ServeMux is a struct, DefaultServeMux is an instance of ServeMux.

Go has a function type named HandlerFunc, which will adapt a function f with the appropriate
signature into a Handler with a method f.
***/
type handlers struct {
  /***
  The 'type HandlerFunc func(ResponseWriter, *Request)' is a function type that has methods
  ('func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)') and satisfies an interface,
  http.Handler. The behavior of its ServeHTTP method is to call the underlying function.
  HandlerFunc is thus an adapter that lets a function value satisfy an interface, where the
  function and the interface's sole method have the same signature.
  ***/
  mux map[string]http.HandlerFunc  //Multiplexer.
}

/***
A Handler responds to an HTTP request.
'ServeHTTP' is the only method of the 'type Handler interface'.

The ServeHTTP function takes two parameters -- the ResponseWriter interface and a pointer to a
Request struct. Since changes to Request by the handler need to be visible to the server, it is
passed by reference. But why is ResponseWriter passed by value? ResponseWriter is an interface to a
nonexported struct response; we're passing the struct by reference (we're passing in a pointer to
response) and not by value.
***/
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

func main() {
  var exists bool = false
  var ev string
  _, exists = os.LookupEnv("NOT_K8S")
  if !exists {
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
  }
  ev, exists = os.LookupEnv("PORT")
  if exists {
    PORT = ev
  }
  fmt.Printf("%s - Using PORT: %s\n", m.DTF(), PORT)
  SERVER += ":" + PORT
  ev, exists = os.LookupEnv("SHUTDOWN_TIMEOUT")
  if exists {
    tm, err := strconv.Atoi(ev)
    if err == nil {
      SHUTDOWN_TIMEOUT = tm
    } else {
      fmt.Printf("%s - '%s' is not an int number.\n", m.DTF(), ev)
    }
  }
  fmt.Printf("%s - Using SHUTDOWN_TIMEOUT: %d\n", m.DTF(), SHUTDOWN_TIMEOUT)
  var wfpages = webfinances.WfPages{}
  var wfadcp = webfinances.NewWfAdCpPages()
  var wfadepp = webfinances.NewWfAdEppPages()
  var wfadfv = webfinances.NewWfAdFvPages()
  var wfadpv = webfinances.NewWfAdPvPages()
  var wfoainterest = webfinances.NewWfOaInterestRatePages()
  var wfoapv = webfinances.NewWfOaPvPages()
  var wfoafv = webfinances.NewWfOaFvPages()
  var wfoacp = webfinances.NewWfOaCpPages()
  var wfoaepp = webfinances.NewWfOaEppPages()
  var wfoaga = webfinances.NewWfOaGaPages()
  var wfoaperpetuity = webfinances.NewWfOaPerpetuityPages()
  var wfmortgage = webfinances.NewWfMortgagePages()
  var wfbonds = webfinances.NewWfBondsPages()
  var wfsia = webfinances.NewWfSiAccuratePages()
  var wfsio = webfinances.NewWfSiOrdinaryPages()
  var wfsib = webfinances.NewWfSiBankersPages()
  var wfmisc = webfinances.NewWfMiscellaneousPages()
  /***
  The Go web server will route requests to different functions depending on the requested URL.
  ***/
  var h handlers = handlers{}
  h.mux = make(map[string]http.HandlerFunc, 64)
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
  //Serve static files; i.e., the server will serve it as it is, without processing it first.
  h.mux["/public/css/home.css"] = wfpages.PublicHomeFile
  h.mux["/favicon.ico"] = faviconHandler
  h.mux["/"] = wfpages.IndexPage
  h.mux["/login"] = wfpages.LoginPage
  h.mux["/logout"] = wfpages.LogoutPage
  h.mux["/welcome"] = wfpages.WelcomePage
  h.mux["/contact"] = wfpages.ContactPage
  h.mux["/about"] = wfpages.AboutPage
  h.mux["/finances"] = wfpages.FinancesPage
  h.mux["/fin/ordinaryannuity"] = wfpages.OrdinaryAnnuityPage
  h.mux["/fin/ordinaryannuity/interestrate"] = wfoainterest.OaInterestRatePages
  h.mux["/fin/ordinaryannuity/fv"] = wfoafv.OaFvPages
  h.mux["/fin/ordinaryannuity/pv"] = wfoapv.OaPvPages
  h.mux["/fin/ordinaryannuity/cp"] = wfoacp.OaCpPages
  h.mux["/fin/ordinaryannuity/epp"] = wfoaepp.OaEppPages
  h.mux["/fin/ordinaryannuity/ga"] = wfoaga.OaGaPages
  h.mux["/fin/ordinaryannuity/perpetuity"] = wfoaperpetuity.OaPerpetuityPages
  h.mux["/fin/annuitydue"] = wfpages.AnnuityDuePage
  h.mux["/fin/annuitydue/cp"] = wfadcp.AdCpPages
  h.mux["/fin/annuitydue/epp"] = wfadepp.AdEppPages
  h.mux["/fin/annuitydue/fv"] = wfadfv.AdFvPages
  h.mux["/fin/annuitydue/pv"] = wfadpv.AdPvPages
  h.mux["/fin/bonds"] = wfbonds.BondsPages
  h.mux["/fin/mortgage"] = wfmortgage.MortgagePages
  h.mux["/fin/simpleinterest"] = wfpages.SimpleInterestPage
  h.mux["/fin/simpleinterest/accurate"] = wfsia.SimpleInterestAccuratePages
  h.mux["/fin/simpleinterest/bankers"] = wfsib.SimpleInterestBankersPages
  h.mux["/fin/simpleinterest/ordinary"] = wfsio.SimpleInterestOrdinaryPages
  h.mux["/fin/miscellaneous"] = wfmisc.MiscellaneousPages
  commonMiddlewares := []middlewares.Middleware {
    middlewares.CorrelationId,
  }
  for idx, f := range h.mux {
    h.mux[idx] = middlewares.ChainMiddlewares(f, commonMiddlewares)
  }
  sessions.Users["jct"] = "pw"
  server := &http.Server {  //https://pkg.go.dev/net/http#ServeMux
    /***
    By not specifying an IP address before the colon, the server will listen on every IP address
    associated with the computer, and it will listen on port PORT.
    ***/
    Addr: ":" + PORT,
    /***
      Connection accepted
      │
      │  Wait for the client to send the request
      │  │
      │  │        If enabled (The TLS handshake doesn't have to be repeated with an already established connection.)
      │  │        │
      │  │        │           Read the headers
      │  │        │           │
      │  │        │           │        Read the body
      │  │        │           │        │
      │  │        │           │        │          Write the response
      │  │        │           │        │          │
      │  │        │           │        │          │
      ╔══════╦═══════════╦═════════╦═════════╦══════════╦════════╗
      ║ Wait ║   TLS     ║ Request ║ Request ║ Response ║  Idle  ║
      ║      ║ handshake ║ headers ║  body   ║          ║        ║
      ╚══════╩═══════════╩═════════╩═════════╩══════════╩════════╝
                                   <-------------------->
                                         HTTP handler
                                                        <-------->
                                                        IdleTimeout
                                                     (Keep-alive only)
                         <--------->
                 http.Server.ReadHeaderTimeout
      <-------------------------------------->
              http.Server.ReadTimeout
                                   <-------------------->
                                    http.TimeoutHandler

      The five steps of an HTTP response and the related timeouts.
    ***/
    //It specifies the maximum amount of time to read the request headers.
    ReadHeaderTimeout: 250 * time.Millisecond,
    /***
    It specifies the maximum amount of time to read the entire request.
    ReadTimeout = ReadHeaderTimeout + TimeoutHandler + Extra time
    **/
    ReadTimeout: 990 * time.Millisecond,
    /***
    If a handler fails to respond on time, the server will reply with "503 Service Unavailable" and
    the specified message; the context passed to the handler will be canceled.
    Note: The http.Server.WriteTimeout is not necessary since http.TimeoutHandler is being used.
    ***/
    Handler: http.TimeoutHandler(&h, 700 * time.Millisecond, "Request timeout."),
    /***
    It configures the maximum amount of time for the next request when keep-alives are enabled.
    Note that if http.Server.IdleTimeout isn't set, the value of http.Server.ReadTimeout is used
    for the idle timeout. If neither is set, there won't be any timeouts, and connections will
    remain open until they are closed by clients.
    ***/
    IdleTimeout: 30 * time.Second,
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
  err := (*server).ListenAndServe()
  if errors.Is(err, http.ErrServerClosed) {
    fmt.Printf("%s - Server has been closed.\n", m.DTF())
  } else if err != nil {
    fmt.Printf("%s - Server error: %+v\n", m.DTF(), err)
    signalChan <- syscall.SIGINT //Let the goroutine finish.
  }
  <- waitMainChan //Block until shutdown is done.
}

func faviconHandler(res http.ResponseWriter, req *http.Request) {
  http.NotFound(res, req)  //404 - page not found.
}
