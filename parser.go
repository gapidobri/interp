package main

import (
	"fmt"
	"github.com/samber/lo"
)

type ParseError struct {
	Token   Token
	Message string
}

func newParseError(token Token, message string) ParseError {
	return ParseError{token, message}
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

func (p *Parser) Parse() (Expr, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = &Binary{expr, operator, right}
	}
	return expr, nil
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

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = Binary{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(Minus, Plus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = Binary{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(Slash, Star) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = Binary{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(Bang, Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return Unary{operator, right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	switch {
	case p.match(False):
		return Literal{false}, nil
	case p.match(True):
		return Literal{true}, nil
	case p.match(Nil):
		return Literal{nil}, nil
	case p.match(Number, String):
		return Literal{*p.previous().Literal}, nil
	case p.match(LeftParen):
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(RightParen, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return Grouping{expr}, nil
	}

	return nil, p.error(p.peek(), "Expect expression.")
}

func (p *Parser) consume(t TokenType, message string) (*Token, error) {
	if p.check(t) {
		return lo.ToPtr(p.advance()), nil
	}
	return nil, newParseError(p.peek(), message)
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
