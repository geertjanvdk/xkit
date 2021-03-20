// Copyright (c) 2019, 2021 Geert JM Vanderkelen

package xlog

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"
)

const (
	FieldError   = "err"
	FieldErrCode = "errCode"
	FieldTime    = "time"
	FieldLevel   = "level"
	FieldMsg     = "msg"
	FieldScope   = "scope"
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
	Message string
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

func (e Entry) getLogTime() time.Time {
	if e.logger.UseUTC {
		return time.Now().UTC()
	}

	return time.Now()
}

func (e *Entry) log(level Level) {
	if e.Time.IsZero() {
		e.Time = e.getLogTime()
	}

	e.Level = level
	if e.Scope == "" {
		e.Scope = e.logger.Scope // which can be empty
	}
	e.write()

	if level == PanicLevel {
		panic(e.Message)
	}
}

// String returns the textual representation serialized by the logger's formatter.
func (e Entry) String() string {
	te := []byte(e.Message)

	if e.logger != nil {
		var err error
		te, err = e.logger.Formatter.Format(e)
		if err != nil {
			return "(failed formatting log entry)"
		}
	}

	return string(te)
}

func (e Entry) write() {
	e.logger.mu.Lock()
	defer e.logger.mu.Unlock()

	te, err := e.logger.Formatter.Format(e)
	if err != nil {
		te = []byte(fmt.Sprintf("failed formatting log entry: %v\n", e))
	}

	_, err = e.logger.Out.Write([]byte(te))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed writing to log: %s\n", err)
	}
}

func (e *Entry) Log(level Level, a ...interface{}) {
	if !e.logger.isLogged(level) {
		return
	}

	if len(a) > 0 {
		if err, ok := a[0].(error); ok {
			m := reMySQLError.FindAllStringSubmatch(err.Error(), -1)
			if m != nil && len(m[0]) == 3 {
				e.Message = m[0][2]
				e.ErrCode = m[0][1]
				e.Scope = "mysql"
			} else {
				e.Message = err.Error()
			}
		} else {
			e.Message = fmt.Sprint(a...)
		}
	} else {
		e.Message = fmt.Sprint(a...)
	}

	e.log(level)
}

func (e *Entry) Logf(level Level, format string, a ...interface{}) {
	if !e.logger.isLogged(level) {
		return
	}

	e.Log(level, fmt.Sprintf(format, a...))
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
	e.Fields[FieldError] = err.Error()

	return e
}

func (e *Entry) Error(a ...interface{}) {
	e.Log(ErrorLevel, a...)
}

func (e *Entry) Errorf(format string, a ...interface{}) {
	e.Logf(ErrorLevel, format, a...)
}

func (e *Entry) Warn(a ...interface{}) {
	e.Log(WarnLevel, a...)
}

func (e *Entry) Warnf(format string, a ...interface{}) {
	e.Logf(WarnLevel, format, a...)
}

func (e *Entry) Info(a ...interface{}) {
	e.Log(InfoLevel, a...)
}

func (e *Entry) Infof(format string, a ...interface{}) {
	e.Logf(InfoLevel, format, a...)
}

func (e *Entry) Debug(a ...interface{}) {
	e.Log(DebugLevel, a...)
}

func (e *Entry) Debugf(format string, a ...interface{}) {
	e.Logf(DebugLevel, format, a...)
}

func (e *Entry) UnmarshalJSON(data []byte) error {
	var res struct {
		Fields Fields
		baseEntry
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	e.Message = res.Message
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
