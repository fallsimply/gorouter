package gorouter

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

//Router is a router based off ServeMux
type Router struct {
	http.ServeMux
	Err    http.HandlerFunc
	Routes []string
}

//New is a returns a new router instance
func New(errHandlers ...http.HandlerFunc) (rtr *Router) {
	rtr = new(Router)
	if len(errHandlers) > 0 {
		rtr.Err = errHandlers[0]
	} else {
		rtr.Err = http.NotFound
	}
	return rtr
}

func (rtr *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rtr.ServeMux.ServeHTTP(rw, req)
}

//Add adds a new route to the router
func (rtr *Router) Add(pattern string, handler http.Handler, serveIndex bool, index ...http.Handler) {
	var routerRegex = regexp.MustCompile(`([:](?P<name>.+))+`)
	var prefix string = "/"
	var params []string = []string{"", "", ""}
	var indexHandler http.Handler = nil

	if routerRegex.MatchString(pattern) {
		prefix = routerRegex.ReplaceAllString(pattern, "")
		params = routerRegex.FindStringSubmatch(pattern)
	}
	fmt.Println(params)
	if len(index) == 1 {
		indexHandler = index[0]
	}

	rtr.ServeMux.HandleFunc(prefix, rtr.routeMiddleware(prefix, pattern, map[string]http.Handler{
		"main":  handler,
		"index": indexHandler,
	}, params, rtr.Err, serveIndex))
}

func (rtr *Router) routeMiddleware(prefix string, pattern string, handlers map[string]http.Handler, params []string, errorHandler http.HandlerFunc, serveIndex bool) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		paramVal := strings.Replace(req.URL.Path, pattern, "", -1)
		fmt.Printf("[gorouter/route] pattern: \"%s\", url: \"%s\", parameter: { name: \"%s\", value: \"%s\" }\n", pattern, req.URL.String(), params[2], paramVal)
		if len(params) != 0 {
			vals, _ := url.ParseQuery(req.URL.RawQuery)
			vals.Add(params[0], strings.ReplaceAll(req.URL.Path, prefix, ""))
			req.URL.RawQuery = vals.Encode()
		}

		switch {
		case (paramVal == "") && serveIndex && (handlers["index"] != nil):
			handlers["index"].ServeHTTP(rw, req)
		case (paramVal != "") && (params[2] != ""):
			handlers["main"].ServeHTTP(rw, req)
		default:
			errorHandler.ServeHTTP(rw, req)
		}
	}
}
