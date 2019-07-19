package string

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var tos = &String{}

func TestSample(t *testing.T) {
	final, _ := tos.Eval(123)
	assert.Equal(t, final, "123")
}

func TestFloat(t *testing.T) {
	final, _ := tos.Eval(float64(1234))
	assert.Equal(t, final, "1234")
}

func TestTostringExpression(t *testing.T) {
	fun, err := factory.NewExpr(`string.tostring(123)`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, "123", v)
	fmt.Println(v)
}
