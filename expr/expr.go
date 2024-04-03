package expr

import "interp/token"

type Expr interface {
	Accept(visitor) (any, error)
}

type visitor interface {
	VisitBinaryExpr(*Binary) (any, error)
	VisitGroupingExpr(*Grouping) (any, error)
	VisitLiteralExpr(*Literal) (any, error)
	VisitUnaryExpr(*Unary) (any, error)
	VisitVariableExpr(*Variable) (any, error)
	VisitAssignExpr(*Assign) (any, error)
	VisitLogicalExpr(*Logical) (any, error)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b Binary) Accept(visitor visitor) (any, error) {
	return visitor.VisitBinaryExpr(&b)
}

type Grouping struct {
	Expression Expr
}

func (g Grouping) Accept(visitor visitor) (any, error) {
	return visitor.VisitGroupingExpr(&g)
}

type Literal struct {
	Value any
}

func (l Literal) Accept(visitor visitor) (any, error) {
	return visitor.VisitLiteralExpr(&l)
}

type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (l Logical) Accept(visitor visitor) (any, error) {
	return visitor.VisitLogicalExpr(&l)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (u Unary) Accept(visitor visitor) (any, error) {
	return visitor.VisitUnaryExpr(&u)
}

type Variable struct {
	Name token.Token
}

func (v Variable) Accept(visitor visitor) (any, error) {
	return visitor.VisitVariableExpr(&v)
}

type Assign struct {
	Name  token.Token
	Value Expr
}

func (a Assign) Accept(visitor visitor) (any, error) {
	return visitor.VisitAssignExpr(&a)
}
