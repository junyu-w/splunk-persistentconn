package persistentconn

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Handler is a handler function that takes a persistentconn request and returns a response or error
// if error is returned, the server will return a 500 (Internal Server Error) response with the returned
// error's message as the response body
type Handler func(Request) (Response, error)

// NoMatchingHandler is the default handler returned when no matching path
// is found from a request
func NoMatchingHandler(re Request) (Response, error) {
	return Response{
		StatusCode: http.StatusNotFound,
		Body:       "The requested path is not found.",
	}, nil
}

// route represents a registered route that has a corresponding handler
type route struct {
	Pattern *regexp.Regexp
	Handler Handler
	Methods []string
}

// newRoute creates a new route object
func newRoute(pathPattern string, handler Handler, allowedMethods []string) *route {
	re := translatePatternToRegexp(pathPattern)
	return &route{
		Pattern: re,
		Handler: handler,
		Methods: allowedMethods,
	}
}

// translatePatternToRegexp translates a path pattern to a regexp
func translatePatternToRegexp(pathPattern string) *regexp.Regexp {
	parts := strings.Split(pathPattern, "/")
	regexpStrParts := make([]string, len(parts))
	for idx, p := range parts {
		if strings.HasPrefix(p, ":") {
			p = fmt.Sprintf(`(?P<%s>[\S|^\/]+)`, p[1:])
		}
		regexpStrParts[idx] = p
	}
	regexpStr := strings.Join(regexpStrParts, "/")
	re := regexp.MustCompile(regexpStr)
	return re
}

// handlerRegistry is where all routes are stored
type handlerRegistry struct {
	routes []*route
}

// gethandler gets the handler based on the input reqeust's path info
func (rg *handlerRegistry) getHandler(req Request) Handler {
	handler := NoMatchingHandler
	for _, rt := range rg.routes {
		if matches := rt.Pattern.FindStringSubmatch(req.Path); len(matches) > 0 && contains(rt.Methods, req.Method) {
			matchGroupNames := rt.Pattern.SubexpNames()
			for idx, name := range matchGroupNames {
				// Since the Regexp as a whole cannot be named, first matched name is always the empty string
				if name != "" {
					req.Params[name] = matches[idx]
				}
			}
			return rt.Handler
		}
	}
	return handler
}

// register func registers a path with a handler
func (rg *handlerRegistry) register(path string, handler Handler, allowedMethods []string) {
	route := newRoute(path, handler, allowedMethods)
	rg.routes = append(rg.routes, route)
}
