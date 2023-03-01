// Test the app.
package main

import (
	"finance/finances"
	// // "app/umath"
	"fmt"
	// "math"
)

// type S struct{}
// func (s *S) addr() { fmt.Printf("%p\n", s) }

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
}
