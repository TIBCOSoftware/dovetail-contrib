package string

import (
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/log"
	"regexp"

	"github.com/project-flogo/core/data/expression/function"
)

type Regex struct {
}

func init() {
	function.Register(&Regex{})
}

func (s *Regex) Name() string {
	return "regex"
}

func (s *Regex) GetCategory() string {
	return "string"
}

func (s *Regex) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString}, false
}

func (s *Regex) Eval(params ...interface{}) (interface{}, error) {

	pattern, err := coerce.ToString(params[0])
	if err != nil {
		return nil, fmt.Errorf("string.regex function first parameter [%+v] must be string", params[0])
	}

	str, err := coerce.ToString(params[1])
	if err != nil {
		return nil, fmt.Errorf("string.regex function first parameter [%+v] must be string", params[0])
	}

	log.RootLogger().Debugf("Checks whether a textual regular expression \"%s\" matches a string \"%s\"", pattern, str)

	matched, err := regexp.MatchString(pattern, str)

	if err != nil {
		log.RootLogger().Debugf("Error occurred during regular expression matching: %s", err)
	}

	return matched, nil
}
