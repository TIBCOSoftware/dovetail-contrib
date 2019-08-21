package float

import (
	"fmt"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/data/expression/script"
	"github.com/project-flogo/core/data/resolve"
	"testing"

	"github.com/stretchr/testify/assert"
)

var s = &Float{}

var resolver = resolve.NewCompositeResolver(map[string]resolve.Resolver{"static": nil, ".": nil, "env": &resolve.EnvResolver{}})
var factory = script.NewExprFactory(resolver)

func init() {
	function.ResolveAliases()
}

func TestSample(t *testing.T) {
	final, err := s.Eval("2.7787654231689989909", 17)
	assert.Nil(t, err)
	assert.Equal(t, float64(2.7787654231689989909), final)
}

func TestExpression(t *testing.T) {
	fun, err := factory.NewExpr(`float.float64("2.778765423168998990922",16)`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, float64(2.778765423168999), v)
	fmt.Println(v)
}

func TestExpression1(t *testing.T) {
	fun, err := factory.NewExpr(`float.float64("2.7787654231689989909")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, float64(2.778765423168999), v)
	fmt.Println(v)
}

func TestExpression2(t *testing.T) {
	fun, err := factory.NewExpr(`float.float64("2.7787654231689989909",10)`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, float64(2.7787654232), v)
	fmt.Println(v)
}
