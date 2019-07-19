package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	sc "github.com/hyperledger/fabric/protos/peer"
)

//Smart Contract struct
type smartContract struct {
}

//Add Event struct
type notifEvent struct {
	AssetUniqueID string `json:"assetUniqueId"`
	Function      string `json:"function"`
	Invoker       string `json:"invoker"`
	Status        string `json:"status"`
}

// Fixed Asset Schema
type fixedAsset struct {
	AssetUniqueID   string   `json:"assetUniqueId"`
	AssetNum        string   `json:"assetNum"`
	AssetTag        string   `json:"assetTag"`
	OrgID           string   `json:"orgId"`
	SiteID          string   `json:"siteId"`
	StatusDate      string   `json:"statusDate"`
	SerialNum       string   `json:"serialNum"`
	Model           string   `json:"model"`
	PurchasePrice   float64  `json:"purchasePrice"`
	InvoicePrice    float64  `json:"invoicePrice"`
	NetBookValue    float64  `json:"netBookValue"`
	Owner           string   `json:"owner"`
	Vendor          string   `json:"vendor"`
	Manufacturer    string   `json:"manufacturer"`
	Value           float64  `json:"value"`
	Location        string   `json:"location"`
	AcquisitionDate string   `json:"acquisitionDate"`
	InstallDate     string   `json:"installDate"`
	InvoiceDate     []string `json:"invoiceDate"` //Dates for multiple invoices
	RetireDate      string   `json:"retireDate"`
	Status          string   `json:"status"`
	Description     string   `json:"description"`
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *smartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {

	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *smartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	//var err error
	if fn == "receiveAsset" {
		return t.receiveAsset(stub, args)
	} else if fn == "installAsset" {
		return t.installAsset(stub, args)
	} else if fn == "receiveInvoice" {
		return t.receiveInvoice(stub, args)
	} else if fn == "faUpdate" {
		return t.faUpdate(stub, args)
	} else if fn == "submitPO" {
		return t.submitPO(stub, args)
	}
	return shim.Error("Invalid function name")
}

// Data entry and validation when asset is received at IBX.
/*
Args List:
1. Asset Unique Identifier
2. Description
3. Acquisition Date
4. Location
5. Make
6. Vendor
7. Model
8. Serial Number
9. Org ID
*/
func (t *smartContract) receiveAsset(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 9 { //For now this is okay, need to think about this...
		return shim.Error("Too few arguments")
	}

	//If asset is already present, peform validation and insert additional attributes. If not present throw error.
	value, _ := stub.GetState(args[0])
	if value == nil {
		retval := emitEvent(stub, "add", args[0], "Receive Asset", "Operations", "FAIL: Missing PO")
		if retval < 0 {
			return shim.Error("Failed to set event")
		}
		return shim.Error("Fail: Missing PO")
	}

	//Validate with invoice data for same asset
	fa := fixedAsset{}
	json.Unmarshal(value, &fa)
	values := make([]string, 3)
	values[0] = args[3]
	values[1] = args[8]
	values[2] = args[5]
	statusStr, res := validate(fa, values, false)

	//Check results - if mismatches, send out event
	if res == true {

		//Emit event
		retval := emitEvent(stub, "add", args[0], "receiveAsset", "Operations", statusStr)
		if retval < 0 {
			return shim.Error("Failed to set event")
		}

		return shim.Error(statusStr)
	}

	//Append new attributes (For now overwrite any conflicts)
	fa.Description = args[1]
	fa.AcquisitionDate = args[2]
	fa.Manufacturer = args[4]
	fa.Model = args[6]
	fa.SerialNum = args[7]
	faBytes, _ := json.Marshal(fa)
	err := stub.PutState(fa.AssetUniqueID, faBytes)
	if err != nil {
		return shim.Error("Failed to set asset")
	}
	return shim.Success(value)
}

// Data entry and validation when asset is installed
/*
Args List:
1. Asset Unique Identifier
2. Description
3. Acquisition Date
4. Location
5. Make
6. Vendor
7. Model
8. Serial Number
9. Org ID
10. Install Date
*/
func (t *smartContract) installAsset(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	value, _ := stub.GetState(args[0])

	//If asset in ledger, validate
	if value != nil {
		fa := fixedAsset{}
		json.Unmarshal(value, &fa)

		//Need to ensure that asset has been received first
		if fa.AcquisitionDate == "" {
			retval := emitEvent(stub, "install", args[0], "Install Asset", "Operations", "FAIL: Asset not received")
			if retval < 0 {
				return shim.Error("Failed to set event")
			}
			return shim.Error("Asset not received yet")
		}

		//Validate attributes
		values := make([]string, 3)
		values[0] = args[3]
		values[1] = args[8]
		values[2] = args[5]
		statusStr, res := validate(fa, values, false)

		//Check results - if mismatches, send out event
		if res == true {

			//Emit event
			retval := emitEvent(stub, "install", args[0], "Install Asset", "Operations", statusStr)
			if retval < 0 {
				return shim.Error("Failed to set event")
			}

			return shim.Error(statusStr)
		}

		//Append Install date to existing asset
		fa.InstallDate = args[9]
		faBytes, _ := json.Marshal(fa)
		err := stub.PutState(fa.AssetUniqueID, faBytes)
		if err != nil {
			return shim.Error("Failed to set asset")
		}

	} else { //Error!
		retval := emitEvent(stub, "install", args[0], "Install Asset", "Operations", "FAIL: PO Missing")
		if retval < 0 {
			return shim.Error("Failed to set event")
		}
		return shim.Error("PO Missing for this Asset")
	}

	return shim.Success(value)
}

// Entry of Invoice data on ledger
/*
Args List:
1. Asset Unique Identifier
2. Description
3. Location
4. Invoice Date
5. Vendor
6. Invoice Price
7. Org ID
*/
func (t *smartContract) receiveInvoice(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	value, _ := stub.GetState(args[0])

	//If asset is already present, peform validation and insert additional attributes. If not present throw error.
	if value != nil {
		fa := fixedAsset{}
		json.Unmarshal(value, &fa)

		//Validate attributes
		values := make([]string, 3)
		values[0] = args[2]
		values[1] = args[4]
		values[2] = args[6]
		statusStr, res := validate(fa, values, false)

		//Check results - if mismatches, send out event
		if res == true {
			//Emit event
			retval := emitEvent(stub, "invoice", args[0], "Invoice", "Finance", statusStr)
			if retval < 0 {
				return shim.Error("Failed to set event")
			}

			return shim.Error(statusStr)
		}

		//Append new data to existing asset
		fa.InvoiceDate = append(fa.InvoiceDate, args[3])
		invPrice, _ := strconv.ParseFloat(args[5], 64)
		fa.InvoicePrice += invPrice

		faBytes, _ := json.Marshal(fa)
		err := stub.PutState(fa.AssetUniqueID, faBytes)
		if err != nil {
			return shim.Error("Failed to set asset")
		}

	} else { //Error
		retval := emitEvent(stub, "invoice", args[0], "Invoice", "Finance", "FAIL: PO Missing")
		if retval < 0 {
			return shim.Error("Failed to set event")
		}
		return shim.Error("PO Missing for this Asset")
	}

	return shim.Success(value)
}

//Update Value of Asset
/*
Args List:
1. Asset Unique Identifier
2. Description
3. Purchase Price
4. Location
5. Model
6. Org_ID
7. Install Date
8. Invoice Date
9. Net Book Value
10. Serial Number
*/
func (t *smartContract) faUpdate(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	value, _ := stub.GetState(args[0])

	//If asset is already present, peform validation and insert additional attributes. If not present throw error.
	if value != nil {
		fa := fixedAsset{}
		json.Unmarshal(value, &fa)

		//Ensure that asset has been installed
		if fa.InstallDate == "" {
			retval := emitEvent(stub, "fa", args[0], "FA Update", "Finance", "FAIL: Asset not installed")
			if retval < 0 {
				return shim.Error("Failed to set event")
			}
			return shim.Error("Asset not installed yet")
		}

		//Validate asset attrs
		values := make([]string, 4)
		values[0] = args[3]
		values[1] = args[5]
		values[2] = fa.Vendor
		values[3] = args[6]
		statusStr, res := validate(fa, values, true)

		//Check results - if mismatches, send out event
		if res == true {
			//Emit event
			retval := emitEvent(stub, "fa", args[0], "FA Update", "Finance", statusStr)
			if retval < 0 {
				return shim.Error("Failed to set event")
			}

			return shim.Error(statusStr)
		}

		//Update Net Book Life of asset
		netBV, _ := strconv.ParseFloat(args[8], 64)
		fa.NetBookValue = netBV
		faBytes, _ := json.Marshal(fa)
		err := stub.PutState(fa.AssetUniqueID, faBytes)
		if err != nil {
			return shim.Error("Failed to set asset")
		}

	} else { //Error
		retval := emitEvent(stub, "fa", args[0], "FA Update", "Finance", "FAIL: PO Missing")
		if retval < 0 {
			return shim.Error("Failed to set event")
		}
		return shim.Error("PO Missing for this Asset")
	}

	return shim.Success(value)
}

//Submit Purchase Order for an asset
/*
Args List:
1. Asset Unique Identifier
2. Description
3. Purchase Price
4. Location
5. Vendor
6. Org_ID
*/
func (t *smartContract) submitPO(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 6 { //Note: Only checks number of args, args themselves can be empty.
		return shim.Error("Too few arguments")
	}

	value, _ := stub.GetState(args[0])

	//If asset in ledger, throw error
	if value != nil {
		retval := emitEvent(stub, "PO", args[0], "Purchase Order", "Finance", "FAIL: Duplicate PO")
		if retval < 0 {
			return shim.Error("Failed to set event")
		}
		return shim.Error("Asset with provided identifier already exists")
	}

	//If not, insert PO data in ledger
	floatVal, _ := strconv.ParseFloat(args[2], 64)
	var fa = fixedAsset{AssetUniqueID: args[0], Description: args[1], PurchasePrice: floatVal, Location: args[3], Vendor: args[4], OrgID: args[5]}
	faBytes, _ := json.Marshal(fa)
	err := stub.PutState(args[0], faBytes)
	if err != nil {
		return shim.Error("Failed to set asset")
	}
	return shim.Success(value)
}

//Emits failure-related events for listeners to consume
func emitEvent(stub shim.ChaincodeStubInterface, eventName string, assetUID string, function string, invoker string, status string) int {
	var errev = notifEvent{AssetUniqueID: assetUID, Function: function, Invoker: invoker, Status: status}
	errevBytes, _ := json.Marshal(errev)
	eventerr := stub.SetEvent(eventName, errevBytes)
	if eventerr != nil {
		return -1
	}
	return 0
}

//Validate asset attributes
/*
Args List:
1. Location
2. Org ID
3. Vendor
4. Install Date (Optional)
*/
func validate(fa fixedAsset, values []string, installIncluded bool) (string, bool) {

	var statusStr = "FAIL: "
	var anyMismatches = false

	if fa.Location != values[0] {
		statusStr += "Location, "
		anyMismatches = true
	}

	if fa.OrgID != values[1] {
		statusStr += "Org ID, "
		anyMismatches = true
	}

	if fa.Vendor != values[2] {
		statusStr += "Vendor"
		anyMismatches = true
	}

	if installIncluded == true {
		if fa.InstallDate != values[3] {
			statusStr += ", Install Date"
			anyMismatches = true
		}
	}

	statusStr += " Mismatch"

	return statusStr, anyMismatches
}

func main() {
	err := shim.Start(new(smartContract))
	if err != nil {
		fmt.Printf("Error starting equinix chaincode: %s", err)
	}
}
