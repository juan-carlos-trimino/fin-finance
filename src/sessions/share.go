package sessions

import (
	"sync"
)

//Grouping together three related variables in a single package-level variable, protect.
var shr = struct {  //Unnamed struct.
  slock sync.RWMutex  //Lock for the sessions map.
  //Store the session information for each user in memory.
  sessions map[string]session_token  //key: sessionToken, value: session
  /***
  Each logical resource in your application should have its own lock that is used to synchronize
  access to any and all parts of the logical resource. You should not have a single lock for all
  logical resources as this reduces scalability if multiple threads (goroutines) are accessing
  different logical resources: only one thread (goroutine) will be allowed to execute at any one
  time.

  Sometimes you'll need to access two (or more) logical resources simultaneously. If each resource
  has its own lock, you have to use both locks to do all of this atomically. For example,

  go func function1() {
    ...
    muFile.Lock()
    defer muFile.Unlock()
    muMemory.Lock()
    defer muMemory.Unlock()
    ...
  }

  Suppose another thread (goroutine) in the process, written as follows, also requires access to
  the two resources:

  go func function2() {
    ...
    muMemory.Lock()
    defer muMemory.Unlock()
    muFile.Lock()
    defer muFile.Unlock()
    ...
  }

  That is, in the preceding function, the order of the calls to the locks has been switched.
  Because of this switch, a deadlock might occur. Suppose that function1 begins executing and gains
  ownership of the muFile lock. At the same time, function2 is executing and gains ownership of the
  muMemory lock. Now there is a deadlock. When either function1 or function2 tries to continue
  executing, neither function can gain ownership of the other lock it requires.

  To solve this problem, you must always enter resource locks in exactly the same order everywhere
  in your code. Notice that order does not matter when you call the unlock functions because these
  functions never causes a thread (goroutine) to enter a wait state.
  ***/
  mu sync.Mutex  //Protect the file.
  /***
  It allows read-only operations to proceed in parallel with each other, but write operations to
  have fully exclusive access; this lock is called a multiple readers, single writer lock.

  It's only profitable to use an RWMutex when most of the goroutines that acquire the lock are
  readers, and the lock is under contention, that is, goroutines routinely have to wait to acquire
  it. An RWMutex requires more complex internal bookkeeping , making it slower than a regular mutex
  for uncontended locks.
  ***/
  sync.RWMutex  //Protect the map; embedded field.
  //Using pointer and non-pointer (see main.go [h.mux = make(map[string]http.HandlerFunc, 128)])
  //Store the username and password for each user.
  // users map[string]*[]byte  //key: username, value: password
  users map[string][]byte  //key: username, value: password
}{
  // users: make(map[string]*[]byte, 16),
  users: make(map[string][]byte, 16),
  sessions: make(map[string]session_token, 16),  //key: sessionToken, value: session
}
