package tabp

import "fmt"

// Env define tabp execution environment.
type Env struct {
	parent *Env
	funcs  map[Symbol]func(*Env, ReadOnlyTable) Value
	macros map[Symbol]func(*Env, ReadOnlyTable) Value
	vars   map[Symbol]Value
	// True if env is dedicated to function execution.
	isFuncEnv bool
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
func NewEnv(parent *Env) Env {
	return Env{
		parent:    parent,
		funcs:     map[Symbol]func(*Env, ReadOnlyTable) Value{},
		macros:    map[Symbol]func(*Env, ReadOnlyTable) Value{},
		vars:      map[Symbol]Value{},
		isFuncEnv: false,
	}
}

func newFuncEnv(parent *Env) Env {
	return Env{
		parent:    parent,
		funcs:     map[Symbol]func(*Env, ReadOnlyTable) Value{},
		macros:    map[Symbol]func(*Env, ReadOnlyTable) Value{},
		vars:      map[Symbol]Value{},
		isFuncEnv: true,
	}
}

func (e *Env) globalEnv() *Env {
	current := e
	for {
		if current == nil || current.isFuncEnv == false {
			return current
		}

		current = e.parent
	}
}

func (e *Env) getFunc(name Symbol) func(*Env, ReadOnlyTable) Value {
	fn, ok := e.funcs[name]
	if ok {
		return fn
	}

	if e.parent != nil {
		return e.parent.getFunc(name)
	}

	return nil
}

func (e *Env) getMacro(name Symbol) func(*Env, ReadOnlyTable) Value {
	fn, ok := e.macros[name]
	if ok {
		return fn
	}

	if e.parent != nil {
		return e.parent.getMacro(name)
	}

	return nil
}

func (e *Env) getVar(name Symbol) Value {
	fn, ok := e.vars[name]
	if ok {
		return fn
	}

	if e.parent != nil {
		return e.parent.getVar(name)
	}

	return NilValue
}

// Defun define a function in the environment.
func (e *Env) Defun(name Symbol, fn func(*Env, ReadOnlyTable) Value) {
	e.funcs[name] = fn
}

// Defmacro define a macro in the environment.
func (e *Env) Defmacro(name Symbol, fn func(*Env, ReadOnlyTable) Value) {
	e.macros[name] = fn
}

// Defvar define a variable in the environment.
func (e *Env) Defvar(name Symbol, v Value) {
	e.vars[name] = v
}

// Eval evaluates the given value within the environment and returns a new value.
func (e *Env) Eval(value Value) Value {
	fn := func() Value {
		if value.Type == TableValueType {
			tab := value.AsTable()
			name := tab.Get(IntValue(0))
			if name.Type == SymbolValueType {
				symbol := name.AsSymbol()
				macro := e.getMacro(symbol)
				// Function.
				if macro == nil {
					return e.evalFunc(tab, symbol)
				}

				// Macro.
				return e.Eval(macro(e, tab))
			}

			return ErrorValue(EvalError{Cause: Error("function/macro name is not a symbol"), Expr: value})
		} else if value.Type == SymbolValueType {
			return e.getVar(value.AsSymbol())
		}

		return value
	}

	res := fn()
	return res
}

func (e *Env) evalFunc(tab ReadOnlyTable, fnName Symbol) Value {
	// TODO: fix Expr in ErrorValue
	fn := e.getFunc(fnName)
	if fn == nil {
		return ErrorValue(EvalError{Cause: Error("function not found"), Expr: SymbolValue(fnName)})
	}

	var args Table
	for k, v := range tab.Iter() {
		// Copy function name symbol.
		if v.Type == SymbolValueType && v.AsSymbol() == fnName {
			args.Set(k, v)
			continue
		}

		arg := e.Eval(v)
		if arg.Type == ErrorValueType {
			return ErrorValue(EvalError{Cause: arg.AsError(), Expr: v})
		}

		args.Set(k, arg)
	}

	result := fn(e.globalEnv(), &args)
	if result.Type == ErrorValueType {
		return ErrorValue(EvalError{Cause: result.AsError(), Expr: SymbolValue(fnName)})
	}

	return result
}
