package dovetail

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"math"
)

var s = &Sum{}

func TestIntArray(t *testing.T) {
	v, err := s.Eval([]int{1, 2, 3})
	assert.NoError(t, err, "sum of [1,2,3] should not raise error %v", err)
	assert.Equal(t, float64(6), v, "sum of [1,2,3] should not be %v", v)
}

func TestStringArray(t *testing.T) {
	v, err := s.Eval([]string{"1.5", "2.5", "foo"})
	assert.NoError(t, err, "sum of [\"1.5\", \"2.5\", \"foo\"] should not raise error %v", err)
	assert.Equal(t, float64(4), v, "sum of [\"1.5\", \"2.5\", \"foo\"] should not be %v", v)
}

func TestInf(t *testing.T) {
	v, err := s.Eval([]interface{}{"1.5", 2.5, math.Inf(1)})
	assert.NoError(t, err, "sum of [\"1.5\", \"2.5\", \"Inf\"] should not raise error %v", err)
	assert.Equal(t, float64(4), v, "sum of [\"1.5\", \"2.5\", \"Inf\"] should not be %v", v)
}

func TestNaN(t *testing.T) {
	v, err := s.Eval([]interface{}{"1.5", 2.5, math.NaN()})
	assert.NoError(t, err, "sum of [\"1.5\", \"2.5\", \"NaN\"] should not raise error %v", err)
	assert.Equal(t, float64(4), v, "sum of [\"1.5\", \"2.5\", \"NaN\"] should not be %v", v)
}
