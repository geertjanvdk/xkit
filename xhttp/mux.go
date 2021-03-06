// Copyright (c) 2021, Geert JM Vanderkelen

package xhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strconv"
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

const defaultCaptureConverter = "str"

var ErrIncorrectCaptureConverter = errors.New("incorrect capture converter")

type Capture struct {
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	Converter string      `json:"converter"`
}

// AsStr returns the captured value as string. This will always return
// as values are retrieved as string.
func (c Capture) AsStr() string {
	return c.Value.(string)
}

// AsInt64 tries to convert the captured value and return it as int64.
// Returns error ErrIncorrectCaptureConverter when converter is not
// `int`.
func (c Capture) AsInt64() (int64, error) {
	if c.Converter != "int" {
		return 0, ErrIncorrectCaptureConverter
	}
	return strconv.ParseInt(c.Value.(string), 10, 64)
}

// AsInt tries to convert the captured value and return it as int.
// Returns error ErrIncorrectCaptureConverter when converter is not
// `int`.
func (c Capture) AsInt() (int, error) {
	if c.Converter != "int" {
		return 0, ErrIncorrectCaptureConverter
	}
	n, err := strconv.ParseInt(c.Value.(string), 10, 32)
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

type Captures map[string]Capture

type reHandler struct {
	pattern  string
	regex    string
	compiled *regexp.Regexp
	handler  http.Handler
	methods  []Method
	captures Captures
}

// setPattern will store pattern into r after validating it and setting up
// optional captures.
// Panics when pattern could not be compiled and is invalid, when an
// unsupported converted type is used for capturing values, and when anything
// is wrong when parsing captures.
func (r *reHandler) setPattern(pattern string) {
	captureRegs := map[string]string{
		"int": `(?P<%s>\d{1,19})`, // digits of 2^64-1
		"str": `(?P<%s>[\w-_]+)`,
	}

	reCaptures := regexp.MustCompile(`<((?:(int|str):)?([0-9A-Za-z_]+))>`)

	matches := reCaptures.FindAllStringSubmatch(pattern, -1)
	newPattern := pattern

	if matches != nil {
		r.captures = Captures{}
		pos := 0

		for _, m := range matches {
			name := m[3]
			if _, have := r.captures[name]; have {
				panic("(xhttp) pattern capture value `" + name + "` specified twice")
			}

			_ = captureRegs
			capConverter := m[2]
			if capConverter == "" {
				capConverter = defaultCaptureConverter
			}
			capReg, have := captureRegs[capConverter]
			if !have {
				panic("(xhttp) invalid pattern capture Converter type; was " + capConverter)
			}
			aOpen := strings.IndexByte(newPattern[pos:], '<') + pos
			aClose := strings.IndexByte(newPattern[pos:], '>') + pos

			subReg := fmt.Sprintf(capReg, name)
			newPattern = newPattern[0:aOpen] + subReg + newPattern[aClose+1:]

			pos = aOpen + len(subReg) + 1

			r.captures[name] = Capture{
				Name:      name,
				Converter: capConverter,
			}
		}
	}

	if strings.ContainsAny(newPattern, "<>") {
		panic("(xhttp) invalid pattern capture with mismatch angle brackets; was " + pattern)
	}

	reg, err := regexp.Compile(newPattern)
	if err != nil {
		panic("(xhttp) invalid pattern; was " + pattern)
	}

	r.pattern = pattern
	r.regex = newPattern
	r.compiled = reg
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
// registered patterns provide as regular expressions.
//
// Every registered pattern can also optionally associated with allowed
// HTTP methods. When no allowed method is provided, it will only allow GET.
//
// A pattern can also capture values using `<>` angle brackets.
// For example, `/blog/<blogUID>` will capture the value where `<blogUID>`
// and it will be made available in the request's context under the key
// `xhttp.CapturesContextKey` as an instance of `xhttp.Capture`.
// It is also possible to add a type, for example `<int:blogID>`, so
// that the methods `AsInt` and `AsInt64` of `xhttp.Capture` can be used
// to get the converted value.
// It is not possible to use the same named capture twice.
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

// Handle registers the handler for the given pattern, which is a regular
// expression with optional captures in angle brackets.
// Panics when handler already exists for pattern, or if pattern could not
// compile the expression.
func (s *ServeReMux) Handle(pattern string, handler http.Handler, methods ...Method) {
	if s.handlers.Has(pattern) {
		panic("(xhttp) pattern `" + pattern + "` already registered")
	}

	h := &reHandler{}
	h.setPattern(pattern)
	h.methods = methods
	h.handler = handler

	s.handlers.Set(pattern, h)
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
	h, pattern, captures := s.Handler(r)
	ctx := context.WithValue(r.Context(), CapturesContextKey, captures)
	ctx = context.WithValue(ctx, RegexpMatchContextKey, pattern)
	r = r.Clone(ctx)
	h.ServeHTTP(w, r)
}

// Handler returns the handler to use for the given request.
// Panics when registered handler is not supported.
func (s *ServeReMux) Handler(r *http.Request) (http.Handler, string, *Captures) {
	p := requestPathCleanUp(r.URL.Path)
	var foundMatchButNotAllowed bool
	for _, v := range s.handlers.Values() {
		h, ok := v.(*reHandler)
		if !ok {
			panic(fmt.Sprintf("(xhttp) ServeReMux has unsupported handler registered; was %v", p))
		}

		if h.compiled.MatchString(p) {
			if h.allowedMethod(Method(r.Method)) {
				var captures *Captures
				if h.captures != nil {
					caps := Captures{}
					matches := h.compiled.FindStringSubmatch(p)
					for i, name := range h.compiled.SubexpNames() {
						if i != 0 {
							caps[name] = Capture{
								Name:      name,
								Value:     matches[i],
								Converter: h.captures[name].Converter,
							}
						}
					}
					captures = &caps
				}
				return h.handler, h.pattern, captures
			}
			foundMatchButNotAllowed = true
		}
	}

	if foundMatchButNotAllowed {
		return MethodNotAllowedHandler(), "", nil
	}

	return NotFoundHandler(), "", nil
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
