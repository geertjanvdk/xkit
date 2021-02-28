// Copyright (c) 2021, Geert JM Vanderkelen

package xutil

import (
	"reflect"
	"testing"
	"time"

	"github.com/geertjanvdk/xkit/xt"
)

func TestUNow(t *testing.T) {
	t.Run("timezone set to UTC", func(t *testing.T) {
		now := time.Now().UTC()
		ts := UNow()
		xt.Eq(t, time.UTC, ts.Location())
		xt.Assert(t, ts.After(now) || ts.Equal(now))
		xt.Assert(t, ts.Sub(now).Seconds() < 1) // if running on a 80486, maybe set to 2
	})

	t.Run("return pointer value", func(t *testing.T) {
		now := time.Now().UTC()
		ts := UNowPtr()
		xt.Eq(t, reflect.Ptr, reflect.ValueOf(ts).Kind())
		xt.Eq(t, time.UTC, ts.Location())
		xt.Assert(t, ts.After(now) || ts.Equal(now))
		xt.Assert(t, ts.Sub(now).Seconds() < 1) // if running on a 80486, maybe set to 2
	})
}
