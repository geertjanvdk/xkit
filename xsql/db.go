// Copyright (c) 2020, Geert JM Vanderkelen

package xsql

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql" // make sure we compile the MySQL driver in

	"github.com/geertjanvdk/xkit/xutil"
)

var supportedDrivers = []string{"mysql"}

// DB is a wrapper around xsql.DB offering a pool of connections.
type DB struct {
	*sql.DB
	DataSourceName string
	Driver         string
}

// Open opens a database specified by its database driver name and a
// driver-specific data source name.
// Arguments are checked, and the database is pinged.
//
// Panics when driver is not supported.
func Open(driverName, dataSourceName string) (*DB, error) {
	if !xutil.HasString(supportedDrivers, driverName) {
		panic(ErrDriverUnsupported.Error())
	}

	if driverName == "mysql" {
		cfg, err := mysql.ParseDSN(dataSourceName)
		if err != nil {
			return nil, fmt.Errorf("failed parsing data source name (%s)", err)
		}
		if cfg.Params == nil {
			cfg.Params = map[string]string{}
		}
		cfg.Params["parseTime"] = "true"
		if _, have := cfg.Params["timeout"]; !have {
			cfg.Params["timeout"] = "2s"
		}
		dataSourceName = cfg.FormatDSN()
	}

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed connecting to data store (%s)", err)
	}

	return &DB{
		DB:             db,
		DataSourceName: dataSourceName,
		Driver:         driverName,
	}, nil
}

// HaveTable checks whether table is available.
func (db *DB) HaveTable(table string) (bool, error) {
	q := `SELECT COUNT(*) as cnt FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = database() AND TABLE_NAME = ?`

	var cnt int

	err := db.QueryRow(q, table).Scan(&cnt)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, fmt.Errorf("failed checking existance of table `%s` (%s)", table, err)
	}

	return cnt > 0, nil
}

// HaveConstraint checks whether constraint is available for table.
func (db *DB) HaveIndex(table string, index string) (bool, error) {
	q := `SELECT COUNT(*) as cnt FROM INFORMATION_SCHEMA.STATISTICS
WHERE TABLE_SCHEMA = database() AND TABLE_NAME = ? AND INDEX_NAME = ?`

	var cnt int

	err := db.QueryRow(q, table, index).Scan(&cnt)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, fmt.Errorf("failed checking existance of table index `%s(%s)` (%s)",
			table, index, err)
	}

	return cnt > 0, nil
}

// HaveConstraint checks whether constraint is available for table.
func (db *DB) HaveConstraint(table string, constraint string) (bool, error) {
	q := `SELECT COUNT(*) as cnt FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
WHERE CONSTRAINT_SCHEMA = database() AND TABLE_NAME = ? AND CONSTRAINT_NAME = ?`

	var cnt int

	err := db.QueryRow(q, table, constraint).Scan(&cnt)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, fmt.Errorf("failed checking existance of table constraint `%s(%s)` (%s)",
			table, constraint, err)
	}

	return cnt > 0, nil
}

// HaveTrigger checks whether trigger is available for table.
func (db *DB) HaveTrigger(table string, trigger string) (bool, error) {
	q := `SELECT COUNT(*) as cnt FROM INFORMATION_SCHEMA.TRIGGERS
WHERE EVENT_OBJECT_SCHEMA = database() AND EVENT_OBJECT_TABLE = ? AND TRIGGER_NAME = ?`

	var cnt int

	err := db.QueryRow(q, table, trigger).Scan(&cnt)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, fmt.Errorf("failed checking existance of table trigger `%s(%s)` (%s)",
			table, trigger, err)
	}

	return cnt > 0, nil
}

type ColumnInfo struct {
	TableName        string
	Name             string
	DataType         string
	Nullable         bool
	CharMaxLength    *int
	CharOctetLength  *int
	PartOfPrimaryKey bool
}

func (db *DB) ColumnInfo(table, column string) (*ColumnInfo, error) {
	q := `SELECT TABLE_NAME, COLUMN_NAME, DATA_TYPE, IF(IS_NULLABLE='YES', TRUE, FALSE),
       CHARACTER_MAXIMUM_LENGTH, CHARACTER_OCTET_LENGTH,
       IF(COLUMN_KEY='PRI', TRUE, FALSE) AS part_of_pk
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = database() AND TABLE_NAME = ? AND COLUMN_NAME = ?`

	col := &ColumnInfo{}

	err := db.QueryRow(q, table, column).Scan(
		&col.TableName,
		&col.Name,
		&col.DataType,
		&col.Nullable,
		&col.CharMaxLength,
		&col.CharOctetLength,
		&col.PartOfPrimaryKey,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed getting column `%s.%s` (%s)", table, column, err)
	}

	return col, nil
}
