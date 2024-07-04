/*
 * Aggregator is a mechanism to aggregate items to bulks and notify when the bulk reached a certain size or after time out
 */

package aggregator

import (
	"sync"
	"time"
)

// bulkCallback is called when bulk reached or exceeded the bulk size
type bulkCallback[T any] func(bulk []T)

// timeoutCallback is called when bulk did not reach the bulk size but a timeout aws triggered
type timeoutCallback[T any] func(bulk []T)

// Aggregator is a synchronized map of items that can auto-expire once stale
type Aggregator[T any] struct {
	mutex           sync.Mutex    // Mutex for sync operations
	timeout         time.Duration // Timeout no notify when bulk was not yet created
	bulkSize        int           // Bulk size
	items           []T
	bulkCallback    bulkCallback[T]
	timeoutCallback timeoutCallback[T]
	shutdownSignal  chan (chan struct{})
	isShutDown      bool
}

// SetBulkCallback sets the callback on bulk creation
func (agg *Aggregator[T]) SetBulkCallback(callback bulkCallback[T]) {
	agg.bulkCallback = callback
}

// SetTimeoutCallback sets the callback on timeout
func (agg *Aggregator[T]) SetTimeoutCallback(callback timeoutCallback[T]) {
	agg.timeoutCallback = callback
}

// Add item to the aggregator
func (agg *Aggregator[T]) Add(item T) {

	agg.mutex.Lock()
	agg.items = append(agg.items, item)

	// return if number of items is less than bulk size
	if len(agg.items) < agg.bulkSize {
		agg.mutex.Unlock()
		return
	}

	// Move items to bulk and invoke callback
	bulk := make([]T, 0)
	bulk = append(bulk, agg.items...)
	agg.items = make([]T, 0)
	agg.mutex.Unlock()

	// Invoke callback
	if agg.bulkCallback != nil {
		agg.bulkCallback(bulk)
	}
}

// Count returns the number of items in the aggregator
func (agg *Aggregator[T]) Count() int {
	agg.mutex.Lock()
	length := len(agg.items)
	agg.mutex.Unlock()
	return length
}

// Close calls Purge, and then stops the goroutine that does ttl checking, for a clean shutdown.
// The cache is no longer cleaning up after the first call to Close, repeated calls are safe though.
func (agg *Aggregator[T]) Close() {
	agg.mutex.Lock()
	if !agg.isShutDown {
		agg.isShutDown = true
		agg.mutex.Unlock()
		feedback := make(chan struct{})
		agg.shutdownSignal <- feedback
		<-feedback
		close(agg.shutdownSignal)
	} else {
		agg.mutex.Unlock()
	}
	agg.Purge()
}

// Purge will remove all entries
func (agg *Aggregator[T]) Purge() {
	agg.mutex.Lock()
	agg.items = make([]T, 0)
	agg.mutex.Unlock()
}

// start the timeout thread
func (agg *Aggregator[T]) startBulkTimeoutProcess() {
	timer := time.NewTimer(agg.timeout)
	for {
		timer.Reset(agg.timeout)
		select {
		case shutdownFeedback := <-agg.shutdownSignal:
			timer.Stop()
			shutdownFeedback <- struct{}{}
			return
		case <-timer.C:
			timer.Stop()
			agg.mutex.Lock()
			if len(agg.items) == 0 {
				agg.mutex.Unlock()
				continue
			} else {
				bulk := make([]T, 0)
				bulk = append(bulk, agg.items...)
				if agg.timeoutCallback != nil {
					agg.items = make([]T, 0)
					agg.timeoutCallback(bulk)
				}
			}
			agg.mutex.Unlock()
		}
	}
}

// NewAggregator is a helper to create instance of the aggregator
func NewAggregator[T any](bulkSize int, timeout time.Duration) *Aggregator[T] {
	shutdownChan := make(chan chan struct{})
	agg := &Aggregator[T]{
		items:          make([]T, 0),
		timeout:        timeout,
		bulkSize:       bulkSize,
		shutdownSignal: shutdownChan,
		isShutDown:     false,
	}
	go agg.startBulkTimeoutProcess()
	return agg
}
