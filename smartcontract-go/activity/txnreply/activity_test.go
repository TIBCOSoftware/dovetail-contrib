/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package txnreply

// Imports
import (
	"encoding/json"
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
		jsonMetadataBytes, err := ioutil.ReadFile("../../../SmartContract/activity/txnreply/activity.json")
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

	tc.SetInput(ivStatus, SUCCESS_WITH_DATA)

	inputmp := make(map[string]interface{})
	inputmp["Value"] = `{"func":"txnreply"}`
	v, err := json.Marshal(inputmp)
	complexObject := &data.ComplexObject{}
	err = json.Unmarshal([]byte(v), complexObject)
	tc.SetInput(ivData, complexObject)

	// Execute the activity
	_, err = act.Eval(tc)

	// We assume there will be no errors
	assert.Nil(t, err)
}

func TestEvalError(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	tc.SetInput(ivStatus, ERROR)
	tc.SetInput(ivMessage, "this is testing")
	// Execute the activity
	_, err := act.Eval(tc)

	// We assume there will be no errors
	assert.Nil(t, err)
}
