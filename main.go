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

	for _, token := range tokens {
		fmt.Printf("%+v\n", token)
	}
}
