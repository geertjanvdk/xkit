// Copyright (c) 2020, Geert JM Vanderkelen

package xsql

import "database/sql"

// DeferCloseRows closes rows. This is used with the `defer` statement.
// The error is conveniently ignored.
func DeferCloseRows(rows *sql.Rows) {
	if rows != nil {
		_ = rows.Close()
	}
}

// DeferRollback rolls back tx. This is used with the `defer` statement.
// The error is conveniently ignored.
func DeferRollback(tx *sql.Tx) {
	if tx != nil {
		_ = tx.Rollback()
	}
}
