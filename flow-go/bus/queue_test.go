package bus

import (
	"fmt"
	"testing"
	"time"
)

func TestSafeQueue(t *testing.T) {
	q := NewSafeQueue[int]()

	// 测试 Enqueue 和 Dequeue
	q.Enqueue(1)
	q.Enqueue(2)

	if q.Dequeue() != 1 {
		t.Errorf("expected 1 but got different value")
	}
	if q.Dequeue() != 2 {
		t.Errorf("expected 2 but got different value")
	}

	// 测试 TryDequeue
	q.Enqueue(3)
	if val, ok := q.TryDequeue(); !ok || val != 3 {
		t.Errorf("expected 3 but got %v", val)
	}

	if _, ok := q.TryDequeue(); ok {
		t.Errorf("expected empty queue but got value")
	}
}

func TestSafeQueueConcurrent(t *testing.T) {
	q := NewSafeQueue[int]()

	// 并发写入
	go func() {
		for i := 1; i <= 5; i++ {
			q.Enqueue(i)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 并发读取
	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Println("Dequeued:", q.Dequeue())
		}
	}()

	time.Sleep(1 * time.Second)
}
