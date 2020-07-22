# Http Limitor

HTTP requests rate limiter for Golang.

## Quick Start

HTTP limitor follows standard golang http package HandleFunc type.

```go

func main() {

    http.HandleFunc("/", httplimitor.LimitInterceptor(func(writer http.ResponseWriter, request *http.Request) {
        writer.Write([]byte("OK"))
    }))
    http.ListenAndServe(":8000", nil)                          

}

```