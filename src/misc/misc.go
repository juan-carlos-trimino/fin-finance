package misc

import (
  "time"
  "unsafe"
)

type Misc struct{}

func (m Misc) DTF() string {
  return time.Now().UTC().Format(time.RFC3339Nano)
}

func (m Misc) Sizeof(i interface{}) (size uint) {
  // switch i := i.(type) {  //Type switch.
  // case struct:
    size = uint(unsafe.Sizeof(i.(struct{})))
  // }
  return
}

