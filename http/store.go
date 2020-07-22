package http

import "sync"

type (
	LimitorStore interface {
		GetValue(ip Ip) ([]int64, bool)
		SetValue(ip Ip, v []int64)
	}

	MemoryStore struct {
		mapStore sync.Map
	}
)

func NewMemoryStore() (store LimitorStore) {

	return &MemoryStore{}
}

func (store *MemoryStore) GetValue(ip Ip) ([]int64, bool) {

	v, ok := store.mapStore.Load(ip)

	if ok {
		return v.([]int64), ok
	} else {
		return []int64{}, ok
	}
}

func (store *MemoryStore) SetValue(ip Ip, v []int64) {
	store.mapStore.Store(ip, v)
}
