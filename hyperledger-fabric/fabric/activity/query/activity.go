package query

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
)

const (
	ivQuery         = "query"
	ivQueryParams   = "queryParams"
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
var log = shim.NewLogger("activity-fabric-query")

func init() {
	common.SetChaincodeLogLevel(log)
}

// FabricQueryActivity is a stub for executing Hyperledger Fabric query operations
type FabricQueryActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new FabricQueryActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &FabricQueryActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *FabricQueryActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *FabricQueryActivity) Eval(ctx activity.Context) (done bool, err error) {
	// check input args
	query, ok := ctx.GetInput(ivQuery).(string)
	if !ok || query == "" {
		log.Error("query statement is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "query statement is not specified")
		return false, errors.New("query statement is not specified")
	}
	log.Debugf("query statement: %s\n", query)
	queryParams, paramTypes, err := getQueryParams(ctx)
	if err != nil {
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, err.Error())
		return false, err
	}
	log.Debugf("query parameters: %+v\n", queryParams)

	queryStatement, err := prepareQueryStatement(query, queryParams, paramTypes)
	if err != nil {
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, err.Error())
		return false, err
	}
	log.Debugf("query statement: %s\n", queryStatement)

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, err.Error())
		return false, err
	}

	if isPrivate, ok := ctx.GetInput(ivIsPrivate).(bool); ok && isPrivate {
		// query private data
		return queryPrivateData(ctx, stub, queryStatement)
	}

	// query state data
	return queryData(ctx, stub, queryStatement)
}

func getQueryParams(ctx activity.Context) (params map[string]interface{}, paramTypes map[string]string, err error) {
	queryParams, ok := ctx.GetInput(ivQueryParams).(*data.ComplexObject)
	if !ok || queryParams == nil || queryParams.Value == nil || queryParams.Metadata == "" {
		log.Debug("no query parameter is specified\n")
		return nil, nil, nil
	}

	// extract parameter definitions from metadata
	paramTypes, err = getQueryParamTypes(queryParams.Metadata)
	if err != nil {
		log.Errorf("failed to parse parameter metadata %+v\n", err)
		return nil, nil, err
	}
	if paramTypes == nil || len(paramTypes) == 0 {
		log.Debugf("no parameters defined in metadata\n")
		return nil, nil, nil
	}

	// verify parameter values
	params, ok = queryParams.Value.(map[string]interface{})
	if !ok {
		log.Errorf("query parameter type %T is not JSON object\n", queryParams.Value)
		return nil, nil, errors.Errorf("query parameter type %T is not JSON object\n", queryParams.Value)
	}
	return params, paramTypes, nil
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

func queryPrivateData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, query string) (bool, error) {
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

	// query data from a private collection
	collection, ok := ctx.GetInput(ivCollection).(string)
	if !ok || collection == "" {
		log.Error("private collection is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "private collection is not specified")
		return false, errors.New("private collection is not specified")
	}

	// query private data
	if pageSize > 0 {
		log.Infof("private data query does not support pagination, so ignore specified page size %d and bookmark %s\n", pageSize, bookmark)
	}
	resultIterator, err := ccshim.GetPrivateDataQueryResult(collection, query)
	if err != nil {
		log.Errorf("failed to execute query %s on private collection %s: %+v\n", query, collection, err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to execute query %s on private collection %s: %+v", query, collection, err))
		return false, errors.Wrapf(err, "failed to execute query %s on private collection %s", query, collection)
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
		log.Infof("no data returned for query %s on private collection %s\n", query, collection)
		ctx.SetOutput(ovCode, 300)
		ctx.SetOutput(ovMessage, fmt.Sprintf("no data returned for query %s on private collection %s", query, collection))
		return true, nil
	}
	log.Debugf("query result from private collection %s: %s\n", collection, string(jsonBytes))

	var value interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to parse JSON data: %+v", err))
		return false, errors.Wrapf(err, "failed to parse JSON data %s", string(jsonBytes))
	}

	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("result of query %s from private collection %s: %s", query, collection, string(jsonBytes)))
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

func queryData(ctx activity.Context, ccshim shim.ChaincodeStubInterface, query string) (bool, error) {
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

	// query state data
	var resultIterator shim.StateQueryIteratorInterface
	var resultMetadata *pb.QueryResponseMetadata
	var err error
	if pageSize > 0 {
		if resultIterator, resultMetadata, err = ccshim.GetQueryResultWithPagination(query, pageSize, bookmark); err != nil {
			log.Errorf("failed to execute query %s with page size %d: %+v\n", query, pageSize, err)
			ctx.SetOutput(ovCode, 500)
			ctx.SetOutput(ovMessage, fmt.Sprintf("failed to execute query %s with page size %d: %+v", query, pageSize, err))
			return false, errors.Wrapf(err, "failed to execute query %s with page size %d", query, pageSize)
		}
	} else {
		if resultIterator, err = ccshim.GetQueryResult(query); err != nil {
			log.Errorf("failed to execute query %s: %+v\n", query, err)
			ctx.SetOutput(ovCode, 500)
			ctx.SetOutput(ovMessage, fmt.Sprintf("failed to execute query %s: %+v", query, err))
			return false, errors.Wrapf(err, "failed to execute query %s", query)
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
		log.Infof("no data returned for query %s\n", query)
		ctx.SetOutput(ovCode, 300)
		ctx.SetOutput(ovMessage, fmt.Sprintf("no data returned for query %s", query))
		return true, nil
	}
	log.Debugf("query returned data: %s\n", string(jsonBytes))

	var value interface{}
	if err := json.Unmarshal(jsonBytes, &value); err != nil {
		log.Errorf("failed to parse JSON data: %+v\n", err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to parse JSON data: %+v", err))
		return false, errors.Wrapf(err, "failed to parse JSON data %s", string(jsonBytes))
	}

	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("data returned for query %s: %s", query, string(jsonBytes)))
	if result, ok := ctx.GetOutput(ovResult).(*data.ComplexObject); ok && result != nil {
		log.Debugf("set activity output result: %+v\n", value)
		result.Value = value
		ctx.SetOutput(ovResult, result)
		if resultMetadata != nil {
			// Note: bookmark of the last page is returned repeatedly for rich query
			// this may be a bug. Pagination of range query returns blank when returned last page
			log.Debugf("set pagination metadata: count=%d, bookmark=%s\n", resultMetadata.FetchedRecordsCount, resultMetadata.Bookmark)
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
