package main

import (
	"httplimitor/internal"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", LimitInterceptor(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("OK"))
	}, 10, 1))

	log.Print("Listening on :8000")
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatal(err)
	}

}



func LimitInterceptor(next http.HandlerFunc, reqLimit, minutesLimit int) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if internal.IsRequestLimit(internal.Ip(r.URL.Host), reqLimit, minutesLimit) {

			log.Print("Request limit reached.")

			w.WriteHeader(429)
			_, _ = w.Write([]byte("Rate limit exceeded."))
			return
		}

		log.Print("Storing request")

		done := make(chan int)

		go  internal.StoreRequest(internal.Ip(r.URL.Host), done)

		next.ServeHTTP(w, r)

		<- done

	}
}
