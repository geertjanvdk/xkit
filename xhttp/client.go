// Copyright (c) 2020, Geert JM Vanderkelen

package xhttp

import (
	"crypto/tls"
	"io"
	"net/http"
)

const (
	defaultContentType = "application/json"
)

// Client is a wrapper around Go's http.Client storing the URI and options.
type Client struct {
	*http.Client
	URI     string
	options *clientOptions
}

// NewClient instantiates Client with uri and optional functional
// options.
func NewClient(uri string, options ...ClientOption) *Client {
	opt := newClientOptions()
	for _, o := range options {
		o(opt)
	}

	c := &Client{
		URI:     uri,
		options: opt,
		Client: &http.Client{
			Transport: &transport{
				trp: http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: opt.TLSInsecureSkipVerify,
					},
				},
				bearer:                opt.AuthBearer,
				contentType:           opt.getContentType(),
				tlsInsecureSkipVerify: opt.TLSInsecureSkipVerify,
			},
			Timeout: opt.Timeout,
		},
	}
	return c
}

// Get uses the HTTP get method to make a request using URI of c.
func (c *Client) Get() (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, c.URI, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post uses the HTTP post method to make a request using URI of c sending body as payload.
func (c *Client) Post(body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, c.URI, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// SetAuthzBearer sets the authorization bearer.
func (c *Client) SetAuthzBearer(s string) {
	c.Client.Transport.(*transport).bearer = s
	c.options.AuthBearer = s
}

// AuthzBearer returns the authorization bearer of c.
func (c Client) AuthzBearer() string {
	return c.Client.Transport.(*transport).bearer
}
