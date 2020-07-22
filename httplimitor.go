package httplimitor

import (
	"httplimitor/internal"
	"net/http"
)

func LimitInterceptor(next http.HandlerFunc) http.HandlerFunc {

	return internal.Limit(next, 100, 60, internal.GetIP, internal.NewMemoryStore())
}

func LimitInterceptorWithCustomLimit(next http.HandlerFunc, reqLimit, minutesLimit int) http.HandlerFunc {

	return internal.Limit(next, reqLimit, minutesLimit, internal.GetIP, internal.NewMemoryStore())
}

func LimitInterceptorWithCustomIp(next http.HandlerFunc, reqLimit, minutesLimit int, getIpFunc internal.GetIpFunc) http.HandlerFunc {

	return internal.Limit(next, reqLimit, minutesLimit, getIpFunc, internal.NewMemoryStore())
}

func LimitInterceptorWithCustomStore(next http.HandlerFunc, reqLimit, minutesLimit int, getIpFunc internal.GetIpFunc, store internal.HTTPLimitorStore) http.HandlerFunc {

	return internal.Limit(next, reqLimit, minutesLimit, getIpFunc, store)
}
