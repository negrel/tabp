package tabp

import (
	"fmt"
	"unsafe"
)

const (
	ptrSize        = 8
	valueUnionSize = 2 * ptrSize
)

// ValueType enumerate possible types of a value.
type ValueType uint8

// Value define a Tabp value.
type ValueUnion [valueUnionSize]byte // 64bit

// Value define a Tabp runtime value.
type Value struct {
	Type ValueType
	Data ValueUnion
}

const (
	NilValueType ValueType = iota
	SymbolValueType
	StringValueType
	ErrorValueType
	IntValueType
	FloatValueType
	TableValueType
)

var (
	NilValue = Value{NilValueType, [valueUnionSize]byte{}}
)

// IntValue creates and returns a Value containing n.
func IntValue(n int) Value {
	data := [valueUnionSize]byte{}
	copy(data[0:ptrSize], (*(*[ptrSize]byte)(unsafe.Pointer(&n)))[:])

	return Value{IntValueType, data}
}

// FloatValue creates and returns a Value containing n.
func FloatValue(n float64) Value {
	data := [valueUnionSize]byte{}
	copy(data[0:ptrSize], (*(*[ptrSize]byte)(unsafe.Pointer(&n)))[:])

	return Value{IntValueType, data}
}

// ErrorValue creates a Value out of an error.
// If error is nil, NilValue is returned.
func ErrorValue(err error) Value {
	if err == nil {
		return NilValue
	}

	return Value{ErrorValueType, *(*[valueUnionSize]byte)(unsafe.Pointer(&err))}
}

// TableValue creates a Value out of a Table.
func TableValue(tab *Table) Value {
	data := [valueUnionSize]byte{}
	copy(data[0:ptrSize], (*(*[ptrSize]byte)(unsafe.Pointer(&tab)))[:])

	return Value{
		Type: TableValueType,
		Data: data,
	}
}

// StringValue creates and returns a Value out of a symbol.
func StringValue(str string) Value {
	stringData := unsafe.StringData(str)
	stringLen := len(str)

	data := [valueUnionSize]byte{}
	copy(data[0:ptrSize], (*(*[ptrSize]byte)(unsafe.Pointer(&stringData)))[:])
	copy(data[ptrSize:], (*(*[ptrSize]byte)(unsafe.Pointer(&stringLen)))[:])

	return Value{
		Type: StringValueType,
		Data: data,
	}
}

// SymbolValue creates and returns a Value out of a symbol.
func SymbolValue(symbol Symbol) Value {
	stringData := unsafe.StringData(string(symbol))
	stringLen := len(symbol)

	data := [valueUnionSize]byte{}
	copy(data[0:ptrSize], (*(*[ptrSize]byte)(unsafe.Pointer(&stringData)))[:])
	copy(data[ptrSize:], (*(*[ptrSize]byte)(unsafe.Pointer(&stringLen)))[:])

	return Value{
		Type: SymbolValueType,
		Data: data,
	}
}

// AsTable returns underlying table pointer.
func (v Value) AsTable() *Table {
	return *(**Table)(unsafe.Pointer(&v.Data[0]))
}

// AsSymbol returns symbol itself.
func (v Value) AsSymbol() Symbol {
	stringDataPtr := *(**byte)(unsafe.Pointer(&v.Data[0]))
	stringLen := *(*int)(unsafe.Pointer(&v.Data[ptrSize]))
	return Symbol(unsafe.String(stringDataPtr, stringLen))
}

// AsString returns symbol itself.
func (v Value) AsString() string {
	stringDataPtr := *(**byte)(unsafe.Pointer(&v.Data[0]))
	stringLen := *(*int)(unsafe.Pointer(&v.Data[ptrSize]))
	return unsafe.String(stringDataPtr, stringLen)
}

// AsError returns error itself.
func (v Value) AsError() error {
	return *(*error)(unsafe.Pointer(&v.Data[0]))
}

// AsInt returns integer value.
func (v Value) AsInt() int {
	return *(*int)(unsafe.Pointer(&v.Data[0]))
}

// AsFloat64 returns integer value.
func (v Value) AsFloat64() float64 {
	return *(*float64)(unsafe.Pointer(&v.Data[0]))
}

// AsAny returns wrapped value as any.
func (v Value) AsAny() any {
	switch v.Type {
	case ErrorValueType:
		return v.AsError()
	case FloatValueType:
		return v.AsFloat64()
	case IntValueType:
		return v.AsInt()
	case NilValueType:
		return nil
	case StringValueType:
		return v.AsString()
	case SymbolValueType:
		return v.AsSymbol()
	case TableValueType:
		return v.AsTable()
	default:
		panic(fmt.Sprintf("unexpected tabp.ValueType: %#v", v.Type))
	}
}

// ToSExpr implements SExpr.
func (v Value) ToSExpr() string {
	return Sexpr(v.AsAny())
}

// // String implements fmt.Stringer.
// func (v Value) String() string {
// 	return fmt.Sprint(v.AsAny())
// }

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
