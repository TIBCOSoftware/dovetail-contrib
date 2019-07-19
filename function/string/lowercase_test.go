package string

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var low = &LowerCase{}

func TestStaticFunc_Lower_case(t *testing.T) {
	final, _ := low.Eval("TIBCO Web Integrator")
	fmt.Println(final)
	assert.Equal(t, "tibco web integrator", final)
}

func TestLowExpression(t *testing.T) {
	fun, err := factory.NewExpr(`string.lowerCase("TIBCO NAME")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}
