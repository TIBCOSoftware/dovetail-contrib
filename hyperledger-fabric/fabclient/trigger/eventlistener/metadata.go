package eventlistener

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings of the trigger
type Settings struct {
}

// HandlerSettings of the trigger
type HandlerSettings struct {
	FabricConnector map[string]interface{} `md:"connectionName,required"`
	EventType       string                 `md:"eventType,required,allowed(Block,Filtered Block,Chaincode)"`
	ChaincodeID     string                 `md:"chaincodeID"`
	EventFilter     string                 `md:"eventFilter"`
	Org             string                 `md:"org"`
	User            string                 `md:"user,required"`
	Validation      bool                   `md:"validation"`
}

// Output of the activity
type Output struct {
	Data interface{} `md:"data"`
}

// Reply of the trigger
type Reply struct {
}

// ToMap converts handler settings to a map
func (h *HandlerSettings) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"connectionName": h.FabricConnector,
		"eventType":      h.EventType,
		"chaincodeID":    h.ChaincodeID,
		"eventFilter":    h.EventFilter,
		"org":            h.Org,
		"user":           h.User,
		"validatioon":    h.Validation,
	}
}

// FromMap sets handler setting values from a map
func (h *HandlerSettings) FromMap(values map[string]interface{}) error {

	var err error
	if h.FabricConnector, err = coerce.ToObject(values["connectionName"]); err != nil {
		return err
	}
	if h.EventType, err = coerce.ToString(values["eventType"]); err != nil {
		return err
	}
	if h.ChaincodeID, err = coerce.ToString(values["chaincodeID"]); err != nil {
		return err
	}
	if h.EventFilter, err = coerce.ToString(values["eventFilter"]); err != nil {
		return err
	}
	if h.Org, err = coerce.ToString(values["org"]); err != nil {
		return err
	}
	if h.User, err = coerce.ToString(values["user"]); err != nil {
		return err
	}
	if h.Validation, err = coerce.ToBool(values["validation"]); err != nil {
		return err
	}

	return nil
}

// ToMap converts trigger output to a map
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data": o.Data,
	}
}

// FromMap sets trigger output values from a map
func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	if o.Data, err = coerce.ToAny(values["data"]); err != nil {
		return err
	}

	return nil
}
