package persistentconn

import (
	"fmt"
	"math/rand"
	"time"
)

type handlerRegistry struct {
}

// Handler is a handler function that takes a persistentconn request and returns a response or error
type Handler func(Request) (Response, error)

func (rg *handlerRegistry) getHandler(path string) Handler {
	// TODO: implement this
	return func(req Request) (Response, error) {
		sleepTime := rand.Intn(1000)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		return Response{
			StatusCode: 200,
			Body:       fmt.Sprintf("hello world %s", req.OutputMode),
		}, nil
	}
}
