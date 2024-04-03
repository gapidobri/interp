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
	VisitWhileStmt(*While) (any, error)
	VisitBlockStmt(*Block) (any, error)
	VisitIfStmt(*If) (any, error)
}

type Expression struct {
	Expression expr.Expr
}

func (e Expression) Accept(visitor visitor) (any, error) {
	return visitor.VisitExpressionStmt(&e)
}

type If struct {
	Condition  expr.Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i If) Accept(visitor visitor) (any, error) {
	return visitor.VisitIfStmt(&i)
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

type While struct {
	Condition expr.Expr
	Body      Stmt
}

func (w While) Accept(visitor visitor) (any, error) {
	return visitor.VisitWhileStmt(&w)
}

type Block struct {
	Statements []Stmt
}

func (b Block) Accept(visitor visitor) (any, error) {
	return visitor.VisitBlockStmt(&b)
}
