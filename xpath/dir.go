// Copyright (c) 2020, Geert JM Vanderkelen

package xpath

import (
	"io/ioutil"
	"os"
	"sort"
)

// IsDir returns whether path is a directory.
func IsDir(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		return fi.Mode().IsDir()
	}
	return false
}

// RegularFilesInDir returns regular files found in directory path.
func RegularFilesInDir(path string) ([]string, error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var l []string
	for _, f := range entries {
		if !f.Mode().IsRegular() {
			continue
		}
		l = append(l, f.Name())
	}

	sort.Strings(l)
	return l, nil
}
