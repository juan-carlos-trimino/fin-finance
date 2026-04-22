package banking

//To fold all block comments:
//  Ctrl+K and Ctrl+/
//To unfold all block comments:
//  Ctrl+K and Ctrl+J

import (
  "context"
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "golang.org/x/crypto/bcrypt"
  "math"
  "math/rand"
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
  //Pass null for OUT parameter in the call.
  SP_AUTHENTICATE_USER = "CALL fin.authenticate_user($1, $2, $3, null)"
)

type AddCustomer struct {
  User_name, Password, First_name, Last_name, Gender, Address1, City, State,
  Country, Email, Phone string
  Marketing bool
  //Nullable types.
  Address2, Middle_name, Zip_code *string
  Birth_date *time.Time
}

func DbAddCustomer(c *AddCustomer, ctx context.Context, correlationId string) error {
  db := GetBsInstance()
  //Use Exec when the stored procedure does not return a result set.
  _, err := db.bsPool.Exec(ctx, SP_ADD_CUSTOMER, c.User_name, c.Password, c.First_name,
    c.Middle_name, c.Last_name, c.Marketing, c.Birth_date, c.Gender, c.Address1, c.Address2,
    c.City, c.State, c.Country, c.Zip_code, c.Email, c.Phone)
  if err != nil {
    logger.LogError(fmt.Sprintf("SP fin.add_customer: %v", err), correlationId)
  } else {
    logger.LogInfo(fmt.Sprintf("SP fin.add_customer succeeded. Username: %s", c.User_name),
      correlationId)
  }
  return err
}

/***
Hash in the backend/application layer, if possible: Hashing within the database can expose
plain-text passwords in query logs. Hashing the password in your application code before sending
it to the database is often a safer practice.
This function should be used during user registration.
***/
func HashAndSaltPassword(password, correlationId string) string {
  /***
  Use strong algorithms: Modern algorithms such as Argon2, bcrypt, or PBKDF2 are recommended over
  older ones like MD5 or SHA-1, which are considered broken for password hashing.
  Use GenerateFromPassword to hash & salt the password. The cost can be any value you want, but
  DefaultCost is a good starting point.
  ***/
  hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    logger.LogError(fmt.Sprintf("Error on DbHashAndSaltPassword: %v", err), correlationId)
    return ""
  }
  /***
  The resulting hash string includes the salt and all necessary parameters, which should be
  stored in the database. GenerateFromPassword returns a byte slice, so convert it to a string
  for storage.
  ***/
  return string(hash)
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
