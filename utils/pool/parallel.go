package pool

import (
	"fmt"
	"sync"
	"time"
)

// Parallel is a concurrency utility enable user to use the same processor (function) to process multiple objects
// in parallel using multiple workers.
// This utility is similar to WorkerPool but the main difference is that it is using the same function to reduce number of allocations
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

func NewParallel[T any](maxWorkers int, capacity int) *Parallel[T] {
	if capacity <= 0 {
		capacity = 100
	}
	if maxWorkers <= 0 {
		maxWorkers = 1
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

// Start pool with processor function to be called for any object in the queue
func (q *Parallel[T]) Start(processor func(T)) error {

	if len(q.workers) > 0 {
		return fmt.Errorf("pool already started")
	}

	// Create and start workers
	q.workers = make([]*parallelWorker[T], q.maxWorkers, q.maxWorkers)

	// create and start workers
	for i := 0; i < q.maxWorkers; i++ {
		q.workers[i] = newParallelWorker[T](i+1, q.readyPool, q.workersStopped)
		q.workers[i].Start(processor)
	}
	go q.dispatch()
	return nil
}

// WaitAll blocks until completion of all tasks in the queue
func (q *Parallel[T]) WaitAll() {
	q.quit <- true
	q.dispatcherStopped.Wait()
}

func (q *Parallel[T]) dispatch() {
	// start the dispatcher
	q.dispatcherStopped.Add(1)
	for {
		select {
		case job := <-q.singleJob:
			workerXChannel := <-q.readyPool // wait for free worker
			workerXChannel <- job           // dispatch job to the free worker

		case job := <-q.queuedJob:
			workerXChannel := <-q.readyPool // wait for free worker
			workerXChannel <- job           // dispatch job to the free worker

		case <-q.quit:
			// free all workers
			for i := 0; i < len(q.workers); i++ {
				q.workers[i].Stop()
			}
			// wait for all workers to finish their job
			q.workersStopped.Wait()

			// stop the dispatcher
			q.dispatcherStopped.Done()
			return
		}
	}
}

// Submit a task to the job queue, blocked if no workers are available
func (q *Parallel[T]) Submit(item T) {
	q.singleJob <- item
}

// Enqueue submits a task to the buffered job queue without blocking, returns false if queue is full
func (q *Parallel[T]) Enqueue(item T) bool {
	select {
	case q.queuedJob <- item:
		return true
	default:
		return false
	}
}

// EnqueueWithTimeout submits a task to the buffered job queue without blocking, returns false if queue is full within the duration
func (q *Parallel[T]) EnqueueWithTimeout(item T, timeout time.Duration) bool {
	if timeout <= 0 {
		timeout = 1 * time.Second
	}

	ch := make(chan bool, 1)
	t := time.AfterFunc(timeout, func() {
		ch <- false
	})
	defer func() {
		t.Stop()
	}()

	select {
	case q.queuedJob <- item:
		return true
	case <-ch:
		return false
	}
}
