package interpreter

import (
	"fmt"
	"interp/environment"
	"interp/stmt"
)

type Function struct {
	declaration *stmt.Function
	closure     *environment.Environment
}

func NewFunction(declaration *stmt.Function, closure *environment.Environment) Callable {
	return Function{
		declaration: declaration,
		closure:     closure,
	}
}

//goland:noinspection GoTypeAssertionOnErrors
func (f Function) call(interpreter *Interpreter, arguments []any) (any, error) {
	env := environment.NewEnvironment(f.closure)
	for i, param := range f.declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	err := interpreter.executeBlock(f.declaration.Body, env)
	if err != nil {
		if returnValue, ok := err.(Return); ok {
			return returnValue.Value, nil
		}
		return nil, err
	}

	return nil, nil
}

func (f Function) arity() int {
	return len(f.declaration.Params)
}

func (f Function) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}
