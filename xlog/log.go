// Copyright (c) 2019, 2021, Geert JM Vanderkelen

package xlog

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/debug"
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
		Formatter: defaultLogger.Formatter,
		Out:       os.Stderr,
		UseUTC:    true,
	}
}

func (l *Logger) NewEntry() *Entry {
	return newEntry(l)
}

// Level returns log level of l.
func (l *Logger) Level() Level {
	return l.level
}

// SetLevel sets the log level of l.
func (l *Logger) SetLevel(level Level) Level {
	if _, ok := levelName[level]; !ok {
		panic(fmt.Sprintf("xlog: invalid log level; was %d", level))
	}
	l.level = level
	return l.level
}

// SetFormatter sets f as formatter for l.
func (l *Logger) SetFormatter(f Formatter) {
	l.Formatter = f
}

// isWritten returns true if level has to be written.
func (l *Logger) isWritten(level Level) bool {
	if l.level == 0 {
		return level >= defaultLogLevel
	}
	return level >= l.level
}

// WithError returns an entry with value of field 'error' set to err.
// This is the same as calling Logger.WithField(FieldError, err).
func (l *Logger) WithError(err error) *Entry {
	return newEntry(l).WithError(err)
}

func (l *Logger) WithField(name string, value interface{}) *Entry {
	return newEntry(l).WithField(name, value)
}

func (l *Logger) WithFields(fields Fields) *Entry {
	return newEntry(l).WithFields(fields)
}

// Logf logs according to a format specifier, and optional arguments, for given level.
func (l *Logger) logf(callDepth int, level Level, format string, a ...interface{}) {
	entry := l.NewEntry()
	entry.Level = level
	entry.setMessagef(format, a...)
	l.output(callDepth, entry)
}

// Log logs then entry according to a level using provided operands.
func (l *Logger) log(callDepth int, level Level, a ...interface{}) {
	entry := l.NewEntry()
	entry.Level = level
	entry.setMessage(a...)
	l.output(callDepth, entry)
}

// Logf logs according to a format specifier, and optional arguments, for given level.
func (l *Logger) Logf(level Level, format string, a ...interface{}) {
	l.logf(2, level, format, a...)
}

// Log logs then entry according to a level using provided operands.
func (l *Logger) Log(level Level, a ...interface{}) {
	l.log(2, level, a...)
}

func (l *Logger) Fatal(a ...interface{}) {
	l.Log(FatalLevel, a...)
}

func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.Logf(FatalLevel, format, a...)
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

func (l *Logger) output(callDepth int, e *Entry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if e.Level == PanicLevel {
		// inspired by Go's log package
		l.mu.Unlock()
		var ok bool
		_, file, line, ok := runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
		e.WithField(FieldFileLine, fmt.Sprintf("%s:%d", file, line))
	}

	if e.Time.IsZero() {
		e.Time = e.getLogTime()
	}

	if e.Scope == "" {
		e.Scope = e.logger.Scope // which can be empty
	}

	writable := l.isWritten(e.Level)

	switch e.Level {
	case PanicLevel:
		if writable {
			defer func() {
				if r := recover(); r == nil {
					lines := bytes.Split(debug.Stack(), []byte("\n"))
					e.WithField(FieldStack, string(bytes.Join(lines, []byte("\\n"))))
					l.write(e)
					panic(e.message)
				}
			}()
		} else {
			panic(e.message)
		}
	case FatalLevel:
		if writable {
			l.write(e)
		}
		os.Exit(1)
	default:
		if writable {
			l.write(e)
		}
	}
}

func (l *Logger) write(e *Entry) {
	te, err := l.Formatter.Format(e)
	if err != nil {
		te = []byte(fmt.Sprintf("failed formatting log entry: %v\n", e))
	}

	_, err = l.Out.Write(te)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed writing to log: %s\n", err)
	}
}
