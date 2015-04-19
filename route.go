package rome

import (
	"net/http"
)

type route struct {
	pattern   string
	handler   http.Handler
	paramKeys []string
}

type param struct {
	name  string
	value string
}

func newRoute(pattern string, handler http.Handler) *route {
	return &route{pattern: pattern, handler: handler}
}
