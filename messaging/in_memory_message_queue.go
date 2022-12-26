// Copyright 2022. Motty Cohen
//
// In-memory implementation of a message queue (IMessageQueue interface)

package messaging

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-yaaf/yaaf-common/utils/collections"
)

// InMemoryMessageQueue represents in memory implementation of IMessageQueue interface
type InMemoryMessageQueue struct {
	queues map[string]collections.Queue
	mu     sync.RWMutex
	subs   map[string][]chan []byte
}

// NewInMemoryMessageQueue Factory method for database
func NewInMemoryMessageQueue() (mq IMessageQueue, err error) {
	return &InMemoryMessageQueue{
		queues: make(map[string]collections.Queue),
		subs:   make(map[string][]chan []byte),
	}, nil
}

// region IMessageQueue methods implementation -------------------------------------------------------------------------

// Ping Test connectivity for retries number of time with time interval (in seconds) between retries
func (m *InMemoryMessageQueue) Ping(retries uint, intervalInSeconds uint) error {
	return nil
}

// Push Append one or multiple messages to a queue
func (m *InMemoryMessageQueue) Push(messages ...IQueueMessage) error {

	// Thread safeguard
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, message := range messages {
		queueName := message.Queue()
		if queue, ok := m.queues[queueName]; ok {
			queue.Push(message)
		} else {
			queue = collections.NewQueue()
			queue.Push(message)
			m.queues[queueName] = queue
		}
	}
	return nil
}

// Pop Remove and get the last message in a queue or block until timeout expires
func (m *InMemoryMessageQueue) Pop(mf QueueMessageFactory, timeout time.Duration, queue ...string) (IQueueMessage, error) {

	if timeout == 0 {
		return m.pop(queue...)
	}

	after := time.After(timeout)
	for {
		select {
		case _ = <-time.Tick(time.Millisecond):
			if message, err := m.pop(queue...); err == nil {
				return message, nil
			}
		case <-after:
			return nil, fmt.Errorf("timeout")
		}
	}
}

// try to pop message from one of the queues
func (m *InMemoryMessageQueue) pop(queue ...string) (IQueueMessage, error) {

	// Thread safeguard
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, qName := range queue {
		if q, ok := m.queues[qName]; ok {
			if msg, exists := q.Pop(); exists {
				return msg.(IQueueMessage), nil
			}
		}
	}
	return nil, fmt.Errorf("not found")
}

// endregion
