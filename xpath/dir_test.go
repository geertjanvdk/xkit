// Copyright (c) 2020, Geert JM Vanderkelen

package xpath

import (
	"io/ioutil"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestIsDir(t *testing.T) {
	t.Run("existing directory", func(t *testing.T) {
		xt.Assert(t, IsDir("../xpath"))
	})

	t.Run("non-existing directory", func(t *testing.T) {
		xt.Assert(t, !IsDir("../xpathxpathxpath"))
	})

	t.Run("regular file is not a dir", func(t *testing.T) {
		xt.Assert(t, !IsDir("dir.go"))
	})
}

func TestRegularFilesInDir(t *testing.T) {
	t.Run("list xpath files", func(t *testing.T) {
		dir := "testdata/regular_files_in_dir"
		readdir, err := ioutil.ReadDir(dir)
		xt.OK(t, err)
		var exp []string
		for _, f := range readdir {
			exp = append(exp, f.Name())
		}

		files, err := RegularFilesInDir(dir)
		xt.OK(t, err)

		xt.Eq(t, exp, files)
	})

	t.Run("list files in xkit", func(t *testing.T) {
		files, err := RegularFilesInDir("../xt")
		xt.OK(t, err)

		exp := []string{"assert.go", "equality.go", "errors.go", "panic.go", "regex.go", "utils.go"}
		xt.Eq(t, exp, files)
	})
}
