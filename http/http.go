package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type (
	Ip string
	GetIpFunc func(r *http.Request) Ip
)

func Limit(next http.HandlerFunc, reqLimit, minutesLimit int, getIpFunc GetIpFunc, store LimitorStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		clientIP := getIpFunc(r)

		isLimit, lastReqTime := IsRequestLimit(store, Ip(clientIP), reqLimit, minutesLimit)

		if isLimit {

			w.WriteHeader(429)

			tLeft := time.Unix(0, lastReqTime).Add(time.Minute * time.Duration(minutesLimit)).Sub(time.Now().UTC())

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

	now := time.Now().UTC().UnixNano()

	ts, ok := store.GetValue(ip)

	if !ok {
		store.SetValue(ip, []int64{now})
	} else {
		store.SetValue(ip, append(ts, now))
	}

	done <- 0
}

func IsRequestLimit(store LimitorStore, ip Ip, limit int, minutes int) (isLimit bool, lastTimestamp int64) {

	tLimit := time.Now().Add(time.Minute * time.Duration(-minutes)).UTC().UnixNano()

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
