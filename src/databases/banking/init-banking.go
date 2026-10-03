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
  "github.com/juan-carlos-trimino/go-os"
  "github.com/juan-carlos-trimino/go-logger"
  "os/exec"
  "regexp"
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
Implement an Advisory Lock in Go

If the pods must all trigger this check from their Go main() functions on startup, we can use Postgres' built-in Distributed
Advisory Locks. Before a pod even looks at the SQL script, it asks the central database for a cluster-wide "lock" tied to a
random 64-bit number. Only one pod can hold the lock at a time. The other pods will wait automatically until the first pod
finishes creating the database.

Using a PostgreSQL Distributed Advisory Lock is the cleanest, industry-standard way to solve this architecture problem. Because advisory locks are managed entirely by the central PostgreSQL engine, they act as a cluster-wide mutex wrapper.
When your pods boot up simultaneously in Kubernetes:
1. Pod 1 acquires the lock and enters your script execution function.
2. Pod 2 and Pod 3 reach the lock statement, pause, and wait cleanly.
3. Pod 1 checks if the database exists (it doesn't), executes psql, creates the user's database layout, and releases the lock.
4. Pod 2 immediately wakes up, takes the lock, checks if the database exists (it now does), skips running the script, and safely proceeds to start your web server.

--------------------
Update your startup logic by adding an existence check and wrapping your ExecuteSqlScript function with a pg_advisory_lock.
---------------------
package db

import (
	"context"
	"fmt"
	"log"
	"strings"
	"://github.com"
)

// ProvisionDatabaseWithLock locks the cluster before checking/running your psql script
func ProvisionDatabaseWithLock(ctx context.Context, adminConn *pgx.Conn, targetDb, host, user, password, defaultDb, sslmode string, port, connect_timeout int, pathToScript, correlationId string) bool {

	// 1. Define a unique 64-bit integer lock key for your app provisioning system
	const provisionLockID int64 = 8492049103

	log.Printf("[%s] Waiting to acquire cluster-wide advisory lock...", correlationId)

	// 2. This statement will block all secondary pods cleanly until the lock is available
	_, err := adminConn.Exec(ctx, "SELECT pg_advisory_lock($1);", provisionLockID)
	if err != nil {
		log.Printf("[%s] ERROR: Failed to acquire advisory lock: %v", correlationId, err)
		return false
	}

	// 3. CRITICAL: Release the lock when this pod completes, waking up the next waiting pod
	defer func() {
		_, err := adminConn.Exec(ctx, "SELECT pg_advisory_unlock($1);", provisionLockID)
		if err != nil {
			log.Printf("[%s] WARNING: Failed to release advisory lock: %v", correlationId, err)
		} else {
			log.Printf("[%s] Released cluster-wide advisory lock.", correlationId)
		}
	}()

	// 4. Check if the database has already been created by a previous pod
	exists, err := CheckIfDatabaseExists(ctx, adminConn, targetDb, correlationId)
	if err != nil {
		log.Printf("[%s] ERROR: Database existence check failed: %v", correlationId, err)
		return false
	}

	// 5. If it already exists, skip running the script entirely!
	if exists {
		log.Printf("[%s] Database '%s' already initialized by another pod. Skipping execution.", correlationId, targetDb)
		return true
	}

	log.Printf("[%s] Database '%s' not found. Launching psql script runner...", correlationId, targetDb)

	// 6. Run your existing psql execution command safely
	return ExecuteSqlScript(host, user, password, defaultDb, targetDb, sslmode, port, connect_timeout, pathToScript, correlationId)
}

// CheckIfDatabaseExists verifies if the dynamic target database is in the pg_database catalog
func CheckIfDatabaseExists(ctx context.Context, adminConn *pgx.Conn, dbName string, correlationId string) (bool, error) {
	var exists bool

	// Enforce lowercase lookups since SanitizeDbName structures identifiers as lowercase
	cleanDbName := strings.ToLower(dbName)

	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1);"
	err := adminConn.QueryRow(ctx, query, cleanDbName).Scan(&exists)

	return exists, err
}

Why this is a resilient production design:
• No Resource Deadlocks: pg_advisory_lock ties the lock context strictly to the connection session. If a pod drops dead or hits an unexpected out-of-memory error mid-script, PostgreSQL will automatically release the lock, ensuring your other cluster pods don't hang indefinitely.
• Deterministic Deployment Log Profiles: Your logs will clearly show Pod 1 processing the execution block, while Pod 2 and 3 output clear Skipping execution messages, keeping your metrics legible.




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

***/


func ExecuteSqlScript(host, user, password, defaultDb, targetDb, sslmode string, port, connect_timeout int, pathToScript,
     correlationId string) bool {
  for _, str := range strings.Split(osu.ShowPermissions(pathToScript, false), "\n") {
    if str != "" {
      logger.LogInfo(str, correlationId)
    }
  }
  //Strict input validation: Allow only letters, numbers, and underscores.
  //It stops malicious SQL characters (like ';', '--', '"') completely (preventing SQL injection).
  isValidName := regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
  if !isValidName(targetDb) {
    logger.LogError(fmt.Sprintf("Invalid database name %q: Names must only contain alphanumeric characters or underscores", targetDb),
      correlationId)
    return false
  }
  //Using a connection URI (recommended).
  //postgresql://[user[:password]@][host[:port]]/[dbname][?option1=value1&option2=value2]
  // var connString string = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbname)
  //Using key-value pairs.
  connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s connect_timeout=%d sslmode=%s", host, port, user, password,
    defaultDb, connect_timeout, sslmode)



  db := GetBsInstance()
  //Choose a unique, random 64-bit ID for the application's setup lock.
  const clusterLockID = 9876543210    //???????????????????????????????????????????????
  //This blocks all other pods until the lock is acquired.
  _, err := db.bsPool.Exec(context.Background(), "SELECT pg_advisory_lock($1);", clusterLockID)
  if err != nil {
    panic(fmt.Sprintf("Failed to acquire cluster lock from Postgres: %v", err))
  }
  //Ensure the lock is released when this pod finishes, allowing the next pod to proceed.
  defer db.bsPool.Exec(context.Background(), "SELECT pg_advisory_unlock($1);", clusterLockID)   ///????????????
  //Now it is safe to check and run the script.
  //Pod 2 will wait here until Pod 1 has fully created the database.
  if !databaseExists(db, targetDb) {
    //ExecuteSqlScript(db, defaultDb)
  }







  // logger.LogInfo(fmt.Sprintf("Connection string: %s", connString), correlationId)
  //psql accepts two distinct connection string formats: URIs and Key-Value.
  cmd := exec.Command("psql", connString, "-f", pathToScript,
    "-v", fmt.Sprintf("ALWAYS_DB_ADMIN=%s", strconv.FormatBool(config.GetAlwaysDbAdmin(correlationId))),
    "-v", fmt.Sprintf("DB_NAME=%s", targetDb))
  //Run the command and returns its combined standard output and standard error.
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
