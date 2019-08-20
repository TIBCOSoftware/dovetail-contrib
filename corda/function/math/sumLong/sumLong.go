package mathfunc

import (
	"math/big"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
)

type Sum struct {
}

func init() {
	function.Registry(&Sum{})
}

func (s *Sum) GetName() string {
	return "sumLong"
}

func (s *Sum) GetCategory() string {
	return "math"
}

func (s *Sum) Eval(vals ...int64) (int64, error) {

	sum := big.NewInt(0)
	for _, v := range vals {
		sum = sum.Add(sum, big.NewInt(v))
	}

	return sum.Int64(), nil
}
