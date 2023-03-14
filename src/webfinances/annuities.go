package webfinances

import (
	"finance/finances"
	"finance/mathutil"
	"fmt"
	"net/http"
	// "strconv"
)

type Annuities struct{}

/***
To execute this function from a browser:
annuities/AverageRateOfReturn?ret=5.0&ret=-3.0&ret=12.0&ret=10
Average Rate of Return = 5.838
***/
func (a Annuities) AverageRateOfReturn(res http.ResponseWriter, req *http.Request) () {
  params := req.URL.Query()
  if len(params) != 1 {
    fmt.Fprintf(res, "Parameters required = 1; parameters provided = %d", len(params))
    return
  }
  var mu mathutil.MathUtil
  //A nil slice is also an empty slice; no allocation of memory.
  var returns []float64
  var err = error(nil)
  //Iterate over all the query parameters.
  for _, v := range params {
    returns, err = mu.ConvertToFloat64(v)
    if err != nil {
      fmt.Fprintf(res, "%s", err)
      return
    }
  }
  var fa finances.Annuities
  var gmr = fa.AverageRateOfReturn(returns) * 100.0
  fmt.Fprintf(res, "Average Rate of Return = %.3f\n", gmr)
}
