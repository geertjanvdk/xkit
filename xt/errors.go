// Copyright (c) 2020, Geert JM Vanderkelen

package xt

import (
	"fmt"
	"testing"
)

func OK(t *testing.T, err error, messages ...string) {
	TestHelper(t)

	if err != nil {
		if len(messages) > 0 {
			messages = append([]string{"--"}, messages...)
		}
		fatal(t, fmt.Sprintf("\u001B[31;1mexpected no error, got:\u001B[0m\n%s", err.Error()),
			messages...)
	}
}

func KO(t *testing.T, err error, messages ...string) {
	TestHelper(t)

	if err == nil {
		if len(messages) > 0 {
			messages = append([]string{"--"}, messages...)
		}
		fatal(t, fmt.Sprintf("\u001B[31;1mexpected error\u001B[0m"), messages...)
	}
}
