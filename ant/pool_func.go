package ant

type PoolWithFunc struct {
	*poolCommon

	// fn is the unified function for processing tasks.
	fn func(any)
}

func NewPoolWithFunc(size int, pf func(any), options ...Option) (*PoolWithFunc, error) {
	if pf == nil {
		return nil, ErrLackPoolFunc
	}

	pc, err := newPool(size, options...)
	if err != nil {
		return nil, err
	}

	pool := &PoolWithFunc{
		poolCommon: pc,
		fn:         pf,
	}

	pool.workerCache.New = func() any {
		return &goWorkerWithFunc{
			pool: pool,
			arg:  make(chan any, workerChanCap),
		}
	}

	return pool, nil
}

// Invoke passes arguments to the pool.
func (p *PoolWithFunc) Invoke(arg any) error {
	if p.IsClosed() {
		return ErrPoolClosed
	}

	w, err := p.retrieveWorker()
	if err != nil {
		return err
	}

	w.inputArg(arg)

	return nil
}
