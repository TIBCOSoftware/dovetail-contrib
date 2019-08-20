/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package aggregate

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

var activityMetadata *activity.Metadata

// getActivityMetadata reads the activity.json file and sets the activityMetadata variable
// if the variable already contains metadata it simply returns the current value rather than reading the file again
func getActivityMetadata() *activity.Metadata {
	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("../../../SmartContract/activity/aggregate/activity.json")
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
func TestEval(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	jsonstring := `[{"data":100},{"data":50},{"data":50}]`
	var input interface{}

	err := json.Unmarshal([]byte(jsonstring), &input)
	if err != nil {
		fmt.Printf("unmarshal error %v\n", err)
		panic(err)
	}

	inputObj := &data.ComplexObject{}
	inputObj.Value = input
	tc.SetInput(ivInput, inputObj)
	tc.SetInput(ivOperation, "SUM")
	tc.SetInput(ivDatatype, "Long")

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
		fmt.Printf("results=%v\n", results.Value.(map[string]interface{})["result"].(int64))
		assert.EqualValues(t, 200, results.Value.(map[string]interface{})["result"].(int64))
	}

	// We assume there will be no errors
	assert.Nil(t, err)
}
