// Copyright 2022. Motty Cohen
//
// In-memory implementation of a message bus (IMessageBus interface)

package messaging

import (
	"encoding/json"
	"fmt"
	"github.com/go-yaaf/yaaf-common/entity"
	"sync"
	"time"

	"github.com/go-yaaf/yaaf-common/utils/collections"
)

// InMemoryMessageBus represents in memory implementation of IMessageBus interface
// topics is a map ot topic -> array of channels (channel per subscriber)
type InMemoryMessageBus struct {
	mu     sync.RWMutex
	topics map[string][]chan []byte
	queues map[string]collections.Queue
}

// NewInMemoryMessageBus Factory method
func NewInMemoryMessageBus() (mq IMessageBus, err error) {
	return &InMemoryMessageBus{
		topics: make(map[string][]chan []byte),
		queues: make(map[string]collections.Queue),
	}, nil
}

// region IMessageQueue methods implementation -------------------------------------------------------------------------

// Ping Test connectivity for retries number of time with time interval (in seconds) between retries
func (m *InMemoryMessageBus) Ping(retries uint, intervalInSeconds uint) error {
	return nil
}

// Publish messages to a channel (topic)
func (m *InMemoryMessageBus) Publish(messages ...IMessage) error {
	// Thread safeguard
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, message := range messages {
		data, err := json.Marshal(message)
		if err != nil {
			return err
		}

		for _, ch := range m.topics[message.Topic()] {
			ch <- data
		}
	}

	return nil
}

// Subscribe on topics
func (m *InMemoryMessageBus) Subscribe(callback SubscriptionCallback, mf MessageFactory, topics ...string) (subscriptionId string) {

	// Validate callback
	if callback == nil {
		return ""
	}

	// Thread safeguard
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create and register channel
	subscriptionId = fmt.Sprintf("%d", entity.Now())
	cn := make(chan []byte, 1000)

	for _, topic := range topics {
		if _, ok := m.topics[topic]; !ok {
			m.topics[topic] = make([]chan []byte, 0)
		}
		m.topics[topic] = append(m.topics[topic], cn)
	}

	go func() {
		for {
			select {
			case data := <-cn:
				message := mf()
				if err := json.Unmarshal(data, &message); err == nil {
					callback(message)
				}
			}
		}
	}()

	return subscriptionId
}

// Unsubscribe with the given subscriber id
func (m *InMemoryMessageBus) Unsubscribe(subscriptionId string) (success bool) {
	// Tdo nothing
	return true
}

// Push Append one or multiple messages to a queue
func (m *InMemoryMessageBus) Push(messages ...IMessage) error {

	// Thread safeguard
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, message := range messages {
		queueName := message.Topic()
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
func (m *InMemoryMessageBus) Pop(mf MessageFactory, timeout time.Duration, queue ...string) (IMessage, error) {

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
func (m *InMemoryMessageBus) pop(queue ...string) (IMessage, error) {

	// Thread safeguard
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, qName := range queue {
		if q, ok := m.queues[qName]; ok {
			if msg, exists := q.Pop(); exists {
				return msg.(IMessage), nil
			}
		}
	}
	return nil, fmt.Errorf("not found")
}

// endregion
