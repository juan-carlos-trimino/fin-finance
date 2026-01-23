package banking

//To fold all block comments:
//  Ctrl+K and Ctrl+/
//To unfold all block comments:
//  Ctrl+K and Ctrl+J

import (
  "context"
  "fmt"
  "time"
  "github.com/juan-carlos-trimino/gplogger"
  "golang.org/x/crypto/bcrypt"
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
  //Pass null for OUT parameter in the call.
  SP_GET_PASSWORD_HASH = "CALL fin.get_password_hash($1, null)"
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
  c.Password_hash = DbHashAndSaltPassword(c.Password_hash, correlationId)
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

// This function should be used during user registration.
func DbHashAndSaltPassword(password, correlationId string) string {
  //Use GenerateFromPassword to hash & salt the password.
  //The cost can be any value you want, but DefaultCost is a good starting point.
  hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    logger.LogError(fmt.Sprintf("Error on DbHashAndSaltPassword: %v", err), correlationId)
    return ""
  }
  //The resulting hash string includes the salt and all necessary parameters, which should be
  //stored in the database. GenerateFromPassword returns a byte slice, so convert it to a string
  //for storage.
  return string(hash)
}

//During login, retrieve the stored hash from the database and use the CompareHashAndPassword
//function to verify the user-provided password. This function extracts the salt from the stored
//hash, applies it to the input password, and compares the resulting hashes securely, mitigating
//timing attacks.
func DbAuthenticateUser(ctx context.Context, userName, candidatePassword,
  correlationId string) bool {
  db := GetBsInstance()
  var password_hash string
  err := db.bsPool.QueryRow(ctx, SP_GET_PASSWORD_HASH, userName).Scan(&password_hash)
  if err != nil {
    logger.LogError(fmt.Sprintf("Error on DbAuthenticateUser: %v", err), correlationId)
    return false
  }
  //CompareHashAndPassword compares a bcrypt hashed password with its possible plaintext
  //equivalent. It returns nil on success, or an error on failure.
  err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(candidatePassword))
  if err != nil {
    //Log the error but do not reveal specific details to the user for security reasons.
    logger.LogError(fmt.Sprintf("Error on DbAuthenticateUser: %v", err), correlationId)
    return false
  }
  return true
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
