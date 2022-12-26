// Copyright 2022. Motty Cohen.
//
// Thread-safe implementation of FIFO queue data structure
//

package collections

import (
	"sync"
)

// Queue functions for manager data items in a stack
type Queue interface {
	// Push item into a queue
	Push(v any)

	// Pop last item
	Pop() (any, bool)

	// Length get length of the queue
	Length() int
}

type defaultQueue struct {
	sync.Mutex
	queue []any
}

// NewQueue get queue functions manager
func NewQueue() Queue {
	return &defaultQueue{
		queue: make([]any, 0),
	}
}

// Push item to queue
func (p *defaultQueue) Push(v any) {
	p.Lock()
	defer p.Unlock()
	p.queue = append(p.queue, v)
}

// Pop item from queue
func (p *defaultQueue) Pop() (v any, exist bool) {
	if p.Length() == 0 {
		return
	}

	p.Lock()
	defer p.Unlock()

	v, p.queue, exist = p.queue[0], p.queue[1:], true
	return
}

// Length get queue length (number of items)
func (p *defaultQueue) Length() int {
	p.Lock()
	defer p.Unlock()
	return len(p.queue)
}
