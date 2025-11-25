// In-memory implementation of a message bus (IMessageBus interface)

package messaging

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-yaaf/yaaf-common/entity"
	"github.com/go-yaaf/yaaf-common/logger"

	"github.com/go-yaaf/yaaf-common/utils/collections"
)

// InMemoryMessageBus represents an in-memory implementation of the IMessageBus interface.
// It simulates a message bus behavior within the application's memory, making it suitable for testing or simple use cases.
//
// Fields:
//
//	mu: A RWMutex to ensure thread-safe access to topics and queues.
//	topics: A map where keys are topic names and values are slices of channels. Each channel represents a subscriber.
//	queues: A map where keys are queue names and values are Queue instances for message storage.
type InMemoryMessageBus struct {
	mu     sync.RWMutex
	topics map[string][]chan []byte
	queues map[string]collections.Queue[IMessage]
}

// NewInMemoryMessageBus is a factory method that creates and returns a new instance of InMemoryMessageBus.
// It initializes the topics and queues maps.
func NewInMemoryMessageBus() (mq IMessageBus, err error) {
	return &InMemoryMessageBus{
		topics: make(map[string][]chan []byte),
		queues: make(map[string]collections.Queue[IMessage]),
	}, nil
}

// region IMessageBus methods implementation ---------------------------------------------------------------------------

// Ping tests the connectivity of the message bus. For the in-memory implementation, it always returns nil,
// indicating that the bus is always available.
//
// Parameters:
//
//	retries: The number of times to retry the ping. (Not used in this implementation).
//	intervalInSeconds: The interval in seconds between retries. (Not used in this implementation).
//
// Returns:
//
//	An error if the ping fails, otherwise nil.
func (m *InMemoryMessageBus) Ping(retries uint, intervalInSeconds uint) error {
	return nil
}

// Close releases any resources held by the message bus. For the in-memory implementation,
// it simply logs a debug message as there are no external connections to close.
//
// Returns:
//
//	An error if closing fails, otherwise nil.
func (m *InMemoryMessageBus) Close() error {
	logger.Debug("In-memory message bus closed")
	return nil
}

// CloneMessageBus returns a clone (copy) of the message bus instance.
// In this in-memory implementation, it returns a pointer to the same instance,
// as the state is shared.
//
// Returns:
//
//	A new IMessageBus instance that is a copy of the original.
//	An error if the cloning process fails.
func (m *InMemoryMessageBus) CloneMessageBus() (IMessageBus, error) {
	return m, nil
}

// Publish sends one or more messages to their respective topics.
// It iterates through the provided messages, marshals them, and sends the data to all subscriber channels for that topic.
//
// Parameters:
//
//	messages: A variadic slice of IMessage to be published.
//
// Returns:
//
//	An error if marshalling or sending fails, otherwise nil.
func (m *InMemoryMessageBus) Publish(messages ...IMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, message := range messages {
		data, err := entity.Marshal(message)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}

		if subscribers, ok := m.topics[message.Topic()]; ok {
			for _, ch := range subscribers {
				ch <- data
			}
		}
	}
	return nil
}

// Subscribe creates a subscription to one or more topics.
// It sets up a channel for the subscriber and starts a goroutine to listen for messages on that channel.
// When a message is received, it's unmarshalled and passed to the provided callback function.
//
// Parameters:
//
//	subscription: A string identifying the subscription (not used in this implementation).
//	mf: A MessageFactory function to create new message instances for unmarshalling.
//	callback: A SubscriptionCallback function to be executed when a message is received.
//	topics: A variadic slice of topic names to subscribe to.
//
// Returns:
//
//	A unique subscription ID.
//	An error if the callback is nil or another issue occurs.
func (m *InMemoryMessageBus) Subscribe(subscription string, mf MessageFactory, callback SubscriptionCallback, topics ...string) (subscriptionId string, error error) {
	if callback == nil {
		return "", fmt.Errorf("callback cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	subscriptionId = fmt.Sprintf("%d", entity.Now())
	ch := make(chan []byte, 1000)

	for _, topic := range topics {
		if _, ok := m.topics[topic]; !ok {
			m.topics[topic] = make([]chan []byte, 0)
		}
		m.topics[topic] = append(m.topics[topic], ch)
	}

	go func() {
		for data := range ch {
			message := mf()
			if err := entity.Unmarshal(data, &message); err == nil {
				callback(message)
			} else {
				logger.Error("failed to unmarshal message: %v", err)
			}
		}
	}()

	return subscriptionId, nil
}

// Unsubscribe removes a subscription. In this in-memory implementation, this is a no-op.
//
// Parameters:
//
//	subscriptionId: The ID of the subscription to remove.
//
// Returns:
//
//	A boolean indicating if the unsubscription was successful (always true).
func (m *InMemoryMessageBus) Unsubscribe(subscriptionId string) (success bool) {
	// In-memory implementation does not actively manage subscriptions by ID, so this is a no-op.
	return true
}

// Push adds one or more messages to a queue. The queue is determined by the message's topic.
// If the queue does not exist, it is created.
//
// Parameters:
//
//	messages: A variadic slice of IMessage to be pushed to the queue.
//
// Returns:
//
//	An error if any issue occurs, otherwise nil.
func (m *InMemoryMessageBus) Push(messages ...IMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, message := range messages {
		queueName := message.Topic()
		if queue, ok := m.queues[queueName]; ok {
			queue.Push(message)
		} else {
			queue = collections.NewQueue[IMessage]()
			queue.Push(message)
			m.queues[queueName] = queue
		}
	}
	return nil
}

// Pop removes and returns a message from one of the specified queues.
// It can operate in two modes:
// 1. If timeout is 0, it attempts to pop a message immediately.
// 2. If timeout is greater than 0, it blocks until a message is available or the timeout is reached.
//
// Parameters:
//
//	mf: A MessageFactory function to create message instances (not used in this implementation but required by the interface).
//	timeout: The duration to wait for a message before timing out.
//	queue: A variadic slice of queue names to pop from.
//
// Returns:
//
//	The popped IMessage.
//	An error if no message is found, the timeout is reached, or another issue occurs.
func (m *InMemoryMessageBus) Pop(mf MessageFactory, timeout time.Duration, queue ...string) (IMessage, error) {
	if timeout == 0 {
		return m.pop(queue...)
	}

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	timeoutChan := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			if message, err := m.pop(queue...); err == nil {
				return message, nil
			}
		case <-timeoutChan:
			return nil, fmt.Errorf("timeout waiting for message in queue(s): %v", queue)
		}
	}
}

// CreateProducer creates a message producer for a specific topic.
// For the in-memory bus, the bus itself acts as the producer.
//
// Parameters:
//
//	topic: The topic for which to create the producer.
//
// Returns:
//
//	An IMessageProducer instance.
//	An error if creation fails.
func (m *InMemoryMessageBus) CreateProducer(topic string) (IMessageProducer, error) {
	return m, nil
}

// CreateConsumer creates a message consumer for a specific topic.
//
// Parameters:
//
//	subscription: A string identifying the subscription (not used in this implementation).
//	mf: A MessageFactory to create message instances.
//	topics: The topics to consume from. Only the first topic is used.
//
// Returns:
//
//	An IMessageConsumer instance.
//	An error if creation fails.
func (m *InMemoryMessageBus) CreateConsumer(subscription string, mf MessageFactory, topics ...string) (IMessageConsumer, error) {
	if len(topics) == 0 {
		return nil, fmt.Errorf("at least one topic is required")
	}
	return &InMemoryMessageConsumer{
		bus:     m,
		topic:   topics[0],
		factory: mf,
	}, nil
}

// pop is an internal helper function that attempts to pop a message from the given queues without blocking.
func (m *InMemoryMessageBus) pop(queue ...string) (IMessage, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, qName := range queue {
		if q, ok := m.queues[qName]; ok {
			if msg, exists := q.Pop(); exists {
				if iMessage, ok2 := msg.(IMessage); ok2 {
					return iMessage, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("queues are empty")
}

// endregion

// region IMessageConsumer methods implementation ----------------------------------------------------------------------

// InMemoryMessageConsumer represents a consumer for the in-memory message bus.
type InMemoryMessageConsumer struct {
	bus     IMessageBus
	topic   string
	factory MessageFactory
}

// Close releases resources used by the consumer. For the in-memory implementation, this is a no-op.
//
// Returns:
//
//	An error if closing fails.
func (m *InMemoryMessageConsumer) Close() error {
	logger.Debug("In-memory message consumer closed")
	return nil
}

// Read retrieves a message from the consumer's topic.
// It blocks until a message is available or the specified timeout is reached.
//
// Parameters:
//
//	timeout: The duration to wait for a message.
//
// Returns:
//
//	The IMessage read from the topic.
//	An error if the timeout is reached or another issue occurs.
func (m *InMemoryMessageConsumer) Read(timeout time.Duration) (IMessage, error) {
	return m.bus.Pop(m.factory, timeout, m.topic)
}

// endregion
