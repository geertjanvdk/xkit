// Copyright (c) 2021, Geert JM Vanderkelen

package xmath

// AbsInt returns the absolute value of integer x.
// Note that when x is math.MinInt32 (or -2147483648), the same
// negative value is returned since there is no 'signed' int 2147483648.
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// AbsInt64 returns the absolute value of 64-bit integer x.
// Note that when x is math.MinInt (or -9223372036854775808), the same
// negative value is returned since there is no 'signed' int64 9223372036854775808.
func AbsInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
