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

type adPvFields struct {
  MenuPage string `json:"menuPage"`
  CurrentPage string `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1N string `json:"fd1N"`
  Fd1TimePeriod string `json:"fd1TimePeriod"`
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1FV string `json:"fd1FV"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2N string `json:"fd2N"`
  Fd2TimePeriod string `json:"fd2TimePeriod"`
  Fd2Interest string `json:"fd2Interest"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2PMT string `json:"fd2PMT"`
  Fd2Result string `json:"fd2Result"`
}

func newAdPvFields(dir1, dir2, correlationId string) *adPvFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  //Default values returned if file is missing, empty, or JSON is corrupt.
  m := adPvFields {
    MenuPage: "",
    CurrentPage: "rhs-ui2",
		CurrentButton: "lhs-button2",
    //
    Fd1N: "1.0",
    Fd1TimePeriod: "year",
    Fd1Interest: "1.00",
    Fd1Compound: "monthly",
    Fd1FV: "1.00",
    Fd1Result: "",
    //
    Fd2N: "1.0",
    Fd2TimePeriod: "year",
    Fd2Interest: "1.00",
    Fd2Compound: "monthly",
    Fd2PMT: "1.00",
    Fd2Result: "",
  }
  obj, err := readFields(dir + "adpv.txt")
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
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "adpv.txt"), correlationId)
  }
  return &m
}

func getAdPvFields(userName string) *adPvFields {
  return currentFields[userName].adPv
}

type WfAdPvPages struct {}

func (a WfAdPvPages) AdPvPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering webfinances.AdPvPages.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res, correlationId)
    return
  }
  //
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    fields := getAdPvFields(userName)
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
    // if strings.EqualFold(p.CurrentPage, "rhs-ui1") {
    //   p.CurrentButton = "lhs-button1"
      // if req.Method == http.MethodPost {
      //   p.fd1N = req.PostFormValue("fd1-n")
      //   p.fd1TimePeriod = req.PostFormValue("fd1-tp")
      //   p.fd1Interest = req.PostFormValue("fd1-interest")
      //   p.fd1Compound = req.PostFormValue("fd1-cp")
      //   p.fd1FV = req.PostFormValue("fd1-fv")
      //   var n float64
      //   var i float64
      //   var fv float64
      //   var err error
      //   if n, err = strconv.ParseFloat(p.fd1N, 64); err != nil {
      //     p.Fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1N, err)
      //   } else if i, err = strconv.ParseFloat(p.fd1Interest, 64); err != nil {
      //     p.Fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1Interest, err)
      //   } else if fv, err = strconv.ParseFloat(p.fd1FV, 64); err != nil {
      //     p.Fd1Result = fmt.Sprintf("Error: %s -- %+v", p.fd1FV, err)
      //   } else {
      //     var oa finances.Annuities
      //     p.Fd1Result = fmt.Sprintf("Present Value: $%.2f", oa.O_PresentValue_FV(fv, i / 100.0,
      //                               oa.GetCompoundingPeriod(p.fd1Compound[0], true),
      //                               n, oa.GetTimePeriod(p.fd1TimePeriod[0], true)))
      //   }
      //   logEntry.Print(INFO, correlationId, []string {
      //     fmt.Sprintf("n = %s, tp = %s, i = %s, cp = %s, fv = %s, %s",
      //                 p.fd1N, p.fd1TimePeriod, p.fd1Interest, p.fd1Compound, p.fd1FV, p.Fd1Result),
      //   })
      // }
      // t := template.Must(template.ParseFiles("webfinances/templates/xxxordinaryannuity/pv/pv.html",
      //                                        "webfinances/templates/header.html",
      //                                        "webfinances/templates/xxxordinaryannuity/pv/n-i-FV.html",
      //                                        "webfinances/templates/footer.html"))
      // t.ExecuteTemplate(res, "oapresentvalue", struct {
      //   Header string
      //   Datetime string
      //   CurrentButton string
      //   Fd1N string
      //   Fd1TimePeriod string
      //   Fd1Interest string
      //   Fd1Compound string
      //   Fd1FV string
      //   Fd1Result string
      // } { "Annuity Due / Present Value", m.DTF(), p.CurrentButton,
      //     p.fd1N, p.fd1TimePeriod, p.fd1Interest, p.fd1Compound, p.fd1FV, p.Fd1Result,
      //   })
    } else*/
    if strings.EqualFold(fields.CurrentPage, "rhs-ui2") {
      fields.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {
        fields.Fd2N = req.FormValue("fd2-n")
        fields.Fd2TimePeriod = req.PostFormValue("fd2-tp")
        fields.Fd2Interest = req.PostFormValue("fd2-interest")
        fields.Fd2Compound = req.PostFormValue("fd2-cp")
        fields.Fd2PMT = req.PostFormValue("fd2-pmt")
        var n float64
        var i float64
        var pmt float64
        var err error
        if n, err = strconv.ParseFloat(fields.Fd2N, 64); err != nil {
          fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2N, err)
        } else if i, err = strconv.ParseFloat(fields.Fd2Interest, 64); err != nil {
          fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2Interest, err)
        } else if pmt, err = strconv.ParseFloat(fields.Fd2PMT, 64); err != nil {
          fields.Fd2Result = fmt.Sprintf("Error: %s -- %+v", fields.Fd2PMT, err)
        } else {
          var oa finances.Annuities
          fields.Fd2Result = fmt.Sprintf("Present Value: $%.5f", oa.D_PresentValue_PMT(pmt, i / 100.0,
            oa.GetCompoundingPeriod(fields.Fd2Compound[0], true), n, oa.GetTimePeriod(fields.Fd2TimePeriod[0], true)))
        }
        logger.LogInfo(fmt.Sprintf("n = %s, tp = %s, interest = %s, cp = %s, pmt = %s, %s", fields.Fd2N,
          fields.Fd2TimePeriod, fields.Fd2Interest, fields.Fd2Compound, fields.Fd2PMT, fields.Fd2Result), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/finances/annuitydue/pv/pv.html",
        "webfinances/templates/finances/annuitydue/pv/n-i-PMT.html",
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
          Fd2N string
          Fd2TimePeriod string
          Fd2Interest string
          Fd2Compound string
          Fd2PMT string
          Fd2Result string
        } { "standard", "Annuity Due / Present Value", logger.DatetimeFormat(), financesMenuPage, fields.CurrentButton, newSession.CsrfToken,
            fields.Fd2N, fields.Fd2TimePeriod, fields.Fd2Interest, fields.Fd2Compound, fields.Fd2PMT, fields.Fd2Result },
      })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", fields.CurrentPage)
      logger.LogError(errString, correlationId)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      logger.LogWarning("*** Request timeout ***", correlationId)
      if strings.EqualFold(fields.CurrentPage, "rhs-ui1") {
        fields.Fd1Result = ""
      } else if strings.EqualFold(fields.CurrentPage, "rhs-ui2") {
        fields.Fd2Result = ""
      }
    }
    //
    if data, err := json.Marshal(fields); err != nil {
      logger.LogError(fmt.Sprintf("%+v", err), correlationId)
    } else {
      filePath := fmt.Sprintf("%s/%s/adpv.txt", mainDir, userName)
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
