//
package finances

import (
  "math"
	"fmt"
)

type Bonds struct {
  Periods
  // mathutil.MathUtil
}

/***
The tax-equivalent yield is the rate an investor would have to get on a taxable bond to match the
tax-free interest on a muni.
***/
func (b *Bonds) TaxableVsTaxFreeYields(taxFreeYield, cityTaxRate, stateTaxRate,
	                                     federalTaxRate float64) float64 {
  // 1. Total your city and state tax rates, expressed as a decimal number.
  cityTaxRate += stateTaxRate
  cityTaxRate /= hundred
  // 2. Multiply the result by 1 minus your federal tax rate.
  federalTaxRate /= hundred
  cityTaxRate *= (one - federalTaxRate)
  // 3. Add the result to your federal tax rate.
  federalTaxRate += cityTaxRate
  // 4. Subtract the sum from 1.
  federalTaxRate = one - federalTaxRate
  // 5. Divide the result into the tax-free yield, expressed as a decimal.
  taxFreeYield /= hundred
  return(taxFreeYield / federalTaxRate)
}

/***
A bond's cash flow is determined by calculating the coupon rate multiplied by the face value. A
$1,000 corporate bond with a 3.0% coupon has an annual cash flow of $30. If it's a 10-year bond
that has five years left until maturity, there would be five coupon payments remaining.

Payment 1 - $30; Payment 2 - $30; and so on.

The final payment would include the face value: $1,000 + $30 = $1,030.

This is important because the closer the bond is to maturity, the higher its value may be.

FV (Face Value)
cp (Compound period or coupon frequency)
***/
func (b *Bonds) CashFlow(FV, couponRate float64, cp int, n float64, tp int) (cashFlow [] float64) {
  couponRate /= hundred
  couponRate = b.PeriodicInterestRate(couponRate, cp)
  var C = FV * couponRate
  var sz int = int(b.NumberOfPeriods(n, tp, float64(Daily365), cp))
  fmt.Printf("szzzzzzzzzzz = %d\n", sz)
  cashFlow = make([]float64, sz, sz) //len=cap=sz
  for idx := range cashFlow {
    cashFlow[idx] = C
  }
  cashFlow[sz - 1] += FV
  return
}






/***
Investors need to be aware of two main risks that can affect a bond's investment value: credit risk (default) and interest rate risk (interest rate fluctuations).

By definition, the current price of a bond is the present value of all its cash flows.
Notes:
(1) As a bond yield decreases, its price rises at an increasing rate whereas a bond's price falls at a decreasing rate as its yield increases. This phenomenon is known as convexity.
(2) The discount rates used to determine the future value of expected coupon rates when yields rise or fall differ and, as a result, have different price impacts.
(3) A bond with more convexity will offer more upside if rates decrease, while promising less downside if rates increase.
(4) A par bond is one where the coupon is the same as the current market yield. If the bond pays a coupon that is higher than the market yield, its price will be higher than, or at a premium to, par. Conversely, if the bond's coupon is lower than the market yield, its price will be lower than, or at a discount to, par.

C - Coupon amount.
Coupon Rate - The stated interest rate paid to bondholders; the nominal yield.
FV - Face Value.
t - Number of time periods.
-------------------------------------------------------------------------------
***/
func (b *Bonds) CurrentPrice(cashFlow []float64, currentRate float64, compoundingPeriod int) (price float64) {
  currentRate /= hundred
  currentRate = b.PeriodicInterestRate(currentRate, compoundingPeriod)
  price = zero
  currentRate += one
  var sz int = len(cashFlow)
  for idx, t := 0, 1.0; idx < sz; idx, t = idx + 1, t + 1.0 {
    price += cashFlow[idx] / math.Pow(currentRate, t)
  }
  return
}




