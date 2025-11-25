package pool

// Task represents a unit of work that can be executed by a worker in the pool.
// It is a generic interface, allowing for any type of task that returns a value of type T.
//
// Type Parameters:
//
//	T: The type of the result returned by the task.
type Task[T any] interface {
	// Run executes the task and returns a result of type T.
	Run() T
}
