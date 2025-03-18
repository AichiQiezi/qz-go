package ant

import (
	"errors"
	"runtime"
)

var (
	// ErrQueueIsFull will be returned when the worker queue is full.
	ErrQueueIsFull = errors.New("the queue is full")

	// ErrLackPoolFunc will be returned when invokers don't provide function for pool.
	ErrLackPoolFunc = errors.New("must provide function for pool")

	// ErrInvalidPoolExpiry will be returned when setting a negative number as the periodic duration to purge goroutines.
	ErrInvalidPoolExpiry = errors.New("invalid expiry for pool")

	// ErrPoolClosed will be returned when submitting task to a closed pool.
	ErrPoolClosed = errors.New("this pool has been closed")

	// ErrPoolOverload will be returned when the pool is full and no workers available.
	ErrPoolOverload = errors.New("too many goroutines blocked on submit or Nonblocking is set")

	// ErrInvalidPreAllocSize will be returned when trying to set up a negative capacity under PreAlloc mode.
	ErrInvalidPreAllocSize = errors.New("can not set up a negative capacity under PreAlloc mode")

	// ErrTimeout will be returned after the operations timed out.
	ErrTimeout = errors.New("operation timed out")

	// ErrInvalidPoolIndex will be returned when trying to retrieve a pool with an invalid index.
	ErrInvalidPoolIndex = errors.New("invalid pool index")

	// ErrInvalidLoadBalancingStrategy will be returned when trying to create a MultiPool with an invalid load-balancing strategy.
	ErrInvalidLoadBalancingStrategy = errors.New("invalid load-balancing strategy")

	// ErrInvalidMultiPoolSize  will be returned when trying to create a MultiPool with an invalid size.
	ErrInvalidMultiPoolSize = errors.New("invalid size for multiple pool")

	// workerChanCap determines whether the channel of a worker should be a buffered channel
	// to get the best performance. Inspired by fasthttp at
	// https://github.com/valyala/fasthttp/blob/master/workerpool.go#L139
	workerChanCap = func() int {
		// Use blocking channel if GOMAXPROCS=1.
		// This switches context from sender to receiver immediately,
		// which results in higher performance (under go1.5 at least).
		if runtime.GOMAXPROCS(0) == 1 {
			return 0
		}

		// Use non-blocking workerChan if GOMAXPROCS>1,
		// since otherwise the sender might be dragged down if the receiver is CPU-bound.
		return 1
	}()
)
