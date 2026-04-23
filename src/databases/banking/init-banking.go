package banking

//To fold all block comments:
//  Ctrl+K and Ctrl+/
//To unfold all block comments:
//  Ctrl+K and Ctrl+J

import (
  "context"
  "fmt"
  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgconn"
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
    //The database connection string can be in URL or keyword/value format.
    config, err := pgxpool.ParseConfig(connString)
    if err != nil {
      logger.LogInfo(fmt.Sprintf("Unable to create the pgxpool.Config: %v", err), correlationId)
    } else {
      //https://pkg.go.dev/github.com/jackc/pgx/v4/pgxpool#Config
      config.MaxConns = 25  //The maximum size of the pool.
      config.MinConns = 5  //The minimum size of the pool.
      //The duration after which an idle connection will be automatically closed by the health check.
      config.MaxConnIdleTime = 15 * time.Minute
      //The duration since creation after which a connection will be automatically closed.
      config.MaxConnLifetime = 1 * time.Hour
      //The duration between checks of the health of idle connections.
      config.HealthCheckPeriod = 1 * time.Minute
      config.PrepareConn = func(ctx context.Context, conn *pgx.Conn) (bool, error) {
        logger.LogInfo("Before acquiring the connection pool.", correlationId)
        return true, nil
      }
      config.AfterRelease = func(conn *pgx.Conn) bool {
        logger.LogInfo("After a connection is released, but before it is returned to the pool.",
          correlationId)
        return true
      }
      config.BeforeClose = func(conn *pgx.Conn) {
        logger.LogInfo("Before a connection is closed and removed from the pool.", correlationId)
      }
      // Set default query execution timeout
      // config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
      //Set the notice handler to capture a RAISE NOTICE (or INFO, WARNING, LOG, DEBUG) message.
      config.ConnConfig.OnNotice = func(conn *pgconn.PgConn, notice *pgconn.Notice) {
        logger.LogInfo(notice.Message, notice.Detail)
      }
      pool, err := pgxpool.NewWithConfig(ctx, config)
      if err != nil {
        logger.LogInfo(fmt.Sprintf("Unable to create connection pool: %v", err), correlationId)
      } else {
        bsInstance = &banking{bsPool: pool}
      }
    }
  })
  //Return the single instance of the Singleton or nil.
  return bsInstance
}

func GetBsInstance() (*banking) {
  //Return the single instance of the Singleton.
  if bsInstance == nil {
    return nil
  } else {
    return bsInstance
  }
}

func ExecuteSqlScript(ctx context.Context, host, user, password, dbname, sslmode string, port,
  connect_timeout int, pathToScript, correlationId string) bool {
  for _, str := range strings.Split(osu.ShowPermissions(pathToScript, false), "\n") {
    if str != "" {
      logger.LogInfo(str, correlationId)
    }
  }
  //postgresql://[user[:password]@][host[:port]]/[dbname][?option1=value1&option2=value2]
  // var connString string = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password,
  //   host, port, dbname)
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s connect_timeout=%d sslmode=%s",
    host, port, user, password, dbname, connect_timeout, sslmode)
  // logger.LogInfo(fmt.Sprintf("Connection string: %s", psqlInfo), correlationId)
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
  if s == "" {
    return nil
  } else {
    return &s
  }
}



func PtrString(s *string) string {
  if s == nil {
    return ""
  } else {
    return *s
  }
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
