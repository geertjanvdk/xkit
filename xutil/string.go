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
	return IndexString(a, x) > -1
}

// IndexString returns the position of x in slice a or -1
// when x is not part of a.
// Note that this is different than Go's `sort.SearchString` which
// requires the slice to be sorted.
func IndexString(a []string, x string) int {
	for i, e := range a {
		if e == x {
			return i
		}
	}

	return -1
}

// RemoveStrings removes first occurrence of x from a and returns
// the result leaving a unmodified.
func RemoveString(a []string, x string) []string {
	index := IndexString(a, x)
	if index == -1 {
		return a
	}

	var n []string
	n = append(n, a[:index]...)
	return append(n, a[index+1:]...)
}

// RemoveStrings removes first occurrence of any x from a and returns
// the result leaving a unmodified.
func RemoveStrings(a []string, x ...string) []string {
	result := append([]string{}, a...)
	for _, needle := range x {
		result = RemoveString(result, needle)
	}

	return result
}
