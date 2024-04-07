package interpreter

import (
	"fmt"
	"interp/errors"
	"interp/token"
)

type Instance struct {
	class  *Class
	fields map[string]any
}

func NewInstance(class *Class) *Instance {
	return &Instance{
		class,
		map[string]any{},
	}
}

func (i *Instance) String() string {
	return i.class.Name + " instance"
}

func (i *Instance) get(name token.Token) (any, error) {
	if field, exists := i.fields[name.Lexeme]; exists {
		return field, nil
	}

	method := i.class.findMethod(name.Lexeme)
	if method != nil {
		return method.bind(i), nil
	}

	return nil, errors.NewRuntimeError(name, fmt.Sprintf("Undefined property '%s'.", name.Lexeme))
}

func (i *Instance) set(name token.Token, value any) {
	i.fields[name.Lexeme] = value
}
