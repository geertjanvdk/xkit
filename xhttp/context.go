// Copyright (c) 2021, Geert JM Vanderkelen

package xhttp

// contextKey is a value for use with context.WithValue.
type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return "<xhttp/context:" + c.name + ">"
}
