package wfadmin

import (
  bank "finance/databases/banking" //Importing a package and assigning it a local alias.
  "context"
  "encoding/json"
  "finance/databases/banking"
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/go-logger"
  sess "github.com/juan-carlos-trimino/go-sessions"
  "net/http"
  "strings"
)

type settingsFields struct{
  CurrentPage string  `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
}

func defaultSettingsFields() *settingsFields {
  return &settingsFields{
    CurrentButton: "lhs-button1",
    CurrentPage: "rhs-ui1",
  }
}

type WfAdminSettingsPages struct{}

func (s WfAdminSettingsPages) AdminSettingsPages(res http.ResponseWriter, req *http.Request) {
  ck := middlewares.MwContextKey{}
  correlationId, _ := ck.GetCorrelationId(req.Context())
  logger.LogInfo("Entering wfadmin.AdminSettingsPages.", correlationId)
  sessInfo, _ := ck.GetSessionInfo(req.Context())
  redisKey := "user:data:" + sessInfo.UserName
  partitionName := "admin_settings"
  var fields *settingsFields = nil
  var jsonBytes []byte
  //Check Redis Cache First.
  jsonBytes, err := sess.HGetRedis(req.Context(), redisKey, partitionName)
  if err == nil {
    logger.LogInfo(fmt.Sprintf("Redis found key(%s)/partition(%s)...", redisKey, partitionName), correlationId)
    //Unmarshal unmarshals the JSON into the value pointed at by v. If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
    err = json.Unmarshal(jsonBytes, &fields)
    if err != nil {
      logger.LogInfo(fmt.Sprintf("Error (Redis) unmarshalling key(%s)/partition(%s): %v. Using default values.", redisKey, partitionName, err), correlationId)
      fields = defaultSettingsFields()
    }
  } else {
    //Cache Miss: Fall back to PostgreSQL database.
    logger.LogInfo(fmt.Sprintf("Redis does not have key(%s)/partition(%s): %v. Fetching from PostgreSQL...", redisKey, partitionName, err), correlationId)
    jsonBytes = bank.DbFetchUserData(req.Context(), sessInfo.UserName, partitionName, correlationId)
    if len(jsonBytes) == 0 {  //On error (nil) or not found (len() == 0), use default values.
      logger.LogInfo(fmt.Sprintf("Database did not find user(%s)/partition(%s). Using default values.", sessInfo.UserName, partitionName), correlationId)
      fields = defaultSettingsFields()
    } else {  //Use the values from the db.
      logger.LogInfo(fmt.Sprintf("Database found user(%s)/partition(%s)...", sessInfo.UserName, partitionName), correlationId)
      //Unmarshal unmarshals the JSON into the value pointed at by v. If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
      err = json.Unmarshal(jsonBytes, &fields)
      if err != nil {
        logger.LogInfo(fmt.Sprintf("Error (db) unmarshalling user(%s)/partition(%s): %v. Using default values.", sessInfo.UserName, partitionName, err), correlationId)
        fields = defaultSettingsFields()
      }
    }
  }
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
      sessInfo.CSRFToken,
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
      // td.CsrfToken = sessInfo.CSRFToken
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
    logger.LogError(fmt.Sprintf("Error marshalling user(%s)/partition(%s): %v", sessInfo.UserName, partitionName, err), correlationId)
  } else {
    go func(userData []byte, uName, redisKey, pName, cId string) {
      ctx := context.Background()
      err := sess.HSetRedis(ctx, uName, pName, userData)
      if err != nil {
        logger.LogInfo(fmt.Sprintf("Error (Redis) writing state for key(%s)/partition(%s).", redisKey, pName), cId)
      }
      ok := bank.DbSaveUserData(ctx, uName, pName, cId, userData)
      if ok {
        logger.LogInfo(fmt.Sprintf("Db successfully persisted state for user(%s)/partition(%s).", uName, pName), cId)
      }
    }(data, sessInfo.UserName, redisKey, partitionName, correlationId)
  }
}
