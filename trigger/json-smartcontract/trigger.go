package transaction

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// log is the default package logger
var log = logger.GetLogger("trigger-smartcontract")

var singleton *ChaincodeTrigger

// ChaincodeTrigger Chaincode trigger struct
type ChaincodeTrigger struct {
	metadata     *trigger.Metadata
	config       *trigger.Config
	handlerInfos []*handlerInfo
	defHandler   *trigger.Handler
}

type handlerInfo struct {
	Invoke  bool
	handler *trigger.Handler
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	log.Debugf("Creating a new factory for id '%s'", md.ID)
	return &ChaincodeTriggerFactory{metadata: md}
}

// ChaincodeTriggerFactory Chaincode Trigger factory
type ChaincodeTriggerFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *ChaincodeTriggerFactory) New(config *trigger.Config) trigger.Trigger {
	log.Debugf("Creating a new singleton trigger")
	singleton = &ChaincodeTrigger{metadata: t.metadata, config: config}

	return singleton
}

// Metadata implements trigger.Trigger.Metadata
func (t *ChaincodeTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

func (t *ChaincodeTrigger) Initialize(ctx trigger.InitContext) error {
	log.SetLogLevel(logger.InfoLevel)

	if len(ctx.GetHandlers()) == 0 {
		return fmt.Errorf("no Handlers found for trigger '%s'", t.config.Id)
	}

	// Init handlers
	for _, handler := range ctx.GetHandlers() {
		aInfo := &handlerInfo{Invoke: false, handler: handler}

		t.handlerInfos = append(t.handlerInfos, aInfo)
	}

	return nil
}

// Start implements util.Managed.Start
func (t *ChaincodeTrigger) Start() error {
	return nil
}

// Stop implements util.Managed.Stop
func (t *ChaincodeTrigger) Stop() error {
	return nil
}

func Invoke(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	log.Debugf("Invoke called!!")
	if len(args) < 2 {
		return nil, fmt.Errorf("Wrong number of parameters passed, expected 2 minimum found %s" len(args))
	}
	// Get Handler to Invoke
	t := args[0]
	if len(t) == 0 {
		return nil, fmt.Errorf("No transaction chosen to start")
	}

	for _, info := range singleton.handlerInfos {

		if handlerT, set := info.handler.GetSetting("transaction"); set {
			if !set {
				continue
			}
			if t == handlerT {
				return singleton.Invoke(info.handler, stub, args)
			}
		}
	}

	return nil, fmt.Errorf("No transaction found to start '%s'", t)
}

func (t *ChaincodeTrigger) Invoke(handler *trigger.Handler, stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	log.Debugf("Trigger Invoke called!!")

	// Get Transaction Input
	ti := args[1]

	// Create transactionInput
	handlerOutput := handler.GetOutput()
	if handlerOutput == nil {
		err := fmt.Errorf("Error Invoking trigger: No handler output found")
		log.Debug(err.Error())
		return "", err
	}

	data := map[string]interface{}{
		"blockchainStub":   stub,
		"transactionInput": ti,
	}

	results, err := handler.Handle(context.Background(), data)

	if err != nil {
		log.Debugf("error: %s", err.Error())
		return "", err
	}

	log.Infof("Result Data: '%v'", results)
	var replyData interface{}

	if len(results) != 0 {
		dataAttr, ok := results["data"]
		if ok {
			replyData = dataAttr.Value()
		}
	}

	if replyData != nil {
		data, err := json.Marshal(replyData)
		if err != nil {
			return "", err
		}
		log.Infof("Returning Result Data: '%s'", string(data))
		return data, nil
	}
	log.Infof("Returning Result Data EMPTY")
	return []byte{}, nil
}
