package dovetail

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
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
	v, err := p.Eval("9/2/2020")
	assert.NoError(t, err, "parse date '9/2/2020' should not raise error %v", err)
	d, ok := v.([]int)
	assert.True(t, ok, "parse-date return type %T should be int array", v)
	assert.Equal(t, 2020, d[0], "year in '9/2/2020' should not be %v", d[0])
	assert.Equal(t, 2, d[1], "month in '9/2/2020' should not be %v", d[1])
	assert.Equal(t, 9, d[2], "day in '9/2/2020' should not be %v", d[2])
}

func TestParseInvalidDate(t *testing.T) {
	v, err := p.Eval("2020-14-01")
	assert.NoError(t, err, "parse date '2020-14-01' should not raise error %v", err)
	d, ok := v.([]int)

	// for invalid month, it took the first digit as month, the second digit as day
	assert.True(t, ok, "parse-date return type %T should be int array", v)
	assert.Equal(t, 2020, d[0], "year in '2020-14-01' should not be %v", d[0])
	assert.Equal(t, 1, d[1], "month in '2020-14-01' should not be %v", d[1])
	assert.Equal(t, 4, d[2], "day in '2020-14-01' should not be %v", d[2])
}

func TestParseNotDate(t *testing.T) {
	v, err := p.Eval("invalid-date")
	assert.NoError(t, err, "parse date 'invalid-date' should not raise error %v", err)
	d, ok := v.([]int)
	assert.True(t, ok, "parse-date return type %T should be int array", v)

	// it returns the current date
	n := time.Now()
	year, month, day := n.Date()
	assert.Equal(t, year, d[0], "'invalid-date' returned year %v, not matching current year %d", d[0], year)
	assert.Equal(t, int(month), d[1], "'invalid-date' returned month %v, not matching current month %d", d[1], int(month))
	assert.Equal(t, day, d[2], "'invalid-date' returned day %v, not matching current day %d", d[2], day)
}