// Copyright (c) 2021, Geert JM Vanderkelen

package xutil

import (
	"time"
)

// UNow returns universal current timestamp (location set to UTC).
func UNow() time.Time {
	return time.Now().UTC()
}

// UNowPtr returns universal current timestamp (location set to UTC)
// as pointer value.
func UNowPtr() *time.Time {
	ts := time.Now().UTC()
	return &ts
}
