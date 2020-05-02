// Copyright (c) 2020, Geert JM Vanderkelen

package xnet

import (
	"fmt"
	"strings"
	"testing"

	"lab.scrum.pub/go/ts"
)

func TestIsEmailAddress(t *testing.T) {
	t.Run("valid email addresses", func(t *testing.T) {
		cases := []string{
			"john@example.com",
			"with-hyphen@example.com",
			"with:colon@example.com",
			"with#hash@example.com",
			"with~tilde@example.com",
			"with_underscore@example.com",
			"with+plus@example.com",
			"with.dot@example.com",

			"with-hyphen-@example.com",
			"with:colon:@example.com",
			"with#hash#@example.com",
			"with~tilde~@example.com",
			"with_underscore_@example.com",
			"with+plus+@example.com",
		}

		for _, c := range cases {
			t.Run("valid_"+c, func(t *testing.T) {
				ts.Assert(t, IsEmailAddress(c), fmt.Sprintf("expected %s to be valid", c))
			})
		}
	})

	t.Run("invalid email addresses", func(t *testing.T) {
		cases := []string{
			"john.example.com",
			"localpart.no.ending.dot.@example.com",
			"no.double..dots@example.com",
			"no.multi.....dots@example.com",
			".no.dot.at.beginning@example.com",
			"max64.in.local.part." + strings.Repeat("a", 44+1) + "@example.com",
			"domain.part.is.max255@" + strings.Repeat("a", 255) + ".com",
		}

		for _, sign := range "!$%&'*/=?^`{|}" {
			cases = append(cases, "with"+string(sign)+"sign@example.com")
		}

		for _, c := range cases {
			t.Run("invalid_"+c, func(t *testing.T) {
				ts.Assert(t, !IsEmailAddress(c), fmt.Sprintf("expected %s to be invalid", c))
			})
		}
	})
}
