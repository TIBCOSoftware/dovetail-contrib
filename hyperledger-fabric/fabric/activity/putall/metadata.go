package putall

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings of the activity
type Settings struct {
}

// Input of the activity
type Input struct {
	StateData         []interface{} `md:"data,required"`
	PrivateCollection string        `md:"privateCollection"`
	CompositeKeys     string        `md:"compositeKeys"`
}

// Output of the activity
type Output struct {
	Code    int           `md:"code"`
	Message string        `md:"message"`
	Count   int           `md:"count"`
	Errors  int           `md:"errors"`
	Result  []interface{} `md:"result"`
}

// ToMap converts activity input to a map
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data":              i.StateData,
		"privateCollection": i.PrivateCollection,
		"compositeKeys":     i.CompositeKeys,
	}
}

// FromMap sets activity input values from a map
func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	if i.StateData, err = coerce.ToArray(values["data"]); err != nil {
		return err
	}
	if i.PrivateCollection, err = coerce.ToString(values["privateCollection"]); err != nil {
		return err
	}
	if i.CompositeKeys, err = coerce.ToString(values["compositeKeys"]); err != nil {
		return err
	}

	return nil
}

// ToMap converts activity output to a map
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code":    o.Code,
		"message": o.Message,
		"count":   o.Count,
		"errors":  o.Errors,
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
	if o.Count, err = coerce.ToInt(values["count"]); err != nil {
		return err
	}
	if o.Errors, err = coerce.ToInt(values["errors"]); err != nil {
		return err
	}
	if o.Result, err = coerce.ToArray(values["result"]); err != nil {
		return err
	}

	return nil
}
