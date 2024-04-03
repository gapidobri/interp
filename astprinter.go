package main

//
//import (
//	"fmt"
//	"interp/expr"
//)
//
//type AstPrinter struct{}
//
//func (a *AstPrinter) print(exp expr.Expr) (*string, error) {
//	value, err := exp.Accept(a)
//	if err != nil {
//		return nil, err
//	}
//	return value.(*string), nil
//}
//
//func (a *AstPrinter) VisitBinaryExpr(expr *expr.Binary) (any, error) {
//	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
//}
//
//func (a *AstPrinter) VisitGroupingExpr(expr *expr.Grouping) (any, error) {
//	return a.parenthesize("group", expr.Expression)
//}
//
//func (a *AstPrinter) VisitLiteralExpr(expr *expr.Literal) (any, error) {
//	if expr.Value == nil {
//		return "nil", nil
//	}
//	return fmt.Sprintf("%v", expr.Value.(interface{})), nil
//}
//
//func (a *AstPrinter) VisitUnaryExpr(expr *expr.Unary) (any, error) {
//	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
//}
//
//func (a *AstPrinter) parenthesize(name string, exprs ...expr.Expr) (*string, error) {
//	str := "(" + name
//	for _, exp := range exprs {
//		value, err := exp.Accept(a)
//		if err != nil {
//			return nil, err
//		}
//		str += fmt.Sprintf(" %v", value)
//	}
//	str += ")"
//	return &str, nil
//}
