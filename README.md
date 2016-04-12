# trace [![Build Status](https://travis-ci.org/vinxi/trace.png)](https://travis-ci.org/vinxi/trace) [![GoDoc](https://godoc.org/github.com/vinxi/trace?status.svg)](https://godoc.org/github.com/vinxi/trace) [![Coverage Status](https://coveralls.io/repos/github/vinxi/trace/badge.svg?branch=master)](https://coveralls.io/github/vinxi/trace?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/vinxi/trace)](https://goreportcard.com/report/github.com/vinxi/trace)

HTTP traffic tracing instrumnetation for your proxies. 

## Installation

```bash
go get -u gopkg.in/vinxi/trace.v0
```

## API

See [godoc](https://godoc.org/github.com/vinxi/trace) reference.

## Example

#### Basic traffic tracing instrumentation

```go
package main

import (
  "fmt"
  "net/http"

  "gopkg.in/vinxi/trace.v0"
  "gopkg.in/vinxi/vinxi.v0"
)

const port = 3100

func main() {
  // Create a new vinxi proxy
  vs := vinxi.NewServer(vinxi.ServerOptions{Port: port})

  // Plugin multiple middlewares writting some logs
  vs.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
    log.Info("[%s] %s", r.Method, r.RequestURI)
    h.ServeHTTP(w, r)
  })

  vs.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
    log.Warnf("%s", "foo bar")
    h.ServeHTTP(w, r)
  })

  // Target server to forward
  vs.Forward("http://httpbin.org")

  fmt.Printf("Server listening on port: %d\n", port)
  err := vs.Listen()
  if err != nil {
    fmt.Errorf("Error: %s\n", err)
  }
}

```

## License

MIT
