package main

import (
	"httplimitor"
	http2 "httplimitor/http"
	"log"
	"net/http"
)

func main() {

	lStore := http2.NewMemoryStore()

	http.HandleFunc("/", httplimitor.LimitInterceptorWithCustomLimit(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("OK"))
	}, lStore, 3, 1))

	log.Print("Listening on :8000")
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatal(err)
	}

}
