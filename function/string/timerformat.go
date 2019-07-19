package string

import (
	"git.tibco.com/git/product/ipaas/wi-contrib.git/function/datetime"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

type TimeFormat struct {
}

func init() {
	function.Register(&TimeFormat{})
}

func (s *TimeFormat) Name() string {
	return "timeFormat"
}

func (s *TimeFormat) GetCategory() string {
	return "string"
}

func (s *TimeFormat) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{}, false
}

func (s *TimeFormat) Eval(d ...interface{}) (interface{}, error) {
	return datetime.GetTimeFormat(), nil
}
