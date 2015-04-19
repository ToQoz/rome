package rome

import (
	"net/http"
)

// ParamsMap is params in URL path.
// WARNING:
//   This is NOT goroutine safe.
//   Now, I'm going to make this goroutine safe.
var ParamsMap = map[*http.Request]Params{}

func setParams(r *http.Request, ps Params) {
	ParamsMap[r] = ps

}
func clearParams(r *http.Request) {
	delete(ParamsMap, r)
}

func PathParams(r *http.Request) Params {
	return ParamsMap[r]
}
