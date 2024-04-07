package resolver

type stack[T any] struct {
	stack []T
}

func (s *stack[T]) push(value T) {
	s.stack = append(s.stack, value)
}

func (s *stack[T]) peek() (T, bool) {
	if len(s.stack) == 0 {
		var empty T
		return empty, false
	}
	return s.stack[len(s.stack)-1], true
}

func (s *stack[T]) pop() (T, bool) {
	if len(s.stack) == 0 {
		var empty T
		return empty, false
	}
	value := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return value, true
}

func (s *stack[T]) get(i int) T {
	return s.stack[i]
}

func (s *stack[T]) isEmpty() bool {
	return len(s.stack) == 0
}

func (s *stack[T]) size() int {
	return len(s.stack)
}
