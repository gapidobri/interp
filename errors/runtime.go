package errors

import (
	"fmt"
	"interp/token"
)

type RuntimeError struct {
	token   token.Token
	message string
}

func NewRuntimeError(token token.Token, message string) RuntimeError {
	return RuntimeError{token: token, message: message}
}

func (r RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", r.message, r.token.Line)
}
