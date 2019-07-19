package string

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/support/log"
)


type Substring struct {
}

func init() {
	function.Register(&Substring{})
}

func (s *Substring) Name() string {
	return "substring"
}

func (s *Substring) GetCategory() string {
	return "string"
}

func (s *Substring) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeInt, data.TypeInt}, false
}

func (s *Substring) Eval(params ...interface{}) (interface{}, error) {

	str, err := coerce.ToString(params[0])
	if err != nil {
		return nil, fmt.Errorf("string.substring function first parameter [%+v] must be string", params[0])
	}

	index, err := coerce.ToInt(params[1])
	if err != nil {
		return nil, fmt.Errorf("string.substring function second parameter [%+v] must be integer", params[1])
	}

	length, err := coerce.ToInt(params[2])
	if err != nil {
		return nil, fmt.Errorf("string.substring function third parameter [%+v] must be integer", params[1])
	}

	log.RootLogger().Debugf("Start substring function with parameters string %s index %d and lenght %d", str, index, length)
	return str[index : index+length], nil
}
