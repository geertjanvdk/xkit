// Copyright (c) 2019, 2021 Geert JM Vanderkelen

package xlog

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/geertjanvdk/xkit/xt"
)

func TestEntry_UnmarshalJSON(t *testing.T) {
	now := time.Now().UTC()
	msg := "this is the message"
	data := fmt.Sprintf(`{"time": "%s", "message": "%s", "fields": {"someTime": "%s", "number": 1234, "valid": true}}`,
		now.Format(time.RFC3339Nano),
		msg,
		now.Format(time.RFC3339Nano))

	e := &Entry{}
	xt.OK(t, json.Unmarshal([]byte(data), e))
	xt.Eq(t, e.Time, now)

	t.Run("field as time.Time", func(t *testing.T) {
		r, ok := e.Fields["someTime"].(time.Time)
		xt.Assert(t, ok, "field is not time.Time")
		xt.Eq(t, now, r)
	})

	t.Run("field as float64", func(t *testing.T) {
		r, ok := e.Fields["number"].(float64)
		xt.Assert(t, ok, "field is not float64")
		xt.Eq(t, 1234, r)
	})

	t.Run("field as bool", func(t *testing.T) {
		r, ok := e.Fields["valid"].(bool)
		xt.Assert(t, ok, "field is not bool")
		xt.Eq(t, true, r)
	})
}
