package string

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/log"
	"strings"

	"github.com/project-flogo/core/data/expression/function"
)

type Contains struct {
}

func init() {
	function.Register(&Contains{})
}

func (s *Contains) Name() string {
	return "contains"
}

func (s *Contains) GetCategory() string {
	return "string"
}

func (s *Contains) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString}, true
}

func (s *Contains) Eval(params ...interface{}) (interface{}, error) {

	str, err := coerce.ToString(params[0])
	if err != nil {
		return nil, fmt.Errorf("string.contains function first parameter [%+v] must be string", params[0])
	}
	substr, err := coerce.ToString(params[1])
	if err != nil {
		return nil, fmt.Errorf("string.contains function second parameter [%+v] must be string", params[1])
	}

	log.RootLogger().Debugf("Reports whether \"%s\" is within \"%s\"", substr, str)

	return strings.Contains(str, substr), nil
}
