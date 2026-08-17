package wfadmin

import (
  "context"
  "encoding/json"
  "finance/databases/banking"
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

type settingsFields struct{
  CurrentPage string  `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
}

func newSettingsFields(dir1, dir2, correlationId string) *settingsFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  //Default values returned if file is missing, empty, or JSON is corrupt.
  m := settingsFields {
    CurrentButton: "lhs-button1",
    CurrentPage: "rhs-ui1",
  }
  obj, err := readFields(dir + "settings.txt")
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
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "settings.txt"), correlationId)
  }
  return &m
}

func getSettingsFields(userName string) *settingsFields {
  return currentFields[userName].settings
}

type WfAdminSettingsPages struct{}

func (s WfAdminSettingsPages) AdminSettingsPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfadmin.AdminSettingsPages.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res, correlationId)
    return
  }
  //
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    fields := getSettingsFields(userName)
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
    if ui := req.FormValue("db"); ui != "" {  //Values from form and URL.
      fields.CurrentPage = ui
    }
    //
    if strings.EqualFold(fields.CurrentPage, "rhs-ui1") {
      fields.CurrentButton = "lhs-button1"
      pd := struct{
        LayoutType string
        Header string
        Datetime string
        CurrentButton string
        CsrfToken string
        Username string
        Old string
        New string
        Confirm string
        ErrMsg string
      } { "std-wo-nav-menu", "Settings - Admin", logger.DatetimeFormat(), fields.CurrentButton, "", "", "", "", "", "" }
      if req.Method == http.MethodPost {
        username := req.PostFormValue("un")
        old := req.PostFormValue("oldpwd")
        new := req.PostFormValue("newpwd")
        confirm := req.PostFormValue("connewpwd")
        if strings.EqualFold(new, confirm) {
          ok := banking.DbChangePassword(req.Context(), username, old, new, correlationId)
          if ok {
            pd.ErrMsg = "Your password has been successfully updated!"
          } else {
            pd.ErrMsg = "Your password was NOT successfully updated!"
          }
        } else {
          pd.ErrMsg = "New password and confirmation password do not match."
        }
        logger.LogInfo(fmt.Sprintf("%s", pd.ErrMsg), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/admin/settings/security/security.html",
        "webfinances/templates/admin/settings/security/password.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/footer.html",
      }
      pd.CsrfToken = newSession.CsrfToken
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{ Data: pd})
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", fields.CurrentPage)
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
      filePath := fmt.Sprintf("%s/%s/settings.txt", mainDir, userName)
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
