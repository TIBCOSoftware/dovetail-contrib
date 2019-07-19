package string

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/support/log"
	"strings"
)

type UpperCase struct {
}

func init() {
	function.Register(&UpperCase{})
}

func (s *UpperCase) Name() string {
	return "upperCase"
}

func (s *UpperCase) GetCategory() string {
	return "string"
}

func (s *UpperCase) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

func (s *UpperCase) Eval(params ...interface{}) (interface{}, error) {

	str, err := coerce.ToString(params[0])
	if err != nil {
		return nil, fmt.Errorf("string.upperCase function first parameter [%+v] must be string", params[0])
	}

	log.RootLogger().Debugf("Returns the upper case of the given string \"%s\"", str)

	return strings.ToUpper(str), nil
}
