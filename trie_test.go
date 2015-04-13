package rome

import (
	"github.com/ToQoz/rome/test_helpers"
	"net/http"
	"testing"
)

func _handler(w http.ResponseWriter, r *http.Request) {}

var (
	handler = http.HandlerFunc(_handler)
)

func TestTrie(t *testing.T) {
	var err error
	var matched *routeWithParam

	tr := newTrie()
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

	_, err = tr.get("/NOT_FOUND")
	test_helpers.AssertNotEqual(t, nil, err)

	matched, _ = tr.get("/users")
	test_helpers.AssertEqual(t, matched.route.pattern, "/users")

	matched, _ = tr.get("/users/")
	test_helpers.AssertEqual(t, matched.route.pattern, "/users")

	matched, _ = tr.get("/users/1")
	test_helpers.AssertEqual(t, matched.route.pattern, "/users/:id")
	test_helpers.AssertEqual(t, "1", matched.params["id"])

	matched, _ = tr.get("/users//2")
	test_helpers.AssertEqual(t, matched.route.pattern, "/users/:id")
	test_helpers.AssertEqual(t, "2", matched.params["id"])

	matched, _ = tr.get("/users/1/posts")
	test_helpers.AssertEqual(t, matched.route.pattern, "/users/:user_id/posts")
	test_helpers.AssertEqual(t, "1", matched.params["user_id"])

	matched, _ = tr.get("/users/1/posts/2")
	test_helpers.AssertEqual(t, matched.route.pattern, "/users/:user_id/posts/:id")
	test_helpers.AssertEqual(t, "1", matched.params["user_id"])
	test_helpers.AssertEqual(t, "2", matched.params["id"])

	matched, _ = tr.get("/users/1/posts/latest")
	test_helpers.AssertEqual(t, matched.route.pattern, "/users/:user_id/posts/latest")
	test_helpers.AssertEqual(t, "1", matched.params["user_id"])

	matched, _ = tr.get("/users/1/books")
	test_helpers.AssertEqual(t, matched.route.pattern, "/users/:user_id/:content_kind")
	test_helpers.AssertEqual(t, "1", matched.params["user_id"])

	matched, _ = tr.get("/users/1/books/2")
	test_helpers.AssertEqual(t, matched.route.pattern, "/users/:user_id/:content_kind/:id")
	test_helpers.AssertEqual(t, "1", matched.params["user_id"])
	test_helpers.AssertEqual(t, "books", matched.params["content_kind"])
	test_helpers.AssertEqual(t, "2", matched.params["id"])

	matched, _ = tr.get("/users/1/books/latest")
	test_helpers.AssertEqual(t, matched.route.pattern, "/users/:user_id/:content_kind/latest")
	test_helpers.AssertEqual(t, "1", matched.params["user_id"])
	test_helpers.AssertEqual(t, "books", matched.params["content_kind"])

	matched, _ = tr.get("/users/1/parts")
	test_helpers.AssertEqual(t, matched.route.pattern, "/users/:user_id/parts")
	test_helpers.AssertEqual(t, "1", matched.params["user_id"])

	matched, _ = tr.get("/users/1/parts/2")
	test_helpers.AssertEqual(t, matched.route.pattern, "/users/:user_id/parts/:id")
	test_helpers.AssertEqual(t, "1", matched.params["user_id"])
	test_helpers.AssertEqual(t, "2", matched.params["id"])

	matched, _ = tr.get("/users/1/parts/latest")
	test_helpers.AssertEqual(t, matched.route.pattern, "/users/:user_id/parts/latest")
	test_helpers.AssertEqual(t, "1", matched.params["user_id"])

	matched, _ = tr.get("/users/login/posts")
	test_helpers.AssertEqual(t, matched.route.pattern, "/users/login/posts")

	matched, _ = tr.get("/users/login/posts/1")
	test_helpers.AssertEqual(t, matched.route.pattern, "/users/login/posts/:id")
	test_helpers.AssertEqual(t, "1", matched.params["id"])

	matched, _ = tr.get("/x/10/y/12")
	test_helpers.AssertEqual(t, matched.route.pattern, "/x/*/y/*")
	test_helpers.AssertEqual(t, "10", matched.splat[0])
	test_helpers.AssertEqual(t, "12", matched.splat[1])
	matched.route.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		test_helpers.AssertEqual(t, len(r.Form["splat"]), 2)
		test_helpers.AssertEqual(t, r.Form["splat"][0], "10")
		test_helpers.AssertEqual(t, r.Form["splat"][1], "12")
	})
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Error(err)
	}
	matched.serveHTTP(nil, req)

	// Middle splat is match until /
	matched, _ = tr.get("/x/10/users/15")
	test_helpers.AssertEqual(t, matched.route.pattern, "/x/*/users/:id")
	test_helpers.AssertEqual(t, "10", matched.splat[0])
	test_helpers.AssertEqual(t, "15", matched.params["id"])

	// Last splat is match until $
	matched, _ = tr.get("/assets/css/app.css")
	test_helpers.AssertEqual(t, matched.route.pattern, "/assets/*")
	test_helpers.AssertEqual(t, "css/app.css", matched.splat[0])

	matched, _ = tr.get("/assets/css/app.css")
	matched.route.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		test_helpers.AssertEqual(t, r.FormValue("splat"), "css/app.css")
		test_helpers.AssertEqual(t, r.FormValue("x"), "y")
	})
	req, err = http.NewRequest("GET", "http://example.com?x=y", nil)
	if err != nil {
		t.Error(err)
	}
	matched.serveHTTP(nil, req)
}

func TestRouteWithParam(t *testing.T) {
}

var (
	tr = newTrie()
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
}

func BenchmarkTrie(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tr.get("/users/login/posts/1")
		tr.get("/users")
		tr.get("/NOT_FOUND")
	}
}
