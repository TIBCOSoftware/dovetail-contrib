package concat

import (
	"fmt"
	"testing"

	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression"
	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/tostring"
	"github.com/stretchr/testify/assert"
)

var s = &Concat{}

func TestStaticFunc_Concat(t *testing.T) {
	final, err := s.Eval("TIBCO", "Web", "Integrator")
	assert.Nil(t, err)
	fmt.Println(final)
	assert.Equal(t, final, "TIBCOWebIntegrator")
}

func TestSample(t *testing.T) {
	result, err := s.Eval("Hello", "World")
	assert.Nil(t, err)
	assert.Equal(t, "HelloWorld", result)
}

func TestOneArgument(t *testing.T) {
	_, err := s.Eval("Hello")
	assert.NotNil(t, err)
}

func TestExpressionDoubleQuotes(t *testing.T) {
	fun, err := expression.ParseExpression(`string.concat('Web',' Inte"grator')`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, `Web Inte"grator`, v)
}

func TestExpressionSingleQuote(t *testing.T) {
	fun, err := expression.ParseExpression(`string.concat("Web"," Inte'gr\a{tor")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, "Web Inte'gr\\a{tor", v)
}

func TestExpressionCombine(t *testing.T) {
	fun, err := expression.ParseExpression(`string.concat('Hello', " 'World'")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, `Hello 'World'`, v)
}

func TestExpressionCombine2(t *testing.T) {
	fun, err := expression.ParseExpression(`string.concat('Hello', ' "World"')`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, `Hello "World"`, v)
}
func TestExpression3(t *testing.T) {
	fun, err := expression.ParseExpression(`string.concat(
	"Web",
	" Integrator"
	)`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, "Web Integrator", v)
}

func TestExpressionSpace(t *testing.T) {
	fun, err := expression.ParseExpression(`string.concat(    "Web"  ,  " Integrator")   `)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, "Web Integrator", v)
}

func TestExpressionSpaceNewLineTab(t *testing.T) {
	fun, err := expression.ParseExpression(`string.concat(    "Web" 
		 ,	" Integrator"	
		 )`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, "Web Integrator", v)
}

func TestExpressionWithMappingRef(t *testing.T) {
	fun, err := expression.ParseExpression(`string.concat("abc", $flow.pathParams.id)`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, "abc$flow.pathParams.id", v)
}

func TestExpressionDoubleDoubleQuotes(t *testing.T) {
	fun, err := expression.ParseExpression(`string.concat("\"abc\"", "dddd")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, `"abc"dddd`, v)
}

func TestExpressionSingleSingleQuote(t *testing.T) {
	fun, err := expression.ParseExpression(`string.concat('\'b\'ac\'', "dddd")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, `'b'ac'dddd`, v)
}

func TestExpressionDolarInref(t *testing.T) {
	fun, err := expression.ParseExpression(`string.concat('\'b\'ac\'', $ActivityName.id.$Class)`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, `'b'ac'$ActivityName.id.$Class`, v)
}

func TestExpressionSpecialField(t *testing.T) {
	fun, err := expression.ParseExpression(`string.concat("name",$[0]["id^1"])`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, `name$[0]["id^1"]`, v)

	fun, err = expression.ParseExpression(`string.concat("name",$RESTInvoke.id["name_&Name"]["id^1"])`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err = fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, `name$RESTInvoke.id["name_&Name"]["id^1"]`, v)

	fun, err = expression.ParseExpression(`string.concat("name",$RESTInvoke.id[0]["god&d"])`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err = fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, `name$RESTInvoke.id[0]["god&d"]`, v)

	fun, err = expression.ParseExpression(`string.concat($RESTInvoke.id[0]["god&d"],$dd[0]["id^1"])`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err = fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, `$RESTInvoke.id[0]["god&d"]$dd[0]["id^1"]`, v)

	fun, err = expression.ParseExpression(`string.concat($RESTInvoke.id[0], "ok")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err = fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, `$RESTInvoke.id[0]ok`, v)

}

func TestExpressionIPAS6460(t *testing.T) {
	fun, err := expression.ParseExpression(`string.concat(string.tostring($flow.body.array1[0].array2[0]["id.2"].name),string.tostring($flow.headers["Content-Type"]))`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, `$flow.body.array1[0].array2[0]["id.2"].name$flow.headers["Content-Type"]`, v)
}

func TestExpressionIPAS6472(t *testing.T) {
	fun, err := expression.ParseExpression(`string.concat($flow.body["[squreal"], "sss")`)
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.Equal(t, `$flow.body["[squreal"]sss`, v)
}

func TestUnescape(t *testing.T) {
	str := strings.Replace(`\"abc\"`, "\\\"", "\"", -1)
	fmt.Println(str)
}
