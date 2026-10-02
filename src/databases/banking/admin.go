package banking

//To fold all block comments:
//  Ctrl+K and Ctrl+/
//To unfold all block comments:
//  Ctrl+K and Ctrl+J

import (
  "context"
  "errors"
  "fmt"
  "github.com/jackc/pgx/v5"
  "github.com/juan-carlos-trimino/go-logger"
  "golang.org/x/crypto/bcrypt"
  "math"
  "math/rand"
  "time"
)

/***
Notes:
In pgx, you can use named arguments for stored procedures by using the pgx.NamedArgs type, where parameter names are prefixed with an @ symbol
in the query string. The library will then automatically rewrite the query to use positional parameters ($1, $2, etc.) before execution.
***/
const (
  //Use placeholder syntax (like $1, $2) to safely pass parameters to the function, preventing SQL injection.
  SP_ADD_CUSTOMER = "CALL fin.add_customer($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)"
  //Pass null for OUT parameter in the call.
  SP_AUTHENTICATE_USER = "CALL fin.authenticate_user($1, $2, $3, null, null)"
  //Change password.
  SP_CHANGE_PASSWORD = "CALL fin.change_password($1, $2, $3, $4, null)"

  QR_UPSERT_USER_DATA =
    `INSERT INTO fin.user_data_vault(user_name, partition_name, data_value)
    VALUES ($1, $2, $3)
    ON CONFLICT (user_name, partition_name)
    DO UPDATE SET
      data_value = EXCLUDED.data_value;`

  QR_SELECT_USER_DATA =
    `SELECT data_value
     FROM fin.user_data_vault
     WHERE user_name = $1 AND partition_name = $2;`

)

type Customer struct {
  User_name, Password, First_name,Last_name, Gender, Address1, City, State, Country, Email, Phone string
  Marketing bool
  //Nullable types.
  Address2, Middle_name, Zip_code *string
  Birth_date *time.Time
}

type CustomersContactDetails struct {  //Struct tags.
  Id int32  `db:"id"`
  Birth_date *time.Time  `db:"birth_date"`
  Gender string  `db:"gender"`
  Address1 string  `db:"address1"`
  Address2 *string  `db:"address2"`
  City_name string  `db:"city_name"`
  State_name string  `db:"state_name"`
  Country_name string  `db:"country_name"`
  Zip_code *string  `db:"zip_code"`
  Email string  `db:"email"`
  Phone string  `db:"phone"`
  Created_at time.Time  `db:"created_at"`
  Updated_at time.Time  `db:"updated_at"`
}

func DbAddCustomer(c *Customer, ctx context.Context, correlationId string) error {
  db := GetBsInstance()
  //Use Exec when the stored procedure does not return a result set.
  _, err := db.bsPool.Exec(ctx, SP_ADD_CUSTOMER, c.User_name, c.Password, c.First_name, c.Middle_name, c.Last_name, c.Marketing,
    c.Birth_date, c.Gender, c.Address1, c.Address2, c.City, c.State, c.Country, c.Zip_code, c.Email, c.Phone)
  if err != nil {
    logger.LogError(fmt.Sprintf("SP fin.add_customer: %v", err), correlationId)
  } else {
    logger.LogInfo(fmt.Sprintf("SP fin.add_customer succeeded. Username: %s", c.User_name), correlationId)
  }
  return err
}

/***
Hash in the backend/application layer, if possible: Hashing within the database can expose plain-text passwords in query logs. Hashing
the password in your application code before sending it to the database is often a safer practice. This function should be used during
user registration.
***/
func HashAndSaltPassword(password, correlationId string) string {
  /***
  Use strong algorithms: Modern algorithms such as Argon2, bcrypt, or PBKDF2 are recommended over older ones like MD5 or SHA-1, which are
  considered broken for password hashing. Use GenerateFromPassword to hash & salt the password. The cost can be any value you want, but
  DefaultCost is a good starting point.
  ***/
  hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    logger.LogError(fmt.Sprintf("Error on DbHashAndSaltPassword: %v", err), correlationId)
    return ""
  }
  /***
  The resulting hash string includes the salt and all necessary parameters, which should be stored in the database. GenerateFromPassword
  returns a byte slice, so convert it to a string for storage.
  ***/
  return string(hash)
}

func DbAuthenticateUser(ctx context.Context, userName, password, correlationId string) (bool, bool) {
  db := GetBsInstance()
  var status int
  var isAdmin bool = false
  var ok bool = true
  err := db.bsPool.QueryRow(ctx, SP_AUTHENTICATE_USER, userName, password, correlationId).Scan(&status, &isAdmin)
  if err != nil {
    logger.LogError(fmt.Sprintf("DbAuthenticateUser: %v", err), correlationId)
    ok = false
  } else if status < 0 {
    ok = false
  }
  return ok, isAdmin
}

func DbGetCustomersContactDetails(ctx context.Context, correlationId string) bool {
  db := GetBsInstance()
  var ok bool = true
  //Call the function returning SETOF/TABLE. Use to read data and inspect rows.
  rows, err := db.bsPool.Query(ctx, "SELECT * FROM fin.get_customers_contact_details()")
  if err != nil {
    logger.LogError(fmt.Sprintf("DbGetCustomersContactDetails: %v", err), correlationId)
    ok = false
  } else {
    var users []CustomersContactDetails
    //Automatically scans all rows into a slice of User structs.
    users, err := pgx.CollectRows(rows, pgx.RowToStructByName[CustomersContactDetails])
    if err != nil {
      logger.LogError(fmt.Sprintf("DbGetCustomersContactDetails: %v", err), correlationId)
      ok = false
    } else {
      for _, user := range users {
        fmt.Printf("Id: %d, Birth_date: %v, Gender: %s, Address1: %s, Address2: %s, City_name: %s, State_name: %s, Country_name: %s, Zip_code: %s, Email: %s, Phone: %s, Created_at: %v, Updated_at: %v\n", user.Id, user.Birth_date.Format("2006-01-02"), user.Gender, user.Address1, PtrString(user.Address2), user.City_name, user.State_name, user.Country_name, PtrString(user.Zip_code), user.Email, user.Phone, user.Created_at, user.Updated_at)
      }
    }
  }
  return ok
}

func DbChangePassword(ctx context.Context, userName, oldPassword, newPassword, correlationId string) bool {
  db := GetBsInstance()
  var ok bool = true
  err := db.bsPool.QueryRow(ctx, SP_CHANGE_PASSWORD, userName, oldPassword, newPassword, correlationId).Scan(&ok)
  if err != nil {
    logger.LogError(fmt.Sprintf("DbGetCustomersContactDetails: %v", err), correlationId)
    ok = false
  }
  return ok
}



func DbSaveUserData(ctx context.Context, userName, partitionName, correlationId string, jsonBytes []byte) bool {
  db := GetBsInstance()
  //Use to modify data or database state.
  result, err := db.bsPool.Exec(ctx, QR_UPSERT_USER_DATA, userName, partitionName, jsonBytes)
  if err != nil {
    logger.LogError(fmt.Sprintf("DbSaveUserData: %v", err), correlationId)
    return false
  }
  logger.LogInfo(fmt.Sprintf("DbSaveUserData: Succeeded. Rows affected %d", result.RowsAffected()), correlationId)
  return true
}


func DbFetchUserData(ctx context.Context, userName, partitionName, correlationId string) []byte {
  db := GetBsInstance()
  var jsonBytes []byte
  //Use if expecting exactly one row (or zero if not found).
  err := db.bsPool.QueryRow(ctx, QR_SELECT_USER_DATA, userName, partitionName).Scan(&jsonBytes)
  if err != nil {
    if errors.Is(err, pgx.ErrNoRows) {
      logger.LogInfo("DbFetchUserData: Succeeded. No record found.", correlationId)
      return []byte{}  //Empty slice.
    }
    logger.LogError(fmt.Sprintf("DbFetchUserData: %v", err), correlationId)
    return nil
  }
  return jsonBytes
}

//////////////////////



//In this really short post, we will demonstrate how to implement a retry mechanism with exponential
// backoff for failed operations in Go. This technique is particularly useful when interacting with
//  external services or APIs that may occasionally fail or respond with errors.
func retryWithExponentialBackoff(retries int, randomFactor float64, baseTime time.Duration, maxBackoff time.Duration) {
  //Calculate the backoff interval using exponential backoff with a base time.
  //Wait Time = (Base Time) * (2 ^ Number of Retries)
  //The exponential factor here is 2^n, where n is the number of retries already made.
  backoff := time.Duration(math.Min(float64(baseTime) * math.Pow(2, float64(retries)), float64(maxBackoff)))
  // Add jitter to the backoff to avoid retry collisions.
  //To prevent the request from retrying at the same interval as other requests, we add "jitter" or
  //randomness to the wait time. A common approach is to add a random amount of time up to the
  //calculated backoff time, which could look something like this:
  //Randomized Wait Time = Wait Time + (Random Factor * Wait Time)
  //If our random factor is 0.5 (or 50%), then after the second failed attempt, the actual wait
  //time could be anywhere from 2 seconds to 3 seconds (2 seconds + 0.5 * 2 seconds).
  //Float64 returns a pseudo-random number in the half-open interval [0.0,1.0) from the default Source.
  jitter := time.Duration(rand.Float64() * float64(backoff) * randomFactor)
  nextBackoff := backoff + jitter
  // Sleep for the backoff interval before retrying.
  time.Sleep(nextBackoff)
}


// ValidatePassword checks password complexity requirements
// func ValidatePassword(password string) error {
//     if len(password) < 8 {
//         return errors.New("password must be at least 8 characters")
//     }
//     // Add more password validation rules as needed
//     return nil
// }


/*
-- 1. Fetch the entire nested JSONB block
SELECT data_value
FROM user_data_vault
WHERE user_name = 'trimino'
  AND partition_name = 'finance_settings';

-- 2. Extract a distinct item property value (e.g., getting "USD") using ->>
SELECT data_value ->> 'currency' AS active_currency
FROM user_data_vault
WHERE user_name = 'trimino'
  AND partition_name = 'finance_settings';



package main

import (
  "context"
  "encoding/json"
  "fmt"
  "log"
  "time"

  "://github.com"
)

// VaultData matches the JSONB structure you want to save
type VaultData struct {
  Theme       string   `json:"theme"`
  Currency    string   `json:"currency"`
  MaxAccounts int      `json:"max_accounts"`
  Languages   []string `json:"languages"`
}

// ExecuteVaultOperations demonstrates both the UPSERT and SELECT queries
func ExecuteVaultOperations(ctx context.Context, conn *pgx.Conn) {
  userName := "trimino"
  partitionName := "finance_settings"

  // Create our structured data object
  myConfig := VaultData{
    Theme:       "dark",
    Currency:    "USD",
    MaxAccounts: 5,
    Languages:   []string{"en", "es"},
  }

  // 1. Marshall the Go struct directly into JSON bytes
  jsonData, err := json.Marshal(myConfig)
  if err != nil {
    log.Fatalf("Failed to marshal JSON: %v", err)
  }

  // ==========================================
  // 2. THE UPSERT QUERY
  // ==========================================
  upsertQuery := `
    INSERT INTO user_data_vault (user_name, partition_name, data_value)
    VALUES ($1, $2, $3)
    ON CONFLICT (user_name, partition_name)
    DO UPDATE SET
      data_value = EXCLUDED.data_value,
      updated_at = CURRENT_TIMESTAMP;
  `

  _, err = conn.Exec(ctx, upsertQuery, userName, partitionName, jsonData)
  if err != nil {
    log.Fatalf("Upsert failed: %v", err)
  }
  fmt.Println("Successfully upserted data!")

  // ==========================================
  // 3. THE SELECT QUERY (Retrieving full JSONB)
  // ==========================================
  selectQuery := `
    SELECT data_value
    FROM user_data_vault
    WHERE user_name = $1 AND partition_name = $2;
  `

  var rawJSON []byte
  err = conn.QueryRow(ctx, selectQuery, userName, partitionName).Scan(&rawJSON)
  if err != nil {
    if err == pgx.ErrNoRows {
      fmt.Println("No record found matching those keys.")
      return
    }
    log.Fatalf("Select failed: %v", err)
  }

  // 4. Unmarshal database bytes back into your Go struct
  var savedConfig VaultData
  if err := json.Unmarshal(rawJSON, &savedConfig); err != nil {
    log.Fatalf("Failed to unmarshal row JSON: %v", err)
  }

  // Output your recovered data fields
  fmt.Printf("Retrieved Config - Theme: %s, Currency: %s\n", savedConfig.Theme, savedConfig.Currency)
}


*/
