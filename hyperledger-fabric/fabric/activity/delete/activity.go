package delete

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
	ivIsPrivate     = "isPrivate"
	ivCollection    = "collection"
	ivKeysOnly      = "keysOnly"
	ivCompositeKeys = "compositeKeys"
	ovCode          = "code"
	ovMessage       = "message"
	ovKey           = "key"
	ovResult        = "result"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-delete")

func init() {
	common.SetChaincodeLogLevel(log)
}

// FabricDeleteActivity is a stub for executing Hyperledger Fabric delete operations
type FabricDeleteActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new FabricDeleteActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &FabricDeleteActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *FabricDeleteActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *FabricDeleteActivity) Eval(ctx activity.Context) (done bool, err error) {
	// check input args
	key, ok := ctx.GetInput(ivKey).(string)
	if !ok || key == "" {
		log.Error("state key is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "state key is not specified")
		return false, errors.New("state key is not specified")
	}
	log.Debugf("state key: %s\n", key)

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, err.Error())
		return false, err
	}

	if isPrivate, ok := ctx.GetInput(ivIsPrivate).(bool); ok && isPrivate {
		// delete data from a private collection
		return deletePrivateData(ctx, stub, key)
	}

	// delete data from the ledger
	return deleteData(ctx, stub, key)
}

func deletePrivateData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, key string) (bool, error) {
	// delete data on a private collection
	collection, ok := ctx.GetInput(ivCollection).(string)
	if !ok || collection == "" {
		log.Error("private collection is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "private collection is not specified")
		return false, errors.New("private collection is not specified")
	}

	// retrieves data for managing composite keys and map to output
	jsonBytes, err := ccshim.GetPrivateData(collection, key)
	if err != nil {
		log.Errorf("failed to get data '%s' from private collection '%s': %+v\n", key, collection, err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to get data '%s' from private collection '%s': %+v", key, collection, err))
		return false, errors.Wrapf(err, "failed to get data '%s' from private collection '%s'", key, collection)
	}
	if jsonBytes == nil {
		log.Infof("no data found for '%s' from private collection '%s'\n", key, collection)
		ctx.SetOutput(ovCode, 300)
		ctx.SetOutput(ovMessage, fmt.Sprintf("no data found for '%s' from private collection '%s'", key, collection))
		return true, nil
	}
	var value interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to parse JSON data: %+v", err))
		return false, errors.Wrapf(err, "failed to parse JSON data %s", string(jsonBytes))
	}

	// delete data if keyOnly is not specified or keyOnly=false
	if keysOnly, ok := ctx.GetInput(ivKeysOnly).(bool); !ok || !keysOnly {
		if err := ccshim.DelPrivateData(collection, key); err != nil {
			log.Errorf("failed to delete data from private collection %s: %+v\n", collection, err)
			ctx.SetOutput(ovCode, 500)
			ctx.SetOutput(ovMessage, fmt.Sprintf("failed to delete data from private collection %s: %+v", collection, err))
			return false, errors.Wrapf(err, "failed to delete data from private collection %s", collection)
		}
		log.Debugf("deleted from private collection %s, data: %s\n", collection, string(jsonBytes))
	}

	// delete composite keys if specified
	if compositeKeyDefs, _ := getCompositeKeyDefinition(ctx); compositeKeyDefs != nil {
		compKeys := common.ExtractCompositeKeys(ccshim, compositeKeyDefs, key, value)
		if compKeys != nil && len(compKeys) > 0 {
			for _, k := range compKeys {
				if err := ccshim.DelPrivateData(collection, k); err != nil {
					log.Errorf("failed to delete composite key %s from collection %s: %+v\n", k, collection, err)
				} else {
					log.Debugf("deleted composite key %s from collection %s\n", k, collection)
				}
			}
		}
	}

	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("deleted from private collection %s, data: %s", collection, string(jsonBytes)))
	if result, ok := ctx.GetOutput(ovResult).(*data.ComplexObject); ok && result != nil {
		log.Debugf("set activity output result: %+v\n", value)
		result.Value = value
		ctx.SetOutput(ovResult, result)
		ctx.SetOutput(ovKey, key)
	}
	return true, nil
}

func deleteData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, key string) (bool, error) {
	// retrieves data for managing composite keys and map to output
	jsonBytes, err := ccshim.GetState(key)
	if err != nil {
		log.Errorf("failed to get data '%s': %+v\n", key, err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to get data '%s': %+v", key, err))
		return false, errors.Wrapf(err, "failed to get data '%s'", key)
	}
	if jsonBytes == nil {
		log.Infof("no data found for '%s'\n", key)
		ctx.SetOutput(ovCode, 300)
		ctx.SetOutput(ovMessage, fmt.Sprintf("no data found for '%s'", key))
		return true, nil
	}
	var value interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to parse JSON data: %+v", err))
		return false, errors.Wrapf(err, "failed to parse JSON data %s", string(jsonBytes))
	}

	// delete data if keyOnly is not specified or keyOnly=false
	if keysOnly, ok := ctx.GetInput(ivKeysOnly).(bool); !ok || !keysOnly {
		if err := ccshim.DelState(key); err != nil {
			log.Errorf("failed to delete data: %+v\n", err)
			ctx.SetOutput(ovCode, 500)
			ctx.SetOutput(ovMessage, fmt.Sprintf("failed to delete data: %+v", err))
			return false, errors.Wrapf(err, "failed to delete data")
		}
		log.Debugf("deleted data: %s\n", string(jsonBytes))
	}

	// delete composite keys if specified
	if compositeKeyDefs, _ := getCompositeKeyDefinition(ctx); compositeKeyDefs != nil {
		compKeys := common.ExtractCompositeKeys(ccshim, compositeKeyDefs, key, value)
		if compKeys != nil && len(compKeys) > 0 {
			for _, k := range compKeys {
				if err := ccshim.DelState(k); err != nil {
					log.Errorf("failed to delete composite key %s: %+v\n", k, err)
				} else {
					log.Debugf("deleted composite key %s\n", k)
				}
			}
		}
	}

	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("deleted data: %s", string(jsonBytes)))
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
