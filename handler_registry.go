package persistentconn

import (
	"net/http"
	"regexp"
)

// Handler is a handler function that takes a persistentconn request and returns a response or error
type Handler func(Request) (Response, error)

// NoMatchingHandler is the default handler returned when no matching path
// is found from a request
func NoMatchingHandler(re Request) (Response, error) {
	return Response{
		StatusCode: http.StatusNotFound,
		Body:       "The requested path is not found.",
	}, nil
}

type route struct {
	Pattern *regexp.Regexp
	Handler Handler
	Methods []string
}

func newRoute(pathPattern string, handler Handler, allowedMethods []string) *route {
	re := regexp.MustCompile(pathPattern)
	return &route{
		Pattern: re,
		Handler: handler,
		Methods: allowedMethods,
	}
}

func (r *route) findMatch(path string) Handler {
	// TODO: implement this
	return nil
}

type handlerRegistry struct {
	routes []*route
}

func (rg *handlerRegistry) getHandler(path string) Handler {
	handler := NoMatchingHandler
	for _, route := range rg.routes {
		handler = route.findMatch(path)
		if handler != nil {
			break
		}
	}
	return handler
	// return func(req Request) (Response, error) {
	// 	sleepTime := rand.Intn(1000)
	// 	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	// 	return Response{
	// 		StatusCode: 200,
	// 		Body:       fmt.Sprintf("hello world %s", req.OutputMode),
	// 	}, nil
	// }
}

func (rg *handlerRegistry) register(path string, handler Handler, allowedMethods []string) {
	route := newRoute(path, handler, allowedMethods)
	rg.routes = append(rg.routes, route)
}
