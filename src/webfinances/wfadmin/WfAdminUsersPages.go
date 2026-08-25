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


  "strconv"
  "math"


)


type Row struct {  //Rows for the accounts.
  AccountName string
  AccountType string
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

type WfAdminUsersPages struct {}

func (u WfAdminUsersPages) AdminUsersPages(res http.ResponseWriter, req *http.Request) {
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
        } { "std-wo-nav-menu", "Register User - Admin", logger.DatetimeFormat(), fields.CurrentButton, "", "", "", "", "", "", "male",
            time.Now().Format("2006-01-02"), "false", "", "", "", "", "", "", "", "", "" }
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
          logger.LogError(err.Error(), correlationId)
          /***
          On error, time.Parse returns a zero time value (0001-01-01 00:00:00 +0000 UTC).
          ***/
          c.Birth_date = bank.TimePtr(newDate)
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
        "webfinances/templates/layout.html",
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
    } else if strings.EqualFold(fields.CurrentPage, "rhs-ui2") {
      fields.CurrentButton = "lhs-button2"


      /*
Using the last selected range (or defaulting to the first range on a fresh login) is an excellent usability pattern. In UX design, this is called Smart Defaults.By predicting what the user wants to see, you eliminate a mandatory extra click every time they visit the page, while still giving them full control to change the range using the slider whenever they want.
      */


      //Always parse the form up front so Go reads both URL query strings and POST bodies cleanly
      if err := req.ParseForm(); err != nil {
        // Handle error if necessary
      }

      // Check if the user explicitly dragged the slider (POST)
      rangeInput := req.PostFormValue("alphabet-range")
      if rangeInput != "" {
        fields.SelectedRange = rangeInput
      }
      // SMART DEFAULT: If the user just landed on the page (GET request)
      // and fields.SelectedRange is blank, automatically fall back to the first range ("1")
      if fields.SelectedRange == "" {
        fields.SelectedRange = "1"
      }
      var minLetter, maxLetter string
      switch fields.SelectedRange {
      case "1":
        minLetter, maxLetter = "A", "G"
      case "2":
        minLetter, maxLetter = "H", "N"
      case "3":
        minLetter, maxLetter = "O", "T"
      case "4":
        minLetter, maxLetter = "U", "Z"
      default:
        minLetter, maxLetter = "A", "G" // Fallback safety
      }
      logger.LogInfo(fmt.Sprintf("min: %s,  max: %s", minLetter, maxLetter), correlationId)

      var rows []Row
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


          rowId := req.PostFormValue("selected_id")

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




      //Extract the page from form body (POST) or URL query string (GET)
      pageStr := req.FormValue("page")
      currentPage, err := strconv.Atoi(pageStr)
      if err != nil || currentPage < 1 {
        currentPage = 1
      }
      //Pagination Math Calculations
      pageSize := 10  // Items per page
      totalItems := len(rows)
      totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
      if totalPages < 1 {
        totalPages = 1

      }
      //
      if currentPage > totalPages {
        currentPage = totalPages
      }
      //Calculate slicing boundaries
      offset := (currentPage - 1) * pageSize
      end := offset + pageSize
      if end > totalItems {
        end = totalItems
      }
      // Slice the data chunk safely
//      var paginatedItems []string
      var paginatedItems []Row
      if offset < totalItems {
//        paginatedItems = mockDatabase[offset:end]
        paginatedItems = rows[offset:end]
      }
      //

      if req.Method == http.MethodPost {
        //Go is designed to look at the request, realize the body hasn't been parsed yet, and automatically call ParseForm() for you under the hood.

        fields.SelectedRange = req.PostFormValue("alphabet-range")
        var minLetter, maxLetter string
        switch fields.SelectedRange {
        case "1":
          minLetter, maxLetter = "A", "G"
        case "2":
          minLetter, maxLetter = "H", "N"
        case "3":
          minLetter, maxLetter = "O", "T"
        case "4":
          minLetter, maxLetter = "U", "Z"
        default:
          minLetter, maxLetter = "A", "G" // Fallback safety
        }
        logger.LogInfo(fmt.Sprintf("min: %s,  max: %s", minLetter, maxLetter), correlationId)
        // 4. Pass these bounded letters into your SQL query execution
        //query := "SELECT id, username FROM users WHERE username >= ? AND username <= ?"



        // alpharange := req.PostFormValue("alphabet-range")
        // fields.SelectedRange = SelectedRange
      }
      newSessionToken, newSession := sessions.UpdateEntryInSessions(sessionToken)
      cookie := sessions.CreateCookie(newSessionToken)
      http.SetCookie(res, cookie)
      templatesNeeded := []string{
        "webfinances/templates/layout.html",
        "webfinances/templates/admin/users/users.html",
        "webfinances/templates/admin/users/unregister.html",
        "webfinances/templates/helpers/slider-alphabet-container.html",
        "webfinances/templates/helpers/scroll-container.html",
        "webfinances/templates/helpers/pagination-container.html",
        "webfinances/templates/title.html",
        "webfinances/templates/datetime.html",
        "webfinances/templates/footer.html",
      }

labels := []string{"A - G", "H - N", "O - T", "U - Z"}

      /***
      Do not read or write sensitive information from the disk; use the database exclusively.
      ***/
      renderer.Render(res, "layout", templatesNeeded, renderer.PageData{
        Data: struct{
          LayoutType string
          Header string
          Datetime string
          CurrentButton string
          CsrfToken string
          SelectedRange string
          ////
//            Items       []string
            Fd2Result       []Row
  CurrentPage int
  TotalPages  int
  PrevPage    int
  NextPage    int
  HasPrev     bool
  HasNext     bool
  RangeLabels   []string
  SliderMax     int

        } { "std-wo-nav-menu", "Unregister User - Admin", logger.DatetimeFormat(), fields.CurrentButton, newSession.CsrfToken,
            fields.SelectedRange,
paginatedItems, currentPage, totalPages, currentPage - 1, currentPage + 1, currentPage > 1, currentPage < totalPages,
  labels, len(labels)         },
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
