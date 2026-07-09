package webfinances

import (
  "context"
  "encoding/json"
  "finance/databases/banking"
  "fmt"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/gposu"
  "github.com/juan-carlos-trimino/gpsessions"
  "html/template"
  "net/http"
  "os"
  "strings"
  "time"
)

type WfSecurityPages struct {
}

type changePassword struct {
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
}

func (s WfSecurityPages) AdminSecurityPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
    return
  }
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering AdminSecurityPages.", correlationId)
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    fields := changePassword {
      "Change Password", logger.DatetimeFormat(), "rhs-ui1", "lhs-button1", "", "", "", "", "", "",
    }
    var errorMsg string = ""
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
    //
    if strings.EqualFold(fields.CurrentPage, "rhs-ui1") {
      fields.CurrentButton = "lhs-button1"
      if req.Method == http.MethodPost {
        fields.Username = req.PostFormValue("un")
        fields.Old = req.PostFormValue("oldpwd")
        fields.New = req.PostFormValue("newpwd")
        fields.Confirm = req.PostFormValue("connewpwd")
        if strings.EqualFold(fields.New, fields.Confirm) {
          ok := banking.DbChangePassword(req.Context(), fields.Username, fields.Old, fields.New, correlationId)
          if ok {
            errorMsg = "Your password has been successfully updated!"
          } else {
            errorMsg = "Your password was NOT successfully updated!"
          }
        } else {
          errorMsg = "New password and confirmation password do not match."
        }
        fields.Username = ""
        fields.Old = ""
        fields.New = ""
        fields.Confirm = ""
        logger.LogInfo(fmt.Sprintf("%s", errorMsg), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      /***
      The Must function wraps around the ParseGlob function that returns a pointer to a template and an error, and it panics if the
      error is not nil.
      ***/
      t := template.Must(template.ParseFiles(
        "webfinances/templates/admin/settings/security/security.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/admin/settings/security/password.html",
        "webfinances/templates/footer.html"))
      err := t.ExecuteTemplate(res, "admin_security", struct {
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
      } { "Change Password", logger.DatetimeFormat(), fields.CurrentPage, fields.CurrentButton, newSession.CsrfToken, fields.Username,
          fields.Old, fields.New, fields.Confirm, errorMsg,
      })
      //
      if err != nil {
        logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
      }
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
