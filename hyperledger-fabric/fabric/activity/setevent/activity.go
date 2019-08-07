package setevent

import (
	"encoding/json"
	"fmt"

	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/project-flogo/core/activity"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-setevent")

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

	if input.Name == "" {
		log.Error("event name is not specified\n")
		output := &Output{Code: 400, Message: "event name is not specified"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	log.Debugf("event name: %s\n", input.Name)

	var jsonBytes []byte
	if input.Payload != nil {
		jsonBytes, err = json.Marshal(input.Payload)
		if err != nil {
			log.Warningf("failed to marshal payload '%+v', error: %+v\n", input.Payload, err)
		}
	}
	log.Debugf("event payload: %+v\n", jsonBytes)

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		log.Errorf("failed to retrieve fabric stub: %+v\n", err)
		output := &Output{Code: 500, Message: err.Error()}
		ctx.SetOutputObject(output)
		return false, err
	}

	// set fabric event
	if err := stub.SetEvent(input.Name, jsonBytes); err != nil {
		log.Errorf("failed to set event %s, error: %+v\n", input.Name, err)
		output := &Output{Code: 500, Message: err.Error()}
		ctx.SetOutputObject(output)
		return false, err
	}

	log.Debugf("set activity output result: %+v\n", input.Payload)
	output := &Output{Code: 200,
		Message: fmt.Sprintf("set event %s, payload: %s", input.Name, string(jsonBytes)),
		Name:    input.Name,
		Result:  input.Payload,
	}
	ctx.SetOutputObject(output)
	return true, nil
}
