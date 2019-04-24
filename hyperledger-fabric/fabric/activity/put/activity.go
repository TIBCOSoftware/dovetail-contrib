package put

import (
	"encoding/json"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
)

const (
	ivKey           = "key"
	ivValueType     = "valueType"
	ivValue         = "value"
	ivData          = "data"
	ivIsPrivate     = "isPrivate"
	ivCollection    = "collection"
	ivCompositeKeys = "compositeKeys"
	ovCode          = "code"
	ovMessage       = "message"
	ovKey           = "key"
	ovResult        = "result"
	objectType      = "object"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-put")

func init() {
	common.SetChaincodeLogLevel(log)
}

// FabricPutActivity is a stub for executing Hyperledger Fabric put operations
type FabricPutActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new FabricPutActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &FabricPutActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *FabricPutActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *FabricPutActivity) Eval(ctx activity.Context) (done bool, err error) {
	// check input args
	key, ok := ctx.GetInput(ivKey).(string)
	if !ok || key == "" {
		log.Error("state key is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "state key is not specified")
		return false, errors.New("state key is not specified")
	}
	log.Debugf("state key: %s\n", key)
	vtype := ctx.GetInput(ivValueType).(string)
	log.Debugf("value type: %s\n", vtype)
	value := ctx.GetInput(ivValue)
	if vtype == objectType {
		if obj, ok := ctx.GetInput(ivData).(*data.ComplexObject); ok {
			value = obj.Value
		} else {
			log.Errorf("input data is not a complex object\n")
			ctx.SetOutput(ovCode, 400)
			ctx.SetOutput(ovMessage, "input data is not a complex object")
			return false, errors.New("input data is not a complex object")
		}
	}
	log.Debugf("input value type %T: %+v\n", value, value)

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, err.Error())
		return false, err
	}

	if isPrivate, ok := ctx.GetInput(ivIsPrivate).(bool); ok && isPrivate {
		// store data on a private collection
		return storePrivateData(ctx, stub, key, value)
	}

	// store data on the ledger
	return storeData(ctx, stub, key, value)
}

func storePrivateData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, key string, value interface{}) (bool, error) {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		log.Errorf("failed to marshal value '%+v', error: %+v\n", value, err)
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to marshal value: %+v", err))
		return false, errors.Wrapf(err, "failed to marshal value: %+v", value)
	}
	// store data on a private collection
	collection, ok := ctx.GetInput(ivCollection).(string)
	if !ok || collection == "" {
		log.Error("private collection is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "private collection is not specified")
		return false, errors.New("private collection is not specified")
	}
	if err := ccshim.PutPrivateData(collection, key, jsonBytes); err != nil {
		log.Errorf("failed to store data in private collection %s: %+v\n", collection, err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to store data in private collection %s: %+v", collection, err))
		return false, errors.Wrapf(err, "failed to store data in private collection %s", collection)
	}
	log.Debugf("stored in private collection %s, data: %s\n", collection, string(jsonBytes))

	// store composite keys if required
	if compositeKeyDefs, _ := getCompositeKeyDefinition(ctx); compositeKeyDefs != nil {
		compKeys := common.ExtractCompositeKeys(ccshim, compositeKeyDefs, key, value)
		if compKeys != nil && len(compKeys) > 0 {
			for _, k := range compKeys {
				cv := []byte{0x00}
				if err := ccshim.PutPrivateData(collection, k, cv); err != nil {
					log.Errorf("failed to store composite key %s on collection %s: %+v\n", k, collection, err)
				} else {
					log.Debugf("stored composite key %s on collection %s\n", k, collection)
				}
			}
		}
	}

	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("stored in private collection %s, data: %s", collection, string(jsonBytes)))
	if result, ok := ctx.GetOutput(ovResult).(*data.ComplexObject); ok && result != nil {
		log.Debugf("set activity output result: %+v\n", value)
		result.Value = value
		ctx.SetOutput(ovResult, result)
		ctx.SetOutput(ovKey, key)
	}
	return true, nil
}

func storeData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, key string, value interface{}) (bool, error) {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		log.Errorf("failed to marshal value '%+v', error: %+v\n", value, err)
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to marshal value: %+v", err))
		return false, errors.Wrapf(err, "failed to marshal value: %+v", value)
	}
	// store data on the ledger
	if err := ccshim.PutState(key, jsonBytes); err != nil {
		log.Errorf("failed to store data on ledger: %+v\n", err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to store data on ledger: %+v", err))
		return false, errors.Errorf("failed to store data on ledger: %+v", err)
	}
	log.Debugf("stored data on ledger: %s\n", string(jsonBytes))

	// store composite keys if required
	if compositeKeyDefs, _ := getCompositeKeyDefinition(ctx); compositeKeyDefs != nil {
		compKeys := common.ExtractCompositeKeys(ccshim, compositeKeyDefs, key, value)
		if compKeys != nil && len(compKeys) > 0 {
			for _, k := range compKeys {
				cv := []byte{0x00}
				if err := ccshim.PutState(k, cv); err != nil {
					log.Errorf("failed to store composite key %s: %+v\n", k, err)
				} else {
					log.Debugf("stored composite key %s\n", k)
				}
			}
		}
	}

	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("stored data on ledger: %s", string(jsonBytes)))
	if result, ok := ctx.GetOutput(ovResult).(*data.ComplexObject); ok && result != nil {
		log.Debugf("set activity output result: %+v\n", value)
		result.Value = value
		ctx.SetOutput(ovResult, result)
		ctx.SetOutput(ovKey, key)
	}
	return true, nil
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
