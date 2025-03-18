package ant

type Pool struct {
	*poolCommon
}

// NewPool instantiates a Pool with customized options.
func NewPool(size int, options ...Option) (*Pool, error) {
	pc, err := newPool(size, options...)
	if err != nil {
		return nil, err
	}

	pool := &Pool{poolCommon: pc}
	// New 用于在 Pool 中调用 Get() 时，如果池中没有可用的对象时，创建一个新的对象。
	pool.workerCache.New = func() any {
		return &goWorker{
			pool: pool,
		}
	}

	return pool, nil
}

func (p *Pool) Submit(task func()) error {
	if p.IsClosed() {
		return ErrPoolClosed
	}

	w, err := p.retrieveWorker()
	if err != nil {
		return err
	}

	w.inputFunc(task)
	return nil
}
