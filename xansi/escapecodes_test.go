// Copyright (c) 2021, Geert JM Vanderkelen

package xansi

import (
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestFgColor256(t *testing.T) {
	t.Run("bad rgb values panics", func(t *testing.T) {
		xt.Panics(t, func() {
			FgColor24(-1, 0, 255)
		})

		xt.Panics(t, func() {
			FgColor24(0, 888, 255)
		})

		xt.Panics(t, func() {
			FgColor24(0, 255, 888)
		})
	})
}
