package uppercase

import (
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("upperCase-function")

type UpperCase struct {
}

func init() {
	function.Registry(&UpperCase{})
}

func (s *UpperCase) GetName() string {
	return "upperCase"
}

func (s *UpperCase) GetCategory() string {
	return "string"
}

func (s *UpperCase) Eval(str string) string {
	log.Debugf("Returns the upper case of the given string \"%s\"", str)

	return strings.ToUpper(str)
}
