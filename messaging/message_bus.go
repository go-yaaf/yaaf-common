// The message bus interface for all messaging implementations
//

package messaging

import (
	"io"
	"time"
)

// IMessageBus defines the interface for a message bus, which provides a unified way to handle messaging
// through various implementations (e.g., in-memory, RabbitMQ, NATS). It supports both publish-subscribe
// and queueing patterns.
type IMessageBus interface {
	// Closer includes the Close() method for releasing resources.
	io.Closer

	// Ping tests the connectivity to the message bus.
	// It attempts to connect a specified number of times with a given interval.
	//
	// Parameters:
	//   retries: The number of times to retry connecting.
	//   intervalInSeconds: The time in seconds to wait between retries.
	//
	// Returns:
	//   An error if the connection cannot be established, otherwise nil.
	Ping(retries uint, intervalInSeconds uint) error

	// CloneMessageBus creates and returns a new instance of the message bus,
	// effectively cloning the current configuration and state.
	//
	// Returns:
	//   A new IMessageBus instance.
	//   An error if the cloning process fails.
	CloneMessageBus() (IMessageBus, error)

	// Publish sends one or more messages to a topic. This is part of the publish-subscribe pattern,
	// where messages are broadcast to all subscribers of a topic.
	//
	// Parameters:
	//   messages: A variadic slice of IMessage to be published.
	//
	// Returns:
	//   An error if publishing fails.
	Publish(messages ...IMessage) error

	// Subscribe creates a subscription to one or more topics, allowing a consumer to receive messages.
	//
	// Parameters:
	//   subscription: A name for the subscription, which can be used for durable subscriptions.
	//   mf: A MessageFactory function to create new message instances for unmarshalling.
	//   callback: The function to be called when a message is received.
	//   topics: A variadic slice of topic names to subscribe to.
	//
	// Returns:
	//   A unique string identifying the subscription.
	//   An error if the subscription fails.
	Subscribe(subscription string, mf MessageFactory, callback SubscriptionCallback, topics ...string) (string, error)

	// Unsubscribe removes a subscription, stopping the flow of messages to it.
	//
	// Parameters:
	//   subscriptionId: The ID of the subscription to be removed.
	//
	// Returns:
	//   A boolean indicating whether the unsubscription was successful.
	Unsubscribe(subscriptionId string) bool

	// Push adds one or more messages to a queue. This is part of the queueing pattern,
	// where messages are processed by one of the competing consumers.
	//
	// Parameters:
	//   messages: A variadic slice of IMessage to be added to the queue.
	//
	// Returns:
	//   An error if the push operation fails.
	Push(messages ...IMessage) error

	// Pop removes and returns a message from a queue. It can block until a message is available
	// or a timeout occurs.
	//
	// Parameters:
	//   mf: A MessageFactory function to create a new message instance for unmarshalling.
	//   timeout: The maximum time to wait for a message.
	//   queue: A variadic slice of queue names to pop from.
	//
	// Returns:
	//   The IMessage that was popped from the queue.
	//   An error if the operation times out or fails.
	Pop(mf MessageFactory, timeout time.Duration, queue ...string) (IMessage, error)

	// CreateProducer creates a message producer for a specific topic. A producer is responsible for
	// sending messages.
	//
	// Parameters:
	//   topic: The name of the topic for which to create the producer.
	//
	// Returns:
	//   An IMessageProducer instance.
	//   An error if the creation fails.
	CreateProducer(topic string) (IMessageProducer, error)

	// CreateConsumer creates a message consumer for one or more topics. A consumer is responsible for
	// receiving messages.
	//
	// Parameters:
	//   subscription: A name for the subscription, useful for durable consumers.
	//   mf: A MessageFactory function to create new message instances.
	//   topics: The topics to consume messages from.
	//
	// Returns:
	//   An IMessageConsumer instance.
	//   An error if the creation fails.
	CreateConsumer(subscription string, mf MessageFactory, topics ...string) (IMessageConsumer, error)
}

// IMessageProducer defines the interface for a message producer, which is responsible for publishing messages
// to a specific topic.
type IMessageProducer interface {
	// Closer includes the Close() method for releasing resources.
	io.Closer

	// Publish sends one or more messages to the producer's topic.
	//
	// Parameters:
	//   messages: A variadic slice of IMessage to be published.
	//
	// Returns:
	//   An error if publishing fails.
	Publish(messages ...IMessage) error
}

// IMessageConsumer defines the interface for a message consumer, which is responsible for reading messages
// from a topic.
type IMessageConsumer interface {
	// Closer includes the Close() method for releasing resources.
	io.Closer

	// Read retrieves a message from the topic, blocking until a new message arrives or a timeout is reached.
	// A timeout of 0 can be used to block indefinitely.
	//
	// A typical usage pattern is an infinite loop:
	//
	//  for {
	//      if msg, err := consumer.Read(time.Second * 5); err != nil {
	//          // Handle error, e.g., log it or break the loop
	//      } else {
	//          // Process the message, often in a separate goroutine
	//          go processThisMessage(msg)
	//      }
	//  }
	//
	// Parameters:
	//   timeout: The maximum duration to wait for a message.
	//
	Read(timeout time.Duration) (IMessage, error)
}
