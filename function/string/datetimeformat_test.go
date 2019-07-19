package string

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatetimeFormatExpression(t *testing.T) {
	fun, err := factory.NewExpr(`string.datetimeFormat()`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}
