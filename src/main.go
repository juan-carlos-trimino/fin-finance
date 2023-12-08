//go:build !windows
// +build !windows

// HTTP server.
package main

/***
*****************************************************************************************
*** To run the app in a K8s environment, do NOT set the environment variable NOT_K8S. ***
*****************************************************************************************

To run the app as a standalone HTTP/HTTPS servers:
Compile and run the app.
$ go build -o finance && HTTPS=true ./finance

To change the PORT.
$ go build -o finance && HTTPS=true PORT=8181 ./finance

Compile and run the app in the background.
$ go build -o finance && HTTPS=true ./finance &
$ go build -o finance && HTTP=true HTTPS=true LE_CERT=false ./finance

Force rebuilding of packages.
$ go build -o finance -a && HTTPS=true ./finance

Compile and run (in the background) at the same time.
$ HTTPS=true go run main.go &

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
	"errors"
	"finance/middlewares"
	"finance/misc"
	"finance/security"
	"finance/sessions"
	"finance/webfinances"
	"fmt"
  "golang.org/x/crypto/acme/autocert"
	"net"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof" //Blank import of pprof.
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (  //Environment variables.
  K8S bool = false
  HTTPS bool = false
  HTTP bool = true
  LE_CERT = false
  MAX_RETRIES int = 10
  SHUTDOWN_TIMEOUT int = 15
  HTTPS_PORT string = "8443"
  HTTP_PORT string = "8080"
  SVC_NAME string
  APP_NAME_VER string
  SERVER string = "localhost"
  USER_NAME string = "a"
  PASSWORD string = "a"
)

const (
  dataDir string = "./dataDir"
  users string = "user.txt"
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
  ev, exists = os.LookupEnv("HTTPS")
  if exists {
    b, err := strconv.ParseBool(ev)
    if err == nil {
      HTTPS = b
    } else {
      fmt.Printf("'%s' is not a boolean.\n", ev)
    }
  }
  //
  ev, exists = os.LookupEnv("HTTP")
  if exists {
    b, err := strconv.ParseBool(ev)
    if err == nil {
      HTTP = b
    } else {
      fmt.Printf("'%s' is not a boolean.\n", ev)
    }
  }
  //
  if HTTP == false && HTTPS == false {
    fmt.Println("You can run only HTTP (default), only HTTPS (set environment variables to:" +
                " HTTP=false and HTTPS=true), or both (set environment variable to: HTTPS=true).")
    return
  }
  ev, exists = os.LookupEnv("LE_CERT")
  if exists {
    b, err := strconv.ParseBool(ev)
    if err == nil {
      LE_CERT = b
    } else {
      fmt.Printf("'%s' is not a boolean.\n", ev)
    }
  }
  ev, exists = os.LookupEnv("K8S")
  if exists {
    b, err := strconv.ParseBool(ev)
    if err == nil {
      K8S = b
    } else {
      fmt.Printf("'%s' is not a boolean.\n", ev)
    }
    //
    if K8S {
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
  }
  //
  if HTTP {
    ev, exists = os.LookupEnv("HTTP_PORT")
    if exists {
      HTTP_PORT = ev
    }
    fmt.Printf("Using HTTP PORT: %s\n", HTTP_PORT)
  }
  //
  if HTTPS {
    ev, exists = os.LookupEnv("HTTPS_PORT")
    if exists {
      HTTPS_PORT = ev
    }
    fmt.Printf("Using HTTPS PORT: %s\n", HTTPS_PORT)
  }
  //SERVER += ":" + PORT
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
  readUsers(dataDir, users)
  /***
  When Shutdown is called, Serve, ListenAndServe, and ListenAndServeTLS immediately return
  ErrServerClosed. Make sure the program doesn't exit and waits instead for Shutdown to return.
  ***/
  var wg sync.WaitGroup
  var httpServer *http.Server
  if HTTP {
    httpServer = makeServer(HTTP_PORT, makeHandlers())
  }
  //https://pkg.go.dev/golang.org/x/crypto/acme/autocert
  var certMan autocert.Manager
  if LE_CERT {
    certMan = autocert.Manager{
      //It always returns true to indicate acceptance of the CA's Terms of Service during account
      //registration.
      Prompt: autocert.AcceptTOS,
      HostPolicy: autocert.HostWhitelist("trimino.xyz", "www.trimino.xyz"), //Domain names.
      Cache: autocert.DirCache(dataDir), //Folder for storing certificates.
    }
  }
  //
  if HTTPS {
    wg.Add(1)
    /***
    A channel is a communication mechanism that lets one goroutine send values to another
    goroutine. Each channel is a conduit for values of a particular type, called the channel's
    element type.

    As with maps, a channel is a reference to the data structure created by make. When we copy a
    channel or pass one as an argument to a function, we are copying a reference, so caller and
    callee refer to the same data structure. As with other reference types, the zero value of a
    channel is nil.
    ***/
    //Buffered channel capacity 1; notifier will not block.
    var signalChan2 chan os.Signal = make(chan os.Signal, 1)
    if HTTP {
      if LE_CERT {
        //https://pkg.go.dev/golang.org/x/crypto/acme/autocert#Manager.HTTPHandler
        httpServer.Handler = certMan.HTTPHandler(nil)
      } else {
        httpServer.Handler = makeHttpToHttpsRedirectHandler(HTTPS_PORT)
      }
    }
    signalChan2 = make(chan os.Signal, 1) //Buffered channel capacity 1; notifier will not block.
    go func() {
      httpsServer := makeServer(HTTPS_PORT, makeHandlers())
      if LE_CERT {
        httpsServer.TLSConfig = &tls.Config{
          MinVersion: tls.VersionTLS13,
          CipherSuites: nil,
          PreferServerCipherSuites: true,
          CurvePreferences: []tls.CurveID{
            tls.CurveP256,
            tls.X25519,
          },
          GetCertificate: certMan.GetCertificate,
        }
      } else {
        httpsServer.TLSConfig = makeTlsConfig()
      }
      go waitForServer(httpsServer, signalChan2, &wg)
      fmt.Printf("%s - Starting the server at port %s...\n", m.DTF(), httpsServer.Addr)
      //Because the paths of the key and cert were set in the TLSConfig field, set the certFile and
      //keyFile arguments to empty strings.
      err := (*httpsServer).ListenAndServeTLS("", "")
      if errors.Is(err, http.ErrServerClosed) {
        fmt.Printf("%s - Server has been closed at port %s.\n", m.DTF(), httpsServer.Addr)
      } else if err != nil {
        fmt.Printf("%s - Server error: %+v\n", m.DTF(), err)
        signalChan2 <- syscall.SIGINT //Let the goroutine finish.
      }
    }()
  }
  //
  if HTTP {
    wg.Add(1)
    signalChan1 := make(chan os.Signal, 1)
    /***
    When Shutdown is called, Serve, ListenAndServe, and ListenAndServeTLS immediately return
    ErrServerClosed. Make sure the program doesn't exit and waits instead for Shutdown to return.
    ***/
    go waitForServer(httpServer, signalChan1, &wg)
    fmt.Printf("%s - Starting the server at port %s...\n", m.DTF(), httpServer.Addr)
    /***
    ListenAndServe runs forever, or until the server fails (or fails to start) with an error,
    always non-nil, which it returns.

    The web server invokes each handler in a new goroutine, so handlers must take precautions such
    as locking when accessing variables that other goroutines, including other requests to the same
    handler, may be accessing.
    ***/
    err := (*httpServer).ListenAndServe()
    if errors.Is(err, http.ErrServerClosed) {
      fmt.Printf("%s - Server has been closed at port %s.\n", m.DTF(), httpServer.Addr)
    } else if err != nil {
      fmt.Printf("%s - Server error: %+v\n", m.DTF(), err)
      signalChan1 <- syscall.SIGINT //Let the goroutine finish.
    }
  }
  wg.Wait() //Block until shutdown is done.
}

func faviconHandler(res http.ResponseWriter, req *http.Request) {
  http.NotFound(res, req)  //404 - page not found.
}

func makeHandlers() *handlers {
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
  h := &handlers{}
  h.mux = make(map[string]http.HandlerFunc, 64)
  //h.mux["/readiness"] =
  //func (res http.ResponseWriter, req *http.Request){
  //  fmt.Printf("\naaaaaaServer not ready. %s\n", SERVER)
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
    //res.WriteHeader(http.StatusOK)
  //}
  //Serve static files; i.e., the server will serve them as they are, without processing it first.
  h.mux["/public/css/home.css"] = wfpages.PublicHomeFile
  h.mux["/public/js/getParams.js"] = wfpages.PublicGetParamsFile
  h.mux["/public/js/mortgage.js"] = wfpages.PublicMortgageFile
  h.mux["/public/js/OaInterestRate.js"] = wfpages.PublicOaInterestRateFile
  h.mux["/public/js/OaPresentValue.js"] = wfpages.PublicOaPresentValueFile
  h.mux["/public/js/OaFutureValue.js"] = wfpages.PublicOaFutureValueFile
  h.mux["/public/js/OaCompoundingPeriods.js"] = wfpages.PublicOaCompoundingPeriodsFile
  h.mux["/public/js/OaEqualPeriodicPayments.js"] = wfpages.PublicOaEqualPeriodicPaymentsFile
  h.mux["/public/js/OaGrowingAnnuity.js"] = wfpages.PublicOaGrowingAnnuityFile
  h.mux["/public/js/OaPerpetuity.js"] = wfpages.PublicOaPerpetuityFile
  h.mux["/public/js/AdCompoundingPeriods.js"] = wfpages.PublicAdCompoundingPeriodsFile
  h.mux["/public/js/AdEqualPeriodicPayments.js"] = wfpages.PublicAdEqualPeriodicPaymentsFile
  h.mux["/public/js/AdFutureValue.js"] = wfpages.PublicAdFutureValueFile
  h.mux["/public/js/AdPresentValue.js"] = wfpages.PublicAdPresentValueFile
  h.mux["/public/js/bonds.js"] = wfpages.PublicBondsFile
  h.mux["/public/js/miscellaneous.js"] = wfpages.PublicMiscellaneousFile
  h.mux["/public/js/SimpleInterestAccurate.js"] = wfpages.PublicSimpleInterestAccurateFile
  h.mux["/public/js/SimpleInterestBankers.js"] = wfpages.PublicSimpleInterestBankersFile
  h.mux["/public/js/SimpleInterestOrdinary.js"] = wfpages.PublicSimpleInterestOrdinaryFile
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
  commonMiddlewares := []middlewares.Middleware{
    middlewares.SecurityHeaders,
    middlewares.CorrelationId,
    //middlewares.ValidateSessions,
  }
  //
  for idx, f := range h.mux{
    h.mux[idx] = middlewares.ChainMiddlewares(f, commonMiddlewares)
  }
  return h
}

func makeHttpToHttpsRedirectHandler(port string) *handlers {
  /***
  The Go web server will route requests to different functions depending on the requested URL.
  ***/
  h := &handlers{}
  h.mux = make(map[string]http.HandlerFunc, 1)
  h.mux["/"] = func(res http.ResponseWriter, req *http.Request) {
    host, _, _ := net.SplitHostPort(req.Host)
    u := req.URL
    u.Host = net.JoinHostPort(host, port)
    u.Scheme="https"
    fmt.Printf("%s - Redirecting to %s\n", m.DTF(), u.String())
    http.Redirect(res, req, u.String(), http.StatusMovedPermanently)
  }
  return h
}

func readUsers(dir, filename string) {
  if dir != "" && filename != "" {
    if _, err := os.Stat(dir); err != nil {
      if os.IsNotExist(err) {
        oldMask := syscall.Umask(0017)
        fmt.Printf("Default mask: %04o\nUsing mask: 0017\n", oldMask)
        err = os.Mkdir(dir, 0777)
        syscall.Umask(oldMask)
        if err != nil {
          panic("Cannot create dir.\n" + err.Error())
        } else if err = sessions.AddUserToFile(dir + "/" + filename, USER_NAME, PASSWORD);
                  err != nil {
          panic(err)
        } else if err = sessions.AddUserToFile(dir + "/" + filename, "jct1", "pw1"); err != nil {
          panic(err)
        }
      } else {
        panic(err)
      }
    }
    //
    if err := sessions.ReadUsersFromFile(dir + "/" + filename); err != nil {
      panic(err)
    }
  } else {
    sessions.AddFromMemory(USER_NAME, PASSWORD)
  }
}



func makeTlsConfig() *tls.Config {

// see Certificate structure at
// http://golang.org/pkg/crypto/x509/#Certificate
// see http://golang.org/pkg/crypto/x509/#KeyUsage

//Addr:      ":4443",
//server would listen on IP address 0.0.0.0 and TCP port 4443.

  rootCert, rootCertPEM, rootPrivKey := security.GenRootCA()
  _, serverCertPEM, serverPrivKeyPEM := security.GenServerCert(rootCert, rootPrivKey)


// fmt.Println(string(serverCertPEM))
// fmt.Println(string(serverPrivKeyPEM))

  //Generate the key/pair for the server.
  serverTlsCert, err := tls.X509KeyPair(serverCertPEM, serverPrivKeyPEM)
  if err != nil {
    panic("Failed to generate X509KeyPair for the server.\n" + err.Error())
  }


// fmt.Println(string(serverTlsCert.Certificate[0]))

  rootCAs := security.SystemCertPool()
  if ok := rootCAs.AppendCertsFromPEM(rootCertPEM); !ok {
    panic("Failed to append CA's certificate.\n" + err.Error())
  }
  tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS13,
    /***
    When using version 1.3 this isn't configurable. Because you want to use the most up to date
    version, keep it empty, which results in a default list of ciphersuites to be used with a
    preference order based on hardware performance.
    ***/
    CipherSuites: nil,
    /***
    Control the server's preferred ciphersuite to use as provided by the CipherSuites. When
    false, it will select the client's preferred ciphersuite. Setting this will ensure that safer
    and faster ciphersuites are used.
    ***/
    PreferServerCipherSuites: true,
    /***
    CurvePreferences contains the elliptic curves that will be used in an ECDHE handshake;
    however, without tls.CurveP384 because a client using tls.CurveP384 would cause up to a
    second of CPU to be consumed on the server.
    ***/
    CurvePreferences: []tls.CurveID{
      tls.CurveP521,
      // tls.CurveP384,
      tls.CurveP256,
    },
    Certificates: []tls.Certificate{serverTlsCert},
    RootCAs: rootCAs,
    // InsecureSkipVerify: true,
  }
  return tlsConfig
}

func waitForServer(server *http.Server, signalChan chan os.Signal, wg *sync.WaitGroup) {
  fmt.Printf("%s - Waiting for notification to shut down the server at %s.\n", m.DTF(), server.Addr)
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
    wg.Done() //Shutdown is done; let the main goroutine terminate.
  }()
  //https://pkg.go.dev/net/http#Server.Shutdown
  if err := server.Shutdown(ctx); err != nil {
    fmt.Printf("%s - Server shutdown failed: %+v\n", m.DTF(), err)  //https://pkg.go.dev/fmt
  }
}

func makeServer(port string, h *handlers) *http.Server {
  server := &http.Server{  //https://pkg.go.dev/net/http#ServeMux
    /***
    By not specifying an IP address before the colon, the server will listen on every IP address
    associated with the computer, and it will listen on port PORT.
    ***/
    Addr: ":" + port,
    /***
    Set timeouts so that a slow or malicious client doesn't hold resources forever.

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
    Handler: http.TimeoutHandler(h, 30 * time.Second, "Request timeout."),
    /***
    It configures the maximum amount of time for the next request when keep-alives are enabled.
    Note that if http.Server.IdleTimeout isn't set, the value of http.Server.ReadTimeout is used
    for the idle timeout. If neither is set, there won't be any timeouts, and connections will
    remain open until they are closed by clients.
    ***/
    IdleTimeout: 120 * time.Second,
    MaxHeaderBytes: 1 << 20,  //1 MB.
    /***
    Setting TLSNextProto to an empty map will disable HTTP/2 for this server. If you want to enable
    HTTP/2, set it to nil or remove the field. Since Go 1.6, it is enabled by default.
    ***/
    TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
  }
  return server
}
