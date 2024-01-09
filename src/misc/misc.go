package misc

import (
  "os"
  "time"
)

type Misc struct{}

func (m Misc) DTF() string {
  /***
  time.Now() returns the current local time; using the current time in UTC.
  ***/
  return time.Now().UTC().Format(time.RFC3339Nano)
}

func CheckFileExists(fileName string) (bool, error) {
  info, err := os.Stat(fileName)
  if err != nil {
    if os.IsNotExist(err) {
      return false, nil
    } else {
      return false, err
    }
  }
  return !info.IsDir(), nil
}

func CheckDirExists(dirName string) (bool, error) {
  info, err := os.Stat(dirName)
  if err != nil {
    if os.IsNotExist(err) {
      return false, nil
    } else {
      return false, err
    }
  }
  return info.IsDir(), nil
}
