package scanner

import (
	"fmt"
	"github.com/samber/lo"
	"interp/token"
	"strconv"
)

type ScanError struct {
	line    int
	where   string
	message string
}

func (s ScanError) Error() string {
	return fmt.Sprintf("[line %d] Error%s: %s\n", s.line, s.where, s.message)
}

type Scanner struct {
	source string
	tokens []token.Token

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

func (s *Scanner) ScanTokens() ([]token.Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}

	s.tokens = append(s.tokens, token.Token{Type: token.EOF, Line: s.line})

	return s.tokens, nil
}

func (s *Scanner) scanToken() error {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(token.LeftParen)
	case ')':
		s.addToken(token.RightParen)
	case '{':
		s.addToken(token.LeftBrace)
	case '}':
		s.addToken(token.RightBrace)
	case '.':
		s.addToken(token.Dot)
	case ';':
		s.addToken(token.Semicolon)
	case ',':
		s.addToken(token.Comma)
	case '+':
		s.addToken(token.Plus)
	case '-':
		s.addToken(token.Minus)
	case '*':
		s.addToken(token.Star)
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
			s.addToken(token.Slash)
		}
	case '!':
		s.addToken(lo.Ternary(s.match('='), token.BangEqual, token.Bang))
	case '=':
		s.addToken(lo.Ternary(s.match('='), token.EqualEqual, token.Equal))
	case '<':
		s.addToken(lo.Ternary(s.match('='), token.LessEqual, token.Less))
	case '>':
		s.addToken(lo.Ternary(s.match('='), token.GreaterEqual, token.Greater))
	case '&':
		s.addToken(lo.Ternary(s.match('&'), token.AndAnd, token.And))
	case '|':
		s.addToken(lo.Ternary(s.match('|'), token.PipePipe, token.Pipe))
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++
	case '"':
		return s.string()
	default:
		switch {
		case s.isDigit(c):
			return s.number()
		case s.isAlpha(c):
			s.identifier()
		default:
			return s.error(s.line, "Unexpected character.")
		}
	}
	return nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) addToken(tokenType token.TokenType) {
	s.tokens = append(s.tokens, token.Token{
		Type:   tokenType,
		Lexeme: s.source[s.start:s.current],
		Line:   s.line,
	})
}

func (s *Scanner) addTokenLiteral(tokenType token.TokenType, literal any) {
	s.tokens = append(s.tokens, token.Token{
		Type:    tokenType,
		Literal: &literal,
		Lexeme:  s.source[s.start:s.current],
		Line:    s.line,
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

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return s.error(s.line, "Unterminated string.")
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenLiteral(token.String, value)

	return nil
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

func (s *Scanner) number() error {
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
		return s.error(s.line, "Failed to Parse float.")
	}

	s.addTokenLiteral(token.Number, value)

	return nil
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, ok := token.Keywords[text]
	if !ok {
		tokenType = token.Identifier
	}

	s.addToken(tokenType)
}

func (s *Scanner) error(line int, message string) ScanError {
	return ScanError{line: line, message: message}
}
