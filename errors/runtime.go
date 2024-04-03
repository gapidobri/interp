package errors

import (
	"fmt"
	"interp/token"
	"strings"
)

type RuntimeError struct {
	token   token.Token
	message string
}

func NewRuntimeError(token token.Token, message string) RuntimeError {
	return RuntimeError{token: token, message: message}
}

func (r RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] %s", r.token.Line, r.message)
}

func (r RuntimeError) Print(source *string) {
	line := r.token.Line - 1
	lines := strings.Split(*source, "\n")

	var newLines []string

	start := max(0, line-lineCount)
	end := min(line+1+lineCount, len(lines))
	spaces := strings.Repeat(" ", r.token.Column+2)

	for i := start; i < end; i++ {
		newLines = append(newLines, grey+fmt.Sprintf("%d", i+1)+none+" "+lines[i])
		if i == line {
			newLines = append(newLines, red+spaces+"^ Runtime error: "+r.message+none)
		}
	}

	fmt.Println(strings.Join(newLines, "\n"))
}
