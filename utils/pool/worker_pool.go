package pool

import (
	"fmt"
	"sync"
	"time"
)

// WorkerPool manages a pool of workers to execute tasks concurrently.
// It provides methods for submitting tasks and managing the lifecycle of the pool.
//
// Type Parameters:
//
//	T: The type of the result returned by the tasks.
type WorkerPool[T any] struct {
	singleJob         chan Task[T]
	queuedJob         chan Task[T]
	readyPool         chan chan Task[T]
	workers           []*worker[T]
	maxWorkers        int
	dispatcherStopped *sync.WaitGroup
	workersStopped    *sync.WaitGroup
	quit              chan bool
}

// NewWorkerPool creates and returns a new WorkerPool.
//
// Parameters:
//
//	maxWorkers: The maximum number of workers in the pool.
//	capacity: The capacity of the buffered job queue.
//
// Returns:
//
//	A new WorkerPool instance.
func NewWorkerPool[T any](maxWorkers int, capacity int) *WorkerPool[T] {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}
	if capacity <= 0 {
		capacity = 100
	}

	return &WorkerPool[T]{
		singleJob:         make(chan Task[T]),
		queuedJob:         make(chan Task[T], capacity),
		readyPool:         make(chan chan Task[T], maxWorkers),
		dispatcherStopped: &sync.WaitGroup{},
		workersStopped:    &sync.WaitGroup{},
		quit:              make(chan bool),
		maxWorkers:        maxWorkers,
	}
}

// Start initializes and starts the worker pool. It creates the specified number of workers
// and starts the dispatcher to distribute tasks.
//
// Parameters:
//
//	callback: An optional function to be called with the result of each completed task.
//
// Returns:
//
//	An error if the pool has already been started.
func (q *WorkerPool[T]) Start(callback func(T)) error {
	if len(q.workers) > 0 {
		return fmt.Errorf("worker pool already started")
	}

	q.workers = make([]*worker[T], q.maxWorkers)
	for i := 0; i < q.maxWorkers; i++ {
		q.workers[i] = NewWorker[T](i+1, q.readyPool, q.workersStopped)
		q.workers[i].Start(callback)
	}

	go q.dispatch()
	return nil
}

// Stop gracefully shuts down the worker pool. It signals the dispatcher to stop,
// waits for all workers to finish their current tasks, and then stops the dispatcher.
func (q *WorkerPool[T]) Stop() {
	q.quit <- true
	q.dispatcherStopped.Wait()
}

// dispatch is the main loop for the dispatcher. It waits for tasks to be submitted
// and assigns them to available workers.
func (q *WorkerPool[T]) dispatch() {
	q.dispatcherStopped.Add(1)
	defer q.dispatcherStopped.Done()

	for {
		select {
		case job := <-q.singleJob:
			// Wait for a free worker and dispatch the job.
			workerChannel := <-q.readyPool
			workerChannel <- job
		case job := <-q.queuedJob:
			// Wait for a free worker and dispatch the job.
			workerChannel := <-q.readyPool
			workerChannel <- job
		case <-q.quit:
			// Stop all workers and wait for them to finish.
			for _, w := range q.workers {
				w.Stop()
			}
			q.workersStopped.Wait()
			return
		}
	}
}

// Submit adds a task to the job queue. This method will block until a worker is available to process the task.
//
// Parameters:
//
//	task: The task to be submitted.
func (q *WorkerPool[T]) Submit(task Task[T]) {
	q.singleJob <- task
}

// Enqueue adds a task to the buffered job queue. This method is non-blocking and will return
// `false` if the queue is full.
//
// Parameters:
//
//	task: The task to be enqueued.
//
// Returns:
//
//	`true` if the task was successfully enqueued, `false` otherwise.
func (q *WorkerPool[T]) Enqueue(task Task[T]) bool {
	select {
	case q.queuedJob <- task:
		return true
	default:
		return false
	}
}

// EnqueueWithTimeout adds a task to the buffered job queue with a timeout.
// It will wait for the specified duration for space to become available in the queue.
//
// Parameters:
//
//	task: The task to be enqueued.
//	timeout: The maximum time to wait for the task to be enqueued.
//
// Returns:
//
//	`true` if the task was successfully enqueued, `false` if the timeout was reached.
func (q *WorkerPool[T]) EnqueueWithTimeout(task Task[T], timeout time.Duration) bool {
	if timeout <= 0 {
		timeout = 1 * time.Second
	}

	select {
	case q.queuedJob <- task:
		return true
	case <-time.After(timeout):
		return false
	}
}
