// Mortgage computes information about mortgages.
package finances

/***
(1) The struct type with no fields is called the empty struct, written struct{}. It has size zero
    and carries no information.
(2) For efficiency, larger struct types are usually passed to or returned from functions indirectly
    using a pointer, and this is required if the function must modify its argument, since in a
    call-by-value language like Go, the called function receives only a copy of an argument, not a
    reference to the original argument.
(3) Go lets us declare a field with a type but no name; such fields are called "anonymous" fields.
    The type of the field must be a named type or a pointer to a named type.
(4) Because "anonymous" fields do have implicit names, you can't have two anonymous fields of the
    same type since their names would conflict. And because the name of the field is implicitly
    determined by its type, so too is the visibility of the field.
***/
type Mortgage struct {
  Annuities
  Periods
}

func (m *Mortgage) CostOfMortgage(mortgage, i float64, compoundingPeriod byte, n float64, timePeriod byte) (payment, totalCost, totalInterest float64) {
  var cp int = (*m).GetCompoundingPeriod(compoundingPeriod, true)
  var tp int = (*m).GetTimePeriod(timePeriod, true)
  payment = (*m).O_Payment_PV(mortgage, i, cp, n, tp)
  totalCost = payment * (*m).NumberOfPeriods(n, tp, float64(Daily365), cp)
  totalInterest = totalCost - mortgage
  return
}

/***
Refinance mortgage and HELOC with one load.
If the blended interest rate is higher than what you could get on a new fixed-rate mortgage,
consider it.
***/
func (m *Mortgage) MortgageHeloc(mortgageBalance, mortgageRate, helocBalance, helocRate float64) (blendedInterestRate float64) {
  blendedInterestRate = m.BlendedInterestRate(mortgageBalance, mortgageRate, helocBalance, helocRate) * hundred
  return
}

type row struct { //Rows for the amortization table.
  Payment, PmtPrincipal, PmtInterest, Balance float64
}

type AmortizationTable struct {
  Payment, TotalCost, TotalInterest float64
  Rows []row
}

/***
                              Amortization Table

  Loan Amount: $300,000.00
  Term of the Loan: 360.00 month(s).
  i (%): 3.375% monthly.
  Payment: $1,326.29
  Total Interest: $177,463.91

    Payment                       Payment Applied to:         Declining
      No.          Payment      Principal       Interest       Balance
   ---------------------------------------------------------------------
        -                -              -              -     300,000.00
        1         1,326.29         482.54         843.75     299,517.46
        2         1,326.29         483.90         842.39     299,033.57
        3         1,326.29         485.26         841.03     298,548.31
      ...
      358         1,326.29       1,315.16          11.13       2,641.43
      359         1,326.29       1,318.86           7.43       1,322.57
      360         1,326.29       1,322.57           3.72           0.00
   ---------------------------------------------------------------------
                477,463.91     300,000.00     177,463.91
***/
func (m *Mortgage) AmortizationTable(mortgage, i float64, compoundingPeriod byte, n float64, timePeriod byte) (at AmortizationTable) {
  var payment, totalCost, totalInterest = m.CostOfMortgage(mortgage, i, compoundingPeriod, n, timePeriod)
  var cp int = m.GetCompoundingPeriod(compoundingPeriod, true)
  var tp int = m.GetTimePeriod(timePeriod, true)
  var periods int = int(m.NumberOfPeriods(n, tp, float64(Daily365), cp))
  var rows = make([]row, 0, periods)
  var balance, pmtPrincipal, pmtInterest float64 = zero, zero, zero
  periods--
  tp = cp
  for pmtNumber := 1; periods > -1; periods-- {
    balance = m.O_PresentValue_PMT(payment, i, cp, float64(periods), tp)
    pmtPrincipal = mortgage - balance
    pmtInterest = payment - pmtPrincipal
    rows = append(rows, row{payment, pmtPrincipal, pmtInterest, balance})
    mortgage -= pmtPrincipal
    pmtNumber++
  }
  at = AmortizationTable{Payment: payment, TotalCost: totalCost, TotalInterest: totalInterest, Rows: rows}
  return
}
