// Copyright (c) 2020, Geert JM Vanderkelen

package xhttp

import "time"

// ClientOption is the functional option type used with
// xhttp.Client.
type ClientOption func(*clientOptions)

type clientOptions struct {
	AuthBearer            string
	TLSInsecureSkipVerify bool
	Timeout               time.Duration
	ContentType           string
}

func (co clientOptions) getContentType() string {
	return co.ContentType
}

func newClientOptions() *clientOptions {
	return &clientOptions{
		AuthBearer:            "",
		TLSInsecureSkipVerify: false,
		Timeout:               0,
		ContentType:           defaultContentType,
	}
}

// WithBearer is a functional option to use with the factory function
// NewClient. It sets the authorization bearer b which will be send
// as HTTP header with each request.
func WithBearer(b string) ClientOption {
	return func(options *clientOptions) {
		options.AuthBearer = b
	}
}

// WithTLSInsecure is a function option for xhttp.NewClient. It sets
// whether the TLS certificates of the connecting server needs to be
// verified.
// This should be used with great care, and not in production
// environments.
func WithTLSInsecure() ClientOption {
	return func(options *clientOptions) {
		options.TLSInsecureSkipVerify = true
	}
}
