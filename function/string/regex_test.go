package string

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var reg = &Regex{}

func TestStaticFunc_Concat(t *testing.T) {
	final, _ := reg.Eval("foo.*", "Tseafood")
	fmt.Println(final)
	assert.Equal(t, true, final)

	final2, _ := reg.Eval("bar.*", "seafood")
	fmt.Println(final2)
	assert.Equal(t, false, final2)

	final3, _ := reg.Eval("a(b", "seafood")
	fmt.Println(final3)
	assert.Equal(t, false, final3)
}

func TestRegexExpression(t *testing.T) {
	fun, err := factory.NewExpr(`string.regex("foo.*","seafood")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}
