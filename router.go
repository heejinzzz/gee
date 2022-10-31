package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	// trie trees' root node. Example: roots["GET"] roots["POST"]
	roots map[string]*node
}

func newRouter() *router {
	return &router{roots: make(map[string]*node)}
}

func splitPattern(pattern string) []string {
	split := strings.Split(pattern, "/")

	var parts []string
	for _, s := range split {
		if s != "" {
			parts = append(parts, s)
			if s[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handlers []HandlerFunc) {
	log.Printf("Register Route: %4s - %s\n", method, pattern)
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	parts := splitPattern(pattern)
	r.roots[method].insert(pattern, parts, 0, handlers)
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	parts := splitPattern(path)
	params := map[string]string{}

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(parts, 0)
	if n == nil {
		return nil, nil
	}

	patternParts := splitPattern(n.pattern)
	for i, part := range patternParts {
		if part[0] == ':' {
			params[part[1:]] = parts[i]
		}
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(parts[i:], "/")
			break
		}
	}
	return n, params
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		c.handlersChain = append(c.handlersChain, n.handlersChain...)
	} else {
		c.handlersChain = append(c.handlersChain, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
