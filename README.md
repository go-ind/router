# router

A tiny HTTP router with grouping and middleware support.

## Installation

```
go get github.com/go-ind/router
```

## Usage

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/go-ind/router"
)

func main() {
    r := router.SetupDefaultRouter()

    r.Get("/", func(w http.ResponseWriter, req *http.Request) {
        fmt.Fprint(w, "hello")
    })

    api := r.Group("/api")
    api.Get("/ping", func(w http.ResponseWriter, req *http.Request) {
        fmt.Fprint(w, "pong")
    })

    // simple middleware example
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("X-Demo", "middleware")
            next.ServeHTTP(w, r)
        })
    })

    http.ListenAndServe(":8080", r)
}
```

