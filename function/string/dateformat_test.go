package string

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDateFormatExpression(t *testing.T) {
	fun, err := factory.NewExpr(`string.dateFormat()`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}
