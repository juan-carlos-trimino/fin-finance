package banking

//To fold all block comments:
//  Ctrl+K and Ctrl+/
//To unfold all block comments:
//  Ctrl+K and Ctrl+J

import (
  "context"
  "fmt"
  "github.com/google/uuid"
  "github.com/jackc/pgx/v5"
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
  // SP_CUSTOMER_INFO = "CALL bs.sp_customer_info(@userName, @password, @customerType, @firstName, " +
  //  "@middleName, @lastName, @dateOfBirth, @taxIdentifier, @address1, @address2, @cityName, " +
  //  "@stateName, @countryName, @zipCode, @primaryEmail, @secondaryEmail, @primaryPhone, " +
  //  "@secondaryPhone)"

  SP_CUSTOMER_INFO = "CALL bs.sp_customer_info(@userName, @password, @firstName, " +
   "@middleName, @lastName, @dateOfBirth, @taxIdentifier, @address1, @address2, @cityName, " +
   "@stateName, @countryName, @zipCode, @primaryEmail, @secondaryEmail, @primaryPhone, " +
   "@secondaryPhone)"
  QR_GET_ALL_CUSTOMERS = "SELECT customer_id, customer_type, username, password_hash, " +
   "created_at, updated_at FROM bs.tbl_customer"
  QR_GET_ALL_CUSTOMERS_INFO = "SELECT customer_info_id, customer_id, first_name, COALESCE(middle_name, ''), " +
   "last_name, date_of_birth, tax_identifier, address_1, COALESCE(address_2, ''), city_name, state_name, " +
   "country_name, COALESCE(zip_code, ''), COALESCE(primary_email, ''), COALESCE(secondary_email, ''), primary_phone, COALESCE(secondary_phone, ''), " +
   "created_at, updated_at FROM bs.tbl_customer_info"
    // fmt.Printf("%s - %s - %s - %s - %s - %s -- %s -- %s - %s - %s - %s - %s - %s - %s - %s- %s - %s - %s\n",
    //  c.UserName, c.Password, c.CustomerType, c.FirstName, c.MiddleName, c.LastName, c.DateOfBirth,
    //  c.TaxIdentifier, c.Address1, c.Address2, c.CityName, c.StateName, c.CountryName, c.ZipCode, c.PrimaryEmail,
    //  c.SecondaryEmail, c.PrimaryPhone, c.SecondaryPhone)
)

type Customer struct {
	Id uuid.UUID // Maps to a PostgreSQL UUID column
  // Id UUID //`json: "customer_id"`
  Username string //`json: "username"`
  Password_hash string //`json: "password_hash"`
  First_name string
  Middle_name string
  Last_name string
  Birth_date time.Time
	Gender byte
  Address_1 string
  Address_2 string
  City_name string
  State_name string
  Country_name string
  Zip_code string
  Primary_email string
  Secondary_email string
  Primary_phone string
  Secondary_phone string
  Created_at time.Time
  Updated_at time.Time
}




func (bs *banking) CustomerInfo(ctx context.Context, userName, password, customerType,
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
