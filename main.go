package main

import (
	"interp/interpreter"
	"interp/parser"
	"interp/scanner"
	"os"
)

func main() {
	path := "test.g"

	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	scan := scanner.NewScanner(string(bytes))
	tokens, err := scan.ScanTokens()
	if err != nil {
		panic(err)
	}

	par := parser.NewParser(tokens)
	statements, err := par.Parse()
	if err != nil {
		panic(err)
	}

	inter := interpreter.NewInterpreter()
	err = inter.Interpret(statements)
	if err != nil {
		panic(err)
	}

}
