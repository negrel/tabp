package tabp

// Value define a Tabp value.
type Value any

// Symbol is a string that implements Value.
type Symbol string

// ToSExpr implements SExpr.
func (s Symbol) ToSExpr() string {
	return string(s)
}

// Error define an error message.
type Error string

// Error implements error.
func (e Error) Error() string {
	return string(e)
}

// ToSExpr implements SExpr.
func (e Error) ToSExpr() string {
	return string(e)
}

func toNumber(x Value) (int, float64, bool) {
	switch value := x.(type) {
	case int:
		return value, 0.0, true
	case int8:
		return int(value), 0.0, true
	case int16:
		return int(value), 0.0, true
	case int32:
		return int(value), 0.0, true
	case int64:
		return int(value), 0.0, true
	case uint:
		return int(value), 0.0, true
	case uint8:
		return int(value), 0.0, true
	case uint16:
		return int(value), 0.0, true
	case uint32:
		return int(value), 0.0, true
	case uint64:
		return int(value), 0.0, true

	case float32:
		return 0, float64(value), true
	case float64:
		return 0, value, true

	default:
		return 0, 0, false
	}
}
