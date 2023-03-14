//HTTP server.
package main

import (
  "context"
  "errors"
  "finance/webfinances"
  "fmt"
  "net/http"
  "os"
  "os/signal"
  "strconv"
  "syscall"
  "time"
)

//Environment variables.
var MAX_RETRIES int = 10
var PORT string = "8001"
var SVC_NAME string
var APP_NAME_VER string
var SHUTDOWN_TIMEOUT int = 15

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
/////////////////////////////////////
  /***
  var m finances.Mortgage
  var payment, totalCost, totalInterest = (&m).CostOfMortgage(300000.00, 2.74 / 100.0, 'm', 15.0, 'y')
  fmt.Printf("Payment = $%.2f Total cost = $%.2f Total interest = $%.2f\n", payment, totalCost, totalInterest)
  var bir = (&m).MortgageHeloc(200000, 0.065, 100000, 0.105)
  fmt.Printf("Blended Interest Rate = %.2f%%\n", bir)
  var table = m.AmortizationTable(300000.00, 0.03375, 'm', 30.0, 'y')
  fmt.Printf("payment = $%.2f total cost = $%.2f total interest = $%.2f\n", table.Payment, table.TotalCost, table.TotalInterest)
  for i, v := range table.Rows {
    fmt.Printf("pmtNumber = %d payment = $%.2f pmtPrincipal = $%.2f pmtInterest = $%.2f balance = $%.2f\n", i + 1, v.Payment, v.PmtPrincipal, v.PmtInterest, v.Balance)
  }
  ***/
  //////////////////////  // fmt.Println("eps = ", math.Nextafter(1.0, 2.0) - 1.0)
  //http://localhost:8001/annuities/AverageRateOfReturn?ret=5.0&ret=-3.0&ret=12.0&ret=10


  var exists bool = false
  // SVC_NAME, exists = os.LookupEnv("SVC_NAME")
  // if !exists {
  //   fmt.Println("Missing environment parameter: SVC_NAME")
  //   return
  // }
  // APP_NAME_VER, exists = os.LookupEnv("APP_NAME_VER")
  // if !exists {
  //   fmt.Println("Missing environment parameter: APP_NAME_VER")
  //   return
  // }
  _, exists = os.LookupEnv("PORT")
  if exists {
    PORT = os.Getenv("PORT")
  }
  fmt.Printf("Using PORT: %s\n", PORT)
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
  var a webfinances.Annuities
  var h handlers = handlers{}
  h.mux = make(map[string]func(http.ResponseWriter, *http.Request), 16)
  h.mux["/annuities/AverageRateOfReturn"] = a.AverageRateOfReturn
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
