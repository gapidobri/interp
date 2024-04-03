package main

import (
	"fmt"
	"interp/interpreter"
	"interp/parser"
	"interp/scanner"
	"os"
)

func main() {
	path := "test.g"

	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	scan := scanner.NewScanner(string(bytes))
	tokens, err := scan.ScanTokens()
	if err != nil {
		fmt.Println(err)
		return
	}

	par := parser.NewParser(tokens)
	statements, err := par.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	inter := interpreter.NewInterpreter()
	err = inter.Interpret(statements)
	if err != nil {
		fmt.Println(err)
		return
	}
}
