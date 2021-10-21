// Copyright (c) 2021, Geert JM Vanderkelen

/*
Package xgrid helps formatting data tabular data.

For example: the following is a result from an RDBMS, with one 1

	┌───────┬──────┬─────────────────────────┐
	│ Level │ Code │ Message                 │
	├───────┼──────┼─────────────────────────┤
	│ Error │ 1049 │ Unknown database 'dual' │
	└───────┴──────┴─────────────────────────┘

The above is produced by:

	grid := xgrid.NewGrid(os.Stdout)
	grid.AddHeaderRow("Level", "Code", "Message")
	grid.Append(xgrid.Row("Error", 1049, "Unknown database 'dual'"))

Using the xansi packate, you can for example set the style how the header is
rendered. To make it blue and bold:

	grid.SetHeaderStyle(&xansi.Render{xansi.Blue, xansi.Bold})

*/
package xgrid
