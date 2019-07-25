package persistentconn

type handlerRegistry struct {
}

// Handler is a handler function that takes a persistentconn request and returns a response or error
type Handler func(Request) (Response, error)

func (rg *handlerRegistry) getHandler(path string) Handler {
	// TODO: implement this
	return func(req Request) (Response, error) {
		return Response{
			statusCode: 200,
			body:       "hello world",
		}, nil
	}
}
