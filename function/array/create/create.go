package create

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"reflect"
)

var log = logger.GetLogger("array-create-function")

type Create struct {
}

func init() {
	function.Registry(&Create{})
}

func (s *Create) GetName() string {
	return "create"
}

func (s *Create) GetCategory() string {
	return "array"
}

func (s *Create) Eval(object ...interface{}) ([]interface{}, error) {
	log.Debugf("Start array function with parameters %v", object)
	if object == nil {
		return nil, nil
	}
	result := make([]interface{}, len(object))
    var shouldcheck bool
	for i, s := range object {
		if shouldcheck == true {
			type1 := reflect.TypeOf(result[0])
			type2 := reflect.TypeOf(s)
			if type1 != type2 {
				return nil, fmt.Errorf("Failed to create an array from [%v]. In-compatible types [%v,%v] found.", object, type1, type2)
			}
		}
		result[i] = s
		shouldcheck = true
	}
	log.Debugf("Done array function with result %v", result)
	return result, nil
}
