package webfinances

import (
  bank "finance/databases/banking" //Importing a package and assigning it a local alias.
  "context"
  "fmt"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/gpsessions"
  "html/template"
  "net/http"
  "strings"
  "time"
)

type WfVerificationPages struct {}

func (s WfVerificationPages) AdminWelcomePage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering AdminWelcomePage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    t := template.Must(template.ParseFiles(
      "webfinances/templates/admin/admin_welcome.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/footer.html"))
    err := t.ExecuteTemplate(res, "admin_welcome_page", struct {
      Header string
      Datetime string
    } { "Investments", logger.DatetimeFormat() })
    //
    if err != nil {
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    }
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (s WfVerificationPages) AdminRegisterPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering AdminRegisterPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else if req.Method == http.MethodPost || req.Method == http.MethodGet {
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    cookie := sessions.CreateCookie(newSessionToken)
    http.SetCookie(res, cookie)
    t := template.Must(template.ParseFiles(
      "webfinances/templates/admin/admin_register.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/footer.html"))
    err := t.ExecuteTemplate(res, "admin_register_page", map[string]string {
      "Header": "Register User",
      "Datetime": logger.DatetimeFormat(),
      "CsrfToken": newSession.CsrfToken,
    })
    //
    if err != nil {
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    logger.LogError(errString, correlationId)
    panic(errString)
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (s WfVerificationPages) AdminSaveRegisterPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering AdminSaveRegisterPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else if req.Method == http.MethodPost || req.Method == http.MethodGet {
    clickedButton := req.FormValue("button_action")  //Return either "back" or "register".
    if clickedButton == "back" {
      t := template.Must(template.ParseFiles(
        "webfinances/templates/admin/admin_welcome.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/footer.html"))
      err := t.ExecuteTemplate(res, "admin_welcome_page", struct {
        Header string
        Datetime string
      } { "Investments", logger.DatetimeFormat() })
      //
      if err != nil {
        logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
      }
    } else {
      c := bank.AddCustomer {
        User_name: req.PostFormValue("uname"),
        Password: req.PostFormValue("pwd"),
        First_name: req.PostFormValue("fname"),
        Last_name: req.PostFormValue("lname"),
        Gender: req.PostFormValue("gender"),
        Address1: req.PostFormValue("address1"),
        City: req.PostFormValue("city"),
        State: req.PostFormValue("state"),
        Country: req.PostFormValue("country"),
        Email: req.PostFormValue("email"),
        Phone: req.PostFormValue("phone"),
      }
      marketing := req.PostFormValue("marketing")
      if strings.EqualFold(marketing, "true") {
        c.Marketing = true
      } else {
        c.Marketing = false
      }
      middle_name := req.PostFormValue("mname")
      c.Middle_name = bank.StringPtr(middle_name)
      address2 := req.PostFormValue("address2")
      c.Address2 = bank.StringPtr(address2)
      zip_code := req.PostFormValue("zip_code")
      c.Zip_code = bank.StringPtr(zip_code)
      originalDate := req.PostFormValue("bdate")
      /***
      Go's time formatting uses a reference date and time: Mon Jan 2 15:04:05 MST 2006. Each component of this reference
      time (e.g., 02 for the day, 01 for the month, 2006 for the year) is used as a placeholder in the layout string to
      match the input format; e.g., "dd/mm/yyyy" is "02/01/2006".
      ***/
      if newDate, err := time.Parse("2006-01-02", originalDate); err != nil {
        fmt.Println("Error parsing date: ", err)
      } else {
        c.Birth_date = bank.TimePtr(newDate)
      }
      ok := bank.DbAddCustomer(&c, context.Background(), correlationId)
      if ok == nil {
        t := template.Must(template.ParseFiles(
          "webfinances/templates/admin/admin_welcome.html",
          "webfinances/templates/title.html",
          "webfinances/templates/datetime.html",
          "webfinances/templates/footer.html"))
        err := t.ExecuteTemplate(res, "admin_welcome_page", struct {
          Header string
          Datetime string
        } { "Investments", logger.DatetimeFormat() })
        //
        if err != nil {
          logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
        }
      } else {
        newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
        cookie := sessions.CreateCookie(newSessionToken)
        http.SetCookie(res, cookie)
        /***
        The Must function wraps around the ParseGlob function that returns a pointer to a template
        and an error, and it panics if the error is not nil.
        ***/
        t := template.Must(template.ParseFiles(
          "webfinances/templates/admin/admin_register.html",
          "webfinances/templates/title.html",
          "webfinances/templates/datetime.html",
          "webfinances/templates/footer.html"))
        err := t.ExecuteTemplate(res, "admin_register_page", struct {
          Header string
          Datetime string
          CsrfToken string
          Username string
          Password string
          Fname string
          Mname string
          Lname string
          Marketing string
          Bdate string
          Gender string
          Address1 string
          Address2 string
          City string
          State string
          Country string
          Zip_Code string
          Email string
          Phone string
          ErrMsg string
        } { "Register User", logger.DatetimeFormat(), newSession.CsrfToken, c.User_name, c.Password, c.First_name, middle_name,
            c.Last_name, marketing, originalDate, c.Gender, c.Address1, address2, c.City, c.State, c.Country, zip_code, c.Email,
            c.Phone, ok.Error(),
        })
        //
        if err != nil {
          logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
        }
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    logger.LogError(errString, correlationId)
    panic(errString)
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (s WfVerificationPages) AdminSettingsPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering AdminSettingsPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    t := template.Must(template.ParseFiles(
      "webfinances/templates/admin/settings/settings.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/footer.html"))
    err := t.ExecuteTemplate(res, "admin_settings", struct {
      Header string
      Datetime string
    } { "Settings", logger.DatetimeFormat() })
    //
    if err != nil {
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    }
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (p WfVerificationPages) PublicSettingsSecurityFile(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering PublicSettingsSecurityFile.", correlationId)
  http.ServeFile(res, req, "./webfinances/public/js/admin/SettingsSecurity.js")
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}
