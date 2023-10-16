package sessions

import (
	"bufio"
	"finance/misc"
	"fmt"
	"os"
	"strings"
	"sync"
)

var m = misc.Misc{}
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

Suppose another thread (goroutine) in the process, written as follows, also requires access to the
two resources:

go func function2() {
  ...
  muMemory.Lock()
  defer muMemory.Unlock()
  muFile.Lock()
  defer muFile.Unlock()
  ...
}

That is, in the preceding function, the order of the calls to the locks has been switched. Because
of this switch, a deadlock might occur. Suppose that function1 begins executing and gains ownership
of the muFile lock. At the same time, function2 is executing and gains ownership of the muMemory
lock. Now there is a deadlock. When either function1 or function2 tries to continue executing,
neither function can gain ownership of the other lock it requires.

To solve this problem, you must always enter resource locks in exactly the same order everywhere in
your code. Notice that order does not matter when you call the unlock functions because these
functions never causes a thread (goroutine) to enter a wait state.
***/
var muFile sync.Mutex  //Protect the file.
/***
It allows read-only operations to proceed in parallel with each other, but write operations to have
fully exclusive access; this lock is called a multiple readers, single writer lock.

It's only profitable to use an RWMutex when most of the goroutines that acquire the lock are
readers, and the lock is under contention, that is, goroutines routinely have to wait to acquire
it. An RWMutex requires more complex internal bookkeeping , making it slower than a regular mutex
for uncontended locks.
***/
var muMemory sync.RWMutex
//Store the username and password for each user.
var users = map[string][]byte{}  //key: username, value: password

func UsersLength() int {
  return len(users)
}

func ValidateUser(username, password string) bool {
  var ok bool = false
  var hashedPassword []byte
  muMemory.RLock()  //Readers lock
  //Get the expected password from the in memory map.
  if hashedPassword, ok = users[username]; !ok {
    muMemory.RUnlock()
    return ok
  }
  muMemory.RUnlock()
  ok, _ = CompareHashAndPassword(hashedPassword, []byte(password))
  return ok
}



func AddFromMemory(username, password string) {
  hashPassword, _ := HashSecret(password)
  users[username] = hashPassword
}




func ReadUsersFromFile() error {
  var f *os.File
  var err error
  f, err = os.OpenFile("./files/user.txt", os.O_RDONLY, 0600)
  if err != nil {
    fmt.Printf("%s - %s\n", m.DTF(), err)
    return err
  }
  defer f.Close()
  builder := strings.Builder{}
  //Grow to a larger size to reduce future resizes of the buffer.
  builder.Grow(2048)
  muFile.Lock()  //Readers lock.
  defer muFile.Unlock()
  scanner := bufio.NewScanner(f)
  for scanner.Scan() {
    if builder.Len() == 0 {
      builder.WriteString(scanner.Text())
    } else {
      muMemory.Lock()  //Writer lock
      users[builder.String()] = scanner.Bytes()
      muMemory.Unlock()
      builder.Reset()
    }
  }
  return nil
}

func AddUserToFile(username, password string) error {
  hashPassword, _ := HashSecret(password)
  var f *os.File
  var err error

  //The leading zero forces a base-8 conversion. 0600

  //If the file doesn't exist, create it; otherwise, append to the file.
  f, err = os.OpenFile("./files/user.txt", os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0600)
//  os.FileMode(0600))

  fmt.Println("***** perm: ", os.FileMode(0600))


  if err != nil {
    fmt.Printf("%s - %s\n", m.DTF(), err)
    panic(err)
  }
  defer f.Close()
  hashPassword = append(hashPassword, 0x0A)  //Add LF.
  muFile.Lock()
  defer muFile.Unlock()
  if _, err = f.WriteString(username + "\n"); err != nil {
    fmt.Printf("%s - %s\n", m.DTF(), err)
  } else if _, err = f.Write(hashPassword); err != nil {
    fmt.Printf("%s - %s\n", m.DTF(), err)
  }
  return err
}
