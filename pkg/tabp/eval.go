package tabp

import (
	"bytes"
	"errors"
	"io"
)

type program struct {
	parser Parser
	env    *Env
}

// Same as Eval but uses the given string as tabp program.
func EvalString(tabp string) Value {
	return Eval(bytes.NewBufferString(tabp))
}

// Eval reads, evaluates and returns a tabp program from the given reader.
func Eval(r io.Reader) Value {
	p := program{
		parser: NewParser(r),
		env:    NewEnv(nil),
	}

	// Variables.
	p.env.Defvar("TABP-VERSION", "0.1.0")

	// Macros.
	p.env.Defmacro("DEFUN", macroDefun)
	p.env.Defmacro("IF", macroIf)

	// Functions.
	p.env.Defun("PROGN", fnProgn)
	p.env.Defun("EQ", fnEq)
	p.env.Defun("LT", fnLt)
	p.env.Defun("LE", fnLe)
	p.env.Defun("GT", fnGt)
	p.env.Defun("GE", fnGe)
	p.env.Defun("PRINTF", fnPrintf)
	p.env.Defun("SPRINTF", fnSprintf)
	p.env.Defun("ADD", fnAdd)
	p.env.Defun("SUB", fnSub)
	// p.env.Defun("MUL", func(tab *Table) Value {})
	// p.env.Defun("DIV", func(tab *Table) Value {})

	for {
		value, parseErr := p.parser.Parse()
		if parseErr.Cause != nil {
			if errors.Is(parseErr, io.EOF) {
				return nil
			}
			return parseErr
		}

		result := p.env.Eval(value)
		if err, isErr := result.(error); isErr {
			return err
		}
	}
}
