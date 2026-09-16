// Package ds holds small generic data structures: a stack, a queue and a
// priority queue.
package ds

// Stack is a last-in, first-out collection backed by a slice.
type Stack[T any] struct {
	items []T
}

// Push adds items to the top of the stack.
func (s *Stack[T]) Push(items ...T) {
	s.items = append(s.items, items...)
}

// Pop removes and returns the top item. It panics when the stack is empty;
// check [Stack.Len] first.
func (s *Stack[T]) Pop() T {
	last := len(s.items) - 1
	item := s.items[last]
	s.items = s.items[:last]

	return item
}

// Peek returns the top item without removing it.
func (s *Stack[T]) Peek() T {
	return s.items[len(s.items)-1]
}

// Len returns the number of items.
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// Empty reports whether the stack has no items.
func (s *Stack[T]) Empty() bool {
	return len(s.items) == 0
}
