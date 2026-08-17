package webfinances

import (
  "context"
  "encoding/json"
  "finance/finances"
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/gposu"
  "github.com/juan-carlos-trimino/gpsessions"
  "net/http"
  "os"
  "strconv"
  "strings"
  "time"
)

type siAccurateFields struct{
  MenuPage string `json:"menuPage"`
  CurrentPage string `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1Time string `json:"fd1Time"`
  Fd1Leap bool `json:"fd1Leap"`
  Fd1TimePeriod string `json:"fd1TimePeriod"`
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1PV string `json:"fd1PV"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2Time string `json:"fd2Time"`
  Fd2TimePeriod string `json:"fd2TimePeriod"`
  Fd2Amount string `json:"fd2Amount"`
  Fd2PV string `json:"fd2PV"`
  Fd2Result string `json:"fd2Result"`
  //
  Fd3Time string `json:"fd3Time"`
  Fd3TimePeriod string `json:"fd3TimePeriod"`
  Fd3Interest string `json:"fd3Interest"`
  Fd3Compound string `json:"fd3Compound"`
  Fd3Amount string `json:"fd3Amount"`
  Fd3Result string `json:"fd3Result"`
  //
  Fd4Interest string `json:"fd4Interest"`
  Fd4Compound string `json:"fd4Compound"`
  Fd4Amount string `json:"fd4Amount"`
  Fd4PV string `json:"fd4PV"`
  Fd4Result string `json:"fd4Result"`
}

func newSiAccurateFields(dir1, dir2, correlationId string) *siAccurateFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  //Default values returned if file is missing, empty, or JSON is corrupt.
  m := siAccurateFields{
    MenuPage: "",
    CurrentPage: "rhs-ui1",
    CurrentButton: "lhs-button1",
    //
    Fd1Time: "1",
    Fd1Leap: false,
    Fd1TimePeriod: "year",
    Fd1Interest: "1.00",
    Fd1Compound: "annually",
    Fd1PV: "1.00",
    Fd1Result: "",
    //
    Fd2Time: "1",
    Fd2TimePeriod: "year",
    Fd2Amount: "1.00",
    Fd2PV: "1.00",
    Fd2Result: "",
    //
    Fd3Time: "1",
    Fd3TimePeriod: "year",
    Fd3Interest: "1.0",
    Fd3Compound: "annually",
    Fd3Amount: "1.00",
    Fd3Result: "",
    //
    Fd4Interest: "1.00",
    Fd4Compound: "annually",
    Fd4Amount: "1.00",
    Fd4PV: "1.00",
    Fd4Result: "",
  }
  obj, err := readFields(dir + "siaccurate.txt")
  if obj != nil {
    /***
    When a file is empty, the readFields function successfully returns a valid slice, but it contains zero bytes. Checking the
    length ensures parsing only files that actually contain data.
    ***/
    if len(obj) != 0 {  //Check if the file contains no data (empty)
      err = json.Unmarshal(obj, &m)
      if err != nil {
        //Write error, but continue with default values.
        logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
      }
    }
  } else if err != nil {
    logger.LogError(fmt.Sprintf("%+v", err), correlationId)
  } else {
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "siaccurate.txt"), correlationId)
  }
  return &m
}

func getSiAccurateFields(userName string) *siAccurateFields {
  return currentFields[userName].siAccurate
}

type WfSiAccuratePages struct{}

func (s WfSiAccuratePages) SimpleInterestAccuratePages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.SimpleInterestAccuratePages.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res, correlationId)
    return
  }
  //
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    fields := getSiAccurateFields(userName)
    /***
    The functions in Request that allow to extract data from the URL and/or the body revolve around the Form, PostForm, and MultipartForm
    fields; the data are in the form of key-value pairs.

    If the form and the URL have the same key name, both of them will be placed in a slice, with the form value always prioritized before
    the URL value.

    Since we want the form key-value pairs, we can ignore the URL key-value pairs. The PostForm field provides key-value pairs only for
    the form and not the URL. The PostForm field supports only application/x-www-form-urlencoded.

    The FormValue method lets you access the key-value pairs just like the Form field, except that it's for a specific key and there is no
    need to call the ParseForm method beforehand -- the FormValue method does it. The PostFormValue method does the same thing, except that
    it's for the PostForm field instead of the Form field.
    ***/
    if ui := req.FormValue("compute"); ui != "" {  //Values from form and URL.
      fields.CurrentPage = ui
    }
    //
    if strings.EqualFold(fields.CurrentPage, "rhs-ui1") {
      fields.CurrentButton = "lhs-button1"
      if req.Method == http.MethodPost {
        fields.Fd1Time = req.PostFormValue("fd1-time")
        fields.Fd1TimePeriod = req.PostFormValue("fd1-tp")
        /***
        In HTML whenever a checkbox is checked and is POSTed, it has a value of "on" by default. If the checkbox is unchecked
        it has no value ("").
        ***/
        var daysInYear int = finances.Daily365
        fields.Fd1Leap = false
        if req.PostFormValue("fd1-leap") != "" {
          fields.Fd1Leap = true
          daysInYear = finances.Daily366
        }
        fields.Fd1TimePeriod = req.PostFormValue("fd1-tp")
        fields.Fd1Interest = req.PostFormValue("fd1-interest")
        fields.Fd1Compound = req.PostFormValue("fd1-compound")
        fields.Fd1PV = req.PostFormValue("fd1-pv")
        var n float64
        var i float64
        var pv float64
        var err error
        if n, err = strconv.ParseFloat(fields.Fd1Time, 64); err != nil {
          fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1Time, err)
        } else if i, err = strconv.ParseFloat(fields.Fd1Interest, 64); err != nil {
          fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1Interest, err)
        } else if pv, err = strconv.ParseFloat(fields.Fd1PV, 64); err != nil {
          fields.Fd1Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd1PV, err)
        } else {
          var si finances.SimpleInterest
          var periods finances.Periods
          fields.Fd1Result = fmt.Sprintf("Amount of Interest: $%.5f",
            si.AccurateInterest(pv, i / 100.0, periods.GetCompoundingPeriod(fields.Fd1Compound[0], true), n, daysInYear))
        }
        logger.LogInfo(fmt.Sprintf("n = %s, leap = %t, tp = %s, i = %s, cp = %s, pv = %s, %s", fields.Fd1Time, fields.Fd1Leap,
          fields.Fd1TimePeriod, fields.Fd1Interest, fields.Fd1Compound, fields.Fd1PV, fields.Fd1Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/simpleinterestaccurate/accurate.html",
        "webfinances/templates/finances/simpleinterestaccurate/amountofinterest.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/navbar.html",
        "webfinances/templates/footer.html",
      }
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
        Data: struct{
          LayoutType string
          Header string
          Datetime string
          MenuPage string
          CurrentButton string
          CsrfToken string
          Fd1Time string
          Fd1Leap bool
          Fd1TimePeriod string
          Fd1Interest string
          Fd1Compound string
          Fd1PV string
          Fd1Result string
        } { "standard", "Simple Interest / Accurate (Exact) Interest", logger.DatetimeFormat(), financesMenuPage, fields.CurrentButton,
            newSession.CsrfToken, fields.Fd1Time, fields.Fd1Leap, fields.Fd1TimePeriod, fields.Fd1Interest, fields.Fd1Compound,
            fields.Fd1PV, fields.Fd1Result },
      })
    } else if strings.EqualFold(fields.CurrentPage, "rhs-ui2") {
      fields.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        fields.Fd2Time = req.PostFormValue("fd2-time")
        fields.Fd2TimePeriod = req.PostFormValue("fd2-tp")
        fields.Fd2Amount = req.PostFormValue("fd2-amount")
        fields.Fd2PV = req.PostFormValue("fd2-pv")
        var n float64
        var a float64
        var pv float64
        var err error
        if n, err = strconv.ParseFloat(fields.Fd2Time, 64); err != nil {
          fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Time, err)
        } else if a, err = strconv.ParseFloat(fields.Fd2Amount, 64); err != nil {
          fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Amount, err)
        } else if pv, err = strconv.ParseFloat(fields.Fd2PV, 64); err != nil {
          fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2PV, err)
        } else {
          var si finances.SimpleInterest
          var periods finances.Periods
          fields.Fd2Result = fmt.Sprintf("Interest Rate: %.5f%%",
            si.AccurateRate(pv, a, n, periods.GetTimePeriod(fields.Fd2TimePeriod[0], true)) * 100.0)
        }
        logger.LogInfo(fmt.Sprintf("n = %s, tp = %s, a = %s, pv = %s, %s", fields.Fd2Time, fields.Fd2TimePeriod, fields.Fd2Amount,
          fields.Fd2PV, fields.Fd2Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/simpleinterestaccurate/accurate.html",
        "webfinances/templates/finances/simpleinterestaccurate/interestrate.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/navbar.html",
        "webfinances/templates/footer.html",
      }
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
        Data: struct{
          Layout string
          Header string
          Datetime string
          MenuPage string
          CurrentButton string
          CsrfToken string
          Fd2Time string
          Fd2TimePeriod string
          Fd2Amount string
          Fd2PV string
          Fd2Result string
        } { "standard", "Simple Interest / Accurate (Exact) Interest", logger.DatetimeFormat(), financesMenuPage, fields.CurrentButton,
            newSession.CsrfToken, fields.Fd2Time, fields.Fd2TimePeriod, fields.Fd2Amount, fields.Fd2PV, fields.Fd2Result },
      })
    } else if strings.EqualFold(fields.CurrentPage, "rhs-ui3") {
      fields.CurrentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        fields.Fd3Time = req.PostFormValue("fd3-time")
        fields.Fd3TimePeriod = req.PostFormValue("fd3-tp")
        fields.Fd3Interest = req.PostFormValue("fd3-interest")
        fields.Fd3Compound = req.PostFormValue("fd3-compound")
        fields.Fd3Amount = req.PostFormValue("fd3-amount")
        var n float64
        var i float64
        var a float64
        var err error
        if n, err = strconv.ParseFloat(fields.Fd3Time, 64); err != nil {
          fields.Fd3Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd3Time, err)
        } else if i, err = strconv.ParseFloat(fields.Fd3Interest, 64); err != nil {
          fields.Fd3Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd3Interest, err)
        } else if a, err = strconv.ParseFloat(fields.Fd3Amount, 64); err != nil {
          fields.Fd3Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd3Amount, err)
        } else {
          var si finances.SimpleInterest
          var periods finances.Periods
          fields.Fd3Result = fmt.Sprintf("Principal: $%.5f", si.AccuratePrincipal(a, i / 100.0,
            periods.GetCompoundingPeriod(fields.Fd3Compound[0], true), n, periods.GetTimePeriod(fields.Fd3TimePeriod[0], true)))
        }
        logger.LogInfo(fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, a = %s, %s", fields.Fd3Time, fields.Fd3TimePeriod,
          fields.Fd3Interest, fields.Fd3Compound, fields.Fd3Amount, fields.Fd3Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/simpleinterestaccurate/accurate.html",
        "webfinances/templates/finances/simpleinterestaccurate/principal.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/navbar.html",
        "webfinances/templates/footer.html",
      }
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
        Data: struct{
          LayoutType string
          Header string
          Datetime string
          MenuPage string
          CurrentButton string
          CsrfToken string
          Fd3Time string
          Fd3TimePeriod string
          Fd3Interest string
          Fd3Compound string
          Fd3Amount string
          Fd3Result string
        } { "standard", "Simple Interest / Accurate (Exact) Interest", logger.DatetimeFormat(), financesMenuPage, fields.CurrentButton,
            newSession.CsrfToken, fields.Fd3Time, fields.Fd3TimePeriod, fields.Fd3Interest, fields.Fd3Compound, fields.Fd3Amount,
            fields.Fd3Result },
      })
    } else if strings.EqualFold(fields.CurrentPage, "rhs-ui4") {
      fields.CurrentButton = "lhs-button4"
      if req.Method == http.MethodPost {
        fields.Fd4Interest = req.FormValue("fd4-interest")
        fields.Fd4Compound = req.FormValue("fd4-compound")
        fields.Fd4Amount = req.FormValue("fd4-amount")
        fields.Fd4PV = req.FormValue("fd4-pv")
        var i float64
        var a float64
        var pv float64
        var err error
        if i, err = strconv.ParseFloat(fields.Fd4Interest, 64); err != nil {
          fields.Fd4Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd4Interest, err)
        } else if a, err = strconv.ParseFloat(fields.Fd4Amount, 64); err != nil {
          fields.Fd4Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd4Amount, err)
        } else if pv, err = strconv.ParseFloat(fields.Fd4PV, 64); err != nil {
          fields.Fd4Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd4PV, err)
        } else {
          var si finances.SimpleInterest
          var periods finances.Periods
          fields.Fd4Result = fmt.Sprintf("Time: %.5f %s",
            si.AccurateTime(pv, a, i / 100.0, periods.GetCompoundingPeriod(fields.Fd4Compound[0], true)),
            periods.TimePeriods(fields.Fd4Compound))
        }
        logger.LogInfo(fmt.Sprintf("i = %s, cp = %s, a = %s, pv = %s, %s", fields.Fd4Interest, fields.Fd4Compound, fields.Fd4Amount,
          fields.Fd4PV, fields.Fd4Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/simpleinterestaccurate/accurate.html",
        "webfinances/templates/finances/simpleinterestaccurate/time.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/navbar.html",
        "webfinances/templates/footer.html",
      }
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
        Data: struct{
          LayoutType string
          Header string
          Datetime string
          MenuPage string
          CurrentButton string
          CsrfToken string
          Fd4Interest string
          Fd4Compound string
          Fd4Amount string
          Fd4PV string
          Fd4Result string
        } { "standard", "Simple Interest / Accurate (Exact) Interest", logger.DatetimeFormat(), financesMenuPage, fields.CurrentButton,
            newSession.CsrfToken, fields.Fd4Interest, fields.Fd4Compound, fields.Fd4Amount, fields.Fd4PV, fields.Fd4Result },
      })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", fields.CurrentPage)
      logger.LogInfo(errString, correlationId)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      logger.LogWarning("*** Request timeout ***", correlationId)
      if strings.EqualFold(fields.CurrentPage, "rhs-ui1") {
        fields.Fd1Result = ""
      } else if strings.EqualFold(fields.CurrentPage, "rhs-ui2") {
        fields.Fd2Result = ""
      } else if strings.EqualFold(fields.CurrentPage, "rhs-ui3") {
        fields.Fd3Result = ""
      } else if strings.EqualFold(fields.CurrentPage, "rhs-ui4") {
        fields.Fd4Result = ""
      }
    }
    //
    if data, err := json.Marshal(fields); err != nil {
      logger.LogError(fmt.Sprintf("%+v", err), correlationId)
    } else {
      filePath := fmt.Sprintf("%s/%s/siaccurate.txt", mainDir, userName)
      if _, err := osu.WriteAllExclusiveLock1(filePath, data, os.O_CREATE | os.O_RDWR | os.O_TRUNC, 0o600); err != nil {
        logger.LogError(fmt.Sprintf("%+v", err), correlationId)
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    logger.LogError(errString, correlationId)
    panic(errString)
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms\n", time.Since(startTime).Microseconds()), correlationId)
}
