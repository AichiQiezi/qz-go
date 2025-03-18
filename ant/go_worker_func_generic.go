package ant

import "time"

// goWorkerWithFunc is the actual executor who runs the tasks,
// it starts a goroutine that accepts tasks and
// performs function calls.
type goWorkerWithFuncGeneric[T any] struct {
	worker

	// pool who owns this worker.
	pool *PoolWithFuncGeneric[T]

	// arg is a job should be done.
	arg chan T

	// exit signals the goroutine to exit.
	exit chan struct{}

	// lastUsed will be updated when putting a worker back into queue.
	lastUsed time.Time
}
