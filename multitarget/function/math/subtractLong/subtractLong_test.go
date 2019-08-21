package mathfunc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var s = &Subtract{}

func Test_sumLong(t *testing.T) {
	final, err := s.Eval(600, 200, 100)
	assert.Nil(t, err)
	fmt.Println(final)
	assert.Equal(t, final, int64(300))
}
