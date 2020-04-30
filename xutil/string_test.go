// Copyright (c) 2020, Geert JM Vanderkelen

package xutil

import (
	"reflect"
	"testing"

	"lab.scrum.pub/go/ts"
)

func TestIsZeroString(t *testing.T) {
	var nullString *string
	cases := []struct {
		have interface{}
		exp  bool
	}{
		{"", true},
		{nullString, true},
		{"false", false},
		{"*false", false},
	}

	for _, c := range cases {
		t.Run("no panic", func(t *testing.T) {
			ts.Eq(t, c.exp, IsZeroString(c.have))
		})
	}

	panicCases := []interface{}{nil, 123, []byte("not bytes")}
	for _, c := range panicCases {
		t.Run("panic", func(t *testing.T) {
			ts.Panics(t, func() {
				IsZeroString(c)
			})
		})
	}
}

func TestStringPtr(t *testing.T) {
	rv := reflect.ValueOf(StringPtr("I should be pointer"))
	ts.Assert(t, rv.Kind() == reflect.Ptr)
}
