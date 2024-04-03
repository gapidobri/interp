package errors

import (
	"fmt"
	"interp/token"
	"strings"
)

type SyntaxError struct {
	token   token.Token
	message string
}

func NewSyntaxError(token token.Token, message string) SyntaxError {
	return SyntaxError{
		token:   token,
		message: message,
	}
}

func (p SyntaxError) Error() string {
	return fmt.Sprintf("[line %d] %s", p.token.Line, p.message)
}

func (p SyntaxError) Print(source *string) {
	line := p.token.Line - 1
	lines := strings.Split(*source, "\n")

	var newLines []string

	start := max(0, line-lineCount)
	end := min(line+1+lineCount, len(lines))
	spaces := strings.Repeat(" ", p.token.Column+2)

	for i := start; i < end; i++ {
		newLines = append(newLines, grey+fmt.Sprintf("%d", i+1)+none+" "+lines[i])
		if i == line {
			newLines = append(newLines, red+spaces+"^ Syntax error: "+p.message+none)
		}
	}

	fmt.Println(strings.Join(newLines, "\n"))
}
