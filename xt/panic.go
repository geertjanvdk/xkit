// Copyright (c) 2020, Geert JM Vanderkelen

package xt

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

func Panics(t *testing.T, f func()) {
	TestHelper(t)

	_, file, lineNr, _ := runtime.Caller(1)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf(fmt.Sprintf("expected panic (in test %s:%d)", filepath.Base(file), lineNr))
		}
	}()

	f()
}
