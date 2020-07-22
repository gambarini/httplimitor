package internal

import (
	"time"
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

func IsRequestLimit(ip Ip, limit int, minutes int) (isLimit bool, lastTimestamp int64) {

	tLimit := time.Now().Add(time.Minute * time.Duration(-minutes)).UTC().UnixNano()

	ts, ok := requestStore[ip]

	if !ok {
		return false, tLimit
	}

	c := 0

	for i := len(ts); i > 0; i = i - 1 {

		if ts[i-1] < tLimit {
			break
		}

		c++
	}

	return c >= limit, ts[len(ts)-1]

}
