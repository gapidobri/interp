package environment

import (
	"fmt"
	"interp/errors"
	"interp/token"
)

type Environment struct {
	Enclosing *Environment
	values    map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Enclosing: enclosing,
		values:    map[string]any{},
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name token.Token) (any, error) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	return nil, errors.NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
}

func (e *Environment) Assign(name token.Token, value any) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.Assign(name, value)
	}

	return errors.NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
}
