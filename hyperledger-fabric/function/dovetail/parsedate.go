package dovetail

import (
	"fmt"
	"time"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/support/log"
)

// ParseDate dummy struct
type ParseDate struct {
}

func init() {
	function.Register(&ParseDate{})
}

// Name of function
func (s *ParseDate) Name() string {
	return "parseDate"
}

// Sig - function signature
func (s *ParseDate) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

// Eval - function implementation
func (s *ParseDate) Eval(params ...interface{}) (interface{}, error) {

	log.RootLogger().Debugf("Start parseDate function with param %+v", params[0])

	dateStr, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("param %T is not a date string", params[0])
	}

	fmtStr := "2006-01-02"
	if len(params) > 1 {
		if layout, ok := params[1].(string); ok && layout != "" {
			fmtStr = layout
		}
	}

	if len(dateStr) > len(fmtStr) {
		dateStr = dateStr[0:len(fmtStr)]
	}
	// use utility of github.com/tkuchiki/parsetime
	// p, _ := parsetime.NewParseTime()
	// t, err := p.Parse(dateStr)
	t, err := time.Parse(fmtStr, dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date string '%s': %v", dateStr, err)
	}
	year, month, day := t.Date()
	return []int{year, int(month), day}, nil
}
