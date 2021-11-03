// Copyright (c) 2021, Geert JM Vanderkelen

package xansi

import (
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestRender_Sprintf(t *testing.T) {
	t.Run("bold text", func(t *testing.T) {
		res := Render{Bold}.Sprintf("I am Bold!")
		xt.Eq(t, esc+"[1mI am Bold!", res)
	})

	t.Run("bold text", func(t *testing.T) {
		res := Render{Bold}.Sprintf("I am Bold!")
		xt.Eq(t, esc+"[1mI am Bold!", res)
	})

	t.Run("red bold text", func(t *testing.T) {
		res := Render{Red, Bold}.Sprintf("I am Red & Bold!")
		xt.Eq(t, esc+"[31;1mI am Red & Bold!", res)
	})

	t.Run("bold red text", func(t *testing.T) {
		res := Render{Bold, Red}.Sprintf("I am Red & Bold!")
		xt.Eq(t, esc+"[1;31mI am Red & Bold!", res)
	})

	t.Run("underlined bold white on bright cyan background", func(t *testing.T) {
		txt := "I am underlined bold white on bright cyan!"
		res := Render{Bold, Yellow, Underline, BgBrightCyan}.Sprintf(txt)
		xt.Eq(t, esc+"[1;33;4;106m"+txt, res)
	})

	t.Run("red on yellow background using 8-bit colors", func(t *testing.T) {
		txt := "I am red on yellow text! "
		res := Render{FgColor8(196), BgColor8(226)}.Sprintf(txt)
		xt.Eq(t, esc+"[38;5;196;48;5;226m"+txt, res)
	})

	t.Run("24-bit (RGB) colors", func(t *testing.T) {
		txt := "I am RGB pink on cyan text!"
		res := Render{FgColor24(255, 0, 222), BgColor24(23, 253, 255)}.Sprintf(txt)
		xt.Eq(t, esc+"[38;2;255;0;222;48;2;23;253;255m"+txt, res)
	})
}
