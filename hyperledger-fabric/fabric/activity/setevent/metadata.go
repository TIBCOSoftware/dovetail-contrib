package setevent

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings of the activity
type Settings struct {
}

// Input of the activity
type Input struct {
	Name    string      `md:"name,required"`
	Payload interface{} `md:"payload"`
}

// Output of the activity
type Output struct {
	Code    int         `md:"code"`
	Message string      `md:"message"`
	Name    string      `md:"name"`
	Result  interface{} `md:"result"`
}

// ToMap converts activity input to a map
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":    i.Name,
		"payload": i.Payload,
	}
}

// FromMap sets activity input values from a map
func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	if i.Name, err = coerce.ToString(values["name"]); err != nil {
		return err
	}
	if i.Payload, err = coerce.ToAny(values["payload"]); err != nil {
		return err
	}

	return nil
}

// ToMap converts activity output to a map
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code":    o.Code,
		"message": o.Message,
		"name":    o.Name,
		"result":  o.Result,
	}
}

// FromMap sets activity output values from a map
func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	if o.Code, err = coerce.ToInt(values["code"]); err != nil {
		return err
	}
	if o.Message, err = coerce.ToString(values["message"]); err != nil {
		return err
	}
	if o.Name, err = coerce.ToString(values["name"]); err != nil {
		return err
	}
	if o.Result, err = coerce.ToAny(values["result"]); err != nil {
		return err
	}

	return nil
}
