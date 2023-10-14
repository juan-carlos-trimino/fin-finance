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
var muFile sync.Mutex  //Protect the file.
//Store the username and password for each user.
var users = map[string][]byte{}  //key: username, value: password

func UsersLength() int {
  return len(users)
}

func ValidateUser(username, password string) bool {
  var ok bool = false
  var hashedPassword []byte
  //Get the expected password from the in memory map.
  if hashedPassword, ok = users[username]; !ok {
    return ok
  }
  ok, _ = CompareHashAndPassword(hashedPassword, []byte(password))
  return ok
}

func ReadUsersFromFile() error {
  var f *os.File
  var err error
  f, err = os.OpenFile("./files/user.txt", os.O_APPEND | os.O_RDONLY, 0600)
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
      users[builder.String()] = scanner.Bytes()
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
  f, err = os.OpenFile("./files/user.txt", os.O_CREATE | os.O_APPEND | os.O_WRONLY,
  os.FileMode(0600))

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
