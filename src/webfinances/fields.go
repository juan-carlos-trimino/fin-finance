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

type miscellaneousFields struct {
  MenuPage string `json:"menuPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1Nominal string `json:"fd1Nominal"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1Result [2]string `json:"fd1Result"`
  //
  Fd2Effective string `json:"fd2Effective"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2Result [3]string `json:"fd2Result"`
  //
  Fd3Nominal string `json:"fd3Nominal"`
  Fd3Inflation string `json:"fd3Inflation"`
  Fd3Result [4]string `json:"fd3Result"`
  //
  Fd4CurrentRate string `json:"fd4CurrentRate"`
  Fd4CurrentCompound string `json:"fd4CurrentCompound"`
  Fd4NewCompound string `json:"fd4NewCompound"`
  Fd4Result string `json:"fd4Result"`
  //
  Fd5Interest string `json:"fd5Interest"`
  Fd5Compound string `json:"fd5Compound"`
  Fd5Factor string `json:"fd5Factor"`
  Fd5Result string `json:"fd5Result"`
  //
  Fd6Values string `json:"fd6Values"`
  Fd6Result [2]string `json:"fd6Result"`
  //
  Fd7Time string `json:"fd7Time"`
  Fd7TimePeriod string `json:"fd7TimePeriod"`
  Fd7Rate string `json:"fd7Rate"`
  Fd7Compound string `json:"fd7Compound"`
  Fd7PV string `json:"fd7PV"`
  Fd7Result string `json:"fd7Result"`
}

func newMiscellaneousFields(dir1, dir2, correlationId string) *miscellaneousFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "miscellaneous.txt")
  if obj != nil {
    var m miscellaneousFields
    err := json.Unmarshal(obj, &m)
    if err != nil {
      //Write error, but continue with default values.
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    } else {
      return &m
    }
  } else if err != nil {
    logger.LogError(fmt.Sprintf("%+v", err), correlationId)
  } else {
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "miscellaneous.txt"), correlationId)
  }
  return &miscellaneousFields {
    MenuPage: "",
    CurrentButton: "lhs-button1",
    //
    Fd1Nominal: "3.5",
    Fd1Compound: "monthly",
    Fd1Result: [2]string { misc_notes[0], "" },
    //
    Fd2Effective: "3.5",
    Fd2Compound: "monthly",
    Fd2Result: [3]string { misc_notes[0], misc_notes[1], "" },
    //
    Fd3Nominal: "2.0",
    Fd3Inflation: "2.0",
    Fd3Result: [4]string { misc_notes[1], misc_notes[2], misc_notes[3], "" },
    //
    Fd4CurrentRate: "9.00",
    Fd4CurrentCompound: "annually",
    Fd4NewCompound: "monthly",
    Fd4Result: "",
    //
    Fd5Interest: "14.87",
    Fd5Compound: "annually",
    Fd5Factor: "2.0",
    Fd5Result: "",
    //
    Fd6Values: "2.0;1.5",
    Fd6Result: [2]string { misc_notes[4], "" },
    //
    Fd7Time: "1.0",
    Fd7TimePeriod: "year",
    Fd7Rate: "15.0",
    Fd7Compound: "annually",
    Fd7PV: "1.00",
    Fd7Result: "",
  }
}

func getMiscellaneousFields(userName string) *miscellaneousFields {
  return currentFields[userName].miscellaneous
}




type oaCpFields struct {
  MenuPage string `json:"menuPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1PV string `json:"fd1PV"`
  Fd1FV string `json:"fd1FV"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2Interest string `json:"fd2Interest"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2Payment string `json:"fd2Payment"`
  Fd2PV string `json:"fd2PV"`
  Fd2Result string `json:"fd2Result"`
  //
  Fd3Interest string `json:"fd3Interest"`
  Fd3Compound string `json:"fd3Compound"`
  Fd3Payment string `json:"fd3Payment"`
  Fd3FV string `json:"fd3FV"`
  Fd3Result string `json:"fd3Result"`
}

func getOaCpFields(userName string) *oaCpFields {
  return currentFields[userName].oaCp
}

func newOaCpFields(dir1, dir2, correlationId string) *oaCpFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "oacp.txt")
  if obj != nil {
    var o oaCpFields
    err := json.Unmarshal(obj, &o)
    if err != nil {
      //Write error, but continue with default values.
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    } else {
      return &o
    }
  } else if err != nil {
    logger.LogError(fmt.Sprintf("%+v", err), correlationId)
  } else {
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "oacp.txt"), correlationId)
  }
  return &oaCpFields {
    MenuPage: "",
    CurrentButton: "lhs-button1",
    //
    Fd1Interest: "1.00",
    Fd1Compound: "annually",
    Fd1PV: "1.00",
    Fd1FV: "1.00",
    Fd1Result: "",
    //
    Fd2Interest: "1.00",
    Fd2Compound: "annually",
    Fd2Payment: "1.00",
    Fd2PV: "1.00",
    Fd2Result: "",
    //
    Fd3Interest: "1.00",
    Fd3Compound: "annually",
    Fd3Payment: "1.00",
    Fd3FV: "1.00",
    Fd3Result: "",
  }
}

type oaEppFields struct {
  MenuPage string `json:"menuPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1N string `json:"fd1N"`
  Fd1TimePeriod string `json:"fd1TimePeriod"`
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1FV string `json:"fd1FV"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2N string `json:"fd2N"`
  Fd2TimePeriod string `json:"fd2TimePeriod"`
  Fd2Interest string `json:"fd2Interest"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2PV string `json:"fd2PV"`
  Fd2Result string `json:"fd2Result"`
}

func getOaEppFields(userName string) *oaEppFields {
  return currentFields[userName].oaEpp
}

func newOaEppFields(dir1, dir2, correlationId string) *oaEppFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "oaepp.txt")
  if obj != nil {
    var o oaEppFields
    err := json.Unmarshal(obj, &o)
    if err != nil {
      //Write error, but continue with default values.
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    } else {
      return &o
    }
  } else if err != nil {
    logger.LogError(fmt.Sprintf("%+v", err), correlationId)
  } else {
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "oaepp.txt"), correlationId)
  }
  return &oaEppFields {
    MenuPage: "",
    CurrentButton: "lhs-button1",
    //
    Fd1N: "1.00",
    Fd1TimePeriod: "year",
    Fd1Interest: "1.00",
    Fd1Compound: "annually",
    Fd1FV: "1.00",
    Fd1Result: "",
    //
    Fd2N: "1.00",
    Fd2TimePeriod: "year",
    Fd2Interest: "1.00",
    Fd2Compound: "annually",
    Fd2PV: "1.00",
    Fd2Result: "",
  }
}

type oaFvFields struct {
  MenuPage string `json:"menuPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1N string `json:"fd1N"`
  Fd1TimePeriod string `json:"fd1TimePeriod"`
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1FV string `json:"fd1FV"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2N string `json:"fd2N"`
  Fd2TimePeriod string `json:"fd2TimePeriod"`
  Fd2Interest string `json:"fd2Interest"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2PMT string `json:"fd2PMT"`
  Fd2Result string `json:"fd2Result"`
}

func newOaFvFields(dir1, dir2, correlationId string) *oaFvFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "oafv.txt")
  if obj != nil {
    var o oaFvFields
    err := json.Unmarshal(obj, &o)
    if err != nil {
      //Write error, but continue with default values.
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    } else {
      return &o
    }
  } else if err != nil {
    logger.LogError(fmt.Sprintf("%+v", err), correlationId)
  } else {
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "oafv.txt"), correlationId)
  }
  return &oaFvFields {
    MenuPage: "",
    CurrentButton: "lhs-button1",
    //
    Fd1N: "1.0",
    Fd1TimePeriod: "year",
    Fd1Interest: "1.00",
    Fd1Compound: "monthly",
    Fd1FV: "1.00",
    Fd1Result: "",
    //
    Fd2N: "1.0",
    Fd2TimePeriod: "year",
    Fd2Interest: "1.00",
    Fd2Compound: "monthly",
    Fd2PMT: "1.00",
    Fd2Result: "",
  }
}

func getOaFvFields(userName string) *oaFvFields {
  return currentFields[userName].oaFv
}

type oaGaFields struct {
  MenuPage string `json:"menuPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1N string `json:"fd1N"`
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1Grow string `json:"fd1Grow"`
  Fd1Pmt string `json:"fd1Pmt"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2N string `json:"fd2N"`
  Fd2Interest string `json:"fd2Interest"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2Grow string `json:"fd2Grow"`
  Fd2Pmt string `json:"fd2Pmt"`
  Fd2Result string `json:"fd2Result"`
}

func newOaGaFields(dir1, dir2, correlationId string) *oaGaFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "oaga.txt")
  if obj != nil {
    var o oaGaFields
    err := json.Unmarshal(obj, &o)
    if err != nil {
      //Write error, but continue with default values.
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    } else {
      return &o
    }
  } else if err != nil {
    logger.LogError(fmt.Sprintf("%+v", err), correlationId)
  } else {
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "oaga.txt"), correlationId)
  }
  return &oaGaFields {
    MenuPage: "",
    CurrentButton: "lhs-button1",
    //
    Fd1N: "1.00",
    Fd1Interest: "1.00",
    Fd1Compound: "annually",
    Fd1Grow: "1.00",
    Fd1Pmt: "1.00",
    Fd1Result: "",
    //
    Fd2N: "1.00",
    Fd2Interest: "1.00",
    Fd2Compound: "annually",
    Fd2Grow: "1.00",
    Fd2Pmt: "1.00",
    Fd2Result: "",
  }
}

func getOaGaFields(userName string) *oaGaFields {
  return currentFields[userName].oaGa
}

type oaInterestRateFields struct {
  MenuPage string `json:"menuPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1N string `json:"fd1N"`
  Fd1TimePeriod string `json:"fd1TimePeriod"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1PV string `json:"fd1PV"`
  Fd1FV string `json:"fd1FV"`
  Fd1Result string `json:"fd1Result"`
}

func newOaInterestRateFields(dir1, dir2, correlationId string) *oaInterestRateFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "oainterestrate.txt")
  if obj != nil {
    var o oaInterestRateFields
    err := json.Unmarshal(obj, &o)
    if err != nil {
      //Write error, but continue with default values.
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    } else {
      return &o
    }
  } else if err != nil {
    logger.LogError(fmt.Sprintf("%+v", err), correlationId)
  } else {
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "oainterestrate.txt"),
      correlationId)
  }
  return &oaInterestRateFields {
    MenuPage: "",
    CurrentButton: "lhs-button1",
    //
    Fd1N: "1.0",
    Fd1TimePeriod: "year",
    Fd1Compound: "monthly",
    Fd1PV: "1.00",
    Fd1FV: "1.07",
    Fd1Result: "",
  }
}

func getOaInterestRateFields(userName string) *oaInterestRateFields {
  return currentFields[userName].oaInterestRate
}

type oaPerpetuityFields struct {
  MenuPage string `json:"menuPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1Pmt string `json:"fd1Pmt"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2Interest string `json:"fd2Interest"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2Grow string `json:"fd2Grow"`
  Fd2Pmt string `json:"fd2Pmt"`
  Fd2Result string `json:"fd2Result"`
}

func newOaPerpetuityFields(dir1, dir2, correlationId string) *oaPerpetuityFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "oaperpetuity.txt")
  if obj != nil {
    var o oaPerpetuityFields
    err := json.Unmarshal(obj, &o)
    if err != nil {
      //Write error, but continue with default values.
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    } else {
      return &o
    }
  } else if err != nil {
    logger.LogError(fmt.Sprintf("%+v", err), correlationId)
  } else {
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "oaperpetuity.txt"), correlationId)
  }
  return &oaPerpetuityFields {
    MenuPage: "",
    CurrentButton: "lhs-button1",
    //
    Fd1Interest: "1.00",
    Fd1Compound: "annually",
    Fd1Pmt: "1.00",
    Fd1Result: "",
    //
    Fd2Interest: "1.00",
    Fd2Compound: "annually",
    Fd2Grow: "1.00",
    Fd2Pmt: "1.00",
    Fd2Result: "",
  }
}

func getOaPerpetuityFields(userName string) *oaPerpetuityFields {
  return currentFields[userName].oaPerpetuity
}

type oaPvFields struct {
  MenuPage string `json:"menuPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1N string `json:"fd1N"`
  Fd1TimePeriod string `json:"fd1TimePeriod"`
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1FV string `json:"fd1FV"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2N string `json:"fd2N"`
  Fd2TimePeriod string `json:"fd2TimePeriod"`
  Fd2Interest string `json:"fd2Interest"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2PMT string `json:"fd2PMT"`
  Fd2Result string `json:"fd2Result"`
}

func newOaPvFields(dir1, dir2, correlationId string) *oaPvFields {
  dir, err := osu.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "oapv.txt")
  if obj != nil {
    var o oaPvFields
    err := json.Unmarshal(obj, &o)
    if err != nil {
      //Write error, but continue with default values.
      logger.LogInfo(fmt.Sprintf("%+v", err), correlationId)
    } else {
      return &o
    }
  } else if err != nil {
    logger.LogError(fmt.Sprintf("%+v", err), correlationId)
  } else {
    logger.LogInfo(fmt.Sprintf("File %s does not exit.", dir + "oapv.txt"), correlationId)
  }
  return &oaPvFields {
    MenuPage: "",
    CurrentButton: "lhs-button1",
    //
    Fd1N: "1.0",
    Fd1TimePeriod: "year",
    Fd1Interest: "1.00",
    Fd1Compound: "monthly",
    Fd1FV: "1.00",
    Fd1Result: "",
    //
    Fd2N: "1.0",
    Fd2TimePeriod: "year",
    Fd2Interest: "1.00",
    Fd2Compound: "monthly",
    Fd2PMT: "1.00",
    Fd2Result: "",
  }
}

func getOaPvFields(userName string) *oaPvFields {
  return currentFields[userName].oaPv
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
