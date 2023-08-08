package misc

import (
  "time"
  "unsafe"
)

type Misc struct{}

func (m Misc) DTF() string {
  /***
  time.Now() returns the current local time; using the current time in UTC.
  ***/
  return time.Now().UTC().Format(time.RFC3339Nano)
}

func (m Misc) Sizeof(i interface{}) (size uint) {
  // switch i := i.(type) {  //Type switch.
  // case struct:
    size = uint(unsafe.Sizeof(i.(struct{})))
  // }
  return
}

