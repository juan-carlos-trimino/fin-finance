//
package mathutil

import (
	"math"
)

type MathUtil struct{}

/***
                             Prologue
Copyright (c) 1992 by Juan Carlos Trimino. All rights reserved.
Purpose:
  Computes (x mod y).
Algorithm:
  If x and y are any real numbers, we define the following binary operation:
    x mod y = x - y * floor(x / y), if y != 0;  x mod 0 = x.                       (1)

  From this definition we can see that when y != 0,

          x           x      x mod y
    0 <= --- - floor(---) = --------- < 1;                                         (2)
          y           y         y

  therefore

  (a) if y > 0, then 0 <= x mod y < y;
  (b) if y < 0, then 0 >= x mod y > y;
  (c) the quantity x - (x mod y) is an integral multiple of y; and so we may think of x mod y as
      the remainder when x is divided by y.

  Thus, "mod" is a familiar operation when x and y are integers:

  5 mod 3 = 2,
  18 mod 3 = 0,
  -2 mod 3 = 1.

  We have x mod y = 0 iff x is a multiple of y; i.e., iff x is divisible by y.

  The "mod" operation is also useful when x and y take arbitrary real values; e.g., with
  trigonometric functions we can write

    tan x = tan (x mod pi).                                                        (3)

  The quantity x mod 1 is the "fractional part" of x; we have, by Eq. (1),

    x = floor(x) + (x mod 1).                                                      (4)
Return Value:
  x mod y.
  Note: The mod (modulus) and rem (remainder) functions return identical results for positive
        quantities, but different results for negative quantities.
Usage
  ...
  var um umath.Umath
  var x float64 = 45256.25
  //Reduce the angle x (in radians) to the range [0, 2 * math.Pi).
  fmt.Printf("mod = %f\n", um.Mod64(x, 2 * math.Pi))
  //Reduce the angle x (in radians) to the range [-math.Pi, math.Pi).
  fmt.Printf("mod = %f\n", um.Mod64(x + math.Pi, 2 * math.Pi) - math.Pi)
  ...
History of Changes:
  Rel   Programmer Date   Description
  ----- ---------- ------ -------------------------------------------------------------------------
  1.00  JC Trimino 060392 Original delivery.
  1.01  JC Trimino 123108 Changed to use template.
  1.02  JC Trimino xxxxxx Translated to Go.
***/
func (m *MathUtil) Mod64(x, y float64) float64 {
  //Epsilon is the smallest value that, when added to one, yields a result different from one.
  var epsilon float64 = math.Nextafter(float64(1), float64(2)) - float64(1)
  //If |y| < EPSILON, then y = 0.
  if math.Abs(y) < epsilon {  // x mod 0 = x
    return x
  } else if math.Abs(x) < epsilon {  // 0 mod y = 0
    return 0.0
  } else {
    return (x - y * math.Floor(x / y))
  }
}
