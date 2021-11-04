// Copyright (c) 2019, 2021, Geert JM Vanderkelen

package xlog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestLogger_Logf(t *testing.T) {
	t.Run("non-panic levels", func(t *testing.T) {
		out := &bytes.Buffer{}

		l := New()
		l.Out = out
		l.level = DebugLevel

		for _, level := range []Level{ErrorLevel, WarnLevel, InfoLevel, DebugLevel} {
			ln := levelName[level]
			t.Run("level "+ln, func(t *testing.T) {
				expFormat := "%s log entry"
				l.Logf(level, expFormat, ln)

				bres, err := ioutil.ReadAll(out)
				xt.OK(t, err)
				res := string(bres)
				xt.Assert(t, strings.Contains(res, `msg="`+fmt.Sprintf(expFormat, ln)+`"`))
				xt.Assert(t, strings.Contains(res, "level="+ln))
				xt.Assert(t, strings.Contains(res, `time=`), "was: ", res)
			})
		}
	})

	t.Run("panic level", func(t *testing.T) {
		out := &bytes.Buffer{}

		l := New()
		l.Out = out

		xt.Panics(t, func() {
			l.Logf(PanicLevel, "this is a panic entry log")
		})
		bres, err := ioutil.ReadAll(out)
		xt.OK(t, err)

		res := string(bres)
		xt.Assert(t, strings.Contains(res, `msg="this is a panic entry log"`))
		xt.Assert(t, strings.Contains(res, "level=panic"))
		xt.Assert(t, strings.Contains(res, `time=`))
	})
}

func TestLogger_WithError(t *testing.T) {
	t.Run("no error field when err is nil", func(t *testing.T) {
		out := &bytes.Buffer{}

		l := New()
		l.Out = out

		l.WithError(nil).Info("no error field")

		exp := `^time=.*?\s{1}level=info msg="no error field"$`
		got := strings.TrimSpace(out.String())
		xt.Match(t, exp, got)
	})

	t.Run("with error field when err is not nil", func(t *testing.T) {
		out := &bytes.Buffer{}

		l := New()
		l.Out = out

		l.WithError(fmt.Errorf("this is an error")).Info("no error field")

		exp := `^time=.*?\s{1}level=info msg="no error field" err="this is an error"$`
		got := strings.TrimSpace(out.String())
		fmt.Println(got)
		xt.Match(t, exp, got)
	})
}
