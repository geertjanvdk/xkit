// Copyright (c) 2021, Geert JM Vanderkelen

package xprompt

import (
	"bufio"
	"io"
	"strings"
)

func readString(r io.Reader) (string, error) {
	s, err := bufio.NewReader(r).ReadString('\n')
	return strings.TrimSpace(s), err
}
