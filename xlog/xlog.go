// Copyright (c) 2019, 2021, Geert JM Vanderkelen

package xlog

import (
	"io"
)

type Level int

const (
	FatalLevel Level = -1 // exits
	PanicLevel Level = 1
	ErrorLevel Level = 2
	WarnLevel  Level = 3
	InfoLevel  Level = 4
	DebugLevel Level = 5
)

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

// GetLevel returns the level of the default logger.
func GetLevel() Level {
	return defaultLogger.Level()
}

// SetOut sets where the output of the default logger goes to.
func SetOut(w io.Writer) {
	defaultLogger.Out = w
}

// GetOut returns where the default logger sends its output.
func GetOut() io.Writer {
	return defaultLogger.Out
}

// GetFormatter returns the formatter of the default logger.
func GetFormatter() Formatter {
	return defaultLogger.Formatter
}

// SetFormatter sets the formatter of the default logger.
func SetFormatter(f Formatter) {
	defaultLogger.SetFormatter(f)
}

// WithField returns an Entry for the default logger
// which has field set using name and value.
func WithField(name string, value interface{}) *Entry {
	return newEntry(defaultLogger).WithField(name, value)
}

// WithFields returns an Entry for the default logger
// which has all fields set.
func WithFields(fields Fields) *Entry {
	return newEntry(defaultLogger).WithFields(fields)
}

// WithError returns an Entry for the default logger
// which has a field set with value of err. The name of the
// field is defined as xlog.FieldError.
func WithError(err error) *Entry {
	return newEntry(defaultLogger).WithError(err)
}

// WithScope returns an Entry for the default logger
// which has a field set with value of scope. The name of the
// field is defined as xlog.FieldScope.
func WithScope(scope string) *Entry {
	return newEntry(defaultLogger).WithScope(scope)
}

// Panic simply panics and formats the message using provided operands.
func Panic(a ...interface{}) {
	defaultLogger.log(3, PanicLevel, a...)
}

// Panicf simply panics and formats the message according to a
// format specifier and operands.
func Panicf(format string, a ...interface{}) {
	defaultLogger.logf(3, PanicLevel, format, a...)
}

// Error logs an error entry using the default logger formatting using provided operands.
func Error(a ...interface{}) {
	defaultLogger.Error(a...)
}

// Errorf logs an error entry using the default logger formatting according to a
// format specifier and operands.
func Errorf(format string, a ...interface{}) {
	defaultLogger.Errorf(format, a...)
}

// Warn logs a warning entry using the default logger formatting using provided operands.
func Warn(a ...interface{}) {
	defaultLogger.Warn(a...)
}

// Warnf logs a warning entry using the default logger formatting according to a
// format specifier and operands.
func Warnf(format string, a ...interface{}) {
	defaultLogger.Warnf(format, a...)
}

// Info logs a informational entry using the default logger formatting using
// provided operands.
func Info(a ...interface{}) {
	defaultLogger.Info(a...)
}

// Infof logs a informational entry using the default logger formatting according to a
// format specifier and operands.
func Infof(format string, a ...interface{}) {
	defaultLogger.Infof(format, a...)
}

// Debug logs a debug entry using the default logger formatting using
// provided operands.
func Debug(a ...interface{}) {
	defaultLogger.Debug(a...)
}

// Debugf logs a debug entry using the default logger formatting according to a
// format specifier and operands.
func Debugf(format string, a ...interface{}) {
	defaultLogger.Debugf(format, a...)
}
