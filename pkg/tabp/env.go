package tabp

import "fmt"

// Env define tabp execution environment.
type Env struct {
	parent *Env
	fun    map[Symbol]func(*Table) Value
	vars   map[Symbol]Value
}

// EvalError define errors returned when evaluating a Tabp S-Expression.
type EvalError struct {
	Cause error
	Expr  Value
}

// Error implements error.
func (ee EvalError) Error() string {
	return fmt.Sprintf("failed to evaluate expression %q: %v", Sexpr(ee.Expr), ee.Cause)
}

// NewEnv creates and returns a new blank environment.
func NewEnv() Env {
	return Env{
		parent: nil,
		fun:    map[Symbol]func(*Table) Value{},
		vars:   map[Symbol]Value{},
	}
}

// Defun define a function in the environment.
func (e *Env) Defun(name Symbol, fn func(*Table) Value) {
	e.fun[name] = fn
}

// Defvar define a variable in the environment.
func (e *Env) Defvar(name Symbol, v Value) {
	e.vars[name] = v
}

// Eval evaluates the given value within the environment and returns a new value.
func (e *Env) Eval(v Value) Value {
	if v == nil {
		return v
	}

	switch value := v.(type) {
	case Symbol:
		return e.vars[value]

	case *Table:
		name := value.Get(0)
		if fnSym, isSymbol := name.(Symbol); isSymbol {
			fn, ok := e.fun[fnSym]
			if !ok {
				return EvalError{Cause: Error("function not found"), Expr: v}
			}

			var (
				err   error
				isErr bool
			)
			value.Map(func(k, v Value) (Value, bool) {
				if k == 0 {
					return v, false
				}
				v = e.Eval(v)
				err, isErr = v.(error)
				return v, isErr
			})

			if err != nil {
				return EvalError{Cause: err, Expr: value}
			}

			result := fn(value)
			if err, isErr := result.(error); isErr {
				return EvalError{Cause: err, Expr: value}
			}

			return result
		}
		return EvalError{Cause: Error("function name is not a symbol"), Expr: v}

	default:
		return v
	}
}
