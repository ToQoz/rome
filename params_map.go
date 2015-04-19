package rome

import (
	"net/http"
)

var paramsMap = map[*http.Request]Params{}

// WARNING:
//   This is NOT goroutine safe.
//   Now, I'm going to make this goroutine safe.

func setParams(r *http.Request, ps Params) {
	paramsMap[r] = ps

}
func clearParams(r *http.Request) {
	delete(paramsMap, r)
}

func PathParams(r *http.Request) Params {
	return paramsMap[r]
}
