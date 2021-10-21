// Copyright (c) 2021, Geert JM Vanderkelen

package xgrid

import "os"

func ExampleNewGrid_rdbms_error() {
	grid := NewGrid(os.Stdout)
	_ = grid.AddHeaderRow("Level", "Code", "Message")
	_ = grid.Append(NewRow("Error", 1049, "Unknown database 'dual'"))
	_ = grid.Render()

	// Output:
	// ┌───────┬──────┬─────────────────────────┐
	// │ Level │ Code │ Message                 │
	// ├───────┼──────┼─────────────────────────┤
	// │ Error │ 1049 │ Unknown database 'dual' │
	// └───────┴──────┴─────────────────────────┘
}
