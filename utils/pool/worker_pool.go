package pool

import (
	"fmt"
	"sync/atomic"
	"time"
)

// WorkerPool process all tasks
type WorkerPool struct {
	queueSize  int
	numWorkers int
	workers    []*Worker
	jobQueue   chan Task
	workerPool chan chan Task
	quit       chan bool
	results    chan any
	refCount   int64
}

// NewWorkerPool factory method to create the worker pool with number of worker threads and limit the size of the task queue
func NewWorkerPool(workers int, size int, results chan any) *WorkerPool {
	workerPool := &WorkerPool{queueSize: size, numWorkers: workers}
	workerPool.jobQueue = make(chan Task, size)
	workerPool.workerPool = make(chan chan Task, workers)
	workerPool.quit = make(chan bool)
	workerPool.results = results
	workerPool.createPool()
	return workerPool
}

// Execute submits the job to the queue, return error if the queue is full
func (t *WorkerPool) Execute(task Task) error {
	if len(t.jobQueue) == t.queueSize {
		return fmt.Errorf("the job queue is full")
	}
	atomic.AddInt64(&t.refCount, 1)
	t.jobQueue <- task
	return nil
}

// WaitForAll for all jobs to complete
func (t *WorkerPool) WaitForAll() {
	for t.isDone() == false {
		time.Sleep(time.Second)
	}
}

func (t *WorkerPool) isDone() bool {

	fmt.Println("Jobs reference count:", t.refCount)

	// Check job queue length
	//if len(t.jobQueue) > 0 {
	//	return false
	//}
	return t.refCount <= 0
}

// Close will close the worker pool and terminate all waiting jobs
// It sends the stop signal to all the worker that are running
func (t *WorkerPool) Close() {
	close(t.quit)       // Stops all the routines
	close(t.workerPool) // Closes the Job worker pool
	close(t.jobQueue)   // Closes the job Queue
}

// createPool creates the workers and start listening on the jobQueue
func (t *WorkerPool) createPool() {
	for i := 0; i < t.numWorkers; i++ {
		worker := NewWorker(t.workerPool, t.quit, &t.refCount, t.results)
		worker.Start()
	}
	go t.startDispatcher()
}

// start listen to the jobs queue and dispatch the jobs to a worker
func (t *WorkerPool) startDispatcher() {
	for {
		select {

		case job := <-t.jobQueue:
			// Got job
			func(job Task) {
				// Find a worker for the job
				jobChannel := <-t.workerPool
				// Submit job to the worker
				jobChannel <- job
			}(job)

		case <-t.quit:
			// Close the worker pool
			return
		}
	}
}
