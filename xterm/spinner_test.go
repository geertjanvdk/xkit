// Copyright (c) 2021, Geert JM Vanderkelen

package xterm

import (
	"bytes"
	"testing"
	"time"

	"github.com/geertjanvdk/xkit/xt"
)

func TestNewSpinner(t *testing.T) {
	t.Run("3 seconds of work", func(t *testing.T) {
		var output bytes.Buffer
		spinner, _ := NewSpinner()
		spinner.Delay = time.Second
		spinner.out = &output
		spinner.MessageBefore = "🍺 "
		spinner.Start()

		spinner.MessageAfter = " get beer from fridge"
		time.Sleep(time.Second)
		spinner.MessageAfter = " enjoy beer"
		time.Sleep(time.Second * 2)

		// we deal with time; so timing can screw up
		got := output.Bytes()
		xt.Assert(t, bytes.Contains(got, []byte{0x20, 0xe2, 0xa3, 0xbe, 0x20}), "expected ⣾")
		xt.Assert(t, bytes.Contains(got, []byte{0x20, 0xe2, 0xa3, 0xb7, 0x20}), "expected ⣷")
		xt.Assert(t, bytes.Contains(got, []byte{0x20, 0xe2, 0xa3, 0xaf, 0x20}), "expected ⣯")
	})
}
