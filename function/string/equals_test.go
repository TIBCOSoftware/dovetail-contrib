package string

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var eq = &Equals{}

func TestStaticFunc_Eq(t *testing.T) {
	final1, _ := eq.Eval("TIBCO Web Integrator", "TIBCO")
	fmt.Println(final1)
	assert.Equal(t, false, final1)

	final2, _ := eq.Eval("123T", "123t")
	fmt.Println(final2)
	assert.Equal(t, false, final2)

}

func TestEQExpression(t *testing.T) {
	fun, err := factory.NewExpr(`string.equals("TIBCO NAME", "TIBCO NAME")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}

func TestQExpressionIgnoreCase(t *testing.T) {
	fun, err := factory.NewExpr(`string.equals("TIBCO name", "TIBCO NAME")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}
