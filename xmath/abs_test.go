// Copyright (c) 2021, Geert JM Vanderkelen

package xmath

import (
	"fmt"
	"math"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestAbsInt(t *testing.T) {
	var cases = []int{0, -1, 1, 2147483647, -2147483648}
	var exp = []int{0, 1, 1, 2147483647, 2147483648}

	for i, c := range cases {
		t.Run(fmt.Sprintf("int(%d)", c), func(t *testing.T) {
			xt.Eq(t, exp[i], AbsInt(c))
		})
	}
}

func TestAbsInt64(t *testing.T) {
	var cases = []int64{0, -1, 1, math.MaxInt, math.MinInt}
	var exp = []int64{0, 1, 1, math.MaxInt, math.MinInt}

	for i, c := range cases {
		t.Run(fmt.Sprintf("int64(%d)", c), func(t *testing.T) {
			xt.Eq(t, exp[i], AbsInt64(c))
		})
	}
}
