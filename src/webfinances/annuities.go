//
package webfinances

import (
  "finance/finances"
  "fmt"
  "net/http"
  "strconv"
)

type Annuities struct{}

func (a Annuities) AverageRateOfReturn(res http.ResponseWriter, req *http.Request) () {
  params := req.URL.Query()
  if len(params) != 1 {
    fmt.Fprintf(res, "Parameters required = 1; parameters provided = %d\n", len(params))
    return
  }
  const length int = 16
  var returns []float64 = nil
  //Iterate over all the query parameters.
  for _, v := range params {
    returns = make([]float64, 0, length)
    for idx := 0; idx < len(v); idx++ {
      f, err := strconv.ParseFloat(v[idx], 64)
      if err != nil {
        fmt.Fprintf(res, "'%s' is not a floating number.\n%s", v[idx], v)
        return
      }
      returns = append(returns, f)
    }
  }
  var fa finances.Annuities
  var gmr = fa.AverageRateOfReturn(returns) * 100.0
  fmt.Fprintf(res, "AverageRateOfReturn = %.3f\n", gmr)
}



