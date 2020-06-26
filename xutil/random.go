// Copyright (c) 2020, Geert JM Vanderkelen

package xutil

import (
	"crypto/rand"
	mathrand "math/rand"
	"time"
)

// RandomBytes returns a n-length slice containing random bytes.
// Panics when n is not >1.
func RandomBytes(n int) ([]byte, error) {
	if n < 1 {
		panic("n must be >1")
	}
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// RandomAlphaNumeric returns a n-length long string containing random
// alphanumeric characters (both lower and upper cased are possible).
// Panics when n is not >1.
func RandomAlphaNumeric(n int) string {
	if n < 1 {
		panic("n must be >1")
	}

	const abc = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	buf := make([]byte, n)
	mathrand.Seed(time.Now().UnixNano())
	for i := range buf {
		buf[i] = abc[mathrand.Intn(len(abc))]
	}
	return string(buf)
}
