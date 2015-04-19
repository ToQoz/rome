package rome

import (
	"net/http"
	"testing"
)

func Benchmark_RouterParam1(b *testing.B) {
	router := NewRouter()
	router.Get("/users/:name", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := PathParams(r)
		params.Value("name")
	}))
	r, _ := http.NewRequest("GET", "/users/toqoz", nil)
	benchRequest(b, router, r)
}

func Benchmark_RouterParam5(b *testing.B) {
	router := NewRouter()
	router.Get("/a/:a/b/:b/c/:c/d/:d/e/:e", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := PathParams(r)
		params.Value("a")
		params.Value("b")
		params.Value("c")
		params.Value("d")
		params.Value("e")
	}))
	r, _ := http.NewRequest("GET", "/a/10/b/hoge/c/8/d/99/e/a", nil)
	benchRequest(b, router, r)
}

func benchRequest(b *testing.B, router http.Handler, r *http.Request) {
	w := new(mockResponseWriter)
	u := r.URL
	rq := u.RawQuery
	r.RequestURI = u.RequestURI()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		u.RawQuery = rq
		router.ServeHTTP(w, r)
	}
}

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}
