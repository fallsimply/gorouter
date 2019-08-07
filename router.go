package gorouter

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

//Router is a router based off ServeMux
type Router struct {
	http.ServeMux
	Err http.HandlerFunc
}

//New is a returns a new router instance
func New(errHandler ...http.HandlerFunc) *Router {
	if errHandler[1] == nil {
		errHandler[1] = http.NotFound
	}
	var rtr = new(Router)
	rtr.Err = errHandler[1]
	return rtr

}

//Add adds a new route to the router
func (rtr *Router) Add(pattern string, handler http.Handler, specific bool) {
	var routerRegex = regexp.MustCompile(`([:](?P<name>.+))+`)
	var prefix = routerRegex.ReplaceAllString(pattern, "")
	var pathName = routerRegex.FindStringSubmatch(pattern)
	rtr.ServeMux.Handle(prefix, routeMiddleware(prefix, pattern, handler, pathName, rtr.Err, specific))
}

func (rtr *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rtr.ServeMux.ServeHTTP(rw, req)
}

func routeMiddleware(prefix string, pattern string, handler http.Handler, pathName []string, errorHandler http.HandlerFunc, specific bool) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if specific && (req.URL.Path == prefix) {
			errorHandler(rw, req)
			return
		}
		vals, _ := url.ParseQuery(req.URL.RawQuery)
		vals.Add(pathName[0], strings.ReplaceAll(req.URL.Path, prefix, ""))
		req.URL.RawQuery = vals.Encode()
		handler.ServeHTTP(rw, req)
	}
}
