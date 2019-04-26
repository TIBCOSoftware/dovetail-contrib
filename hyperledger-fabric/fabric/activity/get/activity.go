package get

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
	ivKey        = "key"
	ivIsPrivate  = "isPrivate"
	ivCollection = "collection"
	ovCode       = "code"
	ovMessage    = "message"
	ovKey        = "key"
	ovResult     = "result"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-get")

func init() {
	common.SetChaincodeLogLevel(log)
}

// FabricGetActivity is a stub for executing Hyperledger Fabric get operations
type FabricGetActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new FabricGetActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &FabricGetActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *FabricGetActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *FabricGetActivity) Eval(ctx activity.Context) (done bool, err error) {
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
		// retrieve data from a private collection
		return retrievePrivateData(ctx, stub, key)
	}

	// retrieve data for the key
	return retrieveData(ctx, stub, key)
}

func retrievePrivateData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, key string) (bool, error) {
	// retrieve data from a private collection
	collection, ok := ctx.GetInput(ivCollection).(string)
	if !ok || collection == "" {
		log.Error("private collection is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "private collection is not specified")
		return false, errors.New("private collection is not specified")
	}
	jsonBytes, err := ccshim.GetPrivateData(collection, key)
	if err != nil {
		log.Errorf("failed to retrieve data from private collection %s: %+v\n", collection, err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to retrieve data from private collection %s: %+v", collection, err))
		return false, errors.Wrapf(err, "failed to retrieve data from private collection %s", collection)
	}
	if jsonBytes == nil {
		log.Infof("no data found for key %s on private collection %s\n", key, collection)
		ctx.SetOutput(ovCode, 300)
		ctx.SetOutput(ovMessage, fmt.Sprintf("no data found for key %s on private collection %s", key, collection))
		ctx.SetOutput(ovKey, key)
		return true, nil
	}
	log.Debugf("retrieved from private collection %s, data: %s\n", collection, string(jsonBytes))

	var value interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to parse JSON data: %+v", err))
		return false, errors.Wrapf(err, "failed to parse JSON data %s", string(jsonBytes))
	}

	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("retrieved data from private collection %s, data: %s", collection, string(jsonBytes)))
	if result, ok := ctx.GetOutput(ovResult).(*data.ComplexObject); ok && result != nil {
		log.Debugf("set activity output result: %+v\n", value)
		result.Value = value
		ctx.SetOutput(ovResult, result)
		ctx.SetOutput(ovKey, key)
	}
	return true, nil
}

func retrieveData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, key string) (bool, error) {
	// retrieve data for the key
	jsonBytes, err := ccshim.GetState(key)
	if err != nil {
		log.Errorf("failed to retrieve data for key %s: %+v\n", key, err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to retrieve data for key %s: %+v", key, err))
		return false, errors.Wrapf(err, "failed to retrieve data for key %s", key)
	}
	if jsonBytes == nil {
		log.Infof("no data found for key %s\n", key)
		ctx.SetOutput(ovCode, 300)
		ctx.SetOutput(ovMessage, fmt.Sprintf("no data found for key %s", key))
		ctx.SetOutput(ovKey, key)
		return true, nil
	}
	log.Debugf("retrieved data from ledger: %s\n", string(jsonBytes))

	var value interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to parse JSON data: %+v", err))
		return false, errors.Wrapf(err, "failed to parse JSON data %s", string(jsonBytes))
	}

	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("retrieved data for key %s: %s", key, string(jsonBytes)))
	if result, ok := ctx.GetOutput(ovResult).(*data.ComplexObject); ok && result != nil {
		log.Debugf("set activity output result: %+v\n", value)
		result.Value = value
		ctx.SetOutput(ovResult, result)
		ctx.SetOutput(ovKey, key)
	}
	return true, nil
}
