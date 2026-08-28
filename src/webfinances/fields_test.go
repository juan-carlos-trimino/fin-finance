package webfinances

import (
  "fmt"
  "math/rand"
  "os"  //Added to clean up test directories afterward
  "sync"
  "testing"
  "time"
)

// =======================================
// CONCURRENCY STRESS TEST (The Race Test)
// =======================================
func TestAddDeleteConcurrency(t *testing.T) {
  // FIX: Override your production directory to point to a safe user-writable path.
  // Make sure 'mainDir' matches the exact case of your global variable name.
  oldDir := mainDir
  mainDir = "/tmp/finances_test"
  defer func() {
    mainDir = oldDir // Restore production setting after the test finishes
    os.RemoveAll("/tmp/finances_test") // Clean up generated test folders
  }()

  var wg sync.WaitGroup
  numGoroutines := 2000
  testUser := "JohnDoe"

  t.Logf("Spawning %d goroutines fighting to Add/Delete '%s' simultaneously...", numGoroutines, testUser)

  for i := 0; i < numGoroutines; i++ {
    wg.Add(1)
    go func(id int) {
      defer wg.Done()

      time.Sleep(time.Duration(rand.Intn(100)) * time.Microsecond)

      if id%2 == 0 {
        AddSessionDataPerUser(testUser, fmt.Sprintf("corr-%d", id))
      } else {
        DeleteSessionDataPerUser(testUser)
      }
    }(i)
  }

  wg.Wait()
  t.Log("Successfully completed 2,000 rapid chaotic operations without panicking!")
}

// ====================================
// PERFORMANCE PROFILER (The Benchmark)
// ====================================
func BenchmarkSessionThroughput(b *testing.B) {
  // FIX: Apply the same folder protection to the benchmark run
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
