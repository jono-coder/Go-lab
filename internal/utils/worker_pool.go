package utils

import (
	"context"
	"sync"
)

type TaskWithError func(ctx context.Context) error

type WorkerPool struct {
	tasks     chan TaskWithError
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	errOnce   sync.Once
	err       error
	closeOnce sync.Once
}

func NewWorkerPool(numWorkers, bufferSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	p := &WorkerPool{
		tasks:  make(chan TaskWithError, bufferSize),
		ctx:    ctx,
		cancel: cancel,
	}

	p.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer p.wg.Done()

			for task := range p.tasks {
				if p.ctx.Err() != nil {
					return
				}

				if err := task(p.ctx); err != nil {
					p.errOnce.Do(func() {
						p.err = err
						p.cancel() // cancel the pool
					})
				}
			}
		}()
	}

	return p
}

func (p *WorkerPool) Submit(task TaskWithError) error {
	// If the pool is already cancelled, stop
	if p.ctx.Err() != nil {
		return p.ctx.Err()
	}

	// Try to send task unless cancelled mid-flight
	select {
	case <-p.ctx.Done():
		return p.ctx.Err()
	case p.tasks <- task:
		return nil
	}
}

func (p *WorkerPool) Wait() error {
	p.closeOnce.Do(func() {
		close(p.tasks)
	})

	p.wg.Wait()
	return p.err
}

func (p *WorkerPool) Context() context.Context {
	return p.ctx
}
