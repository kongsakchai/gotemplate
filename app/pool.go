package app

import "sync"

type pool[T any] struct {
	*sync.Pool
}

func createPool[T any]() *pool[T] {
	return &pool[T]{&sync.Pool{
		New: func() any {
			return new(T)
		},
	}}
}

func (p *pool[T]) Get() *T {
	return p.Pool.Get().(*T)
}

func (p *pool[T]) Put(x *T) {
	p.Pool.Put(x)
}
