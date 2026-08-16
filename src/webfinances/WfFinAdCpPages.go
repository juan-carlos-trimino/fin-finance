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

type adCpFields struct {
  MenuPage string `json:"menuPage"`
  CurrentPage string `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1PV string `json:"fd1PV"`
  Fd1FV string `json:"fd1FV"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2Interest string `json:"fd2Interest"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2Payment string `json:"fd2Payment"`
  Fd2PV string `json:"fd2PV"`
  Fd2Result string `json:"fd2Result"`
  //
  Fd3Interest string `json:"fd3Interest"`
  Fd3Compound string `json:"fd3Compound"`
  Fd3Payment string `json:"fd3Payment"`
  Fd3FV string `json:"fd3FV"`
  Fd3Result string `json:"fd3Result"`
}

func newAdCpFields(dir1, dir2, correlationId string) *adCpFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  //Default values returned if file is missing, empty, or JSON is corrupt.
  m := adCpFields {
    MenuPage: "",
    CurrentPage: "rhs-ui2",
    CurrentButton: "lhs-button2",
    //
    Fd1Interest: "1.00",
    Fd1Compound: "annually",
    Fd1PV: "1.00",
    Fd1FV: "1.00",
    Fd1Result: "",
    //
    Fd2Interest: "1.00",
    Fd2Compound: "annually",
    Fd2Payment: "1.00",
    Fd2PV: "1.00",
    Fd2Result: "",
    //
    Fd3Interest: "1.00",
    Fd3Compound: "annually",
    Fd3Payment: "1.00",
    Fd3FV: "1.00",
    Fd3Result: "",
  }
  obj, err := readFields(dir + "adcp.txt")
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
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "adcp.txt"), correlationId)
  }
  return &m
}

func getAdCpFields(userName string) *adCpFields {
  return currentFields[userName].adCp
}

type WfAdCpPages struct {}

func (a WfAdCpPages) AdCpPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.AdCpPages.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res, correlationId)
    return
  }
  //
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    fields := getAdCpFields(userName)
    /***
    The functions in Request that allow to extract data from the URL and/or the body revolve around the Form, PostForm, and
    MultipartForm fields; the data are in the form of key-value pairs.

    If the form and the URL have the same key name, both of them will be placed in a slice, with the form value always prioritized
    before the URL value.

    Since we want the form key-value pairs, we can ignore the URL key-value pairs. The PostForm field provides key-value pairs only
    for the form and not the URL. The PostForm field supports only application/x-www-form-urlencoded.

    The FormValue method lets you access the key-value pairs just like the Form field, except that it's for a specific key and there
    is no need to call the ParseForm method beforehand -- the FormValue method does it. The PostFormValue method does the same thing,
    except that it's for the PostForm field instead of the Form field.
    ***/
    if ui := req.FormValue("compute"); ui != "" {  //Values from form and URL.
      fields.CurrentPage = ui
    }
    /***
    // if strings.EqualFold(currentRHS, "rhs-ui1") {
    //   p.CurrentButton = "lhs-button1"
    //   if req.Method == http.MethodPost {
    //     p.fd1Interest = req.PostFormValue("fd1-interest")
    //     p.fd1Compound = req.PostFormValue("fd1-cp")
    //     p.fd1PV = req.PostFormValue("fd1-pv")
    //     p.fd1FV = req.PostFormValue("fd1-fv")
    //     var i float64
    //     var pv float64
    //     var fv float64
    //     var err error
    //     if i, err = strconv.ParseFloat(p.fd1Interest, 64); err != nil {
    //       p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1Interest, err)
    //     } else if pv, err = strconv.ParseFloat(p.fd1PV, 64); err != nil {
    //       p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1PV, err)
    //     } else if fv, err = strconv.ParseFloat(p.fd1FV, 64); err != nil {
    //       p.fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1FV, err)
    //     } else {
    //       var oa finances.Annuities
    //       p.fd1Result = fmt.Sprintf("Compounding Period: %.5f %s", oa.O_Periods_PV_FV(pv, fv,
    //                                 i / 100.0, oa.GetCompoundingPeriod(p.fd1Compound[0], true)),
    //                                 oa.TimePeriods(p.fd1Compound))
    //     }
    //     logEntry.Print(INFO, correlationId, []string {
    //       fmt.Sprintf("i = %s, cp = %s, pv = %s, fv = %s, %s",
    //                   p.fd1Interest, p.fd1Compound, p.fd1PV, p.fd1FV, p.fd1Result),
    //     })
    //   }
    //   t := template.Must(template.ParseFiles("webfinances/templates/ordinaryannuity/cp/cp.html",
    //                                          "webfinances/templates/header.html",
    //                                          "webfinances/templates/ordinaryannuity/cp/i-PV-FV.html",
    //                                          "webfinances/templates/footer.html"))
    //   t.ExecuteTemplate(res, "oacompoundingperiods", struct {
    //     Header string
    //     Datetime string
    //     CurrentButton string
    //     Fd1Interest string
    //     Fd1Compound string
    //     Fd1PV string
    //     Fd1FV string
    //     Fd1Result string
    //   } { "Ordinary Annuity / Compounding Periods", m.DTF(), p.CurrentButton,
    //       p.fd1Interest, p.fd1Compound, p.fd1PV, p.fd1FV, p.fd1Result,
    //     })
    } else*/
    if strings.EqualFold(fields.CurrentPage, "rhs-ui2") {
      fields.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        fields.Fd2Interest = req.FormValue("fd2-interest")
        fields.Fd2Compound = req.PostFormValue("fd2-cp")
        fields.Fd2Payment = req.PostFormValue("fd2-payment")
        fields.Fd2PV = req.PostFormValue("fd2-pv")
        var i float64
        var pmt float64
        var pv float64
        var err error
        if i, err = strconv.ParseFloat(fields.Fd2Interest, 64); err != nil {
          fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Interest, err)
        } else if pmt, err = strconv.ParseFloat(fields.Fd2Payment, 64); err != nil {
          fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Payment, err)
        } else if pv, err = strconv.ParseFloat(fields.Fd2PV, 64); err != nil {
          fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2PV, err)
        } else {
          var oa finances.Annuities
          fields.Fd2Result = fmt.Sprintf("Compounding Period: %.5f %s", oa.D_Periods_PMT_PV(pmt, pv, i / 100.0,
            oa.GetCompoundingPeriod(fields.Fd2Compound[0], true)), oa.TimePeriods(fields.Fd2Compound))
        }
        logger.LogInfo(fmt.Sprintf("i = %s, cp = %s, pmt = %s, pv = %s, %s", fields.Fd2Interest, fields.Fd2Compound,
          fields.Fd2Payment, fields.Fd2PV, fields.Fd2Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/annuitydue/cp/cp.html",
        "webfinances/templates/finances/annuitydue/cp/i-PMT-PV.html",
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
          Fd2Interest string
          Fd2Compound string
          Fd2Payment string
          Fd2PV string
          Fd2Result string
        } { "standard", "Annuity Due / Compounding Periods", logger.DatetimeFormat(), financesMenuPage, fields.CurrentButton,
            newSession.CsrfToken, fields.Fd2Interest, fields.Fd2Compound, fields.Fd2Payment, fields.Fd2PV, fields.Fd2Result },
      })
    } else if strings.EqualFold(fields.CurrentPage, "rhs-ui3") {
      fields.CurrentButton = "lhs-button3"
      if req.Method == http.MethodPost {
        fields.Fd3Interest = req.FormValue("fd3-interest")
        fields.Fd3Compound = req.PostFormValue("fd3-cp")
        fields.Fd3Payment = req.PostFormValue("fd3-payment")
        fields.Fd3FV = req.PostFormValue("fd3-fv")
        var i float64
        var pmt float64
        var fv float64
        var err error
        if i, err = strconv.ParseFloat(fields.Fd3Interest, 64); err != nil {
          fields.Fd3Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd3Interest, err)
        } else if pmt, err = strconv.ParseFloat(fields.Fd3Payment, 64); err != nil {
          fields.Fd3Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd3Payment, err)
        } else if fv, err = strconv.ParseFloat(fields.Fd3FV, 64); err != nil {
          fields.Fd3Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd3FV, err)
        } else {
          var oa finances.Annuities
          fields.Fd3Result = fmt.Sprintf("Compounding Period: %.5f %s", oa.D_Periods_PMT_FV(pmt, fv, i / 100.0,
            oa.GetCompoundingPeriod(fields.Fd3Compound[0], true)), oa.TimePeriods(fields.Fd3Compound))
        }
        logger.LogInfo(fmt.Sprintf("i = %s, cp = %s, pmt = %s, fv = %s, %s", fields.Fd3Interest, fields.Fd3Compound,
          fields.Fd3Payment, fields.Fd3FV, fields.Fd3Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/annuitydue/cp/cp.html",
        "webfinances/templates/finances/annuitydue/cp/i-PMT-FV.html",
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
          Fd3Interest string
          Fd3Compound string
          Fd3Payment string
          Fd3FV string
          Fd3Result string
        } { "standard", "Annuity Due / Compounding Periods", logger.DatetimeFormat(), financesMenuPage, fields.CurrentButton,
            newSession.CsrfToken, fields.Fd3Interest, fields.Fd3Compound, fields.Fd3Payment, fields.Fd3FV, fields.Fd3Result },
      })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", fields.CurrentPage)
      logger.LogError(errString, correlationId)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      logger.LogWarning("*** Request timeout ***", correlationId)
      if strings.EqualFold(fields.CurrentPage, "rhs-ui2") {
        fields.Fd2Result = ""
      } else if strings.EqualFold(fields.CurrentPage, "rhs-ui3") {
        fields.Fd3Result = ""
      }
    }
    //
    if data, err := json.Marshal(fields); err != nil {
      logger.LogError(fmt.Sprintf("%+v", err), correlationId)
    } else {
      filePath := fmt.Sprintf("%s/%s/adcp.txt", mainDir, userName)
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
