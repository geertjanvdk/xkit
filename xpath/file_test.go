// Copyright (c) 2020, Geert JM Vanderkelen

package xpath

import (
	"testing"

	"lab.scrum.pub/go/ts"
)

func TestIsRegularFile(t *testing.T) {
	t.Run("existing regular file", func(t *testing.T) {
		ts.Assert(t, IsRegularFile("file.go"))
	})

	t.Run("non-existing regular file", func(t *testing.T) {
		ts.Assert(t, !IsRegularFile("filefilefile.go"))
	})

	t.Run("dir is not a regular file", func(t *testing.T) {
		ts.Assert(t, !IsRegularFile("../xpath"))
	})
}
