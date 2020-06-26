// Copyright (c) 2020, Geert JM Vanderkelen

package xt

import "testing"

func Assert(t *testing.T, condition bool, messages ...string) {
	TestHelper(t)

	Eq(t, true, condition, messages...)
}
