package pool

// Task is interface for any job executed by the worker pool
type Task interface {
	Run() any
}
