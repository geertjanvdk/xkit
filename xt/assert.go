// Copyright (c) 2020, Geert JM Vanderkelen

package xt

import "testing"

// Assert tests if condition is true. When not, messages are displayed.
// This is equivalent as Eq(t, true, condition, messages...).
func Assert(t *testing.T, condition bool, messages ...string) {
	TestHelper(t)

	Eq(t, true, condition, messages...)
}
