// Copyright (c) 2021, Geert JM Vanderkelen

package xgrid

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/geertjanvdk/xkit/xansi"
)

// Grid defines a grid which can present data in tabular form.
type Grid struct {
	writeTo   io.Writer
	rows      []Row
	headers   []Row
	cellWidth []int

	headerRender *xansi.Render
}

// NewGrid instantiates a Grid with data stream w, to which the grid will be
// rendered to.
func NewGrid(w io.Writer) *Grid {
	if w == nil {
		panic("w is required")
	}
	return &Grid{
		writeTo:      w,
		headerRender: nil,
	}
}

// SetHeaderStyle sets the render r to style the header. For example,
// r as `xansi.Render{xansi.Bold}` would show headers bold.
func (g *Grid) SetHeaderStyle(r *xansi.Render) {
	g.headerRender = r
}

func (g *Grid) addRow(dst *[]Row, row Row) error {
	if row.Len() == 0 {
		// nothing to do
		return nil
	}

	if len(g.cellWidth) == 0 {
		g.cellWidth = make([]int, row.Len())
	}

	if row.Len() != len(g.cellWidth) {
		return fmt.Errorf("row length must be %d", len(g.cellWidth))
	}

	for i, col := range row.cells {
		v := fmt.Sprintf("%v", col)
		if g.cellWidth[i] < len(v) {
			g.cellWidth[i] = len(v)
		}
	}
	*dst = append(*dst, row)
	return nil
}

// AddHeaderRow adds a new row as a header using data in heads.
// It is thus possible to have multiple header rows.
func (g *Grid) AddHeaderRow(heads ...interface{}) error {
	row := NewRow(heads...)
	row.isHeader = true
	return g.addRow(&g.headers, row)
}

// Append appends a row at the end of the grid.
func (g *Grid) Append(rows ...Row) error {
	for _, row := range rows {
		if err := g.addRow(&g.rows, row); err != nil {
			return err
		}
	}

	return nil
}

func (g *Grid) top(w io.Writer) {
	MustFprint(w, "┌")
	for i, width := range g.cellWidth {
		MustFprint(w, strings.Repeat("─", width+2))
		if i == len(g.cellWidth)-1 {
			MustFprint(w, "┐\n")
		} else {
			MustFprint(w, "┬")
		}
	}
}

func (g *Grid) bottom(w io.Writer) {
	MustFprint(w, "└")
	for i, width := range g.cellWidth {
		MustFprint(w, strings.Repeat("─", width+2))
		if i == len(g.cellWidth)-1 {
			MustFprint(w, "┘\n")
		} else {
			MustFprint(w, "┴")
		}
	}
}

func (g *Grid) separator(w io.Writer) {
	MustFprint(w, "├")
	for i, width := range g.cellWidth {
		MustFprint(w, strings.Repeat("─", width+2))
		if i == len(g.cellWidth)-1 {
			MustFprint(w, "┤\n")
		} else {
			MustFprint(w, "┼")
		}
	}
}

func (g *Grid) row(w io.Writer, row Row, render *xansi.Render) error {
	for j, col := range row.cells {
		if j == 0 {
			MustFprint(w, "│")
		}

		content := fmt.Sprintf("%-*v", g.cellWidth[j], col)
		if render != nil {
			content = render.Sprint(content)
		}
		MustFprint(w, " "+content+" ")

		if j == len(row.cells)-1 {
			MustFprint(w, "│\n")
		} else {
			MustFprint(w, "│")
		}
	}
	return nil
}

func (g *Grid) render(w io.Writer) error {
	if len(g.headers) > 0 {
		for i := -1; i < len(g.headers)+1; i++ {
			if i == -1 {
				g.top(w)
				continue
			}

			if i == len(g.headers) {
				if len(g.rows) > 0 {
					g.separator(w)
				} else {
					g.bottom(w)
				}
				continue
			}

			if err := g.row(w, g.headers[i], g.headerRender); err != nil {
				return err
			}
		}
	}

	if len(g.rows) > 0 {
		for i := -1; i < len(g.rows)+1; i++ {
			if i == -1 {
				if len(g.headers) == 0 {
					g.top(w)
				}
				continue
			}

			if i == len(g.rows) {
				g.bottom(w)
				continue
			}

			if err := g.row(w, g.rows[i], nil); err != nil {
				return err
			}
		}
	}

	return nil
}

// Render renders the grid to the data stream stored with the grid.
func (g *Grid) Render() error {
	return g.render(g.writeTo)
}

// RenderToString renders the grid and returns it as string.
func (g *Grid) RenderToString() string {
	w := new(bytes.Buffer)
	if err := g.render(w); err != nil {
		return "(grid) " + err.Error()
	}

	buf, err := ioutil.ReadAll(w)
	if err != nil {
		return "(grid) " + err.Error()
	}

	return string(buf)
}
