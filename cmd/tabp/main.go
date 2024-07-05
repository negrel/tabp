package main

import (
	"fmt"
	"os"

	"github.com/negrel/tabp/pkg/tabp"
)

func main() {
	val := tabp.Eval(os.Stdin)
	if val != nil {
		if err, isErr := val.(error); isErr {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		fmt.Println(tabp.Sexpr(val))
	}
}
