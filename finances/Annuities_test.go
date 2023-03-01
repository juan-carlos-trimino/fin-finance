// Testing the functions in Annuities.go.
package finances

/***
To build and run the tests:
$ go test

The -v flag prints the name and execution time of each test in the package:
$ go test -v

The -run flag, whose argument is a regular expression, causes 'go test' to run only those tests
whose function name matches the pattern:
$ go test -v -run="Annuities"
***/

import (
	"fmt"
	"math"
	"testing"
)

func TestAnnuities_O_Interest_PV_PMT(t *testing.T) {
  type test struct {
    PV float64 //Present value
    PMT float64 //Payment
    n float64
    i1 float64
    i2 float64
    cp int
    accurancy float64
    want float64
  }
  var tests = []test {
    //i = 0.7628634% per month
    //i = 9.154323% per year = 0.7628634 * cp(Monthly)
    { PV: 24000.00, PMT: 500.00, n: 60.0, i1: 1.0, i2: 31.0, cp: Monthly, accurancy: 1e-6,
      want: 0.7628634 },
    //i = 0.94007411% per month
    //i = 11.2808893% per year = 0.94007411 * cp(Monthly)
    { PV: 11200.00, PMT: 291.00, n: 48.0, i1: 4.0, i2: 12.0, cp: Monthly, accurancy: 1e-6,
      want: 0.94007411 },
    //i = 10.91616% per year
    { PV: 50000.00, PMT: 13500.00, n: 5.0, i1: 10.0, i2: 15.0, cp: Annually, accurancy: 1e-6,
      want: 10.916174523 },
  }
  var a Annuities
  for _, tc := range tests {
    var i = a.O_Interest_PV_PMT(tc.PV, tc.PMT, tc.n, tc.i1, tc.i2, tc.cp, tc.accurancy) * 100.0
    if math.Abs(i - tc.want) < 1e-5 {
      fmt.Printf("Interest = %.2f%%\n", i)
    } else {
      t.Errorf("Interest = %.10f%%, Want = %.10f%%", i, tc.want)
    }
  }
}

func TestAnnuities_NominalAndEAR(t *testing.T) {
  type test struct {
    nominal float64
    cp int
    want float64
  }
  var tests = []test {
    { nominal: 12.0, cp: Monthly, want: 12.6825030131 /*ear*/ },
  }
  var a Annuities
  for _, tc := range tests {
    var ear = a.NominalRateToEAR(tc.nominal / hundred, tc.cp)
    var nominal = a.EARToNominalRate(ear, tc.cp)
    ear *= hundred; nominal *= hundred
    if (math.Abs(ear - tc.want) < 1e-5) && (math.Abs(tc.nominal - nominal) < 1e5) {
      fmt.Printf("ear = %.2f%%\n", ear)
      fmt.Printf("nominal = %.2f%%\n", nominal)
    } else {
      t.Errorf("ear = %.10f%%, Want = %.10f%%", ear, tc.want)
      t.Errorf("nominal = %.10f%%, Want = %.10f%%", nominal, tc.nominal)
    }
  }
}

func TestAnnuities_AverageRateOfReturn(t *testing.T) {
  type test struct {
    returns []float64
    want float64
  }
  var tests = []test {
    { returns: []float64{5.0, -3.0, 12.0, 10.0}, want: 5.83831944 },
    { returns: []float64{2.0, 8.0, -1.0, 10.0}, want: 4.655715635 },
  }
  var a Annuities
  for _, tc := range tests {
    var gmr = a.AverageRateOfReturn(tc.returns) * 100.0
    if math.Abs(gmr - tc.want) < 1e-5 {
      fmt.Printf("gmr = %.2f%%\n", gmr)
    } else {
      t.Errorf("gmr = %.10f%%, Want = %.10f%%", gmr, tc.want)
    }
  }
}

func TestAnnuities_RealInterestRate(t *testing.T) {
  type test struct {
    nominalRate float64
    inflationRate float64
    want float64
  }
  var tests = []test {
    { nominalRate: 4.50, inflationRate: 6.50, want: -1.87793427 },
  }
  var a Annuities
  for _, tc := range tests {
    var real = a.RealInterestRate(tc.nominalRate / hundred, tc.inflationRate / hundred) * 100.0
    if math.Abs(real - tc.want) < 1e-5 {
      fmt.Printf("real = %.2f%%\n", real)
    } else {
      t.Errorf("real = %.10f%%, Want = %.10f%%", real, tc.want)
    }
  }
}

func TestAnnuities_GrowthDecayOfFunds(t *testing.T) {
  type test struct {
    rate float64
    cp int
    factor float64
    want float64
  }
  var tests = []test {
    { rate: 15.0, cp: Annually, factor: 2.0, want: 4.959484454 },
  }
  var a Annuities
  for _, tc := range tests {
    var f = a.GrowthDecayOfFunds(tc.factor, tc.rate / hundred, tc.cp)
    if math.Abs(f - tc.want) < 1e-5 {
      fmt.Printf("f = %.2f\n", f)
    } else {
      t.Errorf("f = %.10f, Want = %.10f%%", f, tc.want)
    }
  }
}
