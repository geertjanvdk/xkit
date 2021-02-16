// Copyright (c) 2021, Geert JM Vanderkelen

package xsql

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func SetDataSourceNameDatabase(driverName, dsn, dbName string) (string, error) {
	switch driverName {
	case "mysql":
		cfg, err := mysql.ParseDSN(dsn)
		if err != nil {
			return "", fmt.Errorf("failed parsing data source name (%s)", err)
		}
		cfg.DBName = dbName
		return cfg.FormatDSN(), nil
	default:
		panic(ErrDriverUnsupported.Error())
	}
}
