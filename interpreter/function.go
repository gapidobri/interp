package interpreter

import (
	"fmt"
	"interp/ast"
	"interp/environment"
)

type Function struct {
	declaration   *ast.FunctionStmt
	closure       *environment.Environment
	isInitializer bool
}

func NewFunction(declaration *ast.FunctionStmt, closure *environment.Environment, isInitializer bool) *Function {
	return &Function{declaration, closure, isInitializer}
}

func (f *Function) bind(instance *Instance) *Function {
	env := environment.NewEnvironment(f.closure)
	env.Define("this", instance)
	return NewFunction(f.declaration, env, f.isInitializer)
}

//goland:noinspection GoTypeAssertionOnErrors
func (f *Function) call(interpreter *Interpreter, arguments []any) (any, error) {
	env := environment.NewEnvironment(f.closure)
	for i, param := range f.declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	err := interpreter.executeBlock(f.declaration.Body, env)
	if err != nil {
		returnValue, ok := err.(Return)
		if !ok {
			return nil, err
		}
		if f.isInitializer {
			return f.closure.GetAt(0, "this"), nil
		}
		return returnValue.Value, nil
	}

	if f.isInitializer {
		return f.closure.GetAt(0, "this"), nil
	}

	return nil, nil
}

func (f *Function) arity() int {
	return len(f.declaration.Params)
}

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}
