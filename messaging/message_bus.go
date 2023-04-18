// The message bus interface for all messaging implementations
//

package messaging

import (
	"io"
	"time"
)

// IMessageBus Message bus interface
type IMessageBus interface {

	// Closer includes method Close()
	io.Closer

	// Ping Test connectivity for retries number of time with time interval (in seconds) between retries
	Ping(retries uint, intervalInSeconds uint) error

	// Publish messages to a channel (topic)
	Publish(messages ...IMessage) error

	// Subscribe on topics and return subscriberId
	Subscribe(subscription string, mf MessageFactory, callback SubscriptionCallback, topics ...string) (string, error)

	// Unsubscribe with the given subscriber id
	Unsubscribe(subscriptionId string) bool

	// Push Append one or multiple messages to a queue
	Push(messages ...IMessage) error

	// Pop Remove and get the last message in a queue or block until timeout expires
	Pop(mf MessageFactory, timeout time.Duration, queue ...string) (IMessage, error)

	// CreateProducer creates message producer for a specific topic
	CreateProducer(topic string) (IMessageProducer, error)

	// CreateConsumer creates message consumer for a specific topic
	CreateConsumer(subscription string, mf MessageFactory, topics ...string) (IMessageConsumer, error)
}

// IMessageProducer Message bus producer interface
type IMessageProducer interface {

	// Closer includes method Close()
	io.Closer

	// Publish messages to a producer channel (topic)
	Publish(messages ...IMessage) error
}

// IMessageConsumer Message bus consumer (reader) interface
type IMessageConsumer interface {
	// Closer includes method Close()
	io.Closer

	// Read message from topic, blocks until a new message arrive or until timeout expires
	// Use 0 instead of time.Duration for unlimited time
	// The standard way to use Read is by using infinite loop:
	//
	//	for {
	//		if msg, err := consumer.Read(time.Second * 5); err != nil {
	//			// Handle error
	//		} else {
	//			// Process message in a dedicated go routine
	//			go processTisMessage(msg)
	//		}
	//	}
	//
	Read(timeout time.Duration) (IMessage, error)
}
