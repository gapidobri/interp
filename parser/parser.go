package parser

import (
	"fmt"
	"github.com/samber/lo"
	"interp/expr"
	"interp/stmt"
	. "interp/token"
)

type ParseError struct {
	Token   Token
	Message string
}

func (p ParseError) Error() string {
	return fmt.Sprintf("Parse error at line %d: %s", p.Token.Line, p.Message)
}

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) Parser {
	return Parser{tokens: tokens}
}

func (p *Parser) Parse() ([]stmt.Stmt, error) {
	var statements []stmt.Stmt
	for !p.isAtEnd() {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}

		statements = append(statements, declaration)
	}
	return statements, nil
}

func (p *Parser) expression() (expr.Expr, error) {
	return p.assigment()
}

func (p *Parser) declaration() (stmt.Stmt, error) {
	var (
		statement stmt.Stmt
		err       error
	)
	switch {
	case p.match(Var):
		statement, err = p.varDeclaration()
	default:
		statement, err = p.statement()
	}
	if err != nil {
		p.synchronize()
		return nil, nil
	}
	return statement, nil
}

func (p *Parser) statement() (stmt.Stmt, error) {
	switch {
	case p.match(Print):
		return p.printStatement()
	case p.match(LeftBrace):
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return stmt.Block{Statements: statements}, nil
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() (stmt.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(Semicolon, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return stmt.Print{Expression: value}, nil
}

func (p *Parser) expressionStatement() (stmt.Stmt, error) {
	exp, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(Semicolon, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}

	return stmt.Expression{Expression: exp}, nil
}

func (p *Parser) block() ([]stmt.Stmt, error) {
	var statements []stmt.Stmt

	for !p.check(RightBrace) && !p.isAtEnd() {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, declaration)
	}

	_, err := p.consume(RightBrace, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) assigment() (expr.Expr, error) {
	exp, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(Equal) {
		equals := p.previous()
		value, err := p.assigment()
		if err != nil {
			return nil, err
		}

		if variable, ok := exp.(expr.Variable); ok {
			return expr.Assign{Name: variable.Name, Value: value}, nil
		}

		return nil, p.error(equals, "Invalid assigment target.")
	}

	return exp, nil
}

func (p *Parser) varDeclaration() (stmt.Stmt, error) {
	name, err := p.consume(Identifier, "Expec variable name.")
	if err != nil {
		return nil, err
	}

	var initializer expr.Expr
	if p.match(Equal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(Semicolon, "Expect  ';' after variable declaration")
	if err != nil {
		return nil, err
	}

	return stmt.Var{Name: *name, Initializer: initializer}, nil
}

func (p *Parser) equality() (expr.Expr, error) {
	exp, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		exp = expr.Binary{Left: exp, Operator: operator, Right: right}
	}
	return exp, nil
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) error(token Token, message string) ParseError {
	return ParseError{token, message}
}

func (p *Parser) comparison() (expr.Expr, error) {
	exp, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		exp = expr.Binary{Left: exp, Operator: operator, Right: right}
	}
	return exp, nil
}

func (p *Parser) term() (expr.Expr, error) {
	exp, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(Minus, Plus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		exp = expr.Binary{Left: exp, Operator: operator, Right: right}
	}

	return exp, nil
}

func (p *Parser) factor() (expr.Expr, error) {
	exp, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(Slash, Star) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		exp = expr.Binary{Left: exp, Operator: operator, Right: right}
	}

	return exp, nil
}

func (p *Parser) unary() (expr.Expr, error) {
	if p.match(Bang, Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return expr.Unary{Operator: operator, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (expr.Expr, error) {
	switch {
	case p.match(False):
		return expr.Literal{Value: false}, nil
	case p.match(True):
		return expr.Literal{Value: true}, nil
	case p.match(Nil):
		return expr.Literal{}, nil
	case p.match(Number, String):
		return expr.Literal{Value: *p.previous().Literal}, nil
	case p.match(Identifier):
		return expr.Variable{Name: p.previous()}, nil
	case p.match(LeftParen):
		exp, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(RightParen, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return expr.Grouping{Expression: exp}, nil
	}

	return nil, p.error(p.peek(), "Expect expression.")
}

func (p *Parser) consume(t TokenType, message string) (*Token, error) {
	if p.check(t) {
		return lo.ToPtr(p.advance()), nil
	}
	return nil, p.error(p.peek(), message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == Semicolon {
			return
		}
		switch p.peek().Type {
		case Class:
			fallthrough
		case Fun:
			fallthrough
		case Var:
			fallthrough
		case For:
			fallthrough
		case If:
			fallthrough
		case While:
			fallthrough
		case Print:
			fallthrough
		case Return:
			return
		}

		p.advance()
	}
}
