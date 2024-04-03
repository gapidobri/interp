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
	return fmt.Sprintf("[line %d] %s", r.token.Line, r.message)
}

func (r RuntimeError) Token() token.Token {
	return r.token
}
