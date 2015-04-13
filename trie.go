package rome

import (
	"errors"
	"net/http"
	"strings"
)

var (
	segmentSeparater = "/."
)

// ----------------------------------------------------------------------------
// Trie-tree
// ----------------------------------------------------------------------------

type trie struct {
	root *node
}

func newTrie() *trie {
	return &trie{newNode()}
}

func (t *trie) add(r *route) error {
	return t.root.add(r.pattern, r, []string{})
}

func (t *trie) get(path string) (*routeWithParam, error) {
	var err error

	// Treat "/some//path" as "/some/path" in dispatching route
	path = strings.Replace(path, "//", "/", -1)

	// Treat "/some-path/" as "/some-path" in dispatching route.
	// But treat "/" as "/".
	if path != "/" {
		paths := strings.Split(path, "")
		lastIndex := len(paths) - 1

		if paths[lastIndex] == "/" {
			path = strings.Join(paths[0:lastIndex], "")
		}
	}

	matched := t.root.get(path, []string{}, []string{})

	// Route not found
	if matched == nil {
		err = errors.New("Not found")
	} else {
		matched.params = map[string]string{}

		for i, pkey := range matched.route.paramKeys {
			matched.params[pkey] = matched.paramValues[i]
		}
	}

	return matched, err
}

// ----------------------------------------------------------------------------
// Trie-Node
// ----------------------------------------------------------------------------

type node struct {
	route      *route
	children   map[string]*node
	paramChild *node
	splatChild *node
}

func newNode() *node {
	return &node{}
}

func (n *node) add(path string, r *route, paramKeys []string) error {
	// Finish to consume path
	if path == "" {
		if n.route != nil {
			return errors.New("Deplicated routing")
		}

		r.paramKeys = paramKeys
		n.route = r
		return nil
	}

	token, tail := firstToken(path)

	// :param branch
	switch token {
	case ":":
		var pkey string

		if n.paramChild == nil {
			n.paramChild = newNode()
		}

		pkey, tail = firstSegments(tail)

		paramKeys = append(paramKeys, pkey)

		return n.paramChild.add(tail, r, paramKeys)
	case "*":
		if n.splatChild == nil {
			n.splatChild = newNode()
		}

		_, tail = firstSegments(tail)
		return n.splatChild.add(tail, r, paramKeys)
	}

	// Main branches
	if n.children == nil {
		n.children = map[string]*node{}
	}

	if n.children[token] == nil {
		n.children[token] = newNode()
	}

	return n.children[token].add(tail, r, paramKeys)
}

func (n *node) get(path string, paramValues []string, splat []string) *routeWithParam {
	// Finish to consume path
	if path == "" {
		if n.route != nil {
			// Found route!
			return newRouteWithParam(n.route, paramValues, splat)
		} else {
			// Not found!
			return nil
		}
	}

	// Search main branch
	token, tail := firstToken(path)

	if n.children[token] != nil {
		match := n.children[token].get(tail, paramValues, splat)

		if match != nil {
			// Found route in main branches!
			return match
		}
	}

	// Search splat branch
	if n.splatChild != nil {
		var match *routeWithParam

		if n.splatChild.children == nil &&
			n.splatChild.paramChild == nil &&
			n.splatChild.splatChild == nil {
			match = n.splatChild.get("", paramValues, append(splat, path))
		} else {
			var s string
			s, tail = firstSegments(path)
			match = n.splatChild.get(tail, paramValues, append(splat, s))
		}

		if match != nil {
			// Found route in splat branch!
			return match
		}
	}

	// Search param branch
	if n.paramChild != nil {
		var pvalue string

		pvalue, tail = firstSegments(path)

		match := n.paramChild.get(tail, append(paramValues, pvalue), splat)

		if match != nil {
			// Found route in param branch!
			return match
		}
	}

	// Not found!
	return nil
}

// ----------------------------------------------------------------------------
// Util method
// ----------------------------------------------------------------------------

func firstToken(str string) (token string, tail string) {
	tokens := strings.Split(str, "")
	token = tokens[0]
	tail = strings.Join(tokens[1:], "")
	return
}

func firstSegments(path string) (firstSegments string, tail string) {
	i := strings.IndexAny(path, segmentSeparater)

	if i != -1 {
		// e.g. when tail: "id/content" => firstSegments:  "id"
		//                              => tail:           "/content"
		//
		// e.g. when tail: "id.json"    => firstSegments:  "id"
		//                              => tail:           ".json"
		firstSegments = path[:i]
		tail = path[i:]
	} else {
		// e.g. when tail: "id"         => firstSegments:  "id"
		//                              => tail:           ""
		firstSegments = path
		tail = ""
	}

	return
}

// ----------------------------------------------------------------------------
// Matched Route(Trie's leaf Value)
// ----------------------------------------------------------------------------

type routeWithParam struct {
	route       *route
	paramValues []string
	params      map[string]string
	splat       []string
}

func newRouteWithParam(r *route, paramValues []string, splat []string) *routeWithParam {
	return &routeWithParam{route: r, paramValues: paramValues, splat: splat}
}

func (rp *routeWithParam) serveHTTP(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	for pkey, pvalue := range rp.params {
		req.Form.Add(pkey, pvalue)
	}
	for _, s := range rp.splat {
		req.Form.Add("splat", s)
	}
	rp.route.handler.ServeHTTP(w, req)
}
