package dovetail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var p = &ParseDate{}

func TestParseDate(t *testing.T) {
	v, err := p.Eval("2020-10-01")
	assert.NoError(t, err, "parse date '2020-10-01' should not raise error %v", err)
	d, ok := v.([]int)
	assert.True(t, ok, "parse-date return type %T should be int array", v)
	assert.Equal(t, 2020, d[0], "year in '2020-10-01' should not be %v", d[0])
	assert.Equal(t, 10, d[1], "month in '2020-10-01' should not be %v", d[1])
	assert.Equal(t, 1, d[2], "day in '2020-10-01' should not be %v", d[2])
}

func TestParseDateTime(t *testing.T) {
	v, err := p.Eval("2020-10-01T11:10:30")
	assert.NoError(t, err, "parse date '2020-10-01' should not raise error %v", err)
	d, ok := v.([]int)
	assert.True(t, ok, "parse-date return type %T should be int array", v)
	assert.Equal(t, 2020, d[0], "year in '2020-10-01' should not be %v", d[0])
	assert.Equal(t, 10, d[1], "month in '2020-10-01' should not be %v", d[1])
	assert.Equal(t, 1, d[2], "day in '2020-10-01' should not be %v", d[2])
}

func TestParseEuroDate(t *testing.T) {
	v, err := p.Eval("09/02/2020", "01/02/2006")
	assert.NoError(t, err, "parse date '9/2/2020' should not raise error %v", err)
	d, ok := v.([]int)
	assert.True(t, ok, "parse-date return type %T should be int array", v)
	assert.Equal(t, 2020, d[0], "year in '9/2/2020' should not be %v", d[0])
	assert.Equal(t, 9, d[1], "month in '9/2/2020' should not be %v", d[1])
	assert.Equal(t, 2, d[2], "day in '9/2/2020' should not be %v", d[2])
}

func TestParseInvalidDate(t *testing.T) {
	_, err := p.Eval("2020-14-01")
	assert.Errorf(t, err, "parse date '2020-14-01' should raise error 'month out of range'")
}

func TestParseNotDate(t *testing.T) {
	_, err := p.Eval("invalid-date")
	assert.Errorf(t, err, "parse date 'invalid-date' should raise error 'cannot parse'")
}
