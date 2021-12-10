// Copyright (c) 2021, Geert JM Vanderkelen

package main

import (
	"fmt"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestGitHubHeaderID(t *testing.T) {
	cases := []struct {
		have string
		exp  string
	}{
		{have: "café", exp: "caf%C3%A9"},
		{have: "Löffel", exp: "l%C3%B6ffel"},
		{have: "func (c *Client) SetAuthzBearer(s string)",
			exp: "func-c-client-setauthzbearers-string"},
		{have: "func (c *Client) SetAuthzBearer(s string)", // 2nd pass adds -1 suffix
			exp: "func-c-client-setauthzbearers-string-1"},
		{have: "func (c *Client) SetAuthzBearer(s string)", // 3nd pass adds -2 suffix
			exp: "func-c-client-setauthzbearers-string-2"},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%02d#%s", i, c.have), func(t *testing.T) {
			xt.Eq(t, c.exp, gitHubHeaderID(c.have))
		})
	}

}
