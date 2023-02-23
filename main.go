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
	var si finances.SimpleInterest
	var interest = (&si).OrdinaryInterest(100, 0.04, finances.Monthly, 1, finances.Months)
	fmt.Printf("interest = $%.2f (Ordinary Interest)\n", interest)
	rate := si.OrdinaryRate(100, interest, finances.Monthly, 1, finances.Months)
	fmt.Printf("rate = %.2f%% (Ordinary Interest)\n", rate*100.0)
	principal := si.OrdinaryPrincipal(interest, rate, finances.Monthly, 1, finances.Months)
	fmt.Printf("principal = $%.2f (Ordinary Interest)\n", principal)
	time := si.OrdinaryTime(principal, interest, rate, finances.Monthly, finances.Months)
	fmt.Printf("time period = %.2f (Ordinary Interest)\n\n", time)
	//
	interest = (&si).OrdinaryInterest(100, 0.04, finances.Annually, 1, finances.Months)
	fmt.Printf("interest = $%.2f (Ordinary Interest)\n", interest)
	rate = si.OrdinaryRate(100, interest, finances.Annually, 1, finances.Months)
	fmt.Printf("rate = %.2f%% (Ordinary Interest)\n", rate*100.0)
	principal = si.OrdinaryPrincipal(interest, rate, finances.Annually, 1, finances.Months)
	fmt.Printf("principal = $%.2f (Ordinary Interest)\n", principal)
	time = si.OrdinaryTime(principal, interest, rate, finances.Annually, finances.Months)
	fmt.Printf("time period = %.2f (Ordinary Interest)\n\n", time)
	//
	interest = (&si).OrdinaryInterest(10000, 0.09, finances.Annually, 153, finances.Days)
	fmt.Printf("interest = $%.2f (Ordinary Interest)\n", interest)
	rate = si.OrdinaryRate(10000, interest, finances.Annually, 153, finances.Days)
	fmt.Printf("rate = %.2f%% (Ordinary Interest)\n", rate*100.0)
	principal = si.OrdinaryPrincipal(interest, rate, finances.Annually, 153, finances.Days)
	fmt.Printf("principal = $%.2f (Ordinary Interest)\n", principal)
	time = si.OrdinaryTime(principal, interest, rate, finances.Annually, finances.Years) //Days
	fmt.Printf("time period = %.4f (Ordinary Interest)\n\n", time)
	***/
	/***
	  var si finances.SimpleInterest
	  var interest = (&si).BankersInterest(100, 0.04, finances.Monthly, 1, finances.Months)
	  fmt.Printf("interest = $%.2f (Banker's Interest)\n", interest)
	  rate := si.BankersRate(100, interest, finances.Monthly, 1, finances.Months)
	  fmt.Printf("rate = %.2f%% (Banker's Interest)\n", rate * 100.0)
	  principal := si.BankersPrincipal(interest, rate, finances.Monthly, 1, finances.Months)
	  fmt.Printf("principal = $%.2f (Banker's Interest)\n", principal)
	  time := si.BankersTime(principal, interest, rate, finances.Monthly, finances.Months)
	  fmt.Printf("time period = %.2f (Banker's Interest)\n\n", time)
	  //
	  interest = (&si).BankersInterest(100, 0.04, finances.Annually, 1, finances.Months)
	  fmt.Printf("interest = $%.2f (Banker's Interest)\n", interest)
	  rate = si.BankersRate(100, interest, finances.Annually, 1, finances.Months)
	  fmt.Printf("rate = %.2f%% (Banker's Interest)\n", rate * 100.0)
	  principal = si.BankersPrincipal(interest, rate, finances.Annually, 1, finances.Months)
	  fmt.Printf("principal = $%.2f (Banker's Interest)\n", principal)
	  time = si.BankersTime(principal, interest, rate, finances.Annually, finances.Months)
	  fmt.Printf("time period = %.2f (Banker's Interest)\n\n", time)
	  //
	  interest = (&si).BankersInterest(10000, 0.09, finances.Annually, 153, finances.Days)
	  fmt.Printf("interest = $%.2f (Banker's Interest)\n", interest)
	  rate = si.BankersRate(10000, interest, finances.Annually, 153, finances.Days)
	  fmt.Printf("rate = %.2f%% (Banker's Interest)\n", rate * 100.0)
	  principal = si.BankersPrincipal(interest, rate, finances.Annually, 153, finances.Days)
	  fmt.Printf("principal = $%.2f (Banker's Interest)\n", principal)
	  time = si.BankersTime(principal, interest, rate, finances.Annually, finances.Years)  //Days
	  fmt.Printf("time period = %.4f (Banker's Interest)\n\n", time)
	  ***/
	/***
	  var si finances.SimpleInterest
	  var interest = (&si).AccurateInterest(100, 0.04, finances.Monthly, 1, finances.Months)
	  fmt.Printf("interest = $%.2f (Accurate Interest)\n", interest)
	  rate := si.AccurateRate(100, interest, finances.Monthly, 1, finances.Months)
	  fmt.Printf("rate = %.2f%% (Accurate Interest)\n", rate * 100.0)
	  principal := si.AccuratePrincipal(interest, rate, finances.Monthly, 1, finances.Months)
	  fmt.Printf("principal = $%.2f (Accurate Interest)\n", principal)
	  time := si.AccurateTime(principal, interest, rate, finances.Monthly, finances.Months)
	  fmt.Printf("time period = %.2f (Accurate Interest)\n\n", time)
	  //
	  interest = (&si).AccurateInterest(100, 0.04, finances.Annually, 1, finances.Months)
	  fmt.Printf("interest = $%.2f (Accurate Interest)\n", interest)
	  rate = si.AccurateRate(100, interest, finances.Annually, 1, finances.Months)
	  fmt.Printf("rate = %.2f%% (Accurate Interest)\n", rate * 100.0)
	  principal = si.AccuratePrincipal(interest, rate, finances.Annually, 1, finances.Months)
	  fmt.Printf("principal = $%.2f (Accurate Interest)\n", principal)
	  time = si.AccurateTime(principal, interest, rate, finances.Annually, finances.Months)
	  fmt.Printf("time period = %.2f (Accurate Interest)\n\n", time)
	  //
	  interest = (&si).AccurateInterest(10000, 0.09, finances.Annually, 153, finances.Days)
	  fmt.Printf("interest = $%.2f (Accurate Interest)\n", interest)
	  rate = si.AccurateRate(10000, interest, finances.Annually, 153, finances.Days)
	  fmt.Printf("rate = %.2f%% (Accurate Interest)\n", rate * 100.0)
	  principal = si.AccuratePrincipal(interest, rate, finances.Annually, 153, finances.Days)
	  fmt.Printf("principal = $%.2f (Accurate Interest)\n", principal)
	  time = si.AccurateTime(principal, interest, rate, finances.Annually, finances.Days)  //Years
	  fmt.Printf("time period = %.4f (Accurate Interest)\n\n", time)
	  ***/
  /***
  var m finances.Miscellaneous
  var real = (&m).RealInterestRate(0.045, 0.065)
  fmt.Printf("Real Interest Rate = %.2f%%\n", real * 100)
  fmt.Printf("Growth/Decay Of Funds = %.2f\n", m.GrowthDecayOfFunds(.15, finances.Annually, 2.0))
  var ear float64 = m.NominalToEffectiveAnnualRate(0.12, finances.Monthly)
  fmt.Printf("Effective Annual Rate = %.6f%%\n", ear * 100.0)
  fmt.Printf("Nomial Rate = %.6f%%\n", m.EffectiveAnnualToNominalRate(ear, finances.Monthly) * 100.0)
  ***/
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
  /***
  var a finances.Annuities
	var ear = a.NominalToEAR(0.12, finances.Monthly)
  fmt.Printf("ear = %.2f%%\n", ear * 100.0)
  fmt.Printf("nr = %.2f%%\n", a.EARToNominal(ear, finances.Monthly) * 100.0)
  //
  v := []float64{5.0, -3.0, 12.0, 10.0}
  fmt.Printf("gmr = %.2f%%\n", a.AverageRateOfReturn(v) * 100.0)
  v1 := []float64{2.0, 8.0, -1.0, 10.0}
  fmt.Printf("gmr = %.2f%%\n", a.AverageRateOfReturn(v1) * 100.0)
	***/
  /***
  var a finances.Annuities
  var cp = finances.Monthly
  var i = a.O_Interest_PV_PMT(24000.0, 500.0, 60.0, cp, 1.0, 31.0, 1e-6)
  fmt.Printf("i(0.7628634%% per month) = %.8f%%\n", i * 100)
  fmt.Printf("i(9.154323%% per year) = %.8f%%\n", i * 100 * float64(cp))
  //
  i = a.O_Interest_PV_PMT(11200, 291, 48, cp, 4.0, 12.0, 1e-6)
  fmt.Printf("i(0.94007411%% per month) = %.8f%%\n", i * 100)
  fmt.Printf("i(11.28%% per year) = %.8f%%\n", i * 100 * float64(cp))
  //
  cp = finances.Annually
  i = a.O_Interest_PV_PMT(50000, 13500, 5, cp, 10.0, 15.0, 1e-6)
  fmt.Printf("i(10.91616%% per year) = %.8f%%\n", i * 100)
  fmt.Printf("i(10.91616%% per year) = %.8f%%\n", i * 100 * float64(cp))
  ***/
	var b finances.Bonds
	var cashFlow = b.CashFlow(1000, 3, finances.SemiAnnually, 5, finances.Years)
	fmt.Printf("Cash Flow\n")
  for idx := range cashFlow {
    fmt.Printf("Payment[%d] = $%.2f\n", idx, cashFlow[idx])
  }
  //FV = $1,000.00
  //coupon rate = 10%; compounding period semiannually
  //t = 10-year
  //current interest = 10%; compounding period semiannually
  cashFlow = b.CashFlow(1000.00, 10.0, b.GetCompoundingPeriod('s', false), 10.0, b.GetTimePeriod('y', false))
  var price = b.CurrentPrice(cashFlow, 10.0, b.GetCompoundingPeriod('s', false))
  fmt.Printf("price($1,000.00 (par)) = $%.2f\n", price)
  //FV = $1,000.00
  //coupon rate = 10%; compounding period semiannually
  //t = 10-year
  //current interest = 11%; compounding period semiannually
  cashFlow = b.CashFlow(1000.00, 10.0, b.GetCompoundingPeriod('s', false), 10.0, b.GetTimePeriod('y', false))
  price = b.CurrentPrice(cashFlow, 11.0, b.GetCompoundingPeriod('s', false))
  fmt.Printf("price($940.25 (discount)) = $%.2f\n", price)
  //FV = $1,000.00
  //coupon rate = 10%; compounding period semiannually
  //t = 10-year
  //current interest = 9%; compounding period semiannually
  cashFlow = b.CashFlow(1000.00, 10.0, b.GetCompoundingPeriod('s', false), 10.0, b.GetTimePeriod('y', false))
  price = b.CurrentPrice(cashFlow, 9.00, b.GetCompoundingPeriod('s', false))
  fmt.Printf("price($1,065.04 (premium)) = $%.2f\n", price)
  //FV = $1,000.00
  //coupon rate = 10%; compounding period semiannually
  //t = 10-year
  //current interest = 1%; compounding period semiannually
  //price = $1,854.43 (premium)
  cashFlow = b.CashFlow(1000.00, 10.0, b.GetCompoundingPeriod('s', false), 10.0, b.GetTimePeriod('y', false))
  price = b.CurrentPrice(cashFlow, 1.00, b.GetCompoundingPeriod('s', false))
  fmt.Printf("price($1,854.43 (premium)) = $%.2f\n", price)
  //FV = $1,000.00
  //coupon rate = 10%; compounding period semiannually
  //t = 10-year
  //current interest = 19%; compounding period semiannually
  cashFlow = b.CashFlow(1000.00, 10.0, b.GetCompoundingPeriod('s', false), 10.0, b.GetTimePeriod('y', false))
  price = b.CurrentPrice(cashFlow, 19.00, b.GetCompoundingPeriod('s', false))
  fmt.Printf("price($603.44 (discount)) = $%.2f\n", price)
  //FV = $1,000.00
  //coupon rate = 10%; compounding period semiannually
  //t = 20-semiannually
  //current interest = 10%; compounding period semiannually
  cashFlow = b.CashFlow(1000.00, 10.0, b.GetCompoundingPeriod('s', false), 20.0, b.GetTimePeriod('s', false))
  price = b.CurrentPrice(cashFlow, 10.0, b.GetCompoundingPeriod('s', false))
  fmt.Printf("price($1,000.00 (par)) = $%.2f\n", price)
  //FV = $100.00
  //coupon rate = 10%; compounding period annually
  //t = 3-year
  //current interest = 9%; compounding period annually
  cashFlow = b.CashFlow(100.00, 10.0, b.GetCompoundingPeriod('a', false), 3, b.GetTimePeriod('y', false))
  price = b.CurrentPrice(cashFlow, 9.0, b.GetCompoundingPeriod('a', false))
  fmt.Printf("price($102.531) = $%.2f\n", price)
  //FV = $1,000.00
  //coupon rate = 10%; annually
  //time to call = 9 years
  //bond price = $1,494.93
  //call price = $1,100.00
  var ytc float64 = b.YieldToCall(1000.00, 10.0, b.GetCompoundingPeriod('a', false), 9, b.GetTimePeriod('y', false), 1494.93, 1100.00)
  fmt.Printf("ytc(4.21%%) = %.2f%%\n", ytc)
  //FV = $1,000.00
  //coupon rate = 10%; semiannually
  //time to call = 5 years
  //bond price = $1,175.00
  //call price = $1,100.00
  ytc = b.YieldToCall(1000.00, 10.0, b.GetCompoundingPeriod('s', false), 5, b.GetTimePeriod('y', false), 1175.00, 1100.00)
  fmt.Printf("ytc(7.43%%) = %.2f%%\n", ytc)
  //FV = $100.00
  //coupon rate = 10%; compounding period annually
  //t = 3-year
  //current interest = 9%; compounding period annually
  cashFlow = b.CashFlow(100.00, 10.0, b.GetCompoundingPeriod('a', false), 3.0, b.GetTimePeriod('y', false))
  var ytm float64 = b.YieldToMaturity(cashFlow, b.CurrentPrice(cashFlow, 9.0, b.GetCompoundingPeriod('a', false)), b.GetCompoundingPeriod('a', false))
  fmt.Printf("ytm(9.00%%) = %.2f%%\n", ytm)
  //FV = $1,000.00
  //coupon rate = 10%; compounding period annually
  //t = 10-year
  //price = $920.00
  cashFlow = b.CashFlow(1000.00, 10.0, b.GetCompoundingPeriod('a', false), 10, b.GetTimePeriod('y', false))
  ytm = b.YieldToMaturity(cashFlow, 920.00, b.GetCompoundingPeriod('a', false))
  fmt.Printf("ytm(11.3801%%) = %.2f%%\n", ytm)
  //FV = $100.00
  //coupon rate = 5%; compounding period semiannually
  //t = 30-month
  //price = $95.92
  cashFlow = b.CashFlow(100.00, 5.0, b.GetCompoundingPeriod('s', false), 30, b.GetTimePeriod('m', false))
  ytm = b.YieldToMaturity(cashFlow, 95.92, b.GetCompoundingPeriod('s', false))
  fmt.Printf("ytm(6.80223%%) = %.2f%%\n", ytm)
  //FV = $100.00
  //coupon rate = 10%; compounding period annually
  //t = 3-year
  //current interest = 9%; compounding period annually
  cashFlow = b.CashFlow(100.00, 10.0, b.GetCompoundingPeriod('a', false), 3, b.GetTimePeriod('y', false))
  var years = b.Duration(cashFlow, 9, b.CurrentPrice(cashFlow, 9, b.GetCompoundingPeriod('a', false)))
  fmt.Printf("years(2.738954) = %.2f\n", years)
	


  // fmt.Println("eps = ", math.Nextafter(1.0, 2.0) - 1.0)
  fmt.Println("Hello, world.")
}
