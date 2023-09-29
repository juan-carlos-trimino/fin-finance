//Testing the functions in SimpleInterest.go.
package finances

/***
To build and run the tests:
$ go test

The -v flag prints the name and execution time of each test in the package:
$ go test -v

The -run flag, whose argument is a regular expression, causes 'go test' to run only those tests
whose function name matches the pattern:
$ go test -v -run="Simple"
***/

import (
  "fmt"
  "math"
  "testing"
)

func TestSimpleInterest_Interest(t *testing.T) {
  type test struct {
    typeOf string
    p float64 //Principal
    rate float64 //Interest rate
    cp int //Compounding period
    n float64 //Time
    tp int //Time period
    want float64
  }
  var tests = []test {
    { typeOf: "accurate", p: 156.00, rate: 153.846, cp: Monthly, n: 1.0, tp: Months, want: 19.999980 },
    { typeOf: "accurate", p: 100.00, rate: 4.00, cp: Annually, n: 1.0, tp: Months, want: 0.3333333333 },
    { typeOf: "accurate", p: 10000.00, rate: 9.00, cp: Annually, n: 153.0, tp: Daily365, want: 377.26027397 },
    { typeOf: "banker's", p: 100.00, rate: 4.00, cp: Monthly, n: 1, tp: Months, want: 0.3333333333 },
    { typeOf: "banker's", p: 100.00, rate: 4.00, cp: Annually, n: 1, tp: Months, want: 0.3333333333 },
    { typeOf: "banker's", p: 10000.00, rate: 9.00, cp: Annually, n: 153, tp: Daily360, want: 382.50 },
    { typeOf: "ordinary", p: 100.00, rate: 4.00, cp: Monthly, n: 1, tp: Months, want: 0.3333333333 },
    { typeOf: "ordinary", p: 100.00, rate: 4.00, cp: Annually, n: 1, tp: Months, want: 0.3333333333 },
    { typeOf: "ordinary", p: 10000.00, rate: 9.00, cp: Annually, n: 153, tp: Daily360, want: 375.00 },
  }
  var si SimpleInterest
  for _, tc := range tests {
    switch tc.typeOf {
    case "accurate":
      i := si.AccurateInterest(tc.p, tc.rate / hundred, tc.cp, tc.n, tc.tp)
      if math.Abs(i - tc.want) < 1e-5 {
        fmt.Printf("Accurate Interest = $%.2f\n", i)
      } else {
        t.Errorf("Accurate Interest = $%.10f, Want = $%.10f", i, tc.want)
      }
    case "banker's":
      i := si.BankersInterest(tc.p, tc.rate / hundred, tc.cp, tc.n, tc.tp)
      if math.Abs(i - tc.want) < 1e-5 {
        fmt.Printf("Banker's Interest = $%.2f\n", i)
      } else {
        t.Errorf("Banker's Interest = $%.10f, Want = $%.10f", i, tc.want)
      }
    case "ordinary":
      i := si.OrdinaryInterest(tc.p, tc.rate / hundred, tc.cp, tc.n, tc.tp)
      if math.Abs(i - tc.want) < 1e-5 {
        fmt.Printf("Ordinary Interest = $%.2f\n", i)
      } else {
        t.Errorf("Ordinary Interest = $%.10f, Want = $%.10f", i, tc.want)
      }
    default:
      fmt.Printf("Unknown Interest = %s\n", tc.typeOf)
    }
  }
}

func TestSimpleInterest_Rate(t *testing.T) {
  type test struct {
    typeOf string
    p float64 //Principal
    a float64 //Amount of interest
    n float64 //Time
    tp int //Time period
    want float64
  }
  var tests = []test {
    { typeOf: "accurate", p: 156.00, a: 20.00, n: 1.0, tp: Months, want: 153.84615384 },
    { typeOf: "accurate", p: 156.00, a: 20.00, n: 1.0, tp: Years, want: 12.8205128 },
    { typeOf: "accurate", p: 10000.00, a: 377.26, n: 153.0, tp: Daily365, want: 9.0 },
    { typeOf: "banker's", p: 100.00, a: 4.00, n: 1.0, tp: Months, want: 48.00 },
    { typeOf: "banker's", p: 100.00, a: 10.33, n: 1.0, tp: Months, want: 123.960 },
    { typeOf: "banker's", p: 10000.00, a: 382.50, n: 153.0, tp: Daily360, want: 9.0 },
    { typeOf: "ordinary", p: 100.00, a: 4.00, n: 1.0, tp: Months, want: 48.00 },
    { typeOf: "ordinary", p: 100.00, a: 23.33, n: 1.0, tp: Months, want: 279.960 },
    { typeOf: "ordinary", p: 10000.00, a: 375.00, n: 153.0, tp: Daily360, want: 9.00 },
  }
  var si SimpleInterest
  for _, tc := range tests {
    switch tc.typeOf {
    case "accurate":
      i := si.AccurateRate(tc.p, tc.a, tc.n, tc.tp) * 100.0
      if math.Abs(i - tc.want) < 1e-5 {
        fmt.Printf("Accurate Rate = %.2f%%\n", i)
      } else {
        t.Errorf("Accurate Rate = %.10f%%, Want = %.10f%%", i, tc.want)
      }
    case "banker's":
      i := si.BankersRate(tc.p, tc.a, tc.n, tc.tp) * 100.0
      if math.Abs(i - tc.want) < 1e-5 {
        fmt.Printf("Banker's Rate = %.2f%%\n", i)
      } else {
        t.Errorf("Banker's Rate = %.10f%%, Want = %.10f%%", i, tc.want)
      }
    case "ordinary":
      i := si.OrdinaryRate(tc.p, tc.a, tc.n, tc.tp) * 100.0
      if math.Abs(i - tc.want) < 1e-5 {
        fmt.Printf("Ordinary Rate = %.2f%%\n", i)
      } else {
        t.Errorf("Ordinary Rate = %.10f%%, Want = %.10f%%", i, tc.want)
      }
    default:
      fmt.Printf("Unknown Rate = %s\n", tc.typeOf)
    }
  }
}

func TestSimpleInterest_Principal(t *testing.T) {
  type test struct {
    typeOf string
    a float64 //Interest of interest
    rate float64 //Interest rate
    cp int //Compounding period
    n float64 //Time
    tp int //Time period
    want float64
  }
  var tests = []test {
    { typeOf: "accurate", a: 4.00, rate: 4.00, cp: Monthly, n: 1.0, tp: Months, want: 1200.00 },
    { typeOf: "accurate", a: 0.33, rate: 4.00, cp: Annually, n: 1.0, tp: Months, want: 99.00 },
    { typeOf: "accurate", a: 377.26, rate: 9.0, cp: Annually, n: 153.0, tp: Daily365, want: 9999.9927378 },
    { typeOf: "banker's", a: 4.00, rate: 4.00, cp: Monthly, n: 1.0, tp: Months, want: 1200.00 },
    { typeOf: "banker's", a: 0.33, rate: 4.00, cp: Annually, n: 1.0, tp: Months, want: 99.00 },
    { typeOf: "banker's", a: 382.50, rate: 9.0, cp: Annually, n: 153.0, tp: Daily360, want: 10000.00 },
    { typeOf: "ordinary", a: 4.00, rate: 4.00, cp: Monthly, n: 1.0, tp: Months, want: 1200.00 },
    { typeOf: "ordinary", a: 0.33, rate: 4.00, cp: Annually, n: 1.0, tp: Months, want: 99.00 },
    { typeOf: "ordinary", a: 375.00, rate: 9.0, cp: Annually, n: 153.0, tp: Daily360, want: 10000.00 },
  }
  var si SimpleInterest
  for _, tc := range tests {
    switch tc.typeOf {
    case "accurate":
      p := si.AccuratePrincipal(tc.a, tc.rate, tc.cp, tc.n, tc.tp) * 100.0
      if math.Abs(p - tc.want) < 1e-5 {
        fmt.Printf("Accurate Principal = $%.2f\n", p)
      } else {
        t.Errorf("Accurate Principal = $%.10f, Want = $%.10f", p, tc.want)
      }
    case "banker's":
      p := si.BankersPrincipal(tc.a, tc.rate, tc.cp, tc.n, tc.tp) * 100.0
      if math.Abs(p - tc.want) < 1e-5 {
        fmt.Printf("Banker's Principal = $%.2f\n", p)
      } else {
        t.Errorf("Banker's Principal = $%.10f, Want = $%.10f", p, tc.want)
      }
    case "ordinary":
      p := si.OrdinaryPrincipal(tc.a, tc.rate, tc.cp, tc.n, tc.tp) * 100.0
      if math.Abs(p - tc.want) < 1e-5 {
        fmt.Printf("Ordinary Principal = $%.2f\n", p)
      } else {
        t.Errorf("Ordinary Principal = $%.10f, Want = $%.10f", p, tc.want)
      }
    default:
      fmt.Printf("Unknown Rate = %s\n", tc.typeOf)
    }
  }
}

func TestSimpleInterest_Time(t *testing.T) {
  type test struct {
    typeOf string
    p float64 //Principal
    a float64 //Interest of interest
    rate float64 //Interest rate
		cp int
    want float64
  }
  var tests = []test {
    { typeOf: "accurate", p:100, a: 4.00, rate: 4.00, cp: Monthly, want: 12.00 },
    { typeOf: "accurate", p:100, a: 1.33, rate: 4.00, cp: Annually, want: 0.3325 },
    { typeOf: "accurate", p:10000.00, a: 377.26, rate: 9.0, cp: Daily365, want: 152.99988888 },
    { typeOf: "banker's", p:100, a: 4.00, rate: 4.00, cp: Monthly, want: 12.00 },
    { typeOf: "banker's", p:100, a: 0.33, rate: 4.00, cp: Annually, want: 0.08250 },
    { typeOf: "banker's", p:10000.00, a: 382.50, rate: 9.0, cp: Daily360, want: 153.00 },
    { typeOf: "ordinary", p:100, a: 4.00, rate: 4.00, cp: Monthly, want: 12.00 },
    { typeOf: "ordinary", p:100, a: 0.33, rate: 4.00, cp: Annually, want: 0.08250 },
    { typeOf: "ordinary", p:10000.00, a: 375.00, rate: 9.0, cp: Daily360, want: 150.00 },
  }
  var si SimpleInterest
  for _, tc := range tests {
    switch tc.typeOf {
    case "accurate":
      time := si.AccurateTime(tc.p, tc.a, tc.rate / hundred, tc.cp)
      if math.Abs(time - tc.want) < 1e-5 {
        fmt.Printf("Accurate Time = %.2f\n", time)
      } else {
        t.Errorf("Accurate Time = %.10f, Want = %.10f", time, tc.want)
      }
    case "banker's":
      time := si.BankersTime(tc.p, tc.a, tc.rate / hundred, tc.cp)
      if math.Abs(time - tc.want) < 1e-5 {
        fmt.Printf("Banker's Time = %.2f\n", time)
      } else {
        t.Errorf("Banker's Time = %.10f, Want = %.10f", time, tc.want)
      }
    case "ordinary":
      time := si.BankersTime(tc.p, tc.a, tc.rate / hundred, tc.cp)
      if math.Abs(time - tc.want) < 1e-5 {
        fmt.Printf("Ordinary Time = %.2f\n", time)
      } else {
        t.Errorf("Ordinary Time = %.10f, Want = %.10f", time, tc.want)
      }
    default:
      fmt.Printf("Unknown Rate = %s\n", tc.typeOf)
    }
  }
}
