// Copyright (c) 2020, Geert JM Vanderkelen

package ts

import (
	"fmt"
	"testing"
)

func OK(t *testing.T, err error) {
	TestHelper(t)

	if err != nil {
		fatal(t, fmt.Sprintf("\033[31;1mexpected no error, got:\033[39m %s", err.Error()))
	}
}

func KO(t *testing.T, err error) {
	TestHelper(t)

	if err == nil {
		fatal(t, fmt.Sprintf("\033[31;1mexpected error\033[39m"))
	}
}
