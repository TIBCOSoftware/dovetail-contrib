package putall

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/project-flogo/core/activity"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-putall")

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

	if input.StateData == nil {
		log.Errorf("input data is nil\n")
		output := &Output{Code: 400, Message: "input data is nil"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	log.Debugf("input value type %T: %+v\n", input.StateData, input.StateData)

	// get composite key definitions
	compositeKeyDefs, _ := getCompositeKeyDefinition(input.CompositeKeys)

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		log.Errorf("failed to retrieve fabric stub: %+v\n", err)
		output := &Output{Code: 500, Message: err.Error()}
		ctx.SetOutputObject(output)
	}

	var successCount, errorCount int
	var errorKeys []string
	var resultValue []interface{}
	if input.PrivateCollection != "" {
		// store data on a private collection
		for _, v := range input.StateData {
			vmap := v.(map[string]interface{})
			vkey := vmap[common.KeyField].(string)
			if err := storePrivateData(stub, input.PrivateCollection, compositeKeyDefs, vkey, vmap[common.ValueField]); err != nil {
				errorCount++
				errorKeys = append(errorKeys, vkey)
			} else {
				successCount++
				resultValue = append(resultValue, vmap)
			}
		}
	} else {
		// store data on the ledger
		for _, v := range input.StateData {
			vmap := v.(map[string]interface{})
			vkey := vmap[common.KeyField].(string)
			if err := storeData(stub, compositeKeyDefs, vkey, vmap[common.ValueField]); err != nil {
				errorCount++
				errorKeys = append(errorKeys, vkey)
			} else {
				successCount++
				resultValue = append(resultValue, vmap)
			}
		}
	}

	if errorCount > 0 {
		output := &Output{
			Code:    500,
			Message: fmt.Sprintf("failed to store keys: %s", strings.Join(errorKeys, ",")),
			Count:   successCount,
			Errors:  errorCount,
			Result:  resultValue,
		}
		if successCount > 0 {
			// return 300 if partial successs
			output.Code = 300
			ctx.SetOutputObject(output)
			return true, nil
		}
		// return 500 if all failures
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	// return 200 if no errors
	log.Debugf("set activity output result: %+v\n", resultValue)
	output := &Output{
		Code:    200,
		Message: fmt.Sprintf("stored data on ledger: %+v", resultValue),
		Count:   successCount,
		Errors:  errorCount,
		Result:  resultValue,
	}
	ctx.SetOutputObject(output)
	return true, nil
}

func storePrivateData(ccshim shim.ChaincodeStubInterface, collection string, compositeKeyDefs map[string][]string, key string, value interface{}) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		log.Errorf("failed to marshal value '%+v', error: %+v\n", value, err)
		return errors.Wrapf(err, "failed to marshal value: %+v", value)
	}

	// store data on a private collection
	if err := ccshim.PutPrivateData(collection, key, jsonBytes); err != nil {
		log.Errorf("failed to store data in private collection %s: %+v\n", collection, err)
		return errors.Wrapf(err, "failed to store data in private collection %s", collection)
	}
	log.Debugf("stored in private collection %s, data: %s\n", collection, string(jsonBytes))

	// store composite keys if required
	if compositeKeyDefs == nil {
		return nil
	}
	compositeKeys := common.ExtractCompositeKeys(ccshim, compositeKeyDefs, key, value)
	if compositeKeys != nil && len(compositeKeys) > 0 {
		for _, k := range compositeKeys {
			cv := []byte{0x00}
			if err := ccshim.PutPrivateData(collection, k, cv); err != nil {
				log.Errorf("failed to store composite key %s on collection %s: %+v\n", k, collection, err)
			} else {
				log.Debugf("stored composite key %s on collection %s\n", k, collection)
			}
		}
	}
	return nil
}

func storeData(ccshim shim.ChaincodeStubInterface, compositeKeyDefs map[string][]string, key string, value interface{}) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		log.Errorf("failed to marshal value '%+v', error: %+v\n", value, err)
		return errors.Wrapf(err, "failed to marshal value: %+v", value)
	}
	// store data on the ledger
	if err := ccshim.PutState(key, jsonBytes); err != nil {
		log.Errorf("failed to store data on ledger: %+v\n", err)
		return errors.Errorf("failed to store data on ledger: %+v", err)
	}
	log.Debugf("stored data on ledger: %s\n", string(jsonBytes))

	// store composite keys if required
	if compositeKeyDefs == nil {
		return nil
	}
	compositeKeys := common.ExtractCompositeKeys(ccshim, compositeKeyDefs, key, value)
	if compositeKeys != nil && len(compositeKeys) > 0 {
		for _, k := range compositeKeys {
			cv := []byte{0x00}
			if err := ccshim.PutState(k, cv); err != nil {
				log.Errorf("failed to store composite key %s: %+v\n", k, err)
			} else {
				log.Debugf("stored composite key %s\n", k)
			}
		}
	}
	return nil
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
