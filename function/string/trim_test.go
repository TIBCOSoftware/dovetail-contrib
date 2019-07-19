package string

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStaticFunc_Trim(t *testing.T) {
	str := " \t\n TIBCO software Inc \n\t\r\n"
	s := &Trim{}
	subStr, _ := s.Eval(str)
	fmt.Println(subStr)
	assert.Equal(t, "TIBCO software Inc", subStr)
}

func TestTrimExpression(t *testing.T) {
	fun, err := factory.NewExpr(`string.trim("    TIBCO")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}
