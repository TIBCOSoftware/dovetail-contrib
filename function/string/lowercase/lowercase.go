package lowercase

import (
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("lowerCase-function")

type LowerCase struct {
}

func init() {
	function.Registry(&LowerCase{})
}

func (s *LowerCase) GetName() string {
	return "lowerCase"
}

func (s *LowerCase) GetCategory() string {
	return "string"
}

func (s *LowerCase) Eval(str string) string {
	log.Debugf("Returns the lower case of string str", str)

	return strings.ToLower(str)
}
