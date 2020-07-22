package internal

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

func Limit(next http.HandlerFunc, reqLimit, minutesLimit int, getIpFunc GetIpFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		clientIP := getIpFunc(r)

		isLimit, lastReqTime := IsRequestLimit(Ip(clientIP), reqLimit, minutesLimit)

		if isLimit {

			w.WriteHeader(429)

			tLeft := time.Unix(0, lastReqTime).Add(time.Minute * time.Duration(minutesLimit)).Sub(time.Now().UTC())

			_, _ = w.Write([]byte(fmt.Sprintf("Rate limit exceeded. Try again in %v.", tLeft)))

			return
		}

		done := make(chan int)

		go StoreRequest(Ip(clientIP), done)

		next.ServeHTTP(w, r)

		<-done

	}
}

func GetIP(r *http.Request) Ip {

	return Ip(strings.Split(r.RemoteAddr, ":")[0])
}
