package dovetail

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

var z = &NotZeroValue{}

func TestNotZeroString(t *testing.T) {
	// return first value if second is blank
	v, err := z.Eval("foo", "")
	assert.NoError(t, err, "notZeroValue(\"foo\", \"\") should not return error %v", err)
	b, ok := v.(string)
	assert.True(t, ok, "notZeroValue(\"foo\", \"\") should return string, not type %T", v)
	assert.Equal(t, "foo", b, "notZeroValue(\"foo\", \"\") should return 'foo', not %s", b)

	// return first value if second is nil
	v, err = z.Eval("foo", nil)
	assert.NoError(t, err, "notZeroValue(\"foo\", nil) should not return error %v", err)
	b, ok = v.(string)
	assert.True(t, ok, "notZeroValue(\"foo\", nil) should return string, not type %T", v)
	assert.Equal(t, "foo", b, "notZeroValue(\"foo\", nil) should return 'foo', not %s", b)
	
	// return second value
	v, err = z.Eval("foo", "bar")
	assert.NoError(t, err, "notZeroValue(\"foo\", \"bar\") should not return error %v", err)
	b, ok = v.(string)
	assert.True(t, ok, "notZeroValue(\"foo\", \"bar\") should return string, not type %T", v)
	assert.Equal(t, "bar", b, "notZeroValue(\"foo\", \"bar\") should return 'bar', not %s", b)
}

func TestNotZeroNumber(t *testing.T) {
	// return first value if second is 0
	v, err := z.Eval(10, 0)
	assert.NoError(t, err, "notZeroValue(10, 0) should not return error %v", err)
	b, ok := v.(int)
	assert.True(t, ok, "notZeroValue(10, 0) should return int, not type %T", v)
	assert.Equal(t, 10, b, "notZeroValue(10, 0) should return 10, not %s", b)

	// return first value if second is nil
	v, err = z.Eval(10, nil)
	assert.NoError(t, err, "notZeroValue(10, nil) should not return error %v", err)
	b, ok = v.(int)
	assert.True(t, ok, "notZeroValue(10, nil) should return int, not type %T", v)
	assert.Equal(t, 10, b, "notZeroValue(10, nil) should return 10, not %s", b)
	
	// return second value
	v, err = z.Eval(10, 20)
	assert.NoError(t, err, "notZeroValue(10, 20) should not return error %v", err)
	b, ok = v.(int)
	assert.True(t, ok, "notZeroValue(10, 20) should return int, not type %T", v)
	assert.Equal(t, 20, b, "notZeroValue(10, 20) should return 20, not %s", b)
}

func TestNotZeroObject(t *testing.T) {
	var o1, o2 interface{}
	err := json.Unmarshal([]byte(`{"foo": "f"}`), &o1)
	assert.NoError(t, err, "json Unmarshal of o1 should not return error %v", err)
	err = json.Unmarshal([]byte(`{"bar": "b"}`), &o2)
	assert.NoError(t, err, "json Unmarshal of o2 should not return error %v", err)

	// return first value if second is nil
	v, err := z.Eval(o1, nil)
	assert.NoError(t, err, "notZeroValue(o1, nil) should not return error %v", err)
	b, ok := v.(map[string]interface{})
	assert.True(t, ok, "notZeroValue(o1, nil) should return map, not type %T", v)
	assert.Equal(t, "f", b["foo"], "notZeroValue(o1, nil) should contain value 'f', not %v", b)
	
	// return second value
	v, err = z.Eval(o1, o2)
	assert.NoError(t, err, "notZeroValue(o1, o2) should not return error %v", err)
	b, ok = v.(map[string]interface{})
	assert.True(t, ok, "notZeroValue(o1, o2) should return map, not type %T", v)
	assert.Equal(t, "b", b["bar"], "notZeroValue(o1, o2) should return contain value 'b', not %v", b)
}