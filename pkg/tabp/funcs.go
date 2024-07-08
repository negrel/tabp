package tabp

import (
	"fmt"
)

func fnAdd(_ *Env, tab ReadOnlyTable) Value {
	first := tab.Get(IntValue(1))
	second := tab.Get(IntValue(2))

	if !(first.Type == IntValueType || first.Type == FloatValueType) ||
		!(second.Type == IntValueType || second.Type == FloatValueType) {
		return ErrorValue(Error("can't compare non number type"))
	}

	if first.Type == IntValueType && second.Type == IntValueType {
		return IntValue(first.AsInt() + second.AsInt())
	}

	panic(fmt.Errorf("%v %v", first.Type, second.Type))
}

func fnSub(_ *Env, tab ReadOnlyTable) Value {
	first := tab.Get(IntValue(1))
	second := tab.Get(IntValue(2))

	if !(first.Type == IntValueType || first.Type == FloatValueType) ||
		!(second.Type == IntValueType || second.Type == FloatValueType) {
		return ErrorValue(Error("can't compare non number type"))
	}

	if first.Type == IntValueType && second.Type == IntValueType {
		return IntValue(first.AsInt() - second.AsInt())
	}

	panic(fmt.Errorf("%v %v", first.Type, second.Type))
}

func fnEq(_ *Env, tab ReadOnlyTable) Value {
	// first := tab.Get(1)
	// for i := 2; i < tab.SeqLen(); i++ {
	// 	if !reflect.DeepEqual(first, tab.Get(i)) {
	// 		return nil
	// 	}
	// }
	//
	// return true
	return NilValue
}

func fnLt(_ *Env, tab ReadOnlyTable) Value {
	first := tab.Get(IntValue(1))
	second := tab.Get(IntValue(2))

	if !(first.Type == IntValueType || first.Type == FloatValueType) ||
		!(second.Type == IntValueType || second.Type == FloatValueType) {
		return ErrorValue(Error("can't compare non number type"))
	}

	if first.Type == IntValueType && second.Type == IntValueType {
		if first.AsInt() < second.AsInt() {
			return IntValue(1)
		} else {
			return NilValue
		}
	}

	panic(fmt.Errorf("%v %v", first.Type, second.Type))
}

func fnLe(_ *Env, tab ReadOnlyTable) Value {
	// first := tab.Get(1)
	// second := tab.Get(2)
	//
	// a, aF, ok := toNumber(first)
	// if !ok {
	// 	return Error("can't compare non number type")
	// }
	//
	// b, bF, ok := toNumber(second)
	// if !ok {
	// 	return Error("can't compare non number type")
	// }
	//
	// if aF != 0 || bF != 0 {
	// 	return aF+float64(a) <= bF+float64(b)
	// }
	//
	// return a <= b
	return NilValue
}

func fnGt(_ *Env, tab ReadOnlyTable) Value {
	// first := tab.Get(1)
	// second := tab.Get(2)
	//
	// a, aF, ok := toNumber(first)
	// if !ok {
	// 	return Error("can't compare non number type")
	// }
	//
	// b, bF, ok := toNumber(second)
	// if !ok {
	// 	return Error("can't compare non number type")
	// }
	//
	// if aF != 0 || bF != 0 {
	// 	return aF+float64(a) > bF+float64(b)
	// }
	//
	// return a > b
	return NilValue
}

func fnGe(_ *Env, tab ReadOnlyTable) Value {
	// first := tab.Get(1)
	// second := tab.Get(2)
	//
	// a, aF, ok := toNumber(first)
	// if !ok {
	// 	return Error("can't compare non number type")
	// }
	//
	// b, bF, ok := toNumber(second)
	// if !ok {
	// 	return Error("can't compare non number type")
	// }
	//
	// if aF != 0 || bF != 0 {
	// 	return aF+float64(a) >= bF+float64(b)
	// }
	//
	// return a >= b
	return NilValue
}

func fnPrintf(_ *Env, tab ReadOnlyTable) Value {
	format := tab.Get(IntValue(1))
	if format.Type != StringValueType {
		return ErrorValue(Error("format is not a string"))
	}

	args := make([]any, len(tab.Seq()[2:]))
	for i, arg := range tab.Seq()[2:] {
		args[i] = arg.AsAny()
	}
	fmt.Printf(format.AsString(), args...)

	return NilValue
}

func fnSprintf(_ *Env, tab ReadOnlyTable) Value {
	format := tab.Get(IntValue(1))
	if format.Type != StringValueType {
		return ErrorValue(Error("format is not a string"))
	}

	args := make([]any, len(tab.Seq()[2:]))
	for i, arg := range tab.Seq()[2:] {
		args[i] = arg.AsAny()
	}
	str := fmt.Sprintf(format.AsString(), args...)

	return StringValue(str)
}

func fnProgn(_ *Env, tab ReadOnlyTable) Value {
	return tab.Get(IntValue(tab.SeqLen() - 1))
}
