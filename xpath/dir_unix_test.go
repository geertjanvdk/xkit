// Copyright (c) 2020, Geert JM Vanderkelen

package xpath

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"lab.scrum.pub/go/ts"
)

func TestFilesInDir(t *testing.T) {
	t.Run("all files with symbolic link", func(t *testing.T) {
		dir := "testdata/files_in_dir"

		script := dir + "/create_sym_link.sh"
		ts.OK(t, os.Chmod(script, 0700))
		cmd := exec.Command(dir + "/create_sym_link.sh")
		ts.OK(t, cmd.Run())

		readdir, err := ioutil.ReadDir(dir)
		ts.OK(t, err)
		var exp []string
		for _, f := range readdir {
			exp = append(exp, f.Name())
		}

		files, err := FilesInDir(dir)
		ts.OK(t, err)

		ts.Eq(t, exp, files)
	})
}
