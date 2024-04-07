package errors

import (
	"fmt"
	"interp/token"
)

const lineCount = 2

const red = "\033[0;31m"
const grey = "\033[0;37m"
const none = "\033[0m"

var HadError = false

func Error(token token.Token, message string) {
	HadError = true
	fmt.Printf("[line %d] %s\n", token.Line, message)
}

func Warning(token token.Token, message string) {
	fmt.Printf("[line %d] %s\n", token.Line, message)
}
