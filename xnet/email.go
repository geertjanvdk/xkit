// Copyright (c) 2020, Geert JM Vanderkelen

package xnet

import (
	"regexp"
	"strings"
)

var reEmailAddr = regexp.MustCompile(`^(?i)[a-z0-9][a-z0-9\-:#~_+.]*[a-z0-9\-:#~_+]@[a-z0-9\-.]+\.[a-z]{2,64}$`)

// IsEmailAddress returns whether addr is a valid email address.
// What is not supported as valid:
//
// * quoted local-parts
//
// * following allowed characters are not supported: ! $ % & ' * / = ? ^  ` { | }
//
func IsEmailAddress(addr string) bool {
	// first pass
	if !reEmailAddr.Match([]byte(addr)) {
		return false
	}

	// second pass: whatever is hard to check with regular expression
	if strings.Contains(addr, "..") {
		return false
	}

	parts := strings.Split(addr, "@") // there are always 2 parts
	if len(parts[0]) > 64 {
		return false
	}

	if len(parts[1]) > 255 {
		return false
	}

	return true
}
