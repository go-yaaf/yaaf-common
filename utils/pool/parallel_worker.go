package pool

import (
	"sync"
)

type parallelWorker[T any] struct {
	id        int
	done      *sync.WaitGroup
	readyPool chan chan T
	work      chan T
	quit      chan bool
	processor func(T)
}

func newParallelWorker[T any](id int, readyPool chan chan T, done *sync.WaitGroup) *parallelWorker[T] {
	return &parallelWorker[T]{
		id:        id,
		done:      done,
		readyPool: readyPool,
		work:      make(chan T),
		quit:      make(chan bool),
	}
}

// Start wait for tasks with optional
func (w *parallelWorker[T]) Start(processor func(T)) {
	w.processor = processor
	go func() {
		w.done.Add(1)
		for {
			w.readyPool <- w.work
			select {
			case item := <-w.work:
				w.processor(item)
			case <-w.quit:
				w.done.Done()
				return
			}
		}
	}()
}

// Stop notify worker to stop after current process
func (w *parallelWorker[T]) Stop() {
	w.quit <- true
}

// WaitAll notify worker to stop after current process
func (w *parallelWorker[T]) WaitAll() {
	w.quit <- true
}
