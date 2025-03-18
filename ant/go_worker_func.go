package ant

import (
	"fmt"
	"time"
)

// goWorkerWithFunc is the actual executor who runs the tasks,
// it starts a goroutine that accepts tasks and
// performs function calls.
type goWorkerWithFunc struct {
	worker

	// pool who owns this worker.
	pool *PoolWithFunc

	// arg is the argument for the function.
	arg chan any

	// lastUsed will be updated when putting a worker back into queue.
	lastUsed time.Time
}

func (w *goWorkerWithFunc) inputArg(arg any) {
	w.arg <- arg
}

func (w *goWorkerWithFunc) finish() {
	w.arg <- nil
}

// run starts a goroutine to repeat the process
// that performs the function calls.
func (w *goWorkerWithFunc) run() {
	w.pool.addRunning(1)
	defer func() {
		if w.pool.addRunning(-1) == 0 && w.pool.IsClosed() {
			w.pool.once.Do(func() {
				close(w.pool.allDone)
			})
		}

		if r := recover(); r != nil {
			if ph := w.pool.options.PanicHandler; ph != nil {
				ph(w)
			} else {
				fmt.Printf("goWorkerWithFunc run error: %v", r)
			}
		}
		// Call Signal() here in case there are goroutines waiting for available workers.
		w.pool.cond.Signal()
	}()

	for arg := range w.arg {
		if arg == nil {
			return
		}
		w.pool.fn(arg)
		if ok := w.pool.revertWorker(w); !ok {
			return
		}
	}
}

func (w *goWorkerWithFunc) lastUsedTime() time.Time {
	return w.lastUsed
}

func (w *goWorkerWithFunc) setLastUsedTime(t time.Time) {
	w.lastUsed = t
}
