package getrange

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
	ivStartKey      = "startKey"
	ivEndKey        = "endKey"
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
var log = shim.NewLogger("activity-fabric-getrange")

func init() {
	common.SetChaincodeLogLevel(log)
}

// FabricRangeActivity is a stub for executing Hyperledger Fabric get-by-range operations
type FabricRangeActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new FabricRangeActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &FabricRangeActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *FabricRangeActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *FabricRangeActivity) Eval(ctx activity.Context) (done bool, err error) {
	// check input args
	startKey, ok := ctx.GetInput(ivStartKey).(string)
	if !ok || startKey == "" {
		log.Error("start key is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "start key is not specified")
		return false, errors.New("start key is not specified")
	}
	log.Debugf("start key: %s\n", startKey)
	endKey, ok := ctx.GetInput(ivEndKey).(string)
	if !ok || endKey == "" {
		log.Error("end key is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "end key is not specified")
		return false, errors.New("end key is not specified")
	}
	log.Debugf("end key: %s\n", endKey)

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, err.Error())
		return false, err
	}

	if isPrivate, ok := ctx.GetInput(ivIsPrivate).(bool); ok && isPrivate {
		// retrieve data range from a private collection
		return retrievePrivateRange(ctx, stub, startKey, endKey)
	}

	// retrieve data range [startKey, endKey)
	return retrieveRange(ctx, stub, startKey, endKey)
}

func retrievePrivateRange(ctx activity.Context, ccshim shim.ChaincodeStubInterface, startKey, endKey string) (bool, error) {
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

	// retrieve private data range [startKey, endKey)
	if pageSize > 0 {
		log.Infof("private data query does not support pagination, so ignore specified page size %d and bookmark %s\n", pageSize, bookmark)
	}
	resultIterator, err := ccshim.GetPrivateDataByRange(collection, startKey, endKey)
	if err != nil {
		log.Errorf("failed to retrieve data range [%s, %s) from private collection %s: %+v\n", startKey, endKey, collection, err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to retrieve data range [%s, %s) from private collection %s: %+v", startKey, endKey, collection, err))
		return false, errors.Wrapf(err, "failed to retrieve data range [%s, %s) from private collection %s", startKey, endKey, collection)
	}
	defer resultIterator.Close()

	jsonBytes, err := common.ConstructQueryResponse(resultIterator, false, nil)
	if err != nil {
		log.Errorf("failed to collect result from iterator: %+v\n", err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to collect result from iterator: %+v", err))
		return false, errors.Wrapf(err, "failed to collect result from iterator")
	}

	if jsonBytes == nil {
		log.Infof("no data found in key range [%s, %s) from private collection %s\n", startKey, endKey, collection)
		ctx.SetOutput(ovCode, 300)
		ctx.SetOutput(ovMessage, fmt.Sprintf("no data found in key range [%s, %s) from private collection %s", startKey, endKey, collection))
		return true, nil
	}
	log.Debugf("retrieved data range from private collection %s: %s\n", collection, string(jsonBytes))

	var value interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to parse JSON data: %+v", err))
		return false, errors.Wrapf(err, "failed to parse JSON data %s", string(jsonBytes))
	}

	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("retrieved data in key range [%s, %s) from private collection %s: %s", startKey, endKey, collection, string(jsonBytes)))
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

func retrieveRange(ctx activity.Context, ccshim shim.ChaincodeStubInterface, startKey, endKey string) (bool, error) {
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

	// retrieve data range [startKey, endKey)
	var resultIterator shim.StateQueryIteratorInterface
	var resultMetadata *pb.QueryResponseMetadata
	var err error
	if pageSize > 0 {
		if resultIterator, resultMetadata, err = ccshim.GetStateByRangeWithPagination(startKey, endKey, pageSize, bookmark); err != nil {
			log.Errorf("failed to retrieve data range [%s, %s) with page size %d: %+v\n", startKey, endKey, pageSize, err)
			ctx.SetOutput(ovCode, 500)
			ctx.SetOutput(ovMessage, fmt.Sprintf("failed to retrieve data range [%s, %s) with page size %d: %+v", startKey, endKey, pageSize, err))
			return false, errors.Wrapf(err, "failed to retrieve data range [%s, %s) with page size %d", startKey, endKey, pageSize)
		}
	} else {
		if resultIterator, err = ccshim.GetStateByRange(startKey, endKey); err != nil {
			log.Errorf("failed to retrieve data range [%s, %s): %+v\n", startKey, endKey, err)
			ctx.SetOutput(ovCode, 500)
			ctx.SetOutput(ovMessage, fmt.Sprintf("failed to retrieve data range [%s, %s): %+v", startKey, endKey, err))
			return false, errors.Wrapf(err, "failed to retrieve data range [%s, %s)", startKey, endKey)
		}
	}
	defer resultIterator.Close()

	jsonBytes, err := common.ConstructQueryResponse(resultIterator, false, nil)
	if err != nil {
		log.Errorf("failed to collect result from iterator: %+v\n", err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to collect result from iterator: %+v", err))
		return false, errors.Wrapf(err, "failed to collect result from iterator")
	}

	if jsonBytes == nil {
		log.Infof("no data found in key range [%s, %s)\n", startKey, endKey)
		ctx.SetOutput(ovCode, 300)
		ctx.SetOutput(ovMessage, fmt.Sprintf("no data found in key range [%s, %s)", startKey, endKey))
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
	ctx.SetOutput(ovMessage, fmt.Sprintf("retrieved data in key range [%s, %s): %s", startKey, endKey, string(jsonBytes)))
	if result, ok := ctx.GetOutput(ovResult).(*data.ComplexObject); ok && result != nil {
		log.Debugf("set activity output result: %+v\n", value)
		result.Value = value
		ctx.SetOutput(ovResult, result)
		if resultMetadata != nil {
			// Returned bookmark is blank when reached last page
			log.Debugf("set count %v and bookmark %s\n", resultMetadata.FetchedRecordsCount, resultMetadata.Bookmark)
			ctx.SetOutput(ovCount, int(resultMetadata.FetchedRecordsCount))
			ctx.SetOutput(ovBookmark, resultMetadata.Bookmark)
		} else {
			ctx.SetOutput(ovBookmark, "")
			if vArray, ok := value.([]interface{}); ok {
				ctx.SetOutput(ovCount, len(vArray))
			}
		}
	}
	return true, nil
}
