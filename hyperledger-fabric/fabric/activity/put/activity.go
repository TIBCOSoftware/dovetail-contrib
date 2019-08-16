package put

import (
	"encoding/json"
	"fmt"

	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/project-flogo/core/activity"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-put")

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	common.SetChaincodeLogLevel(log)
	_ = activity.Register(&Activity{}, New)
}

// Activity is a stub for executing Hyperledger Fabric put operations
type Activity struct {
}

// New creates a new Activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	return &Activity{}, nil
}

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	// check input args
	input := &Input{}
	if err = ctx.GetInputObject(input); err != nil {
		return false, err
	}
	if input.StateKey == "" {
		log.Error("state key is not specified\n")
		output := &Output{Code: 400, Message: "state key is not specified"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	log.Debugf("state key: %s\n", input.StateKey)

	if input.StateData == nil {
		log.Errorf("input data is nil\n")
		output := &Output{Code: 400, Message: "input data is nil"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	log.Debugf("input value type %T: %+v\n", input.StateData, input.StateData)

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		log.Errorf("failed to retrieve fabric stub: %+v\n", err)
		output := &Output{Code: 500, Message: err.Error()}
		ctx.SetOutputObject(output)
		return false, err
	}

	if input.PrivateCollection != "" {
		// store data on a private collection
		return storePrivateData(ctx, stub, input)
	}

	// store data on the ledger
	return storeData(ctx, stub, input)
}

func storePrivateData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, input *Input) (bool, error) {
	jsonBytes, err := json.Marshal(input.StateData)
	if err != nil {
		log.Errorf("failed to marshal value '%+v', error: %+v\n", input.StateData, err)
		output := &Output{Code: 400, Message: fmt.Sprintf("failed to marshal value: %+v", input.StateData)}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	// store data on a private collection
	if err := ccshim.PutPrivateData(input.PrivateCollection, input.StateKey, jsonBytes); err != nil {
		log.Errorf("failed to store data in private collection %s: %+v\n", input.PrivateCollection, err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to store data in private collection %s", input.PrivateCollection)}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	log.Debugf("stored in private collection %s, data: %s\n", input.PrivateCollection, string(jsonBytes))

	// store composite keys if required
	if compositeKeyDefs, _ := getCompositeKeyDefinition(input.CompositeKeys); compositeKeyDefs != nil {
		compKeys := common.ExtractCompositeKeys(ccshim, compositeKeyDefs, input.StateKey, input.StateData)
		if compKeys != nil && len(compKeys) > 0 {
			for _, k := range compKeys {
				cv := []byte{0x00}
				if err := ccshim.PutPrivateData(input.PrivateCollection, k, cv); err != nil {
					log.Errorf("failed to store composite key %s on collection %s: %+v\n", k, input.PrivateCollection, err)
				} else {
					log.Debugf("stored composite key %s on collection %s\n", k, input.PrivateCollection)
				}
			}
		}
	}

	output := &Output{
		Code:     200,
		Message:  fmt.Sprintf("stored in private collection %s, data: %s", input.PrivateCollection, string(jsonBytes)),
		StateKey: input.StateKey,
		Result:   input.StateData,
	}
	ctx.SetOutputObject(output)
	return true, nil
}

func storeData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, input *Input) (bool, error) {
	jsonBytes, err := json.Marshal(input.StateData)
	if err != nil {
		log.Errorf("failed to marshal value '%+v', error: %+v\n", input.StateData, err)
		output := &Output{Code: 400, Message: fmt.Sprintf("failed to marshal value: %+v", input.StateData)}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	// store data on the ledger
	if err := ccshim.PutState(input.StateKey, jsonBytes); err != nil {
		log.Errorf("failed to store data on ledger: %+v\n", err)
		output := &Output{Code: 500, Message: "failed to store data on ledger"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	log.Debugf("stored data on ledger: %s\n", string(jsonBytes))

	// store composite keys if required
	if compositeKeyDefs, _ := getCompositeKeyDefinition(input.CompositeKeys); compositeKeyDefs != nil {
		compKeys := common.ExtractCompositeKeys(ccshim, compositeKeyDefs, input.StateKey, input.StateData)
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

	output := &Output{
		Code:     200,
		Message:  fmt.Sprintf("stored data on ledger: %s", string(jsonBytes)),
		StateKey: input.StateKey,
		Result:   input.StateData,
	}
	ctx.SetOutputObject(output)
	return true, nil
}

func getCompositeKeyDefinition(compositeKeys string) (map[string][]string, error) {
	if compositeKeys != "" {
		log.Debugf("Got composite key definition: %s\n", compositeKeys)
		ckDefs := make(map[string][]string)
		if err := json.Unmarshal([]byte(compositeKeys), &ckDefs); err != nil {
			log.Warningf("failed to unmarshal composite key definitions: %+v\n", err)
			return nil, err
		}
		log.Debugf("Parsed composite key definitions: %+v\n", ckDefs)
		return ckDefs, nil
	}
	log.Debugf("No composite key is defined")
	return nil, nil
}
