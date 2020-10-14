package getbycompositekey

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings of the activity
type Settings struct {
}

// Input of the activity
type Input struct {
	KeyName           string                 `md:"keyName,required"`
	Attributes        map[string]interface{} `md:"attributes"`
	UsePagination     bool                   `md:"usePagination"`
	PageSize          int32                  `md:"pageSize"`
	Start             string                 `md:"start"`
	PrivateCollection string                 `md:"privateCollection"`
}

// Output of the activity
type Output struct {
	Code     int           `md:"code"`
	Message  string        `md:"message"`
	Bookmark string        `md:"bookmark"`
	Count    int           `md:"count"`
	Result   []interface{} `md:"result"`
}

// ToMap converts activity input to a map
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"keyName":           i.KeyName,
		"attributes":        i.Attributes,
		"usePagination":     i.UsePagination,
		"pageSize":          i.PageSize,
		"start":             i.Start,
		"privateCollection": i.PrivateCollection,
	}
}

// FromMap sets activity input values from a map
func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	if i.KeyName, err = coerce.ToString(values["keyName"]); err != nil {
		return err
	}
	if i.Attributes, err = coerce.ToObject(values["attributes"]); err != nil {
		return err
	}
	if i.UsePagination, err = coerce.ToBool(values["usePagination"]); err != nil {
		return err
	}
	if i.PageSize, err = coerce.ToInt32(values["pageSize"]); err != nil {
		return err
	}
	if i.Start, err = coerce.ToString(values["start"]); err != nil {
		return err
	}
	if i.PrivateCollection, err = coerce.ToString(values["privateCollection"]); err != nil {
		return err
	}

	return nil
}

// ToMap converts activity output to a map
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code":     o.Code,
		"message":  o.Message,
		"bookmark": o.Bookmark,
		"count":    o.Count,
		"result":   o.Result,
	}
}

// FromMap sets activity output values from a map
func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	if o.Code, err = coerce.ToInt(values["code"]); err != nil {
		return err
	}
	if o.Message, err = coerce.ToString(values["message"]); err != nil {
		o.Message = ""
	}
	if o.Bookmark, err = coerce.ToString(values["bookmark"]); err != nil {
		return err
	}
	if o.Count, err = coerce.ToInt(values["count"]); err != nil {
		return err
	}
	if o.Result, err = coerce.ToArray(values["result"]); err != nil {
		return err
	}

	return nil
}
