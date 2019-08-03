package transaction

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings for the trigger
type Settings struct {
}

// HandlerSettings for the trigger
type HandlerSettings struct {
	Name       string `md:"name,required"`
	Validation bool   `md:"validation"`
}

// Output of the trigger
type Output struct {
	Parameters    map[string]interface{} `md:"parameters"`
	Transient     map[string]interface{} `md:"transient"`
	TxID          string                 `md:"txID"`
	TxTime        string                 `md:"txTime"`
	ChaincodeStub interface{}            `md:"_chaincode_stub"`
}

// Reply from the trigger
type Reply struct {
	Status  int         `md:"status"`
	Message string      `md:"message"`
	Returns interface{} `md:"returns"`
}

// FromMap sets trigger output values from a map
func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	if o.Parameters, err = coerce.ToObject(values["parameters"]); err != nil {
		return err
	}
	if o.Transient, err = coerce.ToObject(values["transient"]); err != nil {
		return err
	}
	if o.TxID, err = coerce.ToString(values["txID"]); err != nil {
		return err
	}
	if o.TxTime, err = coerce.ToString(values["txTime"]); err != nil {
		return err
	}
	if o.ChaincodeStub, err = coerce.ToAny(values["_chaincode_stub"]); err != nil {
		return err
	}

	return nil
}

// ToMap converts trigger output to a map
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"parameters":      o.Parameters,
		"transient":       o.Transient,
		"txID":            o.TxID,
		"txTime":          o.TxTime,
		"_chaincode_stub": o.ChaincodeStub,
	}
}

// FromMap sets trigger reply values from a map
func (r *Reply) FromMap(values map[string]interface{}) error {
	var err error
	if r.Status, err = coerce.ToInt(values["status"]); err != nil {
		return err
	}
	if r.Message, err = coerce.ToString(values["message"]); err != nil {
		return err
	}
	if r.Returns, err = coerce.ToAny(values["returns"]); err != nil {
		return err
	}
	return nil
}

// ToMap converts trigger reply to a map
func (r *Reply) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"status":  r.Status,
		"message": r.Message,
		"returns": r.Returns,
	}
}
