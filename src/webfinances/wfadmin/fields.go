package wfadmin

import (
  "encoding/json"
  "errors"
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/gposu"
  "os"
)

var mainDir string = "/admin"

func SetupDirStructure(dir string) {
  mainDir = dir + mainDir
  dirErr, err := osu.CreateDirs(0o077, 0o777, mainDir)
  if err != nil {
    panic("Cannot create directory '" + dirErr + "': " + err.Error())
  }
}

//Store the fields for each user in memory.
var currentFields = map[string]*fields{}  //key: user, value: fields

/***
Browsers default to use the same collection of cookies regardless of whether you are opening a duplicate web page in a
new tab or a new browser instance. Hence, two different tabs or two instances of the same browser will look like the
same session to the server.

Because multiple instances use the same cookies, the server cannot tell requests from them apart, and it will associate
them with the same Session data because they all have the same SessionID.
***/
type fields struct {
  //Make the pointers unexported so that clients can't interact with them directly but only via exported methods.
  manageAccounts *manageAccountsFields
  users *usersFields
}

type usersFields struct {
  CurrentButton string `json:"currentButton"`
  CurrentPage string  `json:"currentPage"`
  SelectedRange string `json:"selectedRange"`
}

func newUsersFields(dir1, dir2, correlationId string) *usersFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  //Default values returned if file is missing, empty, or JSON is corrupt.
  m := usersFields {
    CurrentButton: "lhs-button1",
    CurrentPage: "rhs-ui1",
    SelectedRange: "1",
  }
  obj, err := readFields(dir + "users.txt")
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
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "users.txt"), correlationId)
  }
  return &m
}

func getUsersFields(userName string) *usersFields {
  return currentFields[userName].users
}






type manageAccountsFields struct {
  CurrentButton string `json:"currentButton"`
  //
  Fd1BankName string `json:"fFd1BankName"`
  Fd1AccountType string `json:"fd1AccountType"`
  Fd1AccountName string `json:"fd1AccountName"`
  Fd1AccountNumber string `json:"fd1AccountNumber"`
  Fd1RoutingNumber string `json:"fd1RoutingNumber"`
  //
  // Fd2Result []Row `json:"fd2Result"`
}

func newManageAccountsFields(dir1, dir2, correlationId string) *manageAccountsFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "manageaccounts.txt")
  if obj != nil {
    var m manageAccountsFields
    err := json.Unmarshal(obj, &m)
    if err != nil {
      //Write error, but continue with default values.
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    } else {
      return &m
    }
  } else if err != nil {
    logger.LogError(fmt.Sprintf("%+v", err), correlationId)
  } else {
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "manageaccounts.txt"), correlationId)
  }
  return &manageAccountsFields {
    CurrentButton: "lhs-button1",
    Fd1BankName: "",
    Fd1AccountType: "checking",
    Fd1AccountName: "",
    Fd1AccountNumber: "",
    Fd1RoutingNumber: "",
  }
}

func getManageAccountsFields(userName string) *manageAccountsFields {
  return currentFields[userName].manageAccounts
}






func AddSessionDataPerUser(userName, correlationId string) {
  if _, ok := currentFields[userName]; !ok {
    fd := &fields{
      manageAccounts: newManageAccountsFields(mainDir, userName, correlationId),
      users: newUsersFields(mainDir, userName, correlationId),
    }
    currentFields[userName] = fd
  }
}

func DeleteSessionDataPerUser(userName string) {
  delete(currentFields, userName)
}

/***
Notes about synchronization:
(1) Even though each username has its own exclusive set of files, multiple users can share the same
    username. In this scenario, there is a possibility that two or more users may try to modify a
    file at the same time thereby corrupting the file. In order to prevent file corruption, the
    file is protected with a lock that allows a single writer (exclusive write) or multiple readers
    (share reads).
(2) It's always good practice to protect shared resources even if you're sure that there will be no
    conflict.
***/
func readFields(filePath string) ([]byte, error) {
  exists, err := osu.CheckFileExists(filePath)
  if exists {
    obj, err := osu.ReadAllShareLock1(filePath, os.O_RDONLY, 0o400)
    if err != nil {
      return nil, errors.New(fmt.Sprintf("Couldn't open file %s: ", filePath) + err.Error())
    }
    return obj, nil
  } else {
    return nil, err
  }
}
