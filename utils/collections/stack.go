// Package collections provides a thread-safe implementation of a LIFO (Last-In, First-Out) stack data structure.
package collections

import (
	"sync"
)

// Stack defines the interface for a generic, thread-safe stack.
//
// Type Parameters:
//
//	T: The type of the items stored in the stack.
type Stack[T any] interface {
	// Push adds an item to the top of the stack.
	Push(v T)

	// Pop removes and returns the item from the top of the stack.
	// It also returns a boolean indicating if an item was successfully popped.
	Pop() (T, bool)

	// PopMany removes and returns a specified number of items from the top of the stack.
	// If the stack contains fewer items than requested, it returns all the items in the stack.
	PopMany(count int) ([]T, bool)

	// PopAll removes and returns all items from the stack.
	PopAll() ([]T, bool)

	// Peek returns the item at the top of the stack without removing it.
	Peek() (T, bool)

	// Length returns the number of items in the stack.
	Length() int

	// IsEmpty checks if the stack is empty.
	IsEmpty() bool
}

// defaultStack is the default implementation of the Stack interface.
// It uses a slice and a mutex to provide a thread-safe stack.
type defaultStack[T any] struct {
	sync.Mutex
	stack []T
}

// NewStack creates and returns a new instance of a generic, thread-safe stack.
//
// Type Parameters:
//
//	T: The type of the items to be stored in the stack.
//
// Returns:
//
//	A new Stack instance.
func NewStack[T any]() Stack[T] {
	return &defaultStack[T]{
		stack: make([]T, 0),
	}
}

// Push adds an item to the top of the stack in a thread-safe manner.
func (s *defaultStack[T]) Push(v T) {
	s.Lock()
	defer s.Unlock()
	s.stack = append(s.stack, v)
}

// Pop removes and returns the item from the top of the stack in a thread-safe manner.
// If the stack is empty, it returns the zero value for the type and false.
func (s *defaultStack[T]) Pop() (T, bool) {
	s.Lock()
	defer s.Unlock()

	if len(s.stack) == 0 {
		var zero T
		return zero, false
	}

	index := len(s.stack) - 1
	v := s.stack[index]
	s.stack = s.stack[:index]
	return v, true
}

// PopMany removes and returns a specified number of items from the top of the stack.
func (s *defaultStack[T]) PopMany(count int) ([]T, bool) {
	s.Lock()
	defer s.Unlock()

	if len(s.stack) == 0 {
		return nil, false
	}

	if count > len(s.stack) {
		count = len(s.stack)
	}

	index := len(s.stack) - count
	vs := s.stack[index:]
	s.stack = s.stack[:index]
	return vs, true
}

// PopAll removes and returns all items from the stack.
func (s *defaultStack[T]) PopAll() ([]T, bool) {
	s.Lock()
	defer s.Unlock()

	if len(s.stack) == 0 {
		return nil, false
	}

	all := s.stack
	s.stack = make([]T, 0)
	return all, true
}

// Peek returns the item at the top of the stack without removing it.
func (s *defaultStack[T]) Peek() (T, bool) {
	s.Lock()
	defer s.Unlock()

	if len(s.stack) == 0 {
		var zero T
		return zero, false
	}

	return s.stack[len(s.stack)-1], true
}

// Length returns the number of items in the stack.
func (s *defaultStack[T]) Length() int {
	s.Lock()
	defer s.Unlock()
	return len(s.stack)
}

// IsEmpty checks if the stack is empty.
func (s *defaultStack[T]) IsEmpty() bool {
	return s.Length() == 0
}
