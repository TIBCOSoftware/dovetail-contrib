package getbycompositekey

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
)

const (
	ivKeyName       = "keyName"
	ivAttributes    = "attributes"
	ivUsePagination = "usePagination"
	ivPageSize      = "pageSize"
	ivBookmark      = "start"
	ivIsPrivate     = "isPrivate"
	ivCollection    = "collection"
	ovCode          = "code"
	ovMessage       = "message"
	ovBookmark      = "bookmark"
	ovCount         = "count"
	ovResult        = "result"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-getbycompositekey")

func init() {
	common.SetChaincodeLogLevel(log)
}

// FabricCompositeActivity is a stub for executing Hyperledger Fabric get-by-composite-key operations
type FabricCompositeActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new FabricCompositeActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &FabricCompositeActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *FabricCompositeActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *FabricCompositeActivity) Eval(ctx activity.Context) (done bool, err error) {
	// check input args
	keyName, ok := ctx.GetInput(ivKeyName).(string)
	if !ok || keyName == "" {
		log.Error("composite key name is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "composite key name is not specified")
		return false, errors.New("composite key name is not specified")
	}
	log.Debugf("composite key name: %s\n", keyName)
	attributes, ok := ctx.GetInput(ivAttributes).(*data.ComplexObject)
	if !ok || attributes == nil || attributes.Metadata == "" {
		log.Error("composite key attributes are not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "composite key attributes are not specified")
		return false, errors.New("composite key attributes are not specified")
	}
	log.Debugf("composite key attributes: %+v\n", attributes)

	// extract ordered list of attributes from JSON schema
	compIndex, err := common.OrderedParameters([]byte(attributes.Metadata))
	if err != nil {
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, err.Error())
		return false, err
	}
	if len(compIndex) == 0 {
		log.Error("composite key attribute list is empty\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "composite key attribute list is empty")
		return false, errors.New("composite key attribute list is empty")
	}

	// verify that attribute values are specified as JSON object
	attrValueMap, ok := attributes.Value.(map[string]interface{})
	if !ok {
		log.Errorf("invalid attribute value type %T, data: %+v\n", attributes.Value, attributes.Value)
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, fmt.Sprintf("invalid attribute value type %T", attributes.Value))
		return false, errors.Errorf("invalid attribute value type %T", attributes.Value)
	}

	// ordered composite key values, all attributes must be specified
	var compValues []string
	for _, v := range compIndex {
		attr, ok := attrValueMap[v.Name]
		if !ok {
			log.Errorf("composite key attribute %s is not specified\n", v.Name)
			ctx.SetOutput(ovCode, 400)
			ctx.SetOutput(ovMessage, fmt.Sprintf("composite key attribute %s is not specified", v.Name))
			return false, errors.Errorf("composite key attribute %s is not specified", v.Name)
		}
		compValues = append(compValues, fmt.Sprintf("%v", attr))
	}

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, err.Error())
		return false, err
	}

	if isPrivate, ok := ctx.GetInput(ivIsPrivate).(bool); ok && isPrivate {
		// retrieve data by composite key from a private collection
		return retrievePrivateDataByCompositeKey(ctx, stub, keyName, compValues)
	}

	// retrieve data by composite key
	return retrieveByCompositeKey(ctx, stub, keyName, compValues)
}

func retrievePrivateDataByCompositeKey(ctx activity.Context, ccshim shim.ChaincodeStubInterface, keyName string, values []string) (bool, error) {
	// check pagination
	pageSize := int32(0)
	bookmark := ""
	if usePagination, ok := ctx.GetInput(ivUsePagination).(bool); ok && usePagination {
		if f, err := strconv.ParseFloat(fmt.Sprintf("%v", ctx.GetInput(ivPageSize)), 64); err == nil {
			pageSize = int32(f)
			log.Debugf("pageSize=%d\n", pageSize)
		}
		if pageSize > 0 {
			if bookmark, ok = ctx.GetInput(ivBookmark).(string); ok && bookmark != "" {
				log.Debugf("bookmark=%s\n", bookmark)
			}
		}
	}

	// retrieve data from a private collection
	collection, ok := ctx.GetInput(ivCollection).(string)
	if !ok || collection == "" {
		log.Error("private collection is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "private collection is not specified")
		return false, errors.New("private collection is not specified")
	}

	// retrieve private data range [startKey, endKey]
	if pageSize > 0 {
		log.Infof("private data query does not support pagination, so ignore specified page size %d and bookmark %s\n", pageSize, bookmark)
	}
	resultIterator, err := ccshim.GetPrivateDataByPartialCompositeKey(collection, keyName, values)
	if err != nil {
		log.Errorf("failed to retrieve by composite key %s from private collection %s: %+v\n", keyName, collection, err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to retrieve by composite key %s from private collection %s: %+v\n", keyName, collection, err))
		return false, errors.Wrapf(err, "failed to retrieve by composite key %s from private collection %s", keyName, collection)
	}
	defer resultIterator.Close()

	jsonBytes, err := common.ConstructQueryResponse(resultIterator, true, ccshim)
	if err != nil {
		log.Errorf("failed to collect result from iterator: %+v\n", err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to collect result from iterator: %+v", err))
		return false, errors.Wrapf(err, "failed to collect result from iterator")
	}

	if jsonBytes == nil {
		log.Infof("no data found for composite key %s and value %+v from private collection %s\n", keyName, values, collection)
		ctx.SetOutput(ovCode, 300)
		ctx.SetOutput(ovMessage, fmt.Sprintf("no data found for composite key %s and value %+v from private collection %s\n", keyName, values, collection))
		return true, nil
	}
	log.Debugf("retrieved data from private collection %s: %s\n", collection, string(jsonBytes))

	var value interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to parse JSON data: %+v", err))
		return false, errors.Wrapf(err, "failed to parse JSON data %s", string(jsonBytes))
	}

	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("retrieved data for composite key %s from private collection %s: %s", keyName, collection, string(jsonBytes)))
	if result, ok := ctx.GetOutput(ovResult).(*data.ComplexObject); ok && result != nil {
		log.Debugf("set activity output result: %+v\n", value)
		result.Value = value
		ctx.SetOutput(ovResult, result)
		ctx.SetOutput(ovBookmark, "")
		if vArray, ok := value.([]interface{}); ok {
			ctx.SetOutput(ovCount, len(vArray))
		}
	}
	return true, nil
}

func retrieveByCompositeKey(ctx activity.Context, ccshim shim.ChaincodeStubInterface, keyName string, values []string) (bool, error) {
	// check pagination
	pageSize := int32(0)
	bookmark := ""
	if usePagination, ok := ctx.GetInput(ivUsePagination).(bool); ok && usePagination {
		if f, err := strconv.ParseFloat(fmt.Sprintf("%v", ctx.GetInput(ivPageSize)), 64); err == nil {
			pageSize = int32(f)
			log.Debugf("pageSize=%d\n", pageSize)
		}
		if pageSize > 0 {
			if bookmark, ok = ctx.GetInput(ivBookmark).(string); ok && bookmark != "" {
				log.Debugf("bookmark=%s\n", bookmark)
			}
		}
	}

	// retrieve data by composite key
	var resultIterator shim.StateQueryIteratorInterface
	var resultMetadata *pb.QueryResponseMetadata
	var err error
	if pageSize > 0 {
		if resultIterator, resultMetadata, err = ccshim.GetStateByPartialCompositeKeyWithPagination(keyName, values, pageSize, bookmark); err != nil {
			log.Errorf("failed to retrieve by compsite key %s with page size %d: %+v\n", keyName, pageSize, err)
			ctx.SetOutput(ovCode, 500)
			ctx.SetOutput(ovMessage, fmt.Sprintf("failed to retrieve by composite key %s with page size %d: %+v\n", keyName, pageSize, err))
			return false, errors.Wrapf(err, "failed to retrieve data by composite key %s with page size %d", keyName, pageSize)
		}
	} else {
		if resultIterator, err = ccshim.GetStateByPartialCompositeKey(keyName, values); err != nil {
			log.Errorf("failed to retrieve by composite key %s: %+v\n", keyName, err)
			ctx.SetOutput(ovCode, 500)
			ctx.SetOutput(ovMessage, fmt.Sprintf("failed to retrieve by composite key %s: %+v", keyName, err))
			return false, errors.Wrapf(err, "failed to retrieve by composite key %s", keyName)
		}
	}
	defer resultIterator.Close()

	jsonBytes, err := common.ConstructQueryResponse(resultIterator, true, ccshim)
	if err != nil {
		log.Errorf("failed to collect result from iterator: %+v\n", err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to collect result from iterator: %+v", err))
		return false, errors.Wrapf(err, "failed to collect result from iterator")
	}

	if jsonBytes == nil {
		log.Infof("no data found for composite key %s value %+v\n", keyName, values)
		ctx.SetOutput(ovCode, 300)
		ctx.SetOutput(ovMessage, fmt.Sprintf("no data found for composite key %s value %+v", keyName, values))
		return true, nil
	}
	log.Debugf("retrieved data from ledger: %s\n", string(jsonBytes))

	var value interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to parse JSON data: %+v", err))
		return false, errors.Wrapf(err, "failed to parse JSON data %s", string(jsonBytes))
	}

	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("retrieved data for composite key %s value %+v: %s", keyName, values, string(jsonBytes)))
	if result, ok := ctx.GetOutput(ovResult).(*data.ComplexObject); ok && result != nil {
		log.Debugf("set activity output result: %+v\n", value)
		result.Value = value
		ctx.SetOutput(ovResult, result)
		if resultMetadata != nil {
			log.Debugf("set pagination metadata: count=%d, bookmark=%s\n", resultMetadata.FetchedRecordsCount, resultMetadata.Bookmark)
			ctx.SetOutput(ovCount, int(resultMetadata.FetchedRecordsCount))
			ctx.SetOutput(ovBookmark, resultMetadata.Bookmark)
		} else {
			ctx.SetOutput(ovBookmark, "")
			if vArray, ok := value.([]interface{}); ok {
				log.Debugf("set value array lenth: \n", len(vArray))
				ctx.SetOutput(ovCount, len(vArray))
			} else {
				log.Debug("result value is not array\n")
			}
		}
	}
	return true, nil
}
