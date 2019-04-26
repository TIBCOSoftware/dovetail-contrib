package transaction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	jschema "github.com/xeipuuv/gojsonschema"
	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
)

const (
	// STransaction is the name of the handler setting entry for transaction name
	STransaction = "name"
	sValidation  = "validation"
	oParameters  = "parameters"
	oTransient   = "transient"
	oTxID        = "txID"
	oTxTime      = "txTime"
	rReturns     = "returns"
)

// Create a new logger
var log = shim.NewLogger("trigger-fabric-transaction")

func init() {
	common.SetChaincodeLogLevel(log)
}

// TriggerMap maps transaction name in trigger handler setting to the trigger,
// so we can lookup trigger by transaction name
var triggerMap = map[string]*Trigger{}

// GetTrigger returns the cached trigger for a specified transaction name;
// return false in the second value if no trigger is cached for the specified name
func GetTrigger(name string) (*Trigger, bool) {
	trig, ok := triggerMap[name]
	return trig, ok
}

// TriggerFactory Fabric Trigger factory
type TriggerFactory struct {
	metadata *trigger.Metadata
}

// NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &TriggerFactory{metadata: md}
}

// New Creates a new trigger instance for a given id
func (t *TriggerFactory) New(config *trigger.Config) trigger.Trigger {
	return &Trigger{metadata: t.metadata, config: config, parameters: map[string][]common.ParameterIndex{}}
}

// Trigger is a stub for the Trigger implementation
type Trigger struct {
	metadata   *trigger.Metadata
	config     *trigger.Config
	handlers   []*trigger.Handler
	parameters map[string][]common.ParameterIndex
}

// Initialize implements trigger.Init.Initialize
func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	t.handlers = ctx.GetHandlers()
	for _, handler := range t.handlers {
		name := handler.GetStringSetting(STransaction)
		log.Info("init transaction trigger:", name)
		_, ok := triggerMap[name]
		if ok {
			log.Warningf("transaction name %s used by multiple trigger handlers, only the last handler is effective", name)
		}
		triggerMap[name] = t

		// collect input parameter name and types from metadata
		params, ok := handler.GetOutput()[oParameters].(*data.ComplexObject)
		if ok {
			// cache transaction parameters for each handler.
			// Note: Flogo enterprise uses one handler per flow, but share the same trigger instance
			if index, err := common.OrderedParameters([]byte(params.Metadata)); err == nil {
				if index != nil {
					log.Debugf("cache parameters for flow %s: %+v\n", name, index)
					t.parameters[name] = index
				}
			} else {
				log.Errorf("failed to initialize transaction parameters: %+v", err)
			}
		}

		// verify validation setting, value is not used
		handler.GetSetting(sValidation)
		validate := false
		if v, ok := handler.GetSetting(sValidation); ok {
			if bv, err := data.CoerceToBoolean(v); err == nil {
				validate = bv
			}
		}
		log.Info("validate output:", validate)
	}
	return nil
}

// Metadata implements trigger.Trigger.Metadata
func (t *Trigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *Trigger) Start() error {
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *Trigger) Stop() error {
	// stop the trigger
	return nil
}

// Invoke starts the trigger and invokes the action registered in the handler,
// and returns result as JSON string
func (t *Trigger) Invoke(stub shim.ChaincodeStubInterface, fn string, args []string) (string, error) {
	log.Debugf("fabric.Trigger invokes fn %s with args %+v", fn, args)

	for _, handler := range t.handlers {
		if f := handler.GetStringSetting(STransaction); f != fn {
			log.Debugf("skip handler for transaction %s that is different from requested function %s", f, fn)
			continue
		}
		triggerData := make(map[string]interface{})

		// construct transaction parameters
		paramIndex := t.parameters[fn]
		if paramIndex != nil && len(paramIndex) > 0 {
			paramData, err := prepareParameters(paramIndex, args)
			if err != nil {
				return "", err
			}
			if log.IsEnabledFor(shim.LogDebug) {
				// debug flow data
				paramBytes, _ := json.Marshal(paramData)
				log.Debugf("trigger parameters: %s", string(paramBytes))
			}

			// set trigger parameters
			params, _ := handler.GetOutput()[oParameters].(*data.ComplexObject)
			params.Value = paramData
			triggerData[oParameters] = params
		}

		// construct transient attributes
		if transMap, ok := handler.GetOutput()[oTransient].(*data.ComplexObject); ok && transMap != nil && transMap.Metadata != "" {
			transData, err := prepareTransient(stub)
			if err != nil {
				return "", err
			}
			if log.IsEnabledFor(shim.LogDebug) {
				// debug flow data
				transBytes, _ := json.Marshal(transData)
				log.Debugf("trigger transient attributes: %s", string(transBytes))
			}

			// set trigger transient attributes
			transMap.Value = transData
			triggerData[oTransient] = transMap
		}

		triggerData[common.FabricStub] = stub
		triggerData[oTxID] = stub.GetTxID()
		if ts, err := stub.GetTxTimestamp(); err == nil {
			triggerData[oTxTime] = time.Unix(ts.Seconds, int64(ts.Nanos)).UTC().Format(time.RFC3339Nano)
		}

		// execute flogo flow
		log.Debugf("flogo flow started transaction %s with timestamp %s", triggerData[oTxID], triggerData[oTxTime])
		results, err := handler.Handle(context.Background(), triggerData)
		if err != nil {
			log.Errorf("flogo flow returned error: %+v", err)
			return "", err
		}
		if len(results) != 0 {
			if dataAttr, ok := results[rReturns]; ok {
				// return serialized JSON string
				cobj := dataAttr.Value().(*data.ComplexObject)
				replyData, err := json.Marshal(cobj.Value)
				if err != nil {
					log.Errorf("failed to serialize reply: %+v", err)
					return "", err
				}
				log.Debugf("flogo flow returned data of type %T: %s", cobj.Value, string(replyData))
				return string(replyData), nil
			}
			log.Warningf("flogo flow result does not contain attribute %s", rReturns)
		}
		log.Info("flogo flow did not return any data")
		return "", nil
	}
	log.Warningf("no flogo handler is activated for transaction %s", fn)
	return "", nil
}

// construct trigger output transient attributes
func prepareTransient(stub shim.ChaincodeStubInterface) (map[string]interface{}, error) {
	transient := make(map[string]interface{})
	transMap, err := stub.GetTransient()
	if err != nil {
		// cannot find transient attributes
		log.Warningf("no transient map: %+v", err)
		return transient, nil
	}
	for k, v := range transMap {
		var obj interface{}
		if err := json.Unmarshal(v, &obj); err == nil {
			log.Debugf("received transient data, name: %s, value: %+v", k, obj)
			transient[k] = obj
		} else {
			log.Warningf("failed to unmarshal transient data, name: %s, error: %+v", k, err)
		}
	}
	return transient, nil
}

// construct trigger output parameters for specified parameter index, and values of the parameters
func prepareParameters(paramIndex []common.ParameterIndex, values []string) (interface{}, error) {
	log.Debugf("prepare parameters %+v values %+v", paramIndex, values)
	if paramIndex == nil && len(values) > 0 {
		// unknown parameter schema
		return nil, errors.New("parameter schema is not defined")
	}

	if len(paramIndex) < len(values) {
		// some data values are not defined by parameter index
		return nil, fmt.Errorf("parameter list %d is shorter than data items %d", len(paramIndex), len(values))
	}

	// convert string array to object with name-values as defined by parameter index
	result := make(map[string]interface{})
	if values != nil && len(values) > 0 {
		// populate input args
		for i, v := range values {
			if obj := unmarshalString(v, paramIndex[i].JSONType, paramIndex[i].Name); obj != nil {
				result[paramIndex[i].Name] = obj
			}
		}
	}
	return result, nil
}

// unmarshalString returns unmarshaled object if input is a valid JSON object or array,
// or returns the input string if it is not a valid JSON format
func unmarshalString(data, jsonType, name string) interface{} {
	s := strings.TrimSpace(data)
	switch jsonType {
	case jschema.TYPE_STRING:
		return s
	case jschema.TYPE_ARRAY:
		var result []interface{}
		if err := json.Unmarshal([]byte(data), &result); err != nil {
			log.Warningf("failed to parse parameter %s as JSON array: data '%s' error %+v", name, data, err)
		}
		return result
	case jschema.TYPE_BOOLEAN:
		b, err := strconv.ParseBool(s)
		if err != nil {
			log.Warningf("failed to convert parameter %s to boolean: data '%s' error %+v", name, data, err)
			return false
		}
		return b
	case jschema.TYPE_INTEGER:
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Warningf("failed to convert parameter %s to integer: data '%s' error %+v", name, data, err)
			return 0
		}
		return i
	case jschema.TYPE_NUMBER:
		if !strings.Contains(s, ".") {
			i, err := strconv.Atoi(s)
			if err != nil {
				log.Warningf("failed to convert parameter %s to integer: data '%s' error %+v", name, data, err)
				return 0
			}
			return i
		}
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Warningf("failed to convert parameter %s to float: data '%s' error %+v", name, data, err)
			return 0.0
		}
		return n
	default:
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(data), &result); err != nil {
			log.Warningf("failed to convert parameter %s to object: data '%s' error %+v", name, data, err)
		}
		return result
	}
}
