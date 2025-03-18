package ant

import "time"

type workerLoopQueue struct {
	items  []worker
	expiry []worker
	head   int
	tail   int
	size   int
	isFull bool
}

func newWorkerLoopQueue(size int) *workerLoopQueue {
	if size <= 0 {
		return nil
	}
	return &workerLoopQueue{
		items: make([]worker, size),
		size:  size,
	}
}

func (wq *workerLoopQueue) isEmpty() bool {
	return wq.head == wq.tail && !wq.isFull
}

func (wq *workerLoopQueue) insert(w worker) error {
	if wq.isFull {
		return ErrQueueIsFull
	}

	wq.items[wq.tail] = w
	wq.tail = (wq.tail + 1) % wq.size

	if wq.head == wq.tail {
		wq.isFull = true
	}

	return nil
}

func (wq *workerLoopQueue) detach() worker {
	if wq.isEmpty() {
		return nil
	}

	w := wq.items[wq.head]
	wq.items[wq.head] = nil
	wq.head = (wq.head + 1) % wq.size

	wq.isFull = false

	return w
}

func (wq *workerLoopQueue) refresh(duration time.Duration) []worker {
	return nil
}

func (wq *workerLoopQueue) reset() {
	for !wq.isEmpty() {
		if w := wq.detach(); w != nil {
			w.finish()
		}
	}
	wq.head, wq.tail = 0, 0
}

func (wq *workerLoopQueue) len() int {
	if wq.size == 0 || wq.isEmpty() {
		return 0
	}

	if wq.head == wq.tail && wq.isFull {
		return wq.size
	}

	if wq.tail > wq.head {
		return wq.tail - wq.head
	}

	return wq.size - wq.head + wq.tail
}
