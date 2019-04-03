package tostring

import (
	"strconv"

	"git.tibco.com/git/product/ipaas/wi-contrib.git/engine/conversion"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("tostring-function")

type String struct {
}

func init() {
	function.Registry(&String{})
}

func (s *String) GetName() string {
	return "tostring"
}

func (s *String) GetCategory() string {
	return "string"
}

func (s *String) Eval(in interface{}) string {
	log.Debugf("Start String function with parameters %s", in)

	switch in.(type) {
	case string:
		return in.(string)
	case float64:
		return strconv.FormatFloat(in.(float64), 'f', -1, 64)
	case int32, int8, int16, int64:
		//v := int64(in)
		return strconv.FormatInt(in.(int64), 10)
	case int:
		return strconv.Itoa(in.(int))
	case *int:
		return strconv.Itoa(*in.(*int))
	case uint, uint8, uint16, uint32, uint64:
		//v := int64(in)
		return strconv.FormatInt(in.(int64), 10)
	default:
		str, err := conversion.ConvertToString(in)
		if err != nil {
			log.Errorf("Convert to string error %s", err.Error())
		}
		return str
	}
}
