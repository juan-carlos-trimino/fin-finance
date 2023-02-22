//
package mathutil

import (
  "math"
)

const (
  zero float64 = 0.0
  one float64 = 1.0
  two float64 = 2.0
  hundred float64 = 100.0
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
  1.02  JC Trimino 021023 Translated to Go.
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

/***
Root Finding and Nonlinear Sets of Equations
--------------------------------------------
We now consider that most basic of tasks, solving equations numerically. While most equations are
born with both a right-hand side and a left-hand side, one traditionally moves all terms to the
left, leaving

  f(x) = 0

whose solution or solutions are desired. When there is only one independent variable, the problem
is one-dimensional, namely to find the root or roots of a function.

Except in linear problems, root finding invariably proceeds by iteration, and this is equally true
in one or in many dimensions. Starting from some approximate trial solution, a useful algorithm
will improve the solution until some predetermined convergence criterion is satisfied.

It cannot be overemphasized, however, how crucially success depends on having a good first guess
for the solution, especially for multidimensional problems. This crucial beginning usually depends
on analysis rather than numerics. Carefully crafted initial estimates reward you not only with
reduced computational effort, but also with understanding and increased self-esteem.

Bracketing and Bisection
------------------------
We will say that a root is "bracketed" in the interval (a, b) if f(a) and f(b) have opposite signs.
If the function is continuous, then at least one root must lie in that interval (the "intermediate
value theorem"). If the function is discontinuous, but bounded, then instead of a root there might
be a step discontinuity which crosses zero. For numerical purposes, that might as well be a root,
since the behavior is indistinguishable from the case of a continuous function whose zero crossing
occurs in between two "adjacent" floating-point numbers in a machine's finite-precision
representation. Only for functions with singularities is there the possibility that a bracketed
root is not really there.

Newton-Raphson Method Using Derivative
--------------------------------------
Perhaps the most celebrated of all one-dimensional root-finding routines is Newton's method, also
called the Newton-Raphson method. This method is distinguished from the methods of previous
sections by the fact that it requires the evaluation of both the function f(x), and the derivative
f'(x), at arbitrary points x. The Newton-Raphson formula consists geometrically of extending the
tangent line at a current point x(i) until it crosses zero, then setting the next guess x(i+1) to
the abscissa of that zero-crossing. Algebraically, the method derives from the familiar Taylor
series expansion of a function in the neighborhood of a point.

Newton-Raphson is not restricted to one dimension. The method readily generalizes to multiple
dimensions.

Far from a root, where the higher-order terms in the series are important, the Newton-Raphson
formula can give grossly inaccurate, meaningless corrections. For instance, the initial guess for
the root might be so far from the true root as to let the search interval include a local maximum
or minimum of the function. This can be death to the method. If an iteration places a trial guess
near such a local extremum, so that the first derivative nearly vanishes, then Newton-Raphson sends
its solution off to limbo, with vanishingly small hope of recovery.

Newton does not adjust bounds, and works only on local information at the point x. The bounds are
used only to pick the midpoint as the first guess, and to reject the solution if it wanders outside
of the bounds.

While Newton-Raphson's global convergence properties are poor, it is fairly easy to design a fail-
safe routine that utilizes a combination of bisection and Newton-Raphson. The hybrid algorithm
takes a bisection step whenever Newton-Raphson would take the solution out of bounds, or whenever
Newton-Raphson is not reducing the size of the brackets rapidly enough.
---------------------------------------------------------------------------------------------------
Using a combination of Newton-Raphson and bisection, find the root of a function bracketed between
x1 and x2. The root will be refined until its accuracy is known within +/-accurancy.
EvaluateGivenPoint is a user-supplied routine that returns both the function value and the first
derivative of the function.
***/
func (a *MathUtil) NewtonRaphsonBisection(userFunc func(pv, pmt, n, i float64, f, fPrime *float64) (),
                                          pv, pmt, n, x1, x2, accurancy float64) float64 {
  var maxIterations int = 100 //Maximum allowed number of iterations.
  var (
    fLow = zero
    fHigh = zero
    fPrime = zero
    xLow = zero
    xHigh = zero
  )
  userFunc(pv, pmt, n, x1, &fLow, &fPrime)
  userFunc(pv, pmt, n, x2, &fHigh, &fPrime)
  /***
  The principal difference between one and many dimensions is that, in one dimension, it is
  possible to bracket or "trap" a root between bracketing values, and then hunt it down like a
  rabbit. In multidimensions, you can never be sure that the root is there at all until you have
  found it.
  ***/
  if (fLow > 0.0 && fHigh > 0.0) || (fLow < 0.0 && fHigh < 0.0) {
    return(math.NaN()) //Root must be bracketed.
  } else if fLow == zero {
    return x1
  } else if fHigh == zero {
    return x2
  } else if fLow < zero { //Orient the search so that f(xLow) < 0.
    xLow = x1
    xHigh = x2
  } else {
    xHigh = x1
    xLow = x2
  }
  var (
    guess = 0.5 * (x1 + x2)            //Initialize the guess for root,
    fPrimePrevious = math.Abs(x2 - x1) //the "step-size before last,"
    fPrimeLast = fPrimePrevious        //and the last step.
    f = zero
    tmp = zero
  )
  userFunc(pv, pmt, n, guess, &f, &fPrime)
  for iteration := 0; iteration < maxIterations; iteration++ { //Loop over allowed iterations.
    //Bisect if Newton out of range, or not decreasing fast enough.
    if (((guess - xHigh) * fPrime - f) * ((guess - xLow) * fPrime - f)) > zero ||
       math.Abs(two * f) > math.Abs(fPrimePrevious * fPrime) {
      fPrimePrevious = fPrimeLast
      fPrimeLast = 0.5 * (xHigh - xLow)
      guess = xLow + fPrimeLast
      if xLow == guess { //Change in root is negligible. Newton step acceptable. Take it.
        return guess
      }
    } else {
      /***
                       f(x(i))
      x(i+1) = x(i) - ----------
                       f'(x(i))
      ***/
      fPrimePrevious = fPrimeLast
      fPrimeLast = f / fPrime
      tmp = guess
      guess -= fPrimeLast
      if tmp == guess {
        return guess
      }
    }
    //
    if math.Abs(fPrimeLast) < accurancy { //Convergence criterion.
      return guess
    }
    userFunc(pv, pmt, n, guess, &f, &fPrime) //The one new function evaluation per iteration.
    if f < zero { //Maintain the bracket on the root.
      xLow = guess
    } else {
      xHigh = guess
    }
  }
  return(math.NaN()) //Maximum number of iterations exceeded.
}
