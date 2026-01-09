package banking_system

//To fold all block comments:
//  Ctrl+K and Ctrl+/
//To unfold all block comments:
//  Ctrl+K and Ctrl+J

import (
  "context"
  "fmt"
  "github.com/google/uuid"
  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgxpool"
  "github.com/juan-carlos-trimino/go-os"
  "github.com/juan-carlos-trimino/gplogger"
  "os/exec"
  "strings"
  "sync"
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

/////////////
/////////////
/////////////
/////////////


// type Customer struct {
//   CustomerId int //`json: "customer_id"`
//   CustomerType string //`json: "customer_type"`
//   UserName string //`json: "username"`
//   Password string //`json: "password_hash"`
//   CreatedAt time.Time //`json: "created_at"`
//   UpdatedAt time.Time //`json: "updated_at"`
// }

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


/////////////
/////////////
/////////////
/////////////

/***
The Singleton pattern in Go ensures that a specific type has only one instance throughout the
program's lifecycle and provides a global access point to that instance. This pattern is commonly
used for resources like database connections, loggers, or configuration managers where a single,
shared instance is desired.
bankingSystem represents the structure of the singleton instance.
***/
type bankingSystem struct {
  bsPool *pgxpool.Pool
}

var (
  /***
  To avoid creating multiple connection pools, use the Singleton pattern.
  poolInstance holds the single instance of the singleton; it is initialized to nil.
  ***/
  bsInstance *bankingSystem = nil
  /***
  The sync.Once type is a way to implement a thread-safe Singleton, guaranteeing that a function
  is executed only once, even with multiple concurrent goroutines.
  ***/
  bsOnce sync.Once
)

//Initialize the connection pool.
func InitializeBsPool(ctx context.Context, connString, correlationId string) *bankingSystem {
  /***
  The anonymous function passed to bsOnce.Do will be executed only once across all calls to
  GetBsInstance(), even if multiple goroutines call it concurrently. This ensures thread-safe
  initialization.
  Lazy initialization: The Singleton instance is created only when GetBsInstance() is first called,
  not when the program starts.
  ***/
  bsOnce.Do(func() {
    logger.LogInfo("Initializing connection pool...", correlationId)
    pool, err := pgxpool.New(ctx, connString)
    if err != nil {
      logger.LogInfo(fmt.Sprintf("Unable to create connection pool: %v", err), correlationId)
    } else {
      bsInstance = &bankingSystem{bsPool: pool}
    }
  })
  //Return the single instance of the Singleton.
  if bsInstance == nil {
    return nil
  } else {
    return bsInstance
  }
}

func GetBsInstance() (*bankingSystem) {
  //Return the single instance of the Singleton.
  if bsInstance == nil {
    return nil
  } else {
    return bsInstance
  }
}

func ExecuteSqlScript(ctx context.Context, host, user, password, dbname, sslmode string, port int,
  pathToScript, correlationId string) bool {
  for _, str := range strings.Split(osu.ShowPermissions(pathToScript, false), "\n") {
    if str != "" {
      logger.LogInfo(str, correlationId)
    }
  }
  //postgresql://[user[:password]@][host[:port]]/[dbname][?option1=value1&option2=value2]
  //var connString string = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password,
  //host, port, dbname, sslmode)
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port,
   user, password, dbname, sslmode)
  conn, err := pgx.Connect(ctx, psqlInfo)
  if err != nil {
    logger.LogError(fmt.Sprintf("pgx.Connect failed: %v", err), correlationId)
    return false
  }
  defer conn.Close(ctx)
  cmd := exec.Command("psql", psqlInfo, "-f", pathToScript)
  out, err := cmd.CombinedOutput()
  if err != nil {
    logger.LogError(fmt.Sprintf("SQL script failed: %v", err), correlationId)
    return false
  }
  //
  for _, str := range strings.Split(string(out), "\n") {  //Convert []byte to string.
    if str != "" {
      logger.LogInfo(str, correlationId)
    }
  }
  return true
}

//StringPtr is a helper function to return a pointer to a string.
func StringPtr(s string) *string {
  return &s
}

//TimePtr is a helper function to return a pointer to a time.Time.
func TimePtr(t time.Time) *time.Time {
  return &t
}

func (bs *bankingSystem) VerifyConnection(ctx context.Context, correlationId string) bool {
  err := bs.bsPool.Ping(ctx)
  if err != nil {
    logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    return false
  }
  logger.LogInfo("Connected to PostgreSQL database!", correlationId)
  return true
}

func (bs *bankingSystem) Ping(ctx context.Context) error {
  return bs.bsPool.Ping(ctx)
}

func (bs *bankingSystem) Close() {
  bs.bsPool.Close()
}




/***
func (bs *bankingSystem) CustomerInfo(ctx context.Context, userName, password, customerType,
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
***/
