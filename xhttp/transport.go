// Copyright (c) 2020, Geert JM Vanderkelen

package xhttp

import (
	"net/http"
)

type transport struct {
	trp                   http.Transport
	bearer                string
	tlsInsecureSkipVerify bool
	contentType           string
}

// RoundTrip implements a RoundTripper over HTTP.
//
// When t has a authorization bearer, it will set the HTTP Authorization
// header.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.bearer != "" {
		req.Header.Set(HeaderAuthorization, "Bearer "+t.bearer)
	}
	req.Header.Set(HeaderContentType, t.contentType)
	return t.trp.RoundTrip(req)
}
