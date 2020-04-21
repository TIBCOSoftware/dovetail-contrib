package dovetail

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/support/log"
	"reflect"
)

// NotZeroValue dummy struct
type NotZeroValue struct {
}

func init() {
	function.Register(&NotZeroValue{})
}

// Name of function
func (s *NotZeroValue) Name() string {
	return "notZeroValue"
}

// Sig - function signature
func (s *NotZeroValue) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeAny, data.TypeAny}, false
}

// Eval - function implementation
func (s *NotZeroValue) Eval(params ...interface{}) (interface{}, error) {

	log.RootLogger().Debugf("Start notZeroValue function with params %+v", params)

	if len(params) < 2 {
		return params[0], nil
	}

	x := params[1]
	if x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface()) {
		return params[0], nil
	}
	return x, nil
}
