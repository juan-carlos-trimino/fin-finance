package sessions

import (
  "bufio"
  "os"
  "strings"
  "syscall"
)

// const (
//   hashedPasswordLen uint = 64
//   passwordLen uint = 32
// )

func UsersLength() int {
  return len(shr.users)
}

func ValidateUser(username, password string) bool {
  var ok bool = false
  // var hashedPassword *[]byte
  var hashedPassword []byte
  shr.RLock()  //Readers lock.
  //Get the expected password from the in memory map.
  if hashedPassword, ok = shr.users[username]; !ok {
    shr.RUnlock()
    return ok
  }
  shr.RUnlock()
  // ok, _ = CompareHashAndPassword(*hashedPassword, []byte(password))
  ok, _ = CompareHashAndPassword(hashedPassword, []byte(password))
  return ok
}

func ReadUsersFromFile(filename string) error {
  var f *os.File
  var err error
  shr.mu.Lock()
  f, err = os.OpenFile(filename, os.O_RDONLY, 0440)
  if err != nil {
    shr.mu.Unlock()
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
  shr.mu.Unlock()
  shr.Lock()  //Writer lock
  for k, v := range usersTmp {
    // p := make([]byte, len(v))  //Capacity = len(v)
    // copy(p, v)
    // protect.users[k] = &p
    shr.users[k] = v
  }
  shr.Unlock()
  return nil
}

func AddUserToFile(filename, username, password string) error {
  hashPassword, _ := HashSecret(password)
  hashPassword = append(hashPassword, 0x0A)  //Add LF.
  var f *os.File
  var err error
  shr.mu.Lock()
  defer shr.mu.Unlock()
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
