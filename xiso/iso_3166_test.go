// Copyright (c) 2021, Geert JM Vanderkelen

package xiso

import (
	"testing"
	"time"

	"github.com/geertjanvdk/xkit/xt"
)

func TestISO3166Update(t *testing.T) {
	xt.Assert(t, iso3166LastUpdate.Add(time.Hour*24*90).After(time.Now().UTC()),
		"check ISO 3166 for updates")
}

func TestCountryAlpha2(t *testing.T) {
	t.Run("valid country Alpha-2 codes", func(t *testing.T) {
		cases := []string{"BE", "PL", "DE", "io", "Gb"}
		for _, c := range cases {
			t.Run(c, func(t *testing.T) {
				xt.Assert(t, !CountryAlpha2(c).IsEmpty())
			})
		}
	})

	t.Run("invalid country Alpha-2 codes", func(t *testing.T) {
		cases := []string{"xt", "01", "foo", "B", "foobar"}
		for _, c := range cases {
			t.Run(c, func(t *testing.T) {
				xt.Assert(t, CountryAlpha2(c).IsEmpty())
			})
		}
	})
}

func TestCountryAlpha3(t *testing.T) {
	t.Run("valid country Alpha-3 codes", func(t *testing.T) {
		cases := []string{"BEL", "POL", "DEU", "iot", "GbR"}
		for _, c := range cases {
			t.Run(c, func(t *testing.T) {
				xt.Assert(t, !CountryAlpha3(c).IsEmpty())
			})
		}
	})

	t.Run("invalid country Alpha-3 codes", func(t *testing.T) {
		cases := []string{"xtt", "010", "01", "foo", "B", "foobar"}
		for _, c := range cases {
			t.Run(c, func(t *testing.T) {
				xt.Assert(t, CountryAlpha3(c).IsEmpty())
			})
		}
	})
}
