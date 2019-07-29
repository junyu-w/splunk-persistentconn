# splunk-persistentconn

[![GoDoc](https://godoc.org/github.com/DrakeW/splunk-persistentconn?status.svg)](https://godoc.org/github.com/DrakeW/splunk-persistentconn)

splunk-persistentconn implements the protocol that splunk core uses to communicate with persistent script that's used to power REST endpoints of an app

## How to use

Usage example:

```go
package main

import (
    "fmt"
    "net/http"
    "time"

    conn "github.com/DrakeW/splunk-persistentconn"
)

func main() {
    // server initialization
    server := conn.NewServer()

    // register endpoint using either static path or path pattern
    server.Handle("entity/:id/data", func(req conn.Request) (conn.Response, error) {
        return conn.Response{
            StatusCode: http.StatusOK,
            Body:       fmt.Sprintf("hello world %s", req.Params["id"]),
        }, nil
    }, "GET", "POST")

    // starts server
    server.Run()
}
```

## Performance

This package is built to take advantage of the concurrent processing capability provided by goroutines for each handler registered. However, since Splunk core sends requests to the server process in a **synchronous** manner, only one goroutine is spawned at a time therefore currently there's no performance boost.

Based on rough benchmark with the same 5000 non-blocking requests (where each request takes ~0.5s to finish) hitting both the python REST endpoint based on `PersistentServerConnectionApplication` and the Go REST endpoint using the `splunk-persistentconn` package, performance is roughly the same.

## Disclaimer

This package is not officially supported by Splunk. Please use with caution and feel free to open issue or PR for any kind of bug report or feature enhancement.
