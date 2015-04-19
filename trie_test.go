package rome

import (
	"net/http"
	"testing"
)

var (
	tr      = newTrie()
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
)

func init() {
	tr.add(newRoute("/users", handler))
	tr.add(newRoute("/users/:id", handler))
	tr.add(newRoute("/users/:user_id/posts", handler))
	tr.add(newRoute("/users/:user_id/posts/:id", handler))
	tr.add(newRoute("/users/:user_id/posts/latest", handler))
	tr.add(newRoute("/users/:user_id/:content_kind", handler))
	tr.add(newRoute("/users/:user_id/:content_kind/:id", handler))
	tr.add(newRoute("/users/:user_id/:content_kind/latest", handler))
	tr.add(newRoute("/users/:user_id/parts", handler))
	tr.add(newRoute("/users/:user_id/parts/:id", handler))
	tr.add(newRoute("/users/:user_id/parts/latest", handler))
	tr.add(newRoute("/users/login/posts", handler))
	tr.add(newRoute("/users/login/posts/:id", handler))
	tr.add(newRoute("/x/*/y/*", handler))
	tr.add(newRoute("/x/*/users/:id", handler))
	tr.add(newRoute("/assets/*", handler))
}

func TestTrie(t *testing.T) {
	var err error
	var m *routeWithParam

	_, err = tr.get("/NOT_FOUND")
	neq(t, nil, err)

	m, _ = tr.get("/users")
	eq(t, m.route.pattern, "/users")

	m, _ = tr.get("/users/")
	eq(t, m.route.pattern, "/users")

	m, _ = tr.get("/users/1")
	m.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := PathParams(r)
		eq(t, params.Value("id"), "1")
	})
	eq(t, m.route.pattern, "/users/:id")
	testServe(m)

	m, _ = tr.get("/users//2")
	m.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := PathParams(r)
		eq(t, params.Value("id"), "2")
	})
	eq(t, m.route.pattern, "/users/:id")
	testServe(m)

	m, _ = tr.get("/users/1/posts")
	m.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := PathParams(r)
		eq(t, params.Value("user_id"), "1")
	})
	eq(t, m.route.pattern, "/users/:user_id/posts")
	testServe(m)

	m, _ = tr.get("/users/1/posts/2")
	m.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := PathParams(r)
		eq(t, params.Value("user_id"), "1")
		eq(t, params.Value("id"), "2")
	})
	eq(t, m.route.pattern, "/users/:user_id/posts/:id")
	testServe(m)

	m, _ = tr.get("/users/1/posts/latest")
	m.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := PathParams(r)
		eq(t, params.Value("user_id"), "1")
	})
	eq(t, m.route.pattern, "/users/:user_id/posts/latest")
	testServe(m)

	m, _ = tr.get("/users/1/books")
	m.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := PathParams(r)
		eq(t, params.Value("user_id"), "1")
	})
	eq(t, m.route.pattern, "/users/:user_id/:content_kind")
	testServe(m)

	m, _ = tr.get("/users/1/books/2")
	m.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := PathParams(r)
		eq(t, params.Value("user_id"), "1")
		eq(t, params.Value("content_kind"), "books")
		eq(t, params.Value("id"), "2")
	})
	eq(t, m.route.pattern, "/users/:user_id/:content_kind/:id")
	testServe(m)

	m, _ = tr.get("/users/1/books/latest")
	eq(t, m.route.pattern, "/users/:user_id/:content_kind/latest")

	m, _ = tr.get("/x/10/y/12")
	m.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := PathParams(r)
		eq(t, len(params.Values("splat")), 2)
		eq(t, params.Values("splat")[0], "10")
		eq(t, params.Values("splat")[1], "12")
	})
	eq(t, m.route.pattern, "/x/*/y/*")
	testServe(m)

	// Middle splat is match until /
	m, _ = tr.get("/x/10/users/15")
	m.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := PathParams(r)
		eq(t, params.Value("splat"), "10")
		eq(t, params.Value("id"), "15")
	})
	eq(t, m.route.pattern, "/x/*/users/:id")
	testServe(m)

	// Last splat is match until $
	m, _ = tr.get("/assets/css/app.css")
	m.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := PathParams(r)
		eq(t, params.Value("splat"), "css/app.css")
	})
	eq(t, m.route.pattern, "/assets/*")
	testServe(m)
}

func testServe(m *routeWithParam) {
	r, _ := http.NewRequest("GET", "http://example.com", nil)
	m.serveHTTP(nil, r)
}

func eq(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("expected <%s>, but got <%s>", expected, actual)
	}
}

func neq(t *testing.T, expected, actual interface{}) {
	if expected == actual {
		t.Errorf("not expected <%s>, but got <%s>", expected, actual)
	}
}
