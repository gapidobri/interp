package main

import (
	"os"
)

func main() {
	path := "test.g"

	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	scanner := NewScanner(string(bytes))
	tokens, err := scanner.scanTokens()
	if err != nil {
		panic(err)
	}

	parser := NewParser(tokens)
	expression, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	interpreter := NewInterpreter()
	err = interpreter.Interpret(expression)
	if err != nil {
		panic(err)
	}
}
