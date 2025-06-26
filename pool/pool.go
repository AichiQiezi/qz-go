package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Task func()

type Pool struct {
	tasks chan Task

	wg     sync.WaitGroup
	cancel context.CancelFunc
	ctx    context.Context
	mu     sync.Mutex
	closed bool
}

func NewPool(workerCount int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	p := &Pool{
		tasks:  make(chan Task, workerCount),
		wg:     sync.WaitGroup{},
		mu:     sync.Mutex{},
		ctx:    ctx,
		cancel: cancel,
	}

	// 启动 worker
	for i := 0; i < workerCount; i++ {
		go p.worker(i)
	}

	return p
}

func (p *Pool) Submit(task Task) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	if task == nil {
		return
	}

	p.wg.Add(1)
	p.tasks <- task
}

// worker 是单个工作协程
func (p *Pool) worker(id int) {
	for {
		select {
		case <-p.ctx.Done():
			return
		case task := <-p.tasks:
			task()
			p.wg.Done()
		}
	}
}

// Shutdown 等待所有任务完成并关闭
func (p *Pool) Shutdown() {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	p.mu.Unlock()

	// 等待所有任务完成再关闭
	p.wg.Wait()
	p.cancel()
	close(p.tasks)
}

func main() {
	pool := NewPool(5) // 5个worker

	for i := 0; i < 10; i++ {
		n := i
		pool.Submit(func() {
			fmt.Printf("Task %d start\n", n)
			time.Sleep(time.Second)
			fmt.Printf("Task %d done\n", n)
		})
	}

	pool.Shutdown()
	fmt.Println("All tasks finished")
}
