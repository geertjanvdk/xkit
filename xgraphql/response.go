// Copyright (c) 2020, Geert JM Vanderkelen

package xgraphql

import (
	"encoding/json"
	"fmt"
	"time"
)

// Response represents a GraphQL response.
type Response struct {
	Errors []Error         `json:"errors,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

// Error represents a GraphQL error.
type Error struct {
	Message    string     `json:"message"`
	Extensions *Extension `json:"extensions,omitempty"`
}

// Error returns the error a a canonical string.
func (e Error) Error() string {
	return e.Message
}

// Extension represent the extension entry for GraphQL errors.
type Extension struct {
	Code      int       `json:"code,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface. It encodes
// extension ex making it ready to include in a GraphQL response.
func (ex Extension) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}

	m["code"] = ex.Code
	m["codeHex"] = fmt.Sprintf("0x%04x", ex.Code)
	m["timestamp"] = ex.Timestamp
	return json.Marshal(m)
}

// CodeAsInt tries to return the code of ex as int. If the code is empty
// or conversation is not possible, zero (0) is returned.
func (ex Extension) CodeAsInt() int {
	return ex.Code
}

// CodeAsInt64 tries to return the code of ex as int64. If the code is empty
// or conversation is not possible, zero (0) is returned.
func (ex Extension) CodeAsInt64() int64 {
	return int64(ex.Code)
}
