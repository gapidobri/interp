package main

import (
	"fmt"
	"github.com/samber/lo"
	"strconv"
)

type TokenType string

const (
	// Single-character tokens
	LeftParen  TokenType = "left_paren"
	RightParen TokenType = "right_paren"
	LeftBrace  TokenType = "left_brace"
	RightBrace TokenType = "right_brace"
	Comma      TokenType = "comma"
	Semicolon  TokenType = "semicolon"
	Dot        TokenType = "dot"
	Plus       TokenType = "plus"
	Minus      TokenType = "minus"
	Star       TokenType = "star"
	Slash      TokenType = "slash"

	// One or two character tokens
	Equal        TokenType = "equal"
	BangEqual    TokenType = "not_equal"
	Bang         TokenType = "bang"
	EqualEqual   TokenType = "equal_equal"
	Greater      TokenType = "greater"
	GreaterEqual TokenType = "greater_equal"
	Less         TokenType = "less"
	LessEqual    TokenType = "less_equal"
	And          TokenType = "and"
	Pipe         TokenType = "pipe"
	AndAnd       TokenType = "and_and"
	PipePipe     TokenType = "pipe_pipe"

	// Literals
	String TokenType = "string"
	Number TokenType = "number"

	Identifier TokenType = "identifier"

	// Keywords
	If    TokenType = "if"
	Else  TokenType = "else"
	True  TokenType = "true"
	False TokenType = "false"
	Nil   TokenType = "nil"
)

var keywords = map[string]TokenType{
	"if":    If,
	"else":  Else,
	"true":  True,
	"false": False,
	"nil":   Nil,
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal *any
	Line    int
	Start   int
	End     int
}

func (t Token) GetLiteral() any {
	var literal any
	if t.Literal != nil {
		literal = *t.Literal
	}
	return literal
}

func (t Token) String() string {
	return fmt.Sprintf("%d %s %s %v", t.Line, t.Type, t.Lexeme, t.GetLiteral())
}

type Scanner struct {
	source string
	tokens []Token

	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		line:   1,
	}
}

func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LeftParen)
	case ')':
		s.addToken(RightParen)
	case '{':
		s.addToken(LeftBrace)
	case '}':
		s.addToken(RightBrace)
	case '.':
		s.addToken(Dot)
	case ';':
		s.addToken(Semicolon)
	case ',':
		s.addToken(Comma)
	case '+':
		s.addToken(Plus)
	case '-':
		s.addToken(Minus)
	case '*':
		s.addToken(Star)
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			for s.peek() != '*' && s.peekNext() != '/' && !s.isAtEnd() {
				if s.peek() == '\n' {
					s.line++
				}
				s.advance()
			}
			s.advanceN(2)
		} else {
			s.addToken(Slash)
		}
	case '!':
		s.addToken(lo.Ternary(s.match('='), BangEqual, Bang))
	case '=':
		s.addToken(lo.Ternary(s.match('='), EqualEqual, Equal))
	case '<':
		s.addToken(lo.Ternary(s.match('='), LessEqual, Less))
	case '>':
		s.addToken(lo.Ternary(s.match('='), GreaterEqual, Greater))
	case '&':
		s.addToken(lo.Ternary(s.match('&'), AndAnd, And))
	case '|':
		s.addToken(lo.Ternary(s.match('|'), PipePipe, Pipe))
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		switch {
		case s.isDigit(c):
			s.number()
		case s.isAlpha(c):
			s.identifier()
		default:
			s.error(fmt.Sprintf("Unexpected character %c.", c))
		}
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) error(err string) {
	fmt.Printf("Error at %d:%d: %s\n", s.line, s.start+1, err)
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.tokens = append(s.tokens, Token{
		Type:   tokenType,
		Lexeme: s.source[s.start:s.current],
		Line:   s.line,
		Start:  s.start,
		End:    s.current,
	})
}

func (s *Scanner) addTokenLiteral(tokenType TokenType, literal any) {
	s.tokens = append(s.tokens, Token{
		Type:    tokenType,
		Literal: &literal,
		Lexeme:  s.source[s.start:s.current],
		Line:    s.line,
		Start:   s.start,
		End:     s.current,
	})
}

func (s *Scanner) advance() uint8 {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) advanceN(n int) {
	s.current += n
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() uint8 {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() uint8 {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.error("Unterminated string.")
		return
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenLiteral(String, value)
}

func (s *Scanner) isDigit(c uint8) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) isAlpha(c uint8) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isAlphaNumeric(c uint8) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		s.error("Failed to parse float")
		return
	}

	s.addTokenLiteral(Number, value)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = Identifier
	}

	s.addToken(tokenType)
}
