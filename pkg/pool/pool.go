package pool

import "sync"

type Reseter interface {
	Reset()
}

type Pool[T Reseter] interface {
	Get() T
	Put(t T)
}

type pool[T Reseter] struct {
	mu   *sync.Mutex
	objs []T
	new  func() T
}

func (p *pool[T]) Get() T {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.objs) > 0 {
		last := len(p.objs) - 1
		o := p.objs[last]
		p.objs = p.objs[:last]
		return o
	}

	return p.new()
}

func (p *pool[T]) Put(o T) {
	o.Reset()
	p.mu.Lock()
	p.objs = append(p.objs, o)
	p.mu.Unlock()
}

func New[T Reseter](new func() T) Pool[T] {
	return &pool[T]{new: new}
}
