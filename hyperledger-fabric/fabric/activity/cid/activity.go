package cid

import (
	"encoding/json"

	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/attrmgr"
	ci "github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/pkg/errors"
	"github.com/project-flogo/core/activity"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-cid")

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

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		log.Errorf("failed to retrieve fabric stub: %+v\n", err)
		output := &Output{Code: 500, Message: err.Error()}
		ctx.SetOutputObject(output)
		return false, err
	}

	// retrieve data for the key
	return retrieveCid(ctx, stub)
}

func retrieveCid(ctx activity.Context, ccshim shim.ChaincodeStubInterface) (bool, error) {
	// retrieve data for the key
	id, err := ci.GetID(ccshim)
	if err != nil {
		log.Errorf("failed to extract client ID from stub: %+v\n", err)
		output := &Output{Code: 500, Message: "failed to extract client ID from stub"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	mspid, err := ci.GetMSPID(ccshim)
	if err != nil {
		log.Errorf("failed to extract client MSPID from stub: %+v\n", err)
		output := &Output{Code: 500, Message: "failed to extract client MSPID from stub"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	creator, err := ccshim.GetCreator()
	if err != nil {
		log.Errorf("failed to extract proposal creator from stub: %+v\n", err)
		output := &Output{Code: 500, Message: "failed to extract proposal creator from stub"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	attrs, err := attrmgr.New().GetAttributesFromIdemix(creator)
	if err != nil {
		log.Errorf("failed to extract attributes from cerficate: %+v\n", err)
		output := &Output{Code: 500, Message: "failed to extract attributes from cerficate"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	output := &Output{Code: 200,
		Cid:   id,
		Mspid: mspid,
		Attrs: attrs.Attrs,
	}
	msgBytes, err := json.Marshal(output)
	if err != nil {
		log.Errorf("failed to serialize JSON output: %+v\n", err)
		output := &Output{Code: 500, Message: "failed to serialize JSON output"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	output.Message = string(msgBytes)

	ctx.SetOutputObject(output)
	return true, nil
}
