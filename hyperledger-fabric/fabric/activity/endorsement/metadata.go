package endorsement

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings of the activity
type Settings struct {
}

// Input of the activity
type Input struct {
	StateKey          string `md:"key,required"`
	Operation         string `md:"operation,required,allowed(ADD,DELETE,LIST,SET)"`
	Role              string `md:"role,allowed(MEMBER,ADMIN,CLIENT,PEER)"`
	Organizations     string `md:"organizations"`
	Policy            string `md:"policy"`
	IsPrivate         bool   `md:"isPrivate,required"`
	PrivateCollection string `md:"collection"`
}

// Output of the activity
type Output struct {
	Code     int    `md:"code"`
	Message  string `md:"message"`
	StateKey string `md:"key"`
	Result   string `md:"result"`
}

// ToMap converts activity input to a map
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"key":           i.StateKey,
		"operation":     i.Operation,
		"role":          i.Role,
		"organizations": i.Organizations,
		"policy":        i.Policy,
		"isPrivate":     i.IsPrivate,
		"collection":    i.PrivateCollection,
	}
}

// FromMap sets activity input values from a map
func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	if i.StateKey, err = coerce.ToString(values["key"]); err != nil {
		return err
	}
	if i.Operation, err = coerce.ToString(values["operation"]); err != nil {
		return err
	}
	if i.Role, err = coerce.ToString(values["role"]); err != nil {
		return err
	}
	if i.Organizations, err = coerce.ToString(values["organizations"]); err != nil {
		return err
	}
	if i.Policy, err = coerce.ToString(values["policy"]); err != nil {
		return err
	}
	if i.IsPrivate, err = coerce.ToBool(values["isPrivate"]); err != nil {
		return err
	}
	if i.PrivateCollection, err = coerce.ToString(values["collection"]); err != nil {
		return err
	}

	return nil
}

// ToMap converts activity output to a map
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code":    o.Code,
		"message": o.Message,
		"key":     o.StateKey,
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
	if o.StateKey, err = coerce.ToString(values["key"]); err != nil {
		return err
	}
	if o.Result, err = coerce.ToString(values["result"]); err != nil {
		return err
	}

	return nil
}
