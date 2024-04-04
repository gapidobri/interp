package interpreter

import (
	"interp/ast"
	"interp/environment"
)

type Interpreter struct {
	environment *environment.Environment
	globals     *environment.Environment
	locals      map[ast.Expr]int
}

func NewInterpreter() Interpreter {
	globals := environment.NewEnvironment(nil)

	globals.Define("clock", NewClock())

	return Interpreter{
		globals:     globals,
		environment: globals,
		locals:      map[ast.Expr]int{},
	}
}

func (i *Interpreter) Interpret(statements []ast.Stmt) error {
	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) evaluate(expr ast.Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt ast.Stmt) (any, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) Resolve(expr ast.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment *environment.Environment) error {
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
