# trace [![Build Status](https://travis-ci.org/vinxi/trace.png)](https://travis-ci.org/vinxi/trace) [![GoDoc](https://godoc.org/github.com/vinxi/trace?status.svg)](https://godoc.org/github.com/vinxi/trace) [![Coverage Status](https://coveralls.io/repos/github/vinxi/trace/badge.svg?branch=master)](https://coveralls.io/github/vinxi/trace?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/vinxi/trace)](https://goreportcard.com/report/github.com/vinxi/trace)

Traffic tracing instrumentation for your proxies. 
Designed to be extended to trace custom data or modify the request/response. 

Relies on [log](https://github.com/vinxi/log) package to write structured traces and optionally send them via [hooks](https://github.com/Sirupsen/logrus#hooks) to different storage services.

## Installation

```bash
go get -u gopkg.in/vinxi/trace.v0
```

## API

See [godoc](https://godoc.org/github.com/vinxi/trace) reference.

## Example

#### Default tracing

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

  // Instrument the proxy with trace middleware 
  // Now all the incoming traffic will be registered. 
  vs.Use(trace.Default)

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
