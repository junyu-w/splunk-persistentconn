package persistentconn

import "time"

type handlerRegistry struct {
}

// Handler is a handler function that takes a persistentconn request and returns a response or error
type Handler func(Request) (Response, error)

func (rg *handlerRegistry) getHandler(path string) Handler {
	// TODO: implement this
	return func(req Request) (Response, error) {
		time.Sleep(1 * time.Second)
		return Response{
			StatusCode: 200,
			Body:       "hello world",
		}, nil
	}
}
