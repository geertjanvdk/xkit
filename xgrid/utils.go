// Copyright (c) 2021, Geert JM Vanderkelen

package xgrid

import (
	"fmt"
	"io"

	"github.com/geertjanvdk/xkit/xansi"
)

// boldSprintf formats according to a format specifier and returns
// the resulting string bold.
func boldSprintf(format string, a ...interface{}) string {
	return xansi.Render{xansi.Bold}.Sprintf(format, a...)
}

// MustFprint wrappers around fmt.Fprint ignoring any error.
func MustFprint(w io.Writer, a ...interface{}) {
	_, _ = fmt.Fprint(w, a...)
}
