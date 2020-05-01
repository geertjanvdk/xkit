// Copyright (c) 2020, Geert JM Vanderkelen

package xutil

// IsZeroString returns true if s is a string or pointer to string and empty (zero) or nil.
// Panics when s is not string or *string or nil.
func IsZeroString(s interface{}) bool {
	switch v := s.(type) {
	case string:
		return v == ""
	case *string:
		if v == nil {
			return true
		}
		return *v == ""
	default:
		panic("argument must be (pointer) string")
	}
}

// StringPtr returns s as pointer.
func StringPtr(s string) *string {
	return &s
}

// HasString returns true whether x is in slice a.
func HasString(a []string, x string) bool {
	l := len(a)
	if l == 0 {
		return false
	}

	if l == 1 {
		return a[0] == x
	}

	// this is O(n) but OK for now
	for _, e := range a {
		if e == x {
			return true
		}
	}

	return false
}
