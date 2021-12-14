// Copyright (c) 2020, Geert JM Vanderkelen

/*
Package xhttp contains tools and helper functions around the Go http package.

### HTTP Client

Creating a HTTP client is done using the `NewClient` function. It takes functional
options so that it is easier to configure.

The xhttp.Client struct wraps around Go's http.Client. It is instantiated through
the constructor factory function NewClient:

	c := xhttp.NewClient("http://example.com")

Using functional options, we can also configure the client while creating it:

	c := xhttp.NewClient("http://example.com", xhttp.WithBearer("your_token"))

The following functional options are available:

- WithTLSInsecure()
- WithBearer(string)

*/
package xhttp
