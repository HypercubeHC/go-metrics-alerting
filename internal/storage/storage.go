package storage

import (
	"context"
	"sync"
)

type MemStorage[Name comparable, Val any] struct {
	mx      sync.RWMutex
	storage map[Name]Val
}

func New[Name comparable, Val any]() *MemStorage[Name, Val] {
	ms := MemStorage[Name, Val]{storage: make(map[Name]Val)}
	return &ms
}

func (ms *MemStorage[Name, Val]) Get(name Name) (Val, bool) {
	ms.mx.RLock()
	defer ms.mx.RUnlock()
	val, ok := ms.storage[name]
	return val, ok
}

func (ms *MemStorage[Name, Val]) Set(name Name, val Val) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	ms.storage[name] = val
}

func (ms *MemStorage[Name, Val]) Update(name Name, fn func(prev Val) Val) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	ms.storage[name] = fn(ms.storage[name])
}

func (ms *MemStorage[Name, Val]) Size() int {
	ms.mx.RLock()
	defer ms.mx.RUnlock()
	return len(ms.storage)
}

func (ms *MemStorage[Name, Val]) Keys() []Name {
	ms.mx.RLock()
	defer ms.mx.RUnlock()
	keys := make([]Name, 0, len(ms.storage))
	for name := range ms.storage {
		keys = append(keys, name)
	}
	return keys
}

func (ms *MemStorage[Name, Val]) Iter(ctx context.Context, fn func(name Name, val Val) error) error {
	ms.mx.RLock()
	defer ms.mx.RUnlock()
	for name, val := range ms.storage {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := fn(name, val); err != nil {
			return err
		}
	}
	return nil
}

func (ms *MemStorage[Name, Val]) Delete(name Name) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	delete(ms.storage, name)
}
