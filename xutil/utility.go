// Copyright (c) 2022, Geert JM Vanderkelen

package xutil

import "time"

// Retry will rerun action until it does not return an error any longer. It
// attempts a number of times, with an interval. When unsuccessful, the last
// error produced by action will be returned.
// Panics when attempts is less than 1, or interval is negative.
func Retry(attempts int, interval time.Duration, action func() error) error {
	if attempts < 1 {
		panic("attempts must be 1 or greater")
	}

	if interval < 0 {
		panic("interval must not be negative")
	}

	var err error
	for i := 0; i < attempts; i++ {
		if err = action(); err == nil {
			// all good
			return nil
		}
		time.Sleep(interval)
	}

	return err
}
