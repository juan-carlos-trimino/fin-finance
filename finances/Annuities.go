//
package finances

import (
  "math"
)

type Annuities struct {
  Periods
}

/***
             1 - (1 + i)^(-n)
 PV = PMT * ------------------
                    i
***/
func (a *Annuities) O_PresentValue_PMT(PMT, i float64, cp int, n float64, tp int) float64 {
  i = a.PeriodicInterestRate(i, cp)
  n = a.NumberOfPeriods(n, tp, cp, float64(Daily365))
  var PV float64 = PMT * ((one - math.Pow(one + i, -n)) / i)
  return PV
}

/***
                   i
 PMT = PV * ------------------
             1 - (1 + i)^(-n)
***/
func (a *Annuities) O_Payment_PV(PV, i float64, cp int, n float64, tp int) float64 {
  i = a.PeriodicInterestRate(i, cp)
  n = a.NumberOfPeriods(n, tp, cp, float64(Daily365))
  var PMT float64 = PV * (i / (one - math.Pow(one + i, -n)))
  return PMT
}






/***
When two or more loans on a single property are outstanding, investors must calculate the average
interest rate.

          P1                 P2
 (R1 * ---------) + (R2 * ---------) = Average Rate
        P1 + P2            P1 + P2

 P1 = Principal Balance, Loan1                  P2 = Principal Balance, Loan2
 R1 = Interest Rate, Loan 1                     R2 = Interest Rate, Loan 2

 Example: Loan 1 ---> rate = 12.5%, balance = $90,000
          Loan 2 ---> rate = 7.25%, balance = $32,000
          Average rate ---> 11.12%

To figure out whether it would pay to refinance your mortgage and Home-Equity Line of Credit
(HELOC) with one load, you need to calculate a blended interest rate for your total housing debt.
If the result is higher than what you could get on a new fixed-rate mortgage, consider it.

                   Mortgage Balance                    HELOC Balance     Blended
 (Mortgage Rate * ------------------) + (HELOC Rate * ---------------) = Interest
                     Total Debt                         Total Debt       Rate

 Note: Total debt is mortgage balance plus HELOC balance.
 Example:  Mortgage ---> rate = 6.5%, balance = $200,000
           HELOC ---> rate = 10.5%, balance = $100,000
           Blended interest rate ---> 7.8%
***/
func (a *Annuities) BlendedInterestRate(p1, r1, p2, r2 float64) float64 {
  var tb float64 = p1 + p2
  var avg float64 = (r1 * (p1 / tb)) + (r2 * (p2 / tb))
  return avg
}
















