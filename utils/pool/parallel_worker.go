package pool

import (
	"log"
	"runtime"
	"sync"
)

// parallelWorker represents a single worker in a parallel processing pool.
// It is responsible for receiving items and applying a processor function to them.
//
// Type Parameters:
//
//	T: The type of the items to be processed.
type parallelWorker[T any] struct {
	id        int
	done      *sync.WaitGroup
	readyPool chan chan T
	work      chan T
	quit      chan bool
	processor func(T)
}

// newParallelWorker creates and returns a new parallelWorker instance.
//
// Parameters:
//
//	id: The unique identifier for the worker.
//	readyPool: A channel used to register the worker as ready to receive items.
//	done: A WaitGroup to signal when the worker has finished.
//
// Returns:
//
//	A new parallelWorker instance.
func newParallelWorker[T any](id int, readyPool chan chan T, done *sync.WaitGroup) *parallelWorker[T] {
	return &parallelWorker[T]{
		id:        id,
		done:      done,
		readyPool: readyPool,
		work:      make(chan T),
		quit:      make(chan bool),
	}
}

// Start begins the worker's main loop in a new goroutine. The worker registers itself
// with the ready pool and waits for items to be assigned. It can be stopped via the quit channel.
//
// Parameters:
//
//	processor: The function to be applied to each item.
func (w *parallelWorker[T]) Start(processor func(T)) {
	w.processor = processor
	go func() {
		w.done.Add(1)
		defer w.done.Done()

		for {
			// Register with the ready pool to signal availability.
			w.readyPool <- w.work
			select {
			case item := <-w.work:
				// Received an item, process it.
				w.process(item)
			case <-w.quit:
				// Received a quit signal, exit the loop.
				return
			}
		}
	}()
}

// process applies the processor function to a given item.
// It includes panic recovery to ensure the worker does not crash.
func (w *parallelWorker[T]) process(item T) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("panic in parallel worker %d processing item: %v\n%s\n", w.id, r, buf)
		}
	}()
	w.processor(item)
}

// Stop sends a signal to the worker to stop its processing loop.
// The worker will finish its current item before stopping.
func (w *parallelWorker[T]) Stop() {
	w.quit <- true
}
