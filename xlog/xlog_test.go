// Copyright (c) 2021, Geert JM Vanderkelen

package xlog

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestSetGetLevel(t *testing.T) {
	defer func() { defaultLogger.level = defaultLogLevel }()

	xt.Eq(t, defaultLogLevel, defaultLogger.level)
	SetLevel(DebugLevel)
	xt.Eq(t, DebugLevel, GetLevel())
	SetLevel(defaultLogLevel)
	xt.Eq(t, defaultLogLevel, defaultLogger.level)
}

func TestSetGetOut(t *testing.T) {
	defer func() { defaultLogger.Out = os.Stderr }()
	out := &bytes.Buffer{}
	SetOut(out)
	xt.Eq(t, out, GetOut())
}

func TestSetGetFormatter(t *testing.T) {
	defer func() { defaultLogger.Formatter = &TextFormat{} }()
	f := &JSONFormat{}
	SetFormatter(f)
	xt.Eq(t, f, GetFormatter())
}

func TestWithField(t *testing.T) {
	defer func() { defaultLogger.Out = os.Stderr }()
	out := &bytes.Buffer{}
	SetOut(out)

	exp := "I am value"
	entry := WithField("someField", exp)
	xt.Eq(t, exp, entry.Fields["someField"])
	entry.Info("test TestWithField")

	expNeedles := []string{
		`msg="test TestWithField"`,
		`someField="I am value"`,
	}

	got := out.String()
	for _, exp := range expNeedles {
		xt.Assert(t, strings.Contains(got, exp))
	}
}

func TestWithFields(t *testing.T) {
	defer func() { defaultLogger.Out = os.Stderr }()
	out := &bytes.Buffer{}
	SetOut(out)

	fields := Fields{
		"field1": "value1",
		"field2": 1234,
	}

	entry := WithFields(fields)
	xt.Eq(t, fields["field1"], entry.Fields["field1"])
	xt.Eq(t, fields["field2"], entry.Fields["field2"])
	entry.Info("test TestWithFields")

	expNeedles := []string{
		`msg="test TestWithFields"`,
		`field1="value1"`,
		`field2=1234`,
	}

	got := out.String()
	for _, exp := range expNeedles {
		xt.Assert(t, strings.Contains(got, exp))
	}
}

func TestWithError(t *testing.T) {
	defer func() { defaultLogger.Out = os.Stderr }()
	out := &bytes.Buffer{}
	SetOut(out)

	entry := WithError(fmt.Errorf("I am error"))
	entry.Info("level does not matter")

	expNeedles := []string{
		`msg="level does not matter"`,
		`err="I am error"`,
	}

	got := out.String()
	for _, exp := range expNeedles {
		xt.Assert(t, strings.Contains(got, exp))
	}
}

func TestWithScope(t *testing.T) {
	defer func() { defaultLogger.Out = os.Stderr }()
	out := &bytes.Buffer{}
	SetOut(out)

	entry := WithScope("access")
	entry.Info("level does not matter")

	expNeedles := []string{
		`msg="level does not matter"`,
		`scope="access"`,
	}

	got := out.String()
	for _, exp := range expNeedles {
		xt.Assert(t, strings.Contains(got, exp))
	}
}

func TestPanic(t *testing.T) {
	t.Run("panics are not logged; but panic anyway", func(t *testing.T) {
		defer func() {
			defaultLogger.Out = os.Stderr
			defaultLogger.SetLevel(defaultLogLevel)
		}()

		out := &bytes.Buffer{}
		SetOut(out)
		xt.Panics(t, func() {
			Panic("Don't panic.")
		})
		got := out.String()
		xt.Assert(t, !strings.Contains(got, "Don't panic."))
	})

	t.Run("panics are logged when asked", func(t *testing.T) {
		defer func() {
			defaultLogger.Out = os.Stderr
			defaultLogger.SetLevel(defaultLogLevel)
		}()
		defaultLogger.SetLevel(PanicLevel)

		out := &bytes.Buffer{}
		SetOut(out)
		xt.Panics(t, func() {
			Panic("Don't panic.")
		})
		got := out.String()
		xt.Assert(t, strings.Contains(got, "Don't panic."))
		xt.Match(t, `.* fileInfo=".*/xlog/xlog_test.go:\d+".*`, got)
	})
}

func TestPanicf(t *testing.T) {
	t.Run("panics are not logged; but panic anyway", func(t *testing.T) {
		defer func() {
			defaultLogger.Out = os.Stderr
			defaultLogger.SetLevel(defaultLogLevel)
		}()

		out := &bytes.Buffer{}
		SetOut(out)
		xt.Panics(t, func() {
			Panicf("Don't %s.", "panic")
		})
		got := out.String()
		xt.Assert(t, !strings.Contains(got, "Don't panic."))
	})

	t.Run("panics are logged when asked", func(t *testing.T) {
		defer func() {
			defaultLogger.Out = os.Stderr
			defaultLogger.SetLevel(defaultLogLevel)
		}()
		defaultLogger.SetLevel(PanicLevel)

		out := &bytes.Buffer{}
		SetOut(out)
		xt.Panics(t, func() {
			Panicf("Don't %s.", "panic")
		})
		got := out.String()
		xt.Assert(t, strings.Contains(got, "Don't panic."))
	})
}

func TestError(t *testing.T) {
	defer func() { defaultLogger.Out = os.Stderr }()
	out := &bytes.Buffer{}
	SetOut(out)

	Error("I am error")
	expNeedles := []string{
		`msg="I am error"`,
		`level=error`,
	}

	got := out.String()
	for _, exp := range expNeedles {
		xt.Assert(t, strings.Contains(got, exp), fmt.Sprintf("got: %s", got))
	}
}

func TestErrorf(t *testing.T) {
	defer func() { defaultLogger.Out = os.Stderr }()
	out := &bytes.Buffer{}
	SetOut(out)

	Errorf("I am %s", "error")
	expNeedles := []string{
		`msg="I am error"`,
		`level=error`,
	}

	got := out.String()
	for _, exp := range expNeedles {
		xt.Assert(t, strings.Contains(got, exp), fmt.Sprintf("got: %s", got))
	}
}

func TestWarn(t *testing.T) {
	defer func() { defaultLogger.Out = os.Stderr }()
	out := &bytes.Buffer{}
	SetOut(out)

	Warn("I am warning")
	expNeedles := []string{
		`msg="I am warning"`,
		`level=warn`,
	}

	got := out.String()
	for _, exp := range expNeedles {
		xt.Assert(t, strings.Contains(got, exp), fmt.Sprintf("got: %s", got))
	}
}

func TestWarnf(t *testing.T) {
	defer func() { defaultLogger.Out = os.Stderr }()
	out := &bytes.Buffer{}
	SetOut(out)

	Warnf("I am %s", "warning")
	expNeedles := []string{
		`msg="I am warning"`,
		`level=warn`,
	}

	got := out.String()
	for _, exp := range expNeedles {
		xt.Assert(t, strings.Contains(got, exp), fmt.Sprintf("got: %s", got))
	}
}

func TestDebug(t *testing.T) {
	defer func() { defaultLogger.Out = os.Stderr }()
	out := &bytes.Buffer{}
	SetOut(out)

	Debug("I am debug")
	expNeedles := []string{
		`msg="I am debug"`,
		`level=debug`,
	}

	got := out.String()
	for _, exp := range expNeedles {
		xt.Assert(t, strings.Contains(got, exp), fmt.Sprintf("got: %s", got))
	}
}

func TestDebugf(t *testing.T) {
	defer func() { defaultLogger.Out = os.Stderr }()
	out := &bytes.Buffer{}
	SetOut(out)

	Debugf("I am %s", "debug")
	expNeedles := []string{
		`msg="I am debug"`,
		`level=debug`,
	}

	got := out.String()
	for _, exp := range expNeedles {
		xt.Assert(t, strings.Contains(got, exp), fmt.Sprintf("got: %s", got))
	}
}

func TestInfo(t *testing.T) {
	defer func() { defaultLogger.Out = os.Stderr }()
	out := &bytes.Buffer{}
	SetOut(out)

	Info("I am info")
	expNeedles := []string{
		`msg="I am info"`,
		`level=info`,
	}

	got := out.String()
	for _, exp := range expNeedles {
		xt.Assert(t, strings.Contains(got, exp), fmt.Sprintf("got: %s", got))
	}
}

func TestInfof(t *testing.T) {
	defer func() { defaultLogger.Out = os.Stderr }()
	out := &bytes.Buffer{}
	SetOut(out)

	Infof("I am %s", "info")
	expNeedles := []string{
		`msg="I am info"`,
		`level=info`,
	}

	got := out.String()
	for _, exp := range expNeedles {
		xt.Assert(t, strings.Contains(got, exp), fmt.Sprintf("got: %s", got))
	}
}
