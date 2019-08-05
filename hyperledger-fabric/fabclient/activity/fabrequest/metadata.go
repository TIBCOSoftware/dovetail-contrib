package fabrequest

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings of the activity
type Settings struct {
}

// Input of the activity
type Input struct {
	FabricConnector map[string]interface{} `md:"connectionName,required"`
	RequestType     string                 `md:"requestType,required"`
	OrgName         string                 `md:"orgName"`
	UserName        string                 `md:"userName,required"`
	ChaincodeID     string                 `md:"chaincodeID,required"`
	TransactionName string                 `md:"transactionName,required"`
	Parameters      map[string]interface{} `md:"parameters"`
	Transient       map[string]interface{} `md:"transient"`
}

// Output of the activity
type Output struct {
	Code    int         `md:"code"`
	Message string      `md:"message"`
	Result  interface{} `md:"result"`
}

// ToMap converts activity input to a map
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"connectionName":  i.FabricConnector,
		"requestType":     i.RequestType,
		"orgName":         i.OrgName,
		"userName":        i.UserName,
		"chaincodeID":     i.ChaincodeID,
		"transactionName": i.TransactionName,
		"parameters":      i.Parameters,
		"transient":       i.Transient,
	}
}

// FromMap sets activity input values from a map
func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	if i.FabricConnector, err = coerce.ToObject(values["connectionName"]); err != nil {
		return err
	}
	if i.RequestType, err = coerce.ToString(values["requestType"]); err != nil {
		return err
	}
	if i.OrgName, err = coerce.ToString(values["orgName"]); err != nil {
		return err
	}
	if i.UserName, err = coerce.ToString(values["userName"]); err != nil {
		return err
	}
	if i.ChaincodeID, err = coerce.ToString(values["chaincodeID"]); err != nil {
		return err
	}
	if i.TransactionName, err = coerce.ToString(values["transactionName"]); err != nil {
		return err
	}
	if i.Parameters, err = coerce.ToObject(values["parameters"]); err != nil {
		return err
	}
	if i.Transient, err = coerce.ToObject(values["transient"]); err != nil {
		return err
	}

	return nil
}

// ToMap converts activity output to a map
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code":    o.Code,
		"message": o.Message,
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
	if o.Result, err = coerce.ToAny(values["result"]); err != nil {
		return err
	}

	return nil
}
