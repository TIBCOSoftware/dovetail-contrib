package string

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/support/log"
)


type StringLength struct {
}

func init() {
	function.Register(&StringLength{})
}

func (s *StringLength) Name() string {
	return "length"
}

func (s *StringLength) GetCategory() string {
	return "string"
}

func (s *StringLength) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

func (s *StringLength) Eval(params ...interface{}) (interface{}, error) {

	str, err := coerce.ToString(params[0])
	if err != nil {
		return nil, fmt.Errorf("string.length function first parameter [%+v] must be string", params[0])
	}

	log.RootLogger().Debugf("Return the length of a string \"%s\"", str)
	var l int
	//l = len([]rune(str))
	l = len(str)
	log.RootLogger().Debugf("Done calculating the length %d", l)
	return l, nil
}
