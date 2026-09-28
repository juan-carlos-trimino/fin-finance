package wfbanking

import (
  "encoding/json"
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/go-logger"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/go-os"
  "net/http"
  "os"
  "strings"
  "time"
)

type manageAccountsFields struct {
  CurrentButton string `json:"currentButton"`
  CurrentPage string  `json:"currentPage"`
}

func newManageAccountsFields(dir1, dir2, correlationId string) *manageAccountsFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  //Default values returned if file is missing, empty, or JSON is corrupt.
  m := manageAccountsFields{
    CurrentButton: "lhs-button1",
    CurrentPage: "rhs-ui1",
  }
  obj, err := readFields(dir + "manageaccounts.txt")
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
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "manageaccounts.txt"), correlationId)
  }
  return &m
}

func getManageAccountsFields(userName string) *manageAccountsFields {
  return currentFields[userName].manageAccounts
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
  startTime, _ := ck.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfbanking.ManageAccountsPages.", correlationId)
  sessInfo, _ := ck.GetSessionInfo(req.Context())
  fields := getManageAccountsFields(sessInfo.UserName)
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
    logger.LogError(fmt.Sprintf("Failed to marshal fields to JSON for user %s: %+v", sessInfo.UserName, err), correlationId)
  } else {
    go func(userData []byte, uName, cId string) {
      filePath := fmt.Sprintf("%s/%s/manageaccounts.txt", mainDir, uName)
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
  logger.LogInfo(fmt.Sprintf("Request took %vms\n", time.Since(startTime).Microseconds()), correlationId)
}

//Extraction helper for UI-1 calculations.
func (b WfBankingMngAcctsPages) processUi1Form(req *http.Request, fields *ui1Fields) {
  fields.fd1BankName = req.PostFormValue("fd1-bankname")
  fields.fd1AccountType = req.PostFormValue("fd1-accounttype")
  fields.fd1AccountName = req.PostFormValue("fd1-accountname")
  fields.fd1AccountNumber = req.PostFormValue("fd1-accountnumber")
  fields.fd1RoutingNumber = req.PostFormValue("fd1-routingnumber")
}
