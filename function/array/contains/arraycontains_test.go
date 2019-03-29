package contains

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var s = &Contains{}

func TestStaticFunc_Contains(t *testing.T) {
	
	array := []string{"Cat", "Dog", "Snake"}
	final := s.Eval(array, "Snake")
	assert.Equal(t, true, final)
	final = s.Eval(array, "Foo")
	assert.Equal(t, false, final)
	
	arrayInt := []int{5, 40, 10}
	final = s.Eval(arrayInt, 40)
	assert.Equal(t, true, final)
	final = s.Eval(arrayInt, "Foo")
	assert.Equal(t, false, final)
	
}
