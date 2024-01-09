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
  Fd2N string `json:"fd2N"`
  Fd2TimePeriod string `json:"fd2TimePeriod"`
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
  obj, err := readFields(dir + "mortgage.bin")
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
  currentPage string
  currentButton string
  //
  fd1TaxFree string
  fd1CityTax string
  fd1StateTax string
  fd1FederalTax string
  fd1Result string
  //
  fd2FaceValue string
  fd2Time string
  fd2TimePeriod string
  fd2Coupon string
  fd2Current string
  fd2Compound string
  fd2Result string
  //
  fd3FaceValue string
  fd3TimeCall string
  fd3TimePeriod string
  fd3Coupon string
  fd3Compound string
  fd3BondPrice string
  fd3CallPrice string
  fd3Result string
  //
  fd4FaceValue string
  fd4Time string
  fd4TimePeriod string
  fd4Coupon string
  fd4Compound string
  fd4CurrentRadio string
  fd4CurInterest string
  fd4BondPrice string
  fd4Result string
  //
  fd5FaceValue string
  fd5Time string
  fd5TimePeriod string
  fd5Coupon string
  fd5CurInterest string
  fd5Compound string
  fd5Result string
  //
  fd6FaceValue string
  fd6Time string
  fd6TimePeriod string
  fd6Coupon string
  fd6CurInterest string
  fd6Compound string
  fd6Result [2]string
  //
  fd7FaceValue string
  fd7Time string
  fd7TimePeriod string
  fd7Coupon string
  fd7CurInterest string
  fd7Compound string
  fd7Result [2]string
  //
  fd8FaceValue string
  fd8Time string
  fd8TimePeriod string
  fd8Coupon string
  fd8CurInterest string
  fd8Compound string
  fd8Result [2]string
}

func newBondsFields() *bondsFields {
  return &bondsFields {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1TaxFree: "3.5",
    fd1CityTax: "0.0",
    fd1StateTax: "1.0",
    fd1FederalTax: "23.0",
    fd1Result: "",
    //
    fd2FaceValue: "1000.00",
    fd2Time: "5",
    fd2TimePeriod: "year",
    fd2Coupon: "3.00",
    fd2Current: "3.5",
    fd2Compound: "semiannually",
    fd2Result: "",
    //
    fd3FaceValue: "1000.00",
    fd3TimeCall: "2",
    fd3TimePeriod: "year",
    fd3Coupon: "2.0",
    fd3Compound: "semiannually",
    fd3BondPrice: "990.00",
    fd3CallPrice: "1050.00",
    fd3Result: "",
    //
    fd4FaceValue: "1000.00",
    fd4Time: "3",
    fd4TimePeriod: "year",
    fd4Coupon: "2.5",
    fd4Compound: "semiannually",
    fd4CurrentRadio: "fd4-curinterest",
    fd4CurInterest: "2.3",
    fd4BondPrice: "1000.00",
    fd4Result: "",
    //
    fd5FaceValue: "1000.00",
    fd5Time: "5",
    fd5TimePeriod: "year",
    fd5Coupon: "5.4",
    fd5CurInterest: "7.5",
    fd5Compound: "semiannually",
    fd5Result: "",
    //
    fd6FaceValue: "1000.00",
    fd6Time: "5",
    fd6TimePeriod: "year",
    fd6Coupon: "5.4",
    fd6CurInterest: "7.5",
    fd6Compound: "semiannually",
    fd6Result: [2]string { bond_notes[1], "" },
    //
    fd7FaceValue: "1000.00",
    fd7Time: "5",
    fd7TimePeriod: "year",
    fd7Coupon: "5.4",
    fd7CurInterest: "7.5",
    fd7Compound: "semiannually",
    fd7Result: [2]string { bond_notes[0], "" },
    //
    fd8FaceValue: "1000.00",
    fd8Time: "5",
    fd8TimePeriod: "year",
    fd8Coupon: "5.4",
    fd8CurInterest: "7.5",
    fd8Compound: "semiannually",
    fd8Result: [2]string { bond_notes[2], "" },
  }
}

func getBondsFields(userName string) *bondsFields {
  return currentFields[userName].bonds
}

type adFvFields struct {
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

func newAdFvFields() *adFvFields {
  return &adFvFields {
    currentPage: "rhs-ui2",
    currentButton: "lhs-button2",
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

func getAdFvFields(userName string) *adFvFields {
  return currentFields[userName].adFv
}

type adPvFields struct {
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

func newAdPvFields() *adPvFields {
  return &adPvFields {
    currentPage: "rhs-ui2",
    currentButton: "lhs-button2",
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

func getAdPvFields(userName string) *adPvFields {
  return currentFields[userName].adPv
}

type adCpFields struct {
  currentPage string
  currentButton string
  //
  fd1Interest string
  fd1Compound string
  fd1PV string
  fd1FV string
  fd1Result string
  //
  fd2Interest string
  fd2Compound string
  fd2Payment string
  fd2PV string
  fd2Result string
  //
  fd3Interest string
  fd3Compound string
  fd3Payment string
  fd3FV string
  fd3Result string
}

func newAdCpFields() *adCpFields {
  return &adCpFields {
    currentPage: "rhs-ui2",
    currentButton: "lhs-button2",
    //
    fd1Interest: "1.00",
    fd1Compound: "annually",
    fd1PV: "1.00",
    fd1FV: "1.00",
    fd1Result: "",
    //
    fd2Interest: "1.00",
    fd2Compound: "annually",
    fd2Payment: "1.00",
    fd2PV: "1.00",
    fd2Result: "",
    //
    fd3Interest: "1.00",
    fd3Compound: "annually",
    fd3Payment: "1.00",
    fd3FV: "1.00",
    fd3Result: "",
  }
}

func getAdCpFields(userName string) *adCpFields {
  return currentFields[userName].adCp
}

type adEppFields struct {
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
  fd2PV string
  fd2Result string
}

func newAdEppFields() *adEppFields {
  return &adEppFields {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1N: "1.00",
    fd1TimePeriod: "year",
    fd1Interest: "1.00",
    fd1Compound: "annually",
    fd1FV: "1.00",
    fd1Result: "",
    //
    fd2N: "1.00",
    fd2TimePeriod: "year",
    fd2Interest: "1.00",
    fd2Compound: "annually",
    fd2PV: "1.00",
    fd2Result: "",
  }
}

func getAdEppFields(userName string) *adEppFields {
  return currentFields[userName].adEpp
}

type oaCpFields struct {
  currentPage string
  currentButton string
  //
  fd1Interest string
  fd1Compound string
  fd1PV string
  fd1FV string
  fd1Result string
  //
  fd2Interest string
  fd2Compound string
  fd2Payment string
  fd2PV string
  fd2Result string
  //
  fd3Interest string
  fd3Compound string
  fd3Payment string
  fd3FV string
  fd3Result string
}

func getOaCpFields(userName string) *oaCpFields {
  return currentFields[userName].oaCp
}

func newOaCpFields() *oaCpFields {
  return &oaCpFields {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1Interest: "1.00",
    fd1Compound: "annually",
    fd1PV: "1.00",
    fd1FV: "1.00",
    fd1Result: "",
    //
    fd2Interest: "1.00",
    fd2Compound: "annually",
    fd2Payment: "1.00",
    fd2PV: "1.00",
    fd2Result: "",
    //
    fd3Interest: "1.00",
    fd3Compound: "annually",
    fd3Payment: "1.00",
    fd3FV: "1.00",
    fd3Result: "",
  }
}

type oaEppFields struct {
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
  fd2PV string
  fd2Result string
}

func getOaEppFields(userName string) *oaEppFields {
  return currentFields[userName].oaEpp
}

func newOaEppFields() *oaEppFields {
  return &oaEppFields {
    currentPage: "rhs-ui1",
    currentButton: "lhs-button1",
    //
    fd1N: "1.00",
    fd1TimePeriod: "year",
    fd1Interest: "1.00",
    fd1Compound: "annually",
    fd1FV: "1.00",
    fd1Result: "",
    //
    fd2N: "1.00",
    fd2TimePeriod: "year",
    fd2Interest: "1.00",
    fd2Compound: "annually",
    fd2PV: "1.00",
    fd2Result: "",
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
      bonds: newBondsFields(),
      adFv: newAdFvFields(),
      adPv: newAdPvFields(),
      adCp: newAdCpFields(),
      adEpp: newAdEppFields(),
      oaCp: newOaCpFields(),
      oaEpp: newOaEppFields(),
      oaFv: newOaFvFields(),
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
