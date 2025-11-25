package pool

import (
	"log"
	"runtime"
	"sync"
)

// worker represents a single worker in the pool. It is responsible for receiving tasks,
// executing them, and handling the results.
type worker[T any] struct {
	id        int
	done      *sync.WaitGroup
	readyPool chan chan Task[T]
	work      chan Task[T]
	quit      chan bool
	callback  func(T)
}

// NewWorker creates and returns a new worker instance.
//
// Parameters:
//
//	id: The unique identifier for the worker.
//	readyPool: A channel used to register the worker as ready to receive tasks.
//	done: A WaitGroup to signal when the worker has finished its work.
//
// Returns:
//
//	A new worker instance.
func NewWorker[T any](id int, readyPool chan chan Task[T], done *sync.WaitGroup) *worker[T] {
	return &worker[T]{
		id:        id,
		done:      done,
		readyPool: readyPool,
		work:      make(chan Task[T]),
		quit:      make(chan bool),
	}
}

// Process executes a given task. It includes panic recovery to ensure the worker
// does not crash unexpectedly. If a callback is provided, it is invoked with the task's result.
//
// Parameters:
//
//	task: The task to be processed.
func (w *worker[T]) Process(task Task[T]) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("panic in worker %d processing task: %v\n%s\n", w.id, r, buf)
		}
	}()
	result := task.Run()

	if w.callback != nil {
		go w.callback(result)
	}
}

// Start begins the worker's main loop in a new goroutine. The worker registers itself
// with the ready pool and waits for tasks to be assigned. It can be stopped via the quit channel.
//
// Parameters:
//
//	callback: A function to be called with the result of each processed task.
func (w *worker[T]) Start(callback func(T)) {
	w.callback = callback
	go func() {
		w.done.Add(1)
		defer w.done.Done()

		for {
			// Register with the ready pool to signal availability.
			w.readyPool <- w.work
			select {
			case task := <-w.work:
				// Received a task, process it.
				w.Process(task)
			case <-w.quit:
				// Received a quit signal, exit the loop.
				return
			}
		}
	}()
}

// Stop sends a signal to the worker to stop its processing loop.
// The worker will finish its current task before stopping.
func (w *worker[T]) Stop() {
	w.quit <- true
}
