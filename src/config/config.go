package config

import (
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "os"
  "strconv"
)

func GetServer() string {
  ev, exists := os.LookupEnv("SERVER")
  if !exists {
    return ""  //Default value.
  }
  return ev
}

func GetHttp(correlationId string) bool {
  ev, exists := os.LookupEnv("HTTP")
  if exists {
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
  ev, exists := os.LookupEnv("HTTP_PORT")
  if exists {
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
  ev, exists := os.LookupEnv("HTTPS")
  if exists {
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
  ev, exists := os.LookupEnv("HTTPS_PORT")
  if exists {
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
  ev, exists := os.LookupEnv("K8S")
  if exists {
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
  ev, exists := os.LookupEnv("SHUTDOWN_TIMEOUT")
  if exists {
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
  ev, exists := os.LookupEnv("PPROF")
  if exists {
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
  ev, exists := os.LookupEnv("LE_CERT")
  if exists {
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
  ev, exists := os.LookupEnv("USER")
  if !exists {
    return ""  //Default value.
  }
  return ev
}

func GetPreventProbesOutput(correlationId string) bool {
  ev, exists := os.LookupEnv("PREVENT_PROBES_OUTPUT")
  if exists {
    v, err := strconv.ParseBool(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: false).", err), correlationId)
    }
  }
  return false
}
