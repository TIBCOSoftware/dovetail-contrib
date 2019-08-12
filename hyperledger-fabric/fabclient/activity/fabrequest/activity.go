package fabrequest

import (
	"encoding/json"
	"fmt"
	"strings"

	client "github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabclient/common"
	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
	"github.com/pkg/errors"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/log"
)

const (
	conName          = "name"
	conConfig        = "config"
	conEntityMatcher = "entityMatcher"
	conChannel       = "channelID"
	opInvoke         = "invoke"
	opQuery          = "query"
)

// Create a new logger
var logger = log.ChildLogger(log.RootLogger(), "activity-fabclient-request")

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity fabric request activity struct
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

	if input.ChaincodeID == "" {
		logger.Error("chaincode ID is not specified")
		output := &Output{Code: 400, Message: "chaincode ID is not specified"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	logger.Debugf("chaincode ID: %s", input.ChaincodeID)

	if input.TransactionName == "" {
		logger.Error("transaction name is not specified")
		output := &Output{Code: 400, Message: "transaction name is not specified"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	logger.Debugf("transaction name: %s", input.TransactionName)

	reqType := input.RequestType
	if reqType == "" {
		logger.Warn("request type is not specified, assume `query`")
		reqType = opQuery
	}
	logger.Debugf("request type: %s", reqType)

	params, err := getParameters(ctx, input)
	if err != nil {
		output := &Output{Code: 400, Message: "invalid parameters"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	transientMap := getTransient(input.Transient)

	client, err := getFabricClient(input)
	if err != nil {
		output := &Output{Code: 500, Message: "fabric connector failure"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	// invoke fabric transaction
	var response []byte
	if reqType == opInvoke {
		logger.Debugf("execute chaincode %s transaction %s", input.ChaincodeID, input.TransactionName)
		response, err = client.ExecuteChaincode(input.ChaincodeID, input.TransactionName, params, transientMap)
	} else {
		logger.Debugf("query chaincode %s transaction %s", input.ChaincodeID, input.TransactionName)
		response, err = client.QueryChaincode(input.ChaincodeID, input.TransactionName, params, transientMap)
	}

	if err != nil {
		logger.Errorf("Fabric returned error %+v", err)
		output := &Output{Code: 500, Message: "Fabric request returned error"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	if response == nil {
		logger.Debugf("no data returned from fabric")
		output := &Output{Code: 300, Message: "no data returned from fabric"}
		ctx.SetOutputObject(output)
		return true, nil
	}
	logger.Debugf("Fabric response: %s\n", string(response))

	var value interface{}
	if err := json.Unmarshal(response, &value); err != nil {
		logger.Warnf("failed to unmarshal fabric response %+v, error: %+v", response, err)
		output := &Output{Code: 200,
			Message: fmt.Sprintf("data returned from fabric is not JSON: %s", string(response)),
			Result:  response,
		}
		ctx.SetOutputObject(output)
		return true, nil
	}
	output := &Output{Code: 200,
		Message: string(response),
		Result:  value,
	}
	ctx.SetOutputObject(output)
	return true, nil
}

func getFabricClient(input *Input) (*client.FabricClient, error) {
	if input.UserName == "" {
		logger.Error("user name is not specified")
		return nil, errors.New("user name is not specified")
	}

	configs, err := client.GetSettings(input.FabricConnector)
	if err != nil {
		return nil, err
	}
	networkConfig, err := client.ExtractFileContent(configs[conConfig])
	if err != nil {
		return nil, errors.Wrapf(err, "invalid network config")
	}
	entityMatcher, err := client.ExtractFileContent(configs[conEntityMatcher])
	if err != nil {
		return nil, errors.Wrapf(err, "invalid entity-matchers-override")
	}
	endpoints := []string{}
	if len(input.Endpoints) > 0 {
		endpoints = strings.Split(input.Endpoints, ",")
		for i, s := range endpoints {
			endpoints[i] = strings.TrimSpace(s)
		}
	}
	return client.NewFabricClient(client.ConnectorSpec{
		Name:           configs[conName].(string),
		NetworkConfig:  networkConfig,
		EntityMatchers: entityMatcher,
		OrgName:        input.OrgName,
		UserName:       input.UserName,
		ChannelID:      configs[conChannel].(string),
		TimeoutMillis:  input.TimeoutMillis,
		Endpoints:      endpoints,
	})
}

func getTransient(transData map[string]interface{}) map[string][]byte {
	if transData == nil {
		logger.Debug("no transient data is specified")
		return nil
	}
	transMap := make(map[string][]byte)
	for k, v := range transData {
		if jsonBytes, err := json.Marshal(v); err != nil {
			logger.Infof("failed to marshal transient data %+v", err)
		} else {
			transMap[k] = jsonBytes
		}
	}
	return transMap
}

func getParameters(ctx activity.Context, input *Input) ([][]byte, error) {
	var result [][]byte
	// extract parameter definitions from metadata
	if input.Parameters == nil {
		logger.Debug("no parameter is specified")
		return result, nil
	}

	schema, err := common.GetActivityInputSchema(ctx, "parameters")
	if err != nil {
		logger.Error("schema not defined for parameters\n")
		return nil, errors.New("schema not defined for parameters")
	}

	paramIndex, err := common.OrderedParameters([]byte(schema))
	if err != nil {
		logger.Errorf("failed to extract parameter definition from metadata: %+v\n", err)
		return result, nil
	}
	if paramIndex == nil || len(paramIndex) == 0 {
		logger.Debug("no parameter defined in metadata")
		return result, nil
	}

	// extract parameter values in the order of parameter index
	paramValue := input.Parameters
	for _, p := range paramIndex {
		// TODO: assuming string params here to be consistent with implementaton of trigger and chaincode-shim
		// should change all places to use []byte for best portability
		param := ""
		if v, ok := paramValue[p.Name]; ok && v != nil {
			if param, ok = v.(string); !ok {
				pbytes, err := json.Marshal(v)
				if err != nil {
					logger.Errorf("failed to marshal input: %+v", err)
					param = fmt.Sprintf("%v", v)
				} else {
					param = string(pbytes)
				}
			}
			logger.Infof("add chaincode parameter: %s=%s", p.Name, param)
		}
		result = append(result, []byte(param))
	}
	return result, nil
}
