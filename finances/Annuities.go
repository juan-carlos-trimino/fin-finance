//
package finances

import (
  "math"
  "finance/mathutil"
)

const (
  zero float64 = 0.0
  one float64 = 1.0
  two float64 = 2.0
  hundred float64 = 100.0
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
  mathutil.MathUtil
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
      ln(----)
           PV
 n = -----------
      ln(1 + i)
***/
func (a *Annuities) O_Periods_PV_FV(PV, FV, i float64, cp int) (n float64) {
  i = a.PeriodicInterestRate(i, cp)
  n = math.Log(FV / PV) / math.Log(one + i)
  return
}

/***
               i * PV
      ln((1 - --------)^(-1))
                 PMT
 n = -------------------------
            ln(1 + i)
***/
func (a *Annuities) O_Periods_PMT_PV(PMT, PV, i float64, cp int) (n float64) {
  i = a.PeriodicInterestRate(i, cp)
  n = math.Log(math.Pow(one - ((i * PV) / PMT), -1)) / math.Log(one + i)
  return
}

/***
              i * FV
      ln(1 + --------)
               PMT
 n = ------------------
         ln(1 + i)
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
             1 - (1 + i)^(-n)
 PV = PMT * ------------------ * (1 + i)
                    i
***/
func (a *Annuities) D_PresentValue_PMT(PMT, i float64, cp int, n float64, tp int) (PV float64) {
  PV = a.O_PresentValue_PMT(PMT, i, cp, n, tp)
  i = a.PeriodicInterestRate(i, cp)
  PV *= (one + i)
  return
}

/***
             (1 + i)^n - 1
 FV = PMT * --------------- * (1 + i)
                  i
***/
func (a *Annuities) D_FutureValue_PMT(PMT, i float64, cp int, n float64, tp int) (FV float64) {
  FV = a.O_FutureValue_PMT(PMT, i, cp, n, tp)
  i = a.PeriodicInterestRate(i, cp)
  FV *= (one + i)
  return
}

/***
                   i                 1
 PMT = PV * ------------------ * ---------
             1 - (1 + i)^(-n)     (1 + i)
***/
func (a *Annuities) D_Payment_PV(PV, i float64, cp int, n float64, tp int) (PMT float64) {
  PMT = a.O_Payment_PV(PV, i, cp, n, tp)
  i = a.PeriodicInterestRate(i, cp)
  PMT /= (one + i)
  return
}

/***
                  i               1
 PMT = FV * --------------- * ---------
             (1 + i)^n - 1     (1 + i)
***/
func (a *Annuities) D_Payment_FV(FV, i float64, cp int, n float64, tp int) (PMT float64) {
  PMT = a.O_Payment_FV(FV, i, cp, n, tp)
  i = a.PeriodicInterestRate(i, cp)
  PMT /= (one + i)
  return
}

func (a *Annuities) D_Periods_PMT_PV(PMT, PV, i float64, cp int) float64 {
  return(a.O_Periods_PMT_PV(PMT, PV, i, cp))
}

func (a *Annuities) D_Periods_PMT_FV(PMT, FV, i float64, cp int) float64 {
  return(a.O_Periods_PMT_FV(PMT, FV, i, cp))
}

/***
Exponential decay is geometric growth with a negative rate of growth.
---------------------------------------------------------------------
Let i represent the rate of inflation; what costs $1 now will cost $(1 + i) this time next year.
For example, if the inflation rate were i = 25%, then what costs $1 now would cost $1.25 this time
next year. A dollar next year would buy only 0.8 (=1/1.25) times as much as a dollar buys today. In
other words, a dollar next year would be worth only $0.80 in today's dollars -- by next year, a
dollar would have lost 20% of its purchasing power. We say that the present value of a dollar next
year would be $0.80. Notice that although the inflation rate is 25%, the loss in purchasing power
is 20%. In other words, a dollar a year from now is worth $[1 - i/(i + 1)] today, and the loss in
purchasing power is the fraction i/(i + 1). The quantity -i(1 + i) behaves like a negative interest
rate. We can use the compound interest formula to find the value of PV dollars n years from now as

 FV = PV(1 + r)^n = PV(1 - (i / (i + 1)))^n

 The actual posted price of an item, at any time, is said to be in current dollars. That price can
 be compared with prices at other times by converting all prices to constant dollars, dollars of a
 particular year.
***/
func (a *Annuities) Depreciation(PV, i float64, cp int, n float64, tp int) (FV float64) {
  i = a.PeriodicInterestRate(i, cp)
  n = a.NumberOfPeriods(n, tp, float64(Daily365), cp)
  FV = PV * math.Pow(one - (i / (i + one)), n)
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
func (a *Annuities) BlendedInterestRate(p1, r1, p2, r2 float64) (avg float64) {
  var tb float64 = p1 + p2
  avg = (r1 * (p1 / tb)) + (r2 * (p2 / tb))
  return
}

/***
                              Inflation and Interest Rates
                              ----------------------------
 Nominal (non-inflation-adjusted) vs. Real Interest Rates:
                        1 + Nominal Rate
  Real Interest Rate = -------------------- - 1
                        1 + Inflation Rate
***/
func (a *Annuities) RealInterestRate(nominalRate, inflationRate float64) (realRate float64) {
  realRate = ((one + nominalRate) / (one + inflationRate)) - one
  return
}

/***
  Exponential Growth (positive rate of change) and Decay (negative rate of change)

  FV = PV * (1 + i)^n
  FV / PV = F = (1 + i)^n
  log(F) = n * log(1 + i)
  n = log(F) / log(1 + i)

  FV = PV * e^(n * i)
  FV / PV = F = e^(n * i)
  ln(F) = n * i * ln(e)
  n = ln(F) / i
***/
func (a *Annuities) GrowthDecayOfFunds(factor, rate float64, cp int) (n float64) {
  n = zero
  if Continuously == cp {
    n = math.Log(factor) / rate
  } else {
    rate /= float64(cp)
    n = math.Log10(factor) / math.Log10(one + rate)
  }
  return
}

/***
                                 Types of Interest Rates
Compounding involves three types of interest rates: nominal rates (NR), periodic rates (PR), and
effective annual rates (EAR).

1. Nominal, or quoted, rate (NR). This is the rate that is quoted by banks, brokers, and other
   financial institutions. So, if you talk with a banker, broker, mortgage lender, auto finance
   company, or student loan officer about rates, the nominal rate is the one he or she will
   normally quote you. However, to be meaningful, the quoted nominal rate must also include the
   number of compounding periods per year. For example, a bank might offer 6 percent, compounded
   quarterly, on CDs, or a mutual fund might offer 5 percent, compounded monthly, on its money
   market account.

   The nominal rate on loans to consumers is also called the Annual Percentage Rate (APR). For
   example, if a credit card issuer quotes an annual rate of 18 percent, this is the APR.

2. Periodic Rate (PR). This is the rate charged by a lender or paid by a borrower each period. It
   can be a rate per year, per six-month period, per quarter, per month, per day, or per any other
   time interval. For example, a bank might charge 1.5 percent per month on its credit card loans,
   or a finance company might charge 3 percent per quarter on installment loans. We find the
   periodic rate as follows:

     PR = NR / m

   which implies that

     Nominal annual rate = NR = (PR)(m).

   Here m is the number of compounding periods per year. To illustrate, consider a finance company
   loan at 3 percent per quarter:

     NR = (PR)(m) = (3%)(4) = 12%,

     or

     PR = NR / m = (12%) / 4 = 3% per quarter.

   If there is only one payment per year, or if interest is added only once a year, then m = 1, and
   the periodic rate is equal to the nominal rate.

3. Effective (or equivalent) annual rate (EAR). This is the annual rate that produces the same
   result as if we had compounded at a given periodic rate m times per year. The EAR is found as
   follows:

     EAR = (1 + (NR / m))^(m) - 1

   In the EAR equation, NR/m is the periodic rate, and m is the number of periods per year. For
   example, suppose you could borrow using either a credit card that charges 1% per month or a bank
   loan with a 12% quoted nominal interest rate that is compounded quarterly. Which should you
   choose? To answer this question, the cost rate of each alternative must be expressed
   as an EAR:

     Credit card loan: EAR = (1 + 0.01)^12 - 1.0 = (1.01)^12 - 1.0
                           = 1.126825 - 1.0 = 0.126825 = 12.6825%.
           Bank loan: EAR = (1 + 0.03)^4 - 1.0 = (1.03)^4 - 1.0
                          = 1.125509 - 1.0 = 0.125509 = 12.5509%.

   Thus, the credit card loan is slightly more costly than the bank loan. This result should have
   been intuitive to you -- both loans have the same 12% nominal rate, yet you would have to make
   monthly payments on the credit card versus quarterly payments under the bank loan.

   The EAR rate is not used in calculations. However, it should be used to compare the effective
   cost or rate of return on loans or investments when payment periods differ, as in the credit
   card versus bank loan example.

                    Effective or Equivalent Annual Rate (EAR)
There are many situations where an interest rate with a specified compounding frequency must be
converted into an equivalent rate with a different compounding frequency. Examples include
situations where a need exists to compare alternative interest rates with different compounding
frequencies and where the payment frequency does not match the compounding frequency in an annuity
problem. The basic relationship used to convert interest rates from one compounding frequency to
another is shown below.

  EAR = (1 + (r / m))^(m) - 1     or     r = m * ((1 + EAR)^(1/m) - 1)

  Continuous interest: EAR = e^(r) - 1
                       r = ln(EAR + 1)

  where
    r = nominal rate,
    m = number of compounding periods per year.

The effective rate is an annually compounded interest rate that is equivalent to the nominal rate
compounded more frequently. The nominal rate is the stated rate in a problem, such as 5%,
compounded monthly. The number of periods per year is also stated in most problems. An interest
rate compounded monthly involves 12 periods per year, for example. Using the relationship shown
above, any effective annual rate can be converted to a rate compounded more frequently, and any
rate compounded more frequently than once a year can be converted to an effective annual rate.
***/
func (a *Annuities) NominalToEAR(i float64, cp int) (ear float64) {
  ear = zero
  if cp == Continuously {
    ear = math.Pow(math.E, i) - one
  } else {
    ear = math.Pow(one + a.PeriodicInterestRate(i, cp), float64(cp)) - one
  }
  return
}

func (a *Annuities) EARToNominal(ear float64, cp int) (r float64) {
  r = zero
  if cp == Continuously {
    r = math.Log(ear + one)
  } else {
    r = float64(cp) * (math.Pow(one + ear, one / float64(cp)) - one)
  }
  return
}

/***
                              GEOMETRIC MEAN RETURN
The geometric mean is used to determine the average compound growth rate (or average rate of
return) over a given period. The geometric mean is the most accurate method for determining average
rates of return.

                    THE REASON WE NEED THE GEOMETRIC RETURN
To see why we need to use the geometric mean return rather than an arithmetic return, consider the
following example. Assume that the price of a stock is $100. Then, one year later, the price of the
stock has fallen to $50. However, in the following year, the price rises again to $100. What is the
average rate of return per year on the stock for the two-year period?

If we use an arithmetic return, we would say that the stock fell in price 50% during the first year
and rose 100% in price the second year. Therefore, the average rate of return must be 25% per year,
i.e. (-50% + 100%) divided by 2 years.

Obviously, this is not correct. The average rate of return is zero percent per year since the price
ended where it began two years earlier. So we need a better method to calculate the average rate of
return -- that method is the geometric mean return.

                                  METHODOLOGY
To calculate the geometric mean return, we use a five-step process:
1. Determine the rate of return for each time period.
2. Add one to each of the returns (the result is called a holding period return).
3. Multiply each of the holding period returns together (this is called "chain-linking" the
   returns).
4. Take the root of the product in step 3. The root number is equal to the number of time periods.
5. Subtract one from the result.

                                  AN EXAMPLE
For example, assume that we have already calculated the return on a common stock for each of the
last four years. We used this equation to measure the rate of return for each year:

            (Price1 - Price0) + Dividend
  Return = ------------------------------
                      Price0

Let us assume that those returns are as follows:
  1st year: 5.0%
  2nd year: -3.0%
  3rd year: 12.0%
  4th year: 10.0%

Steps 2-4 instruct us to add one to each value, take the product, and then take the nth root, where
n is the number of time periods.
      Annual Return        Holding Period Return
      -------------        ---------------------
          5.0%                  1.05
         -3.0%                  0.97
         12.0%                  1.12
         10.0%                  1.10
                     ---------------------------
                     Product = 1.254792

      4th root = (1.05 * 0.97 * 1.12 * 1.10)^1/4
               = (1.254792)^1/4
               = 1.058383

(If your calculator has a y^x (i.e. y to the x) function, the above calculation is very easy to do;
e.g., y = 1.254792 and x = 1/4 or 0.25)

Finally, subtract one from the result.
  Geometric Mean = 1.058383 - 1
                 = 0.058383
                 = 5.84% (rounded to two decimal places)

Now, let's look at another example, but add one more column: a column which shows the value of
$1.00 invested at the beginning of the period.
                       Holding Period  Value of $1 at the
      Year     Return      Return      end of the period
    -----------------------------------------------------
    1st year     2.0%       1.02         $1.020000
    2nd year     8.0%       1.08          1.101600
    3rd year    -1.0%       0.99          1.090584
    4th year    10.0%       1.10          1.199642

Notice the ending value ($1.199642) at the end of the 4-year period.

Now forget that procedure for a moment and let us solve for the geometric mean return using the
same data:

  Geometric Mean = (1.02 * 1.08 * 0.99 * 1.10)^1/4 - 1
                 = (1.199642)^1/4 - 1
                 = 0.046557
                 = 4.7%

              A SHORTCUT USING THE BEGINNING AND ENDING VALUES
Notice that the value under the root symbol (1.199642) is the same as the value that a dollar has
grown to at the end of the 4-year period. This fact often gives us a shortcut to solving for the
geometric mean. If we only know the value at the end of the period and the value at the beginning
of the period, we can divide the ending value by the beginning value. This gives us the value that
a dollar invested at the beginning of the period would have grown to by the end of the period.
Taking the root of this and subtracting one gives us the average rate of growth per period.

For example, the current price of a stock is $15.00 per share. You expect the stock to double in
price over the next five years. The stock pays no dividend. What is the average rate of growth per
year?
                           $30
  Ending Value of $1.00 = ----- = 2.00
                           $15

  Average rate of return = (2.00)^1/5 - 1
                         = 0.1487
                         = 14.87%

The geometric mean thus offers us a very short, convenient, and accurate method of calculating the
average compound rate of return over a given period.
***/
func (a *Annuities) AverageRateOfReturn(v [] float64) (d float64) {
  d = zero
  var sz int = len(v)
  if sz != 0 { //Is slice empty?
    d = one + (v[0] / hundred)
    for idx := 1; idx < sz; idx++ {
      d *= (one + (v[idx] / hundred))
    }
    d = math.Pow(d, one / float64(sz)) - one
  }
  return
}

/***
     PMT
i = ----- * ((1 + i)^n - 1)
     FV

        FV
f(x) = ----- * i * (1 + i)^-n + (1 + i)^-n - 1
        PMT


..find 1st derivate for this
f'(x) = TBD

***/





/***
  const char* const pLocale = "english_usa.1252";
  locale::global(locale(pLocale));
  cout << std::setfill(' ') << std::fixed << std::left << std::showpoint;
  std::cin.imbue(locale());  //Register global locale.
  cout.imbue(locale());

  std::unique_ptr<Annuities> spA = std::make_unique<Annuities>();
  int cp = spA->GetCompoundingPeriod(L'm');
  double i = spA->O_Interest_PV_PMT(60, 24000, 500, cp, 1.0, 31.0);
  cout << "i(0.7628634% per month) = " << (i * 100) << endl;
  cout << "i(9.154323% per year) = " << (i * 100 * cp) << endl;

  i = spA->O_Interest_PV_PMT(48, 11200, 291, cp, 4.0, 12.0);
  cout << "i(0.94007411% per month) = " << (i * 100) << endl;
  cout << "i(11.28% per year) = " << (i * 100 * cp) << endl;

  cp = spA->GetCompoundingPeriod(L'a');
  i = spA->O_Interest_PV_PMT(5, 50000, 13500, cp, 10.0, 15.0);
  cout << "i(10.91616% per year) = " << (i * 100) << endl;
  cout << "i(10.91616% per year) = " << (i * 100 * cp) << endl;

     PMT
i = ----- * (1 - (1 + i)^-n)
     PV

For an efficient realization of Newton-Raphson the user provides a routine that evaluates both f(x)
and its first derivative f'(x) at the point x.

         PV
 f(x) = ----- * i * (1 + i)^n - (1 + i)^n + 1
         PMT

          PV
 f'(x) = ----- * ((1 + i)^n + n * i * (1 + i)^(n - 1)) - n * (1 + i)^(n - 1)
          PMT
***/
func /*(a *Annuities)*/ evaluateGivenPoint(pv, pmt, n, i float64, f, fPrime *float64) () {
  var pvDivByPmt = pv / pmt
  var i_To_n = math.Pow(one + i, n)
  var i_To_n_minus_1 = math.Pow(one + i, n - 1)
  *f = (pvDivByPmt * i * i_To_n) - i_To_n + one
  *fPrime = pvDivByPmt * (i_To_n + (n * i * i_To_n_minus_1)) - (n * i_To_n_minus_1)
  return
}

func (a *Annuities) O_Interest_PV_PMT(pv, pmt, n float64, cp int, x1, x2, accurancy float64) (i float64) {
  var mu mathutil.MathUtil
  x1 = a.PeriodicInterestRate(x1 / hundred, cp)
  x2 = a.PeriodicInterestRate(x2 / hundred, cp)
  i = mu.NewtonRaphsonBisection(evaluateGivenPoint, pv, pmt, n, x1, x2, accurancy)
  return
}





