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
  "os"
  "os/exec"
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

const (
  // Define a globally unique 64-bit constant for this service's migrations.
// Any random 64-bit integer fits, up to 9223372036854775807.
 MigrationLockValue int64 = 4829103948572019384
)


func GetMigrationLockID() string {
	return strconv.FormatInt(MigrationLockValue, 10)
}


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

func ExecuteSqlScripts(dirPath, host, user, password, defaultDb, targetDb, sslmode string, port, connect_timeout int, correlationId string) bool {
  if !isDbNameValid(targetDb) {
    logger.LogInfo(fmt.Sprintf("Invalid database name %q.", targetDb), correlationId)
    return false
  } else if !isDbNameValid(defaultDb) {
    logger.LogInfo(fmt.Sprintf("Invalid database name %q.", defaultDb), correlationId)
    return false
  }
  //Read all items in the target migration directory.
  files, err := os.ReadDir(dirPath)
  if err != nil {
    logger.LogError(fmt.Sprintf("Cannot read directory %q: %v", dirPath, err), correlationId)
    return false
  }
  //Pre-allocate slice capacity based on the directory size to prevent resize-copy loops.
  var sqlFiles []string = make([]string, 0, len(files))
  var does000AdminDbExist bool = false
  //Filter out non-SQL files and save SQL file names.
  for _, file := range files {
    if file.IsDir() {
      continue
    }
    /***
    Return a read-only reference to the already-existing string block in memory (no new allocations).
    Go is strictly a pass-by-value language so returning a string from a function in Go really means passing a copy of the 16-byte
    structure representing the string.
    ***/
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
      if (!does000AdminDbExist && strings.EqualFold(name, "000_admin_db.sql")) {
        does000AdminDbExist = true
      }
    }
  }
  /***
  Why Mandating the Bootstrap/Migration Files be Present
  1. It Guarantees "Idempotency" and Reproducibility
     In modern DevOps (CI/CD pipelines, Kubernetes, local staging), you need to be able to tear down a database and recreate it
     exactly the same way every single time. If 000_admin_db.sql is missing, you can never spin up a fresh local development
     environment or a temporary preview environment. Your code becomes dependent on a manual state (someone having created the
     database ahead of time).
  2. It Eliminates the "Ghost Migration" Danger
     If 000_admin_db.sql is missing but the database exists, your program assumes everything is fine and proceeds to run the remaining
     migration files. However, you have no guarantee that the pre-existing database matches the state 000_admin_db.sql was supposed to
     create (it might be missing required schemas, extensions, or system roles).
  3. It Enforces Directory Integrity
     A migration directory should be treated as a single, immutable unit of truth. If individual .sql files are missing, it usually
     means a deployment failed, a git merge went wrong, or a developer forgot to commit a file. Silent fallbacks can accidentally
     mask serious human errors.
  ***/
  if !does000AdminDbExist {  //Enforce directory completeness up front.
    logger.LogError(fmt.Sprintf("The bootstrap script '000_admin_db.sql' is missing from %s. " +
      "Migrations cannot be safely verified or executed without the root bootstrap file.", dirPath), correlationId)
    return false
  }
  /***
  Passing files into an external psql process just to have them hit an internal \if block and exit immediately is a massive waste
  of resources. Spawning an OS subprocess (exec.Command) consumes significant CPU cycles and operating system overhead.

  We can track an execution starting index using a variable (e.g., var unappliedMigrations int).
  * If the database is fresh, unappliedMigrations remains 0, running everything starting with 000_admin_db.sql.
  * If the database is already up-to-date or partially migrated, unappliedMigrations will equal len(dbMigrations) + 1. This safely
    bypasses 000_admin_db.sql (slot 0) plus all previously registered scripts, picking up exactly on the first unapplied file.
  ***/
  var unappliedMigrations int = 0
  //CRUCIAL: Sort files alphabetically so 000_admin_db runs before 001_admin_xxx.
  sort.Strings(sqlFiles)
  db := GetBsInstance()
  ctx := context.Background()
  //Query all applied migrations from the database in ASC order.
  rows, err := db.bsPool.Query(ctx, "SELECT migration_name FROM fin.migration_history ORDER BY migration_name ASC;")
  if err != nil {
    //Zero-Allocation Type Assertion to extract native Postgres Error codes.
    var pgErr *pgconn.PgError
    errExists := errors.As(err, &pgErr)
    /***
    Check native structural errors without any string conversions or memory allocations.
    * If the database does not exist, Postgres returns error code 3D000 (invalid_catalog_name).
    * If the database exists but the table is missing, Postgres returns code 42P01 (undefined_table).
    ***/
    if errExists && pgErr.Code == "3D000" {
      logger.LogInfo("Database does not exist. Proceeding with setup.", correlationId)
    } else if errExists && pgErr.Code == "42P01" {
      unappliedMigrations = 1  //Move pass the 000_admin_db.sql file.
      logger.LogInfo("No migration history table found in the database. Proceeding with setup.", correlationId)
    } else {
      //Performance optimization: Only use fmt.Sprintf inside the failure block where the application is crashing/returning false anyway.
      logger.LogError(fmt.Sprintf("Cannot query database migration history: %v", err), correlationId)
      return false
    }
  } else {
    /***
    Pre-allocate a string slice matching the maximum possible local files. By passing 0 for length and len(sqlFiles) for capacity,
    the slice starts empty but guarantees ZERO heap allocations from resizing as it populates.
    ***/
    expectedMigrations := make([]string, 0, len(sqlFiles))
    /***
    AppendRows automatically handles loop iterations, column parsing, and connection cleanup. Instead of creating a fresh slice
    from scratch, it takes an already allocated slice, fills it directly from the database row stream, and returns the result.
    ***/
    dbMigrations, err := pgx.AppendRows(expectedMigrations, rows, pgx.RowTo[string])
    if err != nil {
      logger.LogError(fmt.Sprintf("Cannot parse migration history rows: %v", err), correlationId)
      return false
    }
    /***
    Check once up front if the database records exceed the disk files. We use (len(sqlFiles) - 1) to account for the unlogged
    000_admin_db.sql.
    ***/
    if len(dbMigrations) > (len(sqlFiles) - 1) {
      logger.LogError("The database has more migrations recorded than files found on disk!", correlationId)
      return false
    }
    //Direct Slice Comparison (Skipping index 0 because 000_admin_db isn't in migrations).
    for idx, dbMig := range dbMigrations {
      //Offset the index by +1 because sqlFiles[0] is 000_admin_db.sql.
      fileIdx := idx + 1
      //Strip the .sql extension so it matches the string stored in the database; using zero-allocation.
      strLen := len(sqlFiles[fileIdx])
      //Compare the file names at the exact same position.
      if !strings.EqualFold(dbMig, sqlFiles[fileIdx][: strLen - 4]) {
        logger.LogError(fmt.Sprintf("Database expected %q at position %d, but disk has %q. A file is missing!",
          dbMig, fileIdx, sqlFiles[fileIdx]), correlationId)
        return false
      }
    }
    //Save the next unapplied file position.
    unappliedMigrations = len(dbMigrations) + 1
    /***
    Check if the environment is fully caught up with no new files pending; we add 1 to dbMigrations to account for 000_admin_db.sql
    which isn't logged in the migrations table.
    ***/
    if len(sqlFiles) == unappliedMigrations {
      logger.LogInfo("Database is fully up-to-date. No new migrations to process.", correlationId)
      return true
    }
  }
  //Track patches 1 to N separate from the root infrastructure bootstrap script.
  totalMigrations := len(sqlFiles) - 1
  //Iterate over the sorted files sequentially.
  for idx, fileName := range sqlFiles[unappliedMigrations:] {
    /***
    The range engine always starts counting its index from 0 for the elements it is actively iterating over, regardless of where
    the slice window began.
    ***/
    currentIdx := idx + unappliedMigrations
    if currentIdx == 0 {
      logger.LogInfo(fmt.Sprintf("Executing %s: creating DB %s", fileName, targetDb), correlationId)
    } else {
      logger.LogInfo(fmt.Sprintf("Executing migration step [%d/%d]: %s against DB %s", currentIdx, totalMigrations, fileName, targetDb), correlationId)
    }
    //Using key-value pairs.
    var cmd *exec.Cmd
    if currentIdx == 0 {
      /***
      Scenario A (First Deploy / New Environment):
      The database does not exist. 000_admin_db.sql runs entirely. It creates the database shell and executes every single standard
      table layout or stored procedure inside it. It finishes successfully. Go then proceeds to look at your incremental patches
      (001_xxx.sql onwards).
      ***/
      connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s connect_timeout=%d sslmode=%s", host, port, user, password,
        defaultDb, connect_timeout, sslmode)
      cmd = exec.Command("psql", "--quiet", "--tuples-only", "--no-align", connString,
        "-f", fileName,
        "-v", fmt.Sprintf("LOCK_ID=%s", GetMigrationLockID()),
        "-v", "LOCK_TIMEOUT='5s'",  //Keep the single quotes in the string value.
        "-v", fmt.Sprintf("ALWAYS_DB_ADMIN=%s", strconv.FormatBool(config.GetAlwaysDbAdmin(correlationId))),
        "-v", fmt.Sprintf("DB_NAME=%s", targetDb))
    } else {
      /***
      Scenario B (Subsequent Application Boots):
      The database is already fully intact. Because it's alive, it implies that the baseline snapshot has already run completely on
      a previous startup. Running 000_admin_db.sql again would be a waste of time and would error out trying to recreate existing
      tables. It hits the \q statement, exits gracefully with code 0, and allows Go to jump instantly to parsing the incremental
      patches.
      ***/
      connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s connect_timeout=%d sslmode=%s", host, port, user, password,
        targetDb, connect_timeout, sslmode)
      cmd = exec.Command("psql", "--quiet", "--tuples-only", "--no-align", connString,
        "-f", fileName,
        "-v", fmt.Sprintf("LOCK_ID=%s", GetMigrationLockID()),
        "-v", "LOCK_TIMEOUT='5s'",  //Keep the single quotes in the string value.
        "-v", fmt.Sprintf("MIGRATION_NAME=%s", fileName[: len(fileName) - 4]))
    }
    //Tell the OS to run the psql process INSIDE the migration directory (0 new heap allocations).
    cmd.Dir = dirPath
    //Run the command and wait for it to finish; return its combined standard output and standard error.
    out, err := cmd.CombinedOutput()
    if err != nil {
      //Catch if the failure was specifically due to the 5s lock timeout.
      if bytes.Contains(out, []byte("55P03")) || bytes.Contains(bytes.ToLower(out), []byte("lock timeout")) {
        logger.LogError(fmt.Sprintf("Postgres could not acquire lock within 5 seconds via psql on step %s.", fileName), correlationId)
        return false
      }
      /***
      When psql throws an error, it doesn't print a clean, single line. It outputs a multi-line stack trace. Modern log aggregators
      (like Datadog, AWS CloudWatch, Splunk, or Grafana Loki) assume that "one log line = one distinct event". Hence, when they see
      the newlines inside the psql error output, they slice the single error report into multiple separate log messages in the
      dashboard. The first line gets tied to the correlationId, the other lines do not get timestamp, correlationId, etc.

      Replace actual newlines with a string representation so it stays on one physical line.
      ***/
      out = bytes.ReplaceAll(out, []byte("\n"), []byte(" | "))
      //Log the raw psql output alongside the error so you can see compile syntax errors.
      logger.LogError(fmt.Sprintf("SQL script failed at step %s: %v | Output: %s", fileName, err, out), correlationId)
      return false
    }
    /***
    In Go, sub-slicing a byte slice (slice[start:end]) merely creates a tiny 24-byte header containing a pointer, a length,
    and a capacity. It points right back to the exact same underlying memory block that was already allocated.

    The only line that allocates heap memory inside the loop is the conversion to a string for the logger.
    ***/
    out = bytes.TrimSpace(out)
    if len(out) > 0 {
      //The commented lines output one long line.
      // out = bytes.ReplaceAll(out, []byte("\n"), []byte(" | "))
      // logger.LogInfo(fmt.Sprintf("Migration step %s completed output: %s", fileName, out), correlationId)
      sep := []byte("\n")
      for {
        idx := bytes.Index(out, sep)
        if idx == -1 {
          //It prevents the system from emitting an empty log line if the output happened to end right on a trailing
          //newline character.
          out = bytes.TrimSpace(out)
          if len(out) > 0 {
            logger.LogInfo(string(out), correlationId)
          }
          break
        }
        line := out[:idx]
        line = bytes.TrimSpace(line)
        if len(line) > 0 {
          logger.LogInfo(string(line), correlationId)
        }
        //Advance past the separator.
        out = out[idx + len(sep):]
      }
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
