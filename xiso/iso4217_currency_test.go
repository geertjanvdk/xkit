// Copyright (c) 2021, Geert JM Vanderkelen

package xiso

import (
	"fmt"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestCurrency(t *testing.T) {
	t.Run("number of currencies", func(t *testing.T) {
		xt.Eq(t, 170, len(iso4217Currency))
	})

	t.Run("currency is always 3 char long", func(t *testing.T) {

		for currency, name := range iso4217Currency {
			xt.Assert(t, reISO4217CurrencyCode.MatchString(currency),
				fmt.Sprintf("bad code for %s; was %s", name, currency))
		}
	})

	t.Run("get unknown currency", func(t *testing.T) {
		_, err := Currency("xxx")
		xt.KO(t, err)
		xt.Eq(t, "unknown currency 'xxx'", err.Error())
	})

	t.Run("invalid code (short)", func(t *testing.T) {
		_, err := Currency("xx")
		xt.KO(t, err)
		xt.Eq(t, "invalid currency 'xx' (must be 3 letters)", err.Error())
	})

	t.Run("invalid code (long)", func(t *testing.T) {
		_, err := Currency("xxyy")
		xt.KO(t, err)
		xt.Eq(t, "invalid currency 'xxyy' (must be 3 letters)", err.Error())
	})

	t.Run("get known currency", func(t *testing.T) {
		currency, err := Currency("eUr")
		xt.OK(t, err)
		xt.Eq(t, "Euro", currency.Name)
		xt.Eq(t, "EUR", currency.Code)
	})
}
