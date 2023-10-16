package misc

import (
	"os/user"
	"runtime"
	"strings"
)

//Is the current user running as root?
func IsRoot() (bool, error) {
  current, err := user.Current()
  if err != nil {
    return false, err
  }
  return strings.EqualFold(current.Username, "root"), nil
}

func DetermineOS() string {
  return runtime.GOOS
}
