package interpreter

import "time"

type Clock struct{}

func NewClock() Callable {
	return Clock{}
}

func (c Clock) arity() int {
	return 0
}

func (c Clock) String() string {
	return "<native fn>"
}

func (c Clock) call(interpreter *Interpreter, arguments []any) (any, error) {
	return float64(time.Now().UnixMilli()) / 1000, nil
}
