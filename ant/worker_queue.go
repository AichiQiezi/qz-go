package ant

import "time"

type workerQueue interface {
	len() int
	isEmpty() bool
	insert(worker) error
	detach() worker
	refresh(duration time.Duration) []worker // clean up the stale workers and return them
	reset()
}

type queueType int

const (
	queueTypeLoopQueue queueType = 1
)

func newWorkerQueue(qType queueType, size int) workerQueue {
	switch qType {
	case queueTypeLoopQueue:
		return newWorkerLoopQueue(size)
	}

	return newWorkerLoopQueue(size)
}
