package tostring

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression"
	"github.com/stretchr/testify/assert"
	"testing"
)

var s = &String{}

func TestSample(t *testing.T) {
	final := s.Eval(123)
	assert.Equal(t, final, "123")
}

func TestFloat(t *testing.T) {
	final := s.Eval(float64(1234))
	assert.Equal(t, final, "1234")
}



func TestExpression(t *testing.T) {
	fun, err := expression.ParseExpression(`string.tostring(123)`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, "123", v)
	fmt.Println(v)
}
