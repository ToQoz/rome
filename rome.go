package rome

import (
	"net/http"
	"strings"
)

type Router struct {
	tries           map[string]*trie
	notFoundHandler notFoundHandler
}

func NewRouter() *Router {
	return &Router{tries: map[string]*trie{}, notFoundHandler: notFoundHandler{}}
}

func (router *Router) Get(pattern string, handler http.Handler) {
	router.handle("GET", pattern, handler)

	router.handle("HEAD", pattern, handler)
}

func (router *Router) Head(pattern string, handler http.Handler) {
	router.handle("HEAD", pattern, handler)
}

func (router *Router) Post(pattern string, handler http.Handler) {
	router.handle("POST", pattern, handler)
}

func (router *Router) Put(pattern string, handler http.Handler) {
	router.handle("PUT", pattern, handler)
}

func (router *Router) Delete(pattern string, handler http.Handler) {
	router.handle("DELETE", pattern, handler)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var err error
	var r *routeWithParam

	if router.tries[req.Method] != nil {
		r, err = router.tries[req.Method].get(req.URL.Path)

		if err == nil {
			r.serveHTTP(w, req)
			return
		}
	}

	router.notFoundHandler.ServeHTTP(w, req)
}

func (router *Router) handle(_method string, pattern string, handler http.Handler) {
	method := strings.ToUpper(_method)

	if router.tries[method] == nil {
		router.tries[method] = newTrie()
	}

	router.tries[method].add(newRoute(pattern, handler))
}
