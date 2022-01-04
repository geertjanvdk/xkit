// Copyright (c) 2022, Geert JM Vanderkelen

package xutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/geertjanvdk/xkit/xt"
)

func TestRetry(t *testing.T) {
	t.Run("retry 3 times", func(t *testing.T) {
		counter := 3
		attempts := counter
		action := func() error {
			counter -= 1
			if counter == 0 {
				return nil
			}
			return fmt.Errorf("not yet")
		}

		xt.OK(t, Retry(attempts, time.Second/2, action))
	})

	t.Run("fail after retry 3 times", func(t *testing.T) {
		counter := 3
		attempts := counter
		exp := fmt.Errorf("fail anyway")

		action := func() error {
			counter -= 1
			if counter == 0 {
				return exp
			}
			return fmt.Errorf("not yet")
		}

		err := Retry(attempts, time.Second/2, action)
		xt.KO(t, err)
		xt.Eq(t, exp, err)
	})

	t.Run("panics when attempts < 1", func(t *testing.T) {
		xt.Panics(t, func() {
			_ = Retry(0, time.Second/2, func() error { return nil })
		})
	})

	t.Run("panics when interval is negative", func(t *testing.T) {
		xt.Panics(t, func() {
			_ = Retry(1, -1, func() error { return nil })
		})
	})
}
