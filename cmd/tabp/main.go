package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/negrel/tabp/pkg/tabp"
)

func main() {
	parser := tabp.NewParser(os.Stdin)

	env := tabp.NewEnv()
	env.Defun("sprintf", func(tab *tabp.Table) tabp.Value {
		format, isString := tab.Get(1).(string)
		if !isString {
			return tabp.Error("format is not a string")
		}

		args := make([]any, 0, tab.SeqLen())
		for i := 2; i < tab.SeqLen(); i++ {
			args = append(args, tab.Get(i))
		}

		return fmt.Sprintf(string(format), args...)
	})

	env.Defun("printf", func(tab *tabp.Table) tabp.Value {
		format, isString := tab.Get(1).(string)
		if !isString {
			return tabp.Error("format is not a string")
		}

		args := make([]any, 0, tab.SeqLen())
		for i := 2; i < tab.SeqLen(); i++ {
			args = append(args, tab.Get(i))
		}

		fmt.Printf(string(format), args...)

		return nil
	})

	for {
		fmt.Print("tabp> ")
		value, err := parser.Parse()
		if err.Cause != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			fmt.Fprintln(os.Stderr, err.Error())
		}

		value = env.Eval(value)
		if value != nil {
			if err, isErr := value.(error); isErr {
				fmt.Fprintln(os.Stderr, err.Error())
			} else {
				fmt.Println(tabp.Sexpr(value))
			}
		}
	}
}
