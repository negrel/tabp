package tabp

import (
	"fmt"
	"reflect"
)

func fnAdd(_ *Env, tab ReadOnlyTable) Value {
	args := unsafeAnySlice(tab.Seq()[1:])
	if len(args) < 1 {
		return Error("no argument provided")
	}

	sum := 0
	sumF := 0.0
	for _, v := range args {
		vInt, vFloat, ok := toNumber(v)
		if !ok {
			return Error("can't add non number type")
		}

		sum += vInt
		sumF += vFloat
	}

	if sumF != 0.0 {
		return sumF + float64(sum)
	}

	return sum
}

func fnSub(_ *Env, tab ReadOnlyTable) Value {
	args := unsafeAnySlice(tab.Seq()[1:])
	if len(args) < 1 {
		return Error("no argument provided")
	}

	base, baseF, ok := toNumber(args[0])
	if !ok {
		return Error("can't substract non number type")
	}

	result := 0
	resultF := 0.0

	for _, v := range args[1:] {
		vInt, vFloat, ok := toNumber(v)
		if !ok {
			return Error("can't substract non number type")
		}

		result -= vInt
		resultF -= vFloat
	}

	if baseF != 0.0 || resultF != 0.0 {
		return (float64(base) + baseF) + (resultF - float64(result))
	}

	return base + result
}

func fnEq(_ *Env, tab ReadOnlyTable) Value {
	first := tab.Get(1)
	for i := 2; i < tab.SeqLen(); i++ {
		if !reflect.DeepEqual(first, tab.Get(i)) {
			return nil
		}
	}

	return true
}

func fnLt(_ *Env, tab ReadOnlyTable) Value {
	first := tab.Get(1)
	second := tab.Get(2)

	a, aF, ok := toNumber(first)
	if !ok {
		return Error("can't compare non number type")
	}

	b, bF, ok := toNumber(second)
	if !ok {
		return Error("can't compare non number type")
	}

	if aF != 0 || bF != 0 {
		return aF+float64(a) < bF+float64(b)
	}

	return a < b
}

func fnLe(_ *Env, tab ReadOnlyTable) Value {
	first := tab.Get(1)
	second := tab.Get(2)

	a, aF, ok := toNumber(first)
	if !ok {
		return Error("can't compare non number type")
	}

	b, bF, ok := toNumber(second)
	if !ok {
		return Error("can't compare non number type")
	}

	if aF != 0 || bF != 0 {
		return aF+float64(a) <= bF+float64(b)
	}

	return a <= b
}

func fnGt(_ *Env, tab ReadOnlyTable) Value {
	first := tab.Get(1)
	second := tab.Get(2)

	a, aF, ok := toNumber(first)
	if !ok {
		return Error("can't compare non number type")
	}

	b, bF, ok := toNumber(second)
	if !ok {
		return Error("can't compare non number type")
	}

	if aF != 0 || bF != 0 {
		return aF+float64(a) > bF+float64(b)
	}

	return a > b
}

func fnGe(_ *Env, tab ReadOnlyTable) Value {
	first := tab.Get(1)
	second := tab.Get(2)

	a, aF, ok := toNumber(first)
	if !ok {
		return Error("can't compare non number type")
	}

	b, bF, ok := toNumber(second)
	if !ok {
		return Error("can't compare non number type")
	}

	if aF != 0 || bF != 0 {
		return aF+float64(a) >= bF+float64(b)
	}

	return a >= b
}

func fnPrintf(_ *Env, tab ReadOnlyTable) Value {
	format, isString := tab.Get(1).(string)
	if !isString {
		return Error("format is not a string")
	}

	args := unsafeAnySlice(tab.Seq()[2:])
	fmt.Printf(string(format), args...)

	return nil
}

func fnSprintf(_ *Env, tab ReadOnlyTable) Value {
	format, isString := tab.Get(1).(string)
	if !isString {
		return Error("format is not a string")
	}

	args := unsafeAnySlice(tab.Seq()[2:])
	return fmt.Sprintf(string(format), args...)
}

func fnProgn(_ *Env, tab ReadOnlyTable) Value {
	return tab.Get(tab.SeqLen() - 1)
}
