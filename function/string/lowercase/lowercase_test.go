package lowercase

import (
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression"
	"github.com/stretchr/testify/assert"
)

var s = &LowerCase{}

func TestStaticFunc_Lower_case(t *testing.T) {
	final := s.Eval("TIBCO Web Integrator")
	fmt.Println(final)
	assert.Equal(t, "tibco web integrator", final)
}

func TestExpression(t *testing.T) {
	fun, err := expression.ParseExpression(`string.lowerCase("TIBCO NAME")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}
