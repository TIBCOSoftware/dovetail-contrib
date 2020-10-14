package getbycompositekey

import (
	"encoding/json"
	"fmt"

	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"github.com/project-flogo/core/activity"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-getbycompositekey")

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

	if input.KeyName == "" {
		log.Error("composite key name is not specified\n")
		output := &Output{Code: 400, Message: "composite key name is not specified"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	log.Debugf("composite key name: %s\n", input.KeyName)
	if input.Attributes == nil {
		log.Error("composite key attributes are not specified\n")
		output := &Output{Code: 400, Message: "composite key attributes are not specified"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	log.Debugf("composite key attributes: %+v\n", input.Attributes)

	// extract ordered list of attributes from JSON schema
	keyName, keyValues, err := extractCompositeKeyValues(input.KeyName, input.Attributes)
	if err != nil {
		log.Errorf("failed to extract composite key values: %+v\n", err)
		output := &Output{Code: 400, Message: err.Error()}
		ctx.SetOutputObject(output)
		return false, err
	}
	log.Debugf("Extracted composite key %s and values %+v", keyName, keyValues)

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		log.Errorf("failed to retrieve fabric stub: %+v\n", err)
		output := &Output{Code: 500, Message: err.Error()}
		ctx.SetOutputObject(output)
		return false, err
	}

	if input.PrivateCollection != "" {
		// retrieve data by composite key from a private collection
		return retrievePrivateDataByCompositeKey(ctx, stub, input, keyName, keyValues)
	}

	// retrieve data by composite key
	return retrieveByCompositeKey(ctx, stub, input, keyName, keyValues)
}

// extract composite key and values from input attributes
func extractCompositeKeyValues(keyDef string, attrs map[string]interface{}) (string, []string, error) {
	// parse composite key def into map[string][]string
	keys, err := common.ParseCompositeKeyDefs(keyDef)
	if err != nil {
		return "", nil, err
	}

	// extract key name and field values from input attributes
	name := ""
	var values []string
	for k, fields := range keys {
		name = k
		for _, f := range fields {
			if v, ok := attrs[f]; ok {
				values = append(values, fmt.Sprintf("%v", v))
			} else {
				// stop if attribute is not defined
				break
			}
		}
		if len(values) > 0 {
			// some key attributes are defined, so return the key and values
			return name, values, nil
		}
	}
	// no attribute value for any composite key, so return error
	return "", nil, errors.Errorf("No attribute value is defined for composite key %s", keyDef)
}

func retrievePrivateDataByCompositeKey(ctx activity.Context, ccshim shim.ChaincodeStubInterface, input *Input, name string, values []string) (bool, error) {
	// check pagination
	pageSize := int32(0)
	bookmark := ""
	if input.UsePagination {
		pageSize = input.PageSize
		bookmark = input.Start
	}

	// retrieve data from a private collection
	if input.PrivateCollection == "" {
		log.Error("private collection is not specified\n")
		output := &Output{Code: 400, Message: "private collection is not specified"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}

	// retrieve private data range [startKey, endKey]
	if pageSize > 0 {
		log.Infof("private data query does not support pagination, so ignore specified page size %d and bookmark %s\n", pageSize, bookmark)
	}
	resultIterator, err := ccshim.GetPrivateDataByPartialCompositeKey(input.PrivateCollection, name, values)
	if err != nil {
		log.Errorf("failed to retrieve by composite key %s from private collection %s: %+v\n", name, input.PrivateCollection, err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to retrieve by composite key %s from private collection %s", name, input.PrivateCollection)}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	defer resultIterator.Close()

	jsonBytes, err := common.ConstructQueryResponse(resultIterator, input.PrivateCollection, true, ccshim)
	if err != nil {
		log.Errorf("failed to collect result from iterator: %+v\n", err)
		output := &Output{Code: 500, Message: "failed to collect result from iterator"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	if jsonBytes == nil {
		log.Infof("no data found for composite key %s and value %+v from private collection %s\n", name, values, input.PrivateCollection)
		output := &Output{Code: 300, Message: fmt.Sprintf("no data found for composite key %s and value %+v from private collection %s", name, values, input.PrivateCollection)}
		ctx.SetOutputObject(output)
		return true, nil
	}
	log.Debugf("retrieved data from private collection %s: %s\n", input.PrivateCollection, string(jsonBytes))

	var value []interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to parse JSON data: %s", string(jsonBytes))}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	//	log.Debugf("result %d values %+v", len(value), value)

	output := &Output{Code: 200,
		Message:  fmt.Sprintf("retrieved data for composite key %s from private collection %s: %s", name, input.PrivateCollection, string(jsonBytes)),
		Count:    len(value),
		Bookmark: "",
		Result:   value,
	}
	ctx.SetOutputObject(output)
	return true, nil
}

func retrieveByCompositeKey(ctx activity.Context, ccshim shim.ChaincodeStubInterface, input *Input, name string, values []string) (bool, error) {
	// check pagination
	pageSize := int32(0)
	bookmark := ""
	if input.UsePagination {
		pageSize = input.PageSize
		log.Debug("pageSize:", pageSize)
		bookmark = input.Start
		log.Debug("bookmark:", bookmark)
	}

	// retrieve data by composite key
	var resultIterator shim.StateQueryIteratorInterface
	var resultMetadata *pb.QueryResponseMetadata
	var err error
	if pageSize > 0 {
		if resultIterator, resultMetadata, err = ccshim.GetStateByPartialCompositeKeyWithPagination(name, values, pageSize, bookmark); err != nil {
			log.Errorf("failed to retrieve by compsite key %s with page size %d: %+v\n", name, pageSize, err)
			output := &Output{Code: 500, Message: fmt.Sprintf("failed to retrieve by composite key %s with page size %d", name, pageSize)}
			ctx.SetOutputObject(output)
			return false, errors.Wrapf(err, output.Message)
		}
	} else {
		if resultIterator, err = ccshim.GetStateByPartialCompositeKey(name, values); err != nil {
			log.Errorf("failed to retrieve by composite key %s: %+v\n", name, err)
			output := &Output{Code: 500, Message: fmt.Sprintf("failed to retrieve by composite key %s", name)}
			ctx.SetOutputObject(output)
			return false, errors.Wrapf(err, output.Message)
		}
	}
	defer resultIterator.Close()

	jsonBytes, err := common.ConstructQueryResponse(resultIterator, "", true, ccshim)
	if err != nil {
		log.Errorf("failed to collect result from iterator: %+v\n", err)
		output := &Output{Code: 500, Message: "failed to collect result from iterator"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	if jsonBytes == nil {
		log.Infof("no data found for composite key %s value %+v\n", name, values)
		output := &Output{Code: 300, Message: fmt.Sprintf("no data found for composite key %s value %+v", name, values)}
		ctx.SetOutputObject(output)
		return true, nil
	}
	log.Debugf("retrieved data from ledger: %s\n", string(jsonBytes))

	var value []interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to parse JSON data: %s", string(jsonBytes))}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	if resultMetadata != nil {
		log.Debugf("set pagination metadata: count=%d, bookmark=%s\n", resultMetadata.FetchedRecordsCount, resultMetadata.Bookmark)
		output := &Output{Code: 200,
			Message:  fmt.Sprintf("retrieved data for composite key %s value %+v: %s", name, values, string(jsonBytes)),
			Count:    int(resultMetadata.FetchedRecordsCount),
			Bookmark: resultMetadata.Bookmark,
			Result:   value,
		}
		ctx.SetOutputObject(output)
	} else {
		output := &Output{Code: 200,
			Message:  fmt.Sprintf("retrieved data for composite key %s value %+v: %s", name, values, string(jsonBytes)),
			Count:    len(value),
			Bookmark: "",
			Result:   value,
		}
		ctx.SetOutputObject(output)
	}

	return true, nil
}
