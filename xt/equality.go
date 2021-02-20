// Copyright (c) 2020, Geert JM Vanderkelen

package xt

import (
	"fmt"
	"reflect"
	"testing"
)

func Eq(t *testing.T, expect, have interface{}, messages ...string) {
	TestHelper(t)
	diff := fmt.Sprintf("\n\u001b[31;1mexpect:\t\u001b[0m%v\n\u001b[31;1mhave:\t\u001b[0m%v", expect, have)

	// we can not compare nil values
	if isNil(expect) || isNil(have) {
		if !(isNil(expect) && isNil(have)) {
			fatal(t, diff, messages...)
		}
		return
	}

	expVal := reflect.ValueOf(expect)
	haveType := reflect.TypeOf(have)

	if !expVal.Type().ConvertibleTo(haveType) {
		messages = append(messages, fmt.Sprintf("\u001b[31;1mcannot compare %v with %v\u001b[0m",
			expVal.Type(), haveType))
		fatal(t, diff, messages...)
	}

	expInt := expVal.Convert(haveType).Interface()

	if !reflect.DeepEqual(expInt, have) {
		fatal(t, diff, messages...)
	}
}
