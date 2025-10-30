package banking_system

//To fold all block comments:
//  Ctrl+K and Ctrl+/
//To unfold all block comments:
//  Ctrl+K and Ctrl+J

import (
  "context"
  "fmt"
  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgxpool"
  "github.com/juan-carlos-trimino/gplogger"
  "os/exec"
  "strings"
  "sync"
  "time"
)

/***
Notes:
In pgx, you can use named arguments for stored procedures by using the pgx.NamedArgs type, where parameter names are prefixed with an @ symbol in the query string. The library will then automatically rewrite the query to use positional parameters ($1, $2, etc.) before execution.
***/
const (
  SP_CUSTOMER_INFO = "CALL bs.sp_customer_info(@userName, @password, @customerType, @firstName, " +
   "@middleName, @lastName, @dateOfBirth, @taxIdentifier, @address1, @address2, @cityName, " +
   "@stateName, @countryName, @zipCode, @primaryEmail, @secondaryEmail, @primaryPhone, @secondaryPhone)"
)

/////////////
type Book struct {
    ID       int
    Title    string
    Author   string
    Quantity int
}

type BookStore interface {
    InsertBookIntoDatabase(Book) error
    GetBookDetailsByID(int) (*Book, error)
    GetAllBookDetails(int) (*[]Book, error)
    UpdateBookDetails(int, Book) error
    DeleteBookFromDatabase(int) error
}
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
func GetBsInstance(ctx context.Context, connString, correlationId string) (*bankingSystem, bool) {
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
    return nil, false
  } else {
    return bsInstance, true
  }
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




func ExecuteSqlScript(ctx context.Context, host, user, password, dbname, sslmode string, port int, correlationId string) bool {
  //postgresql://[user[:password]@][host[:port]]/[dbname][?option1=value1&option2=value2]
  //var connString string = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", "postgres", "postgres", "localhost", "5432", "postgres")
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port,
   user, password, dbname, sslmode)
  conn, err := pgx.Connect(ctx, psqlInfo)
  if err != nil {
    logger.LogError(fmt.Sprintf("%v", err), correlationId)
    return false
  }
  defer conn.Close(ctx)
  cmd := exec.Command("psql", psqlInfo, "-f", "../IaC-app/utilities/postgres/sql/baseline/banking-system.sql")
  out, err := cmd.CombinedOutput()
  if err != nil {
    logger.LogError(fmt.Sprintf("Execute SQL script failed: %v", err), correlationId)
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




func (bs *bankingSystem) CustomerInfo(ctx context.Context, userName, password, customerType,
  firstName, middleName, lastName *string, dateOfBirth *time.Time, taxIdentifier, address1, address2, cityName,
  stateName, countryName, zipCode, primaryEmail, secondaryEmail, primaryPhone, secondaryPhone *string,
  correlationId string) {
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
    // userName, password, customerType, firstName,
    // middleName, lastName, dateOfBirth, taxIdentifier, address1, address2, cityName, stateName,
    // countryName, zipCode, primaryEmail, secondaryEmail, primaryPhone, secondaryPhone)
  if err != nil {
    logger.LogError(fmt.Sprintf("Error calling stored procedure CustomerInfo: %v", err), correlationId)
  } else {
    logger.LogInfo("Stored procedure bs.sp_customer_info called successfully (no return value).", correlationId)
  }
}
