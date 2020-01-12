// Copyright (c) 2020, Geert JM Vanderkelen

package xpath

import "os"

// IsDir returns whether path is a directory.
func IsDir(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		return fi.Mode().IsDir()
	}
	return false
}