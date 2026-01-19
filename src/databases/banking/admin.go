package banking

//To fold all block comments:
//  Ctrl+K and Ctrl+/
//To unfold all block comments:
//  Ctrl+K and Ctrl+J

import (
  "context"
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "time"
)

/***
Notes:
In pgx, you can use named arguments for stored procedures by using the pgx.NamedArgs type, where
parameter names are prefixed with an @ symbol in the query string. The library will then
automatically rewrite the query to use positional parameters ($1, $2, etc.) before execution.
***/
const (
  //Use placeholder syntax (like $1, $2) to safely pass parameters to the function, preventing
  //SQL injection.
  SP_ADD_CUSTOMER = "CALL fin.add_customer($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12," +
    " $13, $14, $15, $16)"

  // SP_CUSTOMER_INFO = "CALL bs.sp_customer_info(@userName, @password, @customerType, @firstName, " +
  //  "@middleName, @lastName, @dateOfBirth, @taxIdentifier, @address1, @address2, @cityName, " +
  //  "@stateName, @countryName, @zipCode, @primaryEmail, @secondaryEmail, @primaryPhone, " +
  //  "@secondaryPhone)"
  // QR_GET_ALL_CUSTOMERS = "SELECT customer_id, customer_type, username, password_hash, " +
  //  "created_at, updated_at FROM bs.tbl_customer"
  // QR_GET_ALL_CUSTOMERS_INFO = "SELECT customer_info_id, customer_id, first_name, COALESCE(middle_name, ''), " +
  //  "last_name, date_of_birth, tax_identifier, address_1, COALESCE(address_2, ''), city_name, state_name, " +
  //  "country_name, COALESCE(zip_code, ''), COALESCE(primary_email, ''), COALESCE(secondary_email, ''), primary_phone, COALESCE(secondary_phone, ''), " +
  //  "created_at, updated_at FROM bs.tbl_customer_info"
    // fmt.Printf("%s - %s - %s - %s - %s - %s -- %s -- %s - %s - %s - %s - %s - %s - %s - %s- %s - %s - %s\n",
    //  c.UserName, c.Password, c.CustomerType, c.FirstName, c.MiddleName, c.LastName, c.DateOfBirth,
    //  c.TaxIdentifier, c.Address1, c.Address2, c.CityName, c.StateName, c.CountryName, c.ZipCode, c.PrimaryEmail,
    //  c.SecondaryEmail, c.PrimaryPhone, c.SecondaryPhone)
)

type AddCustomer struct {
  User_name, Password_hash, First_name, Last_name, Gender, Address1, City, State,
  Country, Email, Phone string
  Marketing bool
  //Nullable types.
  Address2, Middle_name, Zip_code *string
  Birth_date *time.Time
}

func DbAddCustomer(c *AddCustomer, ctx context.Context, correlationId string) error {
  db := GetBsInstance()
  //Use Exec when the stored procedure does not return a result set.
  _, err := db.bsPool.Exec(ctx, SP_ADD_CUSTOMER, c.User_name, c.Password_hash, c.First_name,
    c.Middle_name, c.Last_name, c.Marketing, c.Birth_date, c.Gender, c.Address1, c.Address2,
    c.City, c.State, c.Country, c.Zip_code, c.Email, c.Phone)
  if err != nil {
    logger.LogError(fmt.Sprintf("SP fin.add_customer: %v", err), correlationId)
    return err
  } else {
    logger.LogInfo(fmt.Sprintf("SP fin.add_customer succeeded. Username: %s", c.User_name),
      correlationId)
    return nil
  }
}


/***
func dbSelectAllCustomer() {
  p := bank.GetBsInstance()
  customers := p.GetAllCustomer(context.Background(), falseCorrelationId)
  if customers == nil {
    return
  }
  fmt.Println("Results from stored function:")
  for _, c := range customers {
    fmt.Printf("%d - %s - %s- %s - %s - %s\n", c.CustomerId, c.CustomerType, c.UserName, c.Password, c.CreatedAt, c.UpdatedAt)
  }
}

func dbSelectAllCustomerInfo() {
  p := bank.GetBsInstance()
  customers := p.GetAllCustomerInfo(context.Background(), falseCorrelationId)
  if customers == nil {
    return
  }
  fmt.Println("Results from stored function:")
  for _, c := range customers {
    fmt.Printf("%d - %d - %s - %s - %s - %s - %s- %s - %s - %s - %s - %s - %s - %s - %s - %s - %s - %s - %s\n", c.CustomerInfoId, c.CustomerId, c.FirstName, c.MiddleName, c.LastName,
    c.DateOfBirth.Format("2006-01-02"), c.TaxIdentifier, c.Address_1, c.Address_2, c.CityName, c.StateName, c.CountryName, c.ZipCode, c.PrimaryEmail, c.SecondaryEmail, c.PrimaryPhone, c.SecondaryPhone, c.CreatedAt, c.UpdatedAt)
    //c.CreatedAt.Format("2006-01-02"))
  }
}
***/
