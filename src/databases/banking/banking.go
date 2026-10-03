package banking

//To fold all block comments:
//  Ctrl+K and Ctrl+/
//To unfold all block comments:
//  Ctrl+K and Ctrl+J

import (
  "regexp"
  "strings"


  // "fmt"
  // "github.com/google/uuid"
  // "github.com/jackc/pgx/v5"
  // "github.com/juan-carlos-trimino/go-logger"
  // "time"
)

/***
Notes:
In pgx, you can use named arguments for stored procedures by using the pgx.NamedArgs type, where
parameter names are prefixed with an @ symbol in the query string. The library will then
automatically rewrite the query to use positional parameters ($1, $2, etc.) before execution.
***/
const (
  // SP_CUSTOMER_INFO = "CALL bs.sp_customer_info(@userName, @password, @customerType, @firstName, " +
  //  "@middleName, @lastName, @dateOfBirth, @taxIdentifier, @address1, @address2, @cityName, " +
  //  "@stateName, @countryName, @zipCode, @primaryEmail, @secondaryEmail, @primaryPhone, " +
  //  "@secondaryPhone)"

  // SP_CUSTOMER_INFO = "CALL bs.sp_customer_info(@userName, @password, @firstName, " +
  //  "@middleName, @lastName, @dateOfBirth, @taxIdentifier, @address1, @address2, @cityName, " +
  //  "@stateName, @countryName, @zipCode, @primaryEmail, @secondaryEmail, @primaryPhone, " +
  //  "@secondaryPhone)"
  // QR_GET_ALL_CUSTOMERS = "SELECT customer_id, customer_type, username, password_hash, " +
  //  "created_at, updated_at FROM bs.tbl_customer"
  QR_GET_ALL_CUSTOMERS_INFO = "SELECT customer_info_id, customer_id, first_name, COALESCE(middle_name, ''), " +
   "last_name, date_of_birth, tax_identifier, address_1, COALESCE(address_2, ''), city_name, state_name, " +
   "country_name, COALESCE(zip_code, ''), COALESCE(primary_email, ''), COALESCE(secondary_email, ''), primary_phone, COALESCE(secondary_phone, ''), " +
   "created_at, updated_at FROM bs.tbl_customer_info"
    // fmt.Printf("%s - %s - %s - %s - %s - %s -- %s -- %s - %s - %s - %s - %s - %s - %s - %s- %s - %s - %s\n",
    //  c.UserName, c.Password, c.CustomerType, c.FirstName, c.MiddleName, c.LastName, c.DateOfBirth,
    //  c.TaxIdentifier, c.Address1, c.Address2, c.CityName, c.StateName, c.CountryName, c.ZipCode, c.PrimaryEmail,
    //  c.SecondaryEmail, c.PrimaryPhone, c.SecondaryPhone)
)

// type Customer1 struct {
// 	Id uuid.UUID // Maps to a PostgreSQL UUID column
//   // Id UUID //`json: "customer_id"`
//   Username string //`json: "username"`
//   Password_hash string //`json: "password_hash"`
//   First_name string
//   Middle_name string
//   Last_name string
//   Birth_date time.Time
// 	Gender byte
//   Address_1 string
//   Address_2 string
//   City_name string
//   State_name string
//   Country_name string
//   Zip_code string
//   Primary_email string
//   Secondary_email string
//   Primary_phone string
//   Secondary_phone string
//   Created_at time.Time
//   Updated_at time.Time
// }


/*
func (bs *banking) SaveCustomer(ctx context.Context, userName, password, customerType,
  firstName, middleName, lastName *string, dateOfBirth *time.Time, taxIdentifier, address1,
  address2, cityName, stateName, countryName, zipCode, primaryEmail, secondaryEmail,
  primaryPhone, secondaryPhone *string, correlationId string) bool {
  args := pgx.NamedArgs{
    "userName": userName,
    "password": password,
    "customerType": customerType,
    "firstName": firstName,
    "middleName": middleName,
    "lastName": lastName,
    "dateOfBirth": dateOfBirth,
    "taxIdentifier": taxIdentifier,
    "address1": address1,
    "address2": address2,
    "cityName": cityName,
    "stateName": stateName,
    "countryName": countryName,
    "zipCode": zipCode,
    "primaryEmail": primaryEmail,
    "secondaryEmail": secondaryEmail,
    "primaryPhone": primaryPhone,
    "secondaryPhone": secondaryPhone,
  }
  //Use Exec when the stored procedure does not return a result set.
  _, err := bs.bsPool.Exec(ctx, SP_CUSTOMER_INFO, args)
  if err != nil {
    logger.LogError(fmt.Sprintf("Stored procedure bs.sp_customer_info failed: %v", err),
     correlationId)
    return false
  } else {
    logger.LogInfo("Stored procedure bs.sp_customer_info was successful (no return value).",
     correlationId)
    return true
  }
}
*/




/**
func (bs *bankingSystem) GetAllCustomer(ctx context.Context, correlationId string) []Customer {
  rows, err := bs.bsPool.Query(ctx, QR_GET_ALL_CUSTOMERS)
  if err != nil {
    logger.LogError(fmt.Sprintf("Function bs.fn_get_all_customers failed: %v", err), correlationId)
    return nil
  }
  defer rows.Close()
  var customers []Customer
  for rows.Next() {
    var c Customer
    err = rows.Scan(&c.CustomerId, &c.CustomerType, &c.UserName, &c.Password, &c.CreatedAt, &c.UpdatedAt)
    // err = rows.Scan(&c.UserName, &c.Password, &c.CustomerType, &c.FirstName, &c.MiddleName, &c.LastName,
    //   &c.DateOfBirth, &c.TaxIdentifier, &c.Address1, &c.Address2, &c.CityName, &c.StateName, &c.CountryName, &c.ZipCode, &c.PrimaryEmail,
    //   &c.SecondaryEmail, &c.PrimaryPhone, &c.SecondaryPhone)
    if err != nil {
      logger.LogError(fmt.Sprintf("Error scanning row: %v", err), correlationId)
      return nil
    }
    customers = append(customers, c)
  }
  //
  if err = rows.Err(); err != nil {
    logger.LogError(fmt.Sprintf("Error after iterating rows: %v", err), correlationId)
  }
  return customers
}





func (bs *bankingSystem) GetAllCustomerInfo(ctx context.Context, correlationId string) []CustomerInfo {
  rows, err := bs.bsPool.Query(ctx, QR_GET_ALL_CUSTOMERS_INFO)
  if err != nil {
    logger.LogError(fmt.Sprintf("Function bs.fn_get_all_customers failed: %v", err), correlationId)
    return nil
  }
  defer rows.Close()
  var customers []CustomerInfo
  for rows.Next() {
    var c CustomerInfo
    err = rows.Scan(&c.CustomerInfoId, &c.CustomerId, &c.FirstName, &c.MiddleName, &c.LastName, &c.DateOfBirth, &c.TaxIdentifier, &c.Address_1, &c.Address_2, &c.CityName, &c.StateName, &c.CountryName, &c.ZipCode, &c.PrimaryEmail, &c.SecondaryEmail, &c.PrimaryPhone, &c.SecondaryPhone, &c.CreatedAt, &c.UpdatedAt)

    if err != nil {
      logger.LogError(fmt.Sprintf("Error scanning row: %v", err), correlationId)
      return nil
    }
    customers = append(customers, c)
  }
  //
  if err = rows.Err(); err != nil {
    logger.LogError(fmt.Sprintf("Error after iterating rows: %v", err), correlationId)
  }
  return customers
}
**/




/***
SanitizeDbName converts a raw username into a valid, safe Postgres database name.

Because we are building a multi-tenant architecture using a database-per-tenant pattern, we would like to use the username as the
database name, but Postgres has strict rules for database identifiers. Hence, the username cannot be used directly as the
database name. For example, if a user registers with characters that violate these rules, the CREATE DATABASE statement will
crash; furthermore,
* Length Limit: By default, Postgres limits database name identifiers to 63 characters (NAMEDATALEN - 1). Any characters past this limit are silently truncated.
* Character Set Restrictions: Database names must start with a lowercase letter or an underscore. They can only contain lowercase letters, numbers, and underscores. They cannot contain spaces, dashes (-), upper-case letters, or special characters (like symbols or punctuation).
* Reserved Words: A user name cannot conflict with Postgres reserved keywords (like SELECT, USER, or DATABASE).

To make the design reliable, a sanitization function is required. This function converts usernames to lowercase, replaces forbidden characters with underscores, ensures they don't start with a number, and truncates them safely to 63 characters.
***/
func SanitizeDbName(userName string) string {
  //Force lowercase.
  name := strings.ToLower(userName)
  //Replace all non-alphanumeric characters (like spaces or dashes) with underscores.
  reg := regexp.MustCompile(`[^a-z0-9_]`)
  name = reg.ReplaceAllString(name, "_")
  //Trim any trailing or leading underscores that resulted from cleaning.
  name = strings.Trim(name, "_")
  //Force a clean, uniform prefix.
  //This automatically handles usernames starting with numbers or reserved keywords.
  if !strings.HasPrefix(name, "usr_") {
    name = "usr_" + name
  }
  //Truncate to the PostgreSQL limit of 63 characters.
  if len(name) > 63 {
    name = name[:63]
  }
  /***
  How usernames are processed:
  * If input is Trimino --> Becomes usr_trimino
  * If input is 12345 --> Becomes usr_12345 (Safely starts with an ASCII letter u)
  * If input is usr_finance --> Becomes usr_finance (Leaves the prefix alone and avoids duplication)
  * If input is select --> Becomes usr_select (Safely bypasses SQL keyword collisions)
  ***/
  return name
}
