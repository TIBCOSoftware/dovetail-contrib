package dovetail

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/support/log"
	"reflect"
	"fmt"
)

// ConditionalChoice dummy struct
type ConditionalChoice struct {
}

func init() {
	function.Register(&ConditionalChoice{})
}

// Name of function
func (s *ConditionalChoice) Name() string {
	return "conditionalChoice"
}

// Sig - function signature
func (s *ConditionalChoice) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeBool, data.TypeAny, data.TypeAny, data.TypeAny}, false
}

// Eval - function implementation
func (s *ConditionalChoice) Eval(params ...interface{}) (interface{}, error) {

	log.RootLogger().Debugf("Start conditionalChoice function with params %+v", params)

	if len(params) < 4 {
		return nil, fmt.Errorf("invalid number of parameters for conditionalChoice")
	}

	condition, ok := params[0].(bool)
	if !ok {
		return params[1], fmt.Errorf("condition is not boolean, so return original value %v", params[1])
	}

	if !condition && !reflect.DeepEqual(params[1], params[2]) {
		return params[3], nil
	}
	return params[2], nil
}
