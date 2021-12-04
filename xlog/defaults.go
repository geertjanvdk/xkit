// Copyright (c) 2019, 2021, Geert JM Vanderkelen

package xlog

import "os"

// defaultLogLevel defines the default log level.
var defaultLogLevel = ErrorLevel

// defaultLogger is always available when using xlog.
var defaultLogger = &Logger{
	level:     defaultLogLevel,
	Formatter: &TextFormat{},
	Out:       os.Stderr,
	UseUTC:    true,
}
