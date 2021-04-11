// Copyright (c) 2021, Geert JM Vanderkelen

package xutil

import (
	"encoding/json"
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

func TestUDate(t *testing.T) {
	t.Run("2020-02-28 10:11:12.000000123 +0000 UTC", func(t *testing.T) {
		exp := "2020-02-28 10:11:12.000000123 +0000 UTC"
		ts := UDate(2020, 2, 28, 10, 11, 12, 123)
		xt.Eq(t, time.UTC, ts.Location())
		xt.Eq(t, exp, ts.String())
	})

	t.Run("pointer value 2020-02-28 10:11:12.000000123 +0000 UTC", func(t *testing.T) {
		exp := "2020-02-28 10:11:12.000000123 +0000 UTC"
		ts := UDatePtr(2020, 2, 28, 10, 11, 12, 123)
		xt.Eq(t, reflect.Ptr, reflect.ValueOf(ts).Kind())
		xt.Eq(t, time.UTC, ts.Location())
		xt.Eq(t, exp, ts.String())
	})
}

func TestFractionalTime_MarshalJSON(t *testing.T) {
	berlin, err := time.LoadLocation("Europe/Berlin")
	xt.OK(t, err)

	cases := []struct {
		value FractionalTime
		exp   string
	}{
		{
			value: FractionalTime{time.Date(2019, 4, 2, 4, 23, 44, 0, time.UTC)},
			exp:   `"2019-04-02T04:23:44.0Z"`,
		},
		{
			value: FractionalTime{time.Date(2019, 4, 2, 4, 23, 44, 199292235, time.UTC)},
			exp:   `"2019-04-02T04:23:44.199292Z"`,
		},
		{
			value: FractionalTime{time.Date(2019, 4, 2, 4, 23, 44, 299297685, berlin)},
			exp:   `"2019-04-02T04:23:44.299297+02:00"`,
		},
		{
			value: FractionalTime{time.Date(2019, 4, 2, 4, 23, 44, 0, berlin)},
			exp:   `"2019-04-02T04:23:44.0+02:00"`,
		},
	}

	for _, c := range cases {
		t.Run(c.exp, func(t *testing.T) {
			res, err := json.Marshal(c.value)
			xt.OK(t, err)
			xt.Eq(t, c.exp, string(res))
		})
	}
}
