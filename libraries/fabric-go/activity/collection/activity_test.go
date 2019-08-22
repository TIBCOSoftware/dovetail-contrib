/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package collection

// Imports
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/stretchr/testify/assert"
)

// activityMetadata is the metadata of the activity as described in activity.json
// We'll store it as a variable to reuse it across multiple testcases
var activityMetadata *activity.Metadata

// getActivityMetadata reads the activity.json file and sets the activityMetadata variable
// if the variable already contains metadata it simply returns the current value rather than reading the file again
func getActivityMetadata() *activity.Metadata {
	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("../../../SmartContract/activity/collection/activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}
		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}
	return activityMetadata
}

// TestActivityRegistration checks whether the activity can be registered, and is registered in the engine
func TestActivityRegistration(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

// TestEval tests the Eval function and sends a message to IFTTT
// Make sure that you have updated the values below
func TestDistinctEval(t *testing.T) {
	fmt.Println("-----TestDistinct----")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	jsonstring := `[{"data":"USD"},{"data":"EUR"},{"data":"USD"}]`
	var input interface{}

	err := json.Unmarshal([]byte(jsonstring), &input)
	if err != nil {
		fmt.Printf("unmarshal error %v\n", err)
		panic(err)
	}

	inputObj := &data.ComplexObject{}
	inputObj.Value = input
	tc.SetInput(ivInput, inputObj)
	tc.SetInput(ivOperation, "DISTINCT")
	tc.SetInput(ivDatatype, "String")

	// Execute the activity
	_, err = act.Eval(tc)
	if err != nil {
		fmt.Printf("eval error %v\n", err)
		panic(err)
	}

	output := tc.GetOutput(ovOutput)
	results, err := data.CoerceToComplexObject(output)
	if err != nil {
		fmt.Printf("err = %v", err)
	} else {
		expected := make(map[string]interface{})
		expected["result"] = []interface{}{"USD", "EUR"}
		expected["count"] = 2
		assert.EqualValues(t, expected["count"], results.Value.(map[string]interface{})["count"])
	}
	// We assume there will be no errors
	assert.Nil(t, err)
}

func TestCountEval(t *testing.T) {
	fmt.Println("-----TestCount----")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	jsonstring := `[{"data":"USD"},{"data":"EUR"},{"data":"USD"}]`
	var input interface{}

	err := json.Unmarshal([]byte(jsonstring), &input)
	if err != nil {
		fmt.Printf("unmarshal error %v\n", err)
		panic(err)
	}

	inputObj := &data.ComplexObject{}
	inputObj.Value = input
	tc.SetInput(ivInput, inputObj)
	tc.SetInput(ivOperation, "COUNT")
	tc.SetInput(ivDatatype, "String")

	// Execute the activity
	_, err = act.Eval(tc)
	if err != nil {
		fmt.Printf("eval error %v\n", err)
		panic(err)
	}

	output := tc.GetOutput(ovOutput)
	results, err := data.CoerceToComplexObject(output)
	if err != nil {
		fmt.Printf("err = %v", err)
	} else {
		fmt.Printf("output=%v\n", results)
	}
	// We assume there will be no errors
	assert.Nil(t, err)
}

func TestSeqEval(t *testing.T) {
	fmt.Println("-----TestSequence----")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	jsonstring := `[{"data":"USD"},{"data":"EUR"},{"data":"USD"}]`
	var input interface{}

	err := json.Unmarshal([]byte(jsonstring), &input)
	if err != nil {
		fmt.Printf("unmarshal error %v\n", err)
		panic(err)
	}

	inputObj := &data.ComplexObject{}
	inputObj.Value = input
	tc.SetInput(ivUserInput, inputObj)
	tc.SetInput(ivOperation, "INDEXING")
	tc.SetInput(ivDatatype, "User Defined...")

	// Execute the activity
	_, err = act.Eval(tc)
	if err != nil {
		fmt.Printf("eval error %v\n", err)
		panic(err)
	}

	output := tc.GetOutput(ovOutput)
	results, err := data.CoerceToComplexObject(output)
	if err != nil {
		fmt.Printf("err = %v", err)
	} else {
		fmt.Printf("output=%v\n", results)
	}
	// We assume there will be no errors
	assert.Nil(t, err)
}

func TestNullEval(t *testing.T) {
	fmt.Println("-----TestNUll----")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	tc.SetInput(ivOperation, "DISTINCT")
	tc.SetInput(ivDatatype, "String")

	// Execute the activity
	_, err := act.Eval(tc)
	if err != nil {
		fmt.Printf("eval error %v\n", err)
		panic(err)
	}

	output := tc.GetOutput(ovOutput)
	results, err := data.CoerceToComplexObject(output)
	if err != nil {
		fmt.Printf("err = %v", err)
	} else {
		expected := make(map[string]interface{})
		expected["count"] = 0
		assert.EqualValues(t, expected["count"], results.Value.(map[string]interface{})["count"])
	}
	// We assume there will be no errors
	assert.Nil(t, err)
}

func TestJoinEval(t *testing.T) {
	fmt.Println("-----TestJoin----")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	jsonstring := `[{"field":"USD"},{"field":"EUR"},{"field":"USD"}]`
	var input interface{}

	err := json.Unmarshal([]byte(jsonstring), &input)
	if err != nil {
		fmt.Printf("unmarshal error %v\n", err)
		panic(err)
	}

	inputObj := &data.ComplexObject{}
	inputObj.Value = input
	tc.SetInput(ivInput, inputObj)
	tc.SetInput(ivOperation, "REDUCE-JOIN")
	tc.SetInput(ivDatatype, "String")
	tc.SetInput(ivDelimiter, ",")

	// Execute the activity
	_, err = act.Eval(tc)
	if err != nil {
		fmt.Printf("eval error %v\n", err)
		panic(err)
	}

	output := tc.GetOutput(ovOutput)
	results, err := data.CoerceToComplexObject(output)
	if err != nil {
		fmt.Printf("err = %v", err)
	} else {
		expected := "\"USD\",\"EUR\",\"USD\""
		assert.EqualValues(t, expected, results.Value.(map[string]interface{})["result"])
		fmt.Printf("output=%s\n", results.Value.(map[string]interface{})["result"])

	}
	// We assume there will be no errors
	assert.Nil(t, err)
}

func TestMergeEval(t *testing.T) {
	fmt.Println("-----TestMerge----")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	jsonstring1 := `{"input1":[{"field":"USD"},{"field":"EUR"}], "input2":[{"field":"YEN"}]}`
	jsonstring2 := `{"input1":[{"field":"USD"},{"field":"EUR"}], "input2":[]}`
	var input interface{}

	err := json.Unmarshal([]byte(jsonstring1), &input)
	if err != nil {
		fmt.Printf("unmarshal error %v\n", err)
		panic(err)
	}

	inputObj := &data.ComplexObject{}
	inputObj.Value = input
	tc.SetInput(ivInput, inputObj)
	tc.SetInput(ivOperation, "MERGE")

	// Execute the activity
	_, err = act.Eval(tc)
	if err != nil {
		fmt.Printf("eval error %v\n", err)
		panic(err)
	}

	output := tc.GetOutput(ovOutput)
	results, err := data.CoerceToComplexObject(output)
	if err != nil {
		fmt.Printf("err = %v", err)
	} else {
		v := results.Value.([]interface{})
		for _, vt := range v {
			fmt.Printf("output2=%v\n", vt)
		}
	}

	err = json.Unmarshal([]byte(jsonstring2), &input)
	inputObj.Value = input
	tc.SetInput(ivInput, inputObj)
	_, err = act.Eval(tc)
	output = tc.GetOutput(ovOutput)
	results, err = data.CoerceToComplexObject(output)
	fmt.Printf("output=%v\n", results.Value)
	// We assume there will be no errors
	assert.Nil(t, err)
}

func TestFilterEval(t *testing.T) {
	fmt.Println("-----TestFilter----")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	jsonstring1 := `{"dataset":[{"field":"USD"},{"field":"EUR"}], "filterValue":"USD", "filterField":"$dataset.field"}`
	var input interface{}

	err := json.Unmarshal([]byte(jsonstring1), &input)
	if err != nil {
		fmt.Printf("unmarshal error %v\n", err)
		panic(err)
	}

	inputObj := &data.ComplexObject{}
	inputObj.Value = input
	tc.SetInput(ivInput, inputObj)
	tc.SetInput(ivOperation, "FILTER")
	tc.SetInput(ivFilterType, "String")
	tc.SetInput(ivFilterOp, "==")

	// Execute the activity
	_, err = act.Eval(tc)
	if err != nil {
		fmt.Printf("eval error %v\n", err)
		panic(err)
	}

	output := tc.GetOutput(ovOutput)
	results, err := data.CoerceToComplexObject(output)
	if err != nil {
		fmt.Printf("err = %v", err)
	} else {
		fmt.Printf("output2=%v\n", results.Value)
	}

	jsonstring2 := `{"dataset":[{"field":"USD"},{"field":"EUR"}], "filterValue":"YEN", "filterField":"$dataset.field"}`
	err = json.Unmarshal([]byte(jsonstring2), &input)
	inputObj.Value = input
	tc.SetInput(ivInput, inputObj)
	_, err = act.Eval(tc)
	output = tc.GetOutput(ovOutput)
	results, err = data.CoerceToComplexObject(output)
	fmt.Printf("output=%v\n", results.Value)
	// We assume there will be no errors
	assert.Nil(t, err)
}
