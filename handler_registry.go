package persistentconn

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
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
	re := translatePatternToRegexp(pathPattern)
	return &route{
		Pattern: re,
		Handler: handler,
		Methods: allowedMethods,
	}
}

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

type handlerRegistry struct {
	routes []*route
}

func (rg *handlerRegistry) getHandler(req Request) Handler {
	handler := NoMatchingHandler
	for _, rt := range rg.routes {
		if matches := rt.Pattern.FindStringSubmatch(req.Path); len(matches) > 0 && contains(rt.Methods, req.Method) {
			return rt.Handler
		}
	}
	return handler
}

func (rg *handlerRegistry) register(path string, handler Handler, allowedMethods []string) {
	route := newRoute(path, handler, allowedMethods)
	rg.routes = append(rg.routes, route)
}
