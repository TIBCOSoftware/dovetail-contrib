/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package ledger

// Imports
import (
	"encoding/json"
	"fmt"

	"github.com/TIBCOSoftware/dovetail-contrib/libraries/fabric-go/utils"
	"github.com/TIBCOSoftware/flogo-lib/core/data"

	dtsvc "github.com/TIBCOSoftware/dovetail-contrib/libraries/fabric-go/runtime/services"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

// Constants
const (
	ivAssetName      = "assetName"
	ivData           = "input"
	ivStub           = "containerServiceStub"
	ivOperation      = "operation"
	ivAssetKey       = "identifier"
	ivAssetLookupKey = "compositeKey"
	ivCompositeKeys  = "compositeKeys"
	ivIsArray        = "isArray"
	ovOutput         = "output"
)

// describes the metadata of the activity as found in the activity.json file
type LedgerActivity struct {
	metadata *activity.Metadata
}

// NewActivity will instantiate a new LedgerActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &LedgerActivity{metadata: metadata}
}

// Metadata will return the metadata of the LedgerActivity
func (a *LedgerActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval executes the activity
func (a *LedgerActivity) Eval(context activity.Context) (done bool, err error) {
	stub, err := utils.GetContainerStub(context)
	if err != nil {
		return false, err
	}

	logger := stub.GetLogService()

	logger.Debug("Enter ledger activity...")
	isArray, ok := context.GetInput(ivIsArray).(bool)
	if !ok {
		return false, fmt.Errorf("asset name is not initialized")
	}

	assetName, ok := context.GetInput(ivAssetName).(string)
	if !ok {
		return false, fmt.Errorf("asset name is not initialized")
	}

	assetKey, ok := context.GetInput(ivAssetKey).(string)
	if !ok {
		return false, fmt.Errorf("asset key is not initialized")
	}
	operation, ok := context.GetInput(ivOperation).(string)
	if !ok {
		return false, fmt.Errorf("operation is not initialized")
	}

	inputValue, err := data.CoerceToComplexObject(context.GetInput(ivData))
	if err != nil {
		return false, fmt.Errorf("asset value is not initialized")
	}

	if inputValue == nil {
		return false, fmt.Errorf("asset value is not initialized")
	}
	inputs, err := utils.GetInputData(inputValue, isArray)
	if err != nil {
		return false, err
	}

	compositeKeys := context.GetInput(ivCompositeKeys)
	lookupKey := context.GetInput(ivAssetLookupKey)

	result := make([][]byte, 0)
	output := make([]interface{}, 0)
	complexOutput := &data.ComplexObject{}
	for _, av := range inputs {
		asset, ok := av.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("can not cast unmartialed asset value to instance of map[string]interface{}")
		}
		fmt.Printf("ledger input %v\n", asset)
		if operation == "DELETE" {
			record, err := deleteRecord(assetName, assetKey, asset, compositeKeys, stub)
			if err != nil {
				return false, err
			}

			if record == nil {
				output = append(output, asset)
			} else {
				result = append(result, record)
			}
		} else if operation == "PUT" {
			err = putRecord(assetName, assetKey, asset, compositeKeys, stub)
			if err != nil {
				return false, err
			}
			output = append(output, asset)
		} else if operation == "GET" {
			rawvalue, err := getRecord(assetName, assetKey, asset, stub)
			if err != nil {
				return false, err
			}
			if rawvalue != nil {
				result = append(result, rawvalue)
			}
		} else {
			//LOOKUP
			records, err := lookupRecords(assetName, assetKey, lookupKey.(string), asset, stub)
			if err != nil {
				return false, err
			}

			if records != nil && len(records) > 0 {
				result = append(result, records...)
			}
		}
	}

	for _, v := range result {
		m, err := utils.ParseRecord(v)
		if err != nil {
			return false, err
		}
		output = append(output, m)
	}

	if isArray || operation == "LOOKUP" {
		complexOutput.Value = output
	} else {
		if len(output) > 0 {
			complexOutput.Value = output[0]
		}
	}

	context.SetOutput(ovOutput, complexOutput)
	logger.Debug("Exit ledger activity")
	return true, nil
}

func parseRecord(record []byte) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	err := json.Unmarshal(record, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func deleteRecord(assetName string, identifier string, input map[string]interface{}, compositeKeys interface{}, stub dtsvc.ContainerService) ([]byte, error) {
	datasvc := stub.GetDataService()
	return datasvc.DeleteState(assetName, identifier, input, compositeKeys)
}

func putRecord(assetName string, identifier string, input map[string]interface{}, compositeKeys interface{}, stub dtsvc.ContainerService) error {

	datasvc := stub.GetDataService()
	return datasvc.PutState(assetName, identifier, input, compositeKeys)
}

func lookupRecords(assetName string, identifier string, lookupKey string, input map[string]interface{}, stub dtsvc.ContainerService) ([][]byte, error) {
	datasvc := stub.GetDataService()
	return datasvc.LookupState(assetName, identifier, lookupKey, input)
}

func getRecord(assetName string, identifier string, input map[string]interface{}, stub dtsvc.ContainerService) ([]byte, error) {
	datasvc := stub.GetDataService()
	return datasvc.GetState(assetName, identifier, input)
}
