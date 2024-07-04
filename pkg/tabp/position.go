package tabp

import "fmt"

// Position holds byte, line and column position of a cursor in a byte stream.
type Position struct {
	byte, line, col int
}

// NewPosition creates a new position starting on line 1.
func NewPosition() Position {
	return Position{
		byte: 0,
		line: 1,
		col:  0,
	}
}

// String implements fmt.Stringer.
func (p Position) String() string {
	return fmt.Sprintf("%v:%v (%v bytes)", p.line, p.col, p.byte)
}
