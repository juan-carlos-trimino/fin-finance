package webfinances

//
//To fold all block comments:
//  Ctrl+K and Ctrl+/
//To unfold all block comments:
//  Ctrl+K and Ctrl+J

import (
  "errors"
  "fmt"
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
