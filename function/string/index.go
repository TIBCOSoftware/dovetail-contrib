package string

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/log"
	"strings"

	"github.com/project-flogo/core/data/expression/function"
)


type Index struct {
}

func init() {
	function.Register(&Index{})
}

func (s *Index) Name() string {
	return "index"
}

func (s *Index) GetCategory() string {
	return "string"
}

func (s *Index) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString}, false
}

func (s *Index) Eval(params ...interface{}) (interface{}, error) {

	str, err := coerce.ToString(params[0])
	if err != nil {
		return nil, fmt.Errorf("string.contains function first parameter [%+v] must be string", params[0])
	}
	sep, err := coerce.ToString(params[1])
	if err != nil {
		return nil, fmt.Errorf("string.contains function second parameter [%+v] must be string", params[1])
	}
	log.RootLogger().Debugf("Returns the index of the first instance of \"%s\" in \"%s\", or -1 if not present", sep, str)
	var l int
	l = strings.Index(str, sep)
	log.RootLogger().Debugf("Done calculating the index: %n", l)
	return l, nil
}
