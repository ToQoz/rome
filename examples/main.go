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
		w.Write([]byte("/posts/:id\n"))
		w.Write([]byte(fmt.Sprintf("FromValue(\"id\"): %s\n", r.FormValue("id"))))
	}))

	router.Get("/x/:id/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/x/:id/*\n"))
		w.Write([]byte(fmt.Sprintf("FromValue(\"id\"): %s\n", r.FormValue("id"))))
		w.Write([]byte(fmt.Sprintf("FromValue(\"splat\"): %s\n", r.FormValue("splat"))))
		w.Write([]byte(fmt.Sprintf("From[\"splat\"]: %s\n", r.Form["splat"])))
	}))

	log.Printf("Listen and Serve, %q", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
