// Copyright (c) 2021, Geert JM Vanderkelen

package xhttp

// contextKey is a value for use with context.WithValue.
type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return "<xhttp/context:" + c.name + ">"
}

var (
	// RegexpMatchContextKey is a context key which is used to register
	// within the request which pattern was matched by ServeReMux. It's
	// associated type is string.
	RegexpMatchContextKey = &contextKey{name: "xhttp.ServeReMux.RegexpMatch"}

	// CapturesContextKey is a context key which is used to register
	// within the request the captured values in the path when matching
	// URLs using ServeReMux.
	CapturesContextKey = &contextKey{name: "xhttp.ServeReMux.Captures"}
)
