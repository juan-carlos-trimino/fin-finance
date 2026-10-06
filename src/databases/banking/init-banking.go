package banking

//To fold all block comments:
//  Ctrl+K and Ctrl+/
//To unfold all block comments:
//  Ctrl+K and Ctrl+J

import (
  "context"
  "fmt"
  "finance/config"
  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgconn"
  "github.com/jackc/pgx/v5/pgxpool"
  "github.com/juan-carlos-trimino/go-logger"
  "github.com/juan-carlos-trimino/go-os"
  "os"
  "os/exec"
  "path/filepath"
  "regexp"
  "sort"
  "strconv"
  "strings"
  "sync"
  "time"
)

/***
The Singleton pattern in Go ensures that a specific type has only one instance throughout the program's lifecycle and provides
a global access point to that instance. This pattern is commonly used for resources like database connections, loggers, or
configuration managers where a single, shared instance is desired.
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
  InitializeBsPool(), even if multiple goroutines call it concurrently. This ensures thread-safe
  initialization.
  Lazy initialization: The Singleton instance is created only when InitializeBsPool() is first
  called, not when the program starts.
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



/***

Yes, you are 100% correct. This design works beautifully and equally well with concurrent goroutines running inside the same application.
Whether the competition is coming from completely separate servers (Kubernetes pods), separate processes on your local machine, or multiple goroutines spinning inside a single Go binary, the central PostgreSQL database acts as the single source of truth.
Why it works for both Pods and Goroutines
Because pg_advisory_lock is executed at the database connection layer, PostgreSQL assigns the lock to the specific database connection session (the TCP socket network connection).
• With Pods: Each pod opens its own connection and competes for the lock.
• With Goroutines: If your goroutines are each opening their own independent database connections (or grabbing a fresh connection from a pgxpool.Pool), PostgreSQL treats them exactly like separate pods. It safely blocks the secondary goroutine connections until the first goroutine releases the lock.



A Crucial Performance Optimization for Goroutines
If you are using this strategy heavily inside goroutines, you can optimize your code to avoid hitting the network database engine unnecessarily.
Network calls to PostgreSQL take milliseconds. Go memory operations take nanoseconds. To keep your application running fast, you can add a local sync.Mutex directly inside your Go code as a first line of defense:

package db

import (
  "context"
  "sync"
  "://github.com"
)

// A local mutex to block goroutines inside the SAME application immediately
var localProvisionMutex sync.Mutex

func ProvisionDatabaseWithLock(ctx context.Context, adminConn *pgx.Conn, targetDb string, ...) bool {

  // 1. FIRST DEFENSE: Block other internal goroutines instantly without hitting the network
  localProvisionMutex.Lock()
  defer localProvisionMutex.Unlock()

  // 2. SECOND DEFENSE: Block other external Kubernetes Pods across the network cluster
  const provisionLockID int64 = 8492049103
  _, _ = adminConn.Exec(ctx, "SELECT pg_advisory_lock($1);", provisionLockID)
  defer adminConn.Exec(ctx, "SELECT pg_advisory_unlock($1);", provisionLockID)

  // 3. Check and run your script safely...
  return true
}


By combining Go’s local sync.Mutex with PostgreSQL’s pg_advisory_lock, you build an architecture that is optimized for goroutines locally, while remaining completely bulletproof against multi-pod scale-outs in production.



2. Centralized Bottlenecks
While often labeled as a "distributed lock", standard Postgres advisory locks are centralized inside a single database instance. If your primary database crashes, your lock manager goes down with it. (Note: Distributed database flavors like YugabyteDB and EDB Postgres Distributed support cluster-wide global advisory locks to alleviate this single point of failure).
3. Stackable Accumulation
Advisory locks are re-entrant. If your application code calls pg_advisory_lock(42) three times sequentially inside the same connection, it must call pg_advisory_unlock(42) exactly three times to fully release it.

To implement PostgreSQL advisory locks in Go, you need to handle connection pinning carefully, especially when using the standard database/sql package or the advanced pgx driver.
Because sql.DB manages a connection pool automatically, calling an advisory lock query directly on it is dangerous. The lock will attach to a random connection from the pool, and your subsequent queries might run on a different connection, losing the lock or leaving it orphaned.
***/


func ExecuteSqlScript(dirPath, host, user, password, defaultDb, targetDb, sslmode string, port, connect_timeout int, correlationId string) bool {

  for _, str := range strings.Split(osu.ShowPermissions(dirPath, false), "\n") {
    if str != "" {
      logger.LogInfo(str, correlationId)
    }
  }


  //Strict input validation: Allow only letters, numbers, and underscores.
  //It stops malicious SQL characters (like ';', '--', '"') completely (preventing SQL injection).
  isValidName := regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
  if !isValidName(targetDb) {
    logger.LogInfo(fmt.Sprintf("Invalid database name %q: Names must only contain alphanumeric characters or underscores", targetDb), correlationId)
    return false
  } else if !isValidName(defaultDb) {
    logger.LogInfo(fmt.Sprintf("Invalid database name %q: Names must only contain alphanumeric characters or underscores", defaultDb), correlationId)
    return false
  }
  //Read all items in the target migration directory.
  entries, err := os.ReadDir(dirPath)
  if err != nil {
    logger.LogError(fmt.Sprintf("Failed to read directory %s: %v", dirPath, err), correlationId)
    return false
  }
  //Filter out non-SQL files and collect file names.
  var sqlFiles []string
  var hasBootstrapFile bool = false
  for _, entry := range entries {
    if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".sql") {
      sqlFiles = append(sqlFiles, entry.Name())
      if (!hasBootstrapFile && entry.Name() == "000_admin.sql") {  //The string comparison only triggers until the file is found.
        hasBootstrapFile = true
      }
    }
  }
  //Handle empty or missing update file situations safely.
  if !hasBootstrapFile {
    //Check if the target database already exists on the server.
    checkConnString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s connect_timeout=%d sslmode=%s",
      host, port, user, password, defaultDb, connect_timeout, sslmode)
    //Query pg_database. If it returns text, the database exists.
    checkCmd := exec.Command("psql", checkConnString, "-tA", "-c", fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname='%s';", targetDb))
    dbCheckOut, checkErr := checkCmd.CombinedOutput()
    dbExists := strings.TrimSpace(string(dbCheckOut)) == "1"
    if checkErr != nil {
      logger.LogError(fmt.Sprintf("Failed to verify database existence during fallback check: %v", checkErr), correlationId)
      return false
    }
    //
    if !dbExists {
      //CRITICAL ROADBLOCK: The ecosystem is missing and we don't have the bootstrap script to build it.
      logger.LogError(fmt.Sprintf("CRITICAL ERROR: The target database %q does not exist, and the bootstrap script '000_admin.sql' was not found in %s. Execution halted.", targetDb, dirPath), correlationId)
      return false
    }
    //If the database already exists, having no new migration files is perfectly safe.
    if len(sqlFiles) == 0 {
      logger.LogInfo(fmt.Sprintf("Target database %q exists and no new migration files are pending.", targetDb), correlationId)
      return true
    }
  }
  //CRUCIAL: Sort files alphabetically so 000_admin_xxx runs before 001_admin_xxx.
  sort.Strings(sqlFiles)
  //Iterate over the sorted files sequentially.
  for idx, fileName := range sqlFiles {
    fullPath := filepath.Join(dirPath, fileName)
    //Resolve the active target name purely for explicit log reporting clarity.
    activeDb := defaultDb
    if idx > 0 {
      activeDb = targetDb
    }
    logger.LogInfo(fmt.Sprintf("Executing migration step [%d/%d]: %s against DB: %s", idx + 1, len(sqlFiles), fileName, activeDb), correlationId)
    //Using key-value pairs.
    var cmd *exec.Cmd
    if idx == 0 {
      connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s connect_timeout=%d sslmode=%s", host, port, user, password,
        defaultDb, connect_timeout, sslmode)
      cmd = exec.Command("psql", connString,
        "-f", fullPath,
        "-v", fmt.Sprintf("ALWAYS_DB_ADMIN=%s", strconv.FormatBool(config.GetAlwaysDbAdmin(correlationId))),
        "-v", fmt.Sprintf("DB_NAME=%s", targetDb))
    } else {
      connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s connect_timeout=%d sslmode=%s", host, port, user, password,
        targetDb, connect_timeout, sslmode)
      //We automatically strip out the extension to pass a clean text variable tag down to the psql layer context environment.
      cmd = exec.Command("psql", connString,
        "-f", fullPath,
        "-v", fmt.Sprintf("MIGRATION_NAME=%s", strings.TrimSuffix(fileName, ".sql")))
    }
    //Run the command and returns its combined standard output and standard error.
    out, err := cmd.CombinedOutput()
    if err != nil {
      //Catch if the failure was specifically due to the 5s lock timeout.
      if strings.Contains(string(out), "55P03") || strings.Contains(strings.ToLower(string(out)), "lock timeout") {
        logger.LogError(fmt.Sprintf("[DB Session] Postgres aborted: Could not acquire lock within 5 seconds via psql on step %s.", fileName), correlationId)
        return false
      }
      //Optimization: Log the raw psql output alongside the error so you can see compile syntax errors.
      logger.LogError(fmt.Sprintf("SQL script failed at step %s: %v\nOutput: %s", fileName, err, string(out)), correlationId)
      return false
    }
    //
    for _, str := range strings.Split(string(out), "\n") {  //Convert []byte to string.
      if str != "" {
        logger.LogInfo(str, correlationId)
      }
    }
  }
  return true
}




//Simple helper to safely check if the DB exists while under the lock.
func databaseExists(db *banking, dbName string) bool {
  var exists bool
  query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1);"
  _ = db.bsPool.QueryRow(context.Background(), query, dbName).Scan(&exists)
  return exists
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
