// Copyright (c) 2019, 2021, Geert JM Vanderkelen

package xlog

import "os"

type Level int

// A classic way of describing the logging level. Note that 'Panic' is the
// Go-way for 'Trace'.
const (
	FatalLevel   Level = -1 // exits
	DefaultLevel Level = 0  // use whatever is xlog's default
	ErrorLevel   Level = 1
	WarnLevel    Level = 2
	InfoLevel    Level = 3
	DebugLevel   Level = 4
	PanicLevel   Level = 5
)

const (
	lowestLevel  = FatalLevel
	highestLevel = PanicLevel
)

var levelName = map[Level]string{
	FatalLevel: "fatal",
	PanicLevel: "panic",
	ErrorLevel: "error",
	WarnLevel:  "warn",
	InfoLevel:  "info",
	DebugLevel: "debug",
}

var defaultActiveLevels = activeLevels{
	FatalLevel: true,
	ErrorLevel: true,
	WarnLevel:  true,
	InfoLevel:  true,
	DebugLevel: false,
	PanicLevel: false,
}

// defaultLogger is always available when using xlog.
var defaultLogger = &Logger{
	Formatter:    &TextFormat{},
	Out:          os.Stderr,
	UseUTC:       true,
	activeLevels: newActiveLevels(),
}
