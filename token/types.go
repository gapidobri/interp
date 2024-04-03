package token

type TokenType string

const (
	// Single-character token
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

	// One or two character token
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
	If     TokenType = "if"
	Else   TokenType = "else"
	True   TokenType = "true"
	False  TokenType = "false"
	Nil    TokenType = "nil"
	Class  TokenType = "class"
	Fun    TokenType = "fun"
	Var    TokenType = "var"
	For    TokenType = "for"
	While  TokenType = "while"
	Print  TokenType = "print"
	Return TokenType = "return"
	Or     TokenType = "or"
	Super  TokenType = "super"
	This   TokenType = "this"
	Break  TokenType = "break"

	EOF TokenType = "eof"
)
