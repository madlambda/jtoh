package main

import (
	"fmt"
	"os"

	"github.com/madlambda/jtoh"
)

// Version of the code used to build the jtoh tool
var Version = ""

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <selector>\n", os.Args[0])
		fmt.Printf("example: %s :field1:nested.field2\n", os.Args[0])
		fmt.Printf("jtoh version: %q\n", Version)
		os.Exit(1)
	}
	selector := os.Args[1]
	j, err := jtoh.New(selector)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	j.Do(os.Stdin, os.Stdout)
}
