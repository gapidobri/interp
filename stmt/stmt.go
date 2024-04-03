package stmt

import (
	"interp/expr"
	"interp/token"
)

type Stmt interface {
	Accept(visitor) (any, error)
}

type visitor interface {
	VisitExpressionStmt(*Expression) (any, error)
	VisitPrintStmt(*Print) (any, error)
	VisitVarStmt(*Var) (any, error)
	VisitBlockStmt(*Block) (any, error)
}

type Expression struct {
	Expression expr.Expr
}

func (e Expression) Accept(visitor visitor) (any, error) {
	return visitor.VisitExpressionStmt(&e)
}

type Print struct {
	Expression expr.Expr
}

func (p Print) Accept(visitor visitor) (any, error) {
	return visitor.VisitPrintStmt(&p)
}

type Var struct {
	Name        token.Token
	Initializer expr.Expr
}

func (v Var) Accept(visitor visitor) (any, error) {
	return visitor.VisitVarStmt(&v)
}

type Block struct {
	Statements []Stmt
}

func (b Block) Accept(visitor visitor) (any, error) {
	return visitor.VisitBlockStmt(&b)
}
