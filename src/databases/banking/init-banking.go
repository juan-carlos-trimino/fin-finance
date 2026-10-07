package banking

//To fold all block comments:
//  Ctrl+K and Ctrl+/
//To unfold all block comments:
//  Ctrl+K and Ctrl+J

import (
  "bytes"
  "context"
  "errors"
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
  // Global package level scope (compiled exactly once at startup)
  //Strict input validation: Allow only letters, numbers, and underscores.
  //It stops malicious SQL characters (like ';', '--', '"') completely (preventing SQL injection).
  isDbNameValid = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString


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

This design works well with concurrent goroutines running inside the same application.

Whether the competition is coming from completely separate servers (Kubernetes pods), separate processes on your local machine, or multiple goroutines spinning inside a single Go binary, the central Postgres database acts as the single source of truth.
Why it works for both Pods and Goroutines
Because pg_advisory_lock is executed at the database connection layer, Postgres assigns the lock to the specific database connection session (the TCP socket network connection).
* With Pods: Each pod opens its own connection and competes for the lock.
* With Goroutines: If your goroutines are each opening their own independent database connections (or grabbing a fresh connection from a pgxpool.Pool), Postgres treats them exactly like separate pods. It safely blocks the secondary goroutine connections until the first goroutine releases the lock.

A Crucial Performance Optimization for Goroutines
If you are using this strategy heavily inside goroutines, you can optimize your code to avoid hitting the network database engine unnecessarily.
Network calls to Postgres take milliseconds. Go memory operations take nanoseconds. To keep your application running fast, you can add a local sync.Mutex directly inside your Go code as a first line of defense.

package db

import (
  "context"
  "sync"
  "://github.com"
)

A Crucial Performance Optimization for Goroutines


If a user tries to access tenant_alpha and it triggers the provisioning sequence, the global mutex locks down. If a completely different user tries to access tenant_beta a millisecond later, Tenant Beta is forced to wait in line behind Tenant Alpha, even though they are completely unrelated databases! In a production system under heavy load, one slow tenant provisioning script could cause timeouts across your entire platform.
To fix this, you need to lock per tenant database name, rather than locking the entire function.

Instead of a single global lock and a standard map, you can use a thread-safe sync.Map to dynamically manage individual, lightweight locks for each specific tenant database.

package main

import (
  "context"
  "fmt"
  "sync"
)

// 1. A thread-safe registry to store individual locks per tenant database
var tenantLockRegistry sync.Map

func ProvisionDatabaseWithLock(ctx context.Context, adminConn *pgx.Conn, targetDb string) bool {
  // 2. Fetch or create a specific mutex JUST for this targetDb name
  actualLock, _ := tenantLockRegistry.LoadOrStore(targetDb, &sync.Mutex{})
  tenantMutex := actualLock.(*sync.Mutex)

  // 3. FIRST DEFENSE: Only block goroutines targeting the SAME tenant database
  tenantMutex.Lock()
  defer tenantMutex.Unlock()

  // 4. SECOND DEFENSE: Double-Check Database Existence locally before touching the network.
  // (If a previous internal goroutine just finished building it, exit immediately)
  if checkIfDbExistsLocallyCached(targetDb) {
    return true
  }

  // 5. THIRD DEFENSE: Block other external Kubernetes Pods for this tenant across the network cluster
  // Use a deterministic hash of the targetDb string as the advisory lock ID
  provisionLockID := int64(hashStringToInt(targetDb))
  _, err := adminConn.Exec(ctx, "SELECT pg_advisory_lock($1);", provisionLockID)
  if err != nil {
    return false
  }
  defer adminConn.Exec(ctx, "SELECT pg_advisory_unlock($1);", provisionLockID)

  // 6. Final Re-verification on Postgres (in case a separate server pod provisioned it)
  if !checkIfDbExistsOnCluster(ctx, adminConn, targetDb) {
    runMigrationScripts(targetDb)
    markDbAsLocallyCached(targetDb)
  }

  return true
}



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
  //
  if !isDbNameValid(targetDb) {
    logger.LogInfo(fmt.Sprintf("Invalid database name %q: Names must only contain alphanumeric characters or underscores", targetDb), correlationId)
    return false
  } else if !isDbNameValid(defaultDb) {
    logger.LogInfo(fmt.Sprintf("Invalid database name %q: Names must only contain alphanumeric characters or underscores", defaultDb), correlationId)
    return false
  }
  //Read all items in the target migration directory.
  files, err := os.ReadDir(dirPath)
  if err != nil {
    logger.LogError(fmt.Sprintf("Failed to read directory %q: %v", dirPath, err), correlationId)
    return false
  }
  //Pre-allocate slice capacity based on the directory size to prevent resize-copy loops.
  var sqlFiles []string = make([]string, 0, len(files))
  var does000AdminExist bool = false
  //Filter out non-SQL files and collect file names.
  for _, file := range files {
    if file.IsDir() {
      continue
    }
    //Return a read-only reference to the already-existing string block in memory (no new allocations).
    name := file.Name()
    strLen := len(name)
    /***
    Ensure the filename is long enough to have a 4-character extension. If a file name is exactly 4 characters long and matches
    the condition, its name would literally be ".sql". A file named just ".sql" is a hidden file with no actual base name, which
    is almost certainly a user error or temporary system file -- not a valid migration script.

    Slicing a string in Go does not copy string data or allocate memory. It simply constructs a new lightweight header window
    pointing to the end of the existing string data.
    ***/
    if strLen > 4 && strings.EqualFold(name[strLen - 4:], ".sql") {
      sqlFiles = append(sqlFiles, name)
      //The string comparison only triggers until the file is found; case-insensitive check with 0 allocations.
      if (!does000AdminExist && strings.EqualFold(name, "000_admin.sql")) {
        does000AdminExist = true
      }
    }
  }
  /***
  Why Mandating the Bootstrap/Migration Files be Present
  1. It Guarantees "Idempotency" and Reproducibility
     In modern DevOps (CI/CD pipelines, Kubernetes, local staging), you need to be able to tear down a database and recreate it
     exactly the same way every single time. If 000_admin.sql is missing, you can never spin up a fresh local development
     environment or a temporary preview environment. Your code becomes dependent on a manual state (someone having created the
     database ahead of time).
  2. It Eliminates the "Ghost Migration" Danger
     If 000_admin.sql is missing but the database exists, your program assumes everything is fine and proceeds to run the remaining
     migration files. However, you have no guarantee that the pre-existing database matches the state 000_admin.sql was supposed to
     create (it might be missing required schemas, extensions, or system roles).
  3. It Enforces Directory Integrity
     A migration directory should be treated as a single, immutable unit of truth. If individual .sql files are missing, it usually
     means a deployment failed, a git merge went wrong, or a developer forgot to commit a file. Silent fallbacks can accidentally
     mask serious human errors.
  ***/
  if !does000AdminExist {  //Enforce directory completeness up front.
    logger.LogError(fmt.Sprintf("CRITICAL CONFIGURATION ERROR: The bootstrap script '000_admin.sql' is missing from %s. " +
      "Migrations cannot be safely verified or executed without the root bootstrap file.", dirPath), correlationId)
    return false
  }
  //CRUCIAL: Sort files alphabetically so 000_admin runs before 001_admin_xxx.
  sort.Strings(sqlFiles)
  db := GetBsInstance()
  ctx := context.Background()
  //Query all applied migrations from the database in order.
  rows, err := db.bsPool.Query(ctx, "SELECT migration_name FROM migration_history ORDER BY migration_name ASC;")//????????????????????
  if err != nil {
    //Zero-Allocation Type Assertion to extract native Postgres Error codes.
    var pgErr *pgconn.PgError
    /***
    Check native structural errors without any string conversions or memory allocations.
    * If the database does not exist, Postgres returns error code 3D000 (invalid_catalog_name).
    * If the database exists but the table is missing, Postgres returns code 42P01 (undefined_table).
    ***/
    if errors.As(err, &pgErr) && (pgErr.Code == "42P01" || pgErr.Code == "3D000") {
      logger.LogInfo("No migration history table found or database does not exist. Proceeding with setup.", correlationId)
    } else {
      // 2. Performance optimization: Only use fmt.Sprintf inside the failure block
      // where the application is crashing/returning false anyway.
      logger.LogError(fmt.Sprintf("Failed to query database migration history: %v", err), correlationId)
      return false
    }
  } else {

    // 1. Pre-allocate a string slice matching your maximum possible local files.
    // By passing 0 for length and len(sqlFiles) for capacity, the slice starts empty
    // but guarantees ZERO heap allocations from resizing as it populates.

        // 1. Allocate a string slice with zero length but pre-defined upper capacity.
    // This creates an initial underlying backing array large enough to hold all potential files,
    // guaranteeing ZERO data-copying or array-doubling overhead as the rows stream in.

    preAllocatedSlice := make([]string, 0, len(sqlFiles))

        // 2. Append directly to our pre-sized buffer space.
    // AppendRows automatically handles loop iterations, column parsing, and connection safety cleanup.
//pgx.AppendRows behaves identically to pgx.CollectRows in how it iterates, scans, and automatically closes rows. However, it takes an existing slice as its first argument and appends rows directly to that slice.
    //Collect the database rows into a clean string slice.
    dbMigrations, err := pgx.AppendRows(preAllocatedSlice, rows, pgx.RowTo[string])
    if err != nil {
      logger.LogError(fmt.Sprintf("Failed to parse migration history rows: %v", err), correlationId)
      return false
    }




    /***
    Check once up front if the database records exceed the disk files. We use (len(sqlFiles) - 1) to account for the unlogged
    000_admin.sql.
    ***/
    if len(dbMigrations) > (len(sqlFiles) - 1) {
      logger.LogError("CRITICAL INTEGRITY FAILURE: The database has more migrations recorded than files found on disk!", correlationId)
      return false
    }
    //Direct Slice Comparison (Skipping index 0 because 000_admin isn't in migrations).
    for idx, dbMig := range dbMigrations {
      //Offset the index by +1 because sqlFiles[0] is 000_admin.sql.
      fileIdx := idx + 1
      //Strip the .sql extension so it matches the string stored in the database; using zero-allocation.
      strLen := len(sqlFiles[fileIdx])
      //Compare the file names at the exact same position.
      if !strings.EqualFold(dbMig, sqlFiles[fileIdx][: strLen - 4]) {
        logger.LogError(fmt.Sprintf("CRITICAL FILE MISMATCH: Database expected %q at position %d, but disk has %q. A file is missing!",
          dbMig, fileIdx, sqlFiles[fileIdx]), correlationId)
        return false
      }
    }
    /***
    Check if the environment is fully caught up with no new files pending; we add 1 to dbMigrations to account for 000_admin.sql
    which isn't logged in the migrations table.
    ***/
    if len(sqlFiles) == (len(dbMigrations) + 1) {
      logger.LogInfo("Database is fully up-to-date. No new migrations to process.", correlationId)
      return true
    }
  }
  //Iterate over the sorted files sequentially.
  for idx, fileName := range sqlFiles {
    fullPath := filepath.Join(dirPath, fileName)
    // //Default to the target database for all normal migration updates.
    // activeDb := targetDb
    // if idx == 0 {
    //   //Database may need to be created so let's use defaultDb (postgres).
    //   activeDb = defaultDb
    // }
    logger.LogInfo(fmt.Sprintf("Executing migration step [%d/%d]: %s against DB: %s", idx + 1, len(sqlFiles), fileName, targetDb), correlationId)
    //Using key-value pairs.
    var cmd *exec.Cmd
    if idx == 0 {
      /***
      Scenario A (First Deploy / New Environment):
      The database does not exist. 000_admin.sql runs entirely. It creates the database shell and executes every single standard
      table layout or stored procedure inside it. It finishes successfully. Go then proceeds to look at your incremental patches
      (001_xxx.sql onwards).
      ***/
      connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s connect_timeout=%d sslmode=%s", host, port, user, password,
        defaultDb, connect_timeout, sslmode)
      cmd = exec.Command("psql", connString,
        "-f", fullPath,
        "-v", fmt.Sprintf("ALWAYS_DB_ADMIN=%s", strconv.FormatBool(config.GetAlwaysDbAdmin(correlationId))),
        "-v", fmt.Sprintf("DB_NAME=%s", targetDb))
    } else {
      /***
      Scenario B (Subsequent Application Boots):
      The database is already fully intact. Because it's alive, it implies that the baseline snapshot has already run completely on
      a previous startup. Running 000_admin.sql again would be a waste of time and would error out trying to recreate existing
      tables. It hits the \q statement, exits gracefully with code 0, and allows Go to jump instantly to parsing the incremental
      patches.
      ***/
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
      if bytes.Contains(out, []byte("55P03")) || bytes.Contains(bytes.ToLower(out), []byte("lock timeout")) {
        logger.LogError(fmt.Sprintf("[DB Session] Postgres aborted: Could not acquire lock within 5 seconds via psql on step %s.", fileName), correlationId)
        return false
      }
      /***
      When psql throws an error, it doesn't print a clean, single line. It outputs a multi-line stack trace. Modern log aggregators
      (like Datadog, AWS CloudWatch, Splunk, or Grafana Loki) assume that "one log line = one distinct event". Hence, when they see
      the newlines inside the psql error output, they slice the single error report into multiple separate log messages in the
      dashboard. The first line gets tied to the correlationId, the other lines do not get timestamp, correlationId, etc.

      Replace actual newlines with a string representation so it stays on one physical line.
      ***/
      singleLineOutput := strings.ReplaceAll(string(out), "\n", " | ")
      //Optimization: Log the raw psql output alongside the error so you can see compile syntax errors.
      logger.LogError(fmt.Sprintf("SQL script failed at step %s: %v | Output: %s", fileName, err, singleLineOutput), correlationId)
      return false
    }
    temp := strings.TrimSpace(string(out))
    if temp != "" {
      singleLineOutput := strings.ReplaceAll(temp, "\n", " | ")
      logger.LogInfo(fmt.Sprintf("Migration step %s completed output: %s", fileName, singleLineOutput), correlationId)
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


/***

package main

import (
  "context"
  "fmt"
  "os/exec"
  "strings"
  "sync"
)

// 1. Thread-safe registry to store memory mutexes per tenant database
var tenantLockRegistry sync.Map

// 2. Thread-safe registry to cache which tenants are fully built
var provisionedTenants sync.Map

func ProvisionTenantDatabase(ctx context.Context, connectionString string, targetDb string, sqlFilePath string) bool {
  // ----------------------------------------------------
  // STEP 1: FAST LOCAL EARLY EXIT
  // ----------------------------------------------------
  // If this Go app already fully built this tenant, exit immediately.
  // Zero mutex contention, zero OS forks, zero network calls.
  if _, alreadyDone := provisionedTenants.Load(targetDb); alreadyDone {
    return true
  }

  // ----------------------------------------------------
  // STEP 2: SCOPED MUTEX (No cross-tenant bottlenecks)
  // ----------------------------------------------------
  // Fetch or create a specific mutex JUST for this tenant database string.
  actualLock, _ := tenantLockRegistry.LoadOrStore(targetDb, &sync.Mutex{})
  tenantMutex := actualLock.(*sync.Mutex)

  tenantMutex.Lock()
  // Re-verify the check under the safety of the lock (Double-Checked Locking)
  if _, alreadyDone := provisionedTenants.Load(targetDb); alreadyDone {
    tenantMutex.Unlock()
    return true
  }
  defer tenantMutex.Unlock()

  // ----------------------------------------------------
  // STEP 3: RUN THE PSQL COMMAND (Postgres takes the wheel)
  // ----------------------------------------------------
  // Go spins up the migration script. The script itself will execute
  // its lock logic (e.g., SET lock_timeout = '5s'; SELECT pg_advisory_lock(hash);)
  cmd := exec.Command("psql", connectionString, "-f", sqlFilePath)
  out, err := cmd.CombinedOutput()

  if err != nil {
    outStr := string(out)
    outLower := strings.ToLower(outStr)

    // Catch your 5s lock timeout string or error code returned by the script
    if strings.Contains(outStr, "55P03") || strings.Contains(outLower, "lock_timeout") || strings.Contains(outLower, "lock timeout") {
      logger.LogError(fmt.Sprintf("[DB Session] Postgres aborted: Script could not acquire advisory lock for %s.", targetDb), correlationId)
      return false
    }

    logger.LogError(fmt.Sprintf("Migration failed for tenant %s: %v\nOutput: %s", targetDb, err, outStr), correlationId)
    return false
  }

  // ----------------------------------------------------
  // STEP 4: MARK AS DONE LOCALLY
  // ----------------------------------------------------
  // Save to memory cache so subsequent goroutines bypass this entire process.
  provisionedTenants.Store(targetDb, true)
  return true
}



***/
