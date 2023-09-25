// SimpleInterest computes simple interests (Ordinary Simple Interest, Commercial (Banker's)
// Interest, and Accurate (Exact) Interest).
package finances

import (
	"math"
)

/***
Parameter  Definition
---------  ----------

  p, PV     The present value of an amount of money.
  INT       Interest.
  n         A number of equal-time intervals between the present time and a future time (n for
            number of periods).
  i         Interest rate per period.
  cp        Compounding period.
  tp        Time period; i.e., years, months, days, etc.

Simple interest is the amount paid on money borrowed (principal) where the principal remains
unchanged for the period of time the money is in use.

There are three types of simple interest:
(a) Ordinary Simple Interest, based on "a 30-day month and a 360-day year," is frequently used for
    real estate loans, installment loans, and periodic repayment personal loans;
(b) Commercial (Banker's) Interest, based on "a 360-day year and an exactly specified number of
    days," results in the greatest return for the lender;
(c) Accurate (Exact) Interest, based on "a 365-day year and an exact number of days," is more
    frequently used in commercial transactions.

The basis for a loan determines the way in which time is calculated for computing simple interest.
The time calculation depends on whether a 30-day month or an exact number of days is used. In the
former, it is necessary to determine only the number of months over which the loan is made and to
multiply by 30 days. If there are days to either side of an even month, they are accounted for by
simple addition. Remember that 30-day-month calculations are used only for ordinary simple
interest.

The exact-time basis is used for Commercial or Banker's interest and requires determination of the
exact number of days during the life of a loan. The count usually INCLUDES the last day and
EXCLUDES the first day. If the month of February is included in the time period, an extra day must
be added when accounting for leap year.

The formula for computing simple interest is:

  interest = principal * interest rate * time

or

  INT = PV * i * n

where PV is the present value of the principal or amount loaned (or borrowed), INT is interest, i
is the interest rate per period, and n is the number of days the money is on loan. Since it is
common practice to specify the interest rate in terms of an annual rate, if the agreed-to
compounding period is one month, the annual rate must be divided by 12. If the compounding period
is in days, the annual rate must be divided by either 360 (Ordinary and Banker's interest) or 365
(Accurate interest).

To illustrate the difference between Banker's, Acccurate, and Ordinary interest, consider the case
of $10,000 loaned at 9% from June 1 to November 1.

  INT = 10,000 * 0.09 * (153 / 360) = 382.50 (Banker's)
  INT = 10,000 * 0.09 * (153 / 365) = 377.26 (Accurate)
  INT = 10,000 * 0.09 * (150 / 360) = 375.00 (Ordinary [Five 30-day months])
***/
type SimpleInterest struct {
  Periods  //Composition.
}

func (si *SimpleInterest) OrdinaryInterest(p, i float64, cp int, n float64, tp int) float64 {
  if cp == Continuously {
    return (math.NaN())
  } else if tp == Daily360 {
    m, _ := math.Modf(n / 30.0) //Break n into its fractional (ignore) and integer (m) parts.
    n = m * 30  //30-day months.
  }
  i /= float64(cp)
  n /= float64(tp)
  n *= float64(cp)
  return (p * i * n)
}

func (si *SimpleInterest) OrdinaryRate(p, amount float64, n float64, tp int) float64 {
  if tp == Daily360 {
    m, _ := math.Modf(n / 30.0) //Break n into its fractional (ignore) and integer (m) parts.
    n = m * 30  //30-day months.
  }
  n /= float64(tp)
  return (amount / (p * n))  //Per year.
}

func (si *SimpleInterest) OrdinaryPrincipal(amount, INT float64, cp int, n float64, tp int) float64 {
  if cp == Continuously {
    return (math.NaN())
  } else if tp == Daily360 {
    m, _ := math.Modf(n / 30.0) //Divide n into its fractional (ignore) and integer (m) parts.
    n = m * 30  //30-day months.
  }
  INT /= float64(cp)
  n /= float64(tp)
  n *= float64(cp)
  return (amount / (INT * n))
}

func (si *SimpleInterest) OrdinaryTime(p, amount, INT float64, cp int) float64 {
  if cp == Continuously {
    return (math.NaN())
  }
  amount *= float64(cp)
  return (amount / (INT * p))
}

func (si *SimpleInterest) BankersInterest(p, i float64, cp int, n float64, tp int) float64 {
  if cp == Continuously {
    return (math.NaN())
  }
  i /= float64(cp)
  n /= float64(tp)
  n *= float64(cp)
  return (p * i * n)
}

func (si *SimpleInterest) BankersRate(p, amount float64, n float64, tp int) float64 {
  n /= float64(tp)
  return (amount / (p * n))  //Per year.
}

func (si *SimpleInterest) BankersPrincipal(amount, INT float64, cp int, n float64, tp int) float64 {
  if cp == Continuously {
    return (math.NaN())
  }
  INT /= float64(cp)
  n /= float64(tp)
  n *= float64(cp)
  return (amount / (INT * n))
}

func (si *SimpleInterest) BankersTime(p, amount, INT float64, cp int) (n float64) {
  if cp == Continuously {
    return (math.NaN())
  }
  amount *= float64(cp)
  n = amount / (INT * p)
  return
}

func (si *SimpleInterest) AccurateInterest(p, i float64, cp int, n float64, tp int) float64 {
  if cp == Continuously {
    return (math.NaN())
  }
  i /= float64(cp)
  n /= float64(tp)
  n *= float64(cp);
  return (p * i * n)
}

func (si *SimpleInterest) AccurateRate(p, amount float64, n float64, tp int) float64 {
  n /= float64(tp)
  return (amount / (p * n))  //Per year.
}

func (si *SimpleInterest) AccuratePrincipal(amount, INT float64, cp int, n float64, tp int) float64 {
  if cp == Continuously {
    return (math.NaN())
  }
  INT /= float64(cp)
  n /= float64(tp)
  n *= float64(cp)
  return (amount / (INT * n))
}

func (si *SimpleInterest) AccurateTime(p, amount, INT float64, cp int) (n float64) {
  if cp == Continuously {
    return (math.NaN())
  }
  amount *= float64(cp)
  n = amount / (INT * p)
  return
}
