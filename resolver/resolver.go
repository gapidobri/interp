package resolver

import (
	"fmt"
	"interp/ast"
	"interp/interpreter"
	"interp/token"
)

type Resolver struct {
	interpreter *interpreter.Interpreter
	scopes      stack[map[string]bool]
}

func NewResolver(interpreter *interpreter.Interpreter) Resolver {
	return Resolver{
		interpreter: interpreter,
	}
}

func (r *Resolver) error(token token.Token, message string) {
	fmt.Printf("[line %d] %s", token.Line, message)
}

func (r *Resolver) Resolve(statements []ast.Stmt) error {
	for _, statement := range statements {
		err := r.resolveStmt(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveStmt(stmt ast.Stmt) error {
	_, err := stmt.Accept(r)
	return err
}

func (r *Resolver) resolveExpr(expr ast.Expr) error {
	_, err := expr.Accept(r)
	return err
}

func (r *Resolver) resolveFunction(function *ast.FunctionStmt) error {
	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	err := r.Resolve(function.Body)
	if err != nil {
		return err
	}
	r.endScope()
	return nil
}

func (r *Resolver) resolveLambda(lambda *ast.LambdaExpr) error {
	r.beginScope()
	for _, param := range lambda.Params {
		r.declare(param)
		r.define(param)
	}
	err := r.Resolve(lambda.Body)
	if err != nil {
		return err
	}
	r.endScope()
	return nil
}

func (r *Resolver) beginScope() {
	r.scopes.push(map[string]bool{})
}

func (r *Resolver) endScope() {
	r.scopes.pop()
}

func (r *Resolver) declare(name token.Token) {
	if r.scopes.isEmpty() {
		return
	}
	r.scopes.peek()[name.Lexeme] = false
}

func (r *Resolver) define(name token.Token) {
	if r.scopes.isEmpty() {
		return
	}
	r.scopes.peek()[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr ast.Expr, name token.Token) {
	for i := r.scopes.size() - 1; i >= 0; i-- {
		if _, ok := r.scopes.get(i)[name.Lexeme]; ok {
			r.interpreter.Resolve(expr, r.scopes.size()-1-i)
			return
		}
	}
}
