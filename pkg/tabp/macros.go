package tabp

func macroDefun(env *Env, tab ReadOnlyTable) Value {
	name := tab.Get(IntValue(1))
	if name.Type != SymbolValueType {
		return ErrorValue(Error("function name isn't a symbol"))
	}

	funcArgsTable := tab.Get(IntValue(2))
	if funcArgsTable.Type != TableValueType {
		return ErrorValue(Error("function args isn't a table"))
	}

	type funcArg struct {
		name         Symbol
		defaultValue Value
	}

	var funcArgs []funcArg
	for k, v := range funcArgsTable.AsTable().Iter() {
		if k.Type == SymbolValueType { // Key is symbol.
			funcArgs = append(funcArgs, funcArg{k.AsSymbol(), v})
		} else if v.Type == SymbolValueType { // Value is symbol
			funcArgs = append(funcArgs, funcArg{v.AsSymbol(), NilValue})
		} else {
			return ErrorValue(Error("args list of function in defun call is not a symbol"))
		}
	}

	funBody := tab.Get(IntValue(3))

	env.Defun(name.AsSymbol(), func(env *Env, argsTab ReadOnlyTable) Value {
		funcEnv := newFuncEnv(env)
		args := NewArgsTable(argsTab)

		for _, funcArg := range funcArgs {
			argVal := funcArg.defaultValue
			if v := args.consumeArg(funcArg.name); v.Type != NilValueType {
				argVal = v
			}
			funcEnv.Defvar(funcArg.name, argVal)
		}

		return funcEnv.Eval(funBody)
	})

	return NilValue
}

func macroIf(env *Env, tab ReadOnlyTable) Value {
	cond := env.Eval(tab.Get(IntValue(1)))
	if cond.Type != NilValueType {
		return tab.Get(IntValue(2))
	}

	return tab.Get(IntValue(3))
}
