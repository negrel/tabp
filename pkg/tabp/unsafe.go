//go:build go1.20

package tabp

import (
	"unsafe"
)

// UnsafeBytes returns a byte slice without allocation.
func UnsafeBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeString returns a string without allocation.
func UnsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// unsafeAnySlice converts slice of []T to []any without allocation.
func unsafeAnySlice(v []Value) []any {
	sliceData := unsafe.SliceData(v)
	return unsafe.Slice((*any)(sliceData), len(v))
}
