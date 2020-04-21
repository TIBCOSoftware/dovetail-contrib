package dovetail

import (
	"encoding/json"
	"fmt"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/support/log"
	"strconv"
	"reflect"
	"math"
)

// Sum dummy struct
type Sum struct {
}

func init() {
	function.Register(&Sum{})
}

// Name of function
func (s *Sum) Name() string {
	return "sum"
}

// Sig - function signature
func (s *Sum) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeAny}, false
}

// Eval - function implementation
func (s *Sum) Eval(params ...interface{}) (interface{}, error) {

	log.RootLogger().Debugf("Start sum function with param %+v", params[0])

	items := reflect.ValueOf(params[0])
	if items.Kind() != reflect.Slice {
		return nil, fmt.Errorf("param %T is not an array", params[0])
	}

	total := 0.0
	for i := 0; i < items.Len(); i++ {
		val := items.Index(i).Interface()
		if v, err := coerceToFloat(val); err == nil {
			if !math.IsNaN(v) && !math.IsInf(v, 0) {
				total += v
			}
		}
	}

	return total, nil
}

func coerceToFloat(val interface{}) (float64, error) {
	switch t := val.(type) {
	case int:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case float64:
		return float64(t), nil
	case json.Number:
		i, err := t.Float64()
		return float64(i), err
	case string:
		return strconv.ParseFloat(t, 64)
	case nil:
		return 0.0, nil
	default:
		return 0.0, fmt.Errorf("Unable to coerce %#v to float64", val)
	}
}
