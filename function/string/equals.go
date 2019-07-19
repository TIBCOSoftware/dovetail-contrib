package string

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression/function"
)

type Equals struct {
}

func init() {
	function.Register(&Equals{})
}

func (s *Equals) Name() string {
	return "equals"
}

func (s *Equals) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString}, false
}

func (s *Equals) Eval(params ...interface{}) (interface{}, error) {

	str, err := coerce.ToString(params[0])
	if err != nil {
		return nil, fmt.Errorf("string.equals function first parameter [%+v] must be string", params[0])
	}
	str2, err := coerce.ToString(params[1])
	if err != nil {
		return nil, fmt.Errorf("string.quals function second parameter [%+v] must be string", params[1])
	}
	return str == str2, nil
}
