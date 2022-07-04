// Copyright 2022. Motty Cohen.
//
//
package collections

import (
	"sync"
)

type ConcurrentStringMap struct {
	sync.Mutex
	m map[string]string
}

func (c *ConcurrentStringMap) Get(key string) (val string) {
	var ok bool

	c.Lock()
	defer c.Unlock()

	if val, ok = c.m[key]; ok {
	}

	return
}

func (c *ConcurrentStringMap) Put(key string, val string) {
	c.Lock()
	defer c.Unlock()
	c.m[key] = val
}

func NewConcurrentStringMap() ConcurrentStringMap {
	return ConcurrentStringMap{
		m: make(map[string]string),
	}
}
