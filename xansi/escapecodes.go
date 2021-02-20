// Copyright (c) 2021, Geert JM Vanderkelen

package xansi

import (
	"math"
	"strconv"
	"strings"

	"github.com/geertjanvdk/xkit/xutil"
)

type SGR struct {
	attribute uint8
	variant   uint8
	arguments []uint8
}

func (sgr *SGR) String() string {
	parts := []string{strconv.FormatUint(uint64(sgr.attribute), 10)}

	if sgr.variant > 0 {
		parts = append(parts, strconv.FormatUint(uint64(sgr.variant), 10))
	}

	for _, arg := range sgr.arguments {
		parts = append(parts, strconv.FormatUint(uint64(arg), 10))
	}

	return strings.Join(parts, ";")
}

const esc = "\u001b["

var (
	AllReset  = SGR{attribute: 0}
	Bold      = SGR{attribute: 1}
	Underline = SGR{attribute: 4}
	FgColor   = SGR{attribute: 38}
	BgColor   = SGR{attribute: 48}

	// 3/4-bit colors
	Black           = SGR{attribute: 30}
	Red             = SGR{attribute: 31}
	Green           = SGR{attribute: 32}
	Yellow          = SGR{attribute: 33}
	Blue            = SGR{attribute: 34}
	Magenta         = SGR{attribute: 35}
	Cyan            = SGR{attribute: 36}
	White           = SGR{attribute: 37}
	BrightBlack     = SGR{attribute: 90} // Gray
	BrightRed       = SGR{attribute: 91}
	BrightGreen     = SGR{attribute: 92}
	BrightYellow    = SGR{attribute: 93}
	BrightBlue      = SGR{attribute: 94}
	BrightMagenta   = SGR{attribute: 95}
	BrightCyan      = SGR{attribute: 96}
	BrightWhite     = SGR{attribute: 97}
	BgBlack         = SGR{attribute: 40}
	BgRed           = SGR{attribute: 41}
	BgGreen         = SGR{attribute: 42}
	BgYellow        = SGR{attribute: 43}
	BgBlue          = SGR{attribute: 44}
	BgMagenta       = SGR{attribute: 45}
	BgCyan          = SGR{attribute: 46}
	BgWhite         = SGR{attribute: 47}
	BgBrightBlack   = SGR{attribute: 100} // Gray
	BgBrightRed     = SGR{attribute: 101}
	BgBrightGreen   = SGR{attribute: 102}
	BgBrightYellow  = SGR{attribute: 103}
	BgBrightBlue    = SGR{attribute: 104}
	BgBrightMagenta = SGR{attribute: 105}
	BgBrightCyan    = SGR{attribute: 106}
	BgBrightWhite   = SGR{attribute: 107}
)

// mustRGB panics when r, g, b are not within range 0..255.
func mustRGB(r, g, b int) {
	if !xutil.IntInRange(r, 0, math.MaxUint8) {
		panic("(xansi) r must be in range 0 and 255")
	}

	if !xutil.IntInRange(g, 0, math.MaxUint8) {
		panic("(xansi) g must be in range 0 and 255")
	}

	if !xutil.IntInRange(b, 0, math.MaxUint8) {
		panic("(xansi) g must be in range 0 and 255")
	}
}

// FgColor8 returns the SGR foreground color according to the 256-color
// (8-bit) palette.
// Panics when n is not within range 0..255.
func FgColor8(n int) SGR {
	if !xutil.IntInRange(n, 0, math.MaxUint8) {
		panic("(xansi) n must be in range 0 and 255")
	}

	return SGR{
		attribute: FgColor.attribute,
		variant:   5,
		arguments: []uint8{uint8(n)},
	}
}

// FgColor24 returns the SGR foreground color using the
// RGB values (24-bit).
// Panics when r, g, b are not within range 0..255.
func FgColor24(r, g, b int) SGR {
	mustRGB(r, g, b)

	return SGR{
		attribute: FgColor.attribute,
		variant:   2,
		arguments: []uint8{uint8(r), uint8(g), uint8(b)},
	}
}

// BgColor8 returns the SGR background color according to the 256-color
// (8-bit) palette.
// Panics when n is not not within range 0..255.
func BgColor8(n int) SGR {
	if !xutil.IntInRange(n, 0, math.MaxUint8) {
		panic("(xansi) n must be in range 0 and 255")
	}

	return SGR{
		attribute: BgColor.attribute,
		variant:   5,
		arguments: []uint8{uint8(n)},
	}
}

// BgColor24 returns the SGR background color using the
// RGB values (24-bit).
// Panics when r, g, b are not within range 0..255.
func BgColor24(r, g, b int) SGR {
	mustRGB(r, g, b)

	return SGR{
		attribute: BgColor.attribute,
		variant:   2,
		arguments: []uint8{uint8(r), uint8(g), uint8(b)},
	}
}

// Reset resets any SGR.
func Reset() string {
	return esc + AllReset.String() + "m"
}
