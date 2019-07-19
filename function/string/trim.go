package string

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/log"
	"strings"

	"github.com/project-flogo/core/data/expression/function"
)


type Trim struct {
}

func init() {
	function.Register(&Trim{})
}

func (s *Trim) Name() string {
	return "trim"
}

func (s *Trim) GetCategory() string {
	return "string"
}

func (s *Trim) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

func (s *Trim) Eval(params ...interface{}) (interface{}, error) {
	str, err := coerce.ToString(params[0])
	if err != nil {
		return nil, fmt.Errorf("string.trim function first parameter [%+v] must be string", params[0])
	}
	log.RootLogger().Debugf("Trim a string \"%s\" with all leading and trailing white space removed", str)
	return strings.TrimSpace(str), nil
}
