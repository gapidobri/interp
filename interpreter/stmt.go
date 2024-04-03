package interpreter

import (
	"fmt"
	"interp/environment"
	"interp/stmt"
)

func (i *Interpreter) VisitExpressionStmt(stmt *stmt.Expression) (any, error) {
	_, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt *stmt.If) (any, error) {
	value, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(value) {
		_, err = i.execute(stmt.ThenBranch)
		if err != nil {
			return nil, err
		}
	} else if stmt.ElseBranch != nil {
		_, err = i.execute(stmt.ElseBranch)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitPrintStmt(stmt *stmt.Print) (any, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}

	fmt.Println(i.stringify(value))
	return nil, nil
}

func (i *Interpreter) VisitVarStmt(stmt *stmt.Var) (any, error) {
	var value any
	if stmt.Initializer != nil {
		var err error
		value, err = i.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}

	i.environment.Define(stmt.Name.Lexeme, value)

	return nil, nil
}

func (i *Interpreter) VisitWhileStmt(stmt *stmt.While) (any, error) {
	for {
		value, err := i.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
		if !i.isTruthy(value) {
			break
		}
		_, err = i.execute(stmt.Body)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt *stmt.Block) (any, error) {
	return nil, i.executeBlock(stmt.Statements, environment.NewEnvironment(i.environment))
}
