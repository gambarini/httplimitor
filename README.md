# HTTP Limitor

HTTP requests rate limiter for Golang.

## Quick Start

```go

func main() {

    lStore := limit.NewMemoryStore()

    http.HandleFunc("/", httplimitor.LimitInterceptor(func(writer http.ResponseWriter, request *http.Request) {
        writer.Write([]byte("OK"))
    }, lStore))

    http.ListenAndServe(":8000", nil)                          

}

```




