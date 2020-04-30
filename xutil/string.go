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
