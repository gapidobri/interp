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

func NewBinaryExpr(left Expr, operator token.Token, right Expr) *BinaryExpr {
	return &BinaryExpr{left, operator, right}
}

func (b *BinaryExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitBinaryExpr(b)
}

type CallExpr struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

func NewCallExpr(callee Expr, paren token.Token, arguments []Expr) *CallExpr {
	return &CallExpr{callee, paren, arguments}
}

func (c *CallExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitCallExpr(c)
}

type GroupingExpr struct {
	Expression Expr
}

func NewGroupingExpr(expression Expr) *GroupingExpr {
	return &GroupingExpr{expression}
}

func (g *GroupingExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitGroupingExpr(g)
}

type LambdaExpr struct {
	Params []token.Token
	Body   []Stmt
}

func NewLambdaExpr(params []token.Token, body []Stmt) *LambdaExpr {
	return &LambdaExpr{params, body}
}

func (l *LambdaExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitLambdaExpr(l)
}

type LiteralExpr struct {
	Value any
}

func NewLiteralExpr(value any) *LiteralExpr {
	return &LiteralExpr{value}
}

func (l *LiteralExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitLiteralExpr(l)
}

type LogicalExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func NewLogicalExpr(left Expr, operator token.Token, right Expr) *LogicalExpr {
	return &LogicalExpr{left, operator, right}
}

func (l *LogicalExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitLogicalExpr(l)
}

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

func NewUnaryExpr(operator token.Token, right Expr) *UnaryExpr {
	return &UnaryExpr{operator, right}
}

func (u *UnaryExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitUnaryExpr(u)
}

type VariableExpr struct {
	Name token.Token
}

func NewVariableExpr(name token.Token) *VariableExpr {
	return &VariableExpr{name}
}

func (v *VariableExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitVariableExpr(v)
}

type AssignExpr struct {
	Name  token.Token
	Value Expr
}

func NewAssignExpr(name token.Token, value Expr) *AssignExpr {
	return &AssignExpr{name, value}
}

func (a *AssignExpr) Accept(visitor exprVisitor) (any, error) {
	return visitor.VisitAssignExpr(a)
}
