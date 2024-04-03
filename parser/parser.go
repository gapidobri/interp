package parser

import (
	"fmt"
	"github.com/samber/lo"
	"interp/errors"
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
	loops   int
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
	case p.match(Fun):
		statement, err = p.function("function")
	case p.match(Var):
		statement, err = p.varDeclaration()
	default:
		statement, err = p.statement()
	}
	if err != nil {
		p.synchronize()
		return nil, err
	}
	return statement, nil
}

func (p *Parser) statement() (stmt.Stmt, error) {
	switch {
	case p.match(For):
		return p.forStatement()
	case p.match(If):
		return p.ifStatement()
	case p.match(Print):
		return p.printStatement()
	case p.match(Return):
		return p.returnStatement()
	case p.match(While):
		return p.whileStatement()
	case p.match(Break):
		return p.breakStatement()
	case p.match(LeftBrace):
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return stmt.Block{Statements: statements}, nil
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() (stmt.Stmt, error) {
	p.incrementLoops()
	defer p.decrementLoops()

	_, err := p.consume(LeftParen, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer stmt.Stmt
	switch {
	case p.match(Semicolon):
		initializer = nil
	case p.match(Var):
		initializer, err = p.varDeclaration()
	default:
		initializer, err = p.expressionStatement()
	}
	if err != nil {
		return nil, err
	}

	var condition expr.Expr
	if !p.check(Semicolon) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(Semicolon, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment expr.Expr
	if !p.check(RightParen) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(RightParen, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = stmt.Block{Statements: []stmt.Stmt{
			body,
			stmt.Expression{Expression: increment},
		}}
	}

	if condition == nil {
		condition = expr.Literal{Value: true}
	}
	body = stmt.While{Condition: condition, Body: body}

	if initializer != nil {
		body = stmt.Block{Statements: []stmt.Stmt{
			initializer, body,
		}}
	}

	return body, nil
}

func (p *Parser) ifStatement() (stmt.Stmt, error) {
	_, err := p.consume(LeftParen, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RightParen, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch stmt.Stmt
	if p.match(Else) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return stmt.If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
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

func (p *Parser) returnStatement() (stmt.Stmt, error) {
	keyword := p.previous()
	var value expr.Expr
	if !p.check(Semicolon) {
		var err error
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err := p.consume(Semicolon, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}

	return stmt.Return{Keyword: keyword, Value: value}, nil
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

func (p *Parser) function(kind string) (*stmt.Function, error) {
	name, err := p.consume(Identifier, fmt.Sprintf("Expect %s name.", kind))
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LeftParen, fmt.Sprintf("Expect '(' after %s name.", kind))
	if err != nil {
		return nil, err
	}

	var parameters []Token
	if !p.check(RightParen) {
		for {
			if len(parameters) >= 255 {
				_ = p.error(p.peek(), "Can't have more than 255 parameters.")
			}

			parameter, err := p.consume(Identifier, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, *parameter)

			if !p.match(Comma) {
				break
			}
		}
	}

	_, err = p.consume(RightParen, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LeftBrace, fmt.Sprintf("Expect '{' before %s body.", kind))
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &stmt.Function{Name: *name, Params: parameters, Body: body}, nil
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
	exp, err := p.or()
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

func (p *Parser) or() (expr.Expr, error) {
	exp, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(Or) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		exp = expr.Logical{Left: exp, Operator: operator, Right: right}
	}

	return exp, nil
}

func (p *Parser) and() (expr.Expr, error) {
	exp, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(And) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		exp = expr.Logical{Left: exp, Operator: operator, Right: right}
	}

	return exp, nil
}

func (p *Parser) varDeclaration() (stmt.Stmt, error) {
	name, err := p.consume(Identifier, "Expect variable name.")
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

	_, err = p.consume(Semicolon, "Expect ';' after variable declaration")
	if err != nil {
		return nil, err
	}

	return stmt.Var{Name: *name, Initializer: initializer}, nil
}

func (p *Parser) whileStatement() (stmt.Stmt, error) {
	p.incrementLoops()
	defer p.decrementLoops()

	_, err := p.consume(LeftParen, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RightParen, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return stmt.While{Condition: condition, Body: body}, nil
}

func (p *Parser) breakStatement() (stmt.Stmt, error) {
	if p.loops <= 0 {
		return nil, p.error(p.previous(), "Break outside loop.")
	}
	_, err := p.consume(Semicolon, "Expect ';' after 'break'.")
	if err != nil {
		return nil, err
	}
	return stmt.Break{}, nil
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

func (p *Parser) error(token Token, message string) errors.SyntaxError {
	err := errors.NewSyntaxError(token, message)
	fmt.Println(err)
	return err
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

	return p.call()
}

func (p *Parser) finishCall(callee expr.Expr) (expr.Expr, error) {
	var arguments []expr.Expr
	if !p.check(RightParen) {
		for {
			if len(arguments) >= 255 {
				_ = p.error(p.peek(), "Can't have more than 255 arguments.")
			}
			exp, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, exp)
			if !p.match(Comma) {
				break
			}
		}
	}

	paren, err := p.consume(RightParen, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return expr.Call{Callee: callee, Paren: *paren, Arguments: arguments}, nil
}

func (p *Parser) call() (expr.Expr, error) {
	exp, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LeftParen) {
			exp, err = p.finishCall(exp)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return exp, nil
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
	return nil, p.error(p.previous(), message)
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

func (p *Parser) incrementLoops() {
	p.loops++
}

func (p *Parser) decrementLoops() {
	p.loops--
}
