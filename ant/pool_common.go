package ant

import (
	"math"
	syncx "qz-go/ant/pkg/sync"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// DefaultAntsPoolSize is the default capacity for a default goroutine pool.
	DefaultAntsPoolSize = math.MaxInt32

	// DefaultCleanIntervalTime is the interval time to clean up goroutines.
	DefaultCleanIntervalTime = time.Second
)

type poolCommon struct {
	// capacity of the pool, a negative value means that the capacity of pool is limitless.
	capacity int32

	// running is the number of the currently running goroutines.
	running int32

	// lock for protecting the worker queue.
	lock sync.Locker

	// workers is a slice that store the available workers.
	workers workerQueue

	// state is used to notice the pool to closed itself.
	state int32

	// cond for waiting to get an idle worker.
	cond *sync.Cond

	// done is used to indicate that all workers are done.
	allDone chan struct{}

	// once is used to make sure the pool is closed just once.
	once *sync.Once

	// workerCache speeds up the obtainment of a usable worker in function:retrieveWorker.
	workerCache sync.Pool

	// waiting is the number of goroutines already been blocked on pool.Submit(), protected by pool.lock
	waiting int32

	now atomic.Value

	options *Options
}

func newPool(size int, options ...Option) (*poolCommon, error) {
	if size < 0 {
		size = -1
	}

	opts := loadOptions(options...)

	if !opts.DisablePurge {
		if expiry := opts.ExpiryDuration; expiry < 0 {
			return nil, ErrInvalidPoolExpiry
		} else if expiry == 0 {
			opts.ExpiryDuration = DefaultCleanIntervalTime
		}
	}

	pc := &poolCommon{
		capacity: int32(size),
		allDone:  make(chan struct{}),
		lock:     syncx.NewSpinLock(),
		once:     &sync.Once{},
		options:  opts,
	}

	if pc.options.PreAlloc {
		if size == -1 {
			return nil, ErrInvalidPreAllocSize
		}
		pc.workers = newWorkerQueue(queueTypeLoopQueue, size)
	} else {
		pc.workers = newWorkerQueue(queueTypeLoopQueue, 0)
	}

	pc.cond = sync.NewCond(pc.lock)

	//pc.goPurge()    // 清理大师
	//pc.goTicktock() // ⌛️计时
	return pc, nil
}

// retrieveWorker returns an available worker to run the tasks.
func (p *poolCommon) retrieveWorker() (worker, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for {
		// fast path
		if w := p.workers.detach(); w != nil {
			return w, nil
		}

		if capacity := p.Cap(); capacity == -1 || capacity > p.Running() {
			w := p.workerCache.Get().(worker)
			w.run()
			return w, nil
		}

		if p.options.Nonblocking || (p.options.MaxBlockingTasks != 0 && p.Waiting() > p.options.MaxBlockingTasks) {
			return nil, ErrPoolOverload
		}

		p.addWaiting(1)
		p.cond.Wait()
		p.addWaiting(-1)

		if p.IsClosed() {
			return nil, ErrPoolClosed
		}
	}
}

// revertWorker puts a work into free pool,recycling the goroutines
func (p *poolCommon) revertWorker(w worker) bool {
	if capacity := p.Cap(); (capacity > 0 && p.Running() > capacity) || p.IsClosed() {
		p.cond.Broadcast()
		return false
	}

	// todo
	w.setLastUsedTime(p.nowTime())

	p.lock.Lock()
	defer p.lock.Unlock()

	if p.IsClosed() {
		return false
	}

	if err := p.workers.insert(w); err != nil {
		return false
	}

	// Notify the invoker, there is an available worker in the worker queue
	p.cond.Signal()

	return true
}

// Waiting returns the number of tasks waiting to be executed.
func (p *poolCommon) Waiting() int {
	return int(atomic.LoadInt32(&p.waiting))
}

// IsClosed indicates whether the pool is closed.
func (p *poolCommon) IsClosed() bool {
	return false
}

func (p *poolCommon) addRunning(delta int) int {
	return int(atomic.AddInt32(&p.running, int32(delta)))
}

func (p *poolCommon) addWaiting(delta int) {
	atomic.AddInt32(&p.waiting, int32(delta))
}

// Cap returns the capacity of this pool.
func (p *poolCommon) Cap() int {
	return int(atomic.LoadInt32(&p.capacity))
}

// Running returns the number of workers currently running.
func (p *poolCommon) Running() int {
	return int(atomic.LoadInt32(&p.running))
}

func (p *poolCommon) nowTime() time.Time {
	return p.now.Load().(time.Time)
}
