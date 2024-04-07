package interpreter

import (
	"fmt"
	"interp/ast"
	"interp/environment"
)

func (i *Interpreter) VisitExpressionStmt(stmt *ast.ExpressionStmt) (any, error) {
	_, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *ast.FunctionStmt) (any, error) {
	function := NewFunction(stmt, i.environment, false)
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) (any, error) {
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

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) (any, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}

	fmt.Println(i.stringify(value))
	return nil, nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) (any, error) {
	var value any
	if stmt.Value != nil {
		var err error
		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}

	return nil, Return{Value: value}
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) (any, error) {
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

//goland:noinspection GoTypeAssertionOnErrors
func (i *Interpreter) VisitWhileStmt(stmt *ast.WhileStmt) (any, error) {
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
			if _, ok := err.(Break); ok {
				return nil, nil
			}
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) (any, error) {
	return nil, i.executeBlock(stmt.Statements, environment.NewEnvironment(i.environment))
}

func (i *Interpreter) VisitClassStmt(stmt *ast.ClassStmt) (any, error) {
	i.environment.Define(stmt.Name.Lexeme, nil)

	methods := map[string]*Function{}
	for _, method := range stmt.Methods {
		function := NewFunction(method, i.environment, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = function
	}

	class := NewClass(stmt.Name.Lexeme, methods)
	return nil, i.environment.Assign(stmt.Name, class)
}

func (i *Interpreter) VisitBreakStmt(stmt *ast.BreakStmt) (any, error) {
	return nil, Break{}
}
