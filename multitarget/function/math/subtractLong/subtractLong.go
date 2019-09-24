package mathfunc

import (
	"math/big"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
)

type Subtract struct {
}

func init() {
	function.Registry(&Subtract{})
}

func (s *Subtract) GetName() string {
	return "subtractLong"
}

func (s *Subtract) GetCategory() string {
	return "math"
}

func (s *Subtract) Eval(vals ...int64) (int64, error) {
	input := vals[1:]
	result := big.NewInt(vals[0])
	for _, v := range input {
		result = result.Sub(result, big.NewInt(v))
	}

	return result.Int64(), nil
}
