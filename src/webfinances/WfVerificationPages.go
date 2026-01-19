package webfinances

import (
  bank "finance/databases/banking" //Importing a package and assigning it a local alias.
  "context"
  "fmt"
  "github.com/juan-carlos-trimino/go-middlewares"
  "github.com/juan-carlos-trimino/gplogger"
  "html/template"
  "net/http"
  "strings"
  "time"
)

type WfVerificationPages struct {}

func (s WfVerificationPages) RegistrationPage(res http.ResponseWriter, req *http.Request) {
  ctxKey := middlewares.MwContextKey{}
  correlationId, _ := ctxKey.GetCorrelationId(req.Context())
  startTime, _ := ctxKey.GetStartTime(req.Context())
  logger.LogInfo(fmt.Sprintf("Created correlationId at %s.",
    startTime.UTC().Format(time.RFC3339Nano)), correlationId)
  logger.LogInfo("Entering VericationPage.", correlationId)
  c := bank.AddCustomer {
    User_name: req.PostFormValue("uname"),
    Password_hash: req.PostFormValue("pwd"),
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
  //Go's time formatting uses a reference date and time: Mon Jan 2 15:04:05 MST 2006. Each
  //component of this reference time (e.g., 02 for the day, 01 for the month, 2006 for the year) is
  //used as a placeholder in the layout string to match the input format; e.g., "dd/mm/yyyy" is
  //"02/01/2006".
  if newDate, err := time.Parse("2006-01-02", originalDate); err != nil {
    fmt.Println("Error parsing date: ", err)
  } else {
    c.Birth_date = bank.TimePtr(newDate)
  }
  ok := bank.DbAddCustomer(&c, context.Background(), correlationId)
  if ok == nil {
    tmpl.ExecuteTemplate(res, "login_page", nil)
  } else {
    t := template.Must(template.ParseFiles("webfinances/templates/register.html"))
    t.ExecuteTemplate(res, "register_page", struct {
      Username string
      Password string
      Fname string
      Mname string
      Lname string
      Marketing string
      Bdate string
      Gender string
      Address1 string
      Address2 string
      City string
      State string
      Country string
      Zip_Code string
      Email string
      Phone string
      ErrMsg string
    } { c.User_name, c.Password_hash, c.First_name, middle_name, c.Last_name, marketing,
        originalDate, c.Gender, c.Address1, address2, c.City, c.State, c.Country, zip_code,
        c.Email, c.Phone, ok.Error(),
    })
  }
}
