// Copyright (c) 2021, Geert JM Vanderkelen

package xutil

import (
	"reflect"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestInt64Ptr(t *testing.T) {
	rv := reflect.ValueOf(Int64Ptr(1234))
	xt.Assert(t, rv.Kind() == reflect.Ptr)
}
