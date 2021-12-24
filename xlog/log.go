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
	"sort"
	"sync"
)

type activeLevels map[Level]bool

func newActiveLevels() activeLevels {
	al := activeLevels{}
	for level, active := range defaultActiveLevels {
		al[level] = active
	}

	return al
}

type levelSort []Level

func (s levelSort) Len() int {
	return len(s)
}
func (s levelSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s levelSort) Less(i, j int) bool {
	return int(s[i]) < int(s[j])
}

type Logger struct {
	mu           sync.Mutex
	Out          io.Writer
	Formatter    Formatter
	Scope        string
	UseUTC       bool
	activeLevels activeLevels
}

func New() *Logger {
	return &Logger{
		Formatter:    defaultLogger.Formatter,
		Out:          os.Stderr,
		UseUTC:       true,
		activeLevels: newActiveLevels(),
	}
}

func (l *Logger) NewEntry() *Entry {
	return newEntry(l)
}

// Levels returns the active levels. The result is sorted.
func (l *Logger) Levels() []Level {
	var res []Level
	for level, active := range l.activeLevels {
		if active {
			res = append(res, level)
		}
	}
	sort.Sort(levelSort(res))
	return res
}

// LevelsAsStrings returns the active levels their names. The result
// is sorted.
func (l *Logger) LevelsAsStrings() []string {
	var res []string
	for _, level := range l.Levels() {
		res = append(res, levelName[level])
	}
	sort.Strings(res)
	return res
}

// ActivateLevels is used to activate particular levels of l.
// For example, ActivateLevels(DebugLevel) can be used  if
// the logger logs errors, but debug message are wanted, without
// info messages.
func (l *Logger) ActivateLevels(levels ...Level) {
	if l.activeLevels == nil {
		l.activeLevels = defaultActiveLevels
	}

	for _, level := range levels {
		if !(level >= lowestLevel && level <= highestLevel || level == DefaultLevel) {
			panic(fmt.Sprintf("xlog: invalid log level; was %d", level))
		}
		l.activeLevels[level] = true
	}
}

// DeactivateLevels is used to deactivate particular levels of l.
// For example, DeactivateLevels(DebugLevel) can be used to deactivate all
// debugging messages.
func (l *Logger) DeactivateLevels(levels ...Level) {
	if l.activeLevels == nil {
		l.activeLevels = defaultActiveLevels
	}

	for _, level := range levels {
		if !(level >= lowestLevel && level <= highestLevel) {
			panic(fmt.Sprintf("xlog: invalid log level; was %d", level))
		}
		l.activeLevels[level] = false
	}
}

// SetFormatter sets f as formatter for l.
func (l *Logger) SetFormatter(f Formatter) {
	l.Formatter = f
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

	writable := l.activeLevels[e.Level]

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
