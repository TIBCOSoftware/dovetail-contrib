package string

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var concat = &Concat{}

func TestStatic_Concat(t *testing.T) {
	final, err := concat.Eval("TIBCO", "Web", "Integrator")
	assert.Nil(t, err)
	fmt.Println(final)
	assert.Equal(t, final, "TIBCOWebIntegrator")
}

func TestConcatSample(t *testing.T) {
	result, err := concat.Eval("Hello", "World")
	assert.Nil(t, err)
	assert.Equal(t, "HelloWorld", result)
}

func TestOneArgument(t *testing.T) {
	_, err := concat.Eval("Hello")
	assert.NotNil(t, err)
}

func TestExpressionDoubleQuotes(t *testing.T) {
	fun, err := factory.NewExpr(`string.concat('Web',' Inte"grator')`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, `Web Inte"grator`, v)
}

//
//func TestExpressionSingleQuote(t *testing.T) {
//	fun, err := factory.NewExpr(`string.concat("Web"," Inte'gr\a{tor")`)
//	assert.Nil(t, err)
//	assert.NotNil(t, fun)
//	v, err := fun.Eval(nil)
//	assert.Nil(t, err)
//	assert.Equal(t, "Web Inte'gr\\a{tor", v)
//}

func TestExpressionCombine(t *testing.T) {
	fun, err := factory.NewExpr(`string.concat('Hello', " 'World'")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, `Hello 'World'`, v)
}

func TestConcatExpressionCombine2(t *testing.T) {
	fun, err := factory.NewExpr(`string.concat('Hello', ' "World"')`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, `Hello "World"`, v)
}
func TestConcatExpression3(t *testing.T) {
	fun, err := factory.NewExpr(`string.concat(
	"Web",
	" Integrator"
	)`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, "Web Integrator", v)
}

func TestExpressionSpace(t *testing.T) {
	fun, err := factory.NewExpr(`string.concat(    "Web"  ,  " Integrator")   `)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, "Web Integrator", v)
}

func TestExpressionSpaceNewLineTab(t *testing.T) {
	fun, err := factory.NewExpr(`string.concat(    "Web" 
		 ,	" Integrator"	
		 )`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, "Web Integrator", v)
}

//func TestExpressionWithMappingRef(t *testing.T) {
//	fun, err := factory.NewExpr(`string.concat("abc", $flow.pathParamconcat.id)`)
//	assert.Nil(t, err)
//	assert.NotNil(t, fun)
//	v, err := fun.Eval(nil)
//	assert.Nil(t, err)
//	assert.Equal(t, "abc$flow.pathParamconcat.id", v)
//}

//func TestExpressionDoubleDoubleQuotes(t *testing.T) {
//	fun, err := factory.NewExpr(`string.concat("\"abc\"", "dddd")`)
//	assert.Nil(t, err)
//	assert.NotNil(t, fun)
//	v, err := fun.Eval(nil)
//	assert.Nil(t, err)
//	assert.Equal(t, `"abc"dddd`, v)
//}
//
//func TestExpressionSingleSingleQuote(t *testing.T) {
//	fun, err := factory.NewExpr(`string.concat('\'b\'ac\'', "dddd")`)
//	assert.Nil(t, err)
//	assert.NotNil(t, fun)
//	v, err := fun.Eval(nil)
//	assert.Nil(t, err)
//	assert.Equal(t, `'b'ac'dddd`, v)
//}

//func TestExpressionDolarInref(t *testing.T) {
//	fun, err := factory.NewExpr(`string.concat('\'b\'ac\'', $ActivityName.id.$Class)`)
//	assert.Nil(t, err)
//	assert.NotNil(t, fun)
//	v, err := fun.Eval(nil)
//	assert.Nil(t, err)
//	assert.Equal(t, `'b'ac'$ActivityName.id.$Class`, v)
//}

//func TestExpressionSpecialField(t *testing.T) {
//	fun, err := factory.NewExpr(`string.concat("name",$[0]["id^1"])`)
//	assert.Nil(t, err)
//	assert.NotNil(t, fun)
//	v, err := fun.Eval(nil)
//	assert.Nil(t, err)
//	assert.Equal(t, `name$[0]["id^1"]`, v)
//
//	fun, err = factory.NewExpr(`string.concat("name",$RESTInvoke.id["name_&Name"]["id^1"])`)
//	assert.Nil(t, err)
//	assert.NotNil(t, fun)
//	v, err = fun.Eval(nil)
//	assert.Nil(t, err)
//	assert.Equal(t, `name$RESTInvoke.id["name_&Name"]["id^1"]`, v)
//
//	fun, err = factory.NewExpr(`string.concat("name",$RESTInvoke.id[0]["god&d"])`)
//	assert.Nil(t, err)
//	assert.NotNil(t, fun)
//	v, err = fun.Eval(nil)
//	assert.Nil(t, err)
//	assert.Equal(t, `name$RESTInvoke.id[0]["god&d"]`, v)
//
//	fun, err = factory.NewExpr(`string.concat($RESTInvoke.id[0]["god&d"],$dd[0]["id^1"])`)
//	assert.Nil(t, err)
//	assert.NotNil(t, fun)
//	v, err = fun.Eval(nil)
//	assert.Nil(t, err)
//	assert.Equal(t, `$RESTInvoke.id[0]["god&d"]$dd[0]["id^1"]`, v)
//
//	fun, err = factory.NewExpr(`string.concat($RESTInvoke.id[0], "ok")`)
//	assert.Nil(t, err)
//	assert.NotNil(t, fun)
//	v, err = fun.Eval(nil)
//	assert.Nil(t, err)
//	assert.Equal(t, `$RESTInvoke.id[0]ok`, v)
//
//}

//func TestExpressionIPAS6460(t *testing.T) {
//	fun, err := factory.NewExpr(`string.concat(string.tostring($flow.body.array1[0].array2[0]["id.2"].name),string.tostring($flow.headers["Content-Type"]))`)
//	assert.Nil(t, err)
//	assert.NotNil(t, fun)
//	v, err := fun.Eval(nil)
//	assert.Nil(t, err)
//	assert.Equal(t, `$flow.body.array1[0].array2[0]["id.2"].name$flow.headers["Content-Type"]`, v)
//}
//
//func TestExpressionIPAS6472(t *testing.T) {
//	fun, err := factory.NewExpr(`string.concat($flow.body["[squreal"], "sss")`)
//	assert.Nil(t, err)
//	assert.NotNil(t, fun)
//	v, err := fun.Eval(nil)
//	assert.Nil(t, err)
//	assert.Equal(t, `$flow.body["[squreal"]sss`, v)
//}

func TestUnescape(t *testing.T) {
	str := strings.Replace(`\"abc\"`, "\\\"", "\"", -1)
	fmt.Println(str)
}
