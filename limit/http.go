package limit

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type (
	Ip        string
	GetIpFunc func(r *http.Request) Ip
)

var (
	timeNow = func() time.Time {
		return time.Now().UTC()
	}
)

func Limit(next http.HandlerFunc, reqLimit, minutesLimit int, getIpFunc GetIpFunc, store LimitorStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		clientIP := getIpFunc(r)

		isLimit, lastReqTime := IsRequestLimit(store, Ip(clientIP), reqLimit, minutesLimit)

		if isLimit {

			w.WriteHeader(429)

			tLeft := time.Unix(0, lastReqTime).Add(time.Minute * time.Duration(minutesLimit)).Sub(timeNow())

			_, _ = w.Write([]byte(fmt.Sprintf("Rate limit exceeded. Try again in %v.", tLeft)))

			return
		}

		done := make(chan int)

		go SaveRequest(store, Ip(clientIP), done)

		next.ServeHTTP(w, r)

		<-done

	}
}

func GetIP(r *http.Request) Ip {

	return Ip(strings.Split(r.RemoteAddr, ":")[0])
}

func SaveRequest(store LimitorStore, ip Ip, done chan int) {

	store.SetValue(ip, func(v []int64, ok bool) []int64 {
		now := timeNow().UnixNano()

		if !ok {
			v = []int64{now}
		} else {
			v = append(v, now)
		}

		return v
	})

	done <- 0

}

func IsRequestLimit(store LimitorStore, ip Ip, limit int, minutes int) (isLimit bool, lastTimestamp int64) {

	tLimit := timeNow().Add(time.Minute * time.Duration(-minutes)).UTC().UnixNano()

	ts, ok := store.GetValue(ip)

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
