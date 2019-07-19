package string

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/log"
	"strings"

	"github.com/project-flogo/core/data/expression/function"
)


type EqualsIgnoreCase struct {
}

func init() {
	function.Register(&EqualsIgnoreCase{})
}

func (s *EqualsIgnoreCase) Name() string {
	return "equalsIgnoreCase"
}

func (s *EqualsIgnoreCase) GetCategory() string {
	return "string"
}

func (s *EqualsIgnoreCase) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString}, false
}

func (s *EqualsIgnoreCase) Eval(params ...interface{}) (interface{}, error) {

	str, err := coerce.ToString(params[0])
	if err != nil {
		return nil, fmt.Errorf("string.contains function first parameter [%+v] must be string", params[0])
	}
	str2, err := coerce.ToString(params[1])
	if err != nil {
		return nil, fmt.Errorf("string.contains function second parameter [%+v] must be string", params[1])
	}
	log.RootLogger().Debugf(`Reports whether "%s" equels "%s" with ignore case`, str, str2)
	return strings.EqualFold(str, str2), nil
}
