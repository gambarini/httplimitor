package limit

import "sync"

type (

	// LimitorStore
	// Request storage definition for HTTP Limitor
	// Requests are stored with Ip address and an int64 defined as unix nano seconds datetime
	LimitorStore interface {
		// GetValue Recover the requests datetime for an ip address
		// ip - ip address from the HTTP request
		// tLimit - the datetime limit in unix nano seconds to count the total number requests
		// result - array of datetime in nanoseconds for all requests after the tLimit for the ip specified
		GetValue(ip Ip, tLimit int64) (result []int64)

		// SetValue Store a request fo an ip address
		// ip - ip address from the HTTP request
		// now - time the request happened
		SetValue(ip Ip, now int64)

	}

	// MemoryStore
	// Simple implementation for the LimitorStore interface
	// uses a map to store requests
	// Performance is not assured
	// Use only for single instances
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

func (store *MemoryStore) GetValue(ip Ip, tLimit int64) (result []int64) {

	store.RLock()
	v, ok := store.sMap[ip]
	store.RUnlock()

	if !ok {
		return result
	}

	i := len(v)

	for ; i > 0; i = i - 1 {

		if v[i-1] < tLimit {
			break
		}

	}

	return v[i:]
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
