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
  CurrentPage string `json:"currentPage"`
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

func newMiscellaneousFields(dir1, dir2 string) *miscellaneousFields {
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "miscellaneous.txt")
  if obj != nil {
    var m miscellaneousFields
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
  return &miscellaneousFields {
    CurrentPage: "rhs-ui1",
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
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
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
  Fd2CompoundCoupon string `json:"fd2CompoundCoupon"`
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
  Fd5CompoundCoupon string `json:"fd5CompoundCoupon"`
  Fd5CurInterest string `json:"fd5CurInterest"`
  Fd5Compound string `json:"fd5Compound"`
  Fd5Result string `json:"fd5Result"`
  //
  Fd6FaceValue string `json:"fd6FaceValue"`
  Fd6Time string `json:"fd6Time"`
  Fd6TimePeriod string `json:"fd6TimePeriod"`
  Fd6Coupon string `json:"fd6Coupon"`
  Fd6CompoundCoupon string `json:"fd6CompoundCoupon"`
  Fd6CurInterest string `json:"fd6CurInterest"`
  Fd6Compound string `json:"fd6Compound"`
  Fd6Result [2]string `json:"fd6Result"`
  //
  Fd7FaceValue string `json:"fd7FaceValue"`
  Fd7Time string `json:"fd7Time"`
  Fd7TimePeriod string `json:"fd7TimePeriod"`
  Fd7Coupon string `json:"fd7Coupon"`
  Fd7CompoundCoupon string `json:"fd7CompoundCoupon"`
  Fd7CurInterest string `json:"fd7CurInterest"`
  Fd7Compound string `json:"fd7Compound"`
  Fd7Result [2]string `json:"fd7Result"`
  //
  Fd8FaceValue string `json:"fd8FaceValue"`
  Fd8Time string `json:"fd8Time"`
  Fd8TimePeriod string `json:"fd8TimePeriod"`
  Fd8Coupon string `json:"fd8Coupon"`
  Fd8CompoundCoupon string `json:"fd8CompoundCoupon"`
  Fd8CurInterest string `json:"fd8CurInterest"`
  Fd8Compound string `json:"fd8Compound"`
  Fd8Result [2]string `json:"fd8Result"`
}

func newBondsFields(dir1, dir2 string) *bondsFields {
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
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
    Fd2CompoundCoupon: "annually",
    Fd2Current: "3.5",
    Fd2Compound: "annually",
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
    Fd5CompoundCoupon: "annually",
    Fd5CurInterest: "7.5",
    Fd5Compound: "annually",
    Fd5Result: "",
    //
    Fd6FaceValue: "1000.00",
    Fd6Time: "5",
    Fd6TimePeriod: "year",
    Fd6Coupon: "5.4",
    Fd6CompoundCoupon: "annually",
    Fd6CurInterest: "7.5",
    Fd6Compound: "annually",
    Fd6Result: [2]string { bond_notes[1], "" },
    //
    Fd7FaceValue: "1000.00",
    Fd7Time: "5",
    Fd7TimePeriod: "year",
    Fd7Coupon: "5.4",
    Fd7CompoundCoupon: "annually",
    Fd7CurInterest: "7.5",
    Fd7Compound: "annually",
    Fd7Result: [2]string { bond_notes[0], "" },
    //
    Fd8FaceValue: "1000.00",
    Fd8Time: "5",
    Fd8TimePeriod: "year",
    Fd8Coupon: "5.4",
    Fd8CompoundCoupon: "annually",
    Fd8CurInterest: "7.5",
    Fd8Compound: "annually",
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
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
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
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
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
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
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
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
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
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
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
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
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

func newOaFvFields(dir1, dir2 string) *oaFvFields {
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "oafv.txt")
  if obj != nil {
    var o oaFvFields
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
  return &oaFvFields {
    CurrentPage: "rhs-ui1",
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
  CurrentPage string `json:"currentPage"`
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

func newOaGaFields(dir1, dir2 string) *oaGaFields {
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "oaga.txt")
  if obj != nil {
    var o oaGaFields
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
  return &oaGaFields {
    CurrentPage: "rhs-ui1",
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
  CurrentPage string `json:"currentPage"`
  CurrentButton string `json:"currentButton"`
  //
  Fd1N string `json:"fd1N"`
  Fd1TimePeriod string `json:"fd1TimePeriod"`
  Fd1Compound string `json:"fd1Compound"`
  Fd1PV string `json:"fd1PV"`
  Fd1FV string `json:"fd1FV"`
  Fd1Result string `json:"fd1Result"`
}

func newOaInterestRateFields(dir1, dir2 string) *oaInterestRateFields {
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "oainterestrate.txt")
  if obj != nil {
    var o oaInterestRateFields
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
  return &oaInterestRateFields {
    CurrentPage: "rhs-ui1",
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
  CurrentPage string `json:"currentPage"`
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

func newOaPerpetuityFields(dir1, dir2 string) *oaPerpetuityFields {
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "oaperpetuity.txt")
  if obj != nil {
    var o oaPerpetuityFields
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
  return &oaPerpetuityFields {
    CurrentPage: "rhs-ui1",
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

func newOaPvFields(dir1, dir2 string) *oaPvFields {
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "oapv.txt")
  if obj != nil {
    var o oaPvFields
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
  return &oaPvFields {
    CurrentPage: "rhs-ui1",
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
  CurrentPage string `json:"currentPage"`
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

func newSiAccurateFields(dir1, dir2 string) *siAccurateFields {
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "siaccurate.txt")
  if obj != nil {
    var s siAccurateFields
    err := json.Unmarshal(obj, &s)
    if err != nil {
      //Write error, but continue with default values.
      fmt.Printf("%s - %+v\n", mt.DTF(), err)
    } else {
      return &s
    }
  } else {
    fmt.Printf("%s - %+v\n", mt.DTF(), err)
  }
  return &siAccurateFields {
    CurrentPage: "rhs-ui1",
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

func getSiAccurateFields(userName string) *siAccurateFields {
  return currentFields[userName].siAccurate
}

type siBankersFields struct {
  CurrentPage string `json:"currentPage"`
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

func newSiBankersFields(dir1, dir2 string) *siBankersFields {
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "sibankers.txt")
  if obj != nil {
    var s siBankersFields
    err := json.Unmarshal(obj, &s)
    if err != nil {
      //Write error, but continue with default values.
      fmt.Printf("%s - %+v\n", mt.DTF(), err)
    } else {
      return &s
    }
  } else {
    fmt.Printf("%s - %+v\n", mt.DTF(), err)
  }
  return &siBankersFields {
    CurrentPage: "rhs-ui1",
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
  CurrentPage string `json:"currentPage"`
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

func newSiOrdinaryFields(dir1, dir2 string) *siOrdinaryFields {
  dir, err := misc.CreateDirs(0o077, 0o777, dir1, dir2)
  if err != nil {
    panic("Cannot create directory '" + dir + "': " + err.Error())
  }
  obj, err := readFields(dir + "siordinary.txt")
  if obj != nil {
    var s siOrdinaryFields
    err := json.Unmarshal(obj, &s)
    if err != nil {
      //Write error, but continue with default values.
      fmt.Printf("%s - %+v\n", mt.DTF(), err)
    } else {
      return &s
    }
  } else {
    fmt.Printf("%s - %+v\n", mt.DTF(), err)
  }
  return &siOrdinaryFields {
    CurrentPage: "rhs-ui1",
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

func AddSessionDataPerUser(userName string) {
  if _, ok := currentFields[userName]; !ok {
    fd := &fields{
      miscellaneous: newMiscellaneousFields(mainDir, userName),
      mortgage: newMortgageFields(mainDir, userName),
      bonds: newBondsFields(mainDir, userName),
      adFv: newAdFvFields(mainDir, userName),
      adPv: newAdPvFields(mainDir, userName),
      adCp: newAdCpFields(mainDir, userName),
      adEpp: newAdEppFields(mainDir, userName),
      oaCp: newOaCpFields(mainDir, userName),
      oaEpp: newOaEppFields(mainDir, userName),
      oaFv: newOaFvFields(mainDir, userName),
      oaGa: newOaGaFields(mainDir, userName),
      oaInterestRate: newOaInterestRateFields(mainDir, userName),
      oaPerpetuity: newOaPerpetuityFields(mainDir, userName),
      oaPv: newOaPvFields(mainDir, userName),
      siAccurate: newSiAccurateFields(mainDir, userName),
      siBankers: newSiBankersFields(mainDir, userName),
      siOrdinary: newSiOrdinaryFields(mainDir, userName),
    }
    currentFields[userName] = fd
  }
}

func DeleteSessionDataPerUser(userName string) {
  delete(currentFields, userName)
}

/***
Even though each username has its own exclusive set of files, multiple users can share the same
username. In this scenario, there is a possibility that two or more users may try to modify a file
at the same time thereby corrupting the file. In order to prevent file corruption, the file is
protected with a lock that allows a single writer (exclusive write) or multiple readers (share
reads).
***/
func readFields(filePath string) ([]byte, error) {
  exists, err := misc.CheckFileExists(filePath)
  if exists {
    obj, err := misc.ReadAllShareLock1(filePath, os.O_RDONLY, 0o400)
    if err != nil {
      return nil, errors.New(fmt.Sprintf("Couldn't open file %s: ", filePath) + err.Error())
    }
    return obj, nil
  } else {
    return nil, err
  }
}
