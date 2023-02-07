package umath

/***
A 'go test' (or 'go build') command with no package arguments operates on the package in the
current directory.
$ go test

The -v flag prints the name and execution time of each test in the package.
$ go test -v

The -run flag, whose argument is a regular expression, causes 'go test' to run only those tests
whose function name matches the pattern.
$ go test -v -run="French|Canal"
***/

import (
  "math"
  "testing"
)

func TestMod64(t *testing.T) {
  var tests = []struct {
    input [2] float64  //input[0] = x, input[1] = y
    want float64
  }{
    /***
    The semicolon insertion rules (https://go.dev/ref/spec#Semicolons)
    ***/
    {[2]float64{-2.0, 3.0}, 1.0},
    {[2]float64{math.NaN(), 4.0}, math.NaN()},
    //Reduce the angle x (in radians) to the range [0, 2 * math.Pi).
    {[2]float64{45256.25, 2 * math.Pi}, 4.749418},
  }
  var um Umath
  /***
  'range' on arrays and slices provides both the index and value for each entry. Since we don't
  need the index, we will ignore it with the blank identifier _.
  ***/
  for _, test := range tests {
    got := um.Mod64(test.input[0], test.input[1])
    if math.IsNaN(got) {
      if !math.IsNaN(test.want) {
        t.Errorf("Mod64(%f, %f) = %f; want %f", test.input[0], test.input[1], got, test.want)
      }
    } else if got != test.want {
      t.Errorf("Mod64(%f, %f) = %.10f; want %.10f", test.input[0], test.input[1], got, test.want)
    }
  }
}
