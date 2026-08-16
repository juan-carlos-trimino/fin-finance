package wfbanking

import (
  "context"
  "encoding/json"
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/gposu"
  "github.com/juan-carlos-trimino/gpsessions"
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

func (b WfBankingMngAcctsPages) ManageAccountsPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfbanking.ManageAccountsPages.", correlationId)
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res, correlationId)
    return
  }
  //
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
    if ui := req.FormValue("tablestyle"); ui != "" {  //Values from form and URL.
      fields.CurrentPage = ui
    }
    //
    if strings.EqualFold(fields.CurrentPage, "rhs-ui1") {
      fields.CurrentButton = "lhs-button1"
      var fd1BankName,
          fd1AccountType,
          fd1AccountName,
          fd1AccountNumber,
          fd1RoutingNumber string
      if req.Method == http.MethodPost {
        fd1BankName = req.PostFormValue("fd1-bankname")
        fd1AccountType = req.PostFormValue("fd1-accounttype")
        fd1AccountName = req.PostFormValue("fd1-accountname")
        fd1AccountNumber = req.PostFormValue("fd1-accountnumber")
        fd1RoutingNumber = req.PostFormValue("fd1-routingnumber")
        //save to db
        logger.LogInfo(fmt.Sprintf("bankname = %s, accounttype = %s, accountname = %s, accountnumber = %s, routingnumber = %s",
          fd1BankName, fd1AccountType, fd1AccountName, fd1AccountNumber, fd1RoutingNumber), correlationId)
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/banking/manageaccounts/manageaccounts.html",
        "webfinances/templates/banking/manageaccounts/createaccount.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/navbar.html",
        "webfinances/templates/footer.html",
      }
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
        Data: struct {
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
        } { "standard", "Manage Accounts", logger.DatetimeFormat(), bankingMenuPage, fields.CurrentButton, newSession.CsrfToken,
            fd1BankName, fd1AccountType, fd1AccountName, fd1AccountNumber, fd1RoutingNumber },
      })
    } else if strings.EqualFold(fields.CurrentPage, "rhs-ui2") {
      fields.CurrentButton = "lhs-button2"
      var rows []Row
      // if req.Method == http.MethodPost {



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



      // }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/banking/manageaccounts/manageaccounts.html",
        "webfinances/templates/banking/manageaccounts/deleteaccount.html",
        "webfinances/templates/helpers/table-container.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/navbar.html",
        "webfinances/templates/footer.html",
      }
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
        Data: struct {
          LayoutType string
          Header string
          Datetime string
          MenuPage string
          CurrentButton string
          CsrfToken string
          Fd2Result []Row
        } { "standard", "Manage Accounts", logger.DatetimeFormat(), bankingMenuPage, fields.CurrentButton,
            newSession.CsrfToken, rows },
      })
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
      filePath := fmt.Sprintf("%s/%s/manageaccounts.txt", mainDir, userName)
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
