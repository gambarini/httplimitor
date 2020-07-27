package main

import (
	"httplimitor"
	"httplimitor/limit"
	"log"
	"net/http"
)

func main() {

	lStore := limit.NewMemoryStore()

	http.HandleFunc("/", httplimitor.LimitInterceptorWithCustomLimit(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("OK"))
	}, lStore, 3, 10))

	log.Print("Listening on :8000")
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatal(err)
	}

}
