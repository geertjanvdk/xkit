// Copyright (c) 2022, Geert JM Vanderkelen

package xos

import (
	"os"
	"reflect"
	"strings"
)

// EnvFromStruct uses the tagged fields of struct s to set the field's
// value to what is available from the environment. The tags used are
// `envVar` which sets the name of the environment variable and `default`
// which sets an optional default value.
//
// Only environment variables which are not empty (string) are considered,
// except for variables which have a field typed as 'bool'.
//
// Bool-typed fields are true when
// 1. there is no value (empty string) and the variable is available
// 2. or when string is not empty and is either 'true' or 't' (case insensitive),
//    or '1'.
//
// Panics when s is not a pointer value of struct.
func EnvFromStruct(s any) error {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Pointer {
		panic("s must be struct and pointer value")
	}

	rt := t.Elem()
	if rt.Kind() != reflect.Struct {
		panic("s must be struct and pointer value")
	}

	rv := reflect.ValueOf(s).Elem()

	for i := 0; i < rt.NumField(); i++ {
		rtf := rt.Field(i)

		envVar := rtf.Tag.Get("envVar")
		if envVar == "" {
			continue
		}

		d := rtf.Tag.Get("default")
		var envDefault interface{}
		if d == "" {
			envDefault = nil
		} else {
			envDefault = d
		}

		envValue, have := os.LookupEnv(envVar)

		switch rtf.Type.Kind() {
		case reflect.Bool:
			v := have
			s := strings.TrimSpace(strings.ToLower(envValue))
			if s != "" {
				v = s == "true" || s == "t" || s == "1"
			}
			rv.Field(i).SetBool(v)
		default:
			if envValue == "" && envDefault != nil {
				rv.Field(i).SetString(envDefault.(string))
			} else if envValue != "" {
				rv.Field(i).SetString(envValue)
			}
		}
	}

	return nil
}
