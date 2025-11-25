//go:build ignore
// +build ignore

package messaging

// PubSubConfig provides configuration options for a Pub/Sub consumer.
// This configuration allows for fine-tuning of message consumption behavior,
// which is particularly useful for managing workloads in environments with
// slow or resource-intensive message processing. By adjusting these settings,
// developers can control the flow of messages and prevent consumer applications
// from being overwhelmed.
type PubSubConfig struct {
	// NumGoroutines specifies the number of goroutines that will be used to pull
	// messages from the subscription in parallel. Each goroutine opens a separate
	// StreamingPull stream, which can significantly increase message throughput.
	// However, a higher number of goroutines also increases the load on the system,
	// so this value should be chosen carefully based on the available resources
	// and the desired performance characteristics.
	// If set to 0, it defaults to DefaultReceiveSettings.NumGoroutines.
	NumGoroutines int

	// MaxOutstandingMessages defines the maximum number of unprocessed messages
	// that the client will hold in memory. An unprocessed message is one that has
	// been received but not yet acknowledged (acked) or negatively acknowledged (nacked).
	// This setting helps to prevent the consumer from being overwhelmed by a large
	// volume of incoming messages.
	// If set to 0, it defaults to DefaultReceiveSettings.MaxOutstandingMessages.
	// A negative value indicates no limit.
	MaxOutstandingMessages int

	// MaxOutstandingBytes is the maximum total size of unprocessed messages that the
	// client will hold in memory. This setting is crucial for controlling memory usage,
	// especially when dealing with large messages. It limits the total size of messages
	// that can be held in memory at any given time.
	// If set to 0, it defaults to DefaultReceiveSettings.MaxOutstandingBytes.
	// A negative value indicates no limit on the byte size of unprocessed messages.
	MaxOutstandingBytes int
}
