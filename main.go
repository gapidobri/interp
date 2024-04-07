package main

import (
	"fmt"
	"interp/errors"
	"interp/interpreter"
	"interp/parser"
	"interp/resolver"
	"interp/scanner"
	"os"
)

//goland:noinspection GoTypeAssertionOnErrors,GoTypeAssertionOnErrors
func main() {
	path := "test.g"

	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	source := string(bytes)

	scan := scanner.NewScanner(source)
	tokens, err := scan.ScanTokens()
	if err != nil {
		fmt.Println(err)
		return
	}

	par := parser.NewParser(tokens)
	statements, err := par.Parse()
	if err, ok := err.(errors.SyntaxError); ok {
		err.Print(&source)
	}
	if errors.HadError {
		return
	}

	inter := interpreter.NewInterpreter()

	res := resolver.NewResolver(&inter)
	if err := res.Resolve(statements); err != nil {
		fmt.Println(err)
	}

	if errors.HadError {
		return
	}

	err = inter.Interpret(statements)
	if err, ok := err.(errors.RuntimeError); ok {
		err.Print(&source)
	}
}
