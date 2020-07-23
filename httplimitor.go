package httplimitor

import (
	http2 "httplimitor/limit"
	"net/http"
)

func LimitInterceptor(next http.HandlerFunc, lStore http2.LimitorStore) http.HandlerFunc {

	return http2.Limit(next, 100, 60, http2.GetIP, lStore)
}

func LimitInterceptorWithCustomLimit(next http.HandlerFunc, lStore http2.LimitorStore, reqLimit, minutesLimit int) http.HandlerFunc {

	return http2.Limit(next, reqLimit, minutesLimit, http2.GetIP, lStore)
}

func LimitInterceptorWithCustomIp(next http.HandlerFunc, lStore http2.LimitorStore, reqLimit, minutesLimit int, getIpFunc http2.GetIpFunc) http.HandlerFunc {

	return http2.Limit(next, reqLimit, minutesLimit, getIpFunc, lStore)
}
