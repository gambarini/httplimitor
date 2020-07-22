package main

import (
	"httplimitor"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", httplimitor.LimitInterceptorWithCustomLimit(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("OK"))
	}, 3, 1))

	log.Print("Listening on :8000")
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatal(err)
	}

}
