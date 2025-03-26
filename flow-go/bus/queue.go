package bus

import "sync"

// SafeQueue 是一个并发安全的队列
type SafeQueue[T any] struct {
	items []T
	mu    sync.Mutex
	cond  *sync.Cond
}

// NewSafeQueue 创建一个新的 SafeQueue
func NewSafeQueue[T any]() *SafeQueue[T] {
	q := &SafeQueue[T]{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Enqueue 向队列添加元素
func (q *SafeQueue[T]) Enqueue(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, item)
	q.cond.Signal() // 唤醒一个等待的 Goroutine
}

// Dequeue 阻塞等待直到队列有元素可取，并返回第一个元素
func (q *SafeQueue[T]) Dequeue() T {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == 0 {
		q.cond.Wait() // 释放锁并等待
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item
}

// TryDequeue 立即返回队列中的第一个元素，如果队列为空则返回 false
func (q *SafeQueue[T]) TryDequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		var zero T
		return zero, false
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Len 返回队列的长度
func (q *SafeQueue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}
