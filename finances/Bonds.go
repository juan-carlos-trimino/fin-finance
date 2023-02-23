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
By definition, the current price of a bond is the present value of all its cash flows.

Notes:
(1) Investors need to be aware of two main risks that can affect a bond's investment value: credit
    risk (default) and interest rate risk (interest rate fluctuations).
(2) As a bond yield decreases, its price rises at an increasing rate whereas a bond's price falls
    at a decreasing rate as its yield increases. This phenomenon is known as CONVEXITY.
(3) The discount rates used to determine the future value of expected coupon rates when yields rise
    or fall differ and, as a result, have different price impacts.
(4) A bond with more convexity will offer more upside if rates decrease, while promising less
    downside if rates increase.
(5) A par bond is one where the coupon is the same as the current market yield. If the bond pays a
    coupon that is higher than the market yield, its price will be higher than, or at a premium to,
    par. Conversely, if the bond's coupon is lower than the market yield, its price will be lower
    than, or at a discount to, par.
---------------------------------------------------------------------------------------------------
C - Coupon amount.
Coupon Rate - The stated interest rate paid to bondholders; the nominal yield.
FV - Face Value.
t - Number of time periods.
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
Important: The yield to call is widely deemed to be a more accurate estimate of expected return on
           a bond than the yield to maturity.

YIELD TO CALL (YTC) is a financial term that refers to the return a bondholder receives if the bond
is held until the call date, which occurs sometime before it reaches maturity. This number can be
mathematically calculated as the compound interest rate at which the present value of a bond's
future coupon payments and call price is equal to the current market price of the bond.

Yield to call applies to callable bonds, which are debt instruments that let bond investors redeem
the bonds -- or the bond issuer to repurchase them -- on what is known as the call date, at a price
known as the call price. By definition, the call date of a bond chronologically occurs before the
maturity date.

Generally speaking, bonds are callable over several years. They are normally called at a slight
premium above their face value, though the exact call price is based on prevailing market rates.
***/
func (b *Bonds) YieldToCall(FV, couponRate float64, cp int, timeToCall float64, tp int, bondPrice,
                            callPrice float64) float64 {
  var cashFlow = b.CashFlow(FV, couponRate, cp, timeToCall, tp)
  cashFlow[len(cashFlow) - 1] = callPrice + cashFlow[0]
  return(b.YieldToMaturity(cashFlow, bondPrice, cp))
}

/***
What is the internal rate of return (IRR) on the investment of buying the bond now and holding the
bond to maturity? The answer is the yield to maturity of a bond. IRR assumes reinvestment of coupon
at the bond yield.

Suppose you were offered a 14-year, 10% annual coupon, $1,000 par value bond at a price of
$1,494.93. What rate of interest would you earn on your investment if you bought the bond and held
it to maturity? This rate (5%) is called the bond's YIELD TO MATURITY (YTM); the YTM is identical
to the total rate of return.

Using bisection, find the root of the function CurrentPrice known to lie between 0.0 and 1.0. The
root, returned as r, will be refined until its accuracy is plus or minus ACCURACY.
***/
func (b *Bonds) YieldToMaturity(cashFlow []float64, bondPrice float64, cp int) (r float64) {
  /***
  Since the structure of cash flows is such that there exists only one solution to the equation,
  there is much less likelihood of having multiple solutions when doing this yield estimation for
  bonds.

  Since the bond yield is above zero, set the lower bound to zero. Then find an upper bound on the
  yield by increasing the interest rate until the bond price with this interest rate is negative.
  Finally, bisect the interval between the upper and lower bounds until the desired accuracy is
  obtained.
  ***/
  var bottom float64 = zero
  var top float64 = one
  for b.CurrentPrice(cashFlow, top, cp) > bondPrice {
    top *= 2.0
  }
  r = 0.5 * top
  var diff float64 = zero
  const MAX_BISECTION int = 200 //Maximum allowed number of bisections.
  const ACCURACY float64 = 1e-5
  for idx := 0; idx < MAX_BISECTION; idx++ { //Bisection loop.
    /***
    The bisection method must succeed. Over some interval the function is known to pass through
    zero because it changes sign. Evaluate the function at the interval's midpoint and examine its
    sign. Use the midpoint to replace whichever limit has the same sign. After each iteration the
    bounds containing the root decrease by a factor of two.
    ***/
    diff = b.CurrentPrice(cashFlow, r, cp) - bondPrice
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

/***
Duration measures how long it takes, IN YEARS, for an investor to be repaid the bond's price by the
bond's total cash flows. At the same time, duration is a measure of sensitivity of a bond's or
fixed income portfolio's price to changes in interest rates. In general, the higher the duration,
the more a bond's price will drop as interest rates rise (and the greater the interest rate risk).
As a general rule, for every 1% change in interest rates (increase or decrease), a bond's price
will change approximately 1% in the opposite direction, for every year of duration. If a bond has a
duration of five years and interest rates increase 1%, the bond's price will drop by approximately
5% (1% X 5 years). Likewise, if interest rates fall by 1%, the same bond's price will increase by
about 5% (1% X 5 years).

Certain factors can affect a bond's duration, including:
(1)	TIME TO MATURITY - The longer the maturity, the higher the duration, and the greater the
    interest rate risk. Consider two bonds that each yield 5% and cost $1,000, but have different
    maturities. A bond that matures faster - say, in one year - would repay its true cost faster
    than a bond that matures in 10 years. Consequently, the shorter-maturity bond would have a
    lower duration and less risk.
(2)	COUPON RATE - A bond's coupon rate is a key factor in the calculation of duration. If there are
    two bonds that are identical with the exception on their coupon rates, the bond with the higher
    coupon rate will pay back its original costs faster than the bond with a lower yield. The
    higher the coupon rate, the lower the duration, and the lower the interest rate risk.

The duration of a bond in practice can refer to two different things. The Macaulay duration is the
weighted average time until all the bond's cash flows are paid. By accounting for the present value
of future bond payments, the Macaulay duration helps an investor evaluate and compare bonds
independent of their term or time to maturity.

The second type of duration is called "modified duration" and, unlike Macaulay duration, is not
measured in years. Modified duration measures the expected change in a bond's price to a 1% change
in interest rates. In order to understand modified duration, keep in mind that bond prices are said
to have an inverse relationship with interest rates. Therefore, rising interest rates indicate that
bond prices are likely to fall, while declining interest rates indicate that bond prices are likely
to rise.
Notes:
(1) Unfortunately, duration has limitations when used as a measure of interest rate sensitivity.
    While the statistic calculates a linear relationship between price and yield changes in bonds,
    in reality, the relationship between the changes in price and yield is convex. Hence, the
    larger the change in interest rates, the larger the error in estimating the price change of the
    bond.
(2) When the bond is correctly priced, the duration and Macaulay duration will produce the same
    number.
***/
func (b *Bonds) Duration(cashFlow []float64, currentRate, bondPrice float64) float64 {
  var D float64 = zero
  currentRate /= hundred
  currentRate += one
  var sz int = len(cashFlow)
  for idx, t := 0, 1.0; idx < sz; idx, t = idx + 1, t + 1.0 {
    D += (t * cashFlow[idx]) / math.Pow(currentRate, t)
  }
  return(D / bondPrice)
}



/***
If the bond is priced correctly, the yield to maturity must equal the current interest rate. If current interest rate EQUALS yield to maturity the calculations from Duration and MacaulayDuration will produce the same number.
Notes:
(1) The longer the duration is the more sensitive the bond will be to changes in interest rates.
(2) For a standard bond the Macaulay duration will be between 0 and the maturity of the bond. It is equal to the maturity if and only if the bond is a zero-coupon bond.

FV = $100.00
coupon rate = 10%; compounding period annually
t = 3-year
current interest = 9%; compounding period annually

std::unique_ptr<Bonds> spBond = std::make_unique<Bonds>();
vector<double> cashFlow = spBond->CashFlow(100.00, 10.0, spBond->GetCompoundingPeriod(L'a'), 3.0, spBond->GetTimePeriod(L'y'));
double macyears = spBond->MacaulayDuration(cashFlow, spBond->GetCompoundingPeriod(L'a'), spBond->CurrentPrice(cashFlow, 9, spBond->GetCompoundingPeriod(L'a')));
cout << "macyears(2.738954) = " << macyears << endl;
***/
func (b *Bonds) MacaulayDuration(cashFlow []float64, cp int, bondPrice float64) float64 {
  var ytm = b.YieldToMaturity(cashFlow, bondPrice, cp)
  return(b.Duration(cashFlow, ytm, bondPrice))
}

/***
Modified duration is a formula that expresses the measurable change in the value of a security in response to a change in interest rates. Modified duration follows the concept that interest rates and bond prices move in opposite directions. This formula is used to determine the effect that a 100-basis point (1%) change in interest rates will have on the price of a bond.

How to interpret the result below? The modified duration illustrates the effect of a 100-basis point (1%) change in interest rates on the price of a bond. Therefore,

* If interest rates increase by 1%, the price of the 3-year bond will decrease by 2.513%.
* If interest rates decrease by 1%, the price of the 3-year bond will increase by 2.513%.

The modified duration provides a good measurement of a bond's sensitivity to changes in interest rates.

FV = $100.00
coupon rate = 10%; compounding period annually
t = 3-year
current interest = 9%; compounding period annually

std::unique_ptr<Bonds> spBond = std::make_unique<Bonds>();
vector<double> cashFlow = spBond->CashFlow(100.00, 10.0, spBond->GetCompoundingPeriod(L'a'), 3.0, spBond->GetTimePeriod(L'y'));
double MDuration = spBond->ModifiedDuration(cashFlow, 9, spBond->GetCompoundingPeriod(L'a'), spBond->CurrentPrice(cashFlow, 9, spBond->GetCompoundingPeriod(L'a')));
cout << "MDuration(2.512801%) = " << MDuration << endl;
***/
func (b *Bonds) ModifiedDuration(cashFlow []float64, cp int, bondPrice float64) float64 {
  var ytm = b.YieldToMaturity(cashFlow, bondPrice, cp)
  return(b.Duration(cashFlow, ytm, bondPrice) / (one + (ytm / hundred)))
}


/***
Duration is a linear measure or 1st derivative of how the price of a bond changes in response to interest rate changes. As interest rates change, the price is not likely to change linearly, but instead it would change over some curved function of interest rates. The more curved the price function of the bond is, the more inaccurate duration is as a measure of the interest rate sensitivity.

(Duration can be a good measure of how bond prices may be affected due to small and sudden fluctuations in interest rates. However, the relationship between bond prices and yields is typically more sloped, or convex. Therefore, convexity is a better measure for assessing the impact on bond prices when there are large fluctuations in interest rates.)

Convexity is a measure of the curvature or 2nd derivative of how the price of a bond varies with interest rate; i.e., how the duration of a bond changes as the interest rate changes. Specifically, one assumes that the interest rate is constant across the life of the bond and that changes in interest rates occur evenly. Using these assumptions, duration can be formulated as the first derivative of the price function of the bond with respect to the interest rate in question. Then the convexity would be the second derivative of the price function with respect to the interest rate.

In actual markets, the assumption of constant interest rates and even changes is not correct, and more complex models are needed to actually price bonds. However, these simplifying assumptions allow one to quickly and easily calculate factors which describe the sensitivity of the bond prices to interest rate changes.

Convexity does not assume the relationship between bond value and interest rates to be linear. For large fluctuations in interest rates, it is a better measure than duration.
Notes:
(1) It's important to know how bond prices and market interest rates relate to one another. As interest rates fall, bond prices rise. Conversely, rising market interest rates lead to falling bond prices. This opposite reaction is because as rates rise, the bond may fall behind in the earnings they may offer a potential investor in comparison to other securities.
(2) If a bond's duration increases as yields increase, the bond is said to have negative convexity. In other words, the bond price will decline by a greater rate with a rise in yields than if yields had fallen. Therefore, if a bond has negative convexity, its duration would increase - the price would fall. As interest rates rise, the opposite is true.
(3) If a bond's duration rises and yields fall, the bond is said to have positive convexity. In other words, as yields fall, bond prices rise by a greater rate - or duration - than if yields rose. Positive convexity leads to greater increases in bond prices. If a bond has positive convexity, it would typically experience larger price increases as yields fall, compared to price decreases when yields increase.
(4) Zero-coupon bonds have the highest degree of convexity because they do not offer any coupon payments.

FV = $100.00
coupon rate = 10%; compounding period annually
t = 3-year
current interest = 9%; compounding period annually

std::unique_ptr<Bonds> spBond = std::make_unique<Bonds>();
vector<double> cashFlow = spBond->CashFlow(100.00, 10.0, spBond->GetCompoundingPeriod(L'a'), 3.0, spBond->GetTimePeriod(L'y'));
double Cx = spBond->Convexity(cashFlow, 9, spBond->GetCompoundingPeriod(L'a'));
cout << "Cx(8.932479) = " << Cx << endl;
***/
func (b *Bonds) Convexity(cashFlow []float64, currentRate float64, cp int) float64 {
  var price float64 = b.CurrentPrice(cashFlow, currentRate, cp)
  var Cx float64 = zero
  currentRate /= hundred
  currentRate += one
  var sz int = len(cashFlow)
  for idx, t := 0, 1.0; idx < sz; idx, t = idx + 1, t + 1.0 {
    Cx += ((one + t) * t * cashFlow[idx]) / math.Pow(currentRate, t)
  }
  Cx /= price
  return(Cx / math.Pow(currentRate, 2))
}

/***
When using continuously compounded interest, one does not need the concept of modified duration.
When the bond is correctly priced, the duration and Macaulay duration will produce the same number.

FV = $100.00
coupon rate = 10% (compounding period is annually)
t = 3-year
current interest = 9% (compounding period is continuously)

std::unique_ptr<Bonds> spBond = std::make_unique<Bonds>();
vector<double> cashFlow = spBond->CashFlow(100.00, 10.0, spBond->GetCompoundingPeriod(L'a'), 3.0, spBond->GetTimePeriod(L'y'));
double price = spBond->CurrentPriceContinuous(cashFlow, 9);
cout << "Price(101.463758) = " << price << endl;
double duration = spBond->DurationContinuous(cashFlow, 9, spBond->CurrentPriceContinuous(cashFlow, 9));
cout << "duration(2.737529) = " << duration << endl;
double ytm = spBond->YieldToMaturityContinuous(cashFlow, spBond->CurrentPriceContinuous(cashFlow, 9));
cout << "ytm(9.000000) = " << ytm << endl;
double macdur = spBond->MacaulayDurationContinuous(cashFlow, spBond->CurrentPriceContinuous(cashFlow, 9));
cout << "Macaulay Duration(2.737529) = " << macdur << endl;
double convexity = spBond->ConvexityContinuous(cashFlow, 9, spBond->CurrentPriceContinuous(cashFlow, 9));
cout << "convexity(7.867793) = " << convexity << endl;
price = spBond->CurrentPriceContinuous(cashFlow, 8);
cout << "Price(104.281666 (new rate is 8%)) = " << price << endl;
***/
func (b *Bonds) CurrentPriceContinuous(cashFlow []float64, currentRate float64) (price float64) {
  currentRate /= hundred
  price = zero
  var sz int = len(cashFlow)
  for idx, t := 0, 1.0; idx < sz; idx, t = idx + 1, t + 1.0 {
    price += math.Exp(-currentRate * t) * cashFlow[idx]
  }
  return
}

func (b *Bonds) DurationContinuous(cashFlow []float64, currentRate, bondPrice float64) float64 {
  var D float64 = zero
  currentRate /= hundred
  var sz int = len(cashFlow)
  for idx, t := 0, 1.0; idx < sz; idx, t = idx + 1, t + 1.0 {
    D += math.Exp(-currentRate * t) * t * cashFlow[idx]
  }
  return(D / bondPrice)
}

func (b *Bonds) YieldToMaturityContinuous(cashFlow []float64, bondPrice float64) (r float64) {
  /***
  Since the structure of cash flows is such that there exists only one solution to the equation, there is much less likelihood of having multiple solutions when doing this yield estimation for bonds.

  Since the bond yield is above zero, set the lower bound to zero. Then find an upper bound on the yield by increasing the interest rate until the bond price with this interest rate is negative. Finally, bisect the interval between the upper and lower bounds until the desired accuracy is obtained.
  ***/
  var bottom float64 = zero
  var top float64 = one
  for b.CurrentPriceContinuous(cashFlow, top) > bondPrice {
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
    diff = b.CurrentPriceContinuous(cashFlow, r) - bondPrice
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

func (b *Bonds) MacaulayDurationContinuous(cashFlow []float64, bondPrice float64) float64 {
  var ytm = b.YieldToMaturityContinuous(cashFlow, bondPrice)
  return(b.DurationContinuous(cashFlow, ytm, bondPrice))
}

func (b *Bonds) ConvexityContinuous(cashFlow []float64, currentRate, bondPrice float64) float64 {
  var Cx float64 = zero
  currentRate /= hundred
  var sz int = len(cashFlow)
  for idx, t := 0, 1.0; idx < sz; idx, t = idx + 1, t + 1.0 {
    Cx += t * t * cashFlow[idx] * math.Exp(-currentRate * t)
  }
  return(Cx / bondPrice)
}
