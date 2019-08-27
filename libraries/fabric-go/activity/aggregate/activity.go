/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package aggregate

// Imports
import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"

	hcmath "github.com/TIBCOSoftware/dovetail-contrib/libraries/fabric-go/runtime/functions/math"
)

// Constants
const (
	ivInput     = "input"
	ivOperation = "operation"
	ivDatatype  = "dataType"
	ivPrecsion  = "precision"
	ivScale     = "scale"
	ivRounding  = "rounding"
	ovOutput    = "output"
)

// describes the metadata of the activity as found in the activity.json file
type AggregateActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AggregateActivity{metadata: metadata}
}

func (a *AggregateActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *AggregateActivity) Eval(context activity.Context) (done bool, err error) {
	operation := context.GetInput(ivOperation).(string)
	datatype := context.GetInput(ivDatatype).(string)
	rwinput := context.GetInput(ivInput)
	precision := int32(context.GetInput(ivPrecsion).(int))
	scale := int32(context.GetInput(ivScale).(int))
	rounding := context.GetInput(ivRounding).(string)

	if rwinput == nil {
		return true, nil
	}

	input, err := data.CoerceToComplexObject(rwinput)
	if err != nil {
		return false, err
	}

	value, ok := input.Value.([]interface{})
	if !ok {
		return false, fmt.Errorf("data error: can not case input data to []interface{}")
	}

	values := make([]interface{}, 0)
	for _, v := range value {
		values = append(values, v.(map[string]interface{})["data"])
	}

	switch operation {
	case "SUM":
		var sum interface{}
		switch datatype {
		case "Long":
			sum = hcmath.SumLong(values)
			break
		case "Integer":
			sum = hcmath.SumInt(values)
			break
		case "Double":
			sum, err = hcmath.SumDouble(values, precision, scale, rounding)
			if err != nil {
				return false, err
			}
			break
		default:
			return false, fmt.Errorf("Unsupported data type %s", datatype)
		}
		setOutput(context, sum)
		break
	case "AVG":
		var avg interface{}
		if len(values) == 0 {
			avg = 0
		} else {
			switch datatype {
			case "Long":
				avg = hcmath.AvgLong(values, precision, scale, rounding)
				break
			case "Integer":
				avg = hcmath.AvgInt(values, precision, scale, rounding)
				break
			case "Double":
				avg, err = hcmath.AvgDouble(values, precision, scale, rounding)
				if err != nil {
					return false, err
				}
				break
			default:
				return false, fmt.Errorf("Unsupported data type %s", datatype)
			}
		}

		setOutput(context, avg)
		break
	case "MAX":
		var max interface{}

		switch datatype {
		case "Long":
			max = hcmath.MaxLong(values)
			break
		case "Integer":
			max = hcmath.MaxInt(values)
			break
		case "Double":
			max, err = hcmath.MaxDouble(values, precision, scale, rounding)
			if err != nil {
				return false, err
			}
			break
		default:
			return false, fmt.Errorf("Unsupported data type %s", datatype)
		}

		setOutput(context, max)
		break
	case "MIN":
		var min interface{}

		switch datatype {
		case "Long":
			min = hcmath.MinLong(values)
			break
		case "Integer":
			min = hcmath.MinInt(values)
			break
		case "Double":
			min, err = hcmath.MinDouble(values, precision, scale, rounding)
			if err != nil {
				return false, err
			}
			break
		default:
			return false, fmt.Errorf("Unsupported data type %s", datatype)
		}

		setOutput(context, min)
		break
	default:
		return false, fmt.Errorf("operation %s is not supported", operation)
	}
	return true, nil
}

func setOutput(context activity.Context, arg interface{}) {
	output := &data.ComplexObject{}
	result := make(map[string]interface{})
	result["result"] = arg
	output.Value = result

	context.SetOutput(ovOutput, output)
}
