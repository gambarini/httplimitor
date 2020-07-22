package internal

import (
	"log"
	"time"
)

type (
	Ip string
)

var (
	requestStore = map[Ip][]int64{}
)

func StoreRequest(ip Ip, done chan int) {

	now := time.Now().UTC().UnixNano()

	ts, ok := requestStore[ip]

	if !ok {

		requestStore[ip] = []int64{now}
	} else {

		requestStore[ip] = append(ts, now)
	}

	done <- 0
}

func IsRequestLimit(ip Ip, limit int, minutes int) bool {

	tLimit := time.Now().Add(time.Minute * time.Duration(-minutes)).UTC().UnixNano()

	ts, ok := requestStore[ip]

	log.Printf("Req: %d", len(ts))
	log.Printf("ts: %v", ts)
	log.Printf("ts len: %d", len(ts))

	if !ok {
		return false
	}

	c := 0

	log.Printf("c: %d", c)

	for i := len(ts); i > 0; i = i-1 {

		log.Printf("ts %d - tLimit: %d", ts[i-1], tLimit)

		if ts[i-1] < tLimit {
			break
		}

		c++
	}

	log.Printf("c: %d", c)

	return c >= limit

}


