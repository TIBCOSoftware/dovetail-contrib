package dovetail

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

var c = &ConditionalChoice{}

func TestConditionalChoiceString(t *testing.T) {
	// return second value if condition is true
	v, err := c.Eval(true, "foo", "bar", "opt")
	assert.NoError(t, err, "conditionalChoice(true, \"foo\", \"bar\", \"opt\") should not return error %v", err)
	b, ok := v.(string)
	assert.True(t, ok, "conditionalChoice(true, \"foo\", \"bar\", \"opt\") should return string, not type %T", v)
	assert.Equal(t, "bar", b, "conditionalChoice(true, \"foo\", \"bar\", \"opt\") should return 'bar', not %s", b)

	// return third value if condition is false and different 1st and 2nd values
	v, err = c.Eval(false, "foo", "bar", "opt")
	assert.NoError(t, err, "conditionalChoice(false, \"foo\", \"bar\", \"opt\") should not return error %v", err)
	b, ok = v.(string)
	assert.True(t, ok, "conditionalChoice(false, \"foo\", \"bar\", \"opt\") should return string, not type %T", v)
	assert.Equal(t, "opt", b, "conditionalChoice(false, \"foo\", \"bar\", \"opt\") should return 'opt', not %s", b)

	// return second value if condition is false and same 1st and 2nd values
	v, err = c.Eval(false, "foo", "foo", "opt")
	assert.NoError(t, err, "conditionalChoice(false, \"foo\", \"foo\", \"opt\") should not return error %v", err)
	b, ok = v.(string)
	assert.True(t, ok, "conditionalChoice(false, \"foo\", \"foo\", \"opt\") should return string, not type %T", v)
	assert.Equal(t, "foo", b, "conditionalChoice(false, \"foo\", \"foo\", \"opt\") should return 'foo', not %s", b)
}

func TestConditionalChoiceNumber(t *testing.T) {
	// return second value if condition is true
	v, err := c.Eval(true, 10, 20, 30)
	assert.NoError(t, err, "conditionalChoice(true, 10, 20, 30) should not return error %v", err)
	b, ok := v.(int)
	assert.True(t, ok, "conditionalChoice(true, 10, 20, 30) should return int, not type %T", v)
	assert.Equal(t, 20, b, "conditionalChoice(true, 10, 20, 30) should return 20, not %s", b)

	// return third value if condition is false and different 1st and 2nd values
	v, err = c.Eval(false, 10, 20, 30)
	assert.NoError(t, err, "conditionalChoice(false, 10, 20, 30) should not return error %v", err)
	b, ok = v.(int)
	assert.True(t, ok, "conditionalChoice(false, 10, 20, 30) should return int, not type %T", v)
	assert.Equal(t, 30, b, "conditionalChoice(false, 10, 20, 30) should return 30, not %s", b)

	// return second value if condition is false and same 1st and 2nd values
	v, err = c.Eval(false, 10, 10, 30)
	assert.NoError(t, err, "conditionalChoice(false, 10, 10, 30) should not return error %v", err)
	b, ok = v.(int)
	assert.True(t, ok, "conditionalChoice(false, 10, 10, 30) should return int, not type %T", v)
	assert.Equal(t, 10, b, "conditionalChoice(false, 10, 10, 30) should return 20, not %s", b)
}

func TestConditionalChoiceObject(t *testing.T) {
	var o1, o2, o3 interface{}
	err := json.Unmarshal([]byte(`{"foo": "f"}`), &o1)
	assert.NoError(t, err, "json Unmarshal of o1 should not return error %v", err)
	err = json.Unmarshal([]byte(`{"bar": "b"}`), &o2)
	assert.NoError(t, err, "json Unmarshal of o2 should not return error %v", err)
	err = json.Unmarshal([]byte(`{"opt": "p"}`), &o3)
	assert.NoError(t, err, "json Unmarshal of o3 should not return error %v", err)

	// return second value if condition is true
	v, err := c.Eval(true, o1, o2, o3)
	assert.NoError(t, err, "conditionalChoice(true, o1, o2, o3) should not return error %v", err)
	b, ok := v.(map[string]interface{})
	assert.True(t, ok, "conditionalChoice(true, o1, o2, o3) should return map, not type %T", v)
	assert.Equal(t, "b", b["bar"], "conditionalChoice(true, o1, o2, o3) should contain value 'b', not %v", b)
	
	// return third value if condition is false and different 1st and 2nd values
	v, err = c.Eval(false, o1, o2, o3)
	assert.NoError(t, err, "conditionalChoice(false, o1, o2, o3) should not return error %v", err)
	b, ok = v.(map[string]interface{})
	assert.True(t, ok, "conditionalChoice(false, o1, o2, o3) should return map, not type %T", v)
	assert.Equal(t, "p", b["opt"], "conditionalChoice(false, o1, o2, o3) should contain value 'p', not %v", b)

	// return second value if condition is false and same 1st and 2nd values
	v, err = c.Eval(false, o1, o1, o3)
	assert.NoError(t, err, "conditionalChoice(false, o1, o1, o3) should not return error %v", err)
	b, ok = v.(map[string]interface{})
	assert.True(t, ok, "conditionalChoice(false, o1, o1, o3) should return map, not type %T", v)
	assert.Equal(t, "f", b["foo"], "conditionalChoice(false, o1, o1, o3) should contain value 'f', not %v", b)
}