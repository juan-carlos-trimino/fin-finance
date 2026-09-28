package wfadmin

import (
  "encoding/json"
  "finance/databases/banking"
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/go-logger"
  "github.com/juan-carlos-trimino/go-os"
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
  ck := middlewares.MwContextKey{}
  correlationId, _ := ck.GetCorrelationId(req.Context())
  startTime, _ := ck.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfadmin.AdminSettingsPages.", correlationId)
  sessInfo, _ := ck.GetSessionInfo(req.Context())
  fields := getSettingsFields(sessInfo.UserName)
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
  if ui := req.FormValue("db"); ui != "" {
    fields.CurrentPage = ui
  }
  //Dynamic variables determined by the routing route condition.
  var partialTemplate string
  var templateData interface{}
  //Core application state splitting via explicit page mapping.
  switch strings.ToLower(fields.CurrentPage) {
  case "rhs-ui1":
    fields.CurrentButton = "lhs-button1"
    partialTemplate = "password.html"
    td := struct{
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
    }{
      "std-wo-nav-menu",
      "Settings - Admin",
      logger.DatetimeFormat(),
      fields.CurrentButton,
      "",
      "",
      "",
      "",
      "",
      "",
    }
    //
    if req.Method == http.MethodPost {
      username := req.PostFormValue("un")
      old := req.PostFormValue("oldpwd")
      new := req.PostFormValue("newpwd")
      confirm := req.PostFormValue("connewpwd")
      if strings.EqualFold(new, confirm) {
        ok := banking.DbChangePassword(req.Context(), username, old, new, correlationId)
        if ok {
          td.ErrMsg = "Your password has been successfully updated!"
        } else {
          td.ErrMsg = "Your password was NOT successfully updated!"
        }
      } else {
        td.ErrMsg = "New password and confirmation password do not match."
      }
      td.CsrfToken = sessInfo.CSRFToken
      logger.LogInfo(fmt.Sprintf("%s", td.ErrMsg), correlationId)
    }
    templateData = td
  default:
    errString := fmt.Sprintf("Unsupported page: %s", fields.CurrentPage)
    logger.LogError(errString, correlationId)
    panic(errString)
  }
  //Unified execution of templates.
  templatesNeeded := []string{
    "webfinances/templates/layout.html",
    "webfinances/templates/admin/settings/security/security.html",
    "webfinances/templates/admin/settings/security/" + partialTemplate,
    "webfinances/templates/title.html",
    "webfinances/templates/datetime.html",
    "webfinances/templates/footer.html",
  }
  renderer.Render(res, "layout", templatesNeeded, renderer.PageData{ Data: templateData})
  data, err := json.Marshal(fields)  //Preserve current choices.
  if err != nil {
    //Don't crash the server (panic), but log it clearly so you can debug the serialization.
    logger.LogError(fmt.Sprintf("Failed to marshal fields to JSON for user %s: %+v", sessInfo.UserName, err), correlationId)
  } else {
    go func(userData []byte, uName, cId string) {
      filePath := fmt.Sprintf("%s/%s/adcp.txt", mainDir, uName)
      //The exclusive OS file lock handles goroutine collisions.
      _, err := osu.WriteAllExclusiveLock1(filePath, userData, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o600)
      //Check the error returned from the lock-writing function.
      if err != nil {
        logger.LogError(fmt.Sprintf("Goroutine file system error writing state to %s: %+v", filePath, err), cId)
        return
      }
      logger.LogInfo(fmt.Sprintf("Goroutine successfully persisted state to %s.", filePath), cId)
    }(data, sessInfo.UserName, correlationId) // Pass variables into the closure to prevent scope races
  }
  logger.LogInfo(fmt.Sprintf("Request took %vms", time.Since(startTime).Microseconds()), correlationId)
}
