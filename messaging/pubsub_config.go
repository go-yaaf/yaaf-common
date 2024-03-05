//go:build ignore
// +build ignore

package messaging

// PubSubConfig provides configuration options for a Pub/Sub consumer.
// This configuration allows fine-tuning of the message consumption
// behavior, particularly useful for managing workloads in environments
// with slow message processing capabilities.
type PubSubConfig struct {
	// NumGoroutines specifies the number of goroutines that will be used
	// to pull messages from the subscription in parallel. Each goroutine
	// opens a separate StreamingPull stream. A higher number of goroutines
	// might increase throughput but also increases the system's load.
	// Defaults to DefaultReceiveSettings.NumGoroutines when set to 0.
	NumGoroutines int

	// MaxOutstandingMessages defines the maximum number of unprocessed
	// messages (messages that have been received but not yet acknowledged
	// or expired). Setting this to a lower number can prevent the consumer
	// from being overwhelmed by a large volume of incoming messages.
	// If set to 0, the default is DefaultReceiveSettings.MaxOutstandingMessages.
	// A negative value indicates no limit.
	MaxOutstandingMessages int

	// MaxOutstandingBytes is the maximum total size of unprocessed messages.
	// This setting helps to control memory usage by limiting the total size
	// of messages that can be held in memory at a time. If set to 0, the
	// default is DefaultReceiveSettings.MaxOutstandingBytes. A negative
	// value indicates no limit on the byte size of unprocessed messages.
	MaxOutstandingBytes int
}
