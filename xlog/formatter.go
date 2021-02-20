// Copyright (c) 2019, 2021 Geert JM Vanderkelen

package xlog

type Formatter interface {
	Format(e Entry) ([]byte, error)
}
