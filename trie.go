package rome

import (
	"errors"
	"net/http"
	"strings"
)

var (
	MaxParam = 5
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
		path = strings.TrimSuffix(path, "/")
	}

	matched := t.root.get(path, make([]string, MaxParam), 0)

	// Route not found
	if matched == nil {
		err = errors.New("Not found")
	}

	return matched, err
}

// ----------------------------------------------------------------------------
// Trie-Node
// ----------------------------------------------------------------------------

type node struct {
	route      *route
	children   map[rune]*node
	paramChild *node
	splatChild *node
}

func newNode() *node {
	return &node{}
}

func (n *node) add(path string, r *route, paramKeys []string) error {
	// Finish to consume path
	if len(path) == 0 {
		if n.route != nil {
			return errors.New("Deplicated routing")
		}

		r.paramKeys = paramKeys
		n.route = r
		return nil
	}

	token := rune(path[0])
	tail := path[1:]

	// :param branch
	switch token {
	case ':':
		var pkey string

		if n.paramChild == nil {
			n.paramChild = newNode()
		}

		pkey, tail = firstSegments(tail)
		paramKeys = append(paramKeys, pkey)
		return n.paramChild.add(tail, r, paramKeys)
	case '*':
		if n.splatChild == nil {
			n.splatChild = newNode()
		}

		_, tail = firstSegments(tail)
		paramKeys = append(paramKeys, "splat")
		return n.splatChild.add(tail, r, paramKeys)
	}

	// Main branches
	if n.children == nil {
		n.children = map[rune]*node{}
	}

	if n.children[token] == nil {
		n.children[token] = newNode()
	}

	return n.children[token].add(tail, r, paramKeys)
}

func (n *node) get(path string, paramValues []string, paramValueN int) *routeWithParam {
	// Finish to consume path
	if len(path) == 0 {
		if n.route != nil {
			// Found route!
			return newRouteWithParam(n.route, paramValues)
		} else {
			// Not found!
			return nil
		}
	}

	// Search main branch
	token := rune(path[0])
	tail := path[1:]

	if n.children[token] != nil {
		match := n.children[token].get(tail, paramValues, paramValueN)
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
			paramValues[paramValueN] = path
			paramValueN++
			match = n.splatChild.get("", paramValues, paramValueN)
		} else {
			paramValues[paramValueN], tail = firstSegments(path)
			paramValueN++
			match = n.splatChild.get(tail, paramValues, paramValueN)
		}

		if match != nil {
			// Found route in splat branch!
			return match
		}
	}

	// Search param branch
	if n.paramChild != nil {
		paramValues[paramValueN], tail = firstSegments(path)
		paramValueN++

		match := n.paramChild.get(tail, paramValues, paramValueN)
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

func firstSegments(path string) (head string, tail string) {
	for i := 0; i < len(path); i++ {
		b := path[i]
		if b == '/' || b == '.' {
			head = path[:i]
			tail = path[i:]
			return
		}
	}
	head = path
	return
}

// P is params
type Params []param

func (p Params) Value(key string) string {
	for _, param := range p {
		if param.name == key {
			return param.value
		}
	}
	return ""
}

func (p Params) Values(key string) []string {
	i := 0
	for _, param := range p {
		if param.name == key {
			i++
		}
	}
	ret := make([]string, i)
	i = 0
	for _, param := range p {
		if param.name == key {
			ret[i] = param.value
			i++
		}
	}
	return ret
}

// ----------------------------------------------------------------------------
// Matched Route(Trie's leaf Value)
// ----------------------------------------------------------------------------

type routeWithParam struct {
	*route
	paramValues []string
}

func newRouteWithParam(r *route, paramValues []string) *routeWithParam {
	match := &routeWithParam{route: r, paramValues: paramValues}
	return match
}

func (rp *routeWithParam) serveHTTP(w http.ResponseWriter, req *http.Request) {
	params := make(Params, len(rp.paramKeys))
	for i, key := range rp.paramKeys {
		params[i] = param{name: key, value: rp.paramValues[i]}
	}

	setParams(req, params)
	defer clearParams(req)

	rp.handler.ServeHTTP(w, req)
}
