// Copyright (c) 2019, 2021, Geert JM Vanderkelen

package xlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Logger struct {
	mu        sync.Mutex
	level     Level
	Out       io.Writer
	Formatter Formatter
	Scope     string
	UseUTC    bool
}

func New() *Logger {
	return &Logger{
		level:     defaultLogLevel,
		Formatter: &TextFormat{},
		Out:       os.Stderr,
		UseUTC:    true,
	}
}

func (l *Logger) NewEntry() *Entry {
	return newEntry(l)
}

// Level returns log level of l.
func (l Logger) Level() Level {
	return l.level
}

// Level returns log level of l.
func (l *Logger) SetLevel(level Level) Level {
	if _, ok := levelName[level]; !ok {
		panic(fmt.Sprintf("in"))
	}
	l.level = level
	return l.level
}

// SetFormatter sets f as formatter for l.
func (l *Logger) SetFormatter(f Formatter) {
	l.Formatter = f
}

// isLogged returns true if level has to be logged.
func (l Logger) isLogged(level Level) bool {
	return l.level >= level
}

func (l *Logger) WithField(name string, value interface{}) *Entry {
	return newEntry(l).WithField(name, value)
}

func (l *Logger) WithFields(fields Fields) *Entry {
	return newEntry(l).WithFields(fields)
}

// Logf logs according to a format specifier, and optional arguments, for given level.
func (l *Logger) Logf(level Level, format string, a ...interface{}) {
	if l.isLogged(level) {
		entry := l.NewEntry()
		entry.Logf(level, format, a...)
	}
}

// Logf logs using optional arguments, for given level.
func (l *Logger) Log(level Level, a ...interface{}) {
	if l.isLogged(level) {
		entry := l.NewEntry()
		entry.Log(level, a...)
	}
}

func (l *Logger) Panic(a ...interface{}) {
	l.Log(PanicLevel, a...)
}

func (l *Logger) Panicf(format string, a ...interface{}) {
	l.Logf(PanicLevel, format, a...)
}

func (l *Logger) Error(a ...interface{}) {
	l.Log(ErrorLevel, a...)
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Logf(ErrorLevel, format, a...)
}

func (l *Logger) Warn(a ...interface{}) {
	l.Log(WarnLevel, a...)
}

func (l *Logger) Warnf(format string, a ...interface{}) {
	l.Logf(WarnLevel, format, a...)
}

func (l *Logger) Info(a ...interface{}) {
	l.Log(InfoLevel, a...)
}

func (l *Logger) Infof(format string, a ...interface{}) {
	l.Logf(InfoLevel, format, a...)
}

func (l *Logger) Debug(a ...interface{}) {
	l.Log(DebugLevel, a...)
}

func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Logf(DebugLevel, format, a...)
}

func (l *Logger) Print(v ...interface{}) {
	log.Print(v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
