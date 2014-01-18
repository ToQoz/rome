package rome

import (
	"net/http"
	"strings"
)

type Router struct {
	tries           map[string]*trie
	notFoundHandler http.Handler
}

func NewRouter() *Router {
	return &Router{tries: map[string]*trie{}}
}

// --- Routing Helper ---

// Not Found
func (router *Router) NotFound(handler http.Handler) {
	router.notFoundHandler = handler
}

func (router *Router) NotFoundFunc(handlerFunc func(http.ResponseWriter, *http.Request)) {
	router.NotFound(http.HandlerFunc(handlerFunc))
}

// GET
func (router *Router) Get(pattern string, handler http.Handler) {
	router.handle("GET", pattern, handler)
	router.handle("HEAD", pattern, handler)
}

func (router *Router) GetFunc(pattern string, f func(http.ResponseWriter, *http.Request)) {
	router.Get(pattern, http.HandlerFunc(f))
	router.Head(pattern, http.HandlerFunc(f))
}

// HEAD
func (router *Router) Head(pattern string, handler http.Handler) {
	router.handle("HEAD", pattern, handler)
}

func (router *Router) HeadFunc(pattern string, f func(http.ResponseWriter, *http.Request)) {
	router.Head(pattern, http.HandlerFunc(f))
}

// POST
func (router *Router) Post(pattern string, handler http.Handler) {
	router.handle("POST", pattern, handler)
}

func (router *Router) PostFunc(pattern string, f func(http.ResponseWriter, *http.Request)) {
	router.Post(pattern, http.HandlerFunc(f))
}

// PUT
func (router *Router) Put(pattern string, handler http.Handler) {
	router.handle("PUT", pattern, handler)
}

func (router *Router) PutFunc(pattern string, f func(http.ResponseWriter, *http.Request)) {
	router.Put(pattern, http.HandlerFunc(f))
}

// DELETE
func (router *Router) Delete(pattern string, handler http.Handler) {
	router.handle("DELETE", pattern, handler)
}

func (router *Router) DeleteFunc(pattern string, f func(http.ResponseWriter, *http.Request)) {
	router.Delete(pattern, http.HandlerFunc(f))
}

// --- rome.Router work as http.Handler ---

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var err error
	var r *routeWithParam

	method := strings.ToUpper(req.Method)

	if method == "" {
		method = "GET"
	}

	r, err = router.tries[method].get(req.URL.Path)

	if err != nil {
		if router.notFoundHandler != nil {
			router.notFoundHandler.ServeHTTP(w, req)
			return
		}

		http.NotFoundHandler().ServeHTTP(w, req)
		return
	}

	r.serveHTTP(w, req)
}

// --- Private methods ---

func (router *Router) handle(method string, pattern string, handler http.Handler) {
	method = strings.ToUpper(method)

	if router.tries[method] == nil {
		router.tries[method] = newTrie()
	}

	router.tries[method].add(newRoute(pattern, handler))
}
