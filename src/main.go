package main

import (
  "finance/webfinances"
  "fmt"
  "log"
  "net/http"
  "os"
  "time"
)

// type S struct{}
// func (s *S) addr() { fmt.Printf("%p\n", s) }

/***
How to kill a process using a port on localhost (Windows).
C:\> netstat -ano | findstr :<port>
C:\> taskkill /PID <PID> /F

or

C:\> npx kill-port <port>
***/
func main() {

  // var a, b S
  // a.addr()
  // b.addr()

  // if len(os.Args) > 2 && strings.EqualFold(os.Args[1], "ordinary") {
  //   if os.Args[2] == "interest" && len(os.Args) == 8 {
  //     _, err := strconv.Atoi(os.Args[4])
  //     if err != nil {
  //       panic(err)
  //     }
  //     var si finances.SimpleInterest
  //     var interest = (&si).OrdinaryInterest(100, 0.04, finances.Monthly, 1, finances.Months)
  //     fmt.Printf("interest = $%.2f (Ordinary Interest)\n", interest)
  //   }
  // }

  
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
  // fmt.Println("eps = ", math.Nextafter(1.0, 2.0) - 1.0)


  dns, exists := os.LookupEnv("DNS")
  if !exists {
    fmt.Println("Missing environment parameter: DNS=localhost:8000 go run main.go")
    return
  }
  fmt.Printf("%s - DNS=%s\n", time.Now().UTC().Format(time.RFC3339Nano), dns)
  var w webfinances.Annuities
  /***
  net/http provides ServeMux, a request multiplexer, to simplify the association between URLs and
  handlers. A ServeMux aggregates a collection of http.Handlers into a single http.Handler.
  Different types satisfying the same interface are substitutable: the web server can dispatch
  requests to any http.Handler, regardless of which concrete type is behind it.
  ***/
  // var mux = http.NewServeMux()
  // http.HandleFunc("/bonds", db.bonds)
  http.HandleFunc("/annuities/AverageRateOfReturn", w.AverageRateOfReturn)
  fmt.Printf("%s - Starting the server...\n", time.Now().UTC().Format(time.RFC3339Nano))
  /***
  ListenAndServe runs forever, or until the server fails (or fails to start) with an error,
  always non-nil, which it returns.

  The web server invokes each handler in a new goroutine, so handlers must take precautions such as
  locking when accessing variables that other goroutines, including other requests to the same
  handler, may be accessing.
  ***/
  log.Fatal(http.ListenAndServe(dns, nil)) //DefaultServeMux
}

// type database map[string]int

// func (db database) bonds(res http.ResponseWriter, req *http.Request) () {
//   fmt.Fprintf(res, "URL.Path = %q\n", req.URL.Path)
// }

// func (db database) annuities(res http.ResponseWriter, req *http.Request) () {
//   params := req.URL.Query()
//   if len(params) != 1 {
//     fmt.Fprintf(res, "Parameters required = 1; parameters provided = %d\n", len(params))
//     return
//   }
//   const length int = 16
//   var returns []float64 = nil
//   //Iterate over all the query parameters.
//   for _, v := range params {
//     returns = make([]float64, 0, length)
//     for idx := 0; idx < len(v); idx++ {
//       f, err := strconv.ParseFloat(v[idx], 64)
//       if err != nil {
//         fmt.Fprintf(res, "'%s' is not a floating number.\n%s", v[idx], v)
//         return
//       }
//       returns = append(returns, f)
//     }
//   }
//   var a finances.Annuities
//   var gmr = a.AverageRateOfReturn(returns) * 100.0
//   fmt.Fprintf(res, "AverageRateOfReturn = %.3f\n", gmr)
// }



