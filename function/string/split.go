package string

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/log"
	"strings"

	"github.com/project-flogo/core/data/expression/function"
)

type Split struct {
}

func init() {
	function.Register(&Split{})
}

func (s *Split) Name() string {
	return "split"
}

func (s *Split) GetCategory() string {
	return "string"
}

func (s *Split) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString}, false
}

func (s *Split) Eval(params ...interface{}) (interface{}, error) {

	str, err := coerce.ToString(params[0])
	if err != nil {
		return nil, fmt.Errorf("string.regex function first parameter [%+v] must be string", params[0])
	}

	sep, err := coerce.ToString(params[1])
	if err != nil {
		return nil, fmt.Errorf("string.regex function first parameter [%+v] must be string", params[0])
	}
	log.RootLogger().Debugf("Slices \"%s\" into all substrings separated by \"%s\" and returns a slice of the substrings between those separators", str, sep)
	return strings.Split(str, sep), nil
}
