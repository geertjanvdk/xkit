// Copyright (c) 2021, Geert JM Vanderkelen

package xutil

// IntInRange returns true when n within range lower..upper.
func IntInRange(n, lower, upper int) bool {
	return n >= lower && n <= upper
}

// Int64Ptr returns i as pointer.
func Int64Ptr(i int64) *int64 {
	return &i
}
