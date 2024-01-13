package sessions

import (
  "finance/misc"
//  "bufio"
  "os"
//  "strings"
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

func ReadUsersFromFile(fileName string) error {
  shr.Lock()  //Writer lock
  err := misc.ReadAllShareLock2(fileName, shr.users, os.O_RDONLY, 0400)
  shr.Unlock()
  return err
}

func AddUserToFile(filePath, userName, password string) error {
  hashPassword, _ := HashSecret(password)
  hashPassword = append(hashPassword, 0x0A)  //Add LF.
  _, _, err := misc.WriteAllExclusiveLock2(filePath, userName + "\n", hashPassword,
    os.O_CREATE | os.O_APPEND | os.O_RDWR, 0o600)
  return err
}
