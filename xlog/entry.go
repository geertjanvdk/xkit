// Copyright (c) 2019, 2021 Geert JM Vanderkelen

package xlog

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

const (
	FieldError    = "err"
	FieldErrCode  = "errCode"
	FieldTime     = "time"
	FieldLevel    = "level"
	FieldMsg      = "msg"
	FieldScope    = "scope"
	FieldFileLine = "fileInfo"
	FieldStack    = "debugStack"
)

var reservedFields = map[string]bool{
	FieldErrCode: true,
	FieldTime:    true,
	FieldLevel:   true,
	FieldMsg:     true,
	FieldScope:   true,
}

var reMySQLError = regexp.MustCompile(`^Error (\d{4}): (.*)$`)

type baseEntry struct {
	Level   Level
	message string
	Time    time.Time
	ErrCode string
	Scope   string
}

type Entry struct {
	logger *Logger
	Fields Fields

	baseEntry
}

func newEntry(logger *Logger) *Entry {
	return &Entry{
		logger: logger,
		Fields: make(Fields),
	}
}

func (e *Entry) getLogTime() time.Time {
	if e.logger.UseUTC {
		return time.Now().UTC()
	}

	return time.Now()
}

// String returns the textual representation serialized by the logger's formatter.
func (e *Entry) String() string {
	te := []byte(e.message)

	if e.logger != nil {
		var err error
		te, err = e.logger.Formatter.Format(e)
		if err != nil {
			return "(failed formatting log entry)"
		}
	}

	return string(te)
}

func (e *Entry) setMessage(a ...interface{}) {
	if len(a) > 0 {
		e.message = fmt.Sprint(a...)
	}
}

func (e *Entry) setMessagef(format string, a ...interface{}) {
	if format != "" {
		e.message = fmt.Sprintf(format, a...)
	}
}

func (e *Entry) WithField(name string, value interface{}) *Entry {
	e.Fields[name] = value

	return e
}

func (e *Entry) WithFields(fields Fields) *Entry {
	for k, v := range fields {
		e.Fields[k] = v
	}

	return e
}

func (e *Entry) WithScope(scope string) *Entry {
	e.Scope = scope

	return e
}

func (e *Entry) WithError(err error) *Entry {
	if err == nil {
		return e
	}

	m := reMySQLError.FindAllStringSubmatch(err.Error(), -1)
	if m != nil && len(m[0]) == 3 {
		e.message = m[0][2]
		e.ErrCode = m[0][1]
		e.Scope = "mysql"
	} else {
		e.message = err.Error()
	}
	e.WithField(FieldError, e.message)

	return e
}

func (e *Entry) output(level Level) {
	e.Level = level
	e.logger.output(4, e)
}

func (e *Entry) Error(a ...interface{}) {
	e.setMessage(a...)
	e.output(ErrorLevel)
}

func (e *Entry) Errorf(format string, a ...interface{}) {
	e.setMessagef(format, a...)
	e.output(ErrorLevel)
}

func (e *Entry) Warn(a ...interface{}) {
	e.setMessage(a...)
	e.output(WarnLevel)
}

func (e *Entry) Warnf(format string, a ...interface{}) {
	e.setMessagef(format, a...)
	e.output(WarnLevel)
}

func (e *Entry) Info(a ...interface{}) {
	e.setMessage(a...)
	e.output(InfoLevel)
}

func (e *Entry) Infof(format string, a ...interface{}) {
	e.setMessagef(format, a...)
	e.output(InfoLevel)
}

func (e *Entry) Debug(a ...interface{}) {
	e.setMessage(a...)
	e.output(DebugLevel)
}

func (e *Entry) Debugf(format string, a ...interface{}) {
	e.setMessagef(format, a...)
	e.output(DebugLevel)
}

func (e *Entry) Panic(a ...interface{}) {
	e.setMessage(a...)
	e.output(PanicLevel)
}

func (e *Entry) Panicf(format string, a ...interface{}) {
	e.setMessagef(format, a...)
	e.output(PanicLevel)
}

func (e *Entry) Fatal(a ...interface{}) {
	e.setMessage(a...)
	e.output(FatalLevel)
}

func (e *Entry) Fatalf(format string, a ...interface{}) {
	e.setMessagef(format, a...)
	e.output(FatalLevel)
}

func (e *Entry) UnmarshalJSON(data []byte) error {
	var res struct {
		Fields Fields
		baseEntry
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	e.message = res.message
	e.Time = res.Time
	e.Level = res.Level
	e.Scope = res.Scope

	e.Fields = make(Fields)
	for k, v := range res.Fields {
		switch vv := v.(type) {
		case string:
			if vvv, err := time.Parse(time.RFC3339Nano, vv); err == nil {
				e.Fields[k] = vvv
			} else if vvv, err := time.Parse(time.RFC3339, vv); err == nil {
				e.Fields[k] = vvv
			} else {
				e.Fields[k] = vv
			}
		default:
			e.Fields[k] = vv
		}
	}

	return nil
}
