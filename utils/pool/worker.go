package pool

import (
	"log"
	"runtime"
	"sync"
)

type worker[T any] struct {
	id        int
	done      *sync.WaitGroup
	readyPool chan chan Task[T]
	work      chan Task[T]
	quit      chan bool
	callback  func(T)
}

func NewWorker[T any](id int, readyPool chan chan Task[T], done *sync.WaitGroup) *worker[T] {
	return &worker[T]{
		id:        id,
		done:      done,
		readyPool: readyPool,
		work:      make(chan Task[T]),
		quit:      make(chan bool),
	}
}

func (w *worker[T]) Process(task Task[T]) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("panic running process: %v\n%s\n", r, buf)
		}
	}()
	result := task.Run()

	// If callback defined, invoke callback
	if w.callback != nil {
		go w.callback(result)
	}
}

// Start wait for tasks with optional
func (w *worker[T]) Start(callback func(T)) {
	w.callback = callback
	go func() {
		w.done.Add(1)
		for {
			w.readyPool <- w.work
			select {
			case work := <-w.work:
				w.Process(work)
			case <-w.quit:
				w.done.Done()
				return
			}
		}
	}()
}

// Stop notify worker to stop after current process
func (w *worker[T]) Stop() {
	w.quit <- true
}
