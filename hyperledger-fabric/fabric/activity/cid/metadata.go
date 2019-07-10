package cid

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings of the activity
type Settings struct {
}

// Input of the activity
type Input struct {
}

// Output of the activity
type Output struct {
	Code    int               `md:"code" json:"code"`
	Message string            `md:"message" json:"message,omitempty"`
	Cid     string            `md:"cid" json:"cid"`
	Mspid   string            `md:"mspid" json:"mspid"`
	Name    string            `md:"name" json:"name"`
	Attrs   map[string]string `md:"attrs" json:"attrs"`
}

// ToMap converts activity output to a map
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code":    o.Code,
		"message": o.Message,
		"cid":     o.Cid,
		"mspid":   o.Mspid,
		"name":    o.Name,
		"attrs":   o.Attrs,
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
	if o.Cid, err = coerce.ToString(values["cid"]); err != nil {
		return err
	}
	if o.Mspid, err = coerce.ToString(values["mspid"]); err != nil {
		return err
	}
	if o.Name, err = coerce.ToString(values["name"]); err != nil {
		return err
	}
	if o.Attrs, err = coerce.ToParams(values["attrs"]); err != nil {
		return err
	}
	return nil
}
