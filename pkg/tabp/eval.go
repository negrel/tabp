package tabp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
)

type program struct {
	parser Parser
	env    Env
}

// Same as Eval but uses the given string as tabp program.
func EvalString(tabp string) Value {
	return Eval(bytes.NewBufferString(tabp))
}

// Eval reads, evaluates and returns a tabp program from the given reader.
func Eval(r io.Reader) Value {
	p := program{
		parser: NewParser(r),
		env:    NewEnv(),
	}

	p.env.Defvar("TABP-VERSION", "0.1.0")
	p.env.Defun("PROGN", func(tab *Table) Value {
		return tab.Get(tab.SeqLen())
	})
	p.env.Defun("IF", func(tab *Table) Value {
		cond := tab.Get(1)
		if cond != nil {
			return tab.Get(2)
		}

		return tab.Get(3)
	})
	p.env.Defun("EQ", func(tab *Table) Value {
		first := tab.Get(1)
		for i := 2; i < tab.SeqLen(); i++ {
			if !reflect.DeepEqual(first, tab.Get(i)) {
				return nil
			}
		}

		return true
	})

	p.env.Defun("PRINTF", func(tab *Table) Value {
		format, isString := tab.Get(1).(string)
		if !isString {
			return Error("format is not a string")
		}

		args := unsafeAnySlice(tab.Seq()[2:])
		fmt.Printf(string(format), args...)

		return nil
	})

	for {
		value, parseErr := p.parser.Parse()
		if parseErr.Cause != nil {
			if errors.Is(parseErr, io.EOF) {
				return nil
			}
			return parseErr
		}

		v := p.env.Eval(value)
		if err, isErr := v.(error); isErr {
			return err
		}
	}
}
