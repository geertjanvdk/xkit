// Copyright (c) 2021, Geert JM Vanderkelen

package xgrid

// Row defines a row in a grid. A row consists of 1 or more cells which contain
// the data. A row can also be marked as a header, which is also on top of the
// grid and can be differently styled.
type Row struct {
	isHeader bool
	cells    []interface{}
}

// NewRow instantiate a row.
func NewRow(cells ...interface{}) Row {
	return Row{cells: cells}
}

// Append appends cells to the row. Cells contain the data.
func (r *Row) Append(cells ...interface{}) {
	r.cells = append(r.cells, cells...)
}

// Len returns the number of cells in the row.
func (r *Row) Len() int {
	return len(r.cells)
}
