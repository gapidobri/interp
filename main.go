package main

import (
	"fmt"
	"os"
)

func main() {
	path := "test.g"

	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	scanner := NewScanner(string(bytes))

	tokens := scanner.scanTokens()

	parser := NewParser(tokens)
	expression := parser.parse()

	printer := AstPrinter{}
	fmt.Println(printer.print(expression))
}
