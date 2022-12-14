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

	// Subscribe on topics
	Subscribe(mf MessageFactory, callback SubscriptionCallback, topics ...string) (subscriptionId string)

	// Unsubscribe with the given subscriber id
	Unsubscribe(subscriptionId string) (success bool)

	// Push Append one or multiple messages to a queue
	Push(messages ...IMessage) error

	// Pop Remove and get the last message in a queue or block until timeout expires
	Pop(mf MessageFactory, timeout time.Duration, queue ...string) (IMessage, error)
}
