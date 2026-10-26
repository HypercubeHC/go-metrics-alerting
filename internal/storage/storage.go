package storage

type MemStorage[Name comparable, Val any] struct {
	storage map[Name]Val
}

func New[Name comparable, Val any]() *MemStorage[Name, Val] {
	ms := MemStorage[Name, Val]{storage: make(map[Name]Val)}
	return &ms
}

func (ms *MemStorage[Name, Val]) Get(name Name) (Val, bool) {
	var val Val
	if val, ok := ms.storage[name]; ok {
		return val, true
	}
	return val, false
}

func (ms *MemStorage[Name, Val]) Set(name Name, val Val) {
	ms.storage[name] = val
}
