package interpreter

import (
	"interp/environment"
	"interp/expr"
	"interp/stmt"
)

type Interpreter struct {
	environment *environment.Environment
	globals     *environment.Environment
}

func NewInterpreter() Interpreter {
	globals := environment.NewEnvironment(nil)

	globals.Define("clock", NewClock())

	return Interpreter{
		globals:     globals,
		environment: globals,
	}
}

func (i *Interpreter) Interpret(statements []stmt.Stmt) error {
	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) evaluate(expr expr.Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt stmt.Stmt) (any, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []stmt.Stmt, environment *environment.Environment) error {
	previous := i.environment
	i.environment = environment

	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			i.environment = previous
			return err
		}
	}

	i.environment = previous
	return nil
}
