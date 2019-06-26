package gethistory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/project-flogo/core/activity"
)

const (
	dTxID      = "txID"
	dTxTime    = "txTime"
	dIsDeleted = "isDeleted"
	dValue     = "value"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-gethistory")

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

	// retrieve data for the key
	return retrieveHistory(ctx, stub, input.StateKey)
}

func retrieveHistory(ctx activity.Context, ccshim shim.ChaincodeStubInterface, key string) (bool, error) {
	// retrieve data for the key
	resultsIterator, err := ccshim.GetHistoryForKey(key)
	if err != nil {
		log.Errorf("failed to retrieve history for key %s: %+v\n", key, err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to retrieve history for key %s", key)}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	defer resultsIterator.Close()

	jsonBytes, err := constructHistoryResponse(resultsIterator)
	if jsonBytes == nil {
		log.Infof("no history found for key %s\n", key)
		output := &Output{Code: 300,
			Message:  fmt.Sprintf("no history found for key %s", key),
			StateKey: key,
		}
		ctx.SetOutputObject(output)
		return true, nil
	}
	log.Debugf("retrieved history from ledger: %s\n", string(jsonBytes))

	var value []interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to parse JSON data: %s", string(jsonBytes))}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	output := &Output{
		Code:     200,
		Message:  fmt.Sprintf("retrieved history for key %s: %s", key, string(jsonBytes)),
		StateKey: key,
		Result:   value,
	}
	ctx.SetOutputObject(output)
	return true, nil
}

func constructHistoryResponse(resultsIterator shim.HistoryQueryIteratorInterface) ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString("[")

	isEmpty := true
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if !isEmpty {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"" + dTxID + "\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"" + dValue + "\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"" + dTxTime + "\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).UTC().Format(time.RFC3339Nano))
		buffer.WriteString("\"")

		buffer.WriteString(", \"" + dIsDeleted + "\":")
		//		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		//		buffer.WriteString("\"")

		buffer.WriteString("}")
		isEmpty = false
	}
	buffer.WriteString("]")
	return buffer.Bytes(), nil
}
