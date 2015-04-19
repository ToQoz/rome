package rome

import (
	"net/http"
)

var ParamsMap = map[*http.Request]Params{}

// WARNING:
//   This is NOT goroutine safe.
//   Now, I'm going to make this goroutine safe.

func setParams(r *http.Request, ps Params) {
	ParamsMap[r] = ps

}
func clearParams(r *http.Request) {
	delete(ParamsMap, r)
}

func PathParams(r *http.Request) Params {
	return ParamsMap[r]
}
