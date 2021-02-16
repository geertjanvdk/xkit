// Copyright (c) 2020, Geert JM Vanderkelen

package xutil

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestRandomBytes(t *testing.T) {
	t.Run("some uniqueness", func(t *testing.T) {
		s := map[string]bool{}
		for i := 0; i < 100000; i++ {
			r, err := RandomBytes(16)
			xt.OK(t, err)
			henc := hex.EncodeToString(r)
			xt.Assert(t, !s[henc], "expected at least some uniqueness")
			s[henc] = true
		}
	})

	for _, n := range []int{16, 8, 33} {
		t.Run(fmt.Sprintf("length %d", n), func(t *testing.T) {
			r, err := RandomBytes(n)
			xt.OK(t, err)
			xt.Eq(t, n, len(r))
		})
	}

	t.Run("panics if n < 1", func(t *testing.T) {
		xt.Panics(t, func() {
			_, _ = RandomBytes(0)
		})

		xt.Panics(t, func() {
			_, _ = RandomBytes(-20)
		})
	})
}

func TestRandomAlphaNumeric(t *testing.T) {
	t.Run("some uniqueness", func(t *testing.T) {
		s := map[string]bool{}
		for i := 0; i < 100000; i++ {
			r := RandomAlphaNumeric(16)
			xt.Assert(t, !s[r], "expected at least some uniqueness")
			s[r] = true
		}
	})

	for _, n := range []int{16, 8, 33} {
		t.Run(fmt.Sprintf("length %d", n), func(t *testing.T) {
			r := RandomAlphaNumeric(n)
			xt.Eq(t, n, len(r))
		})
	}

	t.Run("panics if n < 1", func(t *testing.T) {
		xt.Panics(t, func() {
			_ = RandomAlphaNumeric(0)
		})

		xt.Panics(t, func() {
			_ = RandomAlphaNumeric(-20)
		})
	})
}
