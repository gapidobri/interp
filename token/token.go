package token

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal *any
	Line    int
	Column  int
}

func (t Token) GetLiteral() any {
	var literal any
	if t.Literal != nil {
		literal = *t.Literal
	}
	return literal
}

//func (t token) String() string {
//	return fmt.Sprintf("%d %s %s %v", t.Line, t.Type, t.Lexeme, t.GetLiteral())
//}
