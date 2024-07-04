package tabp

import "strconv"

// Value define a Tabp value.
type Value any

// String is a string that implements Value.
type String string

// ToSExpr implements SExpr.
func (s String) ToSExpr() string {
	return strconv.Quote(string(s))
}

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
