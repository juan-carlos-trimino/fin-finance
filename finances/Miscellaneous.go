//
package finances

import (
  "math"
)

type Miscellaneous struct {
  P Periods
}

/***
                             Inflation and Interest Rates
Nominal (non-inflation-adjusted) vs. Real Interest Rates:
                        1 + Nominal Rate
  Real Interest Rate = -------------------- - 1
                        1 + Inflation Rate
***/
func (m *Miscellaneous) RealInterestRate(nominalRate, inflationRate float64) float64 {
  var realRate float64 = ((1.0 + nominalRate) / (1.0 + inflationRate)) - 1.0
  return realRate
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
func (m *Miscellaneous) GrowthDecayOfFunds(rate float64, c int,	factor float64) float64 {
  var n float64 = 0.0;
  if Continuously == c {
    n = math.Log(factor) / rate
  } else {
    rate /= float64(c)
    n = math.Log10(factor) / math.Log10(1.0 + rate)
  }
  return n
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

     EAR = (1 + NR / m)^m - 1

   In the EAR equation, NR/m is the periodic rate, and m is the number of periods per year. For
   example, suppose you could borrow using either a credit card that charges 1 percent per month or
   a bank loan with a 12 percent quoted nominal interest rate that is compounded quarterly. Which
   should you choose? To answer this question, the cost rate of each alternative must be expressed
   as an EAR:

     Credit card loan: EAR = (1 + 0.01)^12 - 1.0 = (1.01)^12 - 1.0
                           = 1.126825 - 1.0 = 0.126825 = 12.6825%.
            Bank loan: EAR = (1 + 0.03)^4 - 1.0 = (1.03)^4 - 1.0
                           = 1.125509 - 1.0 = 0.125509 = 12.5509%.

   Thus, the credit card loan is slightly more costly than the bank loan. This result should have
   been intuitive to you -- both loans have the same 12 percent nominal rate, yet you would have to
   make monthly payments on the credit card versus quarterly payments under the bank loan.

   The EAR rate is not used in calculations. However, it should be used to compare the effective
   cost or rate of return on loans or investments when payment periods differ, as in the credit
   card versus bank loan example.
***************************************************************************************************
                    Effective or Equivalent Annual Rate (EAR)
There are many situations where an interest rate with a specified compounding frequency must be
converted into an equivalent rate with a different compounding frequency. Examples include
situations where a need exists to compare alternative interest rates with different compounding
frequencies and where the payment frequency does not match the compounding frequency in an annuity
problem. The basic relationship used to convert interest rates from one compounding frequency to
another is shown below.

  EAR = (1 + r / m)^(m) - 1     or     r = m * ((1 + EAR)^(1/m) - 1)

  Continuous interest: (1 + r / m)^m = e^r    as m --> infinity
                       or
                       r = ln(EAR + 1)

  where
    r = nominal rate (NR),
    m = number of compounding periods per year.

The effective rate is an annually compounded interest rate that is equivalent to the nominal rate
compounded more frequently. The nominal rate is the stated rate in a problem, such as 5%,
compounded monthly. The number of periods per year is also stated in most problems. An interest
rate compounded monthly involves 12 periods per year, for example. Using the relationship shown
above, any effective annual rate can be converted to a rate compounded more frequently, and any
rate compounded more frequently than once a year can be converted to an effective annual rate.
***/
func (m *Miscellaneous) NominalToEffectiveAnnualRate(i float64, c int) float64 {
  var r float64 = 0.0
  if c == Continuously {
    r = math.Pow(math.E, i) - 1.0
  } else {
    r = math.Pow(1.0 + (*m).P.PeriodicInterestRate(i, c), float64(c)) - 1.0
  }
  return r
}

func (m *Miscellaneous) EffectiveAnnualToNominalRate(i float64, c int) float64 {
  var r float64 = 0.0
  if c == Continuously {
    r = math.Log(1.0 + i)
  } else {
    r = float64(c) * (math.Pow(1.0 + i, 1.0 / float64(c)) - 1.0)
  }
  return r
}
