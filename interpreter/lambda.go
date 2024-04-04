package interpreter

import (
	"interp/ast"
	"interp/environment"
)

type Lambda struct {
	expression *ast.LambdaExpr
	closure    *environment.Environment
}

func NewLambda(expression *ast.LambdaExpr, closure *environment.Environment) Callable {
	return Lambda{expression: expression, closure: closure}
}

func (l Lambda) call(interpreter *Interpreter, arguments []any) (any, error) {
	env := environment.NewEnvironment(l.closure)
	for i, param := range l.expression.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	err := interpreter.executeBlock(l.expression.Body, env)
	if err != nil {
		if returnValue, ok := err.(Return); ok {
			return returnValue.Value, nil
		}
		return nil, err
	}

	return nil, nil
}

func (l Lambda) arity() int {
	return len(l.expression.Params)
}
