// Copyright (c) 2021, Geert JM Vanderkelen

package xid

import (
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestNanoID(t *testing.T) {
	t.Run("default sized and with default alphabet", func(t *testing.T) {
		res := map[string]bool{}
		for i := 0; i < 10000; i++ {
			id := NanoID().String()
			_, have := res[id]
			xt.Assert(t, !have, "expected no duplicates")
			xt.Eq(t, nanoDefaultSize, len(id))
			res[id] = true
		}
	})
}

func TestNanoID_SetSize(t *testing.T) {
	t.Run("size 12 nanoid", func(t *testing.T) {
		res := map[string]bool{}
		for i := 0; i < 10000; i++ {
			id := NanoID().SetSize(12).String()
			_, have := res[id]
			xt.Assert(t, !have, "expected no duplicates")
			xt.Eq(t, 12, len(id))
			res[id] = true
		}
	})

	t.Run("panics on size smaller than minimum", func(t *testing.T) {
		xt.Panics(t, func() {
			NanoID().SetSize(nanoMinSize - 1)
		})
	})

	t.Run("panics on size bigger than maximum", func(t *testing.T) {
		xt.Panics(t, func() {
			NanoID().SetSize(nanoMaxSize + 1)
		})
	})
}

func TestNanoID_SetAlphabet(t *testing.T) {
	t.Run("default alphabet", func(t *testing.T) {
		xt.Eq(t, urlSafeAlphabet, NanoID().alphabet)
		xt.Eq(t, urlSafeAlphabet, NanoID().SetAlphabet("").alphabet)
	})

	t.Run("support base62", func(t *testing.T) {
		xt.Eq(t, base62Alphabet, string(NanoID().SetAlphabet("base62").alphabet))
	})

	t.Run("support base58", func(t *testing.T) {
		xt.Eq(t, base58Alphabet, string(NanoID().SetAlphabet("base58").alphabet))
	})

	t.Run("support base56", func(t *testing.T) {
		xt.Eq(t, base56Alphabet, string(NanoID().SetAlphabet("base56").alphabet))
	})
}
