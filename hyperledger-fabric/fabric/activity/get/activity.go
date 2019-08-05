package get

import (
	"encoding/json"
	"fmt"

	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/project-flogo/core/activity"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-get")

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	common.SetChaincodeLogLevel(log)
	_ = activity.Register(&Activity{}, New)
}

// Activity is a stub for executing Hyperledger Fabric get operations
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

	if input.IsPrivate {
		// retrieve data from a private collection
		return retrievePrivateData(ctx, stub, input)
	}

	// retrieve data for the key
	return retrieveData(ctx, stub, input.StateKey)
}

func retrievePrivateData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, input *Input) (bool, error) {
	// retrieve data from a private collection
	if input.PrivateCollection == "" {
		log.Error("private collection is not specified\n")
		output := &Output{Code: 400, Message: "private collection is not specified"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	jsonBytes, err := ccshim.GetPrivateData(input.PrivateCollection, input.StateKey)
	if err != nil {
		log.Errorf("failed to retrieve data from private collection %s: %+v\n", input.PrivateCollection, err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to retrieve data from private collection %s", input.PrivateCollection)}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	if jsonBytes == nil {
		log.Infof("no data found for key %s on private collection %s\n", input.StateKey, input.PrivateCollection)
		output := &Output{Code: 300,
			Message:  fmt.Sprintf("no data found for key %s on private collection %s", input.StateKey, input.PrivateCollection),
			StateKey: input.StateKey,
		}
		ctx.SetOutputObject(output)
		return true, nil
	}
	log.Debugf("retrieved from private collection %s, data: %s\n", input.PrivateCollection, string(jsonBytes))

	var value map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to parse JSON data: %s", string(jsonBytes))}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	output := &Output{
		Code:     200,
		Message:  fmt.Sprintf("retrieved data from private collection %s, data: %s", input.PrivateCollection, string(jsonBytes)),
		StateKey: input.StateKey,
		Result:   value,
	}
	ctx.SetOutputObject(output)
	return true, nil
}

func retrieveData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, key string) (bool, error) {
	// retrieve data for the key
	jsonBytes, err := ccshim.GetState(key)
	if err != nil {
		log.Errorf("failed to retrieve data for key %s: %+v\n", key, err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to retrieve data for key %s", key)}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	if jsonBytes == nil {
		log.Infof("no data found for key %s\n", key)
		output := &Output{Code: 300,
			Message:  fmt.Sprintf("no data found for key %s", key),
			StateKey: key,
		}
		ctx.SetOutputObject(output)
		return true, nil
	}
	log.Debugf("retrieved data from ledger: %s\n", string(jsonBytes))

	var value map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to parse JSON data: %s", string(jsonBytes))}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	output := &Output{
		Code:     200,
		Message:  fmt.Sprintf("retrieved data for key %s: %s", key, string(jsonBytes)),
		StateKey: key,
		Result:   value,
	}
	ctx.SetOutputObject(output)
	return true, nil
}
