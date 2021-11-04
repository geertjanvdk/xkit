// Copyright (c) 2020, Geert JM Vanderkelen

package xt

import (
	"fmt"
	"regexp"
	"testing"
)

// Match tests whether the string s contains any match of the regular
// expression pattern.
func Match(t *testing.T, pattern, s string, messages ...string) {
	TestHelper(t)

	m, err := regexp.MatchString(pattern, s)
	if err != nil {
		panic(err.Error())
	}

	if !m {
		fatal(t, fmt.Sprintf("\n\033[31;1mstring:\033[39m ```%s```\n\033[31;1mmust match pattern:\033[39m %s", s, pattern), messages...)
	}
}
