package main

import (
	"fmt"
	"os"

	"github.com/negrel/tabp/pkg/tabp"
)

func main() {
	parser := tabp.NewParser(os.Stdin)
	value, err := parser.Parse()
	if err.Cause != nil {
		panic(err)
	}

	env := tabp.NewEnv()
	env.Defun("sprintf", func(tab *tabp.Table) tabp.Value {
		format, isString := tab.Get(1).(tabp.String)
		if !isString {
			return tabp.Error("format is not a string")
		}

		args := make([]any, 0, tab.SequenceLen())
		for i := 2; i < tab.SequenceLen(); i++ {
			args = append(args, tab.Get(i))
		}

		return tabp.String(fmt.Sprintf(string(format), args...))
	})

	env.Defun("printf", func(tab *tabp.Table) tabp.Value {
		format, isString := tab.Get(1).(tabp.String)
		if !isString {
			return tabp.Error("format is not a string")
		}

		args := make([]any, 0, tab.SequenceLen())
		for i := 2; i < tab.SequenceLen(); i++ {
			args = append(args, tab.Get(i))
		}

		fmt.Printf(string(format), args...)

		return nil
	})

	value = env.Eval(value)
	if value != nil {
		fmt.Println(tabp.Sexpr(value))
	}
}
