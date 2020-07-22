package httplimitor

import (
	"httplimitor/internal"
	"net/http"
)



func LimitInterceptor(next http.HandlerFunc) http.HandlerFunc {

	return internal.Limit(next, 100, 60, internal.GetIP)
}

func LimitInterceptorWithCustomLimit(next http.HandlerFunc, reqLimit, minutesLimit int) http.HandlerFunc {

	return internal.Limit(next, reqLimit, minutesLimit, internal.GetIP)
}

func LimitInterceptorWithCustomIp(next http.HandlerFunc, reqLimit, minutesLimit int, getIpFunc internal.GetIpFunc) http.HandlerFunc {

	return internal.Limit(next, reqLimit, minutesLimit, getIpFunc)
}
