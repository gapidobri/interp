package ast

import (
	"interp/token"
	_ "unsafe"
)

type Expr interface {
	Accept(exprVisitor) (any, error)
}

type exprVisitor interface {
	VisitBinaryExpr(*BinaryExpr) (any, error)
	VisitCallExpr(*CallExpr) (any, error)
	VisitGroupingExpr(*GroupingExpr) (any, error)
	VisitLambdaExpr(*LambdaExpr) (any, error)
	VisitLiteralExpr(*LiteralExpr) (any, error)
	VisitUnaryExpr(*UnaryExpr) (any, error)
	VisitVariableExpr(*VariableExpr) (any, error)
	VisitAssignExpr(*AssignExpr) (any, error)
	VisitLogicalExpr(*LogicalExpr) (any, error)
}

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b BinaryExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitBinaryExpr(&b)
}

type CallExpr struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

func (c CallExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitCallExpr(&c)
}

type GroupingExpr struct {
	Expression Expr
}

func (g GroupingExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitGroupingExpr(&g)
}

type LambdaExpr struct {
	Params []token.Token
	Body   []Stmt
}

func (l LambdaExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitLambdaExpr(&l)
}

type LiteralExpr struct {
	Value any
}

func (l LiteralExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitLiteralExpr(&l)
}

type LogicalExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (l LogicalExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitLogicalExpr(&l)
}

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

func (u UnaryExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitUnaryExpr(&u)
}

type VariableExpr struct {
	Name token.Token
}

func (v VariableExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitVariableExpr(&v)
}

type AssignExpr struct {
	Name  token.Token
	Value Expr
}

func (a AssignExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitAssignExpr(&a)
}
