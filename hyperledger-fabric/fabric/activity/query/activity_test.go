package query

import (
	"encoding/json"
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func TestCreate(t *testing.T) {

	mf := mapper.NewFactory(resolve.GetBasicResolver())
	iCtx := test.NewActivityInitContext(Settings{}, mf)
	act, err := New(iCtx)
	assert.Nil(t, err)
	assert.NotNil(t, act, "activity should not be nil")
}

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	mf := mapper.NewFactory(resolve.GetBasicResolver())
	iCtx := test.NewActivityInitContext(Settings{}, mf)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInputObject(&Input{Query: "$query", QueryParams: map[string]interface{}{"query": "{}"}})

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
