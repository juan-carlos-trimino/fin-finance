package webfinances

import (
  "encoding/json"
  "fmt"
  "time"
)

const (
  INFO uint = 1
  WARN uint = 2
  ERROR uint = 3
  FATAL uint = 4
)

type LogEntry struct {  //INFO, WARN, ERROR, and FATAL.
  DateTime string `json:"date_time"`
  CorrelationId string `json:"correlation_id"`
  level string `json:"level"`
  Message []string `json:"msg"`
}

func (le *LogEntry) Print(level uint, cid string, msg []string) {
  if level == INFO {
    le.level = "INFO"
  } else if level == WARN {
    le.level = "WARN"
  } else if level == ERROR {
    le.level = "ERROR"
  } else {
    le.level = "FATAL"
  }
  le.DateTime = time.Now().UTC().Format(time.RFC3339Nano)
  le.CorrelationId = cid
  le.Message = msg
  //data, err := json.MarshalIndent(info, "", "  ")
  json, err := json.Marshal(le)
  if err != nil {
    panic(err)
  }
  fmt.Printf("%s\n", json)
}
