package main

import (
	"fmt"
	"github.com/samber/lo"
)

type AstPrinter struct{}

func (a *AstPrinter) print(expr Expr) (*string, error) {
	value, err := expr.accept(a)
	if err != nil {
		return nil, err
	}
	return lo.ToPtr(value.(string)), nil
}

func (a *AstPrinter) visitBinaryExpr(expr *Binary) (any, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) visitGroupingExpr(expr *Grouping) (any, error) {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) visitLiteralExpr(expr *Literal) (any, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", *expr.Value.(*interface{})), nil
}

func (a *AstPrinter) visitUnaryExpr(expr *Unary) (any, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) (*string, error) {
	str := "(" + name
	for _, expr := range exprs {
		value, err := expr.accept(a)
		if err != nil {
			return nil, err
		}
		str += fmt.Sprintf(" %v", value)
	}
	str += ")"
	return &str, nil
}
