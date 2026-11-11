package wfadmin

import (
  "context"
  "encoding/base64"
  "encoding/json"
  bank "finance/databases/banking" //Importing a package and assigning it a local alias.
  "finance/renderer"
  "fmt"
  // "math"
  "net/http"
  "strconv"
  "strings"
  "time"

  "github.com/juan-carlos-trimino/go-logger"
  "github.com/juan-carlos-trimino/go-middlewares"
  sess "github.com/juan-carlos-trimino/go-sessions"
)

const (
  defaultItemsPerPage = 10
  maxItemsPerPage int = 100  //Hard maximum limit enforced on the server.
)

type usersFields struct{
  Limit int `json:"limit"`
  CurrentButton string `json:"currentButton"`
  CurrentPage string  `json:"currentPage"`
  SelectedRange string `json:"selectedRange"`
}

func defaultUsersFields() *usersFields {
  return &usersFields{
    Limit: defaultItemsPerPage,  //Establish the default baseline.
    CurrentButton: "lhs-button1",
    CurrentPage: "rhs-ui1",
    SelectedRange: "1",
  }
}

type PageInfo struct{
  NextCursor string
  HasNext bool
  Limit int  //Items/Page
}

type WfAdminUsersPages struct{}


/***jct
Cloud-Specific Architectural Considerations

While the pattern is ideal, you need to account for how a cloud environment behaves:

1. Always Order by an Indexed Partition Strategy

In your multi-tenant setup, your queries likely include a tenant filter or user isolation boundary (e.g., WHERE user_id = $1). To ensure performance remains near-instant, you must create a composite B-Tree index that covers your partition keys and your sorting cursor keys together.

2. Keep the Cursor Decoupled from Internal Data Types

Cloud environments are dynamic. If you change your underlying database schema later, you don't want your UI or external mobile apps to break because they were storing raw database values.
* Best Practice: Always stick to the opaque Base64-encoded string pattern. The UI just treats it as a text string, allowing you to rewrite how the backend parses or modifies the cursor internals later without updating your client applications.

3. Handle Empty and Edge Payload States Gracefully

In a distributed cloud setup, components need explicit contracts to prevent infinite loader loops on the front end. Ensure your API payload always handles edge cases cleanly:
* First Page: The client sends no cursor (/api/feed). The backend returns the first batch and a next_cursor.
* Middle Pages: The client passes the string (/api/feed?cursor=ZXleU0M4TW...).
* Last Page: The backend finds no further rows during the look-ahead check. It returns data but sets has_next: false and next_cursor: "". The UI immediately unbinds scroll event listeners.
***/



func (u WfAdminUsersPages) AdminUsersPages(res http.ResponseWriter, req *http.Request) {
  ck := middlewares.MwContextKey{}
  correlationId, _ := ck.GetCorrelationId(req.Context())
  logger.LogInfo("Entering wfadmin.AdminUsersPage.", correlationId)
  sessInfo, _ := ck.GetSessionInfo(req.Context())
  redisKey := "user:data:" + sessInfo.UserName
  partitionName := "admin_users"
  var fields *usersFields = nil
  var jsonBytes []byte
  //Check Redis Cache First.
  jsonBytes, err := sess.HGetRedis(req.Context(), redisKey, partitionName)
  if err == nil {
    logger.LogInfo(fmt.Sprintf("Redis found key(%s)/partition(%s)...", redisKey, partitionName), correlationId)
    //Unmarshal unmarshals the JSON into the value pointed at by v. If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
    err = json.Unmarshal(jsonBytes, &fields)
    if err != nil {
      logger.LogInfo(fmt.Sprintf("Error (Redis) unmarshalling key(%s)/partition(%s): %v. Using default values.", redisKey, partitionName, err), correlationId)
      fields = defaultUsersFields()
    }
  } else {
    //Cache Miss: Fall back to PostgreSQL database.
    logger.LogInfo(fmt.Sprintf("Redis does not have key(%s)/partition(%s): %v. Fetching from PostgreSQL...", redisKey, partitionName, err), correlationId)
    jsonBytes = bank.DbFetchUserData(req.Context(), sessInfo.UserName, partitionName, correlationId)
    if len(jsonBytes) == 0 {  //On error (nil) or not found (len() == 0), use default values.
      logger.LogInfo(fmt.Sprintf("Database did not find user(%s)/partition(%s). Using default values.", sessInfo.UserName, partitionName), correlationId)
      fields = defaultUsersFields()
    } else {  //Use the values from the db.
      logger.LogInfo(fmt.Sprintf("Database found user(%s)/partition(%s)...", sessInfo.UserName, partitionName), correlationId)
      //Unmarshal unmarshals the JSON into the value pointed at by v. If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
      err = json.Unmarshal(jsonBytes, &fields)
      if err != nil {
        logger.LogInfo(fmt.Sprintf("Error (db) unmarshalling user(%s)/partition(%s): %v. Using default values.", sessInfo.UserName, partitionName, err), correlationId)
        fields = defaultUsersFields()
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
  if ui := req.FormValue("db"); ui != "" {  //Values from form and URL.
    fields.CurrentPage = ui
  }
  //Dynamic variables determined by the routing route condition.
  var partialTemplate string
  var templateData interface{}
  var templatesToAdd []string
  var layoutName string = "layout"
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
    }{ "std-wo-nav-menu", "Register User - Admin", logger.DatetimeFormat(), fields.CurrentButton, sessInfo.CSRFToken, "", "", "",
       "", "", "male", time.Now().Format("2006-01-02"), "false", "", "", "", "", "", "", "", "", "" }
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
    /***
    When you provide only the length, the capacity automatically matches it. All elements initialize to their zero values.
    Setting the length to 0 ensures that append starts inserting at index 0.
    ***/
    templatesToAdd = make([]string, 0, 8)
    templatesToAdd = append(templatesToAdd,
      "webfinances/templates/helpers/slider-alphabet-container.html",
      "webfinances/templates/helpers/scroll-container.html",
      "webfinances/templates/helpers/pagination-container-load-more.html",
    )
    var lowerBound string
    var upperBound *string
    var users []bank.Users
    /***
    req.FormValue() retrieves values from BOTH GET and POST requests.
    It is designed to be a universal fallback method. When you call req.FormValue(key), Go searches for the key in two distinct
    places, checking them in a specific order of priority:
    1. The Request Body (POST/PUT/PATCH): It checks the body payload first (data sent via application/x-www-form-urlencoded).
       During a POST request: It automatically parses the body and returns the value from the form field. (If that field isn't
       in the body, it checks the URL query string as a backup).
    2. The URL Query String (GET): It checks the URL parameters second (data sent via ?key=value).
       During a GET request: Because there is no request body, req.FormValue("key") safely checks the URL and returns the query
       parameter value.
    ***/
    if limitStr := req.FormValue("limit"); limitStr != "" {
      if pageLimit, err := strconv.Atoi(limitStr); err != nil {
        logger.LogError(fmt.Sprintf("Limit is invalid: %v", err), correlationId)
      } else if pageLimit > 0 && !(pageLimit > maxItemsPerPage) {
        //Strictly enforce the server-side lower/upper boundary.
        fields.Limit = pageLimit
      } else {
        logger.LogInfo(fmt.Sprintf("Limit (%d) > maximum items per page (%d)", pageLimit, maxItemsPerPage), correlationId)
      }
    }
    /***
    To keep the UI responsive, use the Look-Ahead pattern. When a user requests a limit of 20 records, the database query should
    request 21 items.
    * If 21 items return: A next page exists. You set has_next: true, strip off the 21st item so it isn't displayed, and extract
      your cursor values from the 20th item.
    * If 20 or fewer items return: You are on the final page. Set has_next: false and hide your "Load More" button or stop
      listening to scroll events.
    ***/
    fetchLimit := fields.Limit + 1
    if req.Method == http.MethodPost {
      /***
      Using the last selected range (or defaulting to the first range on a fresh login) is an excellent usability pattern. In UX
      design, this is called Smart Defaults. By predicting what the user wants to see, you eliminate a mandatory extra click every
      time they visit the page, while still giving them full control to change the range using the slider whenever they want.

      Go is designed to look at the request, realize the body hasn't been parsed yet, and automatically call ParseForm() for
      you under the hood.
      ***/
      actionType := req.PostFormValue("action-button")
      targetIdStr := req.PostFormValue("id")
      rangeStr := req.PostFormValue("alphabet-range")  //Check if the user explicitly dragged the slider.
      /***
      -------------------------------------------------
      "Delete User" Button Clicked (Standard Form POST)
      -------------------------------------------------
      When a standard form POST submission happens, the browser completely throws away the current webpage layout and redraws a
      brand-new page from scratch based on the HTML response sent back by Go. (When you click "Delete User", the script explicitly
      calls mainForm.submit(). This bypasses background manipulation, executes a full browser cycle, and builds a fresh,
      standalone DOM tree.) For example, "delete" generates a standard form POST, and the browser wipes the old DOM completely,
      which requires a full page redraw. Consequently, the Go execution pathway must query the database from scratch and render
      the full layout.

      When a JavaScript fetch submission occurs, the JavaScript intercepts the event via event.preventDefault() and dispatches
      a background task. Go returns only the raw text strings for the <tr> tags, and the script appends them. For example,
      "load-more" generates a JavaScript fetch, and the browser keeps the page static and appends rows using Asynchronous Injection.
      Hence, the Go execution pathway returns only raw <tr> row strings.
      ***/
      switch actionType {
      case "delete":
        if targetIdStr != "" {
          targetId, err := strconv.Atoi(targetIdStr)
          if err != nil {
            logger.LogWarning(fmt.Sprintf("Invalid id, resetting to page 1: %v", err), correlationId)
          } else {
            ok := bank.DbDeleteUser(req.Context(), targetId, correlationId)
            if ok {
              logger.LogInfo("User deleted, resetting to page 1.", correlationId)
            } else {
              logger.LogWarning("Database error, resetting to page 1.", correlationId)
            }
          }
        } else {
          logger.LogWarning("Id is empty, resetting to page 1.", correlationId)
        }
        //REFRESH DATA ENGINE: Pull fresh records after the user has been removed.
        lowerBound, upperBound = getRangeBounds(fields.SelectedRange)
        //Pass 0 and "" to pull from the very beginning of the current range.
        users = bank.DbFetchUsers(req.Context(), lowerBound, upperBound, fetchLimit, 0, "", correlationId)
      case "load-more":
        //Instead of executing the outer layout, execute the inner block definition.
        layoutName = "user-rows-fragment"
        /***
        ------------------------------------------------------------
        Asynchronous "Load More" Clicked (via Fetch background task)
        ------------------------------------------------------------
        Confirmed: The "Load More" button was clicked.
        Execute data retrieval using the cursorParam and stream ONLY the rows loop fragment.
        Decode cursor attributes, if present.
        ***/
        cursorParam := req.PostFormValue("cursor")
        var lastUserName string
        var lastId int
        var err error
        if cursorParam == "" {
          logger.LogError("Cursor parameter is empty, resetting to page 1.", correlationId)
          //Reset to zero values so DbFetchUsers fetches from the start.
          lastUserName = ""
          lastId = 0
        } else {
          lastUserName, lastId, err = decodeCursor(cursorParam)
          if err != nil {
            logger.LogError(fmt.Sprintf("Invalid cursor parameter, resetting to page 1: %v", err), correlationId)
            //Reset to zero values so DbFetchUsers fetches from the start.
            lastUserName = ""
            lastId = 0
          }
        }
        lowerBound, upperBound = getRangeBounds(fields.SelectedRange)
        users = bank.DbFetchUsers(req.Context(), "", upperBound, fetchLimit, lastId, lastUserName, correlationId)
      default:
        // ------------------------------------------------------
        // Slider Alphabet Button Clicked (Standard Form POST)
        // ------------------------------------------------------
        if rangeStr != "" {
          fields.SelectedRange = rangeStr
        }
        lowerBound, upperBound = getRangeBounds(fields.SelectedRange)
        users = bank.DbFetchUsers(req.Context(), lowerBound, upperBound, fetchLimit, 0, "", correlationId)
      }
    } else {  //GET
      rangeStr := req.URL.Query().Get("alphabet-range")
      if rangeStr != "" {
        fields.SelectedRange = rangeStr
      }
      lowerBound, upperBound = getRangeBounds(fields.SelectedRange)
      users = bank.DbFetchUsers(req.Context(), lowerBound, upperBound, fetchLimit, 0, "", correlationId)
    }
    //Formulate look-ahead pagination state.
    hasNext := false
    nextCursor := ""
    if len(users) > fields.Limit {
      hasNext = true
      users = users[:fields.Limit]  //Trim look-ahead item.
      lastVisibleUser := users[fields.Limit - 1]
      nextCursor = encodeCursor(lastVisibleUser.UserName, lastVisibleUser.Id)  //////////////////////////////////////////jct
    }
    var pageInfo = PageInfo{
      NextCursor: nextCursor,
      HasNext: hasNext,
      Limit: fields.Limit,
    }
    labels := []string{ "0 - 9", "$..@", "A - G", "H - N", "O - T", "U - Z" }
    templateData = struct{
      LayoutType string
      Header string
      Datetime string
      CurrentButton string
      CsrfToken string
      SelectedRange string
      Fd2Result []bank.Users
      PageInfo PageInfo
      RangeLabels []string
      SliderMax int
    }{
      "std-wo-nav-menu",
      "Unregister User - Admin",
      logger.DatetimeFormat(),
      fields.CurrentButton,
      sessInfo.CSRFToken,
      fields.SelectedRange,
      users,
      pageInfo,
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
  renderer.Render(res, layoutName, templatesNeeded, renderer.PageData{Data: templateData})
  data, err := json.Marshal(fields) //Preserve current choices.
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

func getRangeBounds(selectedRange string) (string, *string) {
  //Helper to return a pointer to a string.
  strPtr := func(s string) *string {
    return &s
  }
  switch selectedRange {
  case "1":
    //Numbers: From "0" up to (but excluding) "a".
    return "0", strPtr("a")
  case "2":  //Matches label "$..@"
    //Under ICU, digits (0-9) immediately follow symbols and punctuation.
    //Reading from the lowest possible string "" up to "0" captures all special characters.
    return "", strPtr("0")
  case "3":  //Matches label "A - G"
    //A through G: Includes letters like À, Ç, É.
    return "a", strPtr("h")
  case "4":  //Matches label "H - N"
    //H through N: Includes letters like Ñ.
    return "h", strPtr("o")
  case "5":  //Matches label "O - T"
    //O through T: Includes letters like Ö.
    return "o", strPtr("u")
  case "6":  //Matches label "U - Z"
    //Open-ended ceiling! No upper string boundary constraint required.
    return "u", nil
  default:
    return "a", strPtr("h")
  }
}

func encodeCursor(userName string, id int) string {
  //Build an opaque string for the UI from the last item properties.
  raw := fmt.Sprintf("%s|%d", userName, id)
  /***
  base64.RawURLEncoding: Using standard URL-safe base64 encoding removes padding characters (=), which prevents errors or
  character stripping when passed inside standard net/http query strings (?cursor=XYZ).
  ***/
  return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

/***
In a multi-column index setup (like sorting by a timestamp where multiple records could share the exact same time), the cursor
must include both the sorting column and a unique tie-breaker column (like the Id) to maintain strict pagination sequence.
***/
func decodeCursor(cursorStr string) (string, int, error) {
  //Parse the string back into database query targets.
  data, err := base64.RawURLEncoding.DecodeString(cursorStr)
  if err != nil {
    return "", 0, err
  }
  parts := strings.Split(string(data), "|")
  if len(parts) != 2 {
    return "", 0, fmt.Errorf("Invalid cursor format.")
  }
  id, err := strconv.Atoi(parts[1])
  if err != nil {
    return "", 0, err
  }
  return parts[0], id, nil  //Return (username, id)
}
