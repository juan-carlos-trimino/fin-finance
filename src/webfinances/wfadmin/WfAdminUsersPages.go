package wfadmin

import (
  bank "finance/databases/banking"  //Importing a package and assigning it a local alias.
  "context"
  "encoding/json"
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/gposu"
  "github.com/juan-carlos-trimino/gpsessions"
  "net/http"
  "os"
  "strings"
  "time"
)

type WfAdminUsersPages struct {}

func (u WfAdminUsersPages) AdminUsersPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfadmin.AdminUsersPage.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res, correlationId)
    return
  }
  //
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    fields := getUsersFields(userName)
    var currentRHS string = "rhs-ui1"  //Default.
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
    if ui := req.FormValue("save"); ui != "" {  //Values from form and URL.
      currentRHS = ui
    }
    //
    if strings.EqualFold(currentRHS, "rhs-ui1") {
      fields.CurrentButton = "lhs-button1"
      pd := struct{
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Username string
        Password string
        Fname string
        Mname string
        Lname string
        Gender string
        Bdate string
        Marketing string
        Address1 string
        Address2 string
        City string
        State string
        Country string
        Zip_Code string
        Email string
        Phone string
        ErrMsg string
        } { "Register User - Admin", logger.DatetimeFormat(), fields.CurrentButton, "", "", "", "", "", "", "male",
            time.Now().Format("2006-01-02"), "false", "", "", "", "", "", "", "", "", "", }
      if req.Method == http.MethodPost {
        c := bank.Customer {
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
        Go's time formatting uses a reference date and time: Mon Jan 2 15:04:05 MST 2006. Each component of this reference time (e.g.,
        02 for the day, 01 for the month, 2006 for the year) is used as a placeholder in the layout string to match the input format;
        e.g., "dd/mm/yyyy" is "02/01/2006".
        ***/
        newDate, err := time.Parse("2006-01-02", originalDate)
        if err != nil {
         // s := fmt.Sprintf("%v", err)
          logger.LogError(err.Error(), correlationId)
          c.Birth_date = bank.TimePtr(newDate)

        //  pd.ErrMsg = err.Error()
        } else {
          c.Birth_date = bank.TimePtr(newDate)
          err = bank.DbAddCustomer(&c, context.Background(), correlationId)
        }
        //
        if err != nil {
          pd.Username = c.User_name
          pd.Password = c.Password
          pd.Fname = c.First_name
          pd.Mname = bank.PtrString(c.Middle_name)
          pd.Lname = c.Last_name
          pd.Gender = c.Gender
          pd.Bdate = c.Birth_date.Format("2006-01-02")
          if c.Marketing {
            pd.Marketing = "true"
          } else {
            pd.Marketing = "false"
          }
          pd.Address1 = c.Address1
          pd.Address2 = bank.PtrString(c.Address2)
          pd.City = c.City
          pd.State = c.State
          pd.Country = c.Country
          pd.Zip_Code = bank.PtrString(c.Zip_code)
          pd.Email = c.Email
          pd.Phone = c.Phone
          pd.ErrMsg = fmt.Sprintf("%v", err)
        }
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout-no-navbar.html",
        "webfinances/templates/admin/users/users.html",
        "webfinances/templates/admin/users/register.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/footer.html",
      }
      pd.CsrfToken = newSession.CsrfToken
      /***
      Do not read or write sensitive information from the disk; use the database exclusively.
      ***/
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{ Data: pd})
    } else if strings.EqualFold(currentRHS, "rhs-ui2") {
      fields.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout-no-navbar.html",
        "webfinances/templates/admin/users/users.html",
        "webfinances/templates/admin/users/unregister.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/footer.html",
      }
      /***
      Do not read or write sensitive information from the disk; use the database exclusively.
      ***/
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
        Data: struct{
          Header string
          Datetime string
          CurrentButton string
          CsrfToken string
        } { "Unregister User - Admin", logger.DatetimeFormat(), fields.CurrentButton, newSession.CsrfToken },
      })


    } else {
      errString := fmt.Sprintf("Unsupported page: %s", currentRHS)
      logger.LogError(errString, correlationId)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      logger.LogWarning("*** Request timeout ***", correlationId)
    }
    //
    if data, err := json.Marshal(fields); err != nil {
      logger.LogError(fmt.Sprintf("%+v", err), correlationId)
    } else {
      filePath := fmt.Sprintf("%s/%s/users.txt", mainDir, userName)
      if _, err := osu.WriteAllExclusiveLock1(filePath, data, os.O_CREATE | os.O_RDWR | os.O_TRUNC, 0o600); err != nil {
        logger.LogError(fmt.Sprintf("%+v", err), correlationId)
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    logger.LogError(errString, correlationId)
    panic(errString)
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}




/*
func (s WfAdminUsersPages) AdminUsersPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
    return
  }
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfadmin.AdminUsersPages.", correlationId)
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    fields := getManageAccountsFields(userName)
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
/*    if ui := req.FormValue("compute"); ui != "" {  //Values from form and URL.
      fields.CurrentPage = ui
    }
    //
    if strings.EqualFold(fields.CurrentPage, "rhs-ui1") {
      fields.CurrentButton = "lhs-button1"
      if req.Method == http.MethodPost {
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout-no-navbar.html",
        "webfinances/templates/admin/settings/security/security.html",
        "webfinances/templates/admin/settings/security/password.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/footer.html",
      }
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
        Data: struct{
          Header string
          Datetime string
          CurrentPage string
          CurrentButton string
          CsrfToken string
          Username string
          Old string
          New string
          Confirm string
          ErrMsg string
        } { "Change Password - Admin", logger.DatetimeFormat(), fields.CurrentPage, fields.CurrentButton, newSession.CsrfToken, fields.Username,
            fields.Old, fields.New, fields.Confirm, errorMsg },
      })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", fields.CurrentPage)
      logger.LogError(errString, "-1")
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      logger.LogWarning("*** Request timeout ***", "-1")
      if strings.EqualFold(fields.CurrentPage, "rhs-ui1") {
        fields.ErrMsg = ""
      }
    }
    //
    if data, err := json.Marshal(fields); err != nil {
      logger.LogError(fmt.Sprintf("%+v", err), "-1")
    } else {
      filePath := fmt.Sprintf("%s/%s/mortgage.txt", mainDir, userName)
      if _, err := osu.WriteAllExclusiveLock1(filePath, data, os.O_CREATE | os.O_RDWR | os.O_TRUNC, 0o600); err != nil {
        logger.LogError(fmt.Sprintf("%+v", err), "-1")
      }
    }
  } else {
    errString := fmt.Sprintf("Unsupported method: %s", req.Method)
    logger.LogError(errString, "-1")
    panic(errString)
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}
*/
