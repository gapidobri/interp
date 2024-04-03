package main

//import (
//	"fmt"
//	"github.com/samber/lo"
//	"interp/expr"
//)
//
//type RPNPrinter struct{}
//
//func (a *RPNPrinter) print(expr expr.Expr) (*string, error) {
//	value, err := expr.Accept(a)
//	if err != nil {
//		return nil, err
//	}
//	return lo.ToPtr(value.(string)), nil
//}
//
//func (a *RPNPrinter) VisitBinaryExpr(expr *expr.Binary) (any, error) {
//	left, err := expr.Left.Accept(a)
//	if err != nil {
//		return nil, err
//	}
//
//	right, err := expr.Right.Accept(a)
//	if err != nil {
//		return nil, err
//	}
//
//	return left.(string) + " " +
//		right.(string) + " " +
//		expr.Operator.Lexeme, nil
//}
//
//func (a *RPNPrinter) VisitGroupingExpr(expr *expr.Grouping) (any, error) {
//	return expr.Expression.Accept(a)
//}
//
//func (a *RPNPrinter) VisitLiteralExpr(expr *expr.Literal) (any, error) {
//	if expr.Value == nil {
//		return "nil", nil
//	}
//	return fmt.Sprintf("%v", expr.Value), nil
//}
//
//func (a *RPNPrinter) VisitUnaryExpr(expr *expr.Unary) (any, error) {
//	return "", nil
//}
