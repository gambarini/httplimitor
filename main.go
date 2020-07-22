package main

import (
	"fmt"
	"httplimitor/internal"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {

	http.HandleFunc("/", LimitInterceptor(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("OK"))
	}, 3, 1))

	log.Print("Listening on :8000")
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatal(err)
	}

}



func LimitInterceptor(next http.HandlerFunc, reqLimit, minutesLimit int) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.RemoteAddr == "" {
			next.ServeHTTP(w, r)
		}

		clientIP := strings.Split(r.RemoteAddr, ":")[0]

		isLimit, lastReqTime := internal.IsRequestLimit(internal.Ip(clientIP), reqLimit, minutesLimit)

		if  isLimit {

			log.Printf("Request limit reached for %s.", clientIP)

			w.WriteHeader(429)

			tLeft := time.Unix(0, lastReqTime).Add(time.Minute * time.Duration(minutesLimit)).Sub(time.Now().UTC())

			_, _ = w.Write([]byte(fmt.Sprintf("Rate limit exceeded. Try again in %v.", tLeft)))


			return
		}

		log.Printf("Storing request for %s", clientIP)

		done := make(chan int)

		go  internal.StoreRequest(internal.Ip(clientIP), done)

		next.ServeHTTP(w, r)

		<- done

	}
}


