package string

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var slen = &StringLength{}

func TestStaticFunc_String_length(t *testing.T) {
	final11, _ := slen.Eval("TIBCO Web Integrator")
	fmt.Println(final11)
	assert.Equal(t, int(20), final11)

	final2, _ := slen.Eval("TIBCO 网路集成器")
	fmt.Println(final2)
	assert.Equal(t, int(21), final2)
}

func TestLenExpression(t *testing.T) {
	fun, err := factory.NewExpr(`string.length("seafood,name")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}
