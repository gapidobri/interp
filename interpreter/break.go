package interpreter

type Break struct{}

func (b Break) Error() string {
	return "break"
}
