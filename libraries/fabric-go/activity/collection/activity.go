/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package collection

// Imports
import (
	"bytes"
	"fmt"

	"github.com/TIBCOSoftware/dovetail-contrib/libraries/fabric-go/utils"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Constants
const (
	ivInput      = "input"
	ivUserInput  = "userInput"
	ivOperation  = "operation"
	ivDatatype   = "dataType"
	ivDelimiter  = "delimiter"
	ivFilterType = "filterFieldType"
	ivFilterOp   = "filterFieldOp"
	userDefined  = "User Defined..."
	ovOutput     = "output"
)

// describes the metadata of the activity as found in the activity.json file
type CollectionActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &CollectionActivity{metadata: metadata}
}

func (a *CollectionActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *CollectionActivity) Eval(context activity.Context) (done bool, err error) {
	operation := context.GetInput(ivOperation).(string)
	datatype := context.GetInput(ivDatatype).(string)
	rwinput := context.GetInput(ivInput)

	if datatype == userDefined {
		rwinput = context.GetInput(ivUserInput)
	}

	if rwinput == nil {
		result := make([]interface{}, 0)
		setOutput(operation, context, result)
		return true, nil
	}

	value, err := utils.GetInputData(rwinput.(*data.ComplexObject), true)
	if err != nil {
		return false, err
	}

	if value == nil || len(value) == 0 {
		setOutput(operation, context, value)
		return true, nil
	}

	switch operation {
	case "DISTINCT":
		set := make(map[string]interface{})
		results := make([]interface{}, 0)

		//unique values
		for _, m := range value {
			for _, v := range m.(map[string]interface{}) {
				set[v.(string)] = v
			}
		}

		for _, v := range set {
			results = append(results, v)
		}

		setOutput(operation, context, results)

		break
	case "COUNT":
		setOutput(operation, context, value)
		break
	case "INDEXING":
		for idx, v := range value {
			sm := v.(map[string]interface{})
			sm["_index_"] = idx
		}

		setOutput(operation, context, value)
		break
	case "REDUCE-JOIN":
		delimiter := context.GetInput(ivDelimiter).(string)
		result := bytes.Buffer{}
		for i, m := range value {
			for _, v := range m.(map[string]interface{}) {
				result.WriteString(fmt.Sprintf("\"%s\"", v.(string)))
				if i < (len(value) - 1) {
					result.WriteString(delimiter)
				}
			}
		}
		setOutput(operation, context, []interface{}{result.String()})
	case "FILTER":
		fobj := &data.ComplexObject{}
		result := make(map[string]interface{})
		filterIn := value[0].(map[string]interface{})
		dataset, ok := filterIn["dataset"].([]interface{})
		if !ok {
			result["trueSet"] = make([]interface{}, 0)
			result["falseSet"] = make([]interface{}, 0)
			context.SetOutput(ovOutput, result)
		} else {
			filterField := filterIn["filterField"].(string)
			filtervalue := filterIn["filterValue"]
			trueset := make([]interface{}, 0)
			falseset := make([]interface{}, 0)
			for _, datain := range dataset {
				datavalue, err := utils.FindValueInMap(datain.(map[string]interface{}), filterField)
				if err != nil {
					return false, err
				}

				if compare(filtervalue, datavalue, context.GetInput(ivFilterOp).(string), context.GetInput(ivFilterType).(string)) {
					trueset = append(trueset, datain)
				} else {
					falseset = append(falseset, datain)
				}
			}
			result["trueSet"] = trueset
			result["falseSet"] = falseset
			fobj.Value = result
			context.SetOutput(ovOutput, fobj)
		}
		break
	case "MERGE":
		mobj := &data.ComplexObject{}
		for _, input := range value {
			inputcolls := input.(map[string]interface{})
			input1, ok1 := inputcolls["input1"]
			input2, ok2 := inputcolls["input2"]

			if !ok1 || !ok2 {
				return false, fmt.Errorf("Invalid input for MERGE")
			}

			if input1 == nil && input2 == nil {
				mobj.Value = make([]interface{}, 0)
			}
			if input1 == nil {
				mobj.Value = input2
			} else if input2 == nil {
				mobj.Value = input1
			} else {
				output := make([]interface{}, 0)
				for _, in1 := range input1.([]interface{}) {
					output = append(output, in1)
				}
				for _, in2 := range input2.([]interface{}) {
					output = append(output, in2)
				}
				mobj.Value = output
			}
			context.SetOutput(ovOutput, mobj)
			break
		}
		break
	default:
		return false, fmt.Errorf("operation %s is not supported", operation)
	}

	return true, nil
}

func setOutput(action string, context activity.Context, arg []interface{}) {
	output := &data.ComplexObject{}

	if action == "DISTINCT" {
		result := make(map[string]interface{})
		if arg == nil {
			result["result"] = make([]string, 0)
			result["count"] = 0
		} else {
			result["result"] = arg
			result["count"] = len(arg)
		}
		output.Value = result
	} else if action == "COUNT" {
		result := make(map[string]interface{})
		if arg == nil {
			result["count"] = 0
		} else {
			result["count"] = len(arg)
		}
		output.Value = result
	} else if action == "INDEXING" {
		output.Value = arg
	} else if action == "REDUCE-JOIN" {
		result := make(map[string]interface{})
		if arg == nil {
			result["result"] = ""
		} else {
			result["result"] = arg[0].(string)
		}
		output.Value = result
	}

	context.SetOutput(ovOutput, output)
}

func compare(expected interface{}, actual interface{}, op, datatype string) bool {
	result := 0
	if op != "IN" {
		switch datatype {
		case "String":
			result = utils.StringCompare(expected.(string), actual.(string))
			break
		case "Integer", "Long":
			result = utils.IntCompare(expected.(int), actual.(int))
			break
		case "Boolean":
			return expected.(bool) == actual.(bool)
		}

		switch op {
		case "==":
			if result == 0 {
				return true
			} else {
				return false
			}
			break
		case "<":
			if result < 0 {
				return true
			} else {
				return false
			}
		case ">":
			if result > 0 {
				return true
			} else {
				return false
			}
		case "<=":
			if result <= 0 {
				return true
			} else {
				return false
			}
		case ">=":
			if result >= 0 {
				return true
			} else {
				return false
			}
		case "!=":
			if result != 0 {
				return true
			} else {
				return false
			}
		}
		return false
	} else {
		expectedarr := expected.([]interface{})
		for _, v := range expectedarr {
			r := compare(v, actual, "==", datatype)
			if r {
				return true
			}
		}
		return false
	}
}
