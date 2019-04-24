package query

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs

	act.Eval(tc)

	//check result attr
}

func TestGetQueryParamsType(t *testing.T) {
	metadata := `{
        "type": "object",
        "properties": {
            "sparam": {
                "type": "string"
            },
            "iparam": {
                "type": "number"
            },
            "bparam": {
                "type": "boolean"
            }
        },
        "required": []}`
	paramType, err := getQueryParamTypes(metadata)
	require.NoError(t, err, "failed to parse parameter metadata")
	assert.Equal(t, "string", paramType["sparam"], "sparam type should be string")
	assert.Equal(t, "number", paramType["iparam"], "iparam type should be number")
	assert.Equal(t, "boolean", paramType["bparam"], "bparam type should be boolean")
}

func TestPrepareQueryStatement(t *testing.T) {
	query := `{
        "selector": {
            "sParam": "$sparam",
            "iParam": {
                "$gt": "$iparam"
            },
            "bParam": "$bparam"
        }}`
	params := `{
        "sparam": "hello",
        "iparam": 100,
        "bparam": true}`
	types := `{
        "sparam": "string",
        "iparam": "number",
        "bparam": "boolean"}`
	result := `{
        "selector": {
            "sParam": "hello",
            "iParam": {
                "$gt": 100
            },
            "bParam": true
        }}`
	var queryParams map[string]interface{}
	err := json.Unmarshal([]byte(params), &queryParams)
	require.NoError(t, err, "failed to parse queryParams")
	var paramTypes map[string]string
	err = json.Unmarshal([]byte(types), &paramTypes)
	require.NoError(t, err, "failed to parse paramTypes")
	stmt, err := prepareQueryStatement(query, queryParams, paramTypes)
	require.NoError(t, err, "failed to prepare query statement")
	assert.Equal(t, result, stmt, "unexpected resulting query statement")
}
