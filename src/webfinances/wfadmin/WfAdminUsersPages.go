package wfadmin

import (
  bank "finance/databases/banking"  //Importing a package and assigning it a local alias.
  "context"
  "encoding/json"
  "finance/renderer"
  "fmt"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/go-logger"
  "github.com/juan-carlos-trimino/go-os"
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
  ck := middlewares.MwContextKey{}
  correlationId, _ := ck.GetCorrelationId(req.Context())
  logger.LogInfo("Entering wfadmin.AdminUsersPage.", correlationId)
  sessInfo, _ := ck.GetSessionInfo(req.Context())
  fields := getUsersFields(sessInfo.UserName)
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
  //Dynamic variables determined by the routing route condition.
  var partialTemplate string
  var templateData interface{}
  var templatesToAdd []string
  //Core application state splitting via explicit page mapping.
  switch strings.ToLower(fields.CurrentPage) {
  case "rhs-ui1":
    fields.CurrentButton = "lhs-button1"
    td := struct{  //Default values.
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
    }{ "std-wo-nav-menu", "Register User - Admin", logger.DatetimeFormat(), fields.CurrentButton, "", "", "", "", "", "", "male",
          time.Now().Format("2006-01-02"), "false", "", "", "", "", "", "", "", "", "" }
    if req.Method == http.MethodPost {
      c := bank.Customer{}
      err := u.processUi1Form(req, &c)
      if err != nil {
        logger.LogError(err.Error(), correlationId)
      } else {
        err = bank.DbAddCustomer(&c, context.Background(), correlationId)  //Db logs its own error.
      }
      //
      if err != nil {
        td.Username = c.User_name
        td.Password = c.Password
        td.Fname = c.First_name
        td.Mname = bank.PtrString(c.Middle_name)
        td.Lname = c.Last_name
        td.Gender = c.Gender
        td.Bdate = c.Birth_date.Format("2006-01-02")
        if c.Marketing {
          td.Marketing = "true"
        } else {
          td.Marketing = "false"
        }
        td.Address1 = c.Address1
        td.Address2 = bank.PtrString(c.Address2)
        td.City = c.City
        td.State = c.State
        td.Country = c.Country
        td.Zip_Code = bank.PtrString(c.Zip_code)
        td.Email = c.Email
        td.Phone = c.Phone
        td.ErrMsg = fmt.Sprintf("%v", err)
      }
    }
    partialTemplate = "register.html"
    templateData = td
  case "rhs-ui2":
    fields.CurrentButton = "lhs-button2"
    partialTemplate = "unregister.html"
    //When you provide only the length, the capacity automatically matches it. All elements initialize to their zero values.
    templatesToAdd = make([]string, 0, 8)
    //Setting the length to 0 ensures that append starts inserting at index 0.
    templatesToAdd = append(templatesToAdd,
      "webfinances/templates/helpers/slider-alphabet-container.html",
      "webfinances/templates/helpers/scroll-container.html",
      "webfinances/templates/helpers/pagination-container.html",
    )

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

    labels := []string{"A - G", "H - N", "O - T", "U - Z"}

    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      CurrentButton string
      CsrfToken string
      SelectedRange string
      Fd2Result []Row
      CurrentPage int
      TotalPages int
      PrevPage int
      NextPage int
      HasPrev bool
      HasNext bool
      RangeLabels []string
      SliderMax int
    }{
      "std-wo-nav-menu",
      "Unregister User - Admin",
      logger.DatetimeFormat(),
      fields.CurrentButton,
      sessInfo.CSRFToken,
      fields.SelectedRange,
      paginatedItems,
      currentPage,
      totalPages,
      currentPage - 1,
      currentPage + 1,
      currentPage > 1,
      currentPage < totalPages,
      labels,
      len(labels),
    }
  default:
    errString := fmt.Sprintf("Unsupported page: %s", fields.CurrentPage)
    logger.LogError(errString, correlationId)
    panic(errString)
  }
  //Unified execution of templates.
  templatesNeeded := []string{
    "webfinances/templates/layout.html",
    "webfinances/templates/admin/users/users.html",
    "webfinances/templates/admin/users/" + partialTemplate,
    "webfinances/templates/title.html",
    "webfinances/templates/datetime.html",
    "webfinances/templates/footer.html",
  }
  //Add if there are additional elements.
  if templatesToAdd != nil {
    /***
    In Go, the three dots (...) are called the unpack operator (or variadic operator). They are needed here because of how Go's
    built-in append function is designed. It does not accept a slice as its second argument; it expects a list of individual elements.
    ***/
    templatesNeeded = append(templatesNeeded, templatesToAdd...)
  }
  //Do not read or write sensitive information from the disk; use the database exclusively.
  renderer.Render(res, "layout", templatesNeeded, renderer.PageData{Data: templateData})
  data, err := json.Marshal(fields) //Preserve current choices.
  if err != nil {
    //Don't crash the server (panic), but log it clearly so you can debug the serialization.
    logger.LogError(fmt.Sprintf("Failed to marshal fields to JSON for user %s: %+v", sessInfo.UserName, err), correlationId)
  } else {
    go func(userData []byte, uName, cId string) {
      filePath := fmt.Sprintf("%s/%s/users.txt", mainDir, uName)
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
}

//Extraction helper for UI-1 calculations.
func (u WfAdminUsersPages) processUi1Form(req *http.Request, c *bank.Customer) error {
  c.User_name = req.PostFormValue("uname")
  c.Password = req.PostFormValue("pwd")
  c.First_name = req.PostFormValue("fname")
  c.Last_name = req.PostFormValue("lname")
  c.Gender = req.PostFormValue("gender")
  c.Address1 = req.PostFormValue("address1")
  c.City = req.PostFormValue("city")
  c.State = req.PostFormValue("state")
  c.Country = req.PostFormValue("country")
  c.Email = req.PostFormValue("email")
  c.Phone = req.PostFormValue("phone")
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
  //On error, time.Parse returns a zero time value (0001-01-01 00:00:00 +0000 UTC).
  c.Birth_date = bank.TimePtr(newDate)
  return err
}
