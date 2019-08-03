package tostring

import (
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/support/log"
)

var logger = log.RootLogger()

type String struct {
}

func init() {
	function.Registry(&String{})
}

func (s *String) GetName() string {
	return "tostring"
}

func (s *String) GetCategory() string {
	return "string"
}

func (s *String) Eval(in interface{}) string {
	logger.Debugf("Start String function with parameters %s", in)

	v, e := coerce.ToString(in)
	if e != nil {
		logger.Errorf("string:tostring error: %v", e)
		return ""
	} else {
		return v
	}
}
