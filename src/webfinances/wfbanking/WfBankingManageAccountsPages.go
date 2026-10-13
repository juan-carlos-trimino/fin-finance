package wfbanking

import (
  bank "finance/databases/banking" //Importing a package and assigning it a local alias.
  "context"
  "encoding/json"
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/go-logger"
  "github.com/juan-carlos-trimino/go-middlewares"
  "net/http"
  sess "github.com/juan-carlos-trimino/go-sessions"
  "strings"
)

type manageAccountsFields struct {
  CurrentButton string `json:"currentButton"`
  CurrentPage string  `json:"currentPage"`
}

func defaultManageAccountsFields() *manageAccountsFields {
  return &manageAccountsFields{
    CurrentButton: "lhs-button1",
    CurrentPage: "rhs-ui1",
  }
}

type WfBankingMngAcctsPages struct {}

type Row struct {  //Rows for the accounts.
  AccountName string
  AccountType string
}

type ui1Fields struct {
  fd1BankName string
  fd1AccountType string
  fd1AccountName string
  fd1AccountNumber string
  fd1RoutingNumber string
}

func (b WfBankingMngAcctsPages) ManageAccountsPages(res http.ResponseWriter, req *http.Request) {
  ck := middlewares.MwContextKey{}
  correlationId, _ := ck.GetCorrelationId(req.Context())
  logger.LogInfo("Entering wfbanking.ManageAccountsPages.", correlationId)
  sessInfo, _ := ck.GetSessionInfo(req.Context())
  redisKey := "user:data:" + sessInfo.UserName
  partitionName := "manage_accounts"
  var fields *manageAccountsFields = nil
  var jsonBytes []byte
  //Check Redis Cache First.
  jsonBytes, err := sess.HGetRedis(req.Context(), redisKey, partitionName)
  if err == nil {
    logger.LogInfo(fmt.Sprintf("Redis found key(%s)/partition(%s)...", redisKey, partitionName), correlationId)
    //Unmarshal unmarshals the JSON into the value pointed at by v. If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
    err = json.Unmarshal(jsonBytes, &fields)
    if err != nil {
      logger.LogInfo(fmt.Sprintf("Error (Redis) unmarshalling key(%s)/partition(%s): %v. Using default values.", redisKey, partitionName, err), correlationId)
      fields = defaultManageAccountsFields()
    }
  } else {
    //Cache Miss: Fall back to PostgreSQL database.
    logger.LogInfo(fmt.Sprintf("Redis does not have key(%s)/partition(%s): %v. Fetching from PostgreSQL...", redisKey, partitionName, err), correlationId)
    jsonBytes = bank.DbFetchUserData(req.Context(), sessInfo.UserName, partitionName, correlationId)
    if len(jsonBytes) == 0 {  //On error (nil) or not found (len() == 0), use default values.
      logger.LogInfo(fmt.Sprintf("Database did not find user(%s)/partition(%s). Using default values.", sessInfo.UserName, partitionName), correlationId)
      fields = defaultManageAccountsFields()
    } else {  //Use the values from the db.
      logger.LogInfo(fmt.Sprintf("Database found user(%s)/partition(%s)...", sessInfo.UserName, partitionName), correlationId)
      //Unmarshal unmarshals the JSON into the value pointed at by v. If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
      err = json.Unmarshal(jsonBytes, &fields)
      if err != nil {
        logger.LogInfo(fmt.Sprintf("Error (db) unmarshalling user(%s)/partition(%s): %v. Using default values.", sessInfo.UserName, partitionName, err), correlationId)
        fields = defaultManageAccountsFields()
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
  if ui := req.FormValue("tablestyle"); ui != "" {
    fields.CurrentPage = ui
  }
  //Dynamic variables determined by the routing route condition.
  var partialTemplate string
  var templateData interface{}
  //Core application state splitting via explicit page mapping.
  switch strings.ToLower(fields.CurrentPage) {
  case "rhs-ui1":
    fields.CurrentButton = "lhs-button1"
    ui1 := ui1Fields{}
    if req.Method == http.MethodPost {
      b.processUi1Form(req, &ui1)
      //save to db... to be implemented
      logger.LogInfo(fmt.Sprintf("bankname = %s, accounttype = %s, accountname = %s, accountnumber = %s, routingnumber = %s",
       ui1.fd1BankName, ui1.fd1AccountType, ui1.fd1AccountName, ui1.fd1AccountNumber, ui1.fd1RoutingNumber), correlationId)
    }
    partialTemplate = "createaccount.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd1BankName string
      Fd1AccountType string
      Fd1AccountName string
      Fd1AccountNumber string
      Fd1RoutingNumber string
    }{
      "standard",
      "Manage Accounts",
      logger.DatetimeFormat(),
      bankingMenuPage,
      fields.CurrentButton,
      sessInfo.CSRFToken,
      ui1.fd1BankName,
      ui1.fd1AccountType,
      ui1.fd1AccountName,
      ui1.fd1AccountNumber,
      ui1.fd1RoutingNumber,
    }
  case "rhs-ui2":
    fields.CurrentButton = "lhs-button2"
    var rows []Row

          rowId := req.PostFormValue("selected_id")

          numberOfRows := 80
          rows = make([]Row, 0, numberOfRows + 1)
          rows = append(rows,
            Row {
              AccountName: "--",
              AccountType: "savings",
            })
          for idx := 0; idx < numberOfRows; idx++ {
            rows = append(rows,
              Row {
                AccountName: fmt.Sprintf("account name%d", idx + 1),
                AccountType: "checking",
              })
          }


          if rowId != "" {
            logger.LogInfo(fmt.Sprintf("Account ID = %s", rowId), correlationId);
            for i, r := range rows {
                if r.AccountName == rowId {
                    // Remove the element and maintain order
                    rows = append(rows[:i], rows[i+1:]...)
                    break // Stop searching after the first match
                }
            }
          }

    partialTemplate = "deleteaccount.html"
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      MenuPage string
      CurrentButton string
      CsrfToken string
      Fd2Result []Row
    }{
      "standard",
      "Manage Accounts",
      logger.DatetimeFormat(),
      bankingMenuPage,
      fields.CurrentButton,
      sessInfo.CSRFToken,
      rows,
    }
  default:
    errString := fmt.Sprintf("Unsupported page: %s", fields.CurrentPage)
    logger.LogError(errString, correlationId)
    panic(errString)
  }
  //Unified execution of templates.
  templatesNeeded := []string{
    "webfinances/templates/layout.html",
    "webfinances/templates/banking/manageaccounts/manageaccounts.html",
    "webfinances/templates/banking/manageaccounts/" + partialTemplate,
    "webfinances/templates/title.html",
    "webfinances/templates/datetime.html",
    "webfinances/templates/navbar.html",
    "webfinances/templates/footer.html",
  }
  renderer.Render(res, "layout", templatesNeeded, renderer.PageData{Data: templateData})
  data, err := json.Marshal(fields)
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

//Extraction helper for UI-1 calculations.
func (b WfBankingMngAcctsPages) processUi1Form(req *http.Request, fields *ui1Fields) {
  fields.fd1BankName = req.PostFormValue("fd1-bankname")
  fields.fd1AccountType = req.PostFormValue("fd1-accounttype")
  fields.fd1AccountName = req.PostFormValue("fd1-accountname")
  fields.fd1AccountNumber = req.PostFormValue("fd1-accountnumber")
  fields.fd1RoutingNumber = req.PostFormValue("fd1-routingnumber")
}
