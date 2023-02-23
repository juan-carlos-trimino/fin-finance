//
package finances

import (
  "math"
	// "fmt"
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
func (b *Bonds) CurrentPrice(cashFlow []float64, currentRate float64, cp int) (price float64) {
  currentRate /= hundred
  currentRate = b.PeriodicInterestRate(currentRate, cp)
  price = zero
  currentRate += one
  var sz int = len(cashFlow)
  for idx, t := 0, 1.0; idx < sz; idx, t = idx + 1, t + 1.0 {
    price += cashFlow[idx] / math.Pow(currentRate, t)
  }
  return
}



/***
Important: The yield to call is widely deemed to be a more accurate estimate of expected return on a bond than the yield to maturity.

YIELD TO CALL (YTC) is a financial term that refers to the return a bondholder receives if the bond is held until the call date, which occurs sometime before it reaches maturity. This number can be mathematically calculated as the compound interest rate at which the present value of a bond's future coupon payments and call price is equal to the current market price of the bond.

Yield to call applies to callable bonds, which are debt instruments that let bond investors redeem the bonds -- or the bond issuer to repurchase them -- on what is known as the call date, at a price known as the call price. By definition, the call date of a bond chronologically occurs before the maturity date.

Generally speaking, bonds are callable over several years. They are normally called at a slight premium above their face value, though the exact call price is based on prevailing market rates.

FV = $1,000.00
coupon rate = 10%; annually
time to call = 9 years
bond price = $1,494.93
call price = $1,100.00

FV = $1,000.00
coupon rate = 10%; semiannually
time to call = 5 years
bond price = $1,175.00
call price = $1,100.00

std::unique_ptr<Bonds> spBond = std::make_unique<Bonds>();
double ytc = spBond->YieldToCall(1000.00, 10.0, spBond->GetCompoundingPeriod(L'a'), 9, spBond->GetTimePeriod(L'y'), 1494.93, 1100.00);
cout << "ytc(4.21%) = " << ytc << endl;
ytc = spBond->YieldToCall(1000.00, 10.0, spBond->GetCompoundingPeriod(L's'), 5, spBond->GetTimePeriod(L'y'), 1175.00, 1100.00);
cout << "ytc(7.43%) = " << ytc << endl;
***/
func (b *Bonds) YieldToCall(FV, couponRate float64, cp int, timeToCall float64, tp int, bondPrice, callPrice float64) float64 {
  var cashFlow = b.CashFlow(FV, couponRate, cp, timeToCall, tp)
  cashFlow[len(cashFlow) - 1] = callPrice + cashFlow[0]
  return(b.YieldToMaturity(cashFlow, bondPrice, cp))
}



/***
What is the internal rate of return (IRR) on the investment of buying the bond now and holding the bond to maturity? The answer is the yield to maturity of a bond. IRR assumes reinvestment of coupon at the bond yield.

Suppose you were offered a 14-year, 10% annual coupon, $1,000 par value bond at a price of $1,494.93. What rate of interest would you earn on your investment if you bought the bond and held it to maturity? This rate (5%) is called the bond�s YIELD TO MATURITY (YTM); the YTM is identical to the total rate of return.

Using bisection, find the root of the function CurrentPrice known to lie between 0.0 and 1.0. The root, returned as r, will be refined until its accuracy is plus or minus ACCURACY.
-------------------------------------------------------------------------------
FV = $100.00
coupon rate = 10%; compounding period annually
t = 3-year
current interest = 9%; compounding period annually

std::unique_ptr<Bonds> spBond = std::make_unique<Bonds>();
vector<double> cashFlow = spBond->CashFlow(100.00, 10.0, spBond->GetCompoundingPeriod(L'a'), 3.0, spBond->GetTimePeriod(L'y'));
double ytm = spBond->YieldToMaturity(cashFlow, spBond->CurrentPrice(cashFlow, 9.0, spBond->GetCompoundingPeriod(L'a')), spBond->GetCompoundingPeriod(L'a'));
cout << "ytm(9.00%) = " << ytm << endl;
-------------------------------------------------------------------------------
FV = $1,000.00
coupon rate = 10%; compounding period annually
t = 10-year
price = $920.00

std::unique_ptr<Bonds> spBond = std::make_unique<Bonds>();
vector<double> cashFlow = spBond->CashFlow(1000.00, 10.0, spBond->GetCompoundingPeriod(L'a'), 10, spBond->GetTimePeriod(L'y'));
double ytm = spBond->YieldToMaturity(cashFlow, 920.00, spBond->GetCompoundingPeriod(L'a'));
cout << "ytm(11.3801%) = " << ytm << endl;
-------------------------------------------------------------------------------
FV = $100.00
coupon rate = 5%; compounding period semiannually
t = 30-month
price = $95.92

std::unique_ptr<Bonds> spBond = std::make_unique<Bonds>();
vector<double> cashFlow = spBond->CashFlow(100.00, 5.0, spBond->GetCompoundingPeriod(L's'), 30, spBond->GetTimePeriod(L'm'));
double ytm = spBond->YieldToMaturity(cashFlow, 95.92, spBond->GetCompoundingPeriod(L's'));
cout << "ytm(6.80223%) = " << ytm << endl;
***/
func (b *Bonds) YieldToMaturity(cashFlow []float64, bondPrice float64, compoundingPeriod int) (r float64) {
  /***
  Since the structure of cash flows is such that there exists only one solution to the equation, there is much less likelihood of having multiple solutions when doing this yield estimation for bonds.

  Since the bond yield is above zero, set the lower bound to zero. Then find an upper bound on the yield by increasing the interest rate until the bond price with this interest rate is negative. Finally, bisect the interval between the upper and lower bounds until the desired accuracy is obtained.
  ***/
  var bottom float64 = zero
  var top float64 = one
  for b.CurrentPrice(cashFlow, top, compoundingPeriod) > bondPrice {
    top *= 2.0
  }
  r = 0.5 * top
  var diff float64 = zero
  const MAX_BISECTION int = 200 //Maximum allowed number of bisections.
  const ACCURACY float64 = 1e-5
  for idx := 0; idx < MAX_BISECTION; idx++ { //Bisection loop.
    /***
    The bisection method must succeed. Over some interval the function is known to pass through zero because it changes sign. Evaluate the function at the interval's midpoint and examine its sign. Use the midpoint to replace whichever limit has the same sign. After each iteration the bounds containing the root decrease by a factor of two.
    ***/
    diff = b.CurrentPrice(cashFlow, r, compoundingPeriod) - bondPrice
    if ACCURACY > math.Abs(diff) {
      break
    } else if diff > zero {
      bottom = r
    } else {
      top = r
    }
    r = 0.5 * (top + bottom)
  }
  return
}

