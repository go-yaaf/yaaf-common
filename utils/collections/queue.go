// Package collections provides a thread-safe implementation of a FIFO (First-In, First-Out) queue data structure.
package collections

import (
	"sync"
)

// Queue defines the interface for a generic, thread-safe queue.
//
// Type Parameters:
//
//	T: The type of the items stored in the queue.
type Queue[T any] interface {
	// Push adds an item to the end of the queue.
	Push(v T)

	// Pop removes and returns the item from the front of the queue.
	// It also returns a boolean indicating if an item was successfully popped.
	Pop() (T, bool)

	// Length returns the number of items in the queue.
	Length() int
}

// defaultQueue is the default implementation of the Queue interface.
// It uses a slice and a mutex to provide a thread-safe queue.
type defaultQueue[T any] struct {
	sync.Mutex
	queue []T
}

// NewQueue creates and returns a new instance of a generic, thread-safe queue.
//
// Type Parameters:
//
//	T: The type of the items to be stored in the queue.
//
// Returns:
//
//	A new Queue instance.
func NewQueue[T any]() Queue[T] {
	return &defaultQueue[T]{
		queue: make([]T, 0),
	}
}

// Push adds an item to the end of the queue in a thread-safe manner.
func (q *defaultQueue[T]) Push(v T) {
	q.Lock()
	defer q.Unlock()
	q.queue = append(q.queue, v)
}

// Pop removes and returns the item from the front of the queue in a thread-safe manner.
// If the queue is empty, it returns the zero value for the type and false.
func (q *defaultQueue[T]) Pop() (T, bool) {
	q.Lock()
	defer q.Unlock()

	if len(q.queue) == 0 {
		var zero T
		return zero, false
	}

	v := q.queue[0]
	q.queue = q.queue[1:]
	return v, true
}

// Length returns the number of items in the queue in a thread-safe manner.
func (q *defaultQueue[T]) Length() int {
	q.Lock()
	defer q.Unlock()
	return len(q.queue)
}
