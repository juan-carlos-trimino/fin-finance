package misc

import (
  "time"
)

type Misc struct{}

func (m Misc) DTF() string {
  /***
  time.Now() returns the current local time; using the current time in UTC.
  ***/
  return time.Now().UTC().Format(time.RFC3339Nano)
}
