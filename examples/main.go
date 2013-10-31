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
		w.Write([]byte("posts"))
	}))

	router.Get("/posts/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("post: %s", r.URL.Query().Get("id"))))
	}))

	log.Printf("Listen and Serve, %q", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
