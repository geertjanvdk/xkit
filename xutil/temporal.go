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

// UDate returns the universal timestamp according to
// year-month-day hour:min:sec.nsec (location always set to UTC).
// This is a wrapper around Go's time.Date.
func UDate(year int, month time.Month, day, hour, min, sec, nsec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, nsec, time.UTC)
}

// UDatePtr returns the universal timestamp according to
// year-month-day hour:min:sec.nsec (location always set to UTC) as pointer value.
// This is a wrapper around Go's time.Date.
func UDatePtr(year int, month time.Month, day, hour, min, sec, nsec int) *time.Time {
	ts := time.Date(year, month, day, hour, min, sec, nsec, time.UTC)
	return &ts
}
