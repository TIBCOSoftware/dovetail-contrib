/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package transform

// Imports
import (
	"fmt"
	"strings"

	hcmath "github.com/TIBCOSoftware/dovetail-contrib/smartcontract-go/runtime/functions/math"
	"github.com/TIBCOSoftware/dovetail-contrib/smartcontract-go/utils"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/shopspring/decimal"
)

// Constants
const (
	ivInput           = "input"
	ivUserInput       = "userInput"
	ivIsArray         = "isArray"
	ivDatatype        = "dataType"
	ivPrecsion        = "precision"
	ivScale           = "scale"
	ivRounding        = "rounding"
	ivDatetimeFormat  = "format"
	ivInputArrayType  = "inputArrayType"
	ivOutputArrayType = "outputArrayType"
	ovOutput          = "output"
)

// describes the metadata of the activity as found in the activity.json file
type MapperActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MapperActivity{metadata: metadata}
}

func (a *MapperActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *MapperActivity) Eval(context activity.Context) (done bool, err error) {
	datatype := context.GetInput(ivDatatype).(string)
	rwinput := context.GetInput(ivInput)

	if datatype == "User Defined..." {
		rwinput = context.GetInput(ivUserInput)
	}
	if rwinput == nil {
		return false, fmt.Errorf("Input is not mapped")
	}

	//	precision := int32(context.GetInput(ivPrecsion).(int))
	//	scale := int32(context.GetInput(ivScale).(int))
	//	rounding := context.GetInput(ivRounding).(string)
	//format := context.GetInput(ivDatimeFormat).(string)
	isArray := context.GetInput(ivIsArray).(bool)
	if utils.IsPrimitive(datatype) && isArray {
		inputArrayType := context.GetInput(ivInputArrayType).(string)
		outputArrayType := context.GetInput(ivOutputArrayType).(string)

		var objs []interface{}
		if inputArrayType == "Object Array" {
			objs, err = utils.GetInputData(rwinput.(*data.ComplexObject), isArray)
			if err != nil {
				return false, err
			}
		} else {
			objs = getPrimitiveArrayInput(rwinput.(*data.ComplexObject))
		}

		if objs == nil || len(objs) == 0 {
			return false, fmt.Errorf("input is not mapped")
		}

		complexObj := data.ComplexObject{}
		if inputArrayType != outputArrayType {
			if outputArrayType == "Object Array" {
				output := make([]interface{}, 0)
				for _, v := range objs {
					m := make(map[string]interface{})
					m["field"] = v
					output = append(output, m)
				}

				complexObj.Value = output
			} else {
				output := make([]interface{}, 0)
				for _, v := range objs {
					m := v.(map[string]interface{})
					output = append(output, m["field"])
				}
				complexObj.Value = output
			}
			context.SetOutput(ovOutput, &complexObj)
		} else {
			context.SetOutput(ovOutput, rwinput)
		}
	} else {
		context.SetOutput(ovOutput, rwinput)
	}

	return true, nil
}

func setOutput(context activity.Context, arg interface{}) {
	output := &data.ComplexObject{}
	result := make(map[string]interface{})
	result["field"] = arg
	output.Value = result

	context.SetOutput(ovOutput, output)
}

func handleDoublValue(value interface{}, precision int32, scale int32, rounding string) (string, error) {
	dec, err := decimal.NewFromString(value.(string))
	if err != nil {
		return "", err
	}

	return hcmath.FormatDecimal(dec, precision, scale, rounding), nil
}

func getPrimitiveArrayInput(rwinput *data.ComplexObject) []interface{} {
	objs := make([]interface{}, 0)

	switch t := rwinput.Value.(type) {
	case []string:
		for _, v := range t {
			objs = append(objs, strings.TrimSpace(v))
		}
	case []int:
		for _, v := range t {
			objs = append(objs, v)
		}
	case []bool:
		for _, v := range t {
			objs = append(objs, v)
		}
	case []float64:
		for _, v := range t {
			objs = append(objs, v)
		}
	default:
		ta := t.([]interface{})
		for _, v := range ta {
			objs = append(objs, v)
		}
	}

	return objs
}
