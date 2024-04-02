package main

type Expr interface {
	accept(Visitor) any
}

type Visitor interface {
	visitBinaryExpr(*Binary) any
	visitGroupingExpr(*Grouping) any
	visitLiteralExpr(*Literal) any
	visitUnaryExpr(*Unary) any
}

type Binary struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func (b Binary) accept(visitor Visitor) any {
	return visitor.visitBinaryExpr(&b)
}

type Grouping struct {
	Expression Expr
}

func (g Grouping) accept(visitor Visitor) any {
	return visitor.visitGroupingExpr(&g)
}

type Literal struct {
	Value any
}

func (l Literal) accept(visitor Visitor) any {
	return visitor.visitLiteralExpr(&l)
}

type Unary struct {
	Operator *Token
	Right    Expr
}

func (u Unary) accept(visitor Visitor) any {
	return visitor.visitUnaryExpr(&u)
}
