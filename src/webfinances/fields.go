package webfinances

//
//To fold all block comments:
//  Ctrl+K and Ctrl+/
//To unfold all block comments:
//  Ctrl+K and Ctrl+J

import (
  "encoding/json"
  "errors"
  "fmt"
  "github.com/juan-carlos-trimino/gplogger"
  "github.com/juan-carlos-trimino/gposu"
  "os"
)

var mainDir string = "/finances"

func SetupDirStructure(dir string) {
  mainDir = dir + mainDir
  dirErr, err := osu.CreateDirs(0o077, 0o777, mainDir)
  if err != nil {
    panic("Cannot create directory '" + dirErr + "': " + err.Error())
  }
}

//Store the fields for each user in memory.
var currentFields = map[string]*fields{}  //key: user, value: fields

/***
Browsers default to use the same collection of cookies regardless of whether you are opening a
duplicate web page in a new tab or a new browser instance. Hence, two different tabs or two
instances of the same browser will look like the same session to the server.

Because multiple instances use the same cookies, the server cannot tell requests from them apart,
and it will associate them with the same Session data because they all have the same SessionID.
***/
type fields struct {
  //Make the pointers unexported so that clients can't interact with them directly but only via
  //exported methods.
  miscellaneous *miscellaneousFields
  mortgage *mortgageFields
  bonds *bondsFields
  adFv *adFvFields
  adPv *adPvFields
  adCp *adCpFields
  adEpp *adEppFields
  oaCp *oaCpFields
  oaEpp *oaEppFields
  oaFv *oaFvFields
  oaGa *oaGaFields
  oaInterestRate *oaInterestRateFields
  oaPerpetuity *oaPerpetuityFields
  oaPv *oaPvFields
  siAccurate *siAccurateFields
  siBankers *siBankersFields
  siOrdinary *siOrdinaryFields
}


type siAccurateFields struct {
  MenuPage string `json:"menuPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1Time string `json:"fd1Time"`
  Fd1Leap bool `json:"fd1Leap"`
  Fd1TimePeriod string `json:"fd1TimePeriod"`
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1PV string `json:"fd1PV"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2Time string `json:"fd2Time"`
  Fd2TimePeriod string `json:"fd2TimePeriod"`
  Fd2Amount string `json:"fd2Amount"`
  Fd2PV string `json:"fd2PV"`
  Fd2Result string `json:"fd2Result"`
  //
  Fd3Time string `json:"fd3Time"`
  Fd3TimePeriod string `json:"fd3TimePeriod"`
  Fd3Interest string `json:"fd3Interest"`
  Fd3Compound string `json:"fd3Compound"`
  Fd3Amount string `json:"fd3Amount"`
  Fd3Result string `json:"fd3Result"`
  //
  Fd4Interest string `json:"fd4Interest"`
  Fd4Compound string `json:"fd4Compound"`
  Fd4Amount string `json:"fd4Amount"`
  Fd4PV string `json:"fd4PV"`
  Fd4Result string `json:"fd4Result"`
}

func newSiAccurateFields(dir1, dir2, correlationId string) *siAccurateFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "siaccurate.txt")
  if obj != nil {
    var s siAccurateFields
    err := json.Unmarshal(obj, &s)
    if err != nil {
      //Write error, but continue with default values.
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    } else {
      return &s
    }
  } else if err != nil {
    logger.LogError(fmt.Sprintf("%+v", err), correlationId)
  } else {
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "siaccurate.txt"), correlationId)
  }
  return &siAccurateFields {
    MenuPage: "",
    CurrentButton: "lhs-button1",
    //
    Fd1Time: "1",
    Fd1Leap: false,
    Fd1TimePeriod: "year",
    Fd1Interest: "1.00",
    Fd1Compound: "annually",
    Fd1PV: "1.00",
    Fd1Result: "",
    //
    Fd2Time: "1",
    Fd2TimePeriod: "year",
    Fd2Amount: "1.00",
    Fd2PV: "1.00",
    Fd2Result: "",
    //
    Fd3Time: "1",
    Fd3TimePeriod: "year",
    Fd3Interest: "1.0",
    Fd3Compound: "annually",
    Fd3Amount: "1.00",
    Fd3Result: "",
    //
    Fd4Interest: "1.00",
    Fd4Compound: "annually",
    Fd4Amount: "1.00",
    Fd4PV: "1.00",
    Fd4Result: "",
  }
}

func getSiAccurateFields(userName string) *siAccurateFields {
  return currentFields[userName].siAccurate
}

type siBankersFields struct {
  MenuPage string `json:"menuPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1Time string `json:"fd1Time"`
  Fd1TimePeriod string `json:"fd1TimePeriod"`
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1PV string `json:"fd1PV"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2Time string `json:"fd2Time"`
  Fd2TimePeriod string `json:"fd2TimePeriod"`
  Fd2Amount string `json:"fd2Amount"`
  Fd2PV string `json:"fd2PV"`
  Fd2Result string `json:"fd2Result"`
  //
  Fd3Time string `json:"fd3Time"`
  Fd3TimePeriod string `json:"fd3TimePeriod"`
  Fd3Interest string `json:"fd3Interest"`
  Fd3Compound string `json:"fd3Compound"`
  Fd3Amount string `json:"fd3Amount"`
  Fd3Result string `json:"fd3Result"`
  //
  Fd4Interest string `json:"fd4Interest"`
  Fd4Compound string `json:"fd4Compound"`
  Fd4Amount string `json:"fd4Amount"`
  Fd4PV string `json:"fd4PV"`
  Fd4Result string `json:"fd4Result"`
}

func newSiBankersFields(dir1, dir2, correlationId string) *siBankersFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "sibankers.txt")
  if obj != nil {
    var s siBankersFields
    err := json.Unmarshal(obj, &s)
    if err != nil {
      //Write error, but continue with default values.
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    } else {
      return &s
    }
  } else if err != nil {
    logger.LogError(fmt.Sprintf("%+v", err), correlationId)
  } else {
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "sibankers.txt"), correlationId)
  }
  return &siBankersFields {
    MenuPage: "",
    CurrentButton: "lhs-button1",
    //
    Fd1Time: "1",
    Fd1TimePeriod: "year",
    Fd1Interest: "1.00",
    Fd1Compound: "annually",
    Fd1PV: "1.00",
    Fd1Result: "",
    //
    Fd2Time: "1",
    Fd2TimePeriod: "year",
    Fd2Amount: "1.00",
    Fd2PV: "1.00",
    Fd2Result: "",
    //
    Fd3Time: "1",
    Fd3TimePeriod: "year",
    Fd3Interest: "1.0",
    Fd3Compound: "annually",
    Fd3Amount: "1.00",
    Fd3Result: "",
    //
    Fd4Interest: "1.00",
    Fd4Compound: "annually",
    Fd4Amount: "1.00",
    Fd4PV: "1.00",
    Fd4Result: "",
  }
}

func getSiBankersFields(userName string) *siBankersFields {
  return currentFields[userName].siBankers
}

type siOrdinaryFields struct {
  MenuPage string `json:"menuPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1Time string `json:"fd1Time"`
  Fd1TimePeriod string `json:"fd1TimePeriod"`
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1PV string `json:"fd1PV"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2Time string `json:"fd2Time"`
  Fd2TimePeriod string `json:"fd2TimePeriod"`
  Fd2Amount string `json:"fd2Amount"`
  Fd2PV string `json:"fd2PV"`
  Fd2Result string `json:"fd2Result"`
  //
  Fd3Time string `json:"fd3Time"`
  Fd3TimePeriod string `json:"fd3TimePeriod"`
  Fd3Interest string `json:"fd3Interest"`
  Fd3Compound string `json:"fd3Compound"`
  Fd3Amount string `json:"fd3Amount"`
  Fd3Result string `json:"fd3Result"`
  //
  Fd4Interest string `json:"fd4Interest"`
  Fd4Compound string `json:"fd4Compound"`
  Fd4Amount string `json:"fd4Amount"`
  Fd4PV string `json:"fd4PV"`
  Fd4Result string `json:"fd4Result"`
}

func newSiOrdinaryFields(dir1, dir2, correlationId string) *siOrdinaryFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "siordinary.txt")
  if obj != nil {
    var s siOrdinaryFields
    err := json.Unmarshal(obj, &s)
    if err != nil {
      //Write error, but continue with default values.
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    } else {
      return &s
    }
  } else if err != nil {
    logger.LogError(fmt.Sprintf("%+v", err), correlationId)
  } else {
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "siordinary.txt"), correlationId)
  }
  return &siOrdinaryFields {
    MenuPage: "",
    CurrentButton: "lhs-button1",
    //
    Fd1Time: "1",
    Fd1TimePeriod: "year",
    Fd1Interest: "1.00",
    Fd1Compound: "annually",
    Fd1PV: "1.00",
    Fd1Result: "",
    //
    Fd2Time: "1",
    Fd2TimePeriod: "year",
    Fd2Amount: "1.00",
    Fd2PV: "1.00",
    Fd2Result: "",
    //
    Fd3Time: "1",
    Fd3TimePeriod: "year",
    Fd3Interest: "1.0",
    Fd3Compound: "annually",
    Fd3Amount: "1.00",
    Fd3Result: "",
    //
    Fd4Interest: "1.00",
    Fd4Compound: "annually",
    Fd4Amount: "1.00",
    Fd4PV: "1.00",
    Fd4Result: "",
  }
}

func getSiOrdinaryFields(userName string) *siOrdinaryFields {
  return currentFields[userName].siOrdinary
}

func AddSessionDataPerUser(userName, correlationId string) {
  if _, ok := currentFields[userName]; !ok {
    fd := &fields{
      miscellaneous: newMiscellaneousFields(mainDir, userName, correlationId),
      mortgage: newMortgageFields(mainDir, userName, correlationId),
      bonds: newBondsFields(mainDir, userName, correlationId),
      adFv: newAdFvFields(mainDir, userName, correlationId),
      adPv: newAdPvFields(mainDir, userName, correlationId),
      adCp: newAdCpFields(mainDir, userName, correlationId),
      adEpp: newAdEppFields(mainDir, userName, correlationId),
      oaCp: newOaCpFields(mainDir, userName, correlationId),
      oaEpp: newOaEppFields(mainDir, userName, correlationId),
      oaFv: newOaFvFields(mainDir, userName, correlationId),
      oaGa: newOaGaFields(mainDir, userName, correlationId),
      oaInterestRate: newOaInterestRateFields(mainDir, userName, correlationId),
      oaPerpetuity: newOaPerpetuityFields(mainDir, userName, correlationId),
      oaPv: newOaPvFields(mainDir, userName, correlationId),
      siAccurate: newSiAccurateFields(mainDir, userName, correlationId),
      siBankers: newSiBankersFields(mainDir, userName, correlationId),
      siOrdinary: newSiOrdinaryFields(mainDir, userName, correlationId),
    }
    currentFields[userName] = fd
  }
}

func DeleteSessionDataPerUser(userName string) {
  delete(currentFields, userName)
}

/***
Notes about synchronization:
(1) Even though each username has its own exclusive set of files, multiple users can share the same
    username. In this scenario, there is a possibility that two or more users may try to modify a
    file at the same time thereby corrupting the file. In order to prevent file corruption, the
    file is protected with a lock that allows a single writer (exclusive write) or multiple readers
    (share reads).
(2) It's always good practice to protect shared resources even if you're sure that there will be no
    conflict.
***/
func readFields(filePath string) ([]byte, error) {
  exists, err := osu.CheckFileExists(filePath)
  if exists {
    obj, err := osu.ReadAllShareLock1(filePath, os.O_RDONLY, 0o400)
    if err != nil {
      return nil, errors.New(fmt.Sprintf("Couldn't open file %s: ", filePath) + err.Error())
    }
    return obj, nil
  } else {
    return nil, err
  }
}
