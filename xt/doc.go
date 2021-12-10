// Copyright (c) 2021, Geert JM Vanderkelen

/*
Package xt offers a number of functions to make testing more comfortable.

The following functions are available:

* Eq(t *testing.T, expect, have interface{}, messages ...string)
* Assert(t *testing.T, condition bool, messages ...string)
* Panics(t *testing.T, f func())
* Match(t *testing.T, pattern, s string, messages ...string)

The above is already plenty as Eq and Panics would suffice.

*/
package xt
