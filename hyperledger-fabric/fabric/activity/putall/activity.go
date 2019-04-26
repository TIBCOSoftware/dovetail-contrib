package putall

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
)

const (
	ivData          = "data"
	ivIsPrivate     = "isPrivate"
	ivCollection    = "collection"
	ivCompositeKeys = "compositeKeys"
	ovCode          = "code"
	ovMessage       = "message"
	ovCount         = "count"
	ovErrors        = "errors"
	ovResult        = "result"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-putall")

func init() {
	common.SetChaincodeLogLevel(log)
}

// FabricPutAllActivity is a stub for executing Hyperledger Fabric put-all operations
type FabricPutAllActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new FabricPutAllActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &FabricPutAllActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *FabricPutAllActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *FabricPutAllActivity) Eval(ctx activity.Context) (done bool, err error) {
	// check input args
	obj, ok := ctx.GetInput(ivData).(*data.ComplexObject)
	if !ok {
		log.Errorf("input data is not a complex object\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "input data is not a complex object")
		return false, errors.New("input data is not a complex object")
	}
	valueArray, ok := obj.Value.([]interface{})
	if !ok {
		log.Errorf("input value %T is not an array of objects\n", obj.Value)
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, fmt.Sprintf("input value %T is not an array of objects", obj.Value))
		return false, errors.Errorf("input value %T is not an array of objects", obj.Value)
	}
	log.Debugf("input value type %T: %+v\n", valueArray, valueArray)

	// get composite key definitions
	compositeKeyDefs, _ := getCompositeKeyDefinition(ctx)

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, err.Error())
		return false, err
	}

	var successCount, errorCount int
	var errorKeys []string
	var resultValue []map[string]interface{}
	if isPrivate, ok := ctx.GetInput(ivIsPrivate).(bool); ok && isPrivate {
		// store data on a private collection
		collection, ok := ctx.GetInput(ivCollection).(string)
		if !ok || collection == "" {
			log.Error("private collection is not specified\n")
			ctx.SetOutput(ovCode, 400)
			ctx.SetOutput(ovMessage, "private collection is not specified")
			return false, errors.New("private collection is not specified")
		}

		for _, v := range valueArray {
			vmap := v.(map[string]interface{})
			vkey := vmap[common.KeyField].(string)
			if err := storePrivateData(stub, collection, compositeKeyDefs, vkey, vmap[common.ValueField]); err != nil {
				errorCount++
				errorKeys = append(errorKeys, vkey)
			} else {
				successCount++
				resultValue = append(resultValue, vmap)
			}
		}
	} else {
		// store data on the ledger
		for _, v := range valueArray {
			vmap := v.(map[string]interface{})
			vkey := vmap[common.KeyField].(string)
			if err := storeData(stub, compositeKeyDefs, vkey, vmap[common.ValueField]); err != nil {
				errorCount++
				errorKeys = append(errorKeys, vkey)
			} else {
				successCount++
				resultValue = append(resultValue, vmap)
			}
		}
	}

	result, ok := ctx.GetOutput(ovResult).(*data.ComplexObject)
	ctx.SetOutput(ovCount, successCount)
	ctx.SetOutput(ovErrors, errorCount)
	if errorCount > 0 {
		errMsg := fmt.Sprintf("failed to store keys: %s", strings.Join(errorKeys, ","))
		ctx.SetOutput(ovMessage, errMsg)
		if successCount > 0 {
			// return 300 if partial successs
			ctx.SetOutput(ovCode, 300)
			if ok && result != nil {
				result.Value = resultValue
				ctx.SetOutput(ovResult, result)
			}
			return true, nil
		}
		// return 500 if all failures
		ctx.SetOutput(ovCode, 500)
		return false, errors.New(errMsg)
	}
	// return 200 if no errors
	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("stored data on ledger: %+v", resultValue))
	if ok && result != nil {
		log.Debugf("set activity output result: %+v\n", resultValue)
		result.Value = resultValue
		ctx.SetOutput(ovResult, result)
	}
	return true, nil
}

func storePrivateData(ccshim shim.ChaincodeStubInterface, collection string, compositeKeyDefs map[string][]string, key string, value interface{}) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		log.Errorf("failed to marshal value '%+v', error: %+v\n", value, err)
		return errors.Wrapf(err, "failed to marshal value: %+v", value)
	}

	// store data on a private collection
	if err := ccshim.PutPrivateData(collection, key, jsonBytes); err != nil {
		log.Errorf("failed to store data in private collection %s: %+v\n", collection, err)
		return errors.Wrapf(err, "failed to store data in private collection %s", collection)
	}
	log.Debugf("stored in private collection %s, data: %s\n", collection, string(jsonBytes))

	// store composite keys if required
	if compositeKeyDefs == nil {
		return nil
	}
	compositeKeys := common.ExtractCompositeKeys(ccshim, compositeKeyDefs, key, value)
	if compositeKeys != nil && len(compositeKeys) > 0 {
		for _, k := range compositeKeys {
			cv := []byte{0x00}
			if err := ccshim.PutPrivateData(collection, k, cv); err != nil {
				log.Errorf("failed to store composite key %s on collection %s: %+v\n", k, collection, err)
			} else {
				log.Debugf("stored composite key %s on collection %s\n", k, collection)
			}
		}
	}
	return nil
}

func storeData(ccshim shim.ChaincodeStubInterface, compositeKeyDefs map[string][]string, key string, value interface{}) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		log.Errorf("failed to marshal value '%+v', error: %+v\n", value, err)
		return errors.Wrapf(err, "failed to marshal value: %+v", value)
	}
	// store data on the ledger
	if err := ccshim.PutState(key, jsonBytes); err != nil {
		log.Errorf("failed to store data on ledger: %+v\n", err)
		return errors.Errorf("failed to store data on ledger: %+v", err)
	}
	log.Debugf("stored data on ledger: %s\n", string(jsonBytes))

	// store composite keys if required
	if compositeKeyDefs == nil {
		return nil
	}
	compositeKeys := common.ExtractCompositeKeys(ccshim, compositeKeyDefs, key, value)
	if compositeKeys != nil && len(compositeKeys) > 0 {
		for _, k := range compositeKeys {
			cv := []byte{0x00}
			if err := ccshim.PutState(k, cv); err != nil {
				log.Errorf("failed to store composite key %s: %+v\n", k, err)
			} else {
				log.Debugf("stored composite key %s\n", k)
			}
		}
	}
	return nil
}

func getCompositeKeyDefinition(ctx activity.Context) (map[string][]string, error) {
	if ckJSON, ok := ctx.GetInput(ivCompositeKeys).(string); ok && ckJSON != "" {
		log.Debugf("Got composite key definition: %s\n", ckJSON)
		ckDefs := make(map[string][]string)
		if err := json.Unmarshal([]byte(ckJSON), &ckDefs); err != nil {
			log.Warningf("failed to unmarshal composite key definitions: %+v\n", err)
			return nil, err
		}
		log.Debugf("Parsed composite key definitions: %+v\n", ckDefs)
		return ckDefs, nil
	}
	log.Debugf("No composite key is defined")
	return nil, nil
}
