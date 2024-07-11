package tabp

func macroQuote(_ *Env, tab ReadOnlyTable) Value {
	return tab.Get(1)
}

func macroQuasiQuote(env *Env, tab ReadOnlyTable) Value {
	qquoted := tab.Get(1)
	if tab, isTab := qquoted.(*Table); isTab {
		qquoted = macroQuasiQuoteRecursive(env, tab)
	}

	return qquoted
}

func macroQuasiQuoteRecursive(env *Env, tab *Table) Value {
	if symbol, isSymbol := tab.Get(0).(Symbol); isSymbol {
		if symbol == Symbol("UNQUOTE") {
			unquoted := env.Eval(tab.Get(1))
			if t, isTab := unquoted.(*Table); isTab {
				tab = t
			} else {
				return unquoted
			}
		}
	}

	for k, v := range tab.Iter() {
		var (
			newK = k
			newV = v
		)

		if tabK, isTab := k.(*Table); isTab {
			newK = macroQuasiQuoteRecursive(env, tabK)
		}

		if tabV, isTab := v.(*Table); isTab {
			newV = macroQuasiQuoteRecursive(env, tabV)
		}

		// Update table.
		if newK != k {
			tab.Set(k, nil)
		}
		tab.Set(newK, newV)
	}

	return tab
}

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
		funcEnv := newFuncEnv(env)
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

	return name
}

func macroDefvar(env *Env, tab ReadOnlyTable) Value {
	name, isSymbol := tab.Get(1).(Symbol)
	if !isSymbol {
		return Error("defvar variable name must be a symbol")
	}

	env.Defvar(name, tab.Get(2))

	return name
}

func macroIf(env *Env, tab ReadOnlyTable) Value {
	cond := env.Eval(tab.Get(1))
	if cond != nil && cond != false {
		return env.Eval(tab.Get(2))
	}

	return env.Eval(tab.Get(3))
}
