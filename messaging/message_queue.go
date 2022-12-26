// Copyright 2022. Shield-IoT Ltd.
//
// The message queue interface for all messaging implementations
//

package messaging

import "time"

// IMessageQueue Message bus interface
type IMessageQueue interface {

	// Ping Test connectivity for retries number of time with time interval (in seconds) between retries
	Ping(retries uint, intervalInSeconds uint) error

	// Push Append one or multiple messages to a queue
	Push(messages ...IQueueMessage) error

	// Pop Remove and get the last message in a queue or block until timeout expires
	Pop(mf QueueMessageFactory, timeout time.Duration, queue ...string) (IQueueMessage, error)
}
