package banking

//To fold all block comments:
//  Ctrl+K and Ctrl+/
//To unfold all block comments:
//  Ctrl+K and Ctrl+J

import (
  "context"
  "fmt"
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
The Singleton pattern in Go ensures that a specific type has only one instance throughout the
program's lifecycle and provides a global access point to that instance. This pattern is commonly
used for resources like database connections, loggers, or configuration managers where a single,
shared instance is desired.
bankingSystem represents the structure of the singleton instance.
***/
type banking struct {
  bsPool *pgxpool.Pool
}

var (
  /***
  To avoid creating multiple connection pools, use the Singleton pattern.
  poolInstance holds the single instance of the singleton; it is initialized to nil.
  ***/
  bsInstance *banking = nil
  /***
  The sync.Once type is a way to implement a thread-safe Singleton, guaranteeing that a function
  is executed only once, even with multiple concurrent goroutines.
  ***/
  bsOnce sync.Once
)

//Initialize the connection pool.
func InitializeBsPool(ctx context.Context, connString, correlationId string) *banking {
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
      bsInstance = &banking{bsPool: pool}
    }
  })
  //Return the single instance of the Singleton.
  if bsInstance == nil {
    return nil
  } else {
    return bsInstance
  }
}

func GetBsInstance() (*banking) {
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
  // var connString string = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password,
  //   host, port, dbname)
  // fmt.Println(connString)
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

//BoolPtr is a helper function to return a pointer to a bool.
func BoolPtr(b bool) *bool {
  return &b
}

//BytePtr is a helper function to return a pointer to a byte.
func BytePtr(b byte) *byte {
  return &b
}

func (bs *banking) VerifyConnection(ctx context.Context, correlationId string) bool {
  err := bs.bsPool.Ping(ctx)
  if err != nil {
    logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    return false
  }
  logger.LogInfo("Connected to PostgreSQL database!", correlationId)
  return true
}

func (bs *banking) Ping(ctx context.Context) error {
  return bs.bsPool.Ping(ctx)
}

func (bs *banking) Close() {
  bs.bsPool.Close()
}
