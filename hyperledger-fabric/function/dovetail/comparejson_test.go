package dovetail

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

var d = &CompareJSON{}

func TestCompareBytes(t *testing.T) {
	o2 := []byte(`{"foo": 2, "bar": "two", "car": "opt"}`)
	o2p := []byte(`{"bar": "two", "foo": 2}`)
	v, err := d.Eval(o2, o2p)
	assert.NoError(t, err, "compareJSON(o2, o2p) should not return error %v", err)
	b, ok := v.(string)
	assert.True(t, ok, "compareJSON(o2,o2p) should return string, not type %T", v)
	assert.Equal(t, "SupersetMatch", b, "compareJSON(o1,o1s) should return 'FullMatch', not %s", b)
}

func TestCompareString(t *testing.T) {
	o2 := `{"foo": 2, "bar": "two", "car": "opt"}`
	o2p := `{"bar": "two", "foo": 2}`
	v, err := d.Eval(o2, o2p)
	assert.NoError(t, err, "compareJSON(o2, o2p) should not return error %v", err)
	b, ok := v.(string)
	assert.True(t, ok, "compareJSON(o2,o2p) should return string, not type %T", v)
	assert.Equal(t, "SupersetMatch", b, "compareJSON(o1,o1s) should return 'FullMatch', not %s", b)
}

func TestCompareJSONObjects(t *testing.T) {
	o1 := []byte(`{"foo": 1, "bar": "one"}`)
	o1s := []byte(`{"bar": "one", "foo": 1}`)
	o2 := []byte(`{"foo": 2, "bar": "two", "car": "opt"}`)
	o2p := []byte(`{"bar": "two", "foo": 2}`)

	var p1, p1s, p2, p2p interface{}

	// compare o1, o2
	err := json.Unmarshal(o1, &p1)
	assert.NoError(t, err, "parse JSON o1 should not return error: %v", err)
	err = json.Unmarshal(o2, &p2)
	assert.NoError(t, err, "parse JSON o2 should not return error: %v", err)
	v, err := d.Eval(p1, p2)
	assert.NoError(t, err, "compareJSON(o1, o2) should not return error %v", err)
	b, ok := v.(string)
	assert.True(t, ok, "compareJSON(o1,o2) should return string, not type %T", v)
	assert.Equal(t, "NoMatch", b, "compareJSON(o1,o2) should return 'NoMatch', not %s", b)

	// compare o1, o1s
	err = json.Unmarshal(o1s, &p1s)
	assert.NoError(t, err, "parse JSON o1s should not return error: %v", err)
	v, err = d.Eval(p1, p1s)
	assert.NoError(t, err, "compareJSON(o1, o1s) should not return error %v", err)
	b, ok = v.(string)
	assert.True(t, ok, "compareJSON(o1,o1s) should return string, not type %T", v)
	assert.Equal(t, "FullMatch", b, "compareJSON(o1,o1s) should return 'FullMatch', not %s", b)

	// compare o2, o2p
	err = json.Unmarshal(o2p, &p2p)
	assert.NoError(t, err, "parse JSON o2p should not return error: %v", err)
	v, err = d.Eval(p2, p2p)
	assert.NoError(t, err, "compareJSON(o2, o2p) should not return error %v", err)
	b, ok = v.(string)
	assert.True(t, ok, "compareJSON(o2,o2p) should return string, not type %T", v)
	assert.Equal(t, "SupersetMatch", b, "compareJSON(o1,o1s) should return 'FullMatch', not %s", b)
}
