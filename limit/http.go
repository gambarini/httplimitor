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

	store.SetValue(ip, timeNow().UnixNano())

	done <- 0

}

func IsRequestLimit(store LimitorStore, ip Ip, limit int, minutes int) (isLimit bool, lastTimestamp int64) {

	tLimit := timeNow().Add(time.Minute * time.Duration(-minutes)).UTC().UnixNano()

	result := store.GetValue(ip ,tLimit)

	if len(result) == 0 {
		return false, tLimit
	}

	return len(result) >= limit, result[len(result)-1]

}
