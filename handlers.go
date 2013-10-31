package rome

import (
	"net/http"
)

type notFoundHandler struct {
	http.Handler
}

func (handler *notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
