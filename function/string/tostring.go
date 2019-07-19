package string

import (
	"git.tibco.com/git/product/ipaas/wi-contrib.git/engine/conversion"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/support/log"
	"strconv"
)


type String struct {
}

func init() {
	function.Register(&String{})
}

func (s *String) Name() string {
	return "tostring"
}

func (s *String) GetCategory() string {
	return "string"
}

func (s *String) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeAny}, false
}

func (s *String) Eval(d ...interface{}) (interface{}, error) {
	log.RootLogger().Debugf("Start String function with parameters %s", d[0])

	switch t := d[0].(type) {
	case string:
		return t, nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	case int32, int8, int16, int64:
		//v := int64(in)
		return strconv.FormatInt(d[0].(int64), 10), nil
	case int:
		return strconv.Itoa(t), nil
	case *int:
		return strconv.Itoa(*t), nil
	case uint, uint8, uint16, uint32, uint64:
		//v := int64(in)
		return strconv.FormatInt(d[0].(int64), 10), nil
	default:
		str, err := conversion.ConvertToString(t)
		if err != nil {
			log.RootLogger().Errorf("Convert to string error %s", err.Error())
		}
		return str, nil
	}
}
