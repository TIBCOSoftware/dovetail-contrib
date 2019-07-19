package string

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/log"
	"strings"

	"github.com/project-flogo/core/data/expression/function"
)


type LowerCase struct {
}

func init() {
	function.Register(&LowerCase{})
}

func (s *LowerCase) Name() string {
	return "lowerCase"
}

func (s *LowerCase) GetCategory() string {
	return "string"
}

func (s *LowerCase) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

func (s *LowerCase) Eval(params ...interface{}) (interface{}, error) {

	str, err := coerce.ToString(params[0])
	if err != nil {
		return nil, fmt.Errorf("string.contains function first parameter [%+v] must be string", params[0])
	}

	log.RootLogger().Debugf("Returns the lower case of string str", str)

	return strings.ToLower(str), nil
}
