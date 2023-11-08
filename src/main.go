//go:build !windows
// +build !windows

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
	"crypto/tls"
	// "crypto/tls"
	"errors"
	"finance/middlewares"
	"finance/misc"
	"finance/security"
	"finance/sessions"
	"finance/webfinances"
	"fmt"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof" //Blank import of pprof.
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (  //Environment variables.
  MAX_RETRIES int = 10
  SHUTDOWN_TIMEOUT int = 15
  PORT string = "8443"
  // PORT string = "8080"
  SVC_NAME string
  APP_NAME_VER string
  SERVER string = "localhost"
  USER_NAME = "jct"
  PASSWORD = "pw"
)

var m = misc.Misc{}

/***
var sessionManager *sessions.SessionManager
//Initialize the session manager.
func init() {
  sessionManager = sessions.NewSessionManager("memory", "session_token", 3600)
}
***/


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
  fmt.Printf("Using PORT: %s\n", PORT)
  SERVER += ":" + PORT
  ev, exists = os.LookupEnv("SHUTDOWN_TIMEOUT")
  if exists {
    tm, err := strconv.Atoi(ev)
    if err == nil {
      SHUTDOWN_TIMEOUT = tm
    } else {
      fmt.Printf("'%s' is not an int number.\n", ev)
    }
  }
  fmt.Printf("Using SHUTDOWN_TIMEOUT: %d\n", SHUTDOWN_TIMEOUT)

////////////////////////////
rootCert, rootCertPEM, rootPrivKey := security.GenRootCA()
fmt.Println("rootCert\n", string(rootCertPEM))
// interCert, interCertPEM, interPrivKey := security.GenIntermediateCA(rootCert, rootPrivKey)
// fmt.Println("Intermediate Cert CA\n", string(interCertPEM))
// security.VerifyIntermediateCA(rootCert, interCert)
// serverCert, serverCertPEM, _ := security.GenServerCert(interCert, interPrivKey)
// fmt.Println("serverCert\n", string(serverCertPEM))
// security.VerifyCertificateChain(rootCert, interCert, serverCert)
/*serverCert*/_, serverCertPEM, serverPrivKeyPEM := security.GenServerCert(rootCert, rootPrivKey)


//Generate the TLS keypair for the server.
// serverTlsCert, err := tls.X509KeyPair(serverCertPEM, serverPrivKeyPEM)
// if err != nil {
//   panic("Failed to generate X509KeyPair for the server.\n" + err.Error())
// }
// var serverTlsConf = &tls.Config{
//   Certificates: []tls.Certificate{serverTlsCert},
// }

//////////////////////////////////////

  /***
  fmt.Println("OS: " + misc.GetOS())
  if userName, err := misc.GetUsername(); err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Username: " + userName)
  }
  //
  if ok, err := misc.IsRoot(); err != nil {
    fmt.Println(err)
  } else if ok {
    fmt.Println("The current user is running as root.")
  } else {
    fmt.Println("The current user is not running as root.")
  }
  //
  if _, err := os.Stat("./files"); err != nil {
    if os.IsNotExist(err) {
      oldMask := syscall.Umask(0017)
      fmt.Printf("Default mask: %04o\nUsing mask: 0017\n", oldMask)
      err = os.Mkdir("./files", 0777)
      syscall.Umask(oldMask)
      if err != nil {
        panic(err)
      } else if err = sessions.AddUserToFile(USER_NAME, PASSWORD); err != nil {
        panic(err)
      } else if err = sessions.AddUserToFile("jct1", "pw1"); err != nil {
        panic(err)
      }
    } else {
      panic(err)
    }
  }
  //
  if err := sessions.ReadUsersFromFile(); err != nil {
    panic(err)
  }
  ***/
sessions.AddFromMemory(USER_NAME, PASSWORD)
//////////////////////////////////////

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
  h.mux["/logout"] = middlewares.ValidateSessions(wfpages.LogoutPage)
  h.mux["/welcome"] = middlewares.ValidateSessions(wfpages.WelcomePage)
  h.mux["/contact"] = middlewares.ValidateSessions(wfpages.ContactPage)
  h.mux["/about"] = middlewares.ValidateSessions(wfpages.AboutPage)
  h.mux["/finances"] = middlewares.ValidateSessions(wfpages.FinancesPage)
  h.mux["/fin/ordinaryannuity"] = middlewares.ValidateSessions(wfpages.OrdinaryAnnuityPage)
  h.mux["/fin/ordinaryannuity/interestrate"] = middlewares.ValidateSessions(wfoainterest.OaInterestRatePages)
  h.mux["/fin/ordinaryannuity/fv"] = middlewares.ValidateSessions(wfoafv.OaFvPages)
  h.mux["/fin/ordinaryannuity/pv"] = middlewares.ValidateSessions(wfoapv.OaPvPages)
  h.mux["/fin/ordinaryannuity/cp"] = middlewares.ValidateSessions(wfoacp.OaCpPages)
  h.mux["/fin/ordinaryannuity/epp"] = middlewares.ValidateSessions(wfoaepp.OaEppPages)
  h.mux["/fin/ordinaryannuity/ga"] = middlewares.ValidateSessions(wfoaga.OaGaPages)
  h.mux["/fin/ordinaryannuity/perpetuity"] = middlewares.ValidateSessions(wfoaperpetuity.OaPerpetuityPages)
  h.mux["/fin/annuitydue"] = middlewares.ValidateSessions(wfpages.AnnuityDuePage)
  h.mux["/fin/annuitydue/cp"] = middlewares.ValidateSessions(wfadcp.AdCpPages)
  h.mux["/fin/annuitydue/epp"] = middlewares.ValidateSessions(wfadepp.AdEppPages)
  h.mux["/fin/annuitydue/fv"] = middlewares.ValidateSessions(wfadfv.AdFvPages)
  h.mux["/fin/annuitydue/pv"] = middlewares.ValidateSessions(wfadpv.AdPvPages)
  h.mux["/fin/bonds"] = middlewares.ValidateSessions(wfbonds.BondsPages)
  h.mux["/fin/mortgage"] = middlewares.ValidateSessions(wfmortgage.MortgagePages)
  h.mux["/fin/simpleinterest"] = middlewares.ValidateSessions(wfpages.SimpleInterestPage)
  h.mux["/fin/simpleinterest/accurate"] = middlewares.ValidateSessions(wfsia.SimpleInterestAccuratePages)
  h.mux["/fin/simpleinterest/bankers"] = middlewares.ValidateSessions(wfsib.SimpleInterestBankersPages)
  h.mux["/fin/simpleinterest/ordinary"] = middlewares.ValidateSessions(wfsio.SimpleInterestOrdinaryPages)
  h.mux["/fin/miscellaneous"] = middlewares.ValidateSessions(wfmisc.MiscellaneousPages)
  /***
  Handlers for pprof.

  One way to enable the Go profiler (pprof) is to use the net/http/pprof package to serve the
  profiling data via HTTP. By using the blank import, it leads to a side effect that allows us to
  reach the pprof URL http://{url}:{port}/debug/pprof. Note that enabling pprof is safe even in
  production (https://go.dev/doc/diagnostics#profiling). The profiles that impact performance, such
  as CPU profiling, aren't enabled by default, nor do they run continuously; they are activated
  only for a specific period.

  To view all available profiles, open your browser and type the following address into the
  browser's address bar:
  http://{url}:{port}/debug/pprof/

  Please note you will need to have graphviz (https://graphviz.org/) installed for web
  visualizations. To install it in a Linux system, run the commands below:
  (If the universe repo is not enabled, enable it.)
  $ sudo add-apt-repository universe
  $ sudo apt update
  $ sudo apt install graphviz

  CPU Profiling
  -------------
  When it is activated, the application asks the OS to interrupt it every 10ms (default). When the
  application is interrupted, it suspends the current activity and transfers the execution to the
  profiler. The profiler collects execution statistics, and then it transfers execution back to the
  application.

  To active the CPU profiling, you access the debug/pprof/profile endpoint. Accessing this endpoint
  will execute CPU profiling for 30 seconds by default. For 30 seconds, the application is
  interrupted every 10ms.
  To write the output to a file, use the command below:
  $ curl http://{url}:{port}/debug/pprof/{prof1}?seconds={x} --output {filename}
    where {prof1} is trace or profile.
  $ curl http://{url}:{port}/debug/pprof/{prof2} --output {filename}
    where {prof2} is heap.
  To inspect a file, use the command below:
  $ go tool pprof {filename}
  To inspect the result using the graphical user interface, use the command below:
  $ go tool pprof -http=:{port1} {filename}
  To directly connect to the debug point, use the command below:
  $ go tool pprof http://{url}:{port}/debug/pprof/{prof1}?seconds={x}
  $ go tool pprof http://{url}:{port}/debug/pprof/{prof2}
  To inspect the result using the graphical user interface, use the command below:
  $ go tool pprof -http=:{port1} http://{url}:{port2}/debug/pprof/{prof1}?seconds={x}
  $ go tool pprof -http=:{port1} http://{url}:{port2}/debug/pprof/{prof2}
  ***/
  h.mux["/debug/pprof/"] = pprof.Index
  // h.mux["/debug/pprof/heap"] = pprof.Index
  h.mux["/debug/pprof/heap"] = pprof.Handler("heap").ServeHTTP
  h.mux["/debug/pprof/block"] = pprof.Handler("block").ServeHTTP
  h.mux["/debug/pprof/goroutine"] = pprof.Handler("goroutine").ServeHTTP
  h.mux["/debug/pprof/cmdline"] = pprof.Cmdline
  h.mux["/debug/pprof/profile"] = pprof.Profile
  h.mux["/debug/pprof/symbol"] = pprof.Symbol
  h.mux["/debug/pprof/trace"] = pprof.Trace
  commonMiddlewares := []middlewares.Middleware {
    middlewares.CorrelationId,
    //middlewares.ValidateSessions,
  }
  for idx, f := range h.mux {
    h.mux[idx] = middlewares.ChainMiddlewares(f, commonMiddlewares)
  }
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
    ReadHeaderTimeout: 3 * time.Second,
    /***
    It specifies the maximum amount of time to read the entire request.
    ReadTimeout = ReadHeaderTimeout + TimeoutHandler + Extra time
    **/
    ReadTimeout: 6 * time.Second,
    /***
    If a handler fails to respond on time, the server will reply with "503 Service Unavailable" and
    the specified message; the context passed to the handler will be canceled.
    Note: The http.Server.WriteTimeout is not necessary since http.TimeoutHandler is being used.
    ***/
    Handler: http.TimeoutHandler(&h, 30 * time.Second, "Request timeout."),
    /***
    It configures the maximum amount of time for the next request when keep-alives are enabled.
    Note that if http.Server.IdleTimeout isn't set, the value of http.Server.ReadTimeout is used
    for the idle timeout. If neither is set, there won't be any timeouts, and connections will
    remain open until they are closed by clients.
    ***/
    IdleTimeout: 120 * time.Second,
    MaxHeaderBytes: 1 << 20,  //1 MB.
    TLSConfig: &tls.Config{
      MinVersion: tls.VersionTLS13,
      /***
      When using version 1.3 this isn't configurable. Because you want to use the most up to date
      version, keep it empty, which results in a default list of ciphersuites to be used with a
      preference order based on hardware performance.
      ***/
      CipherSuites: nil,
      /***
..........      To control the server's preferred ciphersuite to use as provided by the CipherSuites............, when false it will select the client’s preferred ciphersuite. Setting this will ensure that safer and faster ciphersuites are used. [@valsorda2016a]
      ***/
      PreferServerCipherSuites: true,
      CurvePreferences: []tls.CurveID{
        tls.CurveP521,
        tls.CurveP384,
        tls.CurveP256,
      },
//      Certificates: []tls.Certificate{serverTlsCert},

			GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
        // Always get latest localhost.crt and localhost.key 
        // ex: keeping certificates file somewhere in global location where created certificates updated and this closure function can refer that
cert, err := tls.X509KeyPair(serverCertPEM, serverPrivKeyPEM)
if err != nil {
return nil, err
}
return &cert, nil
},


},
TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),

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
  // err := (*server).ListenAndServe()
  err := (*server).ListenAndServeTLS("", "")
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


///////////////////
/***
server := &http.Server{
  Addr:         ":" + *port,
  ReadTimeout:  5 * time.Minute, // 5 min to allow for delays when 'curl' on OSx prompts for username/password
  WriteTimeout: 10 * time.Second,
  TLSConfig:    getTLSConfig(*host, *caCert, tls.ClientAuthType(*certOpt)),
}
****/
//https://youngkin.github.io/post/gohttpsclientserver/


//Redirect http requests to https.
func redirectHttpToHttps(res http.ResponseWriter, req *http.Request) {
  url := *req.URL
  url.Scheme = "https"
  url.Host = req.Host
  http.Redirect(res, req, url.String(), http.StatusMovedPermanently)
}


/***

openssl req -x509 -newkey rsa:4096 -keyout privKey.pem -out pubCert.pem -days 365

Before we can use this code we need a certificate. Run the following command to generate a private key file and a certificate signing request.
openssl req -new -newkey rsa:2048 -nodes -x509 -days 365 -out pubCert.pem -keyout privKey.pem \
  -subj "/C=US/ST=Texas/L=Houston/O=Company Name/OU=Dept A/CN=localhost"




***/