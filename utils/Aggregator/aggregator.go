/*
 * Aggregator is a mechanism to aggregate items to bulks and notify when the bulk reached a certain size or after time out
 */

package aggregator

import (
	"sync"
	"time"
)

// bulkCallback is a function type for handling a full bulk of items.
// It is called when the number of items in the aggregator reaches the defined bulk size.
type bulkCallback[T any] func(bulk []T)

// timeoutCallback is a function type for handling a bulk of items that has timed out.
// It is called when the timeout is triggered and there are items in the aggregator that have not yet formed a full bulk.
type timeoutCallback[T any] func(bulk []T)

// Aggregator provides a mechanism to collect items into bulks of a specified size.
// A bulk is processed either when it's full or when a timeout occurs.
// It is safe for concurrent use.
//
// Type Parameters:
//
//	T: The type of the items to be aggregated.
type Aggregator[T any] struct {
	mutex           sync.Mutex           // Protects access to the items slice and isShutDown flag.
	timeout         time.Duration        // The duration after which a non-full bulk is processed.
	bulkSize        int                  // The target size of a bulk.
	items           []T                  // The collection of items waiting to be processed.
	bulkCallback    bulkCallback[T]      // The callback to execute when a bulk is full.
	timeoutCallback timeoutCallback[T]   // The callback to execute on timeout.
	shutdownSignal  chan (chan struct{}) // A channel to signal the timeout processing goroutine to stop.
	isShutDown      bool                 // A flag to indicate if the aggregator has been shut down.
}

// SetBulkCallback sets the callback function to be executed when a bulk reaches its target size.
//
// Parameters:
//
//	callback: The function to call with the full bulk.
func (agg *Aggregator[T]) SetBulkCallback(callback bulkCallback[T]) {
	agg.bulkCallback = callback
}

// SetTimeoutCallback sets the callback function to be executed when a timeout occurs.
//
// Parameters:
//
//	callback: The function to call with the items that have timed out.
func (agg *Aggregator[T]) SetTimeoutCallback(callback timeoutCallback[T]) {
	agg.timeoutCallback = callback
}

// Add appends an item to the aggregator. If adding the item causes the bulk to reach
// its target size, the bulk is processed using the bulk callback.
//
// Parameters:
//
//	item: The item to add to the aggregator.
func (agg *Aggregator[T]) Add(item T) {
	agg.mutex.Lock()
	defer agg.mutex.Unlock()

	if agg.isShutDown {
		return
	}

	agg.items = append(agg.items, item)

	if len(agg.items) >= agg.bulkSize {
		bulk := agg.items
		agg.items = make([]T, 0)

		if agg.bulkCallback != nil {
			// Invoke callback in a separate goroutine to avoid blocking.
			go agg.bulkCallback(bulk)
		}
	}
}

// Count returns the current number of items in the aggregator.
//
// Returns:
//
//	The number of items.
func (agg *Aggregator[T]) Count() int {
	agg.mutex.Lock()
	defer agg.mutex.Unlock()
	return len(agg.items)
}

// Close gracefully shuts down the aggregator. It stops the timeout processing goroutine
// and purges any remaining items. It is safe to call Close multiple times.
func (agg *Aggregator[T]) Close() {
	agg.mutex.Lock()
	if agg.isShutDown {
		agg.mutex.Unlock()
		return
	}

	agg.isShutDown = true
	agg.mutex.Unlock()

	feedback := make(chan struct{})
	agg.shutdownSignal <- feedback
	<-feedback
	close(agg.shutdownSignal)

	agg.Purge()
}

// Purge removes all items from the aggregator without triggering any callbacks.
func (agg *Aggregator[T]) Purge() {
	agg.mutex.Lock()
	defer agg.mutex.Unlock()
	agg.items = make([]T, 0)
}

// startBulkTimeoutProcess runs in a separate goroutine to handle timeouts.
// It periodically checks if a bulk has timed out and, if so, processes it.
func (agg *Aggregator[T]) startBulkTimeoutProcess() {
	timer := time.NewTimer(agg.timeout)
	for {
		select {
		case shutdownFeedback := <-agg.shutdownSignal:
			timer.Stop()
			shutdownFeedback <- struct{}{}
			return
		case <-timer.C:
			agg.mutex.Lock()
			if len(agg.items) > 0 && agg.timeoutCallback != nil {
				bulk := agg.items
				agg.items = make([]T, 0)
				// Invoke callback in a separate goroutine to avoid blocking.
				go agg.timeoutCallback(bulk)
			}
			agg.mutex.Unlock()
			timer.Reset(agg.timeout)
		}
	}
}

// NewAggregator creates and returns a new Aggregator instance.
// It initializes the aggregator with the specified bulk size and timeout, and starts the timeout processing goroutine.
//
// Type Parameters:
//
//	T: The type of the items to be aggregated.
//
// Parameters:
//
//	bulkSize: The number of items to collect before processing a bulk.
//	timeout: The duration to wait before processing a non-full bulk.
//
// Returns:
//
//	A pointer to the newly created Aggregator.
func NewAggregator[T any](bulkSize int, timeout time.Duration) *Aggregator[T] {
	agg := &Aggregator[T]{
		items:          make([]T, 0),
		timeout:        timeout,
		bulkSize:       bulkSize,
		shutdownSignal: make(chan chan struct{}),
		isShutDown:     false,
	}
	go agg.startBulkTimeoutProcess()
	return agg
}
