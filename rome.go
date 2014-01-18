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
	router := &Router{tries: map[string]*trie{}}
	router.tries["GET"] = newTrie()
	router.tries["HEAD"] = newTrie()
	router.tries["POST"] = newTrie()
	router.tries["PUT"] = newTrie()
	router.tries["DELETE"] = newTrie()

	return router
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
}

func (router *Router) GetFunc(pattern string, f func(http.ResponseWriter, *http.Request)) {
	router.Get(pattern, http.HandlerFunc(f))
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
	method := strings.ToUpper(req.Method)

	if method == "" {
		method = "GET"
	}

	r, err := router.tries[method].get(req.URL.Path)

	// if GET method respond, HEAD method should respond too.
	if err != nil && method == "HEAD" {
		r, err = router.tries["GET"].get(req.URL.Path)
	}

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
	router.tries[strings.ToUpper(method)].add(newRoute(pattern, handler))
}
