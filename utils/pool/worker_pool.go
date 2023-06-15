package pool

import (
	"fmt"
	"sync"
	"time"
)

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

func NewWorkerPool[T any](maxWorkers int, capacity int) *WorkerPool[T] {
	if capacity <= 0 {
		capacity = 100
	}
	if maxWorkers <= 0 {
		maxWorkers = 1
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

// Start pool with optional callback to be invoked on task completion
func (q *WorkerPool[T]) Start(callback func(T)) error {

	if len(q.workers) > 0 {
		return fmt.Errorf("pool already started")
	}

	// Create and start workers
	q.workers = make([]*worker[T], q.maxWorkers, q.maxWorkers)

	// create and start workers
	for i := 0; i < q.maxWorkers; i++ {
		q.workers[i] = NewWorker[T](i+1, q.readyPool, q.workersStopped)
		q.workers[i].Start(callback)
	}
	go q.dispatch()
	return nil
}

func (q *WorkerPool[T]) Stop() {
	q.quit <- true
	q.dispatcherStopped.Wait()
}

func (q *WorkerPool[T]) dispatch() {
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
func (q *WorkerPool[T]) Submit(task Task[T]) {
	q.singleJob <- task
}

// Enqueue submits a task to the buffered job queue without blocking, returns false if queue is full
func (q *WorkerPool[T]) Enqueue(task Task[T]) bool {
	select {
	case q.queuedJob <- task:
		return true
	default:
		return false
	}
}

// EnqueueWithTimeout submits a task to the buffered job queue without blocking, returns false if queue is full within the duration
func (q *WorkerPool[T]) EnqueueWithTimeout(task Task[T], timeout time.Duration) bool {
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
	case q.queuedJob <- task:
		return true
	case <-ch:
		return false
	}
}
