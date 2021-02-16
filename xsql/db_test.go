// Copyright (c) 2020, Geert JM Vanderkelen

package xsql

import (
	"fmt"
	"strings"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestOpen(t *testing.T) {
	t.Run("unsupported driver", func(t *testing.T) {
		xt.Panics(t, func() {
			_, _ = Open("mongodb", "does not matter")
		})
	})
}

func TestOpen_MySQL(t *testing.T) {
	var driver = "mysql"

	t.Run("open bad DSN", func(t *testing.T) {
		dsn := "user_does_not_exists:@tcp(127.0.0.1:3306)/xkit_xsql"
		_, err := Open("mysql", dsn)
		xt.KO(t, err)
	})

	t.Run("successfully connect but with error", func(t *testing.T) {
		_, err := Open("mysql", testDSNMySQL)
		xt.KO(t, err, fmt.Sprintf("using DSN %s", testDSNMySQL))
		xt.Assert(t, strings.Contains(err.Error(), "Error 1049"))
	})

	t.Run("successfully connect", func(t *testing.T) {
		dbName := "xkit_23SKdkWll"

		// create the database
		dsn, err := SetDataSourceNameDatabase(driver, testDSNMySQL, "")
		xt.OK(t, err)

		db, err := Open(driver, dsn)
		xt.OK(t, err, fmt.Sprintf("using DSN %s", testDSNMySQL),
			"user must be able to create databases")

		_, err = db.Exec("CREATE DATABASE " + dbName)
		xt.OK(t, err)

		// connect using database
		dsn, err = SetDataSourceNameDatabase(driver, testDSNMySQL, dbName)
		xt.OK(t, err)

		db, err = Open(driver, dsn)
		xt.OK(t, err, fmt.Sprintf("using DSN %s", testDSNMySQL),
			"user must be able to create databases")

		// clean up
		_, err = db.Exec("DROP DATABASE " + dbName)
		xt.OK(t, err)
	})
}

func TestDB_HaveConstraint_MySQL(t *testing.T) {
	var driver = "mysql"
	dbName := "xkit_TestDB_HaveConstraint_MySQL"
	var db *DB

	{
		// create the database
		dsn, err := SetDataSourceNameDatabase(driver, testDSNMySQL, "")
		xt.OK(t, err)

		db, err = Open(driver, dsn)
		xt.OK(t, err, fmt.Sprintf("using DSN %s", testDSNMySQL),
			"user must be able to create databases")

		_, err = db.Exec("DROP DATABASE IF EXISTS " + dbName)
		xt.OK(t, err)

		_, err = db.Exec("CREATE DATABASE " + dbName)
		xt.OK(t, err)

		// connect using database
		dsn, err = SetDataSourceNameDatabase(driver, testDSNMySQL, dbName)
		xt.OK(t, err)

		db, err = Open(driver, dsn)
		xt.OK(t, err, fmt.Sprintf("using DSN %s", testDSNMySQL),
			"user must be able to create databases")
	}

	t.Run("table has constraint or not", func(t *testing.T) {
		tableName := "t001"
		constraint := tableName + "_c1_nonzero"
		q := fmt.Sprintf("CREATE TABLE %s (c1 INT, CONSTRAINT %s CHECK (c1 <> 0))",
			tableName, constraint)

		_, err := db.Exec(q)
		xt.OK(t, err)

		have, err := db.HaveConstraint(tableName, constraint)
		xt.OK(t, err)
		xt.Assert(t, have)

		have, err = db.HaveConstraint(tableName, "not_existing_constraint")
		xt.OK(t, err)
		xt.Assert(t, !have)
	})
}

func TestDB_HaveTrigger_MySQL(t *testing.T) {
	var driver = "mysql"
	dbName := "xkit_TestDB_HaveTrigger_MySQL"
	var db *DB

	{
		// create the database
		dsn, err := SetDataSourceNameDatabase(driver, testDSNMySQL, "")
		xt.OK(t, err)

		db, err = Open(driver, dsn)
		xt.OK(t, err, fmt.Sprintf("using DSN %s", testDSNMySQL),
			"user must be able to create databases")

		_, err = db.Exec("DROP DATABASE IF EXISTS " + dbName)
		xt.OK(t, err)

		_, err = db.Exec("CREATE DATABASE " + dbName)
		xt.OK(t, err)

		// connect using database
		dsn, err = SetDataSourceNameDatabase(driver, testDSNMySQL, dbName)
		xt.OK(t, err)

		db, err = Open(driver, dsn)
		xt.OK(t, err, fmt.Sprintf("using DSN %s", testDSNMySQL),
			"user must be able to create databases")
	}

	t.Run("table has trigger or not", func(t *testing.T) {
		tableName := "t002"
		trigger := tableName + "_c1_trigger"
		q := fmt.Sprintf("CREATE TABLE %s (c1 INT)", tableName)

		_, err := db.Exec(q)
		xt.OK(t, err)

		q = fmt.Sprintf("CREATE TRIGGER %s BEFORE INSERT ON %s FOR EACH ROW SET @c1 = @c1 * 2",
			trigger, tableName)
		_, err = db.Exec(q)
		xt.OK(t, err)

		have, err := db.HaveTrigger(tableName, trigger)
		xt.OK(t, err)
		xt.Assert(t, have)

		have, err = db.HaveTrigger(tableName, "not_existing_trigger")
		xt.OK(t, err)
		xt.Assert(t, !have)
	})
}
