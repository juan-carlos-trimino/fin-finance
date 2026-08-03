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

type WfBankingMngAcctsPages struct {
}

type Row struct {  //Rows for the accounts.
  AccountName string
  AccountType string
}

func (b WfBankingMngAcctsPages) ManageAccountsPages(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  sessionToken, _ := ctxKey.GetSessionToken(req.Context())
  if sessionToken == "" {
    invalidSession(res)
    return
  }
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.", startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering wfbanking.ManageAccountsPages.", correlationId)
  if req.Method == http.MethodPost || req.Method == http.MethodGet {
    userName := sessions.GetUserName(sessionToken)
    maf := getManageAccountsFields(userName)
    /***
    The functions in Request that allow to extract data from the URL and/or the body revolve around
    the Form, PostForm, and MultipartForm fields; the data are in the form of key-value pairs.

    If the form and the URL have the same key name, both of them will be placed in a slice, with
    the form value always prioritized before the URL value.

    Since we want the form key-value pairs, we can ignore the URL key-value pairs. The PostForm
    field provides key-value pairs only for the form and not the URL. The PostForm field supports
    only application/x-www-form-urlencoded.

    The FormValue method lets you access the key-value pairs just like the Form field, except that
    it's for a specific key and there is no need to call the ParseForm method beforehand -- the
    FormValue method does it. The PostFormValue method does the same thing, except that it's for
    the PostForm field instead of the Form field.
    ***/
    if ui := req.FormValue("manageaccounts"); ui != "" {  //Values from form and URL.
      maf.CurrentPage = ui
    }
    //
    if strings.EqualFold(maf.CurrentPage, "rhs-ui1") {
      maf.CurrentButton = "lhs-button1"
      if req.Method == http.MethodPost {
        maf.Fd1BankName = req.PostFormValue("fd1-bankname")
        maf.Fd1AccountType = req.PostFormValue("fd1-accounttype")
        maf.Fd1AccountName = req.PostFormValue("fd1-accountname")
        maf.Fd1AccountNumber = req.PostFormValue("fd1-accountnumber")
        maf.Fd1RoutingNumber = req.PostFormValue("fd1-routingnumber")
        //save to db
        logger.LogInfo(fmt.Sprintf("bankname = %s, accounttype = %s, accountname = %s, accountnumber = %s, routingnumber = %s",
          maf.Fd1BankName, maf.Fd1AccountType, maf.Fd1AccountName, maf.Fd1AccountNumber, maf.Fd1RoutingNumber), correlationId)
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
          Header string
          Datetime string
          CurrentPage string
          CurrentButton string
          CsrfToken string
          Fd1BankName string
          Fd1AccountType string
          Fd1AccountName string
          Fd1AccountNumber string
          Fd1RoutingNumber string
        } { "Manage Accounts", logger.DatetimeFormat(), "welcome", maf.CurrentButton, newSession.CsrfToken, maf.Fd1BankName,
            maf.Fd1AccountType, maf.Fd1AccountName, maf.Fd1AccountNumber, maf.Fd1RoutingNumber },
      })
    } else if strings.EqualFold(maf.CurrentPage, "rhs-ui2") {
      maf.CurrentButton = "lhs-button2"
      if req.Method == http.MethodPost {

          rowId := req.PostFormValue("selected_id")

          numberOfRows := 80
          maf.Fd2Result = make([]Row, 0, numberOfRows + 1)
          maf.Fd2Result = append(maf.Fd2Result,
            Row {
              AccountName: "--",
              AccountType: "savings",
            })
          for idx := 0; idx < numberOfRows; idx++ {
            maf.Fd2Result = append(maf.Fd2Result,
              Row {
                AccountName: fmt.Sprintf("account name%d", idx + 1),
                AccountType: "checking",
              })
          }


          if rowId != "" {
            logger.LogInfo(fmt.Sprintf("Account ID = %s", rowId), correlationId);
            for i, r := range maf.Fd2Result {
                if r.AccountName == rowId {
                    // Remove the element and maintain order
                    maf.Fd2Result = append(maf.Fd2Result[:i], maf.Fd2Result[i+1:]...)
                    break // Stop searching after the first match
                }
            }
          }



      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/banking/manageaccounts/manageaccounts.html",
        "webfinances/templates/banking/manageaccounts/deleteaccount.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/navbar.html",
        "webfinances/templates/footer.html",
      }
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
        Data: struct {
          Header string
          Datetime string
          CurrentPage string
          CurrentButton string
          CsrfToken string
          Fd2Result []Row
        } { "Manage Accounts", logger.DatetimeFormat(), "welcome", maf.CurrentButton, newSession.CsrfToken, maf.Fd2Result },
      })
    } else {
      errString := fmt.Sprintf("Unsupported page: %s", maf.CurrentPage)
      logger.LogError(errString, correlationId)
      panic(errString)
    }
    //
    if req.Context().Err() == context.DeadlineExceeded {
      logger.LogWarning("*** Request timeout ***", correlationId)
      // if strings.EqualFold(maf.CurrentPage, "rhs-ui1") {
      //   bf.Fd1Result = ""
      // } else if strings.EqualFold(bf.CurrentPage, "rhs-ui2") {
      //   bf.Fd2Result = ""
      // } else if strings.EqualFold(bf.CurrentPage, "rhs-ui3") {
      //   bf.Fd3Result = ""
      // } else if strings.EqualFold(bf.CurrentPage, "rhs-ui4") {
      //   bf.Fd4Result[0] = ""
      //   bf.Fd4Result[1] = ""
      // } else if strings.EqualFold(bf.CurrentPage, "rhs-ui5") {
      //   bf.Fd5Result[0] = ""
      //   bf.Fd5Result[1] = ""
      //   bf.Fd5Result[2] = ""
      //   bf.Fd5Result[3] = ""
      //   bf.Fd5Result[4] = ""
      //   bf.Fd5Result[5] = ""
      //   bf.Fd5Result[6] = ""
      // // } else if strings.EqualFold(bf.CurrentPage, "rhs-ui6") {
      // //   bf.Fd6Result[1] = ""
      // // } else if strings.EqualFold(bf.CurrentPage, "rhs-ui7") {
      // //   bf.Fd7Result[1] = ""
      // // } else if strings.EqualFold(bf.CurrentPage, "rhs-ui8") {
      // //   bf.Fd8Result[1] = ""
      // }
    }
    //
    if data, err := json.Marshal(maf); err != nil {
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
