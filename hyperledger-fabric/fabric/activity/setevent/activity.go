package setevent

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
	ivName    = "name"
	ivPayload = "payload"
	ovCode    = "code"
	ovMessage = "message"
	ovName    = "name"
	ovResult  = "result"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-setevent")

func init() {
	common.SetChaincodeLogLevel(log)
}

// FabricEventActivity is a stub for executing Hyperledger Fabric set-event operations
type FabricEventActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new FabricEventActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &FabricEventActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *FabricEventActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *FabricEventActivity) Eval(ctx activity.Context) (done bool, err error) {
	// check input args
	name, ok := ctx.GetInput(ivName).(string)
	if !ok || name == "" {
		log.Error("event name is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "event name is not specified")
		return false, errors.New("event name is not specified")
	}
	log.Debugf("event name: %s\n", name)

	var jsonBytes []byte
	payload := getEventPayload(ctx)
	if payload != nil {
		jsonBytes, err = json.Marshal(payload)
		if err != nil {
			log.Warningf("failed to marshal payload '%+v', error: %+v\n", payload, err)
		}
	}
	log.Debugf("event payload: %+v\n", jsonBytes)

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, err.Error())
		return false, err
	}

	// set fabric event
	if err := stub.SetEvent(name, jsonBytes); err != nil {
		log.Errorf("failed to set event %s, error: %+v\n", name, err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, err.Error())
		return false, err
	}

	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("set event %s, payload: %s", name, string(jsonBytes)))
	ctx.SetOutput(ovName, name)
	if result, ok := ctx.GetOutput(ovResult).(*data.ComplexObject); ok && result != nil {
		log.Debugf("set activity output result: %+v\n", payload)
		result.Value = payload
		ctx.SetOutput(ovResult, result)
	}
	return true, nil
}

func getEventPayload(ctx activity.Context) interface{} {
	payload, ok := ctx.GetInput(ivPayload).(*data.ComplexObject)
	if !ok {
		log.Debug("payload is not a complex object\n")
		return nil
	}
	return payload.Value
}
