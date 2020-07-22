package internal

type (
	HTTPLimitorStore interface {
		GetValue(ip Ip) ([]int64, bool)
		SetValue(ip Ip, v []int64)
	}

	MemoryStore struct {
		store map[Ip][]int64
	}
)

func NewMemoryStore() (store HTTPLimitorStore){

	store = &MemoryStore {
		store: map[Ip][]int64{},
	}

	return store
}

func (store *MemoryStore) GetValue(ip Ip) ([]int64, bool) {

	v, ok := store.store[ip]

	return v, ok
}

func (store *MemoryStore) SetValue(ip Ip, v []int64) {
	store.store[ip] = v
}


