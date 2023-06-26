package misc

import (
  "time"
)

type Misc struct{}

func (m Misc) DTF() string {
  return time.Now().UTC().Format(time.RFC3339Nano)
}
