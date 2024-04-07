package resolver

import (
	"fmt"
	"interp/ast"
	"interp/errors"
	"interp/interpreter"
	"interp/token"
)

type Resolver struct {
	interpreter     *interpreter.Interpreter
	scopes          stack[map[string]*varState]
	currentFunction FunctionType
	currentLoop     LoopType
}

func NewResolver(interpreter *interpreter.Interpreter) Resolver {
	return Resolver{
		interpreter:     interpreter,
		currentFunction: FunctionTypeNone,
		currentLoop:     LoopTypeNone,
		scopes:          stack[map[string]*varState]{},
	}
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

func (r *Resolver) resolveFunction(function *ast.FunctionStmt, funcType FunctionType) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = funcType

	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	// TODO: Fix for nested returns
	for i, statement := range function.Body {
		if returnStmt, ok := statement.(*ast.ReturnStmt); ok && i != len(function.Body)-1 {
			errors.Warning(returnStmt.Keyword, "Unreachable code after return.")
		}
	}
	err := r.Resolve(function.Body)
	if err != nil {
		return err
	}

	r.endScope()
	r.currentFunction = enclosingFunction
	return nil
}

func (r *Resolver) resolveLambda(lambda *ast.LambdaExpr) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = FunctionTypeLambda

	r.beginScope()
	for _, param := range lambda.Params {
		r.declare(param)
		r.define(param)
	}
	// TODO: Fix for nested returns
	for i, statement := range lambda.Body {
		if returnStmt, ok := statement.(*ast.ReturnStmt); ok && i != len(lambda.Body)-1 {
			errors.Warning(returnStmt.Keyword, "Unreachable code after return.")
		}
	}
	err := r.Resolve(lambda.Body)
	if err != nil {
		return err
	}
	r.endScope()
	r.currentFunction = enclosingFunction
	return nil
}

func (r *Resolver) beginScope() {
	r.scopes.push(map[string]*varState{})
}

func (r *Resolver) endScope() {
	scope, ok := r.scopes.pop()
	if !ok {
		return

	}
	for name, state := range scope {
		if !state.resolved {
			errors.Warning(state.token, fmt.Sprintf("Variable '%s' is declared but never used.", name))
		}
	}
}

func (r *Resolver) declare(name token.Token) {
	if r.scopes.isEmpty() {
		return
	}
	scope, ok := r.scopes.peek()
	if !ok {
		return
	}
	if _, ok = scope[name.Lexeme]; ok {
		errors.Error(name, "Already a variable with this name in this scope.")
	}
	scope[name.Lexeme] = &varState{token: name}
}

func (r *Resolver) define(name token.Token) {
	if r.scopes.isEmpty() {
		return
	}
	scope, ok := r.scopes.peek()
	if !ok {
		return
	}
	scope[name.Lexeme].define()
}

func (r *Resolver) resolveLocal(expr ast.Expr, name token.Token) {
	scope, ok := r.scopes.peek()
	if !ok {
		return
	}
	scope[name.Lexeme].resolve()
	for i := r.scopes.size() - 1; i >= 0; i-- {
		if _, ok = r.scopes.get(i)[name.Lexeme]; ok {
			r.interpreter.Resolve(expr, r.scopes.size()-1-i)
			return
		}
	}
}
