// Copyright (c) 2019, 2021, Geert JM Vanderkelen

package xlog

import (
	"io"
	"os"
)

type Level int

const (
	PanicLevel Level = 0
	ErrorLevel Level = 2
	WarnLevel  Level = 3
	InfoLevel  Level = 4
	DebugLevel Level = 5
)

var defaultLogger = &Logger{
	level:     defaultLogLevel,
	Formatter: &TextFormat{},
	Out:       os.Stderr,
	UseUTC:    true,
}

var levelName = map[Level]string{
	PanicLevel: "panic",
	ErrorLevel: "error",
	WarnLevel:  "warn",
	InfoLevel:  "info",
	DebugLevel: "debug",
}

// SetLevel sets the level of the default logger.
func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// SetOut sets where the output of the default logger goes to.
func SetOut(w io.Writer) {
	defaultLogger.Out = w
}

// GetOut returns where the default logger sends its output.
func GetOut() io.Writer {
	return defaultLogger.Out
}

// GetLevel returns the level of the default logger.
func GetLevel() Level {
	return defaultLogger.Level()
}

// GetFormatter returns the formatter of the default logger.
func GetFormatter() Formatter {
	return defaultLogger.Formatter
}

// SetFormatter sets the formatter of the default logger.
func SetFormatter(f Formatter) {
	defaultLogger.SetFormatter(f)
}

func WithField(name string, value interface{}) *Entry {
	return newEntry(defaultLogger).WithField(name, value)
}

func WithFields(fields Fields) *Entry {
	return newEntry(defaultLogger).WithFields(fields)
}

func WithError(err error) *Entry {
	return newEntry(defaultLogger).WithError(err)
}

func WithScope(scope string) *Entry {
	return newEntry(defaultLogger).WithScope(scope)
}

func Panic(a ...interface{}) {
	defaultLogger.Panic(a...)
}

func Panicf(format string, a ...interface{}) {
	defaultLogger.Panicf(format, a...)
}

func Error(a ...interface{}) {
	defaultLogger.Error(a...)
}

func Errorf(format string, a ...interface{}) {
	defaultLogger.Errorf(format, a...)
}

func Warn(a ...interface{}) {
	defaultLogger.Warn(a...)
}

func Warnf(format string, a ...interface{}) {
	defaultLogger.Warnf(format, a...)
}

func Info(a ...interface{}) {
	defaultLogger.Info(a...)
}

func Infof(format string, a ...interface{}) {
	defaultLogger.Infof(format, a...)
}

func Debug(a ...interface{}) {
	defaultLogger.Debug(a...)
}

func Debugf(format string, a ...interface{}) {
	defaultLogger.Debugf(format, a...)
}

func Print(v ...interface{}) {
	defaultLogger.Print(v...)
}

func Printf(format string, v ...interface{}) {
	defaultLogger.Printf(format, v...)
}
