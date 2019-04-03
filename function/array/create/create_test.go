package create

import (
	"fmt"
	"testing"

	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/boolean/false"
	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/tostring"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression"
	"github.com/stretchr/testify/assert"
)

var s = &Create{}

func TestStaticFunc_ArrayString(t *testing.T) {
	expectedResult := []string{"Cat", "Dog", "Snake"}
	final, err := s.Eval("Cat", "Dog", "Snake")
	assert.Nil(t, err)
	fmt.Println(final)
	for i, item := range final {
		assert.Equal(t, item.(string), expectedResult[i])
	}
}

func TestExpression(t *testing.T) {
	fun, err := expression.ParseExpression(`array.create("123","456")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}

func TestExpression2(t *testing.T) {
	fun, err := expression.ParseExpression(`string.tostring(array.create("adi","shukla",boolean.false()))`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	_, err = fun.Eval()
	assert.NotNil(t, err)
}

func TestExpression3(t *testing.T) {
	fun, err := expression.ParseExpression(`string.tostring(array.create("adi","shukla",true))`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	_, err = fun.Eval()
	assert.NotNil(t, err)
}
