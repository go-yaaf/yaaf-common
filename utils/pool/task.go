package pool

// Task is interface for any job executed by the worker pool
type Task[T any] interface {
	Run() T
}
