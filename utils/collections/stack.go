// Copyright 2022. Motty Cohen.
//
//  Thread-safe implementation of stack data structure
//

package collections

import (
	"sync"
	"sync/atomic"
)

// Stack functions for manager data items in a stack
type Stack interface {
	// push item into a stack
	Push(v interface{})

	// pop last item
	Pop() (interface{}, bool)

	// pop many items
	PopMany(count int64) ([]interface{}, bool)

	// pop all items
	PopAll() ([]interface{}, bool)

	// peek last item
	Peek() (interface{}, bool)

	// get length of stack
	Length() int64

	// is empty stack
	IsEmpty() bool
}

type defaultStack struct {
	sync.Mutex
	length int64
	stack  []interface{}
}

// New get stack functions manager
func NewStack() Stack {
	return &defaultStack{}
}

func (p *defaultStack) Push(v interface{}) {
	p.Lock()
	defer p.Unlock()

	prepend := make([]interface{}, 1)
	prepend[0] = v

	p.stack = append(prepend, p.stack...)
	p.length++
}

func (p *defaultStack) Pop() (v interface{}, exist bool) {
	if p.IsEmpty() {
		return
	}

	p.Lock()
	defer p.Unlock()

	v, p.stack, exist = p.stack[0], p.stack[1:], true
	p.length--

	return
}

func (p *defaultStack) PopMany(count int64) (vs []interface{}, exist bool) {

	if p.IsEmpty() {
		return
	}

	p.Lock()
	defer p.Unlock()

	if count >= p.length {
		count = p.length
	}
	p.length -= count

	vs, p.stack, exist = p.stack[:count-1], p.stack[count:], true
	return
}

func (p *defaultStack) PopAll() (all []interface{}, exist bool) {
	if p.IsEmpty() {
		return
	}
	p.Lock()
	defer p.Unlock()

	all, p.stack, exist = p.stack[:], nil, true
	p.length = 0
	return
}

func (p *defaultStack) Peek() (v interface{}, exist bool) {
	if p.IsEmpty() {
		return
	}

	p.Lock()
	defer p.Unlock()

	return p.stack[0], true
}

func (p *defaultStack) Length() int64 {
	return atomic.LoadInt64(&p.length)
}

func (p *defaultStack) IsEmpty() bool {
	return p.Length() == 0
}
