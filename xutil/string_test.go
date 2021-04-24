// Copyright (c) 2020, Geert JM Vanderkelen

package xutil

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
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
			xt.Eq(t, c.exp, IsZeroString(c.have))
		})
	}

	panicCases := []interface{}{nil, 123, []byte("not bytes")}
	for _, c := range panicCases {
		t.Run("panic", func(t *testing.T) {
			xt.Panics(t, func() {
				IsZeroString(c)
			})
		})
	}
}

func TestStringPtr(t *testing.T) {
	rv := reflect.ValueOf(StringPtr("I should be pointer"))
	xt.Assert(t, rv.Kind() == reflect.Ptr)
}

func TestHasString(t *testing.T) {
	// x is the string we are looking for
	x := "foo"

	cases := []struct {
		have []string
		exp  bool
	}{
		{[]string{}, false},
		{[]string{x}, true},
		{[]string{"bar"}, false},
		{[]string{"bar", "b", "a", "r"}, false},
		{[]string{"bar", "b", x, "r"}, true},
		{[]string{"bar", "b", "a", "r", x}, true},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("list len(%d) exp %v", len(c.have), c.exp), func(t *testing.T) {
			xt.Eq(t, c.exp, HasString(c.have, x))
		})
	}
}

func TestRemoveStrings(t *testing.T) {
	t.Run("unordered slice", func(t *testing.T) {
		a := []string{"orange", "apple", "pear", "grape", "banana"}
		exp := []string{"orange", "apple", "pear", "banana"}

		xt.Eq(t, exp, RemoveStrings(a, "grape"))

		// a is not modified
		xt.Eq(t, []string{"orange", "apple", "pear", "grape", "banana"}, a)
	})

	t.Run("ordered slice", func(t *testing.T) {
		a := []string{"orange", "apple", "pear", "grape", "banana"}
		sort.Strings(a)
		exp := []string{"orange", "apple", "pear", "banana"}
		sort.Strings(exp)

		xt.Eq(t, exp, RemoveStrings(a, "grape"))
	})

	t.Run("remove multiple", func(t *testing.T) {
		a := []string{"orange", "apple", "pear", "grape", "banana"}
		exp := []string{"apple", "pear", "banana"}

		xt.Eq(t, exp, RemoveStrings(a, "grape", "orange"))
	})

	t.Run("ordered slice and remove multiple start", func(t *testing.T) {
		a := []string{"orange", "apple", "pear", "grape", "banana"}
		sort.Strings(a)
		exp := []string{"apple", "pear", "grape"}
		sort.Strings(exp)

		xt.Eq(t, exp, RemoveStrings(a, "banana", "orange"))
	})

	t.Run("nothing removed", func(t *testing.T) {
		a := []string{"orange", "apple", "pear", "grape", "banana"}

		xt.Eq(t, a, RemoveStrings(a, "papaya"))
	})

	t.Run("one string not found", func(t *testing.T) {
		a := []string{"orange", "apple", "pear", "grape", "banana"}
		exp := []string{"orange", "apple", "pear", "grape"}

		xt.Eq(t, exp, RemoveStrings(a, "papaya", "banana"))
	})
}
