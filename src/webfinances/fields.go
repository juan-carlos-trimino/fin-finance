package webfinances

import (
	"encoding/json"
	"errors"
	"finance/misc"
	"fmt"
	"os"
)

const (
  mainDir string = "./dataDir/fields"
)

var mt = misc.Misc{}

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
  currentPage string
  currentButton string
  //
  fd1Nominal string
  fd1Compound string
  fd1Result [2]string
  //
  fd2Effective string
  fd2Compound string
  fd2Result [3]string
  //
  fd3Nominal string
  fd3Inflation string
  fd3Result [4]string
  //
  fd4Interest string
  fd4Compound string
  fd4Factor string
  fd4Result string
  //
  fd5Values string
  fd5Result [2]string
  //
  fd6Time string
  fd6TimePeriod string
  fd6Rate string
  fd6Compound string
  fd6PV string
  fd6Result string
}

func newMiscellaneousFields() *miscellaneousFields {
  return &miscellaneousFields {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1Nominal: "3.5",
    fd1Compound: "monthly",
    fd1Result: [2]string { misc_notes[0], "" },
    //
    fd2Effective: "3.5",
    fd2Compound: "monthly",
    fd2Result: [3]string { misc_notes[0], misc_notes[1], "" },
    //
    fd3Nominal: "2.0",
    fd3Inflation: "2.0",
    fd3Result: [4]string { misc_notes[1], misc_notes[2], misc_notes[3], "" },
    //
    fd4Interest: "14.87",
    fd4Compound: "annually",
    fd4Factor: "2.0",
    fd4Result: "",
    //
    fd5Values: "2.0;1.5",
    fd5Result: [2]string { misc_notes[4], "" },
    //
    fd6Time: "1.0",
    fd6TimePeriod: "year",
    fd6Rate: "15.0",
    fd6Compound: "annually",
    fd6PV: "1.00",
    fd6Result: "",
  }
}

func getMiscellaneousFields(userName string) *miscellaneousFields {
  return currentFields[userName].miscellaneous
}

type mortgageFields struct {
  CurrentPage string `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1N string `json:"fd1N"`
  Fd1TimePeriod string `json:"fd1TimePeriod"`
  Fd1Interest string `json:"fd1Interest"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1Amount string `json:"fd1Amount"`
  Fd1Result [3]string `json:"fd1Result"`
  //
  Fd2N string `json:"Fd2N"`
  Fd2TimePeriod string `json:"Fd2TimePeriod"`
  Fd2Interest string `json:"fd2Interest"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2Amount string `json:"fd2Amount"`
  Fd2TotalCost string `json:"fd2TotalCost"`
  Fd2TotalInterest string `json:"fd2TotalInterest"`
  Fd2Result []Row `json:"fd2Result"`
  //
  Fd3Mrate string `json:"fd3Mrate"`
  Fd3Mbalance string `json:"fd3Mbalance"`
  Fd3Hrate string `json:"fd3Hrate"`
  Fd3Hbalance string `json:"fd3Hbalance"`
  Fd3Result [3]string `json:"fd3Result"`
}

func newMortgageFields(dir1, dir2 string) *mortgageFields {
  dir, err := misc.CreateDirs(0o017, 0o770, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "mortgage.txt")
  if obj != nil {
    var m mortgageFields
    err := json.Unmarshal(obj, &m)
    if err != nil {
      //Write error, but continue with default values.
      fmt.Printf("%s - %+v\n", mt.DTF(), err)
    } else {
      return &m
    }
  } else {
    fmt.Printf("%s - %+v\n", mt.DTF(), err)
  }
  return &mortgageFields {
    CurrentPage: "rhs-ui1",
    CurrentButton: "lhs-button1",
    //
    Fd1N: "30.0",
    Fd1TimePeriod: "year",
    Fd1Interest: "7.50",
    Fd1Compound: "monthly",
    Fd1Amount: "100000.00",
    Fd1Result: [3]string { "", "", "" },
    //
    Fd2N: "30.0",
    Fd2TimePeriod: "year",
    Fd2Interest: "3.00",
    Fd2Compound: "monthly",
    Fd2Amount: "100000.00",
    Fd2TotalCost: "",
    Fd2TotalInterest: "",
    Fd2Result: []Row{},
    //
    Fd3Mrate: "3.375",
    Fd3Mbalance: "300000.00",
    Fd3Hrate: "2.875",
    Fd3Hbalance: "100000.00",
    Fd3Result: [3]string { mortgage_notes[0], mortgage_notes[1], "" },
  }
}

func getMortgageFields(userName string) *mortgageFields {
  return currentFields[userName].mortgage
}

type bondsFields struct {
  CurrentPage string `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1TaxFree string `json:"fd1TaxFree"`
  Fd1CityTax string `json:"fd1CityTax"`
  Fd1StateTax string `json:"fd1StateTax"`
  Fd1FederalTax string `json:"fd1FederalTax"`
  Fd1Result string `json:"fd1Result"`
  //
  Fd2FaceValue string `json:"fd2FaceValue"`
  Fd2Time string `json:"fd2Time"`
  Fd2TimePeriod string `json:"Fd2TimePeriod"`
  Fd2Coupon string `json:"fd2Coupon"`
  Fd2Current string `json:"fd2Current"`
  Fd2Compound string `json:"fd2Compound"`
  Fd2Result string `json:"fd2Result"`
  //
  Fd3FaceValue string `json:"fd3FaceValue"`
  Fd3TimeCall string `json:"fd3TimeCall"`
  Fd3TimePeriod string `json:"fd3TimePeriod"`
  Fd3Coupon string `json:"fd3Coupon"`
  Fd3Compound string `json:"fd3Compound"`
  Fd3BondPrice string `json:"fd3BondPrice"`
  Fd3CallPrice string `json:"fd3CallPrice"`
  Fd3Result string `json:"fd3Result"`
  //
  Fd4FaceValue string `json:"fd4FaceValue"`
  Fd4Time string `json:"fd4Time"`
  Fd4TimePeriod string `json:"fd4TimePeriod"`
  Fd4Coupon string `json:"fd4Coupon"`
  Fd4Compound string `json:"fd4Compound"`
  Fd4CurrentRadio string `json:"fd4CurrentRadio"`
  Fd4CurInterest string `json:"fd4CurInterest"`
  Fd4BondPrice string `json:"fd4BondPrice"`
  Fd4Result string `json:"fd4Result"`
  //
  Fd5FaceValue string `json:"fd5FaceValue"`
  Fd5Time string `json:"fd5Time"`
  Fd5TimePeriod string `json:"fd5TimePeriod"`
  Fd5Coupon string `json:"fd5Coupon"`
  Fd5CurInterest string `json:"fd5CurInterest"`
  Fd5Compound string `json:"fd5Compound"`
  Fd5Result string `json:"fd5Result"`
  //
  Fd6FaceValue string `json:"fd6FaceValue"`
  Fd6Time string `json:"fd6Time"`
  Fd6TimePeriod string `json:"fd6TimePeriod"`
  Fd6Coupon string `json:"fd6Coupon"`
  Fd6CurInterest string `json:"fd6CurInterest"`
  Fd6Compound string `json:"fd6Compound"`
  Fd6Result [2]string `json:"fd6Result"`
  //
  Fd7FaceValue string `json:"fd7FaceValue"`
  Fd7Time string `json:"fd7Time"`
  Fd7TimePeriod string `json:"fd7TimePeriod"`
  Fd7Coupon string `json:"fd7Coupon"`
  Fd7CurInterest string `json:"fd7CurInterest"`
  Fd7Compound string `json:"fd7Compound"`
  Fd7Result [2]string `json:"fd7Result"`
  //
  Fd8FaceValue string `json:"fd8FaceValue"`
  Fd8Time string `json:"fd8Time"`
  Fd8TimePeriod string `json:"fd8TimePeriod"`
  Fd8Coupon string `json:"fd8Coupon"`
  Fd8CurInterest string `json:"fd8CurInterest"`
  Fd8Compound string `json:"fd8Compound"`
  Fd8Result [2]string `json:"fd8Result"`
}

func newBondsFields(dir1, dir2 string) *bondsFields {
  dir, err := misc.CreateDirs(0o017, 0o770, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "bonds.txt")
  if obj != nil {
    var b bondsFields
    err := json.Unmarshal(obj, &b)
    if err != nil {
      //Write error, but continue with default values.
      fmt.Printf("%s - %+v\n", mt.DTF(), err)
    } else {
      return &b
    }
  } else {
    fmt.Printf("%s - %+v\n", mt.DTF(), err)
  }
  return &bondsFields {
    CurrentPage: "rhs-ui1",
    CurrentButton: "lhs-button1",
    //
    Fd1TaxFree: "3.5",
    Fd1CityTax: "0.0",
    Fd1StateTax: "1.0",
    Fd1FederalTax: "23.0",
    Fd1Result: "",
    //
    Fd2FaceValue: "1000.00",
    Fd2Time: "5",
    Fd2TimePeriod: "year",
    Fd2Coupon: "3.00",
    Fd2Current: "3.5",
    Fd2Compound: "semiannually",
    Fd2Result: "",
    //
    Fd3FaceValue: "1000.00",
    Fd3TimeCall: "2",
    Fd3TimePeriod: "year",
    Fd3Coupon: "2.0",
    Fd3Compound: "semiannually",
    Fd3BondPrice: "990.00",
    Fd3CallPrice: "1050.00",
    Fd3Result: "",
    //
    Fd4FaceValue: "1000.00",
    Fd4Time: "3",
    Fd4TimePeriod: "year",
    Fd4Coupon: "2.5",
    Fd4Compound: "semiannually",
    Fd4CurrentRadio: "fd4-curinterest",
    Fd4CurInterest: "2.3",
    Fd4BondPrice: "1000.00",
    Fd4Result: "",
    //
    Fd5FaceValue: "1000.00",
    Fd5Time: "5",
    Fd5TimePeriod: "year",
    Fd5Coupon: "5.4",
    Fd5CurInterest: "7.5",
    Fd5Compound: "semiannually",
    Fd5Result: "",
    //
    Fd6FaceValue: "1000.00",
    Fd6Time: "5",
    Fd6TimePeriod: "year",
    Fd6Coupon: "5.4",
    Fd6CurInterest: "7.5",
    Fd6Compound: "semiannually",
    Fd6Result: [2]string { bond_notes[1], "" },
    //
    Fd7FaceValue: "1000.00",
    Fd7Time: "5",
    Fd7TimePeriod: "year",
    Fd7Coupon: "5.4",
    Fd7CurInterest: "7.5",
    Fd7Compound: "semiannually",
    Fd7Result: [2]string { bond_notes[0], "" },
    //
    Fd8FaceValue: "1000.00",
    Fd8Time: "5",
    Fd8TimePeriod: "year",
    Fd8Coupon: "5.4",
    Fd8CurInterest: "7.5",
    Fd8Compound: "semiannually",
    Fd8Result: [2]string { bond_notes[2], "" },
  }
}

func getBondsFields(userName string) *bondsFields {
  return currentFields[userName].bonds
}

type adFvFields struct {
  CurrentPage string `json:"currentPage"`
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

func newAdFvFields(dir1, dir2 string) *adFvFields {
  dir, err := misc.CreateDirs(0o017, 0o770, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "adfv.txt")
  if obj != nil {
    var a adFvFields
    err := json.Unmarshal(obj, &a)
    if err != nil {
      //Write error, but continue with default values.
      fmt.Printf("%s - %+v\n", mt.DTF(), err)
    } else {
      return &a
    }
  } else {
    fmt.Printf("%s - %+v\n", mt.DTF(), err)
  }
  return &adFvFields {
    CurrentPage: "rhs-ui2",
    CurrentButton: "lhs-button2",
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

func getAdFvFields(userName string) *adFvFields {
  return currentFields[userName].adFv
}

type adPvFields struct {
  CurrentPage string `json:"currentPage"`
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

func newAdPvFields(dir1, dir2 string) *adPvFields {
  dir, err := misc.CreateDirs(0o017, 0o770, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "adpv.txt")
  if obj != nil {
    var a adPvFields
    err := json.Unmarshal(obj, &a)
    if err != nil {
      //Write error, but continue with default values.
      fmt.Printf("%s - %+v\n", mt.DTF(), err)
    } else {
      return &a
    }
  } else {
    fmt.Printf("%s - %+v\n", mt.DTF(), err)
  }
  return &adPvFields {
    CurrentPage: "rhs-ui2",
    CurrentButton: "lhs-button2",
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

func getAdPvFields(userName string) *adPvFields {
  return currentFields[userName].adPv
}

type adCpFields struct {
  CurrentPage string `json:"currentPage"`
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

func newAdCpFields(dir1, dir2 string) *adCpFields {
  dir, err := misc.CreateDirs(0o017, 0o770, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "adcp.txt")
  if obj != nil {
    var a adCpFields
    err := json.Unmarshal(obj, &a)
    if err != nil {
      //Write error, but continue with default values.
      fmt.Printf("%s - %+v\n", mt.DTF(), err)
    } else {
      return &a
    }
  } else {
    fmt.Printf("%s - %+v\n", mt.DTF(), err)
  }
  return &adCpFields {
    CurrentPage: "rhs-ui2",
    CurrentButton: "lhs-button2",
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

func getAdCpFields(userName string) *adCpFields {
  return currentFields[userName].adCp
}

type adEppFields struct {
  CurrentPage string `json:"currentPage"`
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

func newAdEppFields(dir1, dir2 string) *adEppFields {
  dir, err := misc.CreateDirs(0o017, 0o770, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "adepp.txt")
  if obj != nil {
    var a adEppFields
    err := json.Unmarshal(obj, &a)
    if err != nil {
      //Write error, but continue with default values.
      fmt.Printf("%s - %+v\n", mt.DTF(), err)
    } else {
      return &a
    }
  } else {
    fmt.Printf("%s - %+v\n", mt.DTF(), err)
  }
  return &adEppFields {
    CurrentPage: "rhs-ui1",
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

func getAdEppFields(userName string) *adEppFields {
  return currentFields[userName].adEpp
}

type oaCpFields struct {
  CurrentPage string `json:"currentPage"`
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

func newOaCpFields(dir1, dir2 string) *oaCpFields {
  dir, err := misc.CreateDirs(0o017, 0o770, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "oacp.txt")
  if obj != nil {
    var o oaCpFields
    err := json.Unmarshal(obj, &o)
    if err != nil {
      //Write error, but continue with default values.
      fmt.Printf("%s - %+v\n", mt.DTF(), err)
    } else {
      return &o
    }
  } else {
    fmt.Printf("%s - %+v\n", mt.DTF(), err)
  }
  return &oaCpFields {
    CurrentPage: "rhs-ui1",
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
  CurrentPage string `json:"currentPage"`
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

func newOaEppFields(dir1, dir2 string) *oaEppFields {
  dir, err := misc.CreateDirs(0o017, 0o770, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "oaepp.txt")
  if obj != nil {
    var o oaEppFields
    err := json.Unmarshal(obj, &o)
    if err != nil {
      //Write error, but continue with default values.
      fmt.Printf("%s - %+v\n", mt.DTF(), err)
    } else {
      return &o
    }
  } else {
    fmt.Printf("%s - %+v\n", mt.DTF(), err)
  }
  return &oaEppFields {
    CurrentPage: "rhs-ui1",
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
  currentPage string
  currentButton string
  //
  fd1N string
  fd1TimePeriod string
  fd1Interest string
  fd1Compound string
  fd1FV string
  fd1Result string
  //
  fd2N string
  fd2TimePeriod string
  fd2Interest string
  fd2Compound string
  fd2PMT string
  fd2Result string
}

func newOaFvFields() *oaFvFields {
  return &oaFvFields {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1N: "1.0",
    fd1TimePeriod: "year",
    fd1Interest: "1.00",
    fd1Compound: "monthly",
    fd1FV: "1.00",
    fd1Result: "",
    //
    fd2N: "1.0",
    fd2TimePeriod: "year",
    fd2Interest: "1.00",
    fd2Compound: "monthly",
    fd2PMT: "1.00",
    fd2Result: "",
  }
}

func getOaFvFields(userName string) *oaFvFields {
  return currentFields[userName].oaFv
}

type oaGaFields struct {
  currentPage string
  currentButton string
  //
  fd1N string
  fd1Interest string
  fd1Compound string
  fd1Grow string
  fd1Pmt string
  fd1Result string
  //
  fd2N string
  fd2Interest string
  fd2Compound string
  fd2Grow string
  fd2Pmt string
  fd2Result string
}

func newOaGaFields() *oaGaFields {
  return &oaGaFields {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1N: "1.00",
    fd1Interest: "1.00",
    fd1Compound: "annually",
    fd1Grow: "1.00",
    fd1Pmt: "1.00",
    fd1Result: "",
    //
    fd2N: "1.00",
    fd2Interest: "1.00",
    fd2Compound: "annually",
    fd2Grow: "1.00",
    fd2Pmt: "1.00",
    fd2Result: "",
  }
}

func getOaGaFields(userName string) *oaGaFields {
  return currentFields[userName].oaGa
}

type oaInterestRateFields struct {
  currentPage string
  currentButton string
  //
  fd1N string
  fd1TimePeriod string
  fd1Compound string
  fd1PV string
  fd1FV string
  fd1Result string
}

func newOaInterestRateFields() *oaInterestRateFields {
  return &oaInterestRateFields {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1N: "1.0",
    fd1TimePeriod: "year",
    fd1Compound: "monthly",
    fd1PV: "1.00",
    fd1FV: "1.07",
    fd1Result: "",
  }
}

func getOaInterestRateFields(userName string) *oaInterestRateFields {
  return currentFields[userName].oaInterestRate
}

type oaPerpetuityFields struct {
  currentPage string
  currentButton string
  //
  fd1Interest string
  fd1Compound string
  fd1Pmt string
  fd1Result string
  //
  fd2Interest string
  fd2Compound string
  fd2Grow string
  fd2Pmt string
  fd2Result string
}

func newOaPerpetuityFields() *oaPerpetuityFields {
  return &oaPerpetuityFields {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1Interest: "1.00",
    fd1Compound: "annually",
    fd1Pmt: "1.00",
    fd1Result: "",
    //
    fd2Interest: "1.00",
    fd2Compound: "annually",
    fd2Grow: "1.00",
    fd2Pmt: "1.00",
    fd2Result: "",
  }
}

func getOaPerpetuityFields(userName string) *oaPerpetuityFields {
  return currentFields[userName].oaPerpetuity
}

type oaPvFields struct {
  currentPage string
  currentButton string
  //
  fd1N string
  fd1TimePeriod string
  fd1Interest string
  fd1Compound string
  fd1FV string
  fd1Result string
  //
  fd2N string
  fd2TimePeriod string
  fd2Interest string
  fd2Compound string
  fd2PMT string
  fd2Result string
}

func newOaPvFields() *oaPvFields {
  return &oaPvFields {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1N: "1.0",
    fd1TimePeriod: "year",
    fd1Interest: "1.00",
    fd1Compound: "monthly",
    fd1FV: "1.00",
    fd1Result: "",
    //
    fd2N: "1.0",
    fd2TimePeriod: "year",
    fd2Interest: "1.00",
    fd2Compound: "monthly",
    fd2PMT: "1.00",
    fd2Result: "",
  }
}

func getOaPvFields(userName string) *oaPvFields {
  return currentFields[userName].oaPv
}

type siAccurateFields struct {
  currentPage string
  currentButton string
  //
  fd1Time string
  fd1TimePeriod string
  fd1Interest string
  fd1Compound string
  fd1PV string
  fd1Result string
  //
  fd2Time string
  fd2TimePeriod string
  fd2Amount string
  fd2PV string
  fd2Result string
  //
  fd3Time string
  fd3TimePeriod string
  fd3Interest string
  fd3Compound string
  fd3Amount string
  fd3Result string
  //
  fd4Interest string
  fd4Compound string
  fd4Amount string
  fd4PV string
  fd4Result string
}

func newSiAccurateFields() *siAccurateFields {
  return &siAccurateFields {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1Time: "1",
    fd1TimePeriod: "year",
    fd1Interest: "1.00",
    fd1Compound: "annually",
    fd1PV: "1.00",
    fd1Result: "",
    //
    fd2Time: "1",
    fd2TimePeriod: "year",
    fd2Amount: "1.00",
    fd2PV: "1.00",
    fd2Result: "",
    //
    fd3Time: "1",
    fd3TimePeriod: "year",
    fd3Interest: "1.0",
    fd3Compound: "annually",
    fd3Amount: "1.00",
    fd3Result: "",
    //
    fd4Interest: "1.00",
    fd4Compound: "annually",
    fd4Amount: "1.00",
    fd4PV: "1.00",
    fd4Result: "",
  }
}

func getSiAccurateFields(userName string) *siAccurateFields {
  return currentFields[userName].siAccurate
}

type siBankersFields struct {
  currentPage string
  currentButton string
  //
  fd1Time string
  fd1TimePeriod string
  fd1Interest string
  fd1Compound string
  fd1PV string
  fd1Result string
  //
  fd2Time string
  fd2TimePeriod string
  fd2Amount string
  fd2PV string
  fd2Result string
  //
  fd3Time string
  fd3TimePeriod string
  fd3Interest string
  fd3Compound string
  fd3Amount string
  fd3Result string
  //
  fd4Interest string
  fd4Compound string
  fd4Amount string
  fd4PV string
  fd4Result string
}

func newSiBankersFields() *siBankersFields {
  return &siBankersFields {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1Time: "1",
    fd1TimePeriod: "year",
    fd1Interest: "1.00",
    fd1Compound: "annually",
    fd1PV: "1.00",
    fd1Result: "",
    //
    fd2Time: "1",
    fd2TimePeriod: "year",
    fd2Amount: "1.00",
    fd2PV: "1.00",
    fd2Result: "",
    //
    fd3Time: "1",
    fd3TimePeriod: "year",
    fd3Interest: "1.0",
    fd3Compound: "annually",
    fd3Amount: "1.00",
    fd3Result: "",
    //
    fd4Interest: "1.00",
    fd4Compound: "annually",
    fd4Amount: "1.00",
    fd4PV: "1.00",
    fd4Result: "",
  }
}

func getSiBankersFields(userName string) *siBankersFields {
  return currentFields[userName].siBankers
}

type siOrdinaryFields struct {
  currentPage string
  currentButton string
  //
  fd1Time string
  fd1TimePeriod string
  fd1Interest string
  fd1Compound string
  fd1PV string
  fd1Result string
  //
  fd2Time string
  fd2TimePeriod string
  fd2Amount string
  fd2PV string
  fd2Result string
  //
  fd3Time string
  fd3TimePeriod string
  fd3Interest string
  fd3Compound string
  fd3Amount string
  fd3Result string
  //
  fd4Interest string
  fd4Compound string
  fd4Amount string
  fd4PV string
  fd4Result string
}

func newSiOrdinaryFields() *siOrdinaryFields {
  return &siOrdinaryFields {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1Time: "1",
    fd1TimePeriod: "year",
    fd1Interest: "1.00",
    fd1Compound: "annually",
    fd1PV: "1.00",
    fd1Result: "",
    //
    fd2Time: "1",
    fd2TimePeriod: "year",
    fd2Amount: "1.00",
    fd2PV: "1.00",
    fd2Result: "",
    //
    fd3Time: "1",
    fd3TimePeriod: "year",
    fd3Interest: "1.0",
    fd3Compound: "annually",
    fd3Amount: "1.00",
    fd3Result: "",
    //
    fd4Interest: "1.00",
    fd4Compound: "annually",
    fd4Amount: "1.00",
    fd4PV: "1.00",
    fd4Result: "",
  }
}

func getSiOrdinaryFields(userName string) *siOrdinaryFields {
  return currentFields[userName].siOrdinary
}

func AddSessionDataPerUser(userName string) {
  if _, ok := currentFields[userName]; !ok {
    fd := &fields{
      miscellaneous: newMiscellaneousFields(),
      mortgage: newMortgageFields(mainDir, userName),
      bonds: newBondsFields(mainDir, userName),
      adFv: newAdFvFields(mainDir, userName),
      adPv: newAdPvFields(mainDir, userName),
      adCp: newAdCpFields(mainDir, userName),
      adEpp: newAdEppFields(mainDir, userName),
      oaCp: newOaCpFields(mainDir, userName),
      oaEpp: newOaEppFields(mainDir, userName),
      oaFv: newOaFvFields(mainDir, userName),
      oaGa: newOaGaFields(),
      oaInterestRate: newOaInterestRateFields(),
      oaPerpetuity: newOaPerpetuityFields(),
      oaPv: newOaPvFields(),
      siAccurate: newSiAccurateFields(),
      siBankers: newSiBankersFields(),
      siOrdinary: newSiOrdinaryFields(),
    }
    currentFields[userName] = fd
  }
}

func DeleteSessionDataPerUser(userName string) {
  delete(currentFields, userName)
}

func readFields(filePath string) ([]byte, error) {
  exists, err := misc.CheckFileExists(filePath)
  if exists {
    obj, err := misc.ReadAllShareLock(filePath, os.O_RDONLY, 0o660)
    if err != nil {
      return nil, errors.New(fmt.Sprintf("Couldn't open file %s: ", filePath) + err.Error())
    }
    return obj, nil
  } else if err == nil {
    //File doesn't exist, create it.
    f, err := os.OpenFile(filePath, os.O_CREATE | os.O_RDWR, 0o660)
    if err != nil {
      return nil, errors.New(fmt.Sprintf("Couldn't create file %s: ", filePath) + err.Error())
    }
    f.Close()
    return nil, nil
  } else {
    return nil, err
  }
}
