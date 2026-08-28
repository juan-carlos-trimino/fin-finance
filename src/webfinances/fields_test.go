package webfinances

/***
Run the Tests From the Package Directory
----------------------------------------
Open a terminal, navigate directly into the webfinances folder, and execute the test commands:
$ cd path/to/webfinances

* Run the Race Detector *
$ go test -v -race -run=TestAddDeleteConcurrency

Look for a clean PASS output.

* Run the Benchmark *
$ go test -bench=BenchmarkSessionThroughput -benchmem
or (10 seconds; for minutes use 10m)
$ go test -bench=BenchmarkSessionThroughput -benchmem -benchtime=10s
or (exactly 5,000 times)
$ go test -bench=BenchmarkSessionThroughput -benchmem -benchtime=5000x
or (five separate 5-second trials)
$ go test -bench=BenchmarkSessionThroughput -benchmem -benchtime=5s -count=5

Once the benchmark finishes, see the values for:
ns/op (Nanoseconds per operation)
B/op (Memory bytes allocated per operation)

Example:
...
{"date_time":"2026-08-28T02:29:59.240672537Z","level":"INFO","correlation_id":"bench-corr-id","msg":"File /tmp/finances_bench/David/siordinary.txt does not exit."}
    4956            219954 ns/op           29215 B/op        190 allocs/op
PASS
ok      finance/webfinances     2.129s

* Metric Breakdown *
4956 iterations: In this benchmark window (default is 1 second), the system executed nearly 5,000 full lifecycle operations (fully
  instantiating all of the field segments, writing/checking folders, and completely tearing them down via a deletion cycle).

219,954 ns/op (Nanoseconds per Operation): This translates to roughly 0.21 milliseconds per full cycle (Add + Delete). For an engine
  interacting with the local filesystem (/tmp), this is incredibly fast.

29,215 B/op and 190 allocs/op: When creating all of the custom field groups, Go allocates memory pointers for things like
  newBondsFields, newMortgageFields, etc. About 29KB across 190 memory allocations per session initialization is lightweight and
  easily handled by Go's automated garbage collector.

* Why the test took 2.129 seconds *
When you ran your previous benchmark (1 second), your terminal output concluded with:
PASS
ok finance/webfinances 2.129s

This happened because you are running a parallel benchmark using b.RunParallel. The total wall-clock time can be slightly longer
than 1 second because:
* Goroutine Handoff & Scheduling: It takes a moment to spawn and distribute the work across all of your CPU threads.
* Teardown Cleanups: Your code relies on filesystem operations (os.RemoveAll), which must completely finish executing before the
  final ok status is printed to your terminal.
***/

import (
  "fmt"
  "math/rand"  //Using Go's native, thread-safe global rand functions (see documentation).
  "os"  //Added to clean up test directories afterward
  "sync"
  "testing"
  "time"
)

// =======================================
// CONCURRENCY STRESS TEST (The Race Test)
// =======================================
func TestAddDeleteConcurrency(t *testing.T) {
  /***
  When the test code runs AddSessionDataPerUser, the code attempts to create a physical directory at /finances/ (see 'mainDir'
  in fields.go), which is at the root of the Linux filesystem. Because the user account doesn't have root administrator write
  privileges for the root system directory, Linux blocks it, causing the code to trigger a panic.

  Because the code reads the directory from a global variable ('mainDir'), we can safely override it at the very beginning of
  the test function to use a directory the user does have permission to write to, such as the system's temporary directory (/tmp/).
  ***/
  oldDir := mainDir
  mainDir = "/tmp/finances_test"
  defer func() {
    mainDir = oldDir  //Restore production setting after the test finishes.
    os.RemoveAll("/tmp/finances_test")  //Clean up generated test folders
  }()
  numGoroutines := 2000
  testUser := "JohnDoe"
  t.Logf("Spawning %d goroutines fighting to Add/Delete '%s' simultaneously...", numGoroutines, testUser)
  var wg sync.WaitGroup
  for i := 0; i < numGoroutines; i++ {
    id := i  //Capture loop variable for the closure.
    /***
    Go 1.22+ Loop Variable Semantics: Starting in Go 1.22, Go redefined how for loops handle variables. Loop variables are now
    completely fresh memory instances allocated per iteration rather than per loop lifecycle.

    The id := i Explicit Capture: By introducing id := i explicitly inside the loop block, you are creating a block-scoped
    variable that is uniquely bundled to that exact loop pass. When the closure captures it, it binds to that distinct,
    immutable instance.

    wg.Go supports only the use of parameterless functions. wg.Go handles the wg.Add (increment the internal counter),
    goroutine execution, and automatic wg.Done (decrement the counter by 1) signaling under the hood.
    ***/
    wg.Go(func() {
      time.Sleep(time.Duration(rand.Intn(100)) * time.Microsecond)
      if id % 2 == 0 {
        AddSessionDataPerUser(testUser, fmt.Sprintf("corr-%d", id))
      } else {
        DeleteSessionDataPerUser(testUser)
      }
    })
  }
  wg.Wait()  //Block until the counter reaches 0.
  t.Log("Successfully completed 2,000 rapid chaotic operations without panicking!")
}

// ====================================
// PERFORMANCE PROFILER (The Benchmark)
// ====================================
func BenchmarkSessionThroughput(b *testing.B) {
  oldDir := mainDir
  mainDir = "/tmp/finances_bench"
  defer func() {
    mainDir = oldDir
    os.RemoveAll("/tmp/finances_bench")
  }()
  users := []string{"Alice", "Bob", "Charlie", "David", "Emma"}
  b.ResetTimer()
  b.RunParallel(func(pb *testing.PB) {
    for pb.Next() {
      user := users[rand.Intn(len(users))]
      AddSessionDataPerUser(user, "bench-corr-id")
      DeleteSessionDataPerUser(user)
    }
  })
}
