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

func GetHttp() bool {
  ev, exists := os.LookupEnv("HTTP")
  if exists {
    v, err := strconv.ParseBool(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: true).", err), "-1")
    }
  }
  return true
}

func GetHttpPort() int {
  ev, exists := os.LookupEnv("HTTP_PORT")
  if exists {
    v, err := strconv.Atoi(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: 8080).", err), "-1")
    }
  }
  return 8080
}

func GetHttps() bool {
  ev, exists := os.LookupEnv("HTTPS")
  if exists {
    v, err := strconv.ParseBool(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: false).", err), "-1")
    }
  }
  return false
}

func GetHttpsPort() int {
  ev, exists := os.LookupEnv("HTTPS_PORT")
  if exists {
    v, err := strconv.Atoi(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: 8443).", err), "-1")
    }
  }
  return 8443
}

func GetK8s() bool {
  ev, exists := os.LookupEnv("K8S")
  if exists {
    v, err := strconv.ParseBool(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: false).", err), "-1")
    }
  }
  return false
}

func GetShutDownTimeout() int {
  ev, exists := os.LookupEnv("SHUTDOWN_TIMEOUT")
  if exists {
    v, err := strconv.Atoi(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: 15 seconds).", err), "-1")
    }
  }
  return 15  //Seconds.
}

func GetPprof() bool {
  ev, exists := os.LookupEnv("PPROF")
  if exists {
    v, err := strconv.ParseBool(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: false).", err), "-1")
    }
  }
  return false
}

func GetLetsEncryptCert() bool {
  ev, exists := os.LookupEnv("LE_CERT")
  if exists {
    v, err := strconv.ParseBool(ev)
    if err == nil {
      return v
    } else {
      logger.LogInfo(fmt.Sprintf("%s - (Default: false).", err), "-1")
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
