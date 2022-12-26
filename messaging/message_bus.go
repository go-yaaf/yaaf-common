// Copyright 2022. Shield-IoT Ltd.
//
// The message bus interface for all messaging implementations
//

package messaging

// IMessageBus Message bus interface
type IMessageBus interface {

	// Ping Test connectivity for retries number of time with time interval (in seconds) between retries
	Ping(retries uint, intervalInSeconds uint) error

	// Publish messages to a channel (topic)
	Publish(messages ...PubSubMessage) error

	// Subscribe on topics
	Subscribe(callback SubscriptionCallback, mf PubSubMessageFactory, topics ...string) (subscriptionId string)

	// Unsubscribe with the given subscriber id
	Unsubscribe(subscriptionId string) (success bool)
}
