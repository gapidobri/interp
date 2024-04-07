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

func NewExpressionStmt(expression Expr) *ExpressionStmt {
	return &ExpressionStmt{expression}
}

func (e *ExpressionStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitExpressionStmt(e)
}

type FunctionStmt struct {
	Name   token.Token
	Params []token.Token
	Body   []Stmt
}

func NewFunctionStmt(name token.Token, params []token.Token, body []Stmt) *FunctionStmt {
	return &FunctionStmt{name, params, body}
}

func (f *FunctionStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitFunctionStmt(f)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIfStmt(condition Expr, thenBranch Stmt, elseBranch Stmt) *IfStmt {
	return &IfStmt{condition, thenBranch, elseBranch}
}

func (i *IfStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitIfStmt(i)
}

type PrintStmt struct {
	Expression Expr
}

func NewPrintStmt(expression Expr) *PrintStmt {
	return &PrintStmt{expression}
}

func (p *PrintStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitPrintStmt(p)
}

type ReturnStmt struct {
	Keyword token.Token
	Value   Expr
}

func NewReturnStmt(keyword token.Token, value Expr) *ReturnStmt {
	return &ReturnStmt{keyword, value}
}

func (r *ReturnStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitReturnStmt(r)
}

type VarStmt struct {
	Name        token.Token
	Initializer Expr
}

func NewVarStmt(name token.Token, initializer Expr) *VarStmt {
	return &VarStmt{name, initializer}
}

func (v *VarStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitVarStmt(v)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func NewWhileStmt(condition Expr, body Stmt) *WhileStmt {
	return &WhileStmt{condition, body}
}

func (w *WhileStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitWhileStmt(w)
}

type BlockStmt struct {
	Statements []Stmt
}

func NewBlockStmt(statements []Stmt) *BlockStmt {
	return &BlockStmt{statements}
}

func (b *BlockStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitBlockStmt(b)
}

type BreakStmt struct {
	Keyword token.Token
}

func NewBreakStmt(keyword token.Token) *BreakStmt {
	return &BreakStmt{keyword}
}

func (b *BreakStmt) Accept(visitor stmtVisitor) (any, error) {
	return visitor.VisitBreakStmt(b)
}
