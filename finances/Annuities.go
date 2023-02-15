//
package finances

import (
  "math"
)

/***
                                       Future Value
                                       ------------
The process of going from today's values, or present values (PVs), to future values (FVs) is called
COMPOUNDING.

Compounding Process - The interest rate is, in fact, a growth rate: If a sum is deposited and earns
5 percent interest, then the funds on deposit will grow at a rate of 5 percent per period. Note
also that time value concepts can be applied to anything that is growing -- sales, population,
earnings per share, or your future salary.
                                       -------------
                                       Present Value
                                       -------------
Finding present values is called DISCOUNTING, and it is the reverse of compounding -- if you know
the PV, you can compound to find the FV, while if you know the FV, you can discount to find the PV.

Discounting Process - The present value of $1 (or any other sum) to be received in the future
diminishes as the years to receipt and the interest rate increase.
                                         -------
                                         Annuity
                                         -------
An annuity is a series of equal payments made at fixed intervals for a specified number of periods.
The payments can occur at either the beginning or the end of each period. If the payments occur at
the end of each period, as they typically do, the annuity is called ordinary, or deferred, annuity.
Payments on mortgages, car loans, and student loans are typically set up as ordinary annuities. If
payments are made at the beginning of each period, the annuity is an annuity due. Rental payments
for an apartment, life insurance premiums, and lottery payoffs are typically set up as annuities
due.
***/
type Annuities struct {
  Periods
}

/***
Parameter  Definition
---------  ----------------------------------------------------------------------------------------
 PV        Present value of money.
 FV        Future value of money.
 n         Number of compounding periods.
 i         Interest rate per compounding period.
 PMT       Amount of equal periodic payments when payments are made at the end of the compounding
           period.

 The NOMINAL RATE is the stated rate of interest on which any compounding is based.

 The EFFECTIVE RATE is the actual percentage increase in the account. It is also called the
 EQUIVALENT YIELD since an account with this nominal rate would earn the same increase without
 compounding.

 The nominal rate for a year is called the ANNUAL (PERCENTAGE) RATE (APR), and the effective rate
 is called the ANNUAL (EQUIVALENT) YIELD or ANNUAL PERCENTAGE YIELD (APY).

 GEOMETRIC GROWTH (also called EXPONENTIAL GROWTH) is growth proportional to the amount present.

 FV = PV * (1 + i)^n
    = PV * e^(n * i)
***/
func (a *Annuities) O_FutureValue_PV(PV, i float64, cp int, n float64, tp int) (FV float64) {
  FV = zero
  if cp == Continuously {
    n /= float64(tp)
    n *= i
    FV = PV * math.Pow(math.E, n)
  } else {
    i = a.PeriodicInterestRate(i, cp)
    n = a.NumberOfPeriods(n, tp, float64(Daily365), cp)
    FV = PV * math.Pow(one + i, n)
  }
  return
}

/***
             (1 + i)^n - 1
 FV = PMT * ---------------
                  i
***/
func (a *Annuities) O_FutureValue_PMT(PMT, i float64, cp int, n float64, tp int) (FV float64) {
  i = a.PeriodicInterestRate(i, cp)
  n = a.NumberOfPeriods(n, tp, float64(Daily365), cp)
  FV = PMT * ((math.Pow(one + i, n) - one) / i)
  return
}

/***
          FV
PV = -----------
      (1 + i)^n
***/
func (a *Annuities) O_PresentValue_FV(FV, i float64, cp int, n float64, tp int) (PV float64) {
  i = a.PeriodicInterestRate(i, cp)
  n = a.NumberOfPeriods(n, tp, float64(Daily365), cp)
  PV = FV / math.Pow(one + i, n)
  return
}

/***
             1 - (1 + i)^(-n)
 PV = PMT * ------------------
                    i
***/
func (a *Annuities) O_PresentValue_PMT(PMT, i float64, cp int, n float64, tp int) (PV float64) {
  i = a.PeriodicInterestRate(i, cp)
  n = a.NumberOfPeriods(n, tp, float64(Daily365), cp)
  PV = PMT * ((one - math.Pow(one + i, -n)) / i)
  return
}

/***
       FV
 i = (----)^(1/n) - 1
       PV
***/
func (a *Annuities) O_Interest_PV_FV(PV, FV, n float64, tp, cp int) float64 {
  n = a.NumberOfPeriods(n, tp, float64(Daily365), cp)
  var i float64 = math.Pow(FV / PV, one / n) - one
  return(i * float64(cp))
}

/***
                  i
 PMT = FV * ---------------
             (1 + i)^n - 1
***/
func (a *Annuities) O_Payment_FV(FV, i float64, cp int, n float64, tp int) (PMT float64) {
  i = a.PeriodicInterestRate(i, cp)
  n = a.NumberOfPeriods(n, tp, float64(Daily365), cp)
  PMT = FV * (i / (math.Pow(one + i, n) - one))
  return
}

/***
                   i
 PMT = PV * ------------------
             1 - (1 + i)^(-n)
***/
func (a *Annuities) O_Payment_PV(PV, i float64, cp int, n float64, tp int) (PMT float64) {
  i = a.PeriodicInterestRate(i, cp)
  n = a.NumberOfPeriods(n, tp, float64(Daily365), cp)
  PMT = PV * (i / (one - math.Pow(one + i, -n)))
  return
}

/***
           FV
      log(----)
           PV
 n = ------------
      log(1 + i)
***/
func (a *Annuities) O_Periods_PV_FV(PV, FV, i float64, cp int) (n float64) {
  i = a.PeriodicInterestRate(i, cp)
  n = math.Log(FV / PV) / math.Log(one + i)
  return
}

/***
                   1
      log(----------------------)
           1 + ((i * PV) / PMT)
 n = -----------------------------
            log(1 + i)
***/
func (a *Annuities) O_Periods_PMT_PV(PMT, PV, i float64, cp int) (n float64) {
  i = a.PeriodicInterestRate(i, cp)
  n = math.Log(one / (one + ((i * PV) / PMT))) / math.Log(one + i)
  return
}

/***
               i * FV
      log(1 + --------)
                PMT
 n = -------------------
         log(1 + i)
***/
func (a *Annuities) O_Periods_PMT_FV(PMT, FV, i float64, cp int) (n float64) {
  i = a.PeriodicInterestRate(i, cp)
  n = math.Log(one + ((i * FV) / PMT)) / math.Log(one + i)
  return
}

/***
Some annuities go on indefinitely, or perpetually, and these are called perpetuities. The present
value of a perpetuity is

                      Payment        PMT
PV (Perpetuity) = --------------- = ----- (Annual Compounding)
                   Discount rate      d

                          PMT
                = --------------------- (Compounded (m) Times per Year)
                   [(1 + d)^(1/m) - 1]

The value of a perpetuity changes dramatically when interest rates change.
***/
func (a *Annuities) O_Perpetuity(d, PMT float64, cp int) (PV float64) {
  PV = math.NaN()
  if Continuously != cp {
    PV = PMT / (math.Pow(one + d, one / float64(cp)) - one)
  }
  return
}

/***
Some annuities go on indefinitely, or perpetually, and these are called perpetuities. The present
value of a growing perpetuity is

                               Payment               PMT
PV (Perpetuity) = ----------------------------- = -------- (Annual Compounding)
                   Discount rate - growth rate      d - g

                               PMT
         = ------------------------------------------- (Compounded (m) Times per Year)
            [(1 + d)^(1/m) - 1] - [(1 + g)^(1/m) - 1]

It is important to note that the discount rate must be higher than the growth rate when using the
present value of a growing perpetuity formula. This is due to the present value of a growing
perpetuity formula being an infinite geometric series. In theory, if the growth rate is higher than
the discount rate, the growing perpetuity would have an infinite value.

Example:
The expected dividend next year is $1.30, and dividends are expected to grow at 5% forever. If the
discount rate is 10%, what is the value of this promised dividend stream?
Answer: $26.00
***/
func (a *Annuities) O_GrowingPerpetuity(d, g, PMT float64, cp int) (PV float64) {
  PV = math.NaN()
  if Continuously != cp {
    PV = PMT / ((math.Pow(one + d, one / float64(cp)) - one) - (math.Pow(one + g, one / float64(cp)) - one))
  }
  return
}

/***
Growing annuities can be defined as cash flows, equally spaced in time, which grow at a constant
rate over the life of the payment stream and end after a predetermined number of payment or
withdrawal periods.

         (1 + i)^n  -  (1 + g)^n
FV = C (-------------------------)
                 i  -  g

Where:
FV = Future Value of Growing Annuity
 C = Initial deposit, assuming each deposit is made at the end of each period (ordinary annuity).
 i = Interest rate or yield earned on investment.
 g = Growth rate of deposits.
 n = Number of deposits.
***/
func (a *Annuities) O_GrowingAnnuityFutureValue(C, n, g, i float64, cp int) (FV float64) {
  i = a.PeriodicInterestRate(i, cp)
  g = a.PeriodicInterestRate(g, cp)
  FV = C * ((math.Pow(one + i, n) - math.Pow(one + g, n)) / (i - g))
  return
}

/***
        C           1 + g
PV = ------- (1 - (-------)^n)
      i - g         1 + i
***/
func (a *Annuities) O_GrowingAnnuityPresentValue(C, n, g, i float64, cp int) (PV float64) {
  i = a.PeriodicInterestRate(i, cp)
  g = a.PeriodicInterestRate(g, cp)
  PV = (C / (i - g)) * (one - math.Pow((one + g) / (one + i), n))
  return
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
















