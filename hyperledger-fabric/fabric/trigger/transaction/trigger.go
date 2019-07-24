package transaction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/trigger"

	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	jschema "github.com/xeipuuv/gojsonschema"
)

const (
	// STransaction is the name of the handler setting entry for transaction name
	STransaction = "name"
	oParameters  = "parameters"
	oTransient   = "transient"
)

// Create a new logger
var log = shim.NewLogger("trigger-fabric-transaction")

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
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

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{}, &Reply{})

// Factory for trigger
type Factory struct {
}

// New implements trigger.Factory.New
func (t *Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	if err := metadata.MapToStruct(config.Settings, s, true); err != nil {
		return nil, err
	}
	trig := Trigger{
		id:           config.Id,
		handlers:     map[string]trigger.Handler{},
		schemas:      map[string]*trigger.SchemaConfig{},
		parameters:   map[string][]common.ParameterIndex{},
		transientMap: map[string]bool{},
	}
	// extract parameters from handler schema
	for _, hc := range config.Handlers {
		name, ok := hc.Settings[STransaction].(string)
		if !ok {
			return nil, fmt.Errorf("Trigger handler setting does not contain attribute '%s'", STransaction)
		}
		log.Info("set schema config for handler:", name)
		trig.schemas[name] = hc.Schemas

		if _, ok := hc.Schemas.Output[oTransient]; ok {
			log.Info("set transient map for handler:", name)
			trig.transientMap[name] = true
		}

		if schema, ok := hc.Schemas.Output[oParameters].(map[string]interface{}); ok {
			log.Infof("schema config: %+v\n", schema)
			log.Infof("schema value: %T: %+v\n", schema["value"], schema["value"])
			if index, err := common.OrderedParameters([]byte(schema["value"].(string))); err == nil {
				if index != nil {
					log.Debugf("cache parameters for handler %s: %+v\n", name, index)
					trig.parameters[name] = index
				}
			} else {
				log.Errorf("failed to calculate handler parameters: %+v", err)
			}
		}
	}
	return &trig, nil
}

// Metadata implements trigger.Factory.Metadata
func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// Trigger is a stub for the Trigger implementation
type Trigger struct {
	id           string
	handlers     map[string]trigger.Handler
	schemas      map[string]*trigger.SchemaConfig
	parameters   map[string][]common.ParameterIndex
	transientMap map[string]bool
}

// Initialize implements trigger.Init.Initialize
func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	for _, handler := range ctx.GetHandlers() {
		setting := &HandlerSettings{}
		if err := metadata.MapToStruct(handler.Settings(), setting, true); err != nil {
			return err
		}

		name := setting.Name
		log.Info("init transaction trigger:", name)
		_, ok := triggerMap[name]
		if ok {
			log.Warningf("transaction name %s used by multiple trigger handlers, only the last handler is effective", name)
		}
		triggerMap[name] = t
		t.handlers[name] = handler

		// verify validation setting, value is not used
		validate := setting.Validation
		log.Info("validate output:", validate)
	}
	return nil
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
// and returns status code and result as JSON string
func (t *Trigger) Invoke(stub shim.ChaincodeStubInterface, fn string, args []string) (int, string, error) {
	log.Debugf("fabric.Trigger invokes fn %s with args %+v", fn, args)

	handler, ok := t.handlers[fn]
	if !ok {
		return 400, "", fmt.Errorf("Handler not defined for transaction %s", fn)
	}
	triggerData := &Output{}

	// construct transaction parameters
	paramIndex := t.parameters[fn]
	if paramIndex != nil && len(paramIndex) > 0 {
		paramData, err := prepareParameters(paramIndex, args)
		if err != nil {
			return 400, "", err
		}
		if log.IsEnabledFor(shim.LogDebug) {
			// debug flow data
			paramBytes, _ := json.Marshal(paramData)
			log.Debugf("trigger parameters: %s", string(paramBytes))
		}

		// set trigger parameters
		triggerData.Parameters = paramData
	}

	// construct transient attributes
	if trans, ok := t.transientMap[fn]; ok && trans {
		transData, err := prepareTransient(stub)
		if err != nil {
			return 400, "", err
		}
		if log.IsEnabledFor(shim.LogDebug) {
			// debug flow data
			transBytes, _ := json.Marshal(transData)
			log.Debugf("trigger transient attributes: %s", string(transBytes))
		}

		// set trigger transient attributes
		triggerData.Transient = transData
	}

	triggerData.ChaincodeStub = stub
	triggerData.TxID = stub.GetTxID()
	if ts, err := stub.GetTxTimestamp(); err == nil {
		triggerData.TxTime = time.Unix(ts.Seconds, int64(ts.Nanos)).UTC().Format(time.RFC3339Nano)
	}

	// execute flogo flow
	log.Debugf("flogo flow started transaction %s with timestamp %s", triggerData.TxID, triggerData.TxTime)
	results, err := handler.Handle(context.Background(), triggerData.ToMap())
	if err != nil {
		log.Errorf("flogo flow returned error: %+v", err)
		return 500, "", err
	}

	reply := &Reply{}
	if err := reply.FromMap(results); err != nil {
		return 500, "", err
	}

	if reply.Status != 200 {
		log.Infof("flogo flow returned status %d with message %s", reply.Status, reply.Message)
		return reply.Status, reply.Message, nil
	}
	if reply.Returns == nil {
		log.Info("flogo flow did not return any data")
		if reply.Message != "" {
			return 300, reply.Message, nil
		} else {
			return 300, "No data returned", nil
		}
	}

	replyData, err := json.Marshal(reply.Returns)
	if err != nil {
		log.Errorf("failed to serialize reply: %+v", err)
		return 500, "", err
	}
	log.Debugf("flogo flow returned data of type %T: %s", reply.Returns, string(replyData))
	return 200, string(replyData), nil
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
func prepareParameters(paramIndex []common.ParameterIndex, values []string) (map[string]interface{}, error) {
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
