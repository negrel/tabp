package tabp

import (
	"fmt"
	"strconv"
)

// SExpr define any value that can be expressed as an S-Expression.
type SExpr interface {
	ToSExpr() string
}

// Sexpr converts any value to an S-Expression.
func Sexpr(v any) string {
	switch value := v.(type) {
	case SExpr:
		return value.ToSExpr()
	case bool,
		float32, float64,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprint(value)

	case nil:
		return "()"

	case string:
		return strconv.Quote(value)

	case EvalError:
		return value.Error()

	default:
		panic(fmt.Errorf("SExpr not implemented for %T %v", value, value))
	}
}
