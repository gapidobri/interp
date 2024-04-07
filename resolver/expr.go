package resolver

import (
	"interp/ast"
	"interp/errors"
)

func (r *Resolver) VisitVariableExpr(expr *ast.VariableExpr) (any, error) {
	scope, ok := r.scopes.peek()
	if !ok {
		return nil, nil
	}
	if scope[expr.Name.Lexeme].defined == false {
		errors.Error(expr.Name, "Can't read local variable in its own initializer.")
	}
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr *ast.AssignExpr) (any, error) {
	err := r.resolveExpr(expr.Value)
	if err != nil {
		return nil, err
	}

	r.resolveLocal(expr.Value, expr.Name)

	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(expr *ast.BinaryExpr) (any, error) {
	err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}

	return nil, r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitCallExpr(expr *ast.CallExpr) (any, error) {
	err := r.resolveExpr(expr.Callee)
	if err != nil {
		return nil, err
	}

	for _, argument := range expr.Arguments {
		err = r.resolveExpr(argument)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitGetExpr(expr *ast.GetExpr) (any, error) {
	return nil, r.resolveExpr(expr.Object)
}

func (r *Resolver) VisitGroupingExpr(expr *ast.GroupingExpr) (any, error) {
	return nil, r.resolveExpr(expr.Expression)
}

func (r *Resolver) VisitLambdaExpr(expr *ast.LambdaExpr) (any, error) {
	return nil, r.resolveLambda(expr)
}

func (r *Resolver) VisitLiteralExpr(_ *ast.LiteralExpr) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr *ast.LogicalExpr) (any, error) {
	err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}

	return nil, r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitSetExpr(expr *ast.SetExpr) (any, error) {
	err := r.resolveExpr(expr.Value)
	if err != nil {
		return nil, err
	}
	return nil, r.resolveExpr(expr.Object)
}

func (r *Resolver) VisitThisExpr(expr *ast.ThisExpr) (any, error) {
	if r.currentClass == ClassTypeNone {
		errors.Error(expr.Keyword, "Can't use 'this' outside of a class.")
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr *ast.UnaryExpr) (any, error) {
	return nil, r.resolveExpr(expr.Right)
}
