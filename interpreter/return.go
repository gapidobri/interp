package interpreter

import "fmt"

type Return struct {
	Value any
}

func (r Return) Error() string {
	return fmt.Sprintf("return %v", r.Value)
}
