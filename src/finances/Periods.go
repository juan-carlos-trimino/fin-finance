// Periods adjusts compounding periods and time (interest) periods.
package finances

import (
  "math"
)

const (
  //Compounding
  Invalid int = -1
  Annually int = 1
  SemiAnnually int = 2
  Quarterly int = 4
  Monthly int = 12
  Weekly int = 52
  Daily int = -2
  Daily360 int = 360
  Daily365 int = 365
  Continuously int = -100
  //Time Periods
  Years int = 1
  Semiyears int = 2
  Quarters int = 4
  Months int = 12
  Weeks int = 52
  Days int = -2
  //
  semiYearlyPerYear float64 = 2.0
  quartersPerYear float64 = 4.0
  monthsPerYear float64 = 12.0
  weeksPerYear float64 = 52.0
  quartersPerSemiYearly float64 = 2.0
  weeksPerSemiYearly float64 = 26.0
  daysPerSemiYearly float64 = 180.0
  monthsPerSemiYearly float64 = 6.0
  monthsPerQuarter float64 = 3.0
  weeksPerQuarter float64 = 12.0
  weeksPerMonth float64 = 12.0
  daysPerQuarter float64 = 90.0
  daysPerMonth float64 = 30.0
  daysPerWeek float64 = 7.0
)

type Periods struct{}

/***
Adjust if compound period (c) does not equal number of interest periods (t); e.g.,
(1) n is 20, t is years, and c is monthly, then n = 20 * 12 = 240 months.
(2) n is 6.5, t is years, and c is quarterly, then n = 6.5 * 4 = 26 quarters.
***/
func (p Periods) numberOfPeriods(n float64, tp int, forDaysOnly float64, cp int) float64 {
  if tp != cp {
    switch cp {
      case Monthly:
        if tp == Years {
          return (n * monthsPerYear)
        } else if tp == Quarters {
          return (n * monthsPerQuarter)
        } else if tp == Weeks {
          return (n / weeksPerMonth)
        } else if tp == Days {
          return (n / daysPerMonth)
        } else {
          return (math.NaN())
        }
      case Annually:
        if tp == Months {
          return (n / monthsPerYear)
        } else if tp == Quarters {
          return (n / quartersPerYear)
        } else if tp == Weeks {
          return (n / weeksPerYear)
        } else if tp == Days {
          return (n / forDaysOnly)
        } else {
          return (math.NaN())
        }
      case SemiAnnually:
        if tp == Months {
          return (n / monthsPerSemiYearly)
        } else if tp == Years {
          return (n * semiYearlyPerYear)
        } else if tp == Quarters {
          return (n / quartersPerSemiYearly)
        } else if tp == Weeks {
          return (n / weeksPerSemiYearly)
        } else if tp == Days {
          return (n / daysPerSemiYearly)
        } else {
          return (math.NaN())
        }
      case Quarterly:
        if tp == Years {
          return (n * quartersPerYear)
        } else if tp == Months {
          return (n / monthsPerQuarter)
        } else if tp == Weeks {
          return (n / weeksPerQuarter)
        } else if tp == Days {
          return (n / daysPerQuarter)
        } else {
          return (math.NaN())
        }
      case Weekly:
        if tp == Months {
          return (n * weeksPerMonth)
        } else if tp == Years {
          return (n * weeksPerYear)
        } else if tp == Quarters {
          return (n * weeksPerQuarter)
        } else if tp == Days {
          return (n / daysPerWeek)
        } else {
          return (math.NaN())
        }
      case Daily, Daily360, Daily365:
        if tp == Months {
          return (n * daysPerMonth)
        } else if tp == Years {
          return (n * forDaysOnly)
        } else if tp == Quarters {
          return (n * daysPerQuarter)
        } else if tp == Weeks {
          return (n * daysPerWeek)
        } else {
          return (math.NaN())
        }
    }
  } else {
    return n
  }
  return (n * float64(cp))
}

/***
A periodic interest rate is a rate that can be charged on a loan, or realized on an investment over
a specific period of time. Lenders typically quote interest rates on an annual basis, but the
interest compounds more frequently than annually in most cases. The periodic interest rate is the
annual interest rate divided by the number of compounding periods.

Example of a Periodic Interest Rate
-----------------------------------
The interest on a mortgage is compounded or applied on a monthly basis. If the annual interest rate
on that mortgage is 8%, the periodic interest rate used to calculate the interest assessed in any
single month is 0.08 divided by 12, working out to 0.0067 or 0.67%.

The remaining principal balance of the mortgage loan would have a 0.67% interest rate applied to it
each month.
***/
func (p Periods) periodicInterestRate(interestRate float64, compoundingPeriods int) float64 {
  return (interestRate / float64(compoundingPeriods))
}

func (p Periods) GetCompoundingPeriod(compoundingPeriod byte, isDaily365 bool) int {
  /***
  Cases are evaluated from top to bottom, so the first matching one is executed. The optional
  default case matches if none of the other cases does; it may be placed anywhere. Cases do not
  fall through from one to the next as in C-like languages (though there is a rarely used
  fallthrough statement that overrides this behavior).
  ***/
  switch compoundingPeriod {
    case 'm', 'M':  //(M)onthly
      return Monthly
    case 'a', 'A':  //(A)nnually
      return Annually
    case 's', 'S':  //(S)emiannually
      return SemiAnnually
    case 'q', 'Q':  //(Q)uarterly
      return Quarterly
    case 'w', 'W':  //(W)eekly
      return Weekly
    case 'd', 'D':  //(D)aily
      if isDaily365 == true {
        return Daily365
      } else {
        return Daily360
      }
    case 'c', 'C':  //(C)ontinuously
      return Continuously
    default:
      return Invalid
  }
}

func (p Periods) GetTimePeriod(timePeriod byte, isDaily365 bool) int {
  switch timePeriod {
    case 'm', 'M':  //(M)onths
      return Months
    case 'y', 'Y':  //(Y)ears
      return Years
    case 's', 'S':  //(S)emiyears
      return Semiyears
    case 'q', 'Q':  //(Q)uarters
      return Quarters
    case 'w', 'W':  //(W)eeks
      return Weeks
    case 'd', 'D':  //(D)ays
      if isDaily365 == true {
        return Daily365
      } else {
        return Daily360
      }
    default:
      return Invalid
  }
}
