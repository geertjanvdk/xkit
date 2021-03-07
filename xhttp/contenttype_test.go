// Copyright (c) 2020, Geert JM Vanderkelen

package xhttp

import (
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestContentTypes(t *testing.T) {
	xt.Eq(t, "text/plain; charset=utf-8", ContentTypePlain)
	xt.Eq(t, "text/html; charset=utf-8", ContentTypeHTML)
	xt.Eq(t, "application/json; charset=utf-8", ContentTypeJSON)
	xt.Eq(t, "application/octet-stream", ContentTypeBinary)
}
