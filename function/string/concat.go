package string

import (
	"bytes"
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/support/log"
)


type Concat struct {
}

func init() {
	function.Register(&Concat{})
}

func (s *Concat) Name() string {
	return "concat"
}

func (s *Concat) GetCategory() string {
	return "string"
}

func (s *Concat) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, true
}

func (s *Concat) Eval(strs ...interface{}) (interface{}, error) {
	log.RootLogger().Debugf("Start concat function with parameters %s", strs)
	if len(strs) >= 2 {
		var buffer bytes.Buffer

		for _, v := range strs {
			s, err := coerce.ToString(v)
			if err != nil {
				return nil, fmt.Errorf("concat function parameter [%+v] must be string.", v)
			}
			buffer.WriteString(s)
		}
		log.RootLogger().Debugf("Done concat function with result %s", buffer.String())
		return buffer.String(), nil
	}

	return "", fmt.Errorf("Concat function at least have 2 arguments")
}
