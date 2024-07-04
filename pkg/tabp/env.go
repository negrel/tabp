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
	}
}

// Defun define a function in the environment.
func (e *Env) Defun(name string, fn func(*Table) Value) {
	e.fun[Symbol(name)] = fn
}

// Eval evaluates the given value within the environment and returns a new value.
func (e *Env) Eval(v Value) Value {
	if v == nil {
		return v
	}

	switch value := v.(type) {
	case *Table:
		fnName := value.Get(0)
		if fnSym, isSymbol := fnName.(Symbol); isSymbol {
			fn, ok := e.fun[fnSym]
			if !ok {
				return EvalError{Cause: Error("function not found"), Expr: v}
			}

			args := Table{}
			for _, v := range value.kv {
				val := e.Eval(v.Value)
				if err, isErr := val.(error); isErr {
					return EvalError{Cause: err, Expr: value}
				}
				args.Set(v.Key, val)
			}

			result := fn(&args)
			if err, isErr := result.(error); isErr {
				return EvalError{Cause: err, Expr: v}
			}

			return result
		}
		return EvalError{Cause: Error("function name is not a symbol"), Expr: v}

	default:
		return v
	}
}
