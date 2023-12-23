package sessions

import (
  "bufio"
  "os"
  "strings"
  "sync"
  "syscall"
)

// const (
//   hashedPasswordLen uint = 64
//   passwordLen uint = 32
// )

//Grouping together three related variables in a single package-level variable, protect.
var protect = struct {  //Unnamed struct.
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
  //Store the username and password for each user.
  users map[string]*[]byte  //key: username, value: password
} {
  users: make(map[string]*[]byte, 16),
}

func UsersLength() int {
  return len(protect.users)
}

func ValidateUser(username, password string) bool {
  var ok bool = false
  var hashedPassword *[]byte
  protect.RLock()  //Readers lock
  //Get the expected password from the in memory map.
  if hashedPassword, ok = protect.users[username]; !ok {
    protect.RUnlock()
    return ok
  }
  protect.RUnlock()
  ok, _ = CompareHashAndPassword(*hashedPassword, []byte(password))
  return ok
}

func ReadUsersFromFile(filename string) error {
  var f *os.File
  var err error
  protect.mu.Lock()
  f, err = os.OpenFile(filename, os.O_RDONLY, 0440)
  if err != nil {
    protect.mu.Unlock()
    return err
  }
  builder := strings.Builder{}
  //Grow to a larger size to reduce future resizes of the buffer.
  builder.Grow(1024)
  var usersTmp = map[string][]byte{}  //key: username, value: password
  scanner := bufio.NewScanner(f)
  for scanner.Scan() {
    if builder.Len() == 0 {
      builder.WriteString(scanner.Text())
    } else {
      usersTmp[builder.String()] = scanner.Bytes()
      builder.Reset()
    }
  }
  f.Close()
  protect.mu.Unlock()
  protect.Lock()  //Writer lock
  for k, v := range usersTmp {
    p := make([]byte, len(v))  //Capacity = len(v)
    copy(p, v)
    protect.users[k] = &p
  }
  protect.Unlock()
  return nil
}

func AddUserToFile(filename, username, password string) error {
  hashPassword, _ := HashSecret(password)
  hashPassword = append(hashPassword, 0x0A)  //Add LF.
  var f *os.File
  var err error
  protect.mu.Lock()
  defer protect.mu.Unlock()
  oldMask := syscall.Umask(0006)
  //If the file doesn't exist, create it; otherwise, append to the file.
  f, err = os.OpenFile(filename, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
  syscall.Umask(oldMask)
  if err != nil {
    return err
  }
  defer f.Close()
  if _, err = f.WriteString(username + "\n"); err == nil {
    _, err = f.Write(hashPassword)
  }
  return err
}

////////
//for testing
// func AddFromMemory(username, password string) {
//   hashPassword, _ := HashSecret(password)
//   protect.users[username] = hashPassword
// }
