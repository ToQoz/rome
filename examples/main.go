package main

import (
	"fmt"
	"github.com/ToQoz/rome"
	"log"
	"net/http"
)

var (
	addr = ":8877"
)

func main() {
	router := rome.NewRouter()

	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	}))

	router.Get("/posts", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/posts"))
	}))

	router.Get("/posts/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := rome.PathParams(r)
		w.Write([]byte("pattern: /posts/:id"))
		w.Write([]byte{'\n'})
		w.Write([]byte(fmt.Sprintf(`Params.Value("id"): %s`, params.Value("id"))))
		w.Write([]byte{'\n'})
	}))

	router.Get("/x/:id/*/*.*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := rome.PathParams(r)
		w.Write([]byte("pattern /x/:id/*"))
		w.Write([]byte{'\n'})
		w.Write([]byte(fmt.Sprintf(`Params.Value("id"): %s`, params.Value("id"))))
		w.Write([]byte{'\n'})
		w.Write([]byte(fmt.Sprintf(`Params.Value("splat"): %s`, params.Value("splat"))))
		w.Write([]byte{'\n'})
		w.Write([]byte(fmt.Sprintf(`Params.Values("splat"): %q`, params.Values("splat"))))
		w.Write([]byte{'\n'})
	}))

	log.Printf("Listen and Serve, %q", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
