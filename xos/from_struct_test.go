// Copyright (c) 2022, Geert JM Vanderkelen

package xos

import (
	"os"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

type testAppEnv struct {
	BoolVar   bool   `envVar:"BOOL_VAR"`
	StringVar string `envVar:"STRING_VAR" default:"myDefault"`
}

func TestEnvFromStruct(t *testing.T) {
	t.Run("read values from environment", func(t *testing.T) {
		xt.OK(t, os.Setenv("BOOL_VAR", "1"))
		xt.OK(t, os.Setenv("STRING_VAR", "String"))

		e := &testAppEnv{}
		xt.OK(t, EnvFromStruct(e))

		xt.Eq(t, true, e.BoolVar)
		xt.Eq(t, "String", e.StringVar)
	})

	t.Run("support defaults", func(t *testing.T) {
		xt.OK(t, os.Unsetenv("STRING_VAR"))

		e := &testAppEnv{}
		xt.OK(t, EnvFromStruct(e))

		xt.Eq(t, "myDefault", e.StringVar)
	})

	t.Run("boolean variables", func(t *testing.T) {
		var cases = []string{"t", "T", "True", "true", ""}
		for _, c := range cases {
			t.Run("string value "+c+" is true", func(t *testing.T) {
				xt.OK(t, os.Setenv("BOOL_VAR", c))
				e := &testAppEnv{}
				xt.OK(t, EnvFromStruct(e))
				xt.Eq(t, true, e.BoolVar)
			})
		}

		cases = []string{"f", "false", "0", "-1", "234"}
		for _, c := range cases {
			t.Run("string value "+c+" is false", func(t *testing.T) {
				xt.OK(t, os.Setenv("BOOL_VAR", c))
				e := &testAppEnv{}
				xt.OK(t, EnvFromStruct(e))
				xt.Eq(t, false, e.BoolVar)
			})
		}
	})

	t.Run("panics when not struct or pointer", func(t *testing.T) {
		xt.Panics(t, func() {
			_ = EnvFromStruct(1)
		})

		xt.Panics(t, func() {
			_ = EnvFromStruct(struct {
				foo string
			}{})
		})
	})
}
