// Copyright (c) 2021, Geert JM Vanderkelen

package xutil

import (
	"errors"
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

// FractionalTime wraps around time.Time to customize the marshalling
// to JSON to always include a fractional part, even when it is zero.
// Only precision up and till microseconds (0,000001) is considered.
type FractionalTime struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in RFC 3339 format, with sub-second precision added, even
// when zero.
// Code borrowed from Go's time.Time.
func (t FractionalTime) MarshalJSON() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(time.RFC3339Nano)+2)
	b = t.AppendFormat(b, "\"2006-01-02T15:04:05.999999Z07:00")
	if b[20] != '.' { // until 9999, always at this position
		var end []byte
		end = append(end, b[20:]...) // usually 1 byte (Z), but we actually dont know
		b = b[0:20]
		b = append(b, []byte{'.', '0'}...)
		b = append(b, end...)

	}
	b = append(b, '"')
	return b, nil
}
