// Copyright (c) 2020, Geert JM Vanderkelen

package xt

import (
	"reflect"
	"strings"
	"testing"
)

var TestHelper = (*testing.T).Helper

func fatal(t *testing.T, err interface{}, messages ...string) {
	TestHelper(t)

	fm := "test failed"

	switch e := err.(type) {
	case error:
		fm = e.Error()
	case string:
		fm = e
	case *string:
		fm = *e
	default:
		panic("invalid type for err")
	}

	var m string
	if len(messages) > 0 {
		m = "\n\n" + strings.Join(messages, "\n") + "\n"
	}

	t.Fatal(fm + m)
}

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}

	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return rv.IsNil()
	}

	return false
}
