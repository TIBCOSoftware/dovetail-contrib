package string

import (
	"git.tibco.com/git/product/ipaas/wi-contrib.git/function/datetime"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)


type DateFormat struct {
}

func init() {
	function.Register(&DateFormat{})
}

func (s *DateFormat) Name() string {
	return "dateFormat"
}

func (s *DateFormat) GetCategory() string {
	return "string"
}

func (s *DateFormat) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{}, true
}

func (s *DateFormat) Eval(d ...interface{}) (interface{}, error) {
	return datetime.GetDateFormat(), nil
}
