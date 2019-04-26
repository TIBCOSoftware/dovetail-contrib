package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
	trigger "github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/trigger/transaction"
)

const (
	fabricTrigger = "github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/trigger/transaction"
)

// Contract implements chaincode interface for invoking Flogo flows
type Contract struct {
}

var logger = shim.NewLogger("chaincode-shim")

// Init is called during chaincode instantiation to initialize any data,
// and also calls this function to reset or to migrate data.
func (t *Contract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode.
func (t *Contract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	logger.Debugf("invoke transaction fn=%s, args=%+v", fn, args)

	trig, ok := trigger.GetTrigger(fn)
	if !ok {
		return shim.Error(fmt.Sprintf("function %s is not implemented", fn))
	}
	result, err := trig.Invoke(stub, fn, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to execute transaction: %s, error: %+v", fn, err))
	}
	return shim.Success([]byte(result))
}

var cp app.ConfigProvider

// main function starts up the chaincode in the container during instantiate
func main() {
	common.SetChaincodeLogLevel(logger)

	// configure flogo engine
	if cp == nil {
		// Use default config provider
		cp = app.DefaultConfigProvider()
	}

	ac, err := cp.GetApp()
	if err != nil {
		fmt.Printf("failed to read Flogo app config: %+v\n", err)
		os.Exit(1)
	}

	// set mapping to pass fabric stub to activities in the flow
	// this is a workaround until flogo-lib can accept pass-through flow attributes in
	// handler.Handle(context.Background(), triggerData) that bypasses the mapper.
	// see issue: https://github.com/TIBCOSoftware/flogo-lib/issues/267
	inputAssignMap(ac, fabricTrigger, common.FabricStub)
	e, err := engine.New(ac)
	if err != nil {
		fmt.Printf("Failed to create flogo engine instance: %+v\n", err)
		os.Exit(1)
	}

	if err := e.Init(true); err != nil {
		fmt.Printf("Failed to initialize flogo engine: %+v\n", err)
		os.Exit(1)
	}

	if err := shim.Start(new(Contract)); err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

// inputAssignMap sets additional input mapping for a specified trigger ref type;
// this is to ensure the mapping of a trigger output property to avoid user error.
func inputAssignMap(ac *app.Config, triggerRef, name string) {
	// add the name to all flow resources
	prop := map[string]interface{}{"name": name, "type": "any"}
	for _, rc := range ac.Resources {
		var jsonobj map[string]interface{}
		if err := json.Unmarshal(rc.Data, &jsonobj); err != nil {
			logger.Errorf("failed to parse resource data %s: %+v", rc.ID, err)
			continue
		}
		if metadata, ok := jsonobj["metadata"]; ok {
			metaMap := metadata.(map[string]interface{})
			if input, ok := metaMap["input"]; ok {
				inputArray := input.([]interface{})
				done := false
				for _, ip := range inputArray {
					ipMap := ip.(map[string]interface{})
					if ipMap["name"].(string) == name {
						done = true
						continue
					}
				}
				if !done {
					logger.Debugf("add new property %s to resource input of %s", name, rc.ID)
					metaMap["input"] = append(inputArray, prop)
					if jsonbytes, err := json.Marshal(jsonobj); err == nil {
						logger.Debugf("resource data is updated for %s: %s", rc.ID, string(jsonbytes))
						rc.Data = jsonbytes
					} else {
						logger.Debugf("failed to serialize resource %s: %+v", rc.ID, err)
					}
				}
			}
		}
	}
	// add input mapper
	for _, tc := range ac.Triggers {
		if tc.Ref == triggerRef {
			for _, hc := range tc.Handlers {
				ivMap := hc.Action.Mappings.Input
				done := false
				for _, def := range ivMap {
					if def.MapTo == name {
						done = true
						continue
					}
				}
				if !done {
					hc.GetSetting(trigger.STransaction)
					logger.Infof("Add input mapper for %s to handler %+v", name, hc.GetSetting(trigger.STransaction))
					mapDef := data.MappingDef{Type: data.MtAssign, Value: "$." + name, MapTo: name}
					hc.Action.Mappings.Input = append(ivMap, &mapDef)
				}
			}
		}
	}
}
