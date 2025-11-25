package pool

import (
	"fmt"
	"sync"
	"time"
)

// Parallel provides a mechanism to process multiple items in parallel using a fixed number of workers,
// each applying the same processing function. This is useful for scenarios where you have a large number
// of similar items to process and want to leverage concurrency to speed up the work.
//
// Type Parameters:
//
//	T: The type of the items to be processed.
type Parallel[T any] struct {
	singleJob         chan T
	queuedJob         chan T
	readyPool         chan chan T
	workers           []*parallelWorker[T]
	maxWorkers        int
	dispatcherStopped *sync.WaitGroup
	workersStopped    *sync.WaitGroup
	quit              chan bool
}

// NewParallel creates and returns a new Parallel instance.
//
// Parameters:
//
//	maxWorkers: The maximum number of concurrent workers.
//	capacity: The capacity of the buffered queue for items to be processed.
//
// Returns:
//
//	A new Parallel instance.
func NewParallel[T any](maxWorkers int, capacity int) *Parallel[T] {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}
	if capacity <= 0 {
		capacity = 100
	}

	return &Parallel[T]{
		singleJob:         make(chan T),
		queuedJob:         make(chan T, capacity),
		readyPool:         make(chan chan T, maxWorkers),
		dispatcherStopped: &sync.WaitGroup{},
		workersStopped:    &sync.WaitGroup{},
		quit:              make(chan bool),
		maxWorkers:        maxWorkers,
	}
}

// Start initializes and starts the parallel processing pool. It creates the workers
// and starts the dispatcher to distribute items to them.
//
// Parameters:
//
//	processor: The function that will be applied to each item.
//
// Returns:
//
//	An error if the pool has already been started.
func (p *Parallel[T]) Start(processor func(T)) error {
	if len(p.workers) > 0 {
		return fmt.Errorf("parallel pool already started")
	}

	p.workers = make([]*parallelWorker[T], p.maxWorkers)
	for i := 0; i < p.maxWorkers; i++ {
		p.workers[i] = newParallelWorker[T](i+1, p.readyPool, p.workersStopped)
		p.workers[i].Start(processor)
	}

	go p.dispatch()
	return nil
}

// WaitAll blocks until all submitted items have been processed and the pool is stopped.
func (p *Parallel[T]) WaitAll() {
	p.quit <- true
	p.dispatcherStopped.Wait()
}

// dispatch is the main loop for the dispatcher. It waits for items to be submitted
// and assigns them to available workers.
func (p *Parallel[T]) dispatch() {
	p.dispatcherStopped.Add(1)
	defer p.dispatcherStopped.Done()

	for {
		select {
		case item := <-p.singleJob:
			workerChannel := <-p.readyPool
			workerChannel <- item
		case item := <-p.queuedJob:
			workerChannel := <-p.readyPool
			workerChannel <- item
		case <-p.quit:
			for _, w := range p.workers {
				w.Stop()
			}
			p.workersStopped.Wait()
			return
		}
	}
}

// Submit adds an item to be processed. This method will block until a worker is available.
//
// Parameters:
//
//	item: The item to be processed.
func (p *Parallel[T]) Submit(item T) {
	p.singleJob <- item
}

// Enqueue adds an item to the buffered queue to be processed. This method is non-blocking
// and will return `false` if the queue is full.
//
// Parameters:
//
//	item: The item to be enqueued.
//
// Returns:
//
//	`true` if the item was successfully enqueued, `false` otherwise.
func (p *Parallel[T]) Enqueue(item T) bool {
	select {
	case p.queuedJob <- item:
		return true
	default:
		return false
	}
}

// EnqueueWithTimeout adds an item to the buffered queue with a timeout.
// It will wait for the specified duration for space to become available.
//
// Parameters:
//
//	item: The item to be enqueued.
//	timeout: The maximum time to wait.
//
// Returns:
//
//	`true` if the item was successfully enqueued, `false` if the timeout was reached.
func (p *Parallel[T]) EnqueueWithTimeout(item T, timeout time.Duration) bool {
	if timeout <= 0 {
		timeout = 1 * time.Second
	}

	select {
	case p.queuedJob <- item:
		return true
	case <-time.After(timeout):
		return false
	}
}
