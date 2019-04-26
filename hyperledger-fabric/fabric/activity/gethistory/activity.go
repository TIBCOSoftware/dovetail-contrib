package gethistory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
)

const (
	ivKey      = "key"
	ovCode     = "code"
	ovMessage  = "message"
	ovKey      = "key"
	ovResult   = "result"
	dTxID      = "txID"
	dTxTime    = "txTime"
	dIsDeleted = "isDeleted"
	dValue     = "value"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-gethistory")

func init() {
	common.SetChaincodeLogLevel(log)
}

// FabricHistoryActivity is a stub for executing Hyperledger Fabric get-history operations
type FabricHistoryActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new FabricHistoryActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &FabricHistoryActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *FabricHistoryActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *FabricHistoryActivity) Eval(ctx activity.Context) (done bool, err error) {
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

	// retrieve data for the key
	return retrieveHistory(ctx, stub, key)
}

func retrieveHistory(ctx activity.Context, ccshim shim.ChaincodeStubInterface, key string) (bool, error) {
	// retrieve data for the key
	resultsIterator, err := ccshim.GetHistoryForKey(key)
	if err != nil {
		log.Errorf("failed to retrieve history for key %s: %+v\n", key, err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to retrieve history for key %s: %+v", key, err))
		return false, errors.Wrapf(err, "failed to retrieve history for key %s", key)
	}
	defer resultsIterator.Close()

	jsonBytes, err := constructHistoryResponse(resultsIterator)
	if jsonBytes == nil {
		log.Infof("no history found for key %s\n", key)
		ctx.SetOutput(ovCode, 300)
		ctx.SetOutput(ovMessage, fmt.Sprintf("no history found for key %s", key))
		ctx.SetOutput(ovKey, key)
		return true, nil
	}
	log.Debugf("retrieved history from ledger: %s\n", string(jsonBytes))

	var value interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to parse JSON data: %+v", err))
		return false, errors.Wrapf(err, "failed to parse JSON data %s", string(jsonBytes))
	}

	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("retrieved history for key %s: %s", key, string(jsonBytes)))
	if result, ok := ctx.GetOutput(ovResult).(*data.ComplexObject); ok && result != nil {
		log.Debugf("set activity output result: %+v\n", value)
		result.Value = value
		ctx.SetOutput(ovResult, result)
		ctx.SetOutput(ovKey, key)
	}
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
