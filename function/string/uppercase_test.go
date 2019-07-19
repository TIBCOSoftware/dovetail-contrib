package string

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var up = &UpperCase{}

func TestStaticFunc_Upper_case(t *testing.T) {
	final, _ := up.Eval("TIBCO Web Integrator")
	fmt.Println(final)
	assert.Equal(t, "TIBCO WEB INTEGRATOR", final)

	final, _ = up.Eval("212")
	fmt.Println(final)
	assert.Equal(t, "212", final)

}

func TestUpperCaseExpression(t *testing.T) {
	fun, err := factory.NewExpr(`string.upperCase("tibco")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}

func TestUpperCaseExpression2(t *testing.T) {
	fun, err := factory.NewExpr(`string.upperCase(123456789)`)
	assert.Nil(t, err)
	_, err = fun.Eval(nil)
	assert.Nil(t, err)
}
