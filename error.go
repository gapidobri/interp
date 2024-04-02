package main

import "fmt"

func errorLine(line int, message string) {
	report(line, "", message)
}

func errorToken(token *Token, message string) {
	if token.Type == EOF {
		report(token.Line, " at end", message)
	} else {
		report(token.Line, fmt.Sprintf(" at '%s'", token.Lexeme), message)
	}
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
}
