package contains

import (
	"reflect"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("array-contains-function")

type Contains struct {
}

func init() {
	function.Registry(&Contains{})
}

func (s *Contains) GetName() string {
	return "contains"
}

func (s *Contains) GetCategory() string {
	return "array"
}

func (s *Contains) Eval(array, item interface{}) bool {
	log.Infof("Looking for \"%s\" in \"%s\"", item, array)
	if array == nil || item == nil {
		return false
	}
	 arrV := reflect.ValueOf(array)
    if arrV.Kind() == reflect.Slice {
        for i := 0; i < arrV.Len(); i++ {
            if arrV.Index(i).Interface() == item {
                return true
            }
        }
    }
    return false
}
