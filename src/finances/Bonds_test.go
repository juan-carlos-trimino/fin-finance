//Testing the functions in Bonds.go.
package finances

/***
To build and run the tests:
$ go test

The -v flag prints the name and execution time of each test in the package:
$ go test -v

The -run flag, whose argument is a regular expression, causes 'go test' to run only those tests
whose function name matches the pattern:
$ go test -v -run="Bonds"
***/

import (
  "fmt"
  "math"
  "testing"
)

func TestBonds_CurrentPrice(t *testing.T) {
  type test struct {
    FV float64 //Face value
    couponRate float64
    cp byte //Compounding period
    n float64
    tp byte //Time period
    currentRate float64
    currentRate_cp byte //Compounding period
    want float64
  }
  var tests = []test {
    //Price = $940.25 (discount)
    { FV: 1000.00, couponRate: 10.0, cp: 's', n: 10.0, tp: 'y', currentRate: 11.0,
      currentRate_cp: 'S', want: 940.2480875 },
    //Price = $1,065.04 (premium)
    { FV: 1000.00, couponRate: 10.0, cp: 's', n: 10.0, tp: 'y', currentRate: 9.0,
      currentRate_cp: 'S', want: 1065.03968225 },
    //Price = $1,000.00 (par))
    { FV: 1000.00, couponRate: 10.0, cp: 's', n: 10.0, tp: 'y', currentRate: 10.0,
      currentRate_cp: 'S', want: 1000.00 },
    //Price = $1,854.43 (premium)
    { FV: 1000.00, couponRate: 10.0, cp: 's', n: 10.0, tp: 'y', currentRate: 1.0,
      currentRate_cp: 'S', want: 1854.43386160 },
    //Price = $603.44 (discount)
    { FV: 1000.00, couponRate: 10.0, cp: 's', n: 10.0, tp: 'y', currentRate: 19.0,
      currentRate_cp: 'S', want: 603.44280474 },
    //Price = $1,000.00 (par)
    { FV: 1000.00, couponRate: 10.0, cp: 's', n: 20.0, tp: 's', currentRate: 10.0,
      currentRate_cp: 'S', want: 1000.00 },
    //Price = $102.531
    { FV: 100.00, couponRate: 10.0, cp: 'a', n: 3.0, tp: 'y', currentRate: 9.0,
      currentRate_cp: 'a', want: 102.53129466 },
    //Purchased a 5-year $10,000.00, 4.5% interest bearing bond at the end of the third interest
    //date at a price that would yield 6% per annum, compounded semiannually.
    { FV: 10_000.00, couponRate: 4.5, cp: 's', n: 3.5, tp: 'y', currentRate: 6.0,
      currentRate_cp: 's', want: 9_532.7287783 },
    //Purchased a 5-year $1,000.00, 6% interest bearing bond at a price that would yield 5% per
    //annum, compounded semiannually.
    { FV: 1_000.00, couponRate: 6, cp: 's', n: 5, tp: 'y', currentRate: 5.0,
      currentRate_cp: 's', want: 1_043.7603196 },
    //On Jan 1, 19x0, purchased a 4%, 10-year, $1,000.00 bond maturing Jan 1, 19x8, with interest
    //at 6% per annum, compounded semiannually.
    { FV: 1_000.00, couponRate: 4, cp: 's', n: 8, tp: 'y', currentRate: 6.0,
      currentRate_cp: 's', want: 874.3889797 },
  }
  var b Bonds
  for _, tc := range tests {
    var cashFlow = b.CashFlow(tc.FV, tc.couponRate, b.GetCompoundingPeriod(tc.cp, false), tc.n,
                              b.GetTimePeriod(tc.tp, false))
    var price = b.CurrentPrice(cashFlow, tc.currentRate,
                               b.GetCompoundingPeriod(tc.currentRate_cp, false))
    if math.Abs(price - tc.want) < 1e-5 {
      fmt.Printf("Price = $%.2f\n", price)
    } else {
      t.Errorf("Price = $%.10f, Want = $%.10f", price, tc.want)
    }
  }
}

func TestBonds_CurrentPriceContinuous(t *testing.T) {
  type test struct {
    FV float64 //Face value
    couponRate float64
    cp byte //Compounding period
    n float64
    tp byte //Time period
    currentRate float64
    want float64
  }
  var tests = []test {
    //Price = $101.463758
    { FV: 100.00, couponRate: 10.0, cp: 'a', n: 3.0, tp: 'y', currentRate: 9.0, want: 101.463758 },
    //Price = $104.281666
    { FV: 100.00, couponRate: 10.0, cp: 'a', n: 3.0, tp: 'y', currentRate: 8.0, want: 104.281666 },
  }
  var b Bonds
  for _, tc := range tests {
    var cashFlow = b.CashFlow(tc.FV, tc.couponRate, b.GetCompoundingPeriod(tc.cp, false), tc.n,
                              b.GetTimePeriod(tc.tp, false))
    var price = b.CurrentPriceContinuous(cashFlow, tc.currentRate)
    if math.Abs(price - tc.want) < 1e-5 {
      fmt.Printf("Price = $%.2f\n", price)
    } else {
      t.Errorf("Price = $%.10f, Want = $%.10f", price, tc.want)
    }
  }
}

func TestBonds_DurationContinuous(t *testing.T) {
  type test struct {
    FV float64 //Face value
    couponRate float64
    cp byte //Compounding period
    n float64
    tp byte //Time period
    currentRate float64
    want float64
  }
  var tests = []test {
    //Duration = 2.737529
    { FV: 100.00, couponRate: 10.0, cp: 'a', n: 3.0, tp: 'y', currentRate: 9.0, want: 2.737529 },
  }
  var b Bonds
  for _, tc := range tests {
    var cashFlow = b.CashFlow(tc.FV, tc.couponRate, b.GetCompoundingPeriod(tc.cp, false), tc.n,
                              b.GetTimePeriod(tc.tp, false))
    var duration = b.DurationContinuous(cashFlow, tc.currentRate,
                                        b.CurrentPriceContinuous(cashFlow, tc.currentRate))
    if math.Abs(duration - tc.want) < 1e-5 {
      fmt.Printf("Duration = %.2f\n", duration)
    } else {
      t.Errorf("Duration = %.10f, Want = %.10f", duration, tc.want)
    }
  }
}

func TestBonds_YieldToMaturityContinuous(t *testing.T) {
  type test struct {
    FV float64 //Face value
    couponRate float64
    cp byte //Compounding period
    n float64
    tp byte //Time period
    currentRate float64
    want float64
  }
  var tests = []test {
    //Duration = 9.0
    { FV: 100.00, couponRate: 10.0, cp: 'a', n: 3.0, tp: 'y', currentRate: 9.0, want: 9.0 },
  }
  var b Bonds
  for _, tc := range tests {
    var cashFlow = b.CashFlow(tc.FV, tc.couponRate, b.GetCompoundingPeriod(tc.cp, false), tc.n,
                              b.GetTimePeriod(tc.tp, false))
    var ytm = b.YieldToMaturityContinuous(cashFlow,
                                          b.CurrentPriceContinuous(cashFlow, tc.currentRate))
    if math.Abs(ytm - tc.want) < 1e-5 {
      fmt.Printf("YTM = %.2f\n", ytm)
    } else {
      t.Errorf("YTM = %.10f, Want = %.10f", ytm, tc.want)
    }
  }
}

func TestBonds_MacaulayDurationContinuous(t *testing.T) {
  type test struct {
    FV float64 //Face value
    couponRate float64
    cp byte //Compounding period
    n float64
    tp byte //Time period
    currentRate float64
    want float64
  }
  var tests = []test {
    //Duration = 2.737529
    { FV: 100.00, couponRate: 10.0, cp: 'a', n: 3.0, tp: 'y', currentRate: 9.0, want: 2.737529 },
  }
  var b Bonds
  for _, tc := range tests {
    var cashFlow = b.CashFlow(tc.FV, tc.couponRate, b.GetCompoundingPeriod(tc.cp, false), tc.n,
                              b.GetTimePeriod(tc.tp, false))
    var duration = b.MacaulayDurationContinuous(cashFlow,
                                                b.CurrentPriceContinuous(cashFlow, tc.currentRate))
    if math.Abs(duration - tc.want) < 1e-5 {
      fmt.Printf("Duration = %.2f\n", duration)
    } else {
      t.Errorf("Duration = %.10f, Want = %.10f", duration, tc.want)
    }
  }
}

func TestBonds_ConvexityContinuous(t *testing.T) {
  type test struct {
    FV float64 //Face value
    couponRate float64
    cp byte //Compounding period
    n float64
    tp byte //Time period
    currentRate float64
    want float64
  }
  var tests = []test {
    //Convexity = 7.867793
    { FV: 100.00, couponRate: 10.0, cp: 'a', n: 3.0, tp: 'y', currentRate: 9.0, want: 7.867793 },
  }
  var b Bonds
  for _, tc := range tests {
    var cashFlow = b.CashFlow(tc.FV, tc.couponRate, b.GetCompoundingPeriod(tc.cp, false), tc.n,
                              b.GetTimePeriod(tc.tp, false))
    var convexity = b.ConvexityContinuous(cashFlow, tc.currentRate,
                                          b.CurrentPriceContinuous(cashFlow, tc.currentRate))
    if math.Abs(convexity - tc.want) < 1e-5 {
      fmt.Printf("Duration = %.2f\n", convexity)
    } else {
      t.Errorf("Duration = %.10f, Want = %.10f", convexity, tc.want)
    }
  }
}

func TestBonds_YieldToMaturity(t *testing.T) {
  type test struct {
    withPrice bool
    FV float64 //Face value
    couponRate float64
    cp byte //Compounding period
    n float64
    tp byte //Time period
    currentRate float64
    currentRate_cp byte //Compounding period
    price float64
    want float64
  }
  var tests = []test {
    //ytm = 9.00%
    { withPrice: false, FV: 100.00, couponRate: 10.0, cp: 'a', n: 3.0, tp: 'y', currentRate: 9.0,
      currentRate_cp: 'a', want: 9.0 },
    //ytm = 11.380136966%
    { withPrice: true, FV: 1000.00, couponRate: 10.0, cp: 'a', n: 10.0, tp: 'y', price: 920.00,
      want: 11.380136966 },
    //ytm = 6.80223%
    { withPrice: true, FV: 100.00, couponRate: 5.0, cp: 's', n: 30.0, tp: 'm', price: 95.92,
      want: 6.80223 },
    //ytm = 11.341%
    { withPrice: true, FV: 1000.00, couponRate: 10.0, cp: 'm', n: 10.0, tp: 'y', price: 920,
      want: 11.341 },
    //ytm = 11.358871%
    { withPrice: true, FV: 1000.00, couponRate: 10.0, cp: 's', n: 10.0, tp: 'y', price: 920,
      want: 11.358871 },
    }
  var b Bonds
  for _, tc := range tests {
    var cashFlow = b.CashFlow(tc.FV, tc.couponRate, b.GetCompoundingPeriod(tc.cp, false), tc.n,
		                          b.GetTimePeriod(tc.tp, false))
    var ytm float64
    if tc.withPrice {
      ytm = b.YieldToMaturity(cashFlow, tc.price, b.GetCompoundingPeriod(tc.cp, false))
    } else {
      ytm = b.YieldToMaturity(cashFlow, b.CurrentPrice(cashFlow, tc.currentRate,
                              b.GetCompoundingPeriod(tc.currentRate_cp, false)),
                              b.GetCompoundingPeriod(tc.currentRate_cp, false))
    }
    //
    if math.Abs(ytm - tc.want) < 1e-5 {
      fmt.Printf("ytm = %.2f%%\n", ytm)
    } else {
      t.Errorf("ytm = %.10f%%, Want = %.10f%%", ytm, tc.want)
    }
  }
}

func TestBonds_YieldToCall(t *testing.T) {
  type test struct {
    FV float64 //Face value
    couponRate float64
    cp byte //Compounding period
    n float64 //Time to call
    tp byte //Time period
    bondPrice float64
    callPrice float64
    want float64
  }
  var tests = []test {
    //ytc = 4.21485471%
    { FV: 1000.00, couponRate: 10.0, cp: 'a', n: 9.0, tp: 'y', bondPrice: 1494.93,
      callPrice: 1100.00, want: 4.21485471 },
    //ytm = 7.43329954%
    { FV: 1000.00, couponRate: 10.0, cp: 's', n: 5.0, tp: 'y', bondPrice: 1175.00,
      callPrice: 1100.00, want: 7.43329954 },
  }
  var b Bonds
  for _, tc := range tests {
    var ytc = b.YieldToCall(tc.FV, tc.couponRate, b.GetCompoundingPeriod(tc.cp, false), tc.n,
                            b.GetTimePeriod(tc.tp, false), tc.bondPrice, tc.callPrice)
    if math.Abs(ytc - tc.want) < 1e-5 {
      fmt.Printf("ytc = %.2f%%\n", ytc)
    } else {
      t.Errorf("ytc = %.10f%%, Want = %.10f%%", ytc, tc.want)
    }
  }
}

func TestBonds_Convexity(t *testing.T) {
  type test struct {
    FV float64 //Face value
    couponRate float64
    cp byte //Compounding period
    n float64 //Time to call
    tp byte //Time period
    currentInterest float64
    currentInterest_cp byte
    want float64
  }
  var tests = []test {
    //Cx = 8.932479
    { FV: 100.00, couponRate: 10.0, cp: 'a', n: 3.0, tp: 'y', currentInterest: 9.0,
      currentInterest_cp: 'a', want: 8.932479 },
  }
  var b Bonds
  for _, tc := range tests {
    var cashFlow = b.CashFlow(tc.FV, tc.couponRate, b.GetCompoundingPeriod(tc.cp, false),
                              tc.n, b.GetTimePeriod(tc.tp, false))
    var Cx = b.Convexity(cashFlow, tc.currentInterest,
                         b.GetCompoundingPeriod(tc.currentInterest_cp, false))
    if math.Abs(Cx - tc.want) < 1e-5 {
      fmt.Printf("Cx = %.2f\n", Cx)
    } else {
      t.Errorf("Cx = %.10f, Want = %.10f", Cx, tc.want)
    }
  }
}

func TestBonds_ModifiedDuration(t *testing.T) {
  type test struct {
    FV float64 //Face value
    couponRate float64
    cp byte //Compounding period
    n float64 //Time to call
    tp byte //Time period
    currentInterest float64
    currentInterest_cp byte
    want float64
  }
  var tests = []test {
    //MDuration = 2.512801%
    { FV: 100.00, couponRate: 10.0, cp: 'a', n: 3.0, tp: 'y', currentInterest: 9.0,
      want: 2.512801 },
    //MDuration = 2.62144622%
    { FV: 1000.00, couponRate: 10.0, cp: 'a', n: 3.0, tp: 'y', currentInterest: 5.0,
      want: 2.62144622 },
  }
  var b Bonds
  for _, tc := range tests {
    var cashFlow = b.CashFlow(tc.FV, tc.couponRate, b.GetCompoundingPeriod(tc.cp, false), tc.n,
                              b.GetTimePeriod(tc.tp, false))
    var MDuration = b.ModifiedDuration(cashFlow, b.GetCompoundingPeriod(tc.cp, false),
                                       b.CurrentPrice(cashFlow, tc.currentInterest,
                                       b.GetCompoundingPeriod(tc.cp, false)))
    if math.Abs(MDuration - tc.want) < 1e-5 {
      fmt.Printf("MDuration = %.2f%%\n", MDuration)
    } else {
      t.Errorf("MDuration = %.10f%%, Want = %.10f%%", MDuration, tc.want)
    }
  }
}

func TestBonds_MacaulayDuration(t *testing.T) {
  type test struct {
    FV float64 //Face value
    couponRate float64
    cp byte //Compounding period
    n float64 //Time to call
    tp byte //Time period
    currentInterest float64
    currentInterest_cp byte
    want float64
  }
  var tests = []test {
    //macyears = 2.738954
    { FV: 100.00, couponRate: 10.0, cp: 'a', n: 3.0, tp: 'y', currentInterest: 9.0,
      want: 2.738954 },
    //macyears = 2.7525185
    { FV: 1000.00, couponRate: 10.0, cp: 'a', n: 3.0, tp: 'y', currentInterest: 5.0,
      want: 2.7525185 },
  }
  var b Bonds
  for _, tc := range tests {
    var cashFlow = b.CashFlow(tc.FV, tc.couponRate, b.GetCompoundingPeriod(tc.cp, false), tc.n,
                              b.GetTimePeriod(tc.tp, false))
    var macyears = b.MacaulayDuration(cashFlow, b.GetCompoundingPeriod(tc.cp, false),
                              b.CurrentPrice(cashFlow, tc.currentInterest,
                                             b.GetCompoundingPeriod(tc.cp, false)))
    if math.Abs(macyears - tc.want) < 1e-5 {
      fmt.Printf("macyears = %.2f\n", macyears)
    } else {
      t.Errorf("macyears = %.10f, Want = %.10f", macyears, tc.want)
    }
  }
}

func TestBonds_Duration(t *testing.T) {
  type test struct {
    FV float64 //Face value
    couponRate float64
    cp byte //Compounding period
    n float64 //Time to call
    tp byte //Time period
    currentInterest float64
    currentInterest_cp byte
    want float64
  }
  var tests = []test {
    //years = 2.738954
    { FV: 100.00, couponRate: 10.0, cp: 'a', n: 3.0, tp: 'y', currentInterest: 9.0,
      currentInterest_cp: 'a', want: 2.738954 },
  }
  var b Bonds
  for _, tc := range tests {
    var cashFlow = b.CashFlow(tc.FV, tc.couponRate, b.GetCompoundingPeriod(tc.cp, false), tc.n,
                              b.GetTimePeriod(tc.tp, false))
    var years = b.Duration(cashFlow, tc.currentInterest,
                           b.CurrentPrice(cashFlow, tc.currentInterest,
                                          b.GetCompoundingPeriod(tc.currentInterest_cp, false)))
    if math.Abs(years - tc.want) < 1e-5 {
      fmt.Printf("years = %.2f\n", years)
    } else {
      t.Errorf("years = %.10f, Want = %.10f", years, tc.want)
    }
  }
}

func TestBonds_CashFlow(t *testing.T) {
  type test struct {
    FV float64 //Face value
    couponRate float64
    cp byte //Compounding period
    n float64 //Time to call
    tp byte //Time period
    want []float64
  }
  var tests = []test {
    { FV: 1000.00, couponRate: 3.0, cp: 's', n: 5.0, tp: 'y',
      want: []float64{15.0, 15.0, 15.0, 15.0, 15.0, 15.0, 15.0, 15.0, 15.0, 1015.00} },
  }
  var b Bonds
  for _, tc := range tests {
    var cashFlow = b.CashFlow(tc.FV, tc.couponRate, b.GetCompoundingPeriod(tc.cp, false), tc.n,
                              b.GetTimePeriod(tc.tp, false))
    for idx := range cashFlow {
      if math.Abs(cashFlow[idx] - tc.want[idx]) < 1e-5 {
        fmt.Printf("Payment[%d] = $%.2f\n", idx + 1, cashFlow[idx])
      } else {
        t.Errorf("Payment[%d] = $%.10f, Want = $%.10f", idx + 1, cashFlow[idx], tc.want[idx])
      }
    }
  }
}
