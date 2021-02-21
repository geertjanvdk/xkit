// Copyright (c) 2019, 2021, Geert JM Vanderkelen

package xlog

import (
	"encoding/json"
	"time"
)

// Fields is used to add extra, custom fields to a log entry.
type Fields map[string]interface{}

func (f Fields) MarshalJSON() ([]byte, error) {
	res := map[string]interface{}{}

	for k, value := range f {
		switch v := value.(type) {
		case time.Duration:
			// Go does not implement marshalling of time.Duration
			res[k] = v.String()
		default:
			res[k] = value
		}
	}

	return json.Marshal(res)
}
