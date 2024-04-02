package main

import "fmt"

type AstPrinter struct{}

func (a *AstPrinter) print(expr Expr) string {
	return expr.accept(a).(string)
}

func (a *AstPrinter) visitBinaryExpr(expr *Binary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) visitGroupingExpr(expr *Grouping) any {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) visitLiteralExpr(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", *expr.Value.(*interface{}))
}

func (a *AstPrinter) visitUnaryExpr(expr *Unary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	str := "(" + name
	for _, expr := range exprs {
		str += fmt.Sprintf(" %v", expr.accept(a))
	}
	str += ")"
	return str
}
