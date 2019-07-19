package string

import (
	"git.tibco.com/git/product/ipaas/wi-contrib.git/function/datetime"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

type DatetimeFormat struct {
}

func init() {
	function.Register(&DatetimeFormat{})
}

func (s *DatetimeFormat) Name() string {
	return "datetimeFormat"
}

func (s *DatetimeFormat) GetCategory() string {
	return "string"
}

func (s *DatetimeFormat) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{}, true
}

func (s *DatetimeFormat) Eval(d ...interface{}) (interface{}, error) {
	return datetime.GetDatetimeFormat(), nil
}
