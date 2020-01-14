// Copyright (c) 2020, Geert JM Vanderkelen

package xpath

import (
	"io/ioutil"
	"testing"

	ts "lab.scrum.pub/go/ts"
)

func TestIsDir(t *testing.T) {
	t.Run("existing directory", func(t *testing.T) {
		ts.Assert(t, IsDir("../xpath"))
	})

	t.Run("non-existing directory", func(t *testing.T) {
		ts.Assert(t, !IsDir("../xpathxpathxpath"))
	})

	t.Run("regular file is not a dir", func(t *testing.T) {
		ts.Assert(t, !IsDir("dir.go"))
	})
}

func TestRegularFilesInDir(t *testing.T) {
	t.Run("list xpath files", func(t *testing.T) {
		dir := "testdata/regular_files_in_dir"
		readdir, err := ioutil.ReadDir(dir)
		ts.OK(t, err)
		var exp []string
		for _, f := range readdir {
			exp = append(exp, f.Name())
		}

		files, err := RegularFilesInDir(dir)
		ts.OK(t, err)

		ts.Eq(t, exp, files)
	})

	t.Run("list files in xkit", func(t *testing.T) {
		files, err := RegularFilesInDir("..")
		ts.OK(t, err)

		exp := []string{"go.mod", "go.sum"}
		ts.Eq(t, exp, files)
	})
}
