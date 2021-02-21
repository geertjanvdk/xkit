// Copyright (c) 2020, 2021, Geert JM Vanderkelen

package xid

import (
	"fmt"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestNewV4(t *testing.T) {
	u := UUIDv4()
	xt.Assert(t, !u.IsNil())
	xt.Eq(t, 4, u.Version())
	xt.Assert(t, IsUUIDV4(u.String()))
}

func TestNewV5(t *testing.T) {
	u := UUIDv5(UUIDURLNamespace, "http://golang.org")
	xt.Assert(t, !u.IsNil())
	xt.Eq(t, 5, u.Version())
	xt.Assert(t, IsUUIDV5(u.String()))

	j, err := u.MarshalJSON()
	xt.OK(t, err)
	xt.Eq(t, "\"627f0098-88da-5243-af8b-0f3ac4fcd527\"", string(j))

	nu := UUID{}
	err = nu.UnmarshalJSON(j)
	xt.OK(t, err)
	xt.Eq(t, "627f0098-88da-5243-af8b-0f3ac4fcd527", nu.String())
}

func TestIsValid(t *testing.T) {
	var casesV1 = []string{
		"f9f68d4c-0986-11ea-8d71-362b9e155667",
		"02ca5e30-0987-11ea-8d71-362b9e155667",
	}
	var casesV4 = []string{
		"95fc86ec-dc52-43fc-b04c-6a61666fdd34",
		"6ebbac72-6dd1-45dd-8f4a-6814bfe03cda",
	}

	for _, u := range casesV1 {
		xt.Assert(t, UUIDIsValid(u), fmt.Sprintf("expected %s to be valid UUID", u))
		xt.Assert(t, IsUUIDV1(u), fmt.Sprintf("expected %s to be v1", u))
		xt.Assert(t, !IsUUIDV4(u), fmt.Sprintf("expected %s not to be v4", u))
		xt.Eq(t, 1, UUIDFromString(u).Version())
	}

	for _, u := range casesV4 {
		xt.Assert(t, UUIDIsValid(u), fmt.Sprintf("expected %s to be valid UUID", u))
		xt.Assert(t, !IsUUIDV1(u), fmt.Sprintf("expected %s not to be v1", u))
		xt.Assert(t, IsUUIDV4(u), fmt.Sprintf("expected %s to be v4", u))
		xt.Eq(t, 4, UUIDFromString(u).Version())
	}

	for _, u := range append(casesV1, casesV4...) {
		xt.Assert(t, !IsUUIDV5(u), fmt.Sprintf("expected %s not to be v5", u))
	}
}

func TestNewFromString(t *testing.T) {
	exp := UUID{0xf9, 0xf6, 0x8d, 0x4c, 0x9, 0x86, 0x11, 0xea, 0x8d, 0x71, 0x36, 0x2b, 0x9e, 0x15, 0x56, 0x67}
	xt.Eq(t, exp, UUIDFromString("f9f68d4c-0986-11ea-8d71-362b9e155667"))

	xt.Assert(t, UUIDFromString("something not UUIDish").IsNil())
}
