package httplimitor

import (
	limit "httplimitor/limit"
	"net/http"
)

func LimitInterceptor(next http.HandlerFunc, lStore limit.LimitorStore) http.HandlerFunc {

	return limit.Limit(next, 100, 60, limit.GetIP, lStore)
}

func LimitInterceptorWithCustomLimit(next http.HandlerFunc, lStore limit.LimitorStore, reqLimit, minutesLimit int) http.HandlerFunc {

	return limit.Limit(next, reqLimit, minutesLimit, limit.GetIP, lStore)
}

func LimitInterceptorWithCustomIp(next http.HandlerFunc, lStore limit.LimitorStore, reqLimit, minutesLimit int, getIpFunc limit.GetIpFunc) http.HandlerFunc {

	return limit.Limit(next, reqLimit, minutesLimit, getIpFunc, lStore)
}
