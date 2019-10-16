package getrange

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
var log = shim.NewLogger("activity-fabric-getrange")

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

	if input.StartKey == "" {
		output := &Output{Code: 400, Message: "start key is not specified"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	log.Debugf("start key: %s\n", input.StartKey)
	if input.EndKey == "" {
		output := &Output{Code: 400, Message: "end key is not specified"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	log.Debugf("end key: %s\n", input.EndKey)

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		log.Errorf("failed to retrieve fabric stub: %+v\n", err)
		output := &Output{Code: 500, Message: err.Error()}
		ctx.SetOutputObject(output)
		return false, err
	}

	if input.PrivateCollection != "" {
		// retrieve data range from a private collection
		return retrievePrivateRange(ctx, stub, input)
	}

	// retrieve data range [startKey, endKey)
	return retrieveRange(ctx, stub, input)
}

func retrievePrivateRange(ctx activity.Context, ccshim shim.ChaincodeStubInterface, input *Input) (bool, error) {
	// check pagination
	pageSize := int32(0)
	bookmark := ""
	if input.UsePagination {
		pageSize = input.PageSize
		bookmark = input.Start
	}

	// retrieve private data range [startKey, endKey)
	if pageSize > 0 {
		log.Infof("private data query does not support pagination, so ignore specified page size %d and bookmark %s\n", pageSize, bookmark)
	}
	resultIterator, err := ccshim.GetPrivateDataByRange(input.PrivateCollection, input.StartKey, input.EndKey)
	if err != nil {
		log.Errorf("failed to retrieve data range [%s, %s) from private collection %s: %+v\n", input.StartKey, input.EndKey, input.PrivateCollection, err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to retrieve data range [%s, %s) from private collection %s", input.StartKey, input.EndKey, input.PrivateCollection)}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	defer resultIterator.Close()

	jsonBytes, err := common.ConstructQueryResponse(resultIterator, input.PrivateCollection, false, nil)
	if err != nil {
		log.Errorf("failed to collect result from iterator: %+v\n", err)
		output := &Output{Code: 500, Message: "failed to collect result from iterator"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	if jsonBytes == nil {
		log.Infof("no data found in key range [%s, %s) from private collection %s\n", input.StartKey, input.EndKey, input.PrivateCollection)
		output := &Output{Code: 300, Message: fmt.Sprintf("no data found in key range [%s, %s) from private collection %s", input.StartKey, input.EndKey, input.PrivateCollection)}
		ctx.SetOutputObject(output)
		return true, nil
	}
	log.Debugf("retrieved data range from private collection %s: %s\n", input.PrivateCollection, string(jsonBytes))

	var value []interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to parse JSON data: %s", string(jsonBytes))}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	output := &Output{Code: 200,
		Message:  fmt.Sprintf("retrieved data in key range [%s, %s) from private collection %s: %s", input.StartKey, input.EndKey, input.PrivateCollection, string(jsonBytes)),
		Count:    len(value),
		Bookmark: "",
		Result:   value,
	}
	ctx.SetOutputObject(output)
	return true, nil
}

func retrieveRange(ctx activity.Context, ccshim shim.ChaincodeStubInterface, input *Input) (bool, error) {
	// check pagination
	pageSize := int32(0)
	bookmark := ""
	if input.UsePagination {
		pageSize = input.PageSize
		log.Debug("pageSize:", pageSize)
		bookmark = input.Start
		log.Debug("bookmark:", bookmark)
	}

	// retrieve data range [startKey, endKey)
	var resultIterator shim.StateQueryIteratorInterface
	var resultMetadata *pb.QueryResponseMetadata
	var err error
	if pageSize > 0 {
		if resultIterator, resultMetadata, err = ccshim.GetStateByRangeWithPagination(input.StartKey, input.EndKey, pageSize, bookmark); err != nil {
			log.Errorf("failed to retrieve data range [%s, %s) with page size %d: %+v\n", input.StartKey, input.EndKey, pageSize, err)
			output := &Output{Code: 500, Message: fmt.Sprintf("failed to retrieve data range [%s, %s) with page size %d", input.StartKey, input.EndKey, pageSize)}
			ctx.SetOutputObject(output)
			return false, errors.Wrapf(err, output.Message)
		}
	} else {
		if resultIterator, err = ccshim.GetStateByRange(input.StartKey, input.EndKey); err != nil {
			log.Errorf("failed to retrieve data range [%s, %s): %+v\n", input.StartKey, input.EndKey, err)
			output := &Output{Code: 500, Message: fmt.Sprintf("failed to retrieve data range [%s, %s)", input.StartKey, input.EndKey)}
			ctx.SetOutputObject(output)
			return false, errors.Wrapf(err, output.Message)
		}
	}
	defer resultIterator.Close()

	jsonBytes, err := common.ConstructQueryResponse(resultIterator, "", false, nil)
	if err != nil {
		log.Errorf("failed to collect result from iterator: %+v\n", err)
		output := &Output{Code: 500, Message: "failed to collect result from iterator"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	if jsonBytes == nil {
		log.Infof("no data found in key range [%s, %s)\n", input.StartKey, input.EndKey)
		output := &Output{Code: 300, Message: fmt.Sprintf("no data found in key range [%s, %s)", input.StartKey, input.EndKey)}
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
			Message:  fmt.Sprintf("retrieved data in key range [%s, %s): %s", input.StartKey, input.EndKey, string(jsonBytes)),
			Count:    int(resultMetadata.FetchedRecordsCount),
			Bookmark: resultMetadata.Bookmark,
			Result:   value,
		}
		ctx.SetOutputObject(output)
	} else {
		output := &Output{Code: 200,
			Message:  fmt.Sprintf("retrieved data in key range [%s, %s): %s", input.StartKey, input.EndKey, string(jsonBytes)),
			Count:    len(value),
			Bookmark: "",
			Result:   value,
		}
		ctx.SetOutputObject(output)
	}
	return true, nil
}
