package main

import (
	"fmt"
	"httplimitor/internal"
	"log"
	"net/http"
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

		isLimit, lastReqTime := internal.IsRequestLimit(internal.Ip(r.URL.Host), reqLimit, minutesLimit)

		if  isLimit {

			log.Print("Request limit reached.")

			w.WriteHeader(429)

			tLeft := time.Unix(0, lastReqTime).Add(time.Minute * time.Duration(minutesLimit)).Sub(time.Now().UTC())

			_, _ = w.Write([]byte(fmt.Sprintf("Rate limit exceeded. Try again in %v.", tLeft)))


			return
		}

		log.Print("Storing request")

		done := make(chan int)

		go  internal.StoreRequest(internal.Ip(r.URL.Host), done)

		next.ServeHTTP(w, r)

		<- done

	}
}
