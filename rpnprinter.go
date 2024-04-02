package main

import "fmt"

type RPNPrinter struct{}

func (a *RPNPrinter) print(expr Expr) string {
	return expr.accept(a).(string)
}

func (a *RPNPrinter) visitBinaryExpr(expr *Binary) any {
	return expr.Left.accept(a).(string) + " " +
		expr.Right.accept(a).(string) + " " +
		expr.Operator.Lexeme
}

func (a *RPNPrinter) visitGroupingExpr(expr *Grouping) any {
	return expr.Expression.accept(a)
}

func (a *RPNPrinter) visitLiteralExpr(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (a *RPNPrinter) visitUnaryExpr(expr *Unary) any {
	return ""
}
