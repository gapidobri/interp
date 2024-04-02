package main

type Expr interface {
	accept(Visitor) (any, error)
}

type Visitor interface {
	visitBinaryExpr(*Binary) (any, error)
	visitGroupingExpr(*Grouping) (any, error)
	visitLiteralExpr(*Literal) (any, error)
	visitUnaryExpr(*Unary) (any, error)
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b Binary) accept(visitor Visitor) (any, error) {
	return visitor.visitBinaryExpr(&b)
}

type Grouping struct {
	Expression Expr
}

func (g Grouping) accept(visitor Visitor) (any, error) {
	return visitor.visitGroupingExpr(&g)
}

type Literal struct {
	Value any
}

func (l Literal) accept(visitor Visitor) (any, error) {
	return visitor.visitLiteralExpr(&l)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u Unary) accept(visitor Visitor) (any, error) {
	return visitor.visitUnaryExpr(&u)
}
