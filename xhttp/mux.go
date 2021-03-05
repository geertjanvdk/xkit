// Copyright (c) 2021, Geert JM Vanderkelen

package xhttp

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/geertjanvdk/xkit/xutil"
)

// Method is a type for defining the HTTP methods.
type Method string

// HTTP method constants which are exactly the same as Go's http.Method*
// but typed with our own.
const (
	MethodGet  Method = "GET"
	MethodPost Method = "POST"
)

var (
	// RegexpMatchContextKey is a context key which is used to register
	// within the request which pattern was matched by ServeReMux. It's
	// associated type is string.
	RegexpMatchContextKey = &contextKey{name: "x-ServeReMux-match"}
)

type reHandler struct {
	pattern  string
	compiled *regexp.Regexp
	handler  http.Handler
	methods  []Method
}

func (r reHandler) allowedMethod(method Method) bool {
	if len(r.methods) == 0 && method == MethodGet {
		return true
	}

	for _, m := range r.methods {
		if method == m {
			return true
		}
	}

	return false
}

// ServeMux is an HTTP request multiplexer which matches the path
// of the URL of each incoming request against a list of
// registered regular expressions.
//
// Every registered can also operationally associated with allowed
// HTTP methods. When no allowed method is provided, it will
// only allow GET.
//
// When no regular expression matched, 404 is returned. If a pattern
// matched, but it turns out the method was not allowed, the HTTP status
// 405 (method not allowed) is returned.
type ServeReMux struct {
	handlers xutil.OrderedMap
}

// NewServeReMux allocates and returns a new ServeReMux.
func NewServeReMux() *ServeReMux {
	return &ServeReMux{}
}

// Handle registers the handler for the given regular expression.
// Panics when handler already exists for pattern, or if pattern could not
// compile the expression.
func (s *ServeReMux) Handle(regex string, handler http.Handler, methods ...Method) {
	reg, err := regexp.Compile(regex)
	if err != nil {
		panic("(xhttp) invalid regex pattern registering route")
	}

	if s.handlers.Has(regex) {
		panic("(xhttp) pattern `" + regex + "` already registered")
	}

	s.handlers.Set(regex, &reHandler{
		pattern:  regex,
		compiled: reg,
		handler:  handler,
		methods:  methods,
	})
}

// ServeHTTP dispatches the request to the handler whose
// regular expression matches the path of the request URL.
func (s *ServeReMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// copy from Go's ServeMux.ServeHTTP
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h, pattern := s.Handler(r)
	ctx := context.WithValue(r.Context(), RegexpMatchContextKey, pattern)
	r = r.Clone(ctx)
	h.ServeHTTP(w, r)
}

// Handler returns the handler to use for the given request.
func (s *ServeReMux) Handler(r *http.Request) (http.Handler, string) {
	p := requestPathCleanUp(r.URL.Path)
	var foundMatchButNotAllowed bool
	for _, v := range s.handlers.Values() {
		h, ok := v.(*reHandler)
		if !ok {
			panic(fmt.Sprintf("(xhttp) ServeReMux has unsupported handler registered; was %v", p))
		}

		if h.compiled.MatchString(p) {
			if h.allowedMethod(Method(r.Method)) {
				return h.handler, h.pattern
			}
			foundMatchButNotAllowed = true
		}
	}

	if foundMatchButNotAllowed {
		return MethodNotAllowedHandler(), ""
	}

	return NotFoundHandler(), ""
}

// requestPathCleanUp uses Go's path.Clean to clean up p.
// When p is empty, root `/` will be returned. The result will also
// always end with a `/`.
func requestPathCleanUp(p string) string {
	p = strings.TrimSpace(p)
	if p == "" || p == "/" {
		return "/"
	}

	if p[0] != '/' {
		p = "/" + p
	}

	endSlash := p[len(p)-1] == '/'
	p = path.Clean(p)
	if endSlash {
		p += "/"
	}

	return p
}
