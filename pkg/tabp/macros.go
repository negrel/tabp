package tabp

func macroDefun(env *Env, tab ReadOnlyTable) Value {
	name, isSymbol := tab.Get(1).(Symbol)
	if !isSymbol {
		return Error("function name isn't a symbol")
	}

	funcArgsTable, isTable := tab.Get(2).(*Table)
	if !isTable || funcArgsTable == nil {
		return Error("function args isn't a table")
	}

	type funcArg struct {
		name         Symbol
		defaultValue Value
	}

	var funcArgs []funcArg
	for k, v := range funcArgsTable.Iter() {
		if symbol, isSymbol := k.(Symbol); isSymbol { // Key is symbol.
			funcArgs = append(funcArgs, funcArg{symbol, v})
		} else if symbol, isSymbol := v.(Symbol); isSymbol { // Value is symbol
			funcArgs = append(funcArgs, funcArg{symbol, nil})
		} else {
			return Error("args list of function in defun call is not a symbol")
		}
	}

	funBody := tab.Get(3)

	env.Defun(name, func(env *Env, argsTab ReadOnlyTable) Value {
		funcEnv := NewEnv(env)
		args := NewArgsTable(argsTab)

		for _, funcArg := range funcArgs {
			argVal := funcArg.defaultValue
			if v := args.consumeArg(funcArg.name); v != nil {
				argVal = v
			}
			funcEnv.Defvar(funcArg.name, argVal)
		}

		return funcEnv.Eval(funBody)
	})

	return nil
}

func macroIf(env *Env, tab ReadOnlyTable) Value {
	cond := env.Eval(tab.Get(1))
	if cond != nil && cond != false {
		return tab.Get(2)
	}

	return tab.Get(3)
}