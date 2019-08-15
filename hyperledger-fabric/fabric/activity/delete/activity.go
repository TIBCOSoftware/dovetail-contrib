package delete

import (
	"encoding/json"
	"fmt"

	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/project-flogo/core/activity"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-delete")

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

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		log.Errorf("failed to retrieve fabric stub: %+v\n", err)
		output := &Output{Code: 500, Message: err.Error()}
		ctx.SetOutputObject(output)
		return false, err
	}

	if input.PrivateCollection != "" {
		// delete data from a private collection
		return deletePrivateData(ctx, stub, input)
	}

	// delete data from the ledger
	return deleteData(ctx, stub, input)
}

func deletePrivateData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, input *Input) (bool, error) {
	// retrieves data for managing composite keys and map to output
	jsonBytes, err := ccshim.GetPrivateData(input.PrivateCollection, input.StateKey)
	if err != nil {
		log.Errorf("failed to get data '%s' from private collection '%s': %+v\n", input.StateKey, input.PrivateCollection, err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to get data '%s' from private collection '%s'", input.StateKey, input.PrivateCollection)}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	if jsonBytes == nil {
		log.Infof("no data found for '%s' from private collection '%s'\n", input.StateKey, input.PrivateCollection)
		output := &Output{Code: 300, Message: fmt.Sprintf("no data found for '%s' from private collection '%s'", input.StateKey, input.PrivateCollection)}
		ctx.SetOutputObject(output)
		return true, nil
	}
	var value map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to parse JSON data %s", string(jsonBytes))}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	// delete data if keyOnly is not specified or keyOnly=false
	if !input.KeysOnly {
		if err := ccshim.DelPrivateData(input.PrivateCollection, input.StateKey); err != nil {
			log.Errorf("failed to delete data from private collection %s: %+v\n", input.PrivateCollection, err)
			output := &Output{Code: 500, Message: fmt.Sprintf("failed to delete data from private collection %s", input.PrivateCollection)}
			ctx.SetOutputObject(output)
			return false, errors.Wrapf(err, output.Message)
		}
		log.Debugf("deleted from private collection %s, data: %s\n", input.PrivateCollection, string(jsonBytes))
	}

	// delete composite keys if specified
	if compositeKeyDefs, _ := getCompositeKeyDefinition(input.CompositeKeys); compositeKeyDefs != nil {
		compKeys := common.ExtractCompositeKeys(ccshim, compositeKeyDefs, input.StateKey, value)
		if compKeys != nil && len(compKeys) > 0 {
			for _, k := range compKeys {
				if err := ccshim.DelPrivateData(input.PrivateCollection, k); err != nil {
					log.Errorf("failed to delete composite key %s from collection %s: %+v\n", k, input.PrivateCollection, err)
				} else {
					log.Debugf("deleted composite key %s from collection %s\n", k, input.PrivateCollection)
				}
			}
		}
	}

	output := &Output{
		Code:     200,
		Message:  fmt.Sprintf("deleted from private collection %s, data: %s", input.PrivateCollection, string(jsonBytes)),
		StateKey: input.StateKey,
		Result:   value,
	}
	ctx.SetOutputObject(output)
	return true, nil
}

func deleteData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, input *Input) (bool, error) {
	// retrieves data for managing composite keys and map to output
	jsonBytes, err := ccshim.GetState(input.StateKey)
	if err != nil {
		log.Errorf("failed to get data '%s': %+v\n", input.StateKey, err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to get data for key '%s'", input.StateKey)}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	if jsonBytes == nil {
		log.Infof("no data found for '%s'\n", input.StateKey)
		output := &Output{Code: 300, Message: fmt.Sprintf("no data found for '%s'", input.StateKey)}
		ctx.SetOutputObject(output)
		return true, nil
	}
	var value map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to parse JSON data: %s", string(jsonBytes))}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	// delete data if keyOnly is not specified or keyOnly=false
	if !input.KeysOnly {
		if err := ccshim.DelState(input.StateKey); err != nil {
			log.Errorf("failed to delete data: %+v\n", err)
			output := &Output{Code: 500, Message: fmt.Sprintf("failed to delete data for key %s", input.StateKey)}
			return false, errors.Wrapf(err, output.Message)
		}
		log.Debugf("deleted data: %s\n", string(jsonBytes))
	}

	// delete composite keys if specified
	if compositeKeyDefs, _ := getCompositeKeyDefinition(input.CompositeKeys); compositeKeyDefs != nil {
		compKeys := common.ExtractCompositeKeys(ccshim, compositeKeyDefs, input.StateKey, value)
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

	output := &Output{
		Code:     200,
		Message:  fmt.Sprintf("deleted data: %s", string(jsonBytes)),
		StateKey: input.StateKey,
		Result:   value,
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
