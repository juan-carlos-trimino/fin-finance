//Bond pricing with a flat term structure.
package finances

import (
  "math"
)

type Bonds struct {
  /***
  In Go, a struct field is called embedded if it's declared without a name; embedding is about
  composition, not inheritance. Embedding is used to promote the fields and methods of an embedded
  type; the promoted fields and methods are accessible from two different paths.
  ***/
  Periods
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
The bond is a promise to pay a face value F at the maturity date T periods from now. Each period
the bond pays a fixed percentage amount of the face value as coupon C. The cash flows from the bond
are as follows:

 t =               0     1       2       3    ...     T
--------------------------------------------------------------
 Coupon                  C       C       C    ...     C
 Face Value                                           F
--------------------------------------------------------------
Total Cash Flows      C1 = C  C2 = C  C3 = C  ...  CT = C + F

===================================================================================================

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
  couponRate = b.periodicInterestRate(couponRate, cp)
  var C = FV * couponRate
  var sz int = int(b.numberOfCouponPaymentPeriods(n, tp, cp))
  //tp != cp; e.g., tp = semiyear & cp = annually 
  if sz == 0 {
    sz = 1
  }
  //The capacity argument may be omitted, in which case the capacity equals the length.
  cashFlow = make([]float64, sz) //len=cap=sz
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
  currentRate = b.periodicInterestRate(currentRate, cp)
  price = zero
  currentRate += one
  var sz int = len(cashFlow)
  for idx, t := 0, 1.0; idx < sz; idx, t = idx + 1, t + 1.0 {
    price += cashFlow[idx] / math.Pow(currentRate, t)
  }
  return
}

/***
When using continuously compounded interest, one does not need the concept of modified duration.
When the bond is correctly priced, the duration and Macaulay duration will produce the same number.
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
The current yield of a bond calculates the rate of return on a bond by using the current market
price of the bond instead of its face value. It is calculated as the annual coupon payment divided
by the current market price. The current yield is an accurate measure of bond yield as it reflects
the market sentiment and investor expectations from the bond in terms of return.
***/
func (b *Bonds) CurrentYield(annualCoupon, FV, currentPrice float64) (cy float64) {
  cy = ((annualCoupon / hundred) * FV) / currentPrice
  return
}

/***
Yield to Maturity (YTM) is the overall interest rate earned by an investor who buys a bond at the
market price and holds it until maturity. Mathematically, it is the discount rate at which the sum
of all future cash flows (from coupons and principal repayment) equals the price of the bond. YTM
is often quoted in terms of an annual rate and may differ from the bond's coupon rate. It assumes
that coupon and principal payments are made on time. It does not require dividends to be
reinvested, but computations of YTM generally make that assumption. Further, it does not consider
taxes paid by the investor or brokerage costs associated with the purchase.

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
  for idx := 0; idx < MAX_BISECTION; idx++ { //Bisection loop.
    /***
    The bisection method must succeed. Over some interval the function is known to pass through
    zero because it changes sign. Evaluate the function at the interval's midpoint and examine its
    sign. Use the midpoint to replace whichever limit has the same sign. After each iteration the
    bounds containing the root decrease by a factor of two.
    ***/
    diff = b.CurrentPrice(cashFlow, r, cp) - bondPrice
    if Accuracy > math.Abs(diff) {
      break
    } else if diff > zero {
      bottom = r
    } else {
      top = r
    }
    r = 0.5 * (top + bottom)
  }
  /***
  The yield to maturity is the interest rate that makes the present value of the future coupon
  payments equal to the current bond price.
  ***/
  return
}

func (b *Bonds) YieldToMaturityContinuous(cashFlow []float64, bondPrice float64) (r float64) {
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
  for b.CurrentPriceContinuous(cashFlow, top) > bondPrice {
    top *= 2.0
  }
  r = 0.5 * top
  var diff float64 = zero
  const MAX_BISECTION int = 200 //Maximum allowed number of bisections.
  for idx := 0; idx < MAX_BISECTION; idx++ { //Bisection loop.
    /***
    The bisection method must succeed. Over some interval the function is known to pass through
    zero because it changes sign. Evaluate the function at the interval's midpoint and examine its
    sign. Use the midpoint to replace whichever limit has the same sign. After each iteration the
    bounds containing the root decrease by a factor of two.
    ***/
    diff = b.CurrentPriceContinuous(cashFlow, r) - bondPrice
    if Accuracy > math.Abs(diff) {
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
Discount Rate
Discount rate is the rate of return used to discount future cash flows when calculating an
investment's present value. A discount rate is applied to future cash flows because money earned in
the future is less valuable than money earned today. This is based on the principle that money
should make more money over time -- a concept known as the "time value of money".

Duration is an important measure for investors to consider, as bonds with higher durations (given
equal credit, inflation and reinvestment risk) may have greater price volatility than bonds with
lower durations.

Investment theory tells us that the value of a fixed-income investment is the sum of all of its
cash flows discounted at an interest rate that reflects the inherent investment risk. In addition,
due to the time value of money, it assumes that cash flows returned earlier are worth more than
cash flows returned later. In its most basic form, duration measures the weighted average of the
present value of the cash flows of a fixed-income investment.

All of the components of a bond -- price, coupon, maturity, and interest rates -- are used in the
calculation of its duration. Although a bond's price is dependent on many variables apart from
duration, duration can be used to determine how the bond's price may react to changes in interest
rates.

The price of a bond, or any fixed-income investment, is determined by summing the cash flows
discounted by a rate of return. The rate of return can change at any time period and will be
reflected in the calculation of an investment's market price. (The sensitivity of a bond's value to
changing interest rates depends on both the length of time to maturity and on the pattern of cash
flows provided by the bond.)

Macaulay Duration measures the number of years required to recover the true cost of a bond,
considering the present value of all coupon and principal payments received in the future. Thus, it
is the only type of duration quoted in YEARS.

Modified Duration expands or modifies Macaulay duration to measure the responsiveness of a bond's
price to interest rate changes. It is defined as the percentage change in price for a 100 basis
point (1%) change in interest rates. The formula assumes that the cash flows of the bond do not
change as interest rates change (which is not the case for most callable bonds).

Macaulay Duration Formula

         n
       -----
       \          (PV)(CF ) * t
        \	               t
MacD =   \    ----------------------
         /     Market Price of Bond
        /
       /
       -----
        t=1

(PV)(CF ) = Present value of coupon at period t
       t
t = Time to each cash flow (in years)
n = Number of periods to maturity.

Macaulay duration is a measure of the time until the present value of cash flows from an investment
equals the investment's cost. It is calculated by taking the weighted average of the present values
of the cash flows, with the weights being the time until each cash flow is received. Macaulay
duration is often used to compare the risk of different investments. An investment with a longer
Macaulay duration is considered to be riskier because it takes longer for the investment to pay
off; it can also be used to calculate the interest rate sensitivity of an investment. An investment
with a longer Macaulay duration will be more sensitive to changes in interest rates. If interest
rates are at 7% annually, a 3-year bond (face value of $1,000) with a 10% coupon paid annually
would sell for:

Market Price = $100/(1.07)^1 + $100/(1.07)^2 + $1100/(1.07)^3
             = $93.46 + $87.34 + $897.93
             = $1,078.73

               (1 * $93.46 / $1,078.73) +
               (2 * $87.34 / $1,078.73) +
               (3 * $897.93 / $1,078.73)
MacD         = 2.7458
(It takes 2.7458 years to recover the true cost of the bond.)

Note that the duration of a zero-coupon bond equals the maturity of the bond, while the duration of
a coupon bond is less than the maturity because coupons are paid throughout the life of the bond.

The Modified duration is an extension of Macaulay duration because it takes into account interest
rate movements by including the frequency of coupon payments per year.

           Macaulay Duration
ModD = -------------------------
             Yield to maturity
	      1 + -------------------
             Number of coupon
             periods per year

Using the Macaulay example above, yield to maturity is assumed to be 7 percent, there is 1 coupon
period per year, and the Macaulay duration is 2.7458.

       2.7458 / (1 + (0.7 / 1))
       2.7458 / 1.07
ModD = 2.566
(For every 1 percent change in market interest rates, the market value of the bond will move
inversely by 2.566%)

As used in the equations for duration, coupon rate (which determines the size of the periodic cash
flow), interest rates (which determines the present value of the periodic cash flow), and maturity
(which weights each cash flow) all contribute to the duration.

As maturity increases, duration increases and the bond's price becomes more sensitive to interest
rate changes.
  * A decrease in maturity decreases duration and renders the bond less sensitive to changes in
    market yield. Therefore, duration varies directly with maturity.

As the bond coupon increases, its duration decreases and the bond becomes less sensitive to
interest rate changes.
  * Increases in coupon rates raise the present value of each periodic cash flow and therefore the
    market price. This higher market price lowers the duration.

As interest rates increase, duration decreases and the bond becomes less sensitive to further rate
changes.
  * As interest rates increase, all of the net present values of the future cash flows decline as
    their discount factors increase, but the cash flows that are farthest away will show the
    largest proportional decrease. So the early cash flows will have a greater weight relative to
    later cash flows. As yields (interest rates) decline, the opposite will occur.

Convexity
One of the limitations of duration as a measure of interest rate/price sensitivity is that it is a
linear measure. That is, it assumes that for a certain percentage change in interest rates that an
equal percentage change in price will occur. However, as interest rates change, the price of a bond
is not likely to change linearly, but instead would change over some curved, or "convex", function
of interest rates.

For any given bond, a graph of the relationship between price and yield is convex. This means that
the graph forms a curve rather than a straight line. The more convex the relationship the more
inaccurate duration is as a measure of the interest rate sensitivity.

The convexity of a bond is a measure of the curvature of its price/yield relationship. The degree
to which the graph is curved shows how much a bond's yield changes in response to a change in
price.

Used in conjunction with duration, convexity provides a more accurate approximation of the
percentage price change resulting from a specified change in a bond's yield than using duration
alone. In addition to improving the estimate of a bond's price changes to changes in interest
rates, convexity can also be used to compare bonds with the same duration. For example, two bonds
may have the same duration but different convexity values. They may experience different price
changes when there are extraordinary changes in interest rates. For example, if bond A has a higher
convexity than bond B, its price would fall less during rising interest rates and appreciate more
during falling interest rates as compared to bond B. (Zero-coupon bonds, which pay their entire
cash flows at maturity, have the highest convexity. This is because, in general, the more dispersed
the cash flows are, the greater the convexity will be.)

This table shows how the price a of a 10-year bond with a 10% coupon changes at different yields.
The column labeled "Delta ($)" shows the absolute change in price. As you can see, the bond's price
(y-axis) rises at an increasing rate as yields (x-axis) fall, but declines at a decreasing rate as
yields rise. This characteristic causes the line to be convex instead of straight.
Yield	  Change   Price     Delta
 (%)     (%)       ($)      ($)
  1      (9)    1,854.43   854.43
  2      (8)    1,721.82   721.82
  3      (7)    1,600.90   600.90
  4      (6)    1,490.54   490.54
  5      (5)    1,389.73   389.73
  6      (4)    1,297.55   297.55
  7      (3)    1,213.19   213.19
  8      (2)    1,135.90   135.90
  9      (1)    1,065.04    65.04
 10       0     1,000.00     0.00
 11       1       940.25   (59.75)
 12       2       885.30  (114.70)
 13       3       834.72  (165.28)
 14       4       788.12  (211.88)
 15       5       745.14  (254.86)
 16       6       705.46  (294.54)
 17       7       668.78  (331.22)
 18       8       634.86  (365.14)
 19       9       603.44  (396.56)
***/
func (b *Bonds) Duration(cashFlow []float64, cp int, currentRate, bondPrice float64) float64 {
  currentRate /= float64(cp)
  D := b.duration(cashFlow, currentRate, bondPrice)
  return (D / float64(cp))
}

func (b *Bonds) duration(cashFlow []float64, currentRate, bondPrice float64) float64 {
  var D float64 = zero
  currentRate = one + (currentRate / hundred)
  var sz int = len(cashFlow)
  for idx, t := 0, 1.0; idx < sz; idx, t = idx + 1, t + 1.0 {
    D += (t * cashFlow[idx]) / math.Pow(currentRate, t)
  }
  return (D / bondPrice)
}

/*** Silence the message U1000 for unused functions. Uncomment when usage is required.
func (b *Bonds) duration1(cashFlow []float64, currentRate float64) (float64) {
  var D float64 = zero
  var B float64 = zero
  currentRate = one + (currentRate / hundred)
  var sz int = len(cashFlow)
  var tmp float64 = zero
  for idx, t := 0, 1.0; idx < sz; idx, t = idx + 1, t + 1.0 {
    tmp = math.Pow(currentRate, t)
    D += (t * cashFlow[idx]) / tmp
    B += cashFlow[idx] / tmp
  }
  return (D / B)
}
***/

func (b *Bonds) DurationContinuous(cashFlow []float64, currentRate, bondPrice float64) float64 {
  // var B float64 = zero  //Bond Price
  var D float64 = zero
  currentRate /= hundred
  var sz int = len(cashFlow)
  for idx, t := 0, 1.0; idx < sz; idx, t = idx + 1, t + 1.0 {
  //  B += cashFlow[idx] * math.Exp(-currentRate * t)
    D += math.Exp(-currentRate * t) * t * cashFlow[idx]
  }
  // return(D / B)
  return(D / bondPrice)
}

/***
If the bond is priced correctly, the yield to maturity must equal the current interest rate. If
current interest rate EQUALS yield to maturity the calculations from Duration and MacaulayDuration
will produce the same number.
Notes:
(1) The longer the duration is the more sensitive the bond will be to changes in interest rates.
(2) For a standard bond the Macaulay duration will be between 0 and the maturity of the bond. It is
    equal to the maturity if and only if the bond is a zero-coupon bond.
***/
func (b *Bonds) MacaulayDuration(cashFlow []float64, cp int, bondPrice float64) float64 {
  var ytm = b.YieldToMaturity(cashFlow, bondPrice, cp) / float64(cp)
  D := b.duration(cashFlow, ytm, bondPrice)
  return (D / float64(cp))
}

func (b *Bonds) MacaulayDurationContinuous(cashFlow []float64, bondPrice float64) float64 {
  var ytm = b.YieldToMaturityContinuous(cashFlow, bondPrice)
  return(b.DurationContinuous(cashFlow, ytm, bondPrice))
}

/***
Modified duration is a formula that expresses the measurable change in the value of a security in
response to a change in interest rates. Modified duration follows the concept that interest rates
and bond prices move in opposite directions. This formula is used to determine the effect that a
100-basis point (1%) change in interest rates will have on the price of a bond.

How to interpret the result below? The modified duration illustrates the effect of a 100-basis
point (1%) change in interest rates on the price of a bond. Therefore,
(1) If interest rates increase by 1%, the price of the 3-year bond will decrease by 2.513%.
(2) If interest rates decrease by 1%, the price of the 3-year bond will increase by 2.513%.

The modified duration provides a good measurement of a bond's sensitivity to changes in interest
rates.

FV = $100.00
coupon rate = 10%; compounding period annually
t = 3-year
current interest = 9%; compounding period annually
Modified Duration: 2.512801%
***/
func (b *Bonds) ModifiedDuration(cashFlow []float64, cp int, bondPrice float64) float64 {
  var ytm = b.YieldToMaturity(cashFlow, bondPrice, cp) / float64(cp)
  D := b.duration(cashFlow, ytm, bondPrice) / (one + (ytm / hundred))
  return (D / float64(cp))
}

/***
Duration is a linear measure or 1st derivative of how the price of a bond changes in response to
interest rate changes. As interest rates change, the price is not likely to change linearly, but
instead it would change over some curved function of interest rates. The more curved the price
function of the bond is, the more inaccurate duration is as a measure of the interest rate
sensitivity.

(Duration can be a good measure of how bond prices may be affected due to small and sudden
fluctuations in interest rates. However, the relationship between bond prices and yields is
typically more sloped, or convex. Therefore, convexity is a better measure for assessing the
impact on bond prices when there are large fluctuations in interest rates.)

Convexity is a measure of the curvature or 2nd derivative of how the price of a bond varies with
interest rate; i.e., how the duration of a bond changes as the interest rate changes. Specifically,
one assumes that the interest rate is constant across the life of the bond and that changes in
interest rates occur evenly. Using these assumptions, duration can be formulated as the first
derivative of the price function of the bond with respect to the interest rate in question. Then
the convexity would be the second derivative of the price function with respect to the interest
rate.

In actual markets, the assumption of constant interest rates and even changes is not correct, and
more complex models are needed to actually price bonds. However, these simplifying assumptions
allow one to quickly and easily calculate factors which describe the sensitivity of the bond prices
to interest rate changes.

Convexity does not assume the relationship between bond value and interest rates to be linear. For
large fluctuations in interest rates, it is a better measure than duration.
Notes:
(1) It's important to know how bond prices and market interest rates relate to one another. As
    interest rates fall, bond prices rise. Conversely, rising market interest rates lead to falling
    bond prices. This opposite reaction is because as rates rise, the bond may fall behind in the
    earnings they may offer a potential investor in comparison to other securities.
(2) If a bond's duration increases as yields increase, the bond is said to have negative convexity.
    In other words, the bond price will decline by a greater rate with a rise in yields than if
    yields had fallen. Therefore, if a bond has negative convexity, its duration would increase --
    the price would fall. As interest rates rise, the opposite is true.
(3) If a bond's duration rises and yields fall, the bond is said to have positive convexity. In
    other words, as yields fall, bond prices rise by a greater rate -- or duration -- than if
    yields rose. Positive convexity leads to greater increases in bond prices. If a bond has
    positive convexity, it would typically experience larger price increases as yields fall,
    compared to price decreases when yields increase.
(4) Zero-coupon bonds have the highest degree of convexity because they do not offer any coupon
    payments.
***/
func (b *Bonds) Convexity(cashFlow []float64, currentRate float64, cp int) float64 {
  var B float64 = b.CurrentPrice(cashFlow, currentRate, cp)
  var Cx float64 = zero
  currentRate = one + ((currentRate / float64(cp)) / hundred)
  var sz int = len(cashFlow)
  for idx, t := 0, 1.0; idx < sz; idx, t = idx + 1, t + 1.0 {
    Cx += ((one + t) * t * cashFlow[idx]) / math.Pow(currentRate, t)
  }
  Cx /= float64(cp)
  Cx = (Cx / B) / math.Pow(currentRate, 2)
  return (Cx / float64(cp))
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
