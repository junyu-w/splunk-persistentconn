package main

import (
	"encoding/json"
	"net/http"

	conn "github.com/DrakeW/splunk-persistentconn"
)

func main() {
	server := conn.NewServer()
	server.Handle("entity/:id/data", func(req conn.Request) (conn.Response, error) {
		reqContent, err := json.Marshal(req)
		if err != nil {
			return conn.Response{}, err
		}

		return conn.Response{
			StatusCode: http.StatusOK,
			Body:       string(reqContent),
		}, nil
	}, "GET", "POST")

	server.Run()
}
