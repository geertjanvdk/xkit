// Copyright (c) 2020, Geert JM Vanderkelen

package xpath

import "os"

// IsRegularFile returns whether path is a regular file.
func IsRegularFile(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		return fi.Mode().IsRegular()
	}
	return false
}
