package query

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"github.com/project-flogo/core/activity"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-query")

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

	if input.Query == "" {
		log.Error("query statement is not specified\n")
		output := &Output{Code: 400, Message: "query statement is not specified"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	log.Debugf("query statement: %s\n", input.Query)

	// extract query parameter types from queryParams schema
	var paramTypes map[string]string
	if input.QueryParams != nil {
		schema, err := common.GetActivityInputSchema(ctx, "queryParams")
		if err != nil {
			log.Error("schema not defined for queryParams\n")
			output := &Output{Code: 400, Message: "schema not defined for queryParams"}
			ctx.SetOutputObject(output)
			return false, errors.New(output.Message)
		}
		if paramTypes, err = getQueryParamTypes(schema); err != nil {
			log.Errorf("failed to parse parameter schema %+v\n", err)
			output := &Output{Code: 400, Message: "failed to parse parameter schema"}
			ctx.SetOutputObject(output)
			return false, errors.Wrap(err, output.Message)
		}
	}
	log.Debugf("query parameters: %+v\n", input.QueryParams)

	queryStatement, err := prepareQueryStatement(input.Query, input.QueryParams, paramTypes)
	if err != nil {
		output := &Output{Code: 400, Message: err.Error()}
		ctx.SetOutputObject(output)
		return false, err
	}
	log.Debugf("query statement: %s\n", queryStatement)

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		log.Errorf("failed to retrieve fabric stub: %+v\n", err)
		output := &Output{Code: 500, Message: err.Error()}
		ctx.SetOutputObject(output)
		return false, err
	}

	if input.PrivateCollection != "" {
		// query private data
		return queryPrivateData(ctx, stub, queryStatement, input)
	}

	// query state data
	return queryData(ctx, stub, queryStatement, input)
}

func getQueryParamTypes(metadata string) (map[string]string, error) {
	// extract object field name and type from JSON schema
	var objectProps struct {
		Props map[string]struct {
			FieldType string `json:"type"`
		} `json:"properties"`
	}
	if err := json.Unmarshal([]byte(metadata), &objectProps); err != nil {
		log.Errorf("failed to extract properties from metadata: %+v", err)
		return nil, err
	}
	if objectProps.Props == nil {
		log.Debug("no parameter specified in metadata %s\n", metadata)
		return nil, nil
	}

	// collect property name and types
	params := make(map[string]string)
	for k, v := range objectProps.Props {
		log.Debugf("query parameter %s type %s\n", k, v.FieldType)
		params[k] = v.FieldType
	}
	return params, nil
}

func prepareQueryStatement(query string, queryParams map[string]interface{}, paramTypes map[string]string) (string, error) {
	if paramTypes == nil {
		log.Debug("no parameter is defined for query\n")
		return query, nil
	}

	if len(paramTypes) == 1 {
		for k := range paramTypes {
			// check if the single parameter is the query string
			pname := fmt.Sprintf(`"$%s"`, k)
			if pname == strings.TrimSpace(query) {
				log.Debugf("query statement is the first param: %v\n", queryParams[k])
				return fmt.Sprintf("%v", queryParams[k]), nil
			}
		}
	}

	// collect replacer args
	var args []string
	for pname, ptype := range paramTypes {
		value, ok := queryParams[pname]
		if !ok {
			// set default values
			switch ptype {
			case "number":
				value = 0
			case "boolean":
				value = false
			default:
				value = ""
			}
		}

		// collect string replacer args
		param := fmt.Sprintf("%v", value)
		if ptype == "string" {
			if jsonBytes, err := json.Marshal(value); err != nil {
				log.Debugf("failed to marshal value %v: %+v\n", value, err)
				param = "null"
			} else {
				param = string(jsonBytes)
			}
		}
		args = append(args, fmt.Sprintf(`"$%s"`, pname), param)
	}
	log.Debugf("query replacer args %v\n", args)

	// replace query parameters with values
	r := strings.NewReplacer(args...)
	return r.Replace(query), nil
}

func queryPrivateData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, query string, input *Input) (bool, error) {
	// check pagination
	pageSize := int32(0)
	bookmark := ""
	if input.UsePagination {
		pageSize = input.PageSize
		bookmark = input.Start
	}

	// query private data
	if pageSize > 0 {
		log.Infof("private data query does not support pagination, so ignore specified page size %d and bookmark %s\n", pageSize, bookmark)
	}
	resultIterator, err := ccshim.GetPrivateDataQueryResult(input.PrivateCollection, query)
	if err != nil {
		log.Errorf("failed to execute query %s on private collection %s: %+v\n", query, input.PrivateCollection, err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to execute query %s on private collection %s", query, input.PrivateCollection)}
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
		log.Infof("no data returned for query %s on private collection %s\n", query, input.PrivateCollection)
		output := &Output{Code: 300, Message: fmt.Sprintf("no data returned for query %s on private collection %s", query, input.PrivateCollection)}
		ctx.SetOutputObject(output)
		return true, nil
	}
	log.Debugf("query result from private collection %s: %s\n", input.PrivateCollection, string(jsonBytes))

	var value []interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to parse JSON data: %s", string(jsonBytes))}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	output := &Output{Code: 200,
		Message:  fmt.Sprintf("result of query %s from private collection %s: %s", query, input.PrivateCollection, string(jsonBytes)),
		Count:    len(value),
		Bookmark: "",
		Result:   value,
	}
	ctx.SetOutputObject(output)
	return true, nil
}

func queryData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, query string, input *Input) (bool, error) {
	// check pagination
	pageSize := int32(0)
	bookmark := ""
	if input.UsePagination {
		pageSize = input.PageSize
		log.Debug("pageSize:", pageSize)
		bookmark = input.Start
		log.Debug("bookmark:", bookmark)
	}

	// query state data
	var resultIterator shim.StateQueryIteratorInterface
	var resultMetadata *pb.QueryResponseMetadata
	var err error
	if pageSize > 0 {
		if resultIterator, resultMetadata, err = ccshim.GetQueryResultWithPagination(query, pageSize, bookmark); err != nil {
			log.Errorf("failed to execute query %s with page size %d: %+v\n", query, pageSize, err)
			output := &Output{Code: 500, Message: fmt.Sprintf("failed to execute query %s with page size %d", query, pageSize)}
			ctx.SetOutputObject(output)
			return false, errors.Wrapf(err, output.Message)
		}
	} else {
		if resultIterator, err = ccshim.GetQueryResult(query); err != nil {
			log.Errorf("failed to execute query %s: %+v\n", query, err)
			output := &Output{Code: 500, Message: fmt.Sprintf("failed to execute query %s", query)}
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
		log.Infof("no data returned for query %s\n", query)
		output := &Output{Code: 300, Message: fmt.Sprintf("no data returned for query %s", query)}
		ctx.SetOutputObject(output)
		return true, nil
	}
	log.Debugf("query returned data: %s\n", string(jsonBytes))

	var value []interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to parse JSON data: %s", string(jsonBytes))}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	if resultMetadata != nil {
		log.Debugf("set pagination metadata: count=%d, bookmark=%s\n", resultMetadata.FetchedRecordsCount, resultMetadata.Bookmark)
		// Note: bookmark of the last page is returned repeatedly for rich query
		// this may be a fabric bug. Pagination of range query returns blank when returned last page
		output := &Output{Code: 200,
			Message:  fmt.Sprintf("data returned for query %s: %s", query, string(jsonBytes)),
			Count:    int(resultMetadata.FetchedRecordsCount),
			Bookmark: resultMetadata.Bookmark,
			Result:   value,
		}
		ctx.SetOutputObject(output)
	} else {
		output := &Output{Code: 200,
			Message:  fmt.Sprintf("data returned for query %s: %s", query, string(jsonBytes)),
			Count:    len(value),
			Bookmark: "",
			Result:   value,
		}
		ctx.SetOutputObject(output)
	}
	return true, nil
}
