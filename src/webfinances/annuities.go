package webfinances

import (
  "finance/finances"
  "finance/mathutil"
  "fmt"
  "net/http"
  "strings"
)

type Annuities struct{}

/***
To execute this function from a browser:
fin/annuities/AverageRateOfReturn?rets=5.0&rets=-3.0&rets=12.0&rets=10
Average Rate of Return = 5.838
***/
func (a Annuities) AverageRateOfReturn(res http.ResponseWriter, req *http.Request) {
  const paramsRequired int = 1
  params := req.URL.Query()
  if len(params) != paramsRequired {
    fmt.Fprintf(res, "Parameters required = 1; parameters provided = %d", len(params))
    return
  }
  var mu mathutil.MathUtil
  //A nil slice is also an empty slice; no allocation of memory.
  var returns []float64
  var err = error(nil)
  //Iterate over all the query parameters.
  for k, v := range params { //map[string][]string
    switch strings.ToLower(k) {
    case "rets":
      returns, err = mu.ConvertToFloats64(v)
      if err != nil {
        fmt.Fprintf(res, "%s", err)
        return
      }
    default:
      fmt.Fprintf(res, "'%s' is an invalid parameter name.", k)
      return
    }
  }
  var fa finances.Annuities
  var gmr = fa.AverageRateOfReturn(returns) * 100.0
  fmt.Fprintf(res, "Average Rate of Return = %.3f\n", gmr)
}

/***
To execute this function from a browser:
fin/annuities/GrowthDecayOfFunds?factor=2.0&rate=15.0&cp=A
Average Rate of Return = 5.838
***/
func (a Annuities) GrowthDecayOfFunds(res http.ResponseWriter, req *http.Request) {
  const paramsRequired int = 3
  params := req.URL.Query()
  if len(params) != paramsRequired {
    fmt.Fprintf(res, "Parameters required = 3; parameters provided = %d", len(params))
    return
  }
  var cp int = finances.Invalid
  var factor, rate float64
  var err error
  var mu mathutil.MathUtil
  var p finances.Periods
  for k, v := range params { //map[string][]string
    switch strings.ToLower(k) {
    case "factor":
      if len(v) == 1 {
        factor, err = mu.ConvertToFloat64(v[0])
        if err != nil {
          fmt.Fprintf(res, "%s", err)
          return
        }
      } else {
        fmt.Fprintf(res, "'factor' is not an array: '%s'.", v)
        return
      }
    case "rate":
      if len(v) == 1 {
        rate, err = mu.ConvertToFloat64(v[0])
        if err != nil {
          fmt.Fprintf(res, "%s", err)
          return
        }
        rate /= 100.0
      } else {
        fmt.Fprintf(res, "'rate' is not an array: '%s'.", v)
        return
      }
    case "cp":
      if len(v) == 1 {
        if len(v[0]) == 1 {
          cp = p.GetCompoundingPeriod(byte(v[0][0]), true)
          if cp == finances.Invalid {
            fmt.Fprintf(res, "'cp' has an invalid value: '%s'.", v[0])
            return
          }
        } else {
          fmt.Fprintf(res, "'cp' has an invalid value: '%s'.", v[0])
          return
        }
      } else {
        fmt.Fprintf(res, "'cp' is not an array: '%s'.", v)
        return
      }
    default:
      fmt.Fprintf(res, "'%s' is an invalid parameter name.", k)
      return
    }
  }
  var fa finances.Annuities
  gd := fa.GrowthDecayOfFunds(factor, rate, cp)
  fmt.Fprintf(res, "Growth or Decay = %f", gd)
}
