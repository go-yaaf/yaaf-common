package pool

import "sync/atomic"

// Worker type holds the job channel and passed worker threadpool
type Worker struct {
	jobChannel     chan Task
	workerPool     chan chan Task
	quit           chan bool
	resultChannel  chan any
	referenceCount *int64
}

// NewWorker creates the new worker
func NewWorker(workerPool chan chan Task, quit chan bool, rc *int64, results chan any) *Worker {
	return &Worker{workerPool: workerPool, jobChannel: make(chan Task), quit: quit, resultChannel: results, referenceCount: rc}
}

// Start starts the worker by listening to the job channel
func (w Worker) Start() {
	go func() {
		for {
			// Put the job into the worker pool
			w.workerPool <- w.jobChannel

			select {
			// Wait for the tasks in the task channel
			case task := <-w.jobChannel:
				// Got the task and run it
				result := task.Run()
				atomic.AddInt64(w.referenceCount, -1)
				if w.resultChannel != nil {
					w.resultChannel <- result
				}
			case <-w.quit:
				// Exit the go routine when the quit channel is closed
				return
			}
		}
	}()
}
