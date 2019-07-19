package string

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStaticFunc_Substring(t *testing.T) {
	str := "TIBCO software Inc"
	s := &Substring{}
	subStr, _ := s.Eval(str, 0, 5)
	fmt.Println(subStr)
	assert.Equal(t, subStr, "TIBCO")
}

func TestSubStringSample(t *testing.T) {
	sub := &Substring{}
	result, _ := sub.Eval("12345", 2, 3)
	assert.Equal(t, "345", result)
}

func TestSubstringExpression(t *testing.T) {
	fun, err := factory.NewExpr(`string.substring("1999/04/01",2,3)`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}

func TestSubstringExpression2(t *testing.T) {
	fun, err := factory.NewExpr(`string.substring("TIBCO",2,3)`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}

func TestExpression3(t *testing.T) {
	fun, err := factory.NewExpr(`string.substring(datetime.currentDate(), 0, 10)`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}
