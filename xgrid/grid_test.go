// Copyright (c) 2021, Geert JM Vanderkelen

package xgrid

import (
	"fmt"
	"os"
	"testing"

	"github.com/geertjanvdk/xkit/xansi"
	"github.com/geertjanvdk/xkit/xt"
)

func TestGrid_AddRows(t *testing.T) {
	t.Run("not stylish one", func(t *testing.T) {
		g := NewGrid(os.Stdout)
		xt.OK(t, g.Append(NewRow(1)))
		exp := []rune{9484, 9472, 9472, 9472, 9488, 10, 9474, 32, 49, 32, 9474, 10, 9492, 9472, 9472, 9472, 9496, 10}
		xt.Eq(t, exp, []rune(g.RenderToString()))
	})

	t.Run("3 columns, 2 rows, dynamic", func(t *testing.T) {
		exp := `
┌────┬───────────────────────────────────────────────────────────────────┬────┐
│ 1A │ 1B                                                                │ 1C │
│ 2A │ A grid is like underwear, you wear it but it's not to be exposed. │ 2C │
└────┴───────────────────────────────────────────────────────────────────┴────┘
`

		g := NewGrid(os.Stdout)
		xt.OK(t, g.Append([]Row{
			NewRow("1A", "1B", "1C"),
			// actual quote, look it up!
			NewRow("2A", "A grid is like underwear, you wear it but it's not to be exposed.", "2C"),
		}...))

		xt.Eq(t, exp, "\n"+g.RenderToString()) // extra newline for pretty output
	})
}

func TestGrid_AddHeader(t *testing.T) {
	h1 := boldSprintf("A")
	h2 := boldSprintf("B B")
	h3 := boldSprintf("C C C")

	t.Run("header, now rows, bold", func(t *testing.T) {
		exp := `
┌───┬─────┬───────┐
│ ` + h1 + ` │ ` + h2 + ` │ ` + h3 + ` │
└───┴─────┴───────┘
`

		g := NewGrid(os.Stdout)
		g.SetHeaderStyle(&xansi.Render{xansi.Bold})
		xt.OK(t, g.AddHeaderRow("A", "B B", "C C C"))
		xt.Eq(t, exp, "\n"+g.RenderToString()) // extra newline for pretty output
	})

	t.Run("header and 2 rows", func(t *testing.T) {
		h1 := boldSprintf("A ") // matches size of row data

		exp := `
┌────┬─────┬───────┐
│ ` + h1 + ` │ ` + h2 + ` │ ` + h3 + ` │
├────┼─────┼───────┤
│ 1A │ 1B  │ 1C    │
│ 2A │ 2B  │ 2C    │
└────┴─────┴───────┘
`

		g := NewGrid(os.Stdout)
		g.SetHeaderStyle(&xansi.Render{xansi.Bold})
		xt.OK(t, g.AddHeaderRow("A", "B B", "C C C"))
		xt.OK(t, g.Append(
			NewRow("1A", "1B", "1C"),
			NewRow("2A", "2B", "2C"),
		))

		res := "\n" + g.RenderToString()

		xt.Eq(t, []rune(exp), []rune("\n"+g.RenderToString()))

		xt.Eq(t, exp, res) // extra newline for pretty output
	})

	t.Run("lots of headers", func(t *testing.T) {
		g := NewGrid(os.Stdout)
		var headers []interface{}
		for i := 0; i < 60; i++ {
			headers = append(headers, fmt.Sprintf("hh%02d", i))
		}
		xt.OK(t, g.AddHeaderRow(headers...))
	})

	t.Run("content bigger than header", func(t *testing.T) {
		g := NewGrid(os.Stdout)
		xt.OK(t, g.AddHeaderRow("Shorter"))
		xt.OK(t, g.Append(NewRow("Row value of cell 1 is longer")))

		exp := []rune{9484, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472,
			9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472,
			9488, 10, 9474, 32, 83, 104, 111, 114, 116, 101, 114, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32,
			32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 9474, 10, 9500, 9472, 9472, 9472, 9472, 9472, 9472, 9472,
			9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472,
			9472, 9472, 9472, 9472, 9472, 9472, 9472, 9508, 10, 9474, 32, 82, 111, 119, 32, 118, 97, 108, 117,
			101, 32, 111, 102, 32, 99, 101, 108, 108, 32, 49, 32, 105, 115, 32, 108, 111, 110, 103, 101, 114, 32,
			9474, 10, 9492, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472,
			9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472,
			9496, 10}
		xt.Eq(t, exp, []rune(g.RenderToString()))
	})

	t.Run("multi header rows", func(t *testing.T) {
		g := NewGrid(os.Stdout)
		xt.OK(t, g.AddHeaderRow("header row 1"))
		xt.OK(t, g.AddHeaderRow("header row 2"))
		xt.OK(t, g.Append(NewRow(1)))

		exp := []rune{9484, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9488,
			10, 9474, 32, 104, 101, 97, 100, 101, 114, 32, 114, 111, 119, 32, 49, 32, 9474, 10, 9474, 32, 104, 101,
			97, 100, 101, 114, 32, 114, 111, 119, 32, 50, 32, 9474, 10, 9500, 9472, 9472, 9472, 9472, 9472, 9472,
			9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9508, 10, 9474, 32, 49, 32, 32, 32, 32, 32, 32, 32, 32,
			32, 32, 32, 32, 9474, 10, 9492, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472, 9472,
			9472, 9472, 9496, 10}
		xt.Eq(t, exp, []rune(g.RenderToString()))
	})
}
