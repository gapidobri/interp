package resolver

import "interp/ast"

func (r *Resolver) VisitBlockStmt(stmt *ast.BlockStmt) (any, error) {
	r.beginScope()
	err := r.Resolve(stmt.Statements)
	if err != nil {
		return nil, err
	}
	r.endScope()
	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt *ast.FunctionStmt) (any, error) {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	return nil, r.resolveFunction(stmt)
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
	if stmt.Value != nil {
		err := r.resolveExpr(stmt.Value)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt *ast.WhileStmt) (any, error) {
	err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return nil, err
	}

	return nil, r.resolveStmt(stmt.Body)
}

func (r *Resolver) VisitExpressionStmt(stmt *ast.ExpressionStmt) (any, error) {
	return nil, r.resolveStmt(stmt)
}

func (r *Resolver) VisitBreakStmt(stmt *ast.BreakStmt) (any, error) {
	return nil, nil
}
