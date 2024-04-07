package resolver

import (
	"interp/ast"
	"interp/errors"
)

func (r *Resolver) VisitBlockStmt(stmt *ast.BlockStmt) (any, error) {
	r.beginScope()
	err := r.Resolve(stmt.Statements)
	if err != nil {
		return nil, err
	}
	r.endScope()
	return nil, nil
}

func (r *Resolver) VisitClassStmt(stmt *ast.ClassStmt) (any, error) {
	enclosingClass := r.currentClass
	r.currentClass = ClassTypeClass

	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.beginScope()
	scope, ok := r.scopes.peek()
	if !ok {
		return nil, nil
	}

	scope["this"] = &varState{defined: true, resolved: true}

	for _, method := range stmt.Methods {
		declaration := FunctionTypeMethod
		if method.Name.Lexeme == "init" {
			declaration = FunctionTypeInitializer
		}
		err := r.resolveFunction(method, declaration)
		if err != nil {
			return nil, err
		}
	}

	r.endScope()

	r.currentClass = enclosingClass
	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt *ast.FunctionStmt) (any, error) {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	return nil, r.resolveFunction(stmt, FunctionTypeFunction)
}

func (r *Resolver) VisitIfStmt(stmt *ast.IfStmt) (any, error) {
	err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return nil, err
	}

	err = r.resolveStmt(stmt.ThenBranch)
	if err != nil {
		return nil, err
	}

	if stmt.ElseBranch != nil {
		err = r.resolveStmt(stmt.ElseBranch)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitVarStmt(stmt *ast.VarStmt) (any, error) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		err := r.resolveExpr(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}
	r.define(stmt.Name)
	return nil, nil
}

func (r *Resolver) VisitPrintStmt(stmt *ast.PrintStmt) (any, error) {
	return nil, r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitReturnStmt(stmt *ast.ReturnStmt) (any, error) {
	if r.currentFunction == FunctionTypeNone {
		errors.Error(stmt.Keyword, "Can't return from top-level code.")
	}

	if stmt.Value != nil {
		if r.currentFunction == FunctionTypeInitializer {
			errors.Error(stmt.Keyword, "Can't return a value from an initializer.")
		}
		err := r.resolveExpr(stmt.Value)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt *ast.WhileStmt) (any, error) {
	enclosingLoop := r.inLoop
	r.inLoop = true

	err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return nil, err
	}

	err = r.resolveStmt(stmt.Body)
	if err != nil {
		return nil, err
	}

	r.inLoop = enclosingLoop
	return nil, nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ast.ExpressionStmt) (any, error) {
	return nil, r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitBreakStmt(stmt *ast.BreakStmt) (any, error) {
	if !r.inLoop {
		errors.Error(stmt.Keyword, "Can't break outside loop")
	}
	return nil, nil
}
