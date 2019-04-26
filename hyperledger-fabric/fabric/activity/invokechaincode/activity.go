package invokechaincode

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
	ivChaincode   = "chaincodeName"
	ivChannel     = "channelID"
	ivTransaction = "transactionName"
	ivParameters  = "parameters"
	ovCode        = "code"
	ovMessage     = "message"
	ovResult      = "result"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-invokechaincode")

func init() {
	common.SetChaincodeLogLevel(log)
}

// FabricChaincodeActivity is a stub for executing Hyperledger Fabric invoke-chaincode operations
type FabricChaincodeActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new FabricChaincodeActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &FabricChaincodeActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *FabricChaincodeActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *FabricChaincodeActivity) Eval(ctx activity.Context) (done bool, err error) {
	// check input args
	ccName, ok := ctx.GetInput(ivChaincode).(string)
	if !ok || ccName == "" {
		log.Error("chaincode name is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "chaincode name is not specified")
		return false, errors.New("chaincode name is not specified")
	}
	log.Debugf("chaincode name: %s\n", ccName)
	channelID := ""
	if channelID, ok = ctx.GetInput(ivChannel).(string); !ok {
		log.Info("channel ID is not specified\n")
	}
	log.Debugf("channel ID: %s\n", channelID)

	// extract transaction name and parameters
	args, err := constructChaincodeArgs(ctx)
	if err != nil {
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, err.Error())
		return false, err
	}

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, err.Error())
		return false, err
	}

	// invoke chaincode
	response := stub.InvokeChaincode(ccName, args, channelID)
	ctx.SetOutput(ovCode, response.GetStatus())
	ctx.SetOutput(ovMessage, response.GetMessage())
	jsonBytes := response.GetPayload()
	if jsonBytes == nil {
		log.Debugf("no data returned by invoking chaincode\n")
		return true, nil
	}
	var value interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to unmarshal chaincode response %+v, error: %+v\n", jsonBytes, err)
		return true, nil
	}
	if result, ok := ctx.GetOutput(ovResult).(*data.ComplexObject); ok && result != nil {
		log.Debugf("set activity output result: %+v\n", value)
		result.Value = value
		ctx.SetOutput(ovResult, result)
	}
	return true, nil
}

func constructChaincodeArgs(ctx activity.Context) ([][]byte, error) {
	var result [][]byte
	// transaction name from input
	txnName, ok := ctx.GetInput(ivTransaction).(string)
	if !ok || txnName == "" {
		log.Error("transaction name is not specified\n")
		return nil, errors.New("transaction name is not specified")
	}
	log.Debugf("transaction name: %s\n", txnName)
	result = append(result, []byte(txnName))

	// extract parameter definitions from metadata
	paramObj, ok := ctx.GetInput(ivParameters).(*data.ComplexObject)
	if !ok {
		log.Debug("parameter is not a complex object\n")
		return result, nil
	}
	paramIndex, err := common.OrderedParameters([]byte(paramObj.Metadata))
	if err != nil {
		log.Errorf("failed to extract parameter definition from metadata: %+v\n", err)
		return result, nil
	}
	if paramIndex == nil || len(paramIndex) == 0 {
		log.Debug("no parameter defined in metadata\n")
		return result, nil
	}

	// extract parameter values in the order of parameter index
	paramValue, ok := paramObj.Value.(map[string]interface{})
	if !ok {
		log.Debugf("parameter value of type %T is not a JSON object\n", paramObj.Value)
		return result, nil
	}
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
