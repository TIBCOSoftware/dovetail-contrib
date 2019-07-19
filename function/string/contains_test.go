package string

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var con = &Contains{}

func TestStaticFunc_Contains(t *testing.T) {
	final1, _ := con.Eval("TIBCO Web Integrator", "Web")
	fmt.Println(final1)
	assert.Equal(t, true, final1)

	final2, _ := con.Eval("TIBCO 网路 Integrator", "网路")
	fmt.Println(final2)
	assert.Equal(t, true, final2)
}

func TestContainsExpression(t *testing.T) {
	fun, err := factory.NewExpr(`string.contains("TIBCO Web Integrator","Web")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, true, v)
	fmt.Println(v)
}
