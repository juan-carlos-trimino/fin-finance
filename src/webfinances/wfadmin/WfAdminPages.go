package wfadmin

import (
  bank "finance/databases/banking" //Importing a package and assigning it a local alias.
  "context"
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/gpsessions"
  "net/http"
  "strings"
  "time"
)

/***
When handling authentication errors, the application should not disclose which part of the
authentication data was incorrect. Instead of "Invalid username" or "Invalid password", just use
"Invalid username and/or password" interchangeably.
***/
func invalidSession(res http.ResponseWriter) {
  templatesNeeded := []string{
    "webfinances/templates/layout-simple.html",
    "webfinances/templates/login.html",
  }
  renderer.Render(res, "layout-simple", templatesNeeded, renderer.PageData{
    Data: struct{
      Header string
      ErrMsg string
    } { "Login", "Invalid username and/or password" },
  })
}

func displayWelcomePage(res http.ResponseWriter) {
  templatesNeeded := []string{
    "webfinances/templates/layout-no-navbar.html",
    "webfinances/templates/admin/welcome.html",
    "webfinances/templates/title.html",
    "webfinances/templates/datetime.html",
    "webfinances/templates/footer.html",
  }
  renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
    Data: struct{
      Header string
      Datetime string
    } { "Investments - Admin", logger.DatetimeFormat() },
  })
}

type WfAdminPages struct {}

func (s WfAdminPages) WelcomePage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfadmin.WelcomePage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    displayWelcomePage(res)
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}





func (s WfAdminPages) UsersPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfadmin.UsersPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    cookie := sessions.CreateCookie(newSessionToken)
    http.SetCookie(res, cookie)
    templatesNeeded := []string{
      "webfinances/templates/layout-no-navbar.html",
      "webfinances/templates/admin/users.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/footer.html",
    }
    renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
      Data: struct{
        Header string
        Datetime string
        CsrfToken string
      } { "Register User - Admin", logger.DatetimeFormat(), newSession.CsrfToken, },
    })
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}













func (s WfAdminPages) RegisterPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfadmin.RegisterPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
    cookie := sessions.CreateCookie(newSessionToken)
    http.SetCookie(res, cookie)
    templatesNeeded := []string{
      "webfinances/templates/layout-no-navbar.html",
      "webfinances/templates/admin/register.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/footer.html",
    }
    renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
      Data: struct{
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
      } { "Register User - Admin", logger.DatetimeFormat(), newSession.CsrfToken, "", "", "", "",
          "", "false", time.Now().Format("2006-01-02"), "male", "", "", "", "", "", "", "",
          "", "" },
    })
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}





func (s WfAdminPages) SaveRegisterPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfadmin.SaveRegisterPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else if req.Method == http.MethodPost || req.Method == http.MethodGet {
    clickedButton := req.FormValue("button_action")  //Return either "back" or "register".
    if clickedButton == "back" {
      displayWelcomePage(res)
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
        displayWelcomePage(res)
      } else {
        newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
        cookie := sessions.CreateCookie(newSessionToken)
        http.SetCookie(res, cookie)
        templatesNeeded := []string{
          "webfinances/templates/layout-no-navbar.html",
          "webfinances/templates/admin/register.html",
          "webfinances/templates/title.html",
          "webfinances/templates/datetime.html",
          "webfinances/templates/footer.html",
        }
        renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
          Data: struct{
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
          } { "Register User - Admin", logger.DatetimeFormat(), newSession.CsrfToken, c.User_name, c.Password, c.First_name, middle_name,
              c.Last_name, marketing, originalDate, c.Gender, c.Address1, address2, c.City, c.State, c.Country, zip_code, c.Email,
              c.Phone, ok.Error() },
        })
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    logger.LogError(errString, correlationId)
    panic(errString)
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}

func (s WfAdminPages) AdminSettingsPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering AdminSettingsPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
  } else {
    templatesNeeded := []string{
      "webfinances/templates/layout-no-navbar.html",
      "webfinances/templates/admin/settings/settings.html",
      "webfinances/templates/title.html",
      "webfinances/templates/datetime.html",
      "webfinances/templates/footer.html",
    }
    renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
      Data: struct{
        Header string
        Datetime string
      } { "Investments - Admin", logger.DatetimeFormat() },
    })
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}
