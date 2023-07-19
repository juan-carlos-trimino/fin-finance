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

type LogEntry struct {
  DateTime string `json:"date_time"`
  CorrelationId string `json:"correlation_id"`
  Level string `json:"level"`
  Message []string `json:"msg"`
}

func (le *LogEntry) Print(level uint, cid string, msg []string) {
  if level == INFO {
    le.Level = "INFO"
  } else if level == WARN {
    le.Level = "WARN"
  } else if level == ERROR {
    le.Level = "ERROR"
  } else {
    le.Level = "FATAL"
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
