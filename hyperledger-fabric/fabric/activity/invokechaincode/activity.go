package invokechaincode

import (
	"encoding/json"
	"fmt"

	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/project-flogo/core/activity"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-invokechaincode")

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

	if input.ChaincodeName == "" {
		log.Error("chaincode name is not specified\n")
		output := &Output{Code: 400, Message: "chaincode name is not specified"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	log.Debugf("chaincode name: %s\n", input.ChaincodeName)
	log.Debugf("channel ID: %s\n", input.ChannelID)

	// extract transaction name and parameters
	args, err := constructChaincodeArgs(ctx, input)
	if err != nil {
		output := &Output{Code: 400, Message: err.Error()}
		ctx.SetOutputObject(output)
		return false, err
	}

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		log.Errorf("failed to retrieve fabric stub: %+v\n", err)
		output := &Output{Code: 500, Message: err.Error()}
		ctx.SetOutputObject(output)
		return false, err
	}

	// invoke chaincode
	response := stub.InvokeChaincode(input.ChaincodeName, args, input.ChannelID)
	output := &Output{Code: int(response.GetStatus()), Message: response.GetMessage()}
	jsonBytes := response.GetPayload()
	if jsonBytes == nil {
		log.Debugf("no data returned by invoking chaincode\n")
		ctx.SetOutputObject(output)
		return true, nil
	}
	var value interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to unmarshal chaincode response %+v, error: %+v\n", jsonBytes, err)
		ctx.SetOutputObject(output)
		return true, nil
	}
	output.Result = value
	ctx.SetOutputObject(output)
	return true, nil
}

func constructChaincodeArgs(ctx activity.Context, input *Input) ([][]byte, error) {
	var result [][]byte
	// transaction name from input
	if input.TransactionName == "" {
		log.Error("transaction name is not specified\n")
		return nil, errors.New("transaction name is not specified")
	}
	log.Debugf("transaction name: %s\n", input.TransactionName)
	result = append(result, []byte(input.TransactionName))

	if input.Parameters == nil {
		log.Debug("no parameter is specified\n")
		return result, nil
	}

	// extract parameter definitions from metadata
	schema, err := common.GetActivityInputSchema(ctx, "parameters")
	if err != nil {
		log.Error("schema not defined for parameters\n")
		return nil, errors.New("schema not defined for parameters")
	}

	paramIndex, err := common.OrderedParameters([]byte(schema))
	if err != nil {
		log.Errorf("failed to extract parameter definition from schema: %+v\n", err)
		return result, nil
	}
	if paramIndex == nil || len(paramIndex) == 0 {
		log.Debug("no parameter defined in schema\n")
		return result, nil
	}

	// extract parameter values in the order of parameter index
	paramValue := input.Parameters
	for _, p := range paramIndex {
		// TODO: assuming string params here to be consistent with implementaton of trigger and chaincode-shim
		// should change all places to use []byte for best portability
		param := ""
		if v, ok := paramValue[p.Name]; ok && v != nil {
			param = fmt.Sprintf("%v", v)
			log.Debugf("add chaincode parameter: %s=%s", p.Name, param)
		}
		result = append(result, []byte(param))
	}
	return result, nil
}
