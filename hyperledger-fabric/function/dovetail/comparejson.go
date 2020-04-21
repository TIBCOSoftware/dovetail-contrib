package dovetail

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/support/log"
	"github.com/nsf/jsondiff"
	"encoding/json"
	"fmt"
)

// Note: alternative implementation could use reflect.DeepEqual(x, y) to check for FullMatch

// CompareJSON dummy struct
type CompareJSON struct {
}

func init() {
	function.Register(&CompareJSON{})
}

// Name of function
func (s *CompareJSON) Name() string {
	return "compareJSON"
}

// Sig - function signature
func (s *CompareJSON) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeAny, data.TypeAny}, false
}

// Eval - function implementation
func (s *CompareJSON) Eval(params ...interface{}) (interface{}, error) {

	log.RootLogger().Debugf("Start compareJSON function with params %+v", params)

	if len(params) < 2 {
		return nil, fmt.Errorf("need 2 input objects to compare")
	}
	
	p1, err := toJSONBytes(params[0])
	if err != nil {
		return nil, fmt.Errorf("first param is invalid JSON: %v", err)
	}
	p2, err := toJSONBytes(params[1])
	if err != nil {
		return nil, fmt.Errorf("second param is invalid JSON: %v", err)
	}

	opts := jsondiff.DefaultConsoleOptions()
	opts.PrintTypes = false
	diff, _ := jsondiff.Compare(p1, p2, &opts)

	switch diff {
	case jsondiff.FullMatch:
		return "FullMatch", nil
	case jsondiff.SupersetMatch:
		return "SupersetMatch", nil
	default:
		return "NoMatch", nil
	}
}

func toJSONBytes(param interface{}) ([]byte, error) {
	switch v := param.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return json.Marshal(param)
	}
}
