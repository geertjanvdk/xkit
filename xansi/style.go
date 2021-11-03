// Copyright (c) 2021, Geert JM Vanderkelen

package xansi

import (
	"fmt"
	"strings"
)

type Render []SGR

func (s Render) join() string {
	if len(s) == 0 {
		return ""
	}

	parts := make([]string, len(s))
	for i, sgr := range s {
		parts[i] = sgr.String()
	}

	return strings.Join(parts, ";") + "m"
}

func (s Render) Sprintf(format string, a ...interface{}) string {
	return esc + "[" + s.join() + fmt.Sprintf(format, a...)
}
