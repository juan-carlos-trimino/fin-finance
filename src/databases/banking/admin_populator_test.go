package banking

import (
  "context"
  "fmt"
  "math/rand"
  "strings"
  "testing"
  "time"

  "github.com/jackc/pgx/v5/pgxpool"
  //go get github.com/brianvoe/gofakeit/v7
  "github.com/brianvoe/gofakeit/v7"
)

/***
Why an automated test file is the best choice:
1. Production Isolation: Files ending in _test.go are completely ignored by the Go compiler during standard application builds (go build).
   It guarantees that zero mock code or third-party fake data packages bloat or compromise your production cloud binaries.
2. Instant Developer Execution: You can spin up thousands of realistic user records inside your isolated local database schema anytime
   with a single command line flag.
3. No Route or UI Pollution: You don't have to build hidden administrative endpoints or buttons in your UI layout forms that someone could
   accidentally find and click.

---

How to use it in your terminal

When you want to spin up data inside your local development machine:
1. Open your terminal to the package directory.
2. Open populator_test.go and comment out the t.Skip(...) line.
3. Run this command:
   go test -v -run=TestGenerateMockData
4. Once execution finishes, uncomment t.Skip(...) to lock the file back down safely.
***/

// TestGenerateMockData acts as an executable command line task rather than a traditional assertions test.
// Run it with: go test -v -run=TestGenerateMockData
func TestGenerateMockData(t *testing.T) {
  //SAFETY CHECK: Comment this line out when you want to actively run the populator script locally.
  t.Skip("Skipping mock generation to prevent accidental local database overwrites")
  ctx := context.Background()
  //1. Initialize your local isolated database connection details
  //Replace with your actual development environment connection string
  //connStr := "postgres://postgres:secret@localhost:5432/finances?sslmode=disable"
  connStr := "host=localhost port=5432 user=admin_user password=12345 dbname=finances connect_timeout=4 sslmode=disable"
  pool, err := pgxpool.New(ctx, connStr)
  if err != nil {
    t.Fatalf("Failed to connect to local database pool: %v", err)
  }
  defer pool.Close()
  //2. Explicitly force the isolated schema execution path context
  _, err = pool.Exec(ctx, `SET search_path TO fin, public;`)
  if err != nil {
    t.Fatalf("Failed to bind target isolated search_path: %v", err)
  }
  //3. Define total generation count
  rowCount := 150
  gofakeit.Seed(time.Now().UnixNano())
  specialChars := []string{"%", "#", "@", "!", "$", "_", "-"}
  digitChars := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
  //Pool of common European accented/special characters to test ICU collation
  euroChars := []string{"é", "è", "à", "ç", "ù", "ñ", "í", "ó", "á", "ü", "ö", "ä", "ß", "ø", "å", "æ"}
  euroCaps := []string{"É", "È", "À", "Ç", "Ñ", "Í", "Ó", "Á", "Ü", "Ö", "Ä", "Ø", "Å", "Æ"}
  t.Logf("Executing stored procedure %s for %d rows...", SP_ADD_CUSTOMER, rowCount)
  for idx := 0; idx < rowCount; idx++ {
    //Generate standard name fields
    firstName := gofakeit.FirstName()
    lastName := gofakeit.LastName()
    // 40% chance to inject European characters into the display names to test mixed-case persistence.
    if rand.Float32() < 0.40 {
      firstName = euroCaps[rand.Intn(len(euroCaps))] + firstName[1:]
      lastName = lastName + euroChars[rand.Intn(len(euroChars))]
    }
    var middleName *string
    if rand.Float32() > 0.30 {  //70% chance of middle name presence
      mName := gofakeit.FirstName()
      if rand.Float32() < 0.40 {
        mName = euroCaps[rand.Intn(len(euroCaps))] + mName[1:]
      }
      middleName = &mName
    }
    //Generate a username that matches your slider rules
    //Handle username logic supporting the special character slider criteria
    username := strings.ToLower(fmt.Sprintf("%s%d", firstName, gofakeit.Number(10, 999)))
    //Decide the starting character rule based on explicit percentage weights
    roll := rand.Float32()
    switch {
    case roll < 0.15:
      // 15% chance to strictly start with a special character (e.g., "@john52")
      username = specialChars[rand.Intn(len(specialChars))] + username
    case roll < 0.30:
      // 15% chance to strictly start with a digit character (e.g., "7john52")
      username = digitChars[rand.Intn(len(digitChars))] + username
    case roll < 0.60:
      // 30% chance to explicitly start with an international European letter (e.g., "ñatasha41")
      // This directly validates that your "A-G", "H-N" slider code routes them correctly via ICU
      username = euroChars[rand.Intn(len(euroChars))] + username
    default:
      // Remaining 40% keep the normal alpha start character (e.g., "john52")
    }

    //Generate all realistic profile properties required by your SP
    password := gofakeit.Password(true, true, true, false, false, 12)  //Secure mock text password
    marketing := gofakeit.Bool()
    birthDate := gofakeit.DateRange(time.Now().AddDate(-65, 0, 0), time.Now().AddDate(-18, 0, 0))  //18-65 years old
    gender := gofakeit.Gender()
    email := gofakeit.Email()
    phone := gofakeit.Phone()
    //Generate a valid physical address matrix
    addressInfo := gofakeit.Address()
    address1 := addressInfo.Address
    address2 := ""
    if rand.Float32() > 0.75 {
      address2 = fmt.Sprintf("Apt %d", gofakeit.Number(1, 400))
    }
    city := addressInfo.City
    state := addressInfo.State
    country := "USA"
    zipCode := addressInfo.Zip
    //Execute the stored procedure directly against the connection pool
    //The parameter positions map character-for-character to your code definition snippet
    _, err = pool.Exec(ctx, SP_ADD_CUSTOMER,
      username,    // c.User_name
      password,    // c.Password
      firstName,   // c.First_name
      middleName,  // c.Middle_name (safely handles *string nulls)
      lastName,    // c.Last_name
      marketing,   // c.Marketing
      birthDate,   // c.Birth_date
      gender,      // c.Gender
      address1,    // c.Address1
      address2,    // c.Address2
      city,        // c.City
      state,       // c.State
      country,     // c.Country
      zipCode,     // c.Zip_code
      email,       // c.Email
      phone,       // c.Phone
    )
    //
    if err != nil {
      t.Fatalf("Stored procedure failed on execution item %d: %v", idx, err)
    }
  }
  t.Logf("Success! Generated %d complex nested customer accounts securely via SP validation loops.", rowCount)
}
