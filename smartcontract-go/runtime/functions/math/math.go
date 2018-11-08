package math

import "github.com/shopspring/decimal"

func SumLong(args []interface{}) int64 {
	sum := int64(0)
	for _, v := range args {
		sum = sum + int64(v.(float64))
	}

	return sum
}

func SumInt(args []interface{}) int32 {
	sum := int32(0)
	for _, v := range args {
		sum = sum + int32(v.(float64))
	}

	return sum
}

func SumDouble(args []interface{}, precision int32, scale int32, rounding string) (string, error) {
	values := make([]decimal.Decimal, len(args))
	for _, v := range args {
		dec, err := decimal.NewFromString(v.(string))
		if err != nil {
			return "", err
		}

		values = append(values, dec)
	}

	sum := decimal.Sum(decimal.Zero, values...)

	return FormatDecimal(sum, precision, scale, rounding), nil
}
func avg(sum decimal.Decimal, count int, precision int32, scale int32, rounding string) string {
	return FormatDecimal(sum.DivRound(decimal.RequireFromString(string(count)), precision), precision, scale, rounding)
}
func AvgLong(args []interface{}, precision int32, scale int32, rounding string) string {
	sum := SumLong(args)
	return avg(decimal.RequireFromString(string(sum)), len(args), precision, scale, rounding)
}

func AvgInt(args []interface{}, precision int32, scale int32, rounding string) string {
	sum := SumInt(args)
	return avg(decimal.RequireFromString(string(sum)), len(args), precision, scale, rounding)
}

func AvgDouble(args []interface{}, precision int32, scale int32, rounding string) (result string, err error) {
	sum, err := SumDouble(args, precision, scale, rounding)
	if err != nil {
		return "", err
	}

	return avg(decimal.RequireFromString(string(sum)), len(args), precision, scale, rounding), nil
}

func MinLong(args []interface{}) int64 {
	result := args[0].(int64)
	for _, v := range args {
		value := v.(int64)
		if value < result {
			result = value
		}
	}

	return result
}

func MinInt(args []interface{}) int32 {
	result := args[0].(int32)
	for _, v := range args {
		value := v.(int32)
		if value < result {
			result = value
		}
	}

	return result
}

func MinDouble(args []interface{}, precision int32, scale int32, rounding string) (result string, err error) {
	min, _ := decimal.NewFromString(args[0].(string))
	for _, v := range args {
		dec, err := decimal.NewFromString(v.(string))
		if err != nil {
			return "", err
		}

		if dec.LessThan(min) {
			min = dec
		}
	}

	return FormatDecimal(min, precision, scale, rounding), nil
}

func MaxLong(args []interface{}) int64 {
	result := args[0].(int64)
	for _, v := range args {
		value := v.(int64)
		if value > result {
			result = value
		}
	}

	return result
}

func MaxInt(args []interface{}) int32 {
	result := args[0].(int32)
	for _, v := range args {
		value := v.(int32)
		if value > result {
			result = value
		}
	}

	return result
}

func MaxDouble(args []interface{}, precision int32, scale int32, rounding string) (result string, err error) {
	max, _ := decimal.NewFromString(args[0].(string))
	for _, v := range args {
		dec, err := decimal.NewFromString(v.(string))
		if err != nil {
			return "", err
		}

		if dec.GreaterThan(max) {
			max = dec
		}
	}

	return FormatDecimal(max, precision, scale, rounding), nil
}

func FormatDecimal(dec decimal.Decimal, precision int32, scale int32, rounding string) string {
	switch rounding {
	case "CEILING":
		dec = dec.Ceil()
		break
	case "FLOOR":
		dec = dec.Floor()
		break
	case "HALF_EVEN":
		dec = dec.RoundBank(scale)
	default:
		dec = dec.Round(scale)
	}
	dec = dec.Truncate(precision)
	return dec.StringFixed(scale)
}
