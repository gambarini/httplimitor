package limit

import "sync"

type (
	LimitorStore interface {
		GetValue(ip Ip) ([]int64, bool)
		SetValue(ip Ip, now int64)

	}

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

func (store *MemoryStore) SetValue(ip Ip, now int64) {

	store.Lock()

	v, ok := store.sMap[ip]

	if !ok {
		v = []int64{now}
	} else {
		v = append(v, now)
	}

	store.sMap[ip] = v

	store.Unlock()

}
