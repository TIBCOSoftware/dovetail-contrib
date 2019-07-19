package string

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimeExpression(t *testing.T) {
	fun, err := factory.NewExpr(`string.timeFormat()`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Println(v)
}
