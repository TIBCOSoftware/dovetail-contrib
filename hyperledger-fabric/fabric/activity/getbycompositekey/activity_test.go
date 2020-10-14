package getbycompositekey

import (
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
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
	tc.SetInputObject(&Input{KeyName: "testKey=attr1,attr2", Attributes: map[string]interface{}{"attr1": "value1", "attr2": "value2"}})

	input := &Input{}
	err = tc.GetInputObject(input)
	assert.NoError(t, err, "GetInputObject should not return error")

	keyName, keyValues, err := extractCompositeKeyValues(input.KeyName, input.Attributes)
	assert.NoError(t, err, "extractCompositeKeyValues should not return error")

	assert.Equal(t, "testKey", keyName, "extracted composite key %s should be 'tetKey'", keyName)
	assert.Equal(t, 2, len(keyValues), "extracted composite key length %d should be 2", len(keyValues))
	assert.Equal(t, "value2", keyValues[1], "second attribute of composite key %s is not 'value2'", keyValues[1])
}
