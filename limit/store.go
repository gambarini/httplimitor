package limit

import "sync"

type (
	LimitorStore interface {
		GetValue(ip Ip) ([]int64, bool)
		SetValue(ip Ip, updateFunc UpdateFunc)

	}

	UpdateFunc func(v []int64, ok bool) []int64

	MemoryStore struct {
		sMap map[Ip][]int64
		sync.RWMutex
	}
)

func NewMemoryStore() (store LimitorStore) {

	return &MemoryStore{
		sMap: map[Ip][]int64{},
	}
}

func (store *MemoryStore) GetValue(ip Ip) ([]int64, bool) {

	store.RLock()
	v, ok := store.sMap[ip]
	store.RUnlock()

	return v, ok
}

func (store *MemoryStore) SetValue(ip Ip, updateFunc UpdateFunc) {

	store.Lock()
	v, ok := store.sMap[ip]
	store.sMap[ip] = updateFunc(v, ok)
	store.Unlock()

}
