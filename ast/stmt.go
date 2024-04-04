package ast

import (
	"interp/token"
)

type Stmt interface {
	Accept(stmtVisitor) (any, error)
}

type stmtVisitor interface {
	VisitExpressionStmt(*ExpressionStmt) (any, error)
	VisitFunctionStmt(*FunctionStmt) (any, error)
	VisitPrintStmt(*PrintStmt) (any, error)
	VisitReturnStmt(*ReturnStmt) (any, error)
	VisitVarStmt(*VarStmt) (any, error)
	VisitWhileStmt(*WhileStmt) (any, error)
	VisitBlockStmt(*BlockStmt) (any, error)
	VisitIfStmt(*IfStmt) (any, error)
	VisitBreakStmt(*BreakStmt) (any, error)
}

type ExpressionStmt struct {
	Expression Expr
}

func (e ExpressionStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitExpressionStmt(&e)
}

type FunctionStmt struct {
	Name   token.Token
	Params []token.Token
	Body   []Stmt
}

func (f FunctionStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitFunctionStmt(&f)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i IfStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitIfStmt(&i)
}

type PrintStmt struct {
	Expression Expr
}

func (p PrintStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitPrintStmt(&p)
}

type ReturnStmt struct {
	Keyword token.Token
	Value   Expr
}

func (r ReturnStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitReturnStmt(&r)
}

type VarStmt struct {
	Name        token.Token
	Initializer Expr
}

func (v VarStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitVarStmt(&v)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (w WhileStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitWhileStmt(&w)
}

type BlockStmt struct {
	Statements []Stmt
}

func (b BlockStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitBlockStmt(&b)
}

type BreakStmt struct {
	Keyword token.Token
}

func (b BreakStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitBreakStmt(&b)
}
