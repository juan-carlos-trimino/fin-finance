package config

import (
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "os"
  "strconv"
)

func GetServer() string {
  if ev, exists := os.LookupEnv("SERVER"); exists {
    return ev
  }
  return ""  //Default value.
}

func GetHttp(correlationId string) bool {
  if ev, exists := os.LookupEnv("HTTP"); exists {
    v, err := strconv.ParseBool(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: true).", err), correlationId)
    }
  }
  return true
}

func GetHttpPort(correlationId string) int {
  if ev, exists := os.LookupEnv("HTTP_PORT"); exists {
    v, err := strconv.Atoi(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: 8080).", err), correlationId)
    }
  }
  return 8080
}

func GetHttps(correlationId string) bool {
  if ev, exists := os.LookupEnv("HTTPS"); exists {
    v, err := strconv.ParseBool(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: false).", err), correlationId)
    }
  }
  return false
}

func GetHttpsPort(correlationId string) int {
  if ev, exists := os.LookupEnv("HTTPS_PORT"); exists {
    v, err := strconv.Atoi(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: 8443).", err), correlationId)
    }
  }
  return 8443
}

func GetK8s(correlationId string) bool {
  if ev, exists := os.LookupEnv("K8S"); exists {
    v, err := strconv.ParseBool(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: false).", err), correlationId)
    }
  }
  return false
}

func GetShutDownTimeout(correlationId string) int {
  if ev, exists := os.LookupEnv("SHUTDOWN_TIMEOUT"); exists {
    v, err := strconv.Atoi(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: 15 seconds).", err), correlationId)
    }
  }
  return 15  //Seconds.
}

func GetPprof(correlationId string) bool {
  if ev, exists := os.LookupEnv("PPROF"); exists {
    v, err := strconv.ParseBool(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: false).", err), correlationId)
    }
  }
  return false
}

func GetLetsEncryptCert(correlationId string) bool {
  if ev, exists := os.LookupEnv("LE_CERT"); exists {
    v, err := strconv.ParseBool(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: false).", err), correlationId)
    }
  }
  return false
}

func GetUser() string {
  if ev, exists := os.LookupEnv("USER"); exists {
    return ev
  }
  return ""  //Default value.
}

func GetPreventProbesOutput(correlationId string) bool {
  if ev, exists := os.LookupEnv("PREVENT_PROBES_OUTPUT"); exists {
    v, err := strconv.ParseBool(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: false).", err), correlationId)
    }
  }
  return false
}
