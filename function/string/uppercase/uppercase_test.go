package uppercase

import (
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression"
	"github.com/stretchr/testify/assert"
)

var s = &UpperCase{}

func TestStaticFunc_Upper_case(t *testing.T) {
	final := s.Eval("TIBCO Web Integrator")
	fmt.Println(final)
	assert.Equal(t, "TIBCO WEB INTEGRATOR", final)

	final = s.Eval("212")
	fmt.Println(final)
	assert.Equal(t, "212", final)

}

func TestExpression(t *testing.T) {
	fun, err := expression.ParseExpression(`string.upperCase("tibco")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}

func TestExpression2(t *testing.T) {
	fun, err := expression.ParseExpression(`string.upperCase(123456789)`)
	assert.Nil(t, err)
	_, err = fun.Eval()
	assert.Nil(t, err)
}
