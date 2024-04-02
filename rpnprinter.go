package main

import (
	"fmt"
	"github.com/samber/lo"
)

type RPNPrinter struct{}

func (a *RPNPrinter) print(expr Expr) (*string, error) {
	value, err := expr.accept(a)
	if err != nil {
		return nil, err
	}
	return lo.ToPtr(value.(string)), nil
}

func (a *RPNPrinter) visitBinaryExpr(expr *Binary) (any, error) {
	left, err := expr.Left.accept(a)
	if err != nil {
		return nil, err
	}

	right, err := expr.Right.accept(a)
	if err != nil {
		return nil, err
	}

	return left.(string) + " " +
		right.(string) + " " +
		expr.Operator.Lexeme, nil
}

func (a *RPNPrinter) visitGroupingExpr(expr *Grouping) (any, error) {
	return expr.Expression.accept(a)
}

func (a *RPNPrinter) visitLiteralExpr(expr *Literal) (any, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.Value), nil
}

func (a *RPNPrinter) visitUnaryExpr(expr *Unary) (any, error) {
	return "", nil
}
